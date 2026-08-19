# SSH Host Key Verification Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace `ssh.InsecureIgnoreHostKey()` at both SSH dial sites with trust-on-first-use host key pinning, backed by a per-node column and an operator-facing reset action.

**Architecture:** A single `ssh.HostKeyCallback` factory in `pkg/node/hostkey.go` reads and writes the pinned key through a narrow `HostKeyStore` interface. `*repository.NodeRepository` satisfies that interface; both `pkg/node/remote.go` (SFTP + remote exec) and `pkg/hub/ssh.go` (web terminal) use the same callback. Recovery after a legitimate rebuild goes through `POST /api/node/trust/reset` and a button on the Nodes page.

**Tech Stack:** Go 1.25.4, `golang.org/x/crypto/ssh`, sqlx + SQLite, golang-migrate (embedded), CloudWeGo Hertz, Vue 3 + PrimeVue.

**Spec:** `docs/superpowers/specs/2026-08-19-ssh-host-key-verification-design.md`

## Global Constraints

- `go.mod` declares `go 1.25.4`; CI (`.github/workflows/go.yml`) pins Go 1.25.4. Do not raise the language version.
- `make test` must pass at the end of every task that touches Go. It depends on the `dashboard-dist` target, so run `make test`, never a bare `go test ./...`.
- Commit messages follow the repo convention: English, conventional-commits subject ≤70 chars, a body whose bullets pair *what changed* with *why*, and a `Signed-off-by: benshi <807629978@qq.com>` trailer.
- **No environment-variable bypass of host key checking.** The spec rejects it explicitly; do not add one "for testing".
- Store the host key as `authorized_keys` text from `ssh.MarshalAuthorizedKey()`. Display it as `ssh.FingerprintSHA256()`. The host key is public — do not route it through `pkg/crypto`.
- Nodes with `is_local = true` never open an SSH connection; no task should add host key handling to the local path.

## File Structure

| File | Responsibility |
|---|---|
| `migrations/upgrade/000009_add_node_ssh_host_key.up.sql` (create) | Adds the `nodes.ssh_host_key` column |
| `pkg/model/node.go` (modify) | `SSHHostKey` field so `SELECT *` still scans; fingerprint helper; `NodeView` field |
| `pkg/repository/node_repository.go` (modify) | Three host key accessors, including the conditional write that makes pinning race-safe |
| `pkg/repository/node_repository_test.go` (create) | Proves the conditional `UPDATE` refuses to overwrite a pinned key |
| `pkg/node/hostkey.go` (create) | `HostKeyStore` interface + `TOFUHostKeyCallback` — the entire security decision lives here |
| `pkg/node/hostkey_test.go` (create) | Five behaviours from the spec, against a fake store |
| `pkg/node/state.go`, `pkg/node/manager.go`, `pkg/node/remote.go` (modify) | Thread the store to the SFTP/exec dial site |
| `pkg/hub/ssh.go` (modify) | Thread the store to the web-terminal dial site |
| `pkg/handler/node.go` (modify) | `ResetNodeHostKeyHandler` |
| `cmd/serve/main.go` (modify) | Register `POST /api/node/trust/reset` |
| `dashboard/src/api/node.ts`, `dashboard/src/views/Nodes.vue` (modify) | Fingerprint display + reset action |

---

### Task 1: Storage — column, model field, repository accessors

**Files:**
- Create: `migrations/upgrade/000009_add_node_ssh_host_key.up.sql`
- Modify: `pkg/model/node.go`
- Modify: `pkg/repository/node_repository.go`
- Test: `pkg/repository/node_repository_test.go` (create)

**Interfaces:**
- Consumes: nothing.
- Produces: `(*NodeRepository).GetSSHHostKey(id int64) (string, error)`, `(*NodeRepository).SetSSHHostKeyIfEmpty(id int64, key string) error`, `(*NodeRepository).ResetSSHHostKey(id int64) error`. Field `model.Node.SSHHostKey string`.

**Why the repository gets its own test:** the spec's test list covers the callback only, but the callback's concurrency guarantee is delegated to the `WHERE ssh_host_key = ''` clause. A typo there — say `WHERE id = :id` alone — would let a second connection silently overwrite a pinned key, and the callback's fake-store test would still pass. The SQL needs its own proof.

- [ ] **Step 1: Write the migration**

Create `migrations/upgrade/000009_add_node_ssh_host_key.up.sql`:

```sql
ALTER TABLE nodes ADD COLUMN ssh_host_key TEXT NOT NULL DEFAULT '';
```

Existing migrations in this directory are `.up.sql` only — do not add a `.down.sql`.

- [ ] **Step 2: Add the model field**

In `pkg/model/node.go`, add to the `Node` struct, after `SSHPassword`:

```go
	SSHHostKey  string    `db:"ssh_host_key" json:"-"`
```

This is not cosmetic. `GetByID`, `GetByNodeName`, `List` and `Page` all run `SELECT *` into `model.Node`; sqlx fails with `missing destination name ssh_host_key` on every one of them if the field is absent. `json:"-"` keeps the raw key out of API responses — the fingerprint is exposed instead, in Task 4.

- [ ] **Step 3: Write the failing repository test**

Create `pkg/repository/node_repository_test.go`:

```go
package repository

import (
	"path/filepath"
	"testing"

	"github.com/benlocal/lai-panel/migrations"
	"github.com/benlocal/lai-panel/pkg/database"
	"github.com/benlocal/lai-panel/pkg/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

// newTestNodeRepository points the package-level DB at a fresh migrated
// SQLite file and returns a repository bound to it. A file (not ":memory:")
// is used because sqlx pools connections and each pooled connection to
// ":memory:" would get its own empty database.
func newTestNodeRepository(t *testing.T) *NodeRepository {
	t.Helper()

	db, err := sqlx.Connect("sqlite3", filepath.Join(t.TempDir(), "test.db"))
	require.NoError(t, err)
	require.NoError(t, migrations.RunMigrations(db))
	t.Cleanup(func() { _ = db.Close() })

	previous := database.DB
	database.DB = db
	t.Cleanup(func() { database.DB = previous })

	return NewNodeRepository()
}

func newTestNode(t *testing.T, r *NodeRepository, name string) *model.Node {
	t.Helper()
	node := &model.Node{
		Name:        name,
		Address:     "10.0.0.5",
		SSHPort:     22,
		SSHUser:     "root",
		SSHPassword: "encrypted",
		Status:      "online",
	}
	require.NoError(t, r.Create(node))
	require.NotZero(t, node.ID)
	return node
}

func TestSetSSHHostKeyIfEmpty_PinsOnceAndRefusesOverwrite(t *testing.T) {
	r := newTestNodeRepository(t)
	node := newTestNode(t, r, "node-a")

	got, err := r.GetSSHHostKey(node.ID)
	require.NoError(t, err)
	require.Equal(t, "", got, "a new node must start unpinned")

	require.NoError(t, r.SetSSHHostKeyIfEmpty(node.ID, "ssh-ed25519 AAAAkeyA\n"))
	got, err = r.GetSSHHostKey(node.ID)
	require.NoError(t, err)
	require.Equal(t, "ssh-ed25519 AAAAkeyA\n", got)

	// The whole point of the conditional UPDATE: a second write must not win.
	require.NoError(t, r.SetSSHHostKeyIfEmpty(node.ID, "ssh-ed25519 AAAAkeyB\n"))
	got, err = r.GetSSHHostKey(node.ID)
	require.NoError(t, err)
	require.Equal(t, "ssh-ed25519 AAAAkeyA\n", got, "pinned key must survive a second write attempt")
}

func TestResetSSHHostKey_AllowsRepinning(t *testing.T) {
	r := newTestNodeRepository(t)
	node := newTestNode(t, r, "node-a")

	require.NoError(t, r.SetSSHHostKeyIfEmpty(node.ID, "ssh-ed25519 AAAAkeyA\n"))
	require.NoError(t, r.ResetSSHHostKey(node.ID))

	got, err := r.GetSSHHostKey(node.ID)
	require.NoError(t, err)
	require.Equal(t, "", got)

	require.NoError(t, r.SetSSHHostKeyIfEmpty(node.ID, "ssh-ed25519 AAAAkeyB\n"))
	got, err = r.GetSSHHostKey(node.ID)
	require.NoError(t, err)
	require.Equal(t, "ssh-ed25519 AAAAkeyB\n", got)
}

func TestUpdate_PreservesPinnedHostKey(t *testing.T) {
	r := newTestNodeRepository(t)
	node := newTestNode(t, r, "node-a")
	require.NoError(t, r.SetSSHHostKeyIfEmpty(node.ID, "ssh-ed25519 AAAAkeyA\n"))

	node.DisplayName = nil
	node.Address = "10.0.0.9"
	require.NoError(t, r.Update(node))

	got, err := r.GetSSHHostKey(node.ID)
	require.NoError(t, err)
	require.Equal(t, "ssh-ed25519 AAAAkeyA\n", got,
		"editing a node must not silently clear its pinned key; resetting is an explicit action")
}
```

- [ ] **Step 4: Run the test to verify it fails**

Run: `go test ./pkg/repository/ -run TestSetSSHHostKeyIfEmpty -v`
Expected: FAIL to compile — `r.GetSSHHostKey undefined (type *NodeRepository has no field or method GetSSHHostKey)`.

- [ ] **Step 5: Add the repository methods**

Append to `pkg/repository/node_repository.go`, after `UpdateNodeStatus`:

```go
func (r *NodeRepository) GetSSHHostKey(id int64) (string, error) {
	var key string
	err := r.db.Get(&key, "SELECT ssh_host_key FROM nodes WHERE id = ?", id)
	if err != nil {
		return "", err
	}
	return key, nil
}

// SetSSHHostKeyIfEmpty pins a host key only while the node has none. The
// conditional WHERE is what makes concurrent first connections safe: exactly
// one of them wins, and the callers read back to find out which.
func (r *NodeRepository) SetSSHHostKeyIfEmpty(id int64, key string) error {
	query := `UPDATE nodes SET ssh_host_key = :ssh_host_key,
	 updated_at = CURRENT_TIMESTAMP
	 WHERE id = :id AND ssh_host_key = ''`
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"id":           id,
		"ssh_host_key": key,
	})
	return err
}

func (r *NodeRepository) ResetSSHHostKey(id int64) error {
	query := `UPDATE nodes SET ssh_host_key = '',
	 updated_at = CURRENT_TIMESTAMP
	 WHERE id = :id`
	_, err := r.db.NamedExec(query, map[string]interface{}{"id": id})
	return err
}
```

- [ ] **Step 6: Run the tests to verify they pass**

Run: `go test ./pkg/repository/ -v`
Expected: PASS for all three tests.

- [ ] **Step 7: Run the full suite**

Run: `make test`
Expected: exit 0. If it fails with `missing destination name ssh_host_key`, Step 2 was skipped.

- [ ] **Step 8: Commit**

```bash
git add migrations/upgrade/000009_add_node_ssh_host_key.up.sql pkg/model/node.go pkg/repository/node_repository.go pkg/repository/node_repository_test.go
git commit -m "$(cat <<'EOF'
feat(node): store a pinned SSH host key per node

- Add nodes.ssh_host_key plus the matching model field: every node query is a
  SELECT * into model.Node, so the column and the struct field have to land
  together or sqlx fails to scan any node at all
- Make SetSSHHostKeyIfEmpty conditional on the column still being empty: this
  is what lets two concurrent first connections resolve to one pinned key
  instead of racing, and it is covered by its own test because a fake store
  in the callback tests could not catch a mistake in the WHERE clause
- Assert that Update leaves the pinned key alone, so that editing a node's
  address never silently re-trusts a different host

Signed-off-by: benshi <807629978@qq.com>
EOF
)"
```

---

### Task 2: The TOFU callback

**Files:**
- Create: `pkg/node/hostkey.go`
- Test: `pkg/node/hostkey_test.go` (create)

**Interfaces:**
- Consumes: nothing from Task 1 at compile time — the interface is declared here and `*repository.NodeRepository` satisfies it structurally.
- Produces: `node.HostKeyStore` (interface with `GetSSHHostKey(nodeID int64) (string, error)` and `SetSSHHostKeyIfEmpty(nodeID int64, key string) error`), and `node.TOFUHostKeyCallback(nodeID int64, nodeName string, store HostKeyStore) ssh.HostKeyCallback`.

- [ ] **Step 1: Write the failing test**

Create `pkg/node/hostkey_test.go`:

```go
package node

import (
	"crypto/ed25519"
	"crypto/rand"
	"net"
	"strings"
	"sync"
	"testing"

	"golang.org/x/crypto/ssh"
)

type fakeHostKeyStore struct {
	mu  sync.Mutex
	key string
}

func (f *fakeHostKeyStore) GetSSHHostKey(int64) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.key, nil
}

// SetSSHHostKeyIfEmpty mirrors the repository's conditional UPDATE.
func (f *fakeHostKeyStore) SetSSHHostKeyIfEmpty(_ int64, key string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.key == "" {
		f.key = key
	}
	return nil
}

func (f *fakeHostKeyStore) reset() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.key = ""
}

func testHostKey(t *testing.T) ssh.PublicKey {
	t.Helper()
	pub, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate ed25519 key: %v", err)
	}
	sshPub, err := ssh.NewPublicKey(pub)
	if err != nil {
		t.Fatalf("wrap ssh public key: %v", err)
	}
	return sshPub
}

var testRemoteAddr = &net.TCPAddr{IP: net.IPv4(10, 0, 0, 5), Port: 22}

func TestTOFU_FirstConnectionPinsKey(t *testing.T) {
	store := &fakeHostKeyStore{}
	key := testHostKey(t)

	cb := TOFUHostKeyCallback(7, "web-01", store)
	if err := cb("10.0.0.5:22", testRemoteAddr, key); err != nil {
		t.Fatalf("first connection must succeed, got %v", err)
	}

	want := string(ssh.MarshalAuthorizedKey(key))
	if got, _ := store.GetSSHHostKey(7); got != want {
		t.Fatalf("key not pinned: got %q want %q", got, want)
	}
}

func TestTOFU_SameKeyOnLaterConnection(t *testing.T) {
	store := &fakeHostKeyStore{}
	key := testHostKey(t)
	cb := TOFUHostKeyCallback(7, "web-01", store)

	if err := cb("10.0.0.5:22", testRemoteAddr, key); err != nil {
		t.Fatalf("first connection: %v", err)
	}
	if err := cb("10.0.0.5:22", testRemoteAddr, key); err != nil {
		t.Fatalf("second connection with the same key must succeed, got %v", err)
	}
}

func TestTOFU_DifferentKeyIsRejected(t *testing.T) {
	store := &fakeHostKeyStore{}
	cb := TOFUHostKeyCallback(7, "web-01", store)

	pinned := testHostKey(t)
	if err := cb("10.0.0.5:22", testRemoteAddr, pinned); err != nil {
		t.Fatalf("first connection: %v", err)
	}

	imposter := testHostKey(t)
	err := cb("10.0.0.5:22", testRemoteAddr, imposter)
	if err == nil {
		t.Fatal("a different host key must be rejected")
	}
	msg := err.Error()
	for _, want := range []string{
		"web-01",
		ssh.FingerprintSHA256(pinned),
		ssh.FingerprintSHA256(imposter),
		"reset",
	} {
		if !strings.Contains(msg, want) {
			t.Errorf("error message must contain %q, got: %s", want, msg)
		}
	}
}

func TestTOFU_RepinsAfterReset(t *testing.T) {
	store := &fakeHostKeyStore{}
	cb := TOFUHostKeyCallback(7, "web-01", store)

	if err := cb("10.0.0.5:22", testRemoteAddr, testHostKey(t)); err != nil {
		t.Fatalf("first connection: %v", err)
	}

	// What ResetSSHHostKey does in production.
	store.reset()

	rebuilt := testHostKey(t)
	if err := cb("10.0.0.5:22", testRemoteAddr, rebuilt); err != nil {
		t.Fatalf("after a reset the new key must pin, got %v", err)
	}
	want := string(ssh.MarshalAuthorizedKey(rebuilt))
	if got, _ := store.GetSSHHostKey(7); got != want {
		t.Fatalf("new key not pinned: got %q want %q", got, want)
	}
}

func TestTOFU_ConcurrentFirstConnectionsAgreeOnOneKey(t *testing.T) {
	store := &fakeHostKeyStore{}
	cb := TOFUHostKeyCallback(7, "web-01", store)

	keys := []ssh.PublicKey{testHostKey(t), testHostKey(t)}
	errs := make([]error, len(keys))

	var wg sync.WaitGroup
	for i, key := range keys {
		wg.Add(1)
		go func(i int, key ssh.PublicKey) {
			defer wg.Done()
			errs[i] = cb("10.0.0.5:22", testRemoteAddr, key)
		}(i, key)
	}
	wg.Wait()

	succeeded := 0
	for _, err := range errs {
		if err == nil {
			succeeded++
		}
	}
	if succeeded != 1 {
		t.Fatalf("exactly one concurrent first connection may succeed, got %d (%v)", succeeded, errs)
	}
}
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `go test ./pkg/node/ -run TestTOFU -v`
Expected: FAIL to compile — `undefined: TOFUHostKeyCallback`.

- [ ] **Step 3: Write the implementation**

Create `pkg/node/hostkey.go`:

```go
package node

import (
	"fmt"
	"net"

	"golang.org/x/crypto/ssh"
)

// HostKeyStore persists a node's pinned SSH host key. It is deliberately
// narrower than NodeRepository so that this package and pkg/hub stay
// independent of the concrete repository, and so the trust decision below can
// be tested without a database.
type HostKeyStore interface {
	GetSSHHostKey(nodeID int64) (string, error)
	SetSSHHostKeyIfEmpty(nodeID int64, key string) error
}

// TOFUHostKeyCallback pins a node's host key on the first connection and
// requires an exact match on every connection after that.
//
// The panel sends a decrypted SSH password to whatever answers at the node's
// address, and the health check will push and execute an agent installer over
// the same connection, so accepting an unverified key hands both to anyone who
// can intercept the traffic.
func TOFUHostKeyCallback(nodeID int64, nodeName string, store HostKeyStore) ssh.HostKeyCallback {
	return func(hostname string, _ net.Addr, key ssh.PublicKey) error {
		presented := string(ssh.MarshalAuthorizedKey(key))

		pinned, err := store.GetSSHHostKey(nodeID)
		if err != nil {
			return fmt.Errorf("read pinned host key for node %q: %w", nodeName, err)
		}

		if pinned == "" {
			if err := store.SetSSHHostKeyIfEmpty(nodeID, presented); err != nil {
				return fmt.Errorf("pin host key for node %q: %w", nodeName, err)
			}
			// Read back rather than assuming the write landed: a concurrent
			// first connection may have won the conditional update with a
			// different key, and that one is now the pinned key.
			pinned, err = store.GetSSHHostKey(nodeID)
			if err != nil {
				return fmt.Errorf("read back pinned host key for node %q: %w", nodeName, err)
			}
		}

		if pinned == presented {
			return nil
		}

		return fmt.Errorf(
			"host key mismatch for node %q (%s): pinned %s, got %s; "+
				"if this node was rebuilt, reset its host key on the Nodes page first",
			nodeName, hostname, fingerprintOf(pinned), ssh.FingerprintSHA256(key),
		)
	}
}

// fingerprintOf renders a stored authorized_keys line as a SHA256 fingerprint
// for operator-facing messages.
func fingerprintOf(authorizedKey string) string {
	key, _, _, _, err := ssh.ParseAuthorizedKey([]byte(authorizedKey))
	if err != nil {
		return "<unparseable>"
	}
	return ssh.FingerprintSHA256(key)
}
```

- [ ] **Step 4: Run the tests to verify they pass**

Run: `go test ./pkg/node/ -run TestTOFU -v`
Expected: PASS for all five tests.

- [ ] **Step 5: Run the race detector**

Run: `go test ./pkg/node/ -run TestTOFU_Concurrent -race -v`
Expected: PASS with no race reported. The concurrency test is meaningless without this.

- [ ] **Step 6: Commit**

```bash
git add pkg/node/hostkey.go pkg/node/hostkey_test.go
git commit -m "$(cat <<'EOF'
feat(node): add a TOFU SSH host key callback

- Pin a node's host key on first connection and require an exact match after
  that: the panel sends a decrypted node password to whatever answers at the
  address, and the health check pushes and runs an agent installer over the
  same connection, so an unverified key costs both the credential and remote
  code execution on what the panel believes is the node
- Read the pinned key back after writing it instead of trusting the write:
  the store's write is conditional, so a concurrent first connection may have
  pinned a different key, and that one has to win for both callers
- Take a narrow HostKeyStore rather than *repository.NodeRepository, so pkg/hub
  does not gain a repository dependency and the trust decision is testable
  without a database

Not wired into either dial site yet; that follows so the switch to enforcement
is a reviewable change on its own.

Signed-off-by: benshi <807629978@qq.com>
EOF
)"
```

---

### Task 3: Enforce at both dial sites

**Files:**
- Modify: `pkg/node/remote.go` (struct, constructor, `Init` at line 41)
- Modify: `pkg/node/state.go` (struct, `GetExec` at line 74)
- Modify: `pkg/node/manager.go` (`addNode`)
- Modify: `pkg/hub/ssh.go` (`createRemoteSession` at line 198, call site at line 63)

**Interfaces:**
- Consumes: `node.TOFUHostKeyCallback`, `node.HostKeyStore` from Task 2; the three repository methods from Task 1.
- Produces: `NewRemoteNodeExec(node *model.Node, hostKeyStore HostKeyStore) *RemoteNodeExec` (signature change) and `createRemoteSession(target *model.Node, cols int, rows int, hostKeyStore node.HostKeyStore)` (signature change).

This task has no new unit test — it is wiring, and its correctness is "the two `InsecureIgnoreHostKey` calls are gone and everything still compiles and passes". Step 1 is a grep assertion that stands in for that.

- [ ] **Step 1: Confirm the starting state**

Run: `grep -rn 'InsecureIgnoreHostKey' --include='*.go' .`
Expected: exactly two hits — `pkg/node/remote.go:41` and `pkg/hub/ssh.go:214`. If there are more, they need the same treatment and this task's file list is incomplete.

- [ ] **Step 2: Thread the store through RemoteNodeExec**

In `pkg/node/remote.go`, change the struct and constructor:

```go
type RemoteNodeExec struct {
	node         *model.Node
	hostKeyStore HostKeyStore

	sftpClient *sftp.Client
	sshClient  *ssh.Client
}

func NewRemoteNodeExec(node *model.Node, hostKeyStore HostKeyStore) *RemoteNodeExec {
	return &RemoteNodeExec{
		node:         node,
		hostKeyStore: hostKeyStore,
	}
}
```

and in `Init()` replace the `HostKeyCallback` line:

```go
	config := &ssh.ClientConfig{
		User: r.node.SSHUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: TOFUHostKeyCallback(r.node.ID, r.node.Name, r.hostKeyStore),
	}
```

- [ ] **Step 3: Give NodeState the store**

In `pkg/node/state.go`, add the field:

```go
type NodeState struct {
	info         model.Node
	hostKeyStore HostKeyStore
	exec         NodeExec
	dockerClient *dockerClient.Client

	execMu         sync.RWMutex
	dockerClientMu sync.RWMutex
}
```

and update the remote branch of `GetExec`:

```go
	if n.info.IsLocal {
		exec = NewLocalNodeExec()
	} else {
		exec = NewRemoteNodeExec(&n.info, n.hostKeyStore)
	}
```

In `pkg/node/manager.go`, populate it in `addNode`:

```go
	state := NodeState{
		info:         *node,
		hostKeyStore: m.nodeRepository,
	}
```

`*repository.NodeRepository` satisfies `HostKeyStore` after Task 1, and `manager.go` already imports the repository package, so no new import and no import cycle.

- [ ] **Step 4: Thread the store into the web terminal**

In `pkg/hub/ssh.go`, add the import:

```go
	"github.com/benlocal/lai-panel/pkg/node"
```

Change the signature — note the parameter is renamed from `node` to `target`, because the existing name would shadow the package and make `node.HostKeyStore` fail to compile:

```go
func createRemoteSession(target *model.Node, cols int, rows int, hostKeyStore node.HostKeyStore) (*ssh.Client, *ssh.Session, io.WriteCloser, io.ReadCloser, io.ReadCloser, error) {
```

Inside the function, rename every remaining `node.` reference on the parameter to `target.` — there are three: `node.GetDecryptedSSHPassword()`, `node.SSHUser`, and the `fmt.Sprintf("%s:%d", node.Address, node.SSHPort)` in the dial. Replace the callback:

```go
	clientConfig := &ssh.ClientConfig{
		User:            target.SSHUser,
		Auth:            []ssh.AuthMethod{ssh.Password(pwd)},
		HostKeyCallback: node.TOFUHostKeyCallback(target.ID, target.Name, hostKeyStore),
		Timeout:         10 * time.Second,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", target.Address, target.SSHPort), clientConfig)
```

Update the call site at line 63 — `SimpleHub` already holds `nodeRepository`:

```go
		sshClient, session, stdin, stdout, stderr, err := createRemoteSession(targetNode, cols, rows, h.nodeRepository)
```

- [ ] **Step 5: Verify no bypass remains**

Run: `grep -rn 'InsecureIgnoreHostKey' --include='*.go' .`
Expected: no output.

- [ ] **Step 6: Build and test**

Run: `make test`
Expected: exit 0.

- [ ] **Step 7: Commit**

```bash
git add pkg/node/remote.go pkg/node/state.go pkg/node/manager.go pkg/hub/ssh.go
git commit -m "$(cat <<'EOF'
fix(ssh): verify node host keys instead of accepting any key

- Replace ssh.InsecureIgnoreHostKey() at both dial sites — the SFTP/exec path
  in pkg/node and the web terminal in pkg/hub — with the shared TOFU callback,
  so a substituted host no longer receives the node's decrypted SSH password
- Carry the store from NodeManager through NodeState into RemoteNodeExec
  rather than reaching for a global: NodeManager already owns the repository,
  and keeping the dependency explicit is what let the callback stay testable
- Rename createRemoteSession's `node` parameter to `target`: it shadowed the
  newly imported pkg/node and would not compile otherwise

Existing nodes are unaffected at connect time — they hold no pinned key, so
the next connection pins whatever they present, and enforcement starts from
the connection after that.

Signed-off-by: benshi <807629978@qq.com>
EOF
)"
```

---

### Task 4: Reset endpoint and fingerprint exposure

**Files:**
- Modify: `pkg/model/node.go` (`NodeView`, `ToView`, fingerprint helper)
- Modify: `pkg/handler/node.go` (new handler)
- Modify: `cmd/serve/main.go` (route registration)

**Interfaces:**
- Consumes: `(*NodeRepository).ResetSSHHostKey` from Task 1; `model.Node.SSHHostKey` from Task 1.
- Produces: `POST /api/node/trust/reset` accepting `{"id": <int64>}`; `NodeView.SSHHostKeyFingerprint` serialised as `ssh_host_key_fingerprint`.

- [ ] **Step 1: Expose the fingerprint on the view**

In `pkg/model/node.go`, add the import `"golang.org/x/crypto/ssh"`, add a method next to `GetDecryptedSSHPassword`:

```go
// SSHHostKeyFingerprint renders the pinned host key for display. It returns
// an empty string when nothing is pinned yet, which the UI shows as
// "not yet trusted".
func (n *Node) SSHHostKeyFingerprint() string {
	if n.SSHHostKey == "" {
		return ""
	}
	key, _, _, _, err := ssh.ParseAuthorizedKey([]byte(n.SSHHostKey))
	if err != nil {
		return ""
	}
	return ssh.FingerprintSHA256(key)
}
```

add the field to `NodeView` after `AgentPort`:

```go
	SSHHostKeyFingerprint string `json:"ssh_host_key_fingerprint"`
```

and set it in `ToView`:

```go
		SSHHostKeyFingerprint: n.SSHHostKeyFingerprint(),
```

- [ ] **Step 2: Add the handler**

Append to `pkg/handler/node.go`:

```go
func (h *BaseHandler) ResetNodeHostKeyHandler(ctx context.Context, c *app.RequestContext) {
	type resetNodeHostKeyRequest struct {
		ID int64 `json:"id"`
	}

	var req resetNodeHostKeyRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.Error(err)
		return
	}

	if req.ID <= 0 {
		c.Error(errors.New("ID is required"))
		return
	}

	if err := h.NodeRepository().ResetSSHHostKey(req.ID); err != nil {
		c.Error(err)
		return
	}

	// Drop the cached NodeState: it holds a live ssh.Client opened under the
	// old key, so without this the reset would not take effect until the
	// panel restarted.
	h.NodeManager().RemoveNode(req.ID)

	c.JSON(http.StatusOK, EmptyResponse())
}
```

- [ ] **Step 3: Register the route**

In `cmd/serve/main.go`, immediately after the `api.POST("/node/page", h.GetNodePageHandler)` line:

```go
		api.POST("/node/trust/reset", h.ResetNodeHostKeyHandler)
```

It goes inside the `api` group so it inherits `h.AuthMiddleware`. Do not add it to `cmd/agent/main.go` — agents have no node repository.

- [ ] **Step 4: Build and test**

Run: `make test`
Expected: exit 0.

- [ ] **Step 5: Verify the endpoint by hand**

Run:

```bash
make serve && ./bin/serve &
sleep 3
TOKEN=$(curl -s -X POST localhost:8080/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"'"$(printf %s 'admin' | base64)"'"}' \
  | python3 -c 'import sys,json; print(json.load(sys.stdin)["data"]["token"])')
curl -s -X POST localhost:8080/api/node/list -H "Authorization: Bearer $TOKEN" | head -c 400
curl -s -X POST localhost:8080/api/node/trust/reset \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' -d '{"id":1}'
```

Expected: the node list contains `"ssh_host_key_fingerprint"`, and the reset call returns `{"code":0,"message":"success","data":null}`. Substitute the real admin password if it has been changed. Kill the server afterwards.

- [ ] **Step 6: Commit**

```bash
git add pkg/model/node.go pkg/handler/node.go cmd/serve/main.go
git commit -m "$(cat <<'EOF'
feat(node): add an endpoint to reset a node's pinned host key

- Add POST /api/node/trust/reset: a rebuilt node presents a new host key and
  would otherwise be permanently unreachable, including to the health check
  that would have reinstalled its agent, with editing SQLite as the only way out
- Evict the cached NodeState on reset: it holds a live ssh.Client dialled
  under the old key, so clearing only the database row would leave the reset
  ineffective until the next panel restart
- Surface the fingerprint on NodeView rather than the key itself, so the UI can
  show what is currently trusted without putting the raw key in API responses

Signed-off-by: benshi <807629978@qq.com>
EOF
)"
```

---

### Task 5: Nodes page — show the fingerprint, offer the reset

**Files:**
- Modify: `dashboard/src/api/node.ts`
- Modify: `dashboard/src/views/Nodes.vue`

**Interfaces:**
- Consumes: `POST /api/node/trust/reset` and `ssh_host_key_fingerprint` from Task 4.
- Produces: nothing downstream.

- [ ] **Step 1: Extend the API client**

In `dashboard/src/api/node.ts`, add to the `Node` interface after `agent_port`:

```ts
  ssh_host_key_fingerprint?: string;
```

and add a method to `nodeApi`, after `delete`:

```ts
  async resetHostKey(id: number): Promise<ApiResponse<void>> {
    return post<void>("/api/node/trust/reset", { id });
  },
```

- [ ] **Step 2: Add the reset flow to Nodes.vue**

In the `<script setup>` block, next to the existing delete state, add:

```ts
const isResetHostKeyDialogOpen = ref(false);
const nodeToResetHostKey = ref<Node | null>(null);

const openResetHostKeyDialog = (node: Node) => {
  nodeToResetHostKey.value = node;
  isResetHostKeyDialogOpen.value = true;
};

const confirmResetHostKey = async () => {
  if (!nodeToResetHostKey.value) return;
  loading.value = true;
  const res = await nodeApi.resetHostKey(nodeToResetHostKey.value.id);
  loading.value = false;
  if (ApiResponseHelper.isSuccess(res)) fetchNodes();
  else showToast(res.message ?? "Failed to reset host key", "error");
  isResetHostKeyDialogOpen.value = false;
  nodeToResetHostKey.value = null;
};
```

`ref`, `Node`, `nodeApi`, `ApiResponseHelper`, `showToast`, `loading` and `fetchNodes` are all already imported or defined in this file — do not re-import them.

- [ ] **Step 3: Add the column and the action button**

Add a column before `<Column header="Actions">`:

```vue
        <Column header="Host Key">
          <template #body="{ data }">
            <span v-if="data.is_local" class="text-muted-foreground">—</span>
            <span v-else-if="data.ssh_host_key_fingerprint" class="font-mono text-xs">{{ data.ssh_host_key_fingerprint }}</span>
            <Tag v-else value="Not yet trusted" severity="warn" />
          </template>
        </Column>
```

and add a button inside the existing `.action-btns` div, between Edit and Delete:

```vue
              <Button text rounded size="small" :disabled="data.is_local || !data.ssh_host_key_fingerprint" @click="openResetHostKeyDialog(data)" v-tooltip.top="'Reset host key'">
                <i class="pi pi-shield"></i>
              </Button>
```

It is disabled for local nodes (they never use SSH) and for nodes with nothing pinned yet (there is nothing to reset).

- [ ] **Step 4: Add the confirmation dialog**

After the existing delete `<Dialog>`:

```vue
    <Dialog v-model:visible="isResetHostKeyDialogOpen" modal header="Reset Host Key" :style="{ width: '480px' }">
      <p>
        Node "{{ nodeToResetHostKey?.display_name || nodeToResetHostKey?.name }}" currently trusts
        <span class="font-mono text-xs">{{ nodeToResetHostKey?.ssh_host_key_fingerprint }}</span>.
      </p>
      <p class="text-muted-foreground">
        The next connection will trust whatever host key it is offered. Only do this if you know the node was rebuilt.
      </p>
      <template #footer>
        <Button outlined @click="isResetHostKeyDialogOpen = false">Cancel</Button>
        <Button severity="danger" @click="confirmResetHostKey" :disabled="loading">{{ loading ? "Resetting..." : "Reset" }}</Button>
      </template>
    </Dialog>
```

- [ ] **Step 5: Type-check and build**

Run: `cd dashboard && npm run build`
Expected: exit 0. `Tag` and `Dialog` are already imported in this file; if `vue-tsc` reports either as undefined, add the missing `import` alongside the existing PrimeVue imports.

- [ ] **Step 6: Commit**

```bash
git add dashboard/src/api/node.ts dashboard/src/views/Nodes.vue
git commit -m "$(cat <<'EOF'
feat(dashboard): show node host key fingerprints and allow resetting them

- Add a Host Key column so operators can compare what the panel trusts against
  what `ssh-keyscan` reports, and can tell an unverified node from a pinned one
- Add a reset action guarded by a dialog that names the currently trusted
  fingerprint and states plainly that the next connection trusts whatever it
  is offered: this is the supported recovery path after a rebuild, so it has
  to be obvious that it reopens the trust window
- Disable the action for local nodes and for nodes with nothing pinned, since
  neither has a host key to reset

Signed-off-by: benshi <807629978@qq.com>
EOF
)"
```

---

## Final verification

- [ ] `grep -rn 'InsecureIgnoreHostKey' --include='*.go' .` returns nothing.
- [ ] `make test` exits 0.
- [ ] `cd dashboard && npm run build` exits 0.
- [ ] `git log --oneline` shows five commits, one per task.
- [ ] Against a real remote node: the first connection succeeds and the Nodes page then shows a fingerprint; a second connection succeeds; changing the node's `address` to a different host produces the mismatch error naming both fingerprints; resetting and reconnecting pins the new host.
