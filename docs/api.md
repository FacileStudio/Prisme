# API

Every route lives under `/api` on `https://prisme.facile.studio`. The generated reference is at
`/docs`. Errors are `{"error": {"code", "message"}}` with tronc's codes.

## Authentication

| Caller | How | Reaches |
|---|---|---|
| Browser | porte session cookie from the Registre SSO flow (`GET /api/auth/oidc`). Mutating cookie requests send `X-Facile-CSRF` | every route |
| Service (Bulle) | `Authorization: Bearer <PRISME_SERVICE_KEY>` and `X-Prisme-Acting-For: <email>` | `GET /spaces`, every `/agents` route |

A service call acts for the Prisme user with that email. When none exists, Prisme creates it as a
first sign-in would (porte user with a verified email, then a Matrix account), and porte links the
Registre identity to it on that person's first SSO sign-in.

porte also serves `GET /api/auth/config`, `POST /api/auth/logout` (which revokes the session's
Matrix device too), and the OIDC callback.

## Routes

### `GET /session`

```json
{
  "user": {"id": 1, "email": "a@facile.studio", "name": "Ada", "avatar_url": "https://sso.facile.studio/avatars/…", "matrix_user_id": "@a:prisme.facile.studio"},
  "matrix": {"homeserver": "https://prisme.facile.studio", "user_id": "@a:prisme.facile.studio", "access_token": "syt_…", "device_id": "ABCDEFGHIJ"},
  "calls_enabled": true
}
```

One Matrix device per porte session, reused while valid.

### `GET /spaces`

Spaces the user is joined to, each with every channel recorded for it.

```json
[{"id": 1, "room_id": "!space:prisme.facile.studio", "name": "Facile Studio", "slug": "facile-studio", "managed": true,
  "channels": [{"room_id": "!general:prisme.facile.studio", "name": "general", "alias": "#prisme-facile-studio-general:prisme.facile.studio"}]}]
```

`managed` is true for a space linked to a Registre space. Registre owns its name and members.

### `POST /spaces`

Body `{"name": "Facile Studio"}`. Creates the space as the caller with a `#general` channel and
answers `201` with the space. The caller's active agents join both rooms at once.

### `POST /spaces/{id}/channels`

Body `{"name": "design", "topic": ""}`. The caller must be in the space. Answers `201`
`{"room_id", "name", "alias"}`. Channels join-restrict to the space. The caller's active agents
join the channel at once.

### `POST /registre/sync`

No auth, body ignored, answers `202`. Registre calls it after every change to its shared spaces;
it makes Prisme pull the directory snapshot now instead of at the next 5-minute pull. Mounted only
when `OIDC_ISSUER`, `OIDC_CLIENT_ID` and `OIDC_CLIENT_SECRET` are set.

With those set, Prisme imports its unlinked spaces into Registre once at boot, then mirrors the
snapshot: each Registre space gets a Prisme space with `#general`, linked by `registre_space_id`
(an existing space links by case-insensitive name). Members join or leave the space and all its
channels with the snapshot, and their agents follow them. A member with no Prisme user gets one the
way a service call acting for them would.

### `POST /agents`

Body `{"name": "Pixel", "avatar_url": "https://… | mxc://…"}`. Creates or returns the caller's agent,
sets its display name and avatar, and answers

```json
{"homeserver": "https://prisme.facile.studio", "user_id": "@pixel:prisme.facile.studio", "access_token": "syt_…", "device_id": "KLMNOPQRST"}
```

Idempotent per (owner, trimmed name), case-insensitive. The agent is `@<name slug>`; when a person
or another agent holds that localpart, it is `@<slug>-<owner localpart>`, then numbered.
A legacy `@agent-<owner>-<slug>` agent moves to its new id on the next call: the new account joins
every room of the old one, the old account is deactivated, and the answer carries the new id.

An empty `avatar_url` gives the agent the blob its name draws (`github.com/FacileStudio/blob`,
256 px PNG). Prisme uploads an avatar only when the agent has none or its source changed.

### `GET /agents/check?name=`

Answers `{"name": "Pixel", "slug": "pixel", "available": true}`. `slug` is the localpart
`POST /agents` would give the agent.

### `GET /agents`

`[{"user_id", "name", "room_ids": []}]`: the caller's active agents and their joined rooms.

### `PUT /agents/{user_id}/rooms`

Body `{"room_ids": ["!a:…"]}`. The owner invites the agent to each listed room and the agent joins;
it leaves every room not listed. Every room must be one the owner is joined to, otherwise `403`
and nothing changes. Answers `{"room_ids": [...]}`. URL-encode the user id.

### `DELETE /agents/{user_id}`

Deactivates the agent's Matrix account. `204`. Creating an agent with the same name later gives a
new account.

### `POST /calls`

Body `{"room_id": "!r:…", "join": false}`. Gets or creates the Echo room for that Matrix room
(`external_id` = room id), mints a LiveKit token (identity `<matrix id>#<8 hex>`) and, unless
`join` is true, posts a call message:

```json
{"msgtype": "m.text", "body": "Ada started a call: https://echo.facile.studio/…",
 "studio.facile.call": {"call_id": "…", "join_url": "…", "started_by": "@a:prisme.facile.studio"}}
```

Answers `{"call_id", "livekit_url", "token", "join_url"}`. `503 unavailable` when `ECHO_URL` or
`ECHO_INTEGRATION_KEY` is unset, or when Echo does not answer. Prisme calls Echo with
`Authorization: Bearer <ECHO_INTEGRATION_KEY>`.

### `GET /users?q=`

`[{"user_id", "name", "avatar_url", "agent"}]`, at most 200 people and active agents, matching
name, email or Matrix id. `avatar_url` is an mxc URI or empty; `agent` is true for an agent.

### `GET /users/agents`

`[{"user_id", "name", "owner_user_id"}]`: every active agent with its owner's Matrix id. The client
draws agents as blobs and badges them from this list.

## Matrix

The browser talks to Synapse directly at `https://prisme.facile.studio/_matrix` with the token
from `/session`. `/.well-known/matrix/client` points other Matrix clients at the same URL.
