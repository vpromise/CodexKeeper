# CodexKeeper

[简体中文](README.zh.md)

A standalone Codex and Claude usage dashboard for CodexProxy, based on
[CPA Usage Keeper](https://github.com/Willxup/cpa-usage-keeper).
The original dashboard design, charts, light/dark themes and responsive layout are preserved.

## Scope

- Overview, realtime monitoring, analysis, request details and exports.
- Tokens, cache usage, latency, failures and estimated costs.
- Codex / Claude OAuth accounts and native API key credentials, with names and aliases.
- Native quota and subscription information, Codex quota history.
- One administrator password, session management, SQLite storage and backups.

Removed: community/local leaderboards, ranking submission and profiles, API Key user
login and separate user pages, non-native provider integrations, legacy queue-key
probing and HTTP ingestion fallback. A Claude/GPT model routed through another
provider is not treated as native Claude/Codex usage.

This fork reads the current CodexProxy management protocol. It does not add a
module to CodexProxy. Pricing catalogs remain generic data sources; only native
OpenAI/Anthropic model families have special price-selection rules.

## Connect to CodexProxy with Docker

Use one Keeper collector per CPA instance. Enable `usage-statistics-enabled: true`
in CodexProxy, allow management access from its Docker network, and use its
management key (not a client inference key).

```sh
cp .env.example .env
# Set CPA_MANAGEMENT_KEY and a private LOGIN_PASSWORD in .env.
# Set CPA_DOCKER_NETWORK if the existing Docker network has another name.
docker compose --project-directory . -f deploy/docker-compose.example.yml up -d --build
```

The example joins the existing `sub2api-production_default` network and contacts
CodexProxy at `cpa:8317`. Change the service DNS name/port if needed. It builds this
fork locally; it does not pull the upstream Keeper image. Persist `./data`.
Use a fresh data directory for this fork; historical Keeper import is not included.
Existing databases are not purged, and historical schema migrations remain for
safe upgrades, so a reused database can still contain earlier provider data.

No domain is required for collection. The UI binds to host loopback port 8080 by
default. Use an SSH tunnel for private access or reverse-proxy it with HTTPS for
remote browser access. A dedicated hostname is optional. For a `/keeper` subpath,
set `APP_BASE_PATH=/keeper` and preserve that prefix in the reverse proxy. Set
`CPA_PUBLIC_URL` to the browser-facing CPA address for the Back to CPA link.
A normal HTTP reverse proxy does not forward the RESP stream: keep
`REDIS_QUEUE_ADDR=cpa:8317` on the internal Docker network.

To use the same administrator password as CodexProxy without storing its plaintext, leave `LOGIN_PASSWORD` empty and set `LOGIN_PASSWORD_HASH` to its bcrypt `remote-management.secret-key`. Single-quote the hash in `.env` to preserve `$` characters. Configure exactly one password option. This copies the current password hash; later CPA password changes must be applied to Keeper separately.

## Collection and reliability

Keeper subscribes to `usage`, drains the fixed `usage` backlog with `LPOP`, and
repeats this sequence on reconnect. Supported events are persisted to a SQLite
inbox before aggregation; consumed batches are retried on database write failure.
`REDIS_QUEUE_RETRY_INTERVAL` sets the initial retry delay (default 1s).
The upstream stream has no durable acknowledgement: a process crash between
receiving an event and writing it can still lose that event. Error events are
live-only and have no backlog. Keep the collector running.

## Develop and verify

Requires Go 1.26+, a C compiler for SQLite/CGO, and Node.js 24.

```sh
npm --prefix web ci
npm --prefix web run build
go run ./cmd/server --env .env --host 127.0.0.1

go test ./cmd/... ./internal/...
go build -o /tmp/codexkeeper ./cmd/server
npm --prefix web run typecheck
npm --prefix web run lint
npm --prefix web test
```

On Node.js 26, run tests with
`NODE_OPTIONS=--no-experimental-webstorage npm --prefix web test` to prevent Node's
native localStorage from conflicting with the existing happy-dom test environment.

The Go module name and internal binary name remain `cpa-usage-keeper` to avoid
unnecessary import/build churn. The project URL and update checks point to
[vpromise/CodexKeeper](https://github.com/vpromise/CodexKeeper).

## License

MIT. Original work by CPA Usage Keeper contributors; see [LICENSE](LICENSE).
