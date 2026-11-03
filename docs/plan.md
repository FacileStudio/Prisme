# Plan: Prisme, the Facile suite chat app

## Goal

Prisme is the suite's team chat: Slack-shaped, Matrix underneath. Every Facile project is a
space, every space has channels, and agents (Bulle assistants) sit in those channels as
teammates you `@mention`. Calls start from any channel or DM and run on Echo.

It lives at `https://prisme.facile.studio`, repo `FacileStudio/Prisme`, deployed by Dokploy
(`gare.facile.studio`, project Facile, environment production) like every suite app.

Nothing is configured by hand to connect an agent: Bulle asks Prisme for an agent account and
gets one, already in the right rooms.

## Architecture

```
browser (SvelteKit + muse + matrix-js-sdk)
   │  /api/*  (Prisme API, porte session)        │  /_matrix/*  (client-server API)
   ▼                                              ▼
Prisme API (Go) ── application service ──▶ Synapse (federation off)
   │  provisions users, spaces, rooms, agents        │
   ├── Echo integration API (calls)                  Postgres (synapse db)
   └── Postgres (prisme db)
```

- **Homeserver: Synapse**, Postgres-backed, federation disabled, registration disabled, served
  under `prisme.facile.studio/_matrix` by the same Traefik router as the app. Synapse is chosen for
  its mature application-service and admin APIs, which the whole provisioning story rests on.
- **Prisme API is a Matrix application service** with an exclusive user namespace
  (`@[a-z0-9._=-]+:prisme.facile.studio` for humans and agents) and alias namespace
  (`#prisme-.*`). It creates accounts with `POST /register` (`m.login.application_service`) and
  mints access tokens with `POST /login` (`m.login.application_service`). No password exists for
  any account.
- **Sign-in is porte/Registre SSO only.** On first sign-in the API creates the Matrix user for the
  porte subject (localpart from the email's local part, deduplicated), stores the mapping, and
  hands the browser a Matrix access token for that user. The browser then talks Matrix directly
  through matrix-js-sdk. Logout revokes the token.
- **Channels are unencrypted in v1** (like Slack: search, agents and admin tooling read them).
  E2EE for DMs is a later step (matrix-js-sdk rust crypto is ready for it).

## Product (v1)

Everything you expect from a Facile app and a modern team chat:

- Spaces sidebar (Matrix spaces), channels list per space, DMs, unread and mention counters.
- Create, rename, archive spaces and channels; invite members of the suite (directory search
  over Prisme users); join public channels of a space.
- Messages: markdown, code blocks, edits, deletes, replies, **threads** (`m.thread`), reactions,
  `@mentions` with completion (humans and agents), file and image upload and preview (media repo),
  typing indicators, read receipts, link previews off by default.
- Search within a space (server-side via the Synapse search API).
- Notifications: unread badges, browser notifications on mention or DM, per-channel mute.
- Profile: name and avatar from Registre's `picture` claim, presence.
- Agents: a badge next to agent members, an "Add an agent" entry that lists the user's Bulle
  agents, and threads where agents answer.
- **Calls (Echo)**: a call button in every channel and DM. It asks Prisme API for a call; Prisme
  asks Echo to create or reuse the room for that channel and returns a join token; the call opens
  in a panel inside Prisme (LiveKit client against Echo's LiveKit, or Echo's embed route). A call
  message (`studio.facile.call` content in an `m.room.message` with a fallback body and link) is
  posted to the channel and shows who is in the call.
- UI: `@facile/muse` components and tokens only, Svelte 5 + SvelteKit, light and dark, mobile
  layout. Same chassis as the other suite apps (GoSvelteBoilerplate).

## Prisme API (contract)

Auth: porte session cookie for the browser; for service callers (Bulle), a Registre token-exchange
JWT with audience `prisme` and scopes `prisme.read prisme.agents`, acting for the user.

| Route | Purpose |
|---|---|
| `GET /api/session` | The signed-in user and a Matrix login: `{user, matrix: {homeserver, user_id, access_token, device_id}}` |
| `GET /api/spaces` | The user's spaces: `[{id, room_id, name, slug, channels: [{room_id, name, alias}]}]` |
| `POST /api/spaces`, `POST /api/spaces/{id}/channels` | Create a space (with a `#general` channel) or a channel, provisioned through the appservice |
| `POST /api/agents` | Body `{name, avatar_url?}`. Creates or returns the caller's agent `@agent-<owner>-<name>`, sets its display name and avatar, returns `{homeserver, user_id, access_token, device_id}`. Idempotent per (owner, name) |
| `PUT /api/agents/{user_id}/rooms` | Body `{room_ids: []}`. Invites the agent to each room and joins it (appservice masquerade), leaves rooms not listed. Only rooms the owner is a member of |
| `DELETE /api/agents/{user_id}` | Deactivates the agent |
| `POST /api/calls` | Body `{room_id}`. Gets or creates the Echo room for that Matrix room, posts the call message, returns `{call_id, livekit_url, token, join_url}` |

## Bulle side (interconnection)

- Bulle's Matrix bot settings gain a **"Use Prisme"** mode, on by default when the user has no
  hand-made bot: Bulle exchanges the user's token for a `prisme` audience token, calls
  `POST /api/agents` with the assistant's name and avatar, stores the returned access token (sealed,
  token auth, no password), and calls `PUT /api/agents/{id}/rooms` with the `#general` room of each
  linked space.
- Bulle spaces link to Prisme spaces by picking from `GET /api/spaces` in the Matrix settings, so
  the room id field is a picker, not something to copy from room settings.
- The bot behaviour is the one already planned in Bulle's `docs/plan-matrix-spaces.md`: answer only
  when mentioned, in a thread, one Bulle session per thread, work in the linked space.

## Echo side (integration)

Echo gains a server-to-server integration API so another suite app can run calls inside itself:

- `POST /api/integrations/rooms` `{external_id, title}` → `{room_id, join_url}`, idempotent per
  `(app, external_id)`.
- `POST /api/integrations/rooms/{id}/tokens` `{identity, name, avatar_url?}` → `{livekit_url, token}`.
- `GET /api/integrations/rooms/{id}` → participants, for the call message.
- Auth: a per-app integration secret (`ECHO_INTEGRATION_KEYS=prisme:<secret>`), or a Registre
  token with audience `echo`.
- `/embed/{room}` route without Echo chrome, allowed in frames from the apps listed in
  `ECHO_EMBED_ORIGINS` (CSP `frame-ancestors`), and LiveKit webhooks forwarded to the app's
  `callback_url` so Prisme can update the call message when a call ends.

## Deploy

- `docker-compose.yml`: `app` (Prisme, Go API + static client), `synapse`, `postgres` (two
  databases), Traefik labels routing `prisme.facile.studio` `/_matrix` and `/_synapse/client` to
  Synapse and everything else to the app; healthchecks with a short `start_period` (Field lesson).
- Secrets via Dokploy env: `DATABASE_URL`, `SYNAPSE_*`, `PRISME_AS_TOKEN`, `PRISME_HS_TOKEN`,
  porte OIDC client, `ECHO_URL`, `ECHO_INTEGRATION_KEY`. Registre seed gains the `prisme` client
  (redirect `https://prisme.facile.studio/api/auth/oidc/callback`) and the `prisme` audience for
  token exchange with Bulle.
- DNS: `prisme.facile.studio` exists on Cloudflare.

## Tracks (parallel)

| Track | Repo | Work | Exit |
|---|---|---|---|
| P | Prisme | repo from GoSvelteBoilerplate, Synapse + appservice provisioning, API above, SvelteKit muse client, compose, Registre client, Dokploy deploy | gates green, `https://prisme.facile.studio` signs in with SSO, two users chat in a channel, threads and mentions work |
| E | Echo | integration API, embed route, frame-ancestors, webhooks forward | tests green, a token from the API joins a LiveKit room, deployed |
| B | Bulle | "Use Prisme" agent provisioning and the space room picker on top of the server Matrix bot | tests green, an agent appears in a Prisme channel and answers a mention in a thread |

Then: Prisme calls wired to Echo, end-to-end test with two browsers (Playwright, fake media), ship
Prisme, Echo and Bulle.

## Risks

- Registre deploys production on push to main; its seed change must be correct the first time.
- Synapse memory on the shared server; cap caches.
- Appservice tokens are root over the namespace: they live only in Dokploy env and never reach the
  browser.
- Another session may be editing Bulle at the same time. Stage explicit paths only.

## Decisions taken while building track P

- **Bulle reaches Prisme with a service key, not token exchange.** At the suite lead's request,
  Bulle sends `Authorization: Bearer <PRISME_SERVICE_KEY>` and `X-Prisme-Acting-For: <email>` on
  `GET /api/spaces` and every `/api/agents` route. Prisme creates the user for an unknown email,
  so an agent works before its owner ever opens Prisme. Registre carries only the `prisme` OIDC
  client; no `prisme` exchange audience was added.
- **Additive routes.** `GET /api/agents` (the "Add an agent" picker), `GET /api/users?q=` (Synapse's
  user directory never lists appservice users), `calls_enabled` in `GET /api/session`, and
  `join: true` on `POST /api/calls` to join without posting a second call message.
- **Echo auth.** Prisme calls Echo's integration API with `Authorization: Bearer
  <ECHO_INTEGRATION_KEY>` and passes LiveKit identities as `<matrix id>#<8 hex>` so one person on
  two devices does not kick themselves out.
- **Tokens.** One Matrix device per porte session and one per agent, stored sealed with
  `PRISME_SECRET_KEY` and reused while `/account/whoami` accepts them.
- **Rooms.** Room version 12 (Synapse default): creators are not listed in power levels; the
  bot `@prisme` is the listed admin of every space and channel. Channels are `restricted` to
  their space.
- **Archive** of spaces and channels is not in v1; rename goes through Matrix directly.
