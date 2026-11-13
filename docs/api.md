# API

Every route lives under `/api` on `https://prisme.facile.studio`. The generated reference is at
`/docs`. Errors are `{"error": {"code", "message"}}` with tronc's codes.

## Authentication

| Caller | How | Reaches |
|---|---|---|
| Browser | porte session cookie from the Registre SSO flow (`GET /api/auth/oidc`). Mutating cookie requests send `X-Facile-CSRF` | every route |
| Service (Bulle) | `Authorization: Bearer <PRISME_SERVICE_KEY>` and `X-Prisme-Acting-For: <email>` | `GET /spaces`, every `/agents` route |

### Rate limits (since v0.17.0)

Every `/api` route spends from a token bucket. Past it, Prisme answers `429` with the
`rate_limited` code and a `Retry-After` in seconds.

| Bucket | Key | Budget |
|---|---|---|
| Every request | the signed-in user, else the visitor's address | 240 per minute, burst 120 |
| Writes (anything but `GET`, `HEAD`, `OPTIONS`), on top | the signed-in user, else the visitor's address | 30 per minute, burst 30 |
| Sign-in (`/api/auth/*`), on top | the visitor's address | at most 30 in any minute (burst 15) |
| Service key | which key in `PRISME_SERVICE_KEY` | 1200 per minute, burst 600, instead of the visitor buckets |

Since v0.19.0 a request whose session cookie or bearer token resolves to a live session spends from
that user's buckets, wherever it comes from, so people behind one office address no longer share
a budget. A request with no session, or one that does not resolve, spends from its address's
buckets. porte's own routes (`/api/auth/*`) always spend from the address's buckets and never
look the session up. Each limiter remembers at most 50,000 keys and forgets the longest idle
first.

The visitor's address is what tronc's real-IP middleware resolved from `TRUSTED_PROXIES`, never a
header read by hand. Production runs `TRUSTED_PROXIES=private,cloudflare`, so it is
`Cf-Connecting-Ip` read from a trusted chain. Without that setting, every visitor shares the
Cloudflare edge's bucket. Synapse's appservice transactions (`/_matrix/app/...`), `/health`,
`/ready` and `/version` are not limited.

A service call acts for the Prisme user with that email. When none exists, Prisme creates it as a
first sign-in would (porte user with a verified email, then a Matrix account), and porte links the
Registre identity to it on that person's first SSO sign-in.

### `X-Prisme-Acting-Subject` (since v0.14.0)

A service call may also send `X-Prisme-Acting-Subject: <Registre subject>`, the person's `sub`
under Prisme's own `OIDC_ISSUER`. Prisme then links that subject to the account for the email,
so a Registre deletion notice can find the person even if they never signed in to Prisme
(see [erasure](erasure.md)). Without the header, a call behaves exactly as before.

Prisme reads the header only on a call authenticated by the service key. A browser session or a
delegated bearer token that sends it is unaffected, and nothing is validated or written.

| The account for the email | The subject | Answer |
|---|---|---|
| does not exist | is nobody's | Created, subject linked, then the route runs |
| holds no subject of the issuer | is nobody's | Subject linked, then the route runs |
| holds this subject (either spelling of the issuer) | | Left alone, the route runs |
| holds another subject of the issuer | | `409 acting_subject_mismatch` |
| exists or not | belongs to another account | `409 acting_subject_taken` |
| exists or not | was erased by Prisme (since v0.14.2) | `409 acting_subject_erased` |

Refusals write nothing: no account is created and no identity changes. A service key can never
move or merge an identity. Prisme links the subject through porte's identity store, the same
`Save` a sign-in calls, with the issuer spelled as `OIDC_ISSUER` is.

Every refusal uses the error body above (`{"error": {"code", "message"}}`), with these codes:

| Status | `error.code` | Meaning | Retry? |
|---|---|---|---|
| 400 | `acting_subject_invalid` | The header is empty, longer than 64 bytes (Registre's limit), not UTF-8, holds a control character, or is sent twice | No, fix the value |
| 409 | `acting_subject_mismatch` | The account for this email holds a different subject of the issuer | No, the stored data conflicts |
| 409 | `acting_subject_taken` | Another account holds this subject | No, the stored data conflicts |
| 409 | `acting_subject_erased` | Prisme erased this subject after a Registre deletion notice; it never holds an account again | No, never: Registre does not reissue a subject |
| 412 | `acting_subject_unsupported` | This Prisme has no `OIDC_ISSUER`, so there is nothing to file a subject under | No |
| 500 | `internal` | The database failed | Yes |

A 409 is final for the data as stored: Prisme reads the rows before answering it, and a store
refusal the rows do not explain is a 500.

**A changed address is refused, never merged.** If the person changed their email in Registre,
Bulle sends the new address with the old subject. The subject belongs to the account under the
old address, so the call answers `409 acting_subject_taken` and no account is created for the new
address. The person's next Prisme sign-in finds the account by subject and porte rewrites its
verified email; Bulle's next call then succeeds. Prisme logs each link at info and each refusal at
warn, with the account id and the service key's position in `PRISME_SERVICE_KEY`, never the
email or the subject.

porte also serves `GET /api/auth/config`, `POST /api/auth/logout` (which revokes the session's
Matrix device too), and the OIDC callback.

## Routes

### `GET /session`

```json
{
  "user": {"id": 1, "email": "a@facile.studio", "name": "Ada", "avatar_url": "https://sso.facile.studio/avatars/…", "matrix_user_id": "@a:prisme.facile.studio", "locale": "fr"},
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

### `POST /registre/deletions`

No session: Registre posts a signed `user.deleted` security event, and Prisme verifies the RS256
signature, issuer, audience and time window itself (a notice that does not verify is a `400`).
Prisme also polls Registre's deletion list with the same client credentials as the snapshot sync, so
a missed push is caught. For each erased person it settles their rooms, deletes their uploads and
deactivates their Matrix account and their agents' with erasure, through Synapse's admin API, then
deletes their account, agents, sessions and porte user in one transaction. It answers success only
when it erased the person or can prove they were never here; otherwise it answers an error and
Registre retries. A replayed notice is answered as done. [Erasure](erasure.md) has the rules and
what remains. Set it as the client's `deletion_uri` in Registre. Mounted with `POST /registre/sync`.

### `GET /me/invitations`

Session only. The signed-in person's pending Registre space invitations, for the user menu:
`{"count": n, "href": url}`, where `href` is the one invitation's page when `n` is 1 and Registre's
spaces page otherwise. Any failure answers `{"count": 0}` with no `href`. porte/kit caches each
answer for 60 seconds. Mounted with `POST /registre/sync`.

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

Since v0.17.0 the list holds only the caller, the people and agents joined to a space the caller
is joined to, and the caller's own agents. Nobody else is found, by name, email or Matrix id.
Prisme reads the caller's spaces and their members from Synapse on every call.

### `GET /users/agents`

`[{"user_id", "name", "owner_user_id"}]`: the active agents the caller may see, with the owner's
Matrix id. Since v0.17.0 that is the caller's own agents, and the agents joined to a shared space
whose owner is a space peer. The client draws agents as blobs and badges them from this list.

### `GET /push`

`{"native"}`: true when this server can push to the Android app (`FCM_SERVICE_ACCOUNT` and
`PRISME_HS_TOKEN` are set). The Android app only registers with Firebase when it is true.

### `POST /push/devices`

`{"token"}`, the Android app's Firebase registration token. Answers `204`. The row belongs to the
calling porte session and goes when it signs out. `503 unavailable` while push is off.

## Matrix

The browser talks to Synapse directly at `https://prisme.facile.studio/_matrix` with the token
from `/session`. `/.well-known/matrix/client` points other Matrix clients at the same URL.

Synapse calls the API as its appservice at `PUT /_matrix/app/v1/transactions/{txnId}` with
`Authorization: Bearer <PRISME_HS_TOKEN>`, over the compose network only: Traefik sends every
public `/_matrix` path to Synapse. Every other `/_matrix` path on the API answers `404`, so
Synapse's user and room queries get a clear no.
