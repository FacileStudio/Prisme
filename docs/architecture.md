# Architecture

```
browser (SvelteKit + muse + matrix-js-sdk)
   │ /api/*  porte session                │ /_matrix/*, /_synapse/client/*
   ▼                                       ▼
api container (Go + SPA) ── appservice ──▶ synapse container (federation off)
   │                                       │
   └──────────── db container: Postgres 16, databases prisme and synapse
   │
   └── Echo integration API (calls)
```

## Request flows

- **Sign-in.** porte runs the OIDC code flow against Registre and sets the session cookie. The
  client calls `GET /api/session`; the API ensures the Matrix account (`POST /register` with
  `inhibit_login`), reuses or mints a device (`POST /login`, `m.login.application_service`),
  stores its token sealed, and returns it. The client then syncs with Synapse directly.
- **Logout.** A middleware in front of porte's `POST /api/auth/logout` revokes the session's
  Matrix device, then porte ends the session.
- **Spaces and channels.** Created as the user through identity assertion, with the bot
  `@prisme:prisme.facile.studio` invited and joined as admin. The bot writes `m.space.child`.
  Channels are `restricted` to their space.
- **Agents.** `@<name>` accounts, created and placed by the API; the owner
  invites, the agent joins. Bulle reaches these routes with the service key.
- **Calls.** The API asks Echo for the room bound to the Matrix room id, mints a LiveKit token and
  posts a `studio.facile.call` message.

## Data

`porte_*` tables hold identity and sessions. Prisme adds `accounts`, `matrix_sessions`,
`spaces`, `channels` and `agents` (`apps/api/schemas/models.go`). Everything else lives in
Synapse.

The schema of the `prisme` database is ordered goose SQL in `apps/api/migrations/`, embedded in
the binary and applied at boot by `tronc/migrate` behind a Postgres advisory lock. GORM is the
query layer only: `AutoMigrate` never runs, and `scripts/check.sh` fails on an `.AutoMigrate(`
outside a test file. A schema change is a new numbered file with `-- +goose Up` and
`-- +goose Down`, plus the matching struct field in `schemas/`;
`TestGORMModelsMatchTheMigratedSchema` fails when a field has no column. porte's tables are in
the migrations too, and `TestPorteSchemaIsInTheMigrations` fails when a porte bump changes
`pg.Schema`. `docker run <image> migrate status` (also `up`, `version`, `baseline <n>`) answers
without a shell. The `synapse` database belongs to Synapse and is not migrated by the API.
