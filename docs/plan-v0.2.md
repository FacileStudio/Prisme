# Plan: Prisme v0.2, an integrated Facile chat

## Goal

Prisme reads as a Facile app (muse shell, settings page, Registre spaces) and as a real chat
(mention suggestions and highlighting, live typing, live messages with the right faces). Bulle
agents appear in every space of their owner with no setup, under their own name.

## State on 2026-11-03

- `main` (v0.1.0) shipped a thin client. A richer modular client sat unmerged on `feat/client`.
  Commit `841ef81` on `feat/integrated-chat` ports it onto `main`; every gate is green.
- New agents get `@<slug>`, but `agents.find(slug)` returns rows made before that change, so
  `@agent-yann-pixel` never becomes `@pixel`.
- Agent slugs are global: a second user whose assistant is also called "Bulle" gets a 409 on
  every Bulle provisioning pass.
- Prisme spaces are local. Registre owns the suite's shared spaces (`porte/directory`), and ten
  apps already sync them. Prisme does not.
- Bulle sends no avatar for agents, and deleting a Bulle job leaves its running runs alive.
- Blob avatars come from `@facile/blob` (Svelte) and `github.com/FacileStudio/blob` (Go).
  blobatar.dev is no longer used anywhere.

## API contract changes

| Route | Change |
|---|---|
| `GET /api/users?q=` | Items gain `avatar_url` (mxc or empty) and `agent` (bool). Agents are listed too |
| `GET /api/users/agents` | New. Every active agent: `[{user_id, name, owner_user_id}]`. The client uses it to draw blobs and agent badges |
| `POST /api/agents` | Migrates a legacy `@agent-<owner>-<slug>` row to `@<slug>`; falls back to `<slug>-<owner localpart>` when another owner holds the slug; sets a blob PNG avatar when `avatar_url` is empty |
| `POST /api/registre/sync` | New, unauthenticated nudge (body ignored) that triggers a directory pull |

## Tracks (parallel, disjoint files)

### A. API (`apps/api`)

1. Agent identity. Look agents up by `(owner_id, lower(name))`, not by slug alone. When the row's
   Matrix id is not `@<wanted localpart>`, register the new account, join it to every room the old
   one is in, deactivate the old account, and update the row. Localpart: slug, else
   `<slug>-<owner localpart>` when taken by another owner or a human. Test the choice function.
2. Blob avatar. When `avatar_url` is empty, render `blob.New(name, blob.Options{}).Pixels(256)` to
   PNG with `image/png`, upload it to the media repo and set it as the agent avatar. Only when the
   agent has no avatar or the name changed.
3. Registre spaces. Adopt `porte/directory` the way Sablier does
   (`apps/api/modules/spaces/registre_*.go`): import once, pull at boot and every 5 minutes, nudge
   route. Every Registre space gets a Prisme space with `#general`, linked by `registre_space_id`
   (existing spaces link by case-insensitive name). Members join and leave with the snapshot; a
   member without a Prisme account gets one (the same path Bulle's acting-for uses). Linked spaces
   refuse rename. Unconfigured (no OIDC) deploys keep local spaces.
4. Agents follow their owner. When a space or channel is created, or the sync adds an owner to a
   space, join the owner's agents to the new rooms right away instead of waiting for Bulle's
   10-minute pass.
5. `GET /api/users` and `GET /api/users/agents` as above, documented in `Documentation()`.

### B1. Client shell and settings (`components/shell`, `routes`, `dialogs/SettingsDialog`, `state/ui`, `state/theme`, `app.css`)

- Rebuild the shell on muse organisms like Courrier and Opus: `SideBar` (spaces through its
  `spaces` and `SpaceSwitcher`), `Topbar` and `MobileNav` on phones, `PageTransition`.
- A `/settings` route built from `SettingsSection` and `SettingsRow`: profile, appearance,
  notifications, agents, and Sign out at the bottom. Remove sign-out from everywhere else.
- Spaces managed by Registre link to `https://sso.facile.studio/spaces` instead of a create form.
- Bump `@facile/muse` to the newest tag. Tokens only, 44px targets, 360px wide, reduced motion.

### B2. Client conversation (`components/room`, `components/common`, `chat/*`, `matrix/*`, `api/prisme.ts`)

- Mentions: `@` opens a suggestion menu of room members and agents (keyboard and pointer),
  inserts a pill, and the composer highlights `@name` tokens while typing (overlay behind the
  textarea). Sent messages carry `m.mentions` and a `matrix.to` pill; received mentions render as
  highlighted chips, and the row is tinted when it mentions me.
- Typing: send typing notices (debounced, stop on send and blur), and show live "X is typing"
  with avatars from `RoomMember.typing` events.
- Live messages with correct faces: avatars resolve from the room member state at the event,
  update on `RoomMember.membership` and `User.avatarUrl`, and use authenticated media.
- One `Avatar` component: agents draw `<Blob name>` from `@facile/blob` (from
  `GET /api/users/agents`) with an agent badge; humans draw their Matrix avatar with initials
  fallback.

### C. Bulle (separate worktree off `origin/main`)

- Deleting a job (HTTP handler and the `delete_background_job` tool) stops its running runs
  through `Runner.Stop` first, through one shared function.
- Send the assistant's profile avatar to `Prisme.Agent` when set; else send nothing and let Prisme
  draw the blob.
- Replace every blobatar.dev URL with `@facile/blob` (client) or the Go `blob` package (server).

## Exit

`sh scripts/check.sh` and `filet check .` green in both repos; Prisme ships as v0.2.0 and Bulle as
a patch release; on prod, `@pixel:prisme.facile.studio` exists, sits in every space of its
owner, and `@agent-yann-pixel` is deactivated.

## Risks

- Another session has uncommitted Registre-spaces work in Bulle's main tree. Track C works in its
  own worktree and stages explicit paths.
- Registre's seed lacks `directory_sync_uri` for `prisme`; until it is added the 5-minute pull is
  the only trigger. Registre deploys on push, so that change goes through its owner session.
- Dokploy sometimes misses the push webhook: compare the latest deployment title after push.
