# SSH Host Key Verification

**Date:** 2026-08-19
**Status:** Approved, not yet implemented

## Problem

Both SSH dial sites in the panel accept any host key:

- `pkg/node/remote.go:41` — `RemoteNodeExec.Init()`, used for file transfer (SFTP) and remote command execution, including the deploy pipeline and the health check's agent auto-install.
- `pkg/hub/ssh.go:214` — `createRemoteSession()`, used by the web terminal.

Both use `ssh.InsecureIgnoreHostKey()`. Both authenticate with a password that the panel decrypts from `nodes.ssh_password` (AES-GCM, `pkg/crypto`) and sends to whatever answers on `address:ssh_port`.

Anyone able to intercept panel-to-node traffic therefore receives the node's SSH password in cleartext, and can then serve arbitrary command output back to the panel. The health check compounds this: when a node looks offline it pushes `bin/agent` and `install.sh` over SFTP and executes the installer, so a successful MITM gets remote code execution on whatever the panel believes is the node.

Dependabot does not and will not report this — it is a code defect, not a vulnerable dependency.

## Trust model: TOFU with pinning

The first successful connection to a node records that node's host key. Every later connection must present the same key, or the connection fails.

This was chosen over requiring operators to supply a fingerprint when adding a node. Strict entry is stronger — it closes the first-connection window that TOFU leaves open — but it disconnects every node already in the database until someone fills in fingerprints by hand, and it interrupts the health check's auto-install path at the same time. TOFU takes existing deployments from "never verified" to "verified after the first connection" with no operator action, which is the larger share of the available protection.

A `known_hosts` file under the data path was also considered. It reuses OpenSSH semantics and tooling (`ssh-keyscan`), but it keeps trust state outside the node record, so the UI cannot show or reset a node's pinned key. Rejected for that reason.

## Design

### Shared callback

New file `pkg/node/hostkey.go`:

```go
// HostKeyStore persists a node's pinned SSH host key.
type HostKeyStore interface {
    GetSSHHostKey(nodeID int64) (string, error)
    SetSSHHostKeyIfEmpty(nodeID int64, key string) error
}

func TOFUHostKeyCallback(nodeID int64, nodeName string, store HostKeyStore) ssh.HostKeyCallback
```

The callback takes an interface rather than `*repository.NodeRepository` so that `pkg/hub` and `pkg/node` stay decoupled from the concrete repository, and so the security logic is unit-testable without a database.

Callback behaviour on each connection:

1. Read the stored key for the node.
2. **Empty** — call `SetSSHHostKeyIfEmpty`, then read back and compare. The conditional write means only one of several concurrent first connections wins; the losers read back the winner's key and compare against it like any other connection. Without the read-back, two simultaneous first connections to different endpoints would both succeed.
3. **Non-empty** — compare byte-for-byte against the presented key. On mismatch, return an error naming the node, both fingerprints, and the recovery step:
   `host key mismatch for node "web-01" (10.0.0.5:22): pinned SHA256:aaa…, got SHA256:bbb…; if this node was rebuilt, reset its host key on the Nodes page first`

Keys are stored as the `authorized_keys` text produced by `ssh.MarshalAuthorizedKey()`. Fingerprints for display come from `ssh.FingerprintSHA256()`.

Nodes with `is_local = true` never reach this code — they use `LocalNodeExec` and do not open an SSH connection.

### Storage

Migration `migrations/upgrade/000009_add_node_ssh_host_key.up.sql`:

```sql
ALTER TABLE nodes ADD COLUMN ssh_host_key TEXT NOT NULL DEFAULT '';
```

`NodeRepository` gains three methods:

| Method | SQL |
|---|---|
| `GetSSHHostKey(id)` | `SELECT ssh_host_key FROM nodes WHERE id = ?` |
| `SetSSHHostKeyIfEmpty(id, key)` | `UPDATE nodes SET ssh_host_key = ? WHERE id = ? AND ssh_host_key = ''` |
| `ResetSSHHostKey(id)` | `UPDATE nodes SET ssh_host_key = '' WHERE id = ?` |

The host key is not a secret and is not encrypted.

### Wiring the two call sites

`NodeState` currently builds `NewRemoteNodeExec(&n.info)` with no repository reference. `NodeManager` already holds a `*repository.NodeRepository`, so it passes that into `NodeState` at construction, and `NodeState` passes it to `NewRemoteNodeExec`.

`SimpleHub` already holds a `*repository.NodeRepository`. `createRemoteSession` takes an extra `HostKeyStore` parameter.

### Reset path

A node that is legitimately rebuilt presents a new host key, and a hard failure would leave it permanently unreachable — including to the health check that would otherwise reinstall its agent. Recovery must not require editing SQLite.

- `POST /api/node/trust/reset` with `{"id": <nodeID>}`, registered inside the authenticated `/api` group in `cmd/serve/main.go`, handled by `BaseHandler.ResetNodeHostKeyHandler`.
- `NodeView` gains `ssh_host_key_fingerprint` (empty string when nothing is pinned yet) so the UI can show the current fingerprint.
- `dashboard/src/api/node.ts` gains the corresponding call; `Nodes.vue` adds "view fingerprint" and "reset host key" to the per-node actions, with a confirmation noting that the next connection will pin whatever key it is offered.

### Deliberately excluded

- **No environment-variable bypass.** A `PANEL_SKIP_HOST_KEY_CHECK` style switch becomes the production default in practice, which would leave the original defect in place behind a flag. Operators who need to accept a new key use the reset action, one node at a time.
- **No change to the agent's `/docker.proxy`.** That is plain HTTP between panel and agent, a separate transport with a separate threat model.

## Testing

`pkg/node/hostkey_test.go`, using a fake in-memory `HostKeyStore`. No database, so it runs inside `make test`:

1. First connection with an empty store pins the presented key and succeeds.
2. Second connection presenting the same key succeeds.
3. Connection presenting a different key fails, and the error names both fingerprints.
4. After the stored key is cleared (what `ResetSSHHostKey` does in production; the fake store is cleared directly, since `ResetSSHHostKey` is called by the handler and is deliberately not part of `HostKeyStore`), the next connection pins the new key and succeeds.
5. Concurrent first connections presenting *different* keys: exactly one succeeds, the other fails.

## Impact on existing deployments

None at connect time. Every node already in the database has an empty `ssh_host_key` and pins on its next connection. The deploy pipeline, web terminal, and health-check auto-install continue to work without operator action.

The change takes effect from the second connection onward: from then on a substituted host key is refused rather than silently trusted.
