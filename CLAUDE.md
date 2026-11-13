# Prisme

Team chat for the Facile suite on Matrix: a Go API that is a Synapse application service, a
SvelteKit client on matrix-js-sdk, and Synapse itself, all in one Dokploy compose. Read
`docs/plan.md` for the product and `README.md` for configuration and design notes.

## Commands

```sh
docker compose -p prisme-dev -f docker-compose.dev.yml up -d --build   # Postgres :5557, Synapse :18008
mise run dev                  # API :4000 and client :5173
sh scripts/check.sh           # gofmt, vet, go test, svelte-check, client tests
filet check .                 # must report no error before any push
filet gate .                  # the pre-push hook: format, filet, vet, staticcheck, tests, douane
cd apps/client && bun run build
cd apps/desktop && bunx @tauri-apps/cli@2 build --no-bundle     # desktop app (needs the client build)
cd apps/client && bunx cap sync android                          # needs Node 22: mise x node@22 --
```

## Layout

- `apps/api/modules/<name>/` holds `router.go` (filet requires it), `service.go`, `types.go`,
  `documentation.go`. Every route must appear in the module's `Documentation()`, or
  `TestEveryRouteIsDocumented` fails. porte's own routes (`/api/auth/*`) are exempt.
- `apps/api/internal/matrix` is the appservice client: every call carries the appservice token,
  and a call for a user passes `user_id` (identity assertion).
- `apps/api/internal/accounts` maps a porte user to exactly one Matrix account.
- `deploy/synapse` renders `homeserver.yaml` and `appservice.yaml` from the environment at start.

## Rules that bite

- Auth is porte SSO only (porte mounts no password route while `OIDC_ISSUER` is set), plus the service-key path for Bulle on the routes that
  call `Require(scope)` with a non-empty scope. `Require("")` keeps a route to browser sessions.
- Synapse 1.164 creates room version 12: never put the room creator in
  `power_level_content_override.users`.
- Synapse's user directory never returns Prisme users (they are all appservice users); search
  through `GET /api/users`, which shows only the caller's space peers and own agents.
- Every `/api` route sits behind `ratelimit.go` (tronc `Ratelimit`, keyed on RemoteAddr after
  RealIP; service keys get their own bucket). Production needs `TRUSTED_PROXIES=private,cloudflare`.
- Synapse runs no push rules for appservice users, so Matrix pushers never fire for Prisme
  users. Android push comes from the API's appservice transactions (`modules/push`).
- The appservice token is root over every account and, since `deploy/synapse/render.py` makes
  `@prisme` a server admin, a Synapse admin credential too. It lives in Dokploy env only.
- A Registre erasure goes through `internal/erasure`: Synapse's admin API first, then one
  transaction. Answer a deletion notice with success only when the person was erased or provably
  never here; `docs/erasure.md` lists what remains and why. Matrix ids are never reused.
- Every compose variable is forwarded explicitly in `docker-compose.yml`.
- Deployment shape: one app container (Go API + SPA) behind Traefik, Synapse as a second
  service owning `/_matrix` and `/_synapse/client` with a higher router priority. Read the
  header of `docker-compose.yml` before changing routing.
- No inline comments and no comments inside function bodies; filet limits are in `filet.yml`.
