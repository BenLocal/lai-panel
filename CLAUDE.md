# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
make all                 # build both binaries into bin/ (serve + agent)
make serve               # build bin/serve only     (DEBUG=true keeps symbols)
make agent               # build bin/agent only
make run-serve           # build + run the master on :8080
make run-agent           # build + run an agent on :8081
make test                # go test -v ./...
go test ./pkg/tmpl/ -run TestParseDockerComposeWithEnv -v   # single test

cd dashboard
npm install
npm run dev              # Vite dev server on :5173, proxies /api -> :8080
npm run build            # vue-tsc type-check + build into dashboard/dist
```

`dashboard/dashboard.go` does `//go:embed dist/*`, so **any Go build fails if `dashboard/dist` is empty**. The `dashboard-dist` Make target (a prerequisite of `serve` and `test`) creates a `.keep` placeholder for this reason — invoke `go build`/`go test` through Make, or create the directory yourself. To ship a real UI, run `npm run build` before `make serve`.

There is no linter configured; CI (`.github/workflows/go.yml`) only runs `make` and `make test`.

Env vars read at startup: `PANEL_PORT`, `PANEL_MASTER_HOST`, `PANEL_MASTER_PORT` (serve), `LAI_PANEL_ENCRYPTION_KEY` (both). Agent takes CLI flags instead: `--master-host --master-port --name --address`.

## Architecture

### One codebase, two binaries

`cmd/serve` (master/panel) and `cmd/agent` (per-node daemon) share all of `pkg/`. They differ in three places:

1. **Route table** — each `main.go` calls `api.DefaultRegistry.Add(func(h *handler.BaseHandler, router *route.Engine){...})` in `init()`; `api.ApiServer.Start` replays those bindings onto Hertz. Adding an endpoint means editing the relevant `main.go`.
2. **`ctx.AppCtx`** (`pkg/ctx/app.go`) is the DI container and is built differently per mode. In agent mode it holds *only* options, the Docker proxy and the server store — **every repository, `NodeManager`, and the SignalR server are nil**. A handler that touches them will panic if you also register it on the agent.
3. **Background services** — both runtimes register `gracefulshutdown.Service` implementations (`Name/Start/Shutdown`) with `g.Add(...)`; see `cmd/*/runtime.go` for the exact set.

`options.IOptions.Agent()` is the runtime discriminator used throughout `pkg/`.

### Master ↔ agent control plane

- **Registration**: `service.RegistryService` POSTs `/registry` to the master every 30s (`pkg/client/registry.go`). The master upserts the `nodes` row and returns its ID, which the agent stores in the global `ctx.ServerStore`. The master's own node row is created locally as `is_local = true, name = "local"`.
- **Docker over HTTP**: the agent exposes `/docker.proxy/*path`, a `httputil.ReverseProxy` onto the local Docker socket (`pkg/docker/docker_proxy.go`). The master talks to remote Docker by pointing a normal Docker SDK client at `tcp://<addr>:<agentPort>/docker.proxy` (`docker.AgentDockerClient`). No Docker daemon is ever exposed directly.
- **Self-healing**: `service.HealthCheckService` pings each remote node's `/healthz` every 30s. On failure it marks the node offline **and pushes `bin/agent` + `install.sh` over SFTP from `<dataPath>/static/install/` and runs the installer over SSH** — so a dead agent reinstalls itself as long as SSH credentials are valid.
- These control-plane routes (`/healthz`, `/registry`, `/docker_event`) are deliberately outside the auth middleware.

### Node abstraction

`node.NodeManager` caches one `node.NodeState` per node ID. `NodeState` lazily builds two things and memoises them:

- a `node.NodeExec` — `LocalNodeExec` (os/exec) or `RemoteNodeExec` (SSH + SFTP, password decrypted from the DB via `pkg/crypto`), same interface for file I/O and streaming command execution;
- a Docker client — local socket or the agent proxy above.

All deployment work goes through `NodeExec`, so pipeline steps are agnostic about local vs remote.

### Deploy pipeline

`pkg/pipe` defines a tiny generic `Processor[T]` with `Sequence(...)`; on error it calls `Cancel` on the already-run processors **in reverse order**. `DeployPipeline.Up` chains: cleanup workspace → copy workspace → download installer artifact → parse compose template → load image → `docker compose up -d --build`. `Down` is just compose down.

Progress is streamed to the browser as SSE — `DeployCtx.Send(event, data)` writes to a `sse.Writer` held by the handler (`HandleDockerComposeDeploy`); the frontend consumes it via `stream()` in `dashboard/src/api/base.ts`, treating `event: done` as completion.

An app's `docker_compose` column is a **Go `text/template`**, rendered with the service's QA values as the data map plus builtin funcs registered in `deploypipe.builtinFuncMap`: `panel_env`, `is_agent`, `master_host`, `master_port`. After rendering, `com.lai-panel.{managed-by,owner,service}` labels are injected into every service (`pkg/constant`). Those labels are the only link back to the panel — `service.ServicesStateUpdater` filters containers by them every 25s to reconcile service status.

### Web terminals

`pkg/hub` mounts a SignalR server at `/api/signalr` (WebSockets + SSE transports). `SimpleHub` keeps SSH-PTY and docker-exec sessions in maps keyed by SignalR connection ID and tears them down in `OnDisconnected`. Frontend: `@microsoft/signalr` + xterm.js.

### Persistence

sqlx over SQLite, wrapped by an `ngrok/sqlmw` interceptor registered as the `sqlite3-mw` driver — **every statement, row scan and tx is logged**, which is noisy but is the intended debugging surface. Migrations live in `migrations/upgrade/NNNNNN_*.up.sql`, are `embed.FS`-bundled and run by golang-migrate on every startup; an admin user is ensured afterwards (`pkg/database/create_user.go`). Repositories in `pkg/repository` are thin hand-written SQL over `database.GetDB()` and are constructed once in `AppCtx`.

Models carry both `db:` and `json:` tags. Where a column stores JSON (`apps.qa`, `apps.metadata`, `services.metadata`, `services.deploy_info`) the model keeps it as `*string` and exposes `ToView()`/`ToModel()` converters — handlers deal in the `*View` types.

### HTTP conventions

- Everything under `/api` is **POST** (including reads: `/api/node/list`, `/api/application/page`, …) and passes through `BaseHandler.AuthMiddleware` (JWT bearer, `/api/auth/login` exempted). `/open/*` is unauthenticated (static files, uploads).
- Handlers don't return errors; they call `c.Error(err)` and return. `ErrorHandlerMiddleware` converts the last error into **HTTP 200** with `{code, message, data}` where `code` is an HTTP-ish status. Success uses `SuccessResponse(data)` → `code: 0`. The frontend only checks `code === 0`.
- Login sends the password base64-encoded; it is compared against an AES-GCM ciphertext in the DB (`pkg/crypto`, key from `LAI_PANEL_ENCRYPTION_KEY`, otherwise a hardcoded default). Node SSH passwords are stored the same way. The JWT secret is currently a hardcoded const in `pkg/handler/auth.go`.
- Unknown paths fall through to `ui()` in `cmd/serve/main.go`, which serves the embedded SPA's `index.html` — client-side routes work on refresh.

### Data layout on disk

`options.InitOptions` creates `<dataPath>/{log,workspace,service,static}` where `dataPath` is `/var/lai-panel/data` on Linux (`pkg/options/data_path_linux.go`) and `$HOME/.lai-panel/data` elsewhere. Deployed services live in `<dataPath>/service/<serviceName>/`, per-app workspace templates in `<dataPath>/workspace/<appName>/`. On a remote node the base comes from `nodes.data_path` reported at registration, not from the master's own config.

### Frontend

Vue 3 `<script setup>` + PrimeVue 4 + Tailwind 4, built by rolldown-vite. Routes in `src/router/index.ts` guard on `meta.requiresAuth`; the token and user live in `localStorage` behind `src/auth/index.ts`. All backend calls go through the `post`/`get`/`stream` helpers in `src/api/base.ts`, grouped per domain in `src/api/*.ts`. `@` aliases `src/`.

### Adding a feature end to end

migration in `migrations/upgrade/` → model in `pkg/model` → repository in `pkg/repository` (and wire it into `ctx.AppCtx` + a `BaseHandler` accessor) → handler method on `*BaseHandler` → route in `cmd/serve/main.go`'s `init()` → typed client in `dashboard/src/api/` → view.
