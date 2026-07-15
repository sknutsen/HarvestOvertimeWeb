# AGENTS.md

## Cursor Cloud specific instructions

This is a single Go web app (`harvestovertimeweb`): Echo + [templ](https://templ.guide/) + HTMX,
server-rendered. It computes Harvest overtime. Entry point is `main.go`; it starts an HTTP
server on `$PORT` (see `.env`). Standard commands live in `.air.toml` and
`.github/workflows/go.yml`.

### Services & how to run them

| Service | Command | Notes |
|---|---|---|
| Web app (dev, hot reload) | `air` | Uses `.air.toml`; runs `templ generate && go build` then serves on `$PORT` (8080). |
| Web app (plain) | `go run .` | Uses committed `*_templ.go`; no regeneration. |
| Build | `go build ./...` | CI equivalent: `go build -v ./...`. |
| Test | `go test ./...` | No `*_test.go` files exist yet, so this is currently a no-op (CI still runs it). |
| Vet | `go vet ./...` | No project-specific linter configured. |

The dev tools `templ` and `air` are installed via `go install` into `$(go env GOPATH)/bin`
(`/home/ubuntu/go/bin`), which is added to `PATH` in `~/.bashrc`. If `air`/`templ` are not
found, run `export PATH="$PATH:$(go env GOPATH)/bin"`.

### Non-obvious caveats

- **`air` regenerates `view/*_templ.go` on every run** (its build cmd is `templ generate && ...`).
  The committed files were generated with an older templ version, so running `air`/`templ generate`
  produces a 1-line version-comment diff in `view/details_templ.go`, `view/index_templ.go`,
  `view/shared_templ.go`. This is harmless — do NOT commit that churn (`git checkout -- view/*_templ.go`).
- **`.env` is gitignored** and must exist for the app to start cleanly. It needs:
  `PORT`, `HARVEST_CLIENT_ID`, `HARVEST_CLIENT_SECRET`, `USER_AGENT`, `SECRET`, and the DB vars
  `DATABASE`, `DB_HOST`, `DB_USER`, `DB_PASS`, `DB_PORT`. Note `example.env` is missing the DB vars.
- **PostgreSQL is optional but improves fidelity.** It only supplies Norwegian holidays and the
  "years" list (table `dates`, columns `calendar,date,description`, filtered on `calendar='NO'`).
  All DB errors are swallowed — the app still renders without a DB, just without holidays/years.
  There is no migration/seed file in the repo; the `dates` table must be created and seeded manually.
  Postgres is installed in this VM (cluster `16 main`, db `harvest`, user/pass `harvest/harvest`,
  seeded `dates`). Start it with `sudo pg_ctlcluster 16 main start` (the update script does NOT
  start services). `lib/pq` uses `sslmode=require` by default (the connection string has no
  `sslmode`); the local cluster has SSL enabled, so it connects fine.
- **Auth is external:** `/` requires an `accesstoken` cookie, otherwise it 307-redirects to `/auth`
  → Harvest OAuth (`id.getharvest.com`), which needs a real Harvest account + registered OAuth app.
  `POST /hours` and the OAuth callback cannot be exercised end-to-end without real Harvest
  credentials. To view the rendered form locally without Harvest, set a dummy cookie
  (e.g. `document.cookie="accesstoken=devtoken"` in the browser console, or `curl -b accesstoken=x`);
  the page renders with DB-backed years but live time-entry data / the "Get Hours" calculation
  require valid Harvest auth.
