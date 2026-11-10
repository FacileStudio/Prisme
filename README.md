# Prisme

Team chat for the Facile suite, at https://prisme.facile.studio. It looks like Slack and runs on
Matrix: every project is a space, every space has channels, and Bulle agents sit in channels as
teammates you `@mention`. Calls start from any channel or DM and run on Echo.

## What it does

- Signs people in through Registre (porte SSO only) and gives each person a Matrix account on
  first sign-in, named after their email
- Runs its own Synapse homeserver beside the app, federation and registration off
- Acts as a Matrix application service: it creates users, spaces, channels and agent accounts,
  and no account ever has a password
- Hands the browser a Matrix token for its session; the client then talks Matrix directly with
  matrix-js-sdk
- Gives Bulle zero-config agent accounts (`@<name>`) already joined to the right rooms
- Starts Echo calls from a room and posts a call message other members can join from

## Stack

| Layer | Tech |
|---|---|
| API | Go 1.26, chi, GORM + Postgres 16, [porte](https://github.com/FacileStudio/porte) v0.8, [tronc](https://github.com/FacileStudio/tronc) |
| Homeserver | Synapse v1.164.1 on Postgres, configured from the environment at start |
| Client | SvelteKit 2 (static), Svelte 5, `@facile/muse`, matrix-js-sdk, livekit-client |
| Deploy | Dokploy compose: `api` (Go API + SPA in one distroless image), `synapse`, `db` |

## Quick start

```sh
mise trust && mise install
docker compose -p prisme-dev -f docker-compose.dev.yml up -d --build   # Postgres + Synapse on localhost
mise run dev                                                           # API :4000, client :5173
mise run check                                                         # gofmt, vet, test, svelte-check
filet check .
```

The dev compose serves Synapse at `http://localhost:18008` with server name `localhost` and
throwaway secrets. Run the API with `MATRIX_INTERNAL_URL=http://127.0.0.1:18008`,
`MATRIX_PUBLIC_URL=http://localhost:18008`, `MATRIX_SERVER_NAME=localhost`,
`PRISME_AS_TOKEN=dev-as-token` and a 32-character `PRISME_SECRET_KEY`.

## Configuration

| Variable | Service | What it does |
|---|---|---|
| `POSTGRES_PASSWORD` | db, synapse, api | One Postgres cluster, two databases: `prisme` and `synapse` (C collation) |
| `SYNAPSE_SIGNING_KEY` | synapse | The homeserver's ed25519 key. Generated into the volume when unset |
| `SYNAPSE_MACAROON_SECRET`, `SYNAPSE_FORM_SECRET` | synapse | Synapse's own secrets |
| `PRISME_AS_TOKEN` | synapse, api | Appservice token: root over the user namespace. Never reaches a browser |
| `PRISME_HS_TOKEN` | synapse, api | Homeserver token: Synapse sends it on every appservice transaction, and the API refuses a transaction without it |
| `PRISME_AS_URL` | synapse | Where Synapse sends appservice transactions. Defaults to `http://prisme-api:4000`, the API's alias on the compose network; empty renders `null` |
| `FCM_SERVICE_ACCOUNT` | api | Firebase service account (raw or base64 JSON) for Android push. Prisme's Android app lives in the `courrier-de1ea` Firebase project, so this is Courrier's account. Empty turns Android push off |
| `OIDC_MOBILE_REDIRECT` | api | Where porte's app sign-in hands back its code. Defaults to `studio.facile.prisme://auth` while OIDC is on |
| `PRISME_SECRET_KEY` | api | Seals the Matrix tokens Prisme stores (AES-256-GCM), at least 32 characters |
| `PRISME_SERVICE_KEY` | api | Comma-separated keys for service callers (Bulle). Empty turns that path off |
| `MATRIX_SERVER_NAME`, `MATRIX_PUBLIC_URL` | synapse, api | `prisme.facile.studio` and its https URL |
| `OIDC_ISSUER`, `OIDC_CLIENT_ID`, `OIDC_CLIENT_SECRET`, `OIDC_REDIRECT_URL`, `OIDC_SUCCESS_URL` | api | porte SSO against Registre, client id `prisme` |
| `ECHO_URL`, `ECHO_INTEGRATION_KEY` | api | Echo's integration API. Either empty turns calls off with a clear message |
| `LIVEKIT_URL` | api | The LiveKit server Echo hands out, so the security headers can allow its `wss://` and `https://` origins in `connect-src`. Calls still work without it, but the CSP reports a violation |
| `PRISME_PUBLIC_URL` | api | Where Echo posts call events (`/api/echo/events`). Defaults to `MATRIX_PUBLIC_URL` |
| `TRUSTED_PROXIES`, `JOURNAL_URL`, `JOURNAL_TOKEN` | api | Suite defaults |

Every variable is forwarded explicitly in `docker-compose.yml`: one set in Dokploy but missing
there never reaches a container.

## API

Full reference at `/docs`; summary in [docs/api.md](docs/api.md).

| Route | Auth | Purpose |
|---|---|---|
| `GET /api/session` | session | The user, their locale and their Matrix login |
| `GET /api/spaces` | session or service | Spaces the user is in, with every channel |
| `POST /api/spaces`, `POST /api/spaces/{id}/channels` | session | Create a space (with `#general`) or a channel |
| `GET /api/agents`, `POST /api/agents` | session or service | List, or create-or-return, the user's agents |
| `PUT /api/agents/{user_id}/rooms`, `DELETE /api/agents/{user_id}` | session or service | Place or deactivate an agent |
| `POST /api/calls` | session | Start or join the Echo call of a room; returns the LiveKit token and the room's `e2ee_key` |
| `POST /api/echo/events` | Echo signature | Call events: the call message turns into "Call ended" when the last person leaves |
| `GET /api/users?q=` | session | People directory |
| `GET /api/push`, `POST /api/push/devices` | session | Whether Android push is on; register the app's Firebase token |
| `PUT /_matrix/app/v1/transactions/{txnId}` | hs_token | Synapse's appservice events, which drive Android push |
| `POST /api/registre/sync`, `POST /api/registre/deletions` | none, signed | Registre's space nudge, and its signed GDPR `user.deleted` notice |
| `GET /version` | none | The version stamped at build |

Service callers send `Authorization: Bearer <PRISME_SERVICE_KEY>` and
`X-Prisme-Acting-For: <email>`. Prisme creates the user for that email when it does not exist,
the same way a first sign-in would.

## Design notes

- **Appservice provisioning.** Users are created with `POST /register` and
  `m.login.application_service`, with `inhibit_login: true` so registering leaves no device;
  tokens come from `POST /login` with the same login type
  (https://spec.matrix.org/v1.16/application-service-api/). This follows MSC4190, which makes
  `inhibit_login` mandatory on homeservers that move to OAuth2 and blocks appservice `/login`
  there (https://github.com/matrix-org/matrix-spec-proposals/pull/4190). Prisme does not use
  MAS, so `/login` stays available; moving to MAS would mean switching browser logins to
  MSC4190 device creation and identity assertion.
- **One device per porte session.** `GET /api/session` reuses the stored device while
  `/account/whoami` accepts its token, and porte's logout revokes it through a small middleware.
  Agent tokens are reused the same way, so Bulle's stored token stays valid.
- **Deactivation without UIA.** Synapse skips user-interactive auth on
  `/account/deactivate` for appservice requests
  (https://github.com/element-hq/synapse/blob/develop/synapse/rest/client/account.py).
- **Room version 12.** Synapse 1.164 creates v12 rooms, where the creator has unlimited power and
  must not appear in `power_level_content_override.users`. The Prisme bot (`@prisme`, the
  appservice sender) is the only listed admin, so it can link channels into any space.
- **Directory.** Synapse leaves every user of an exclusive appservice namespace out of its user
  directory (`should_include_local_user_in_dir` in
  https://github.com/element-hq/synapse/blob/v1.164.1/synapse/storage/databases/main/user_directory.py),
  and every Prisme human is one. `GET /api/users` answers search instead.
- **Channels.** Rooms with join rule `restricted` to their space, history `shared`, unencrypted in
  v1 so search and agents can read them.
- **Media.** Authenticated media stays on (the Synapse default), so the client fetches media with
  its token and shows blob URLs.
- **Resources.** Synapse runs on a shared host: `caches.global_factor: 0.5`, cache autotuning
  capped at 256M (jemalloc ships in the image), `event_cache_size: 5K`, presence on, URL
  previews off
  (https://element-hq.github.io/synapse/latest/usage/configuration/config_documentation.html#caches).
  `rc_message` is raised to 5/s because browser users post with their own tokens.
- **Matrix client.** The browser builds matrix-js-sdk with `createClient({baseUrl, accessToken,
  userId, deviceId})` from `/api/session` and starts it with thread support
  (https://matrix-org.github.io/matrix-js-sdk/interfaces/matrix.ICreateClientOpts.html).
- **Call encryption.** Every call is end-to-end encrypted, always: media, screen shares and
  in-call data. The key is per Matrix room, `base64url(HMAC-SHA256(SHA-256("prisme-call-e2ee:" +
  PRISME_SECRET_KEY), "prisme-call-e2ee:" + room_id))`, so it needs no table and every member gets
  the same one. The API derives it and returns it from `POST /api/calls` only after the
  membership check; it never sends it to Echo and never logs it. Echo and LiveKit relay frames
  they cannot decrypt, so captions, recording and Echo's own join link cannot work on a Prisme
  call. Rotating `PRISME_SECRET_KEY` changes every room key. A browser without encoded-frame
  transforms (`isE2EESupported()` false) is refused rather than joined in plaintext.
- **Installed apps.** `apps/client` is also the Android app (Capacitor 8, `studio.facile.prisme`)
  and `apps/desktop` wraps the same build in Tauri 2. Both sign in with porte's `?flow=mobile`
  PKCE flow through the `studio.facile.prisme://auth` deep link and call the API with a bearer
  token, so CORS always allows `https://localhost`, `tauri://localhost` and
  `http://tauri.localhost` while OIDC is on. First launch shows a short intro that ends on the
  sign-in slide, then a welcome step that asks for notifications.
- **Notifications.** The browser and the desktop app notify from their own Matrix sync while
  they are out of sight. The Android app cannot, so the server pushes: Synapse never runs push
  rules for users an exclusive appservice owns (`get_if_app_services_interested_in_user` in
  https://github.com/element-hq/synapse/blob/v1.164.1/synapse/push/bulk_push_rule_evaluator.py),
  and every Prisme user is one, so Matrix pushers never fire. Instead Synapse sends every event to
  the appservice (`PUT /_matrix/app/v1/transactions/{txnId}`), and the API pushes a new message
  to the other person of a two-member room and to the users its `m.mentions` names, skipping
  rooms they muted, through Firebase HTTP v1. Device tokens live in `push_devices`, tied to the
  porte session, so signing out removes them. Events older than five minutes are never pushed.
- **Releases of the apps.** Publishing a GitHub release runs `android.yml` (signed
  `Prisme-android.apk`) and `desktop.yml` (Windows, macOS and Linux installers plus the
  updater's `latest.json`). The app version is `apps/client/package.json`, checked against the
  tag. Signing keys live in tiroir as `PRISME_ANDROID_*` and `PRISME_TAURI_SIGNING_*`; neither
  can ever be rotated, or installed apps stop updating.
- **Bulle access.** The plan's Registre token exchange (`prisme.read prisme.agents`) was replaced
  by the service key, at the suite lead's request, so a Bulle user's agent works before they ever
  open Prisme. porte's `OIDC_MACHINE_AUDIENCE` path still works if it is ever configured.

## Structure

```
apps/api/      Go API: modules/ (auth, session, spaces, agents, calls, users, push), internal/ (matrix, accounts, echo, seal)
apps/client/   SvelteKit client, and the Capacitor Android project in android/
apps/desktop/  Tauri desktop shell around the client build
deploy/        synapse/ (templates and render.py), postgres/ (creates the synapse database)
docs/          plan, api, configuration, architecture
```

---

Part of the [Facile Suite](https://facile.studio): self-hosted tools for creative studios and
freelancers. One login, zero cloud dependency.

## License

Prisme is licensed under the GNU Affero General Public License v3.0 only (`AGPL-3.0-only`); see [LICENSE](LICENSE). A commercial license is available from Facile Studio for use that cannot meet the AGPL's terms.
