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
[{"id": 1, "room_id": "!space:prisme.facile.studio", "name": "Facile Studio", "slug": "facile-studio",
  "channels": [{"room_id": "!general:prisme.facile.studio", "name": "general", "alias": "#prisme-facile-studio-general:prisme.facile.studio"}]}]
```

### `POST /spaces`

Body `{"name": "Facile Studio"}`. Creates the space as the caller with a `#general` channel and
answers `201` with the space.

### `POST /spaces/{id}/channels`

Body `{"name": "design", "topic": ""}`. The caller must be in the space. Answers `201`
`{"room_id", "name", "alias"}`. Channels join-restrict to the space.

### `POST /agents`

Body `{"name": "Bulle", "avatar_url": "https://… | mxc://…"}`. Creates or returns the caller's agent
`@agent-<owner localpart>-<name slug>`, sets its display name and avatar, and answers

```json
{"homeserver": "https://prisme.facile.studio", "user_id": "@agent-ada-bulle:prisme.facile.studio", "access_token": "syt_…", "device_id": "KLMNOPQRST"}
```

Idempotent per (owner, name). The token is stored sealed and returned again on later calls; a new
one is minted only when the old one was revoked.

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

`[{"user_id", "name"}]`, at most 20, matching name, email or Matrix id.

## Matrix

The browser talks to Synapse directly at `https://prisme.facile.studio/_matrix` with the token
from `/session`. `/.well-known/matrix/client` points other Matrix clients at the same URL.
