# Changelog

## [Unreleased]

## [0.10.0] - 2026-11-07

### Added

- The interface comes in French, English, Spanish and German, in the locale of the person's Registre profile (`GET /api/session` returns `user.locale`). Dates, times and file sizes follow it through `Intl`. Spanish and German need native review
- The sign-in screen is muse's `AuthScreen` with `SignInWithFacile` and the language switch
- GDPR deletion: Registre's signed `POST /api/registre/deletions` and a poll of its deletion list erase a person's Matrix account and agents in Synapse, their sessions and their porte user. Spaces they created stay with no creator. The state lives in the new `directory_deletion_state` table
- `GET /version`, served from the `APP_VERSION` build argument
- `LIVEKIT_URL` (optional) names the LiveKit origin for the CSP `connect-src`
- CI: build, tests against Postgres, svelte-check, client tests and build, and filet
- Security headers from tronc: `Permissions-Policy` allows camera, microphone, screen capture and speaker selection for calls. The full CSP stays report-only and allows the homeserver, LiveKit, blob media and the E2EE worker

### Changed

- porte v0.8.0 to v0.12.0, porte/directory v0.1.0 to v0.3.0, tronc v0.16.0 to v0.17.1, muse v0.8.1 to v0.13.0
- porte now adds `porte_users.locale` and `porte_login_codes.provider` and `.subject` at boot (idempotent, in `pg.Schema`)
- No borders in the client: reactions, the composer, menus, panels and cards use muse fills

### Removed

- `SSO_ONLY` from the environment and the compose file: porte mounts no password route while `OIDC_ISSUER` is set

## [0.9.1] - 2026-11-03

### Fixed

- An agent no longer leaves its owner's direct message when Bulle restarts: `PUT /api/agents/{user_id}/rooms` only leaves spaces and channels, and puts the agent back into any DM its owner still has with it

## [0.9.0] - 2026-11-03

### Added

- A call message turns into "Call ended" with its duration once the last person leaves; Prisme hears it from Echo's call events at `POST /api/echo/events`
- A Messages page lists every conversation with avatars (agents show their blob), the last message, its time and the unread count, plus everyone you can message, searchable and filtered by people or agents
- `PRISME_PUBLIC_URL` sets where Echo sends call events; it defaults to `MATRIX_PUBLIC_URL`

### Changed

- Direct messages leave the sidebar: one Messages entry sits above Overview
- Starting a call in a room that already has one running joins it instead of posting a second call message
- `GET /api/users` returns up to 200 entries instead of 20

### Removed

- The "New message" dialog; the Messages page replaces it

## [0.8.1] - 2026-11-03

### Changed

- The send button uses Bulle's up arrow instead of the paper plane

## [0.8.0] - 2026-11-03

### Added

- `POST /api/agents` returns `owner_user_id`, the Matrix id of the agent's owner, so Bulle can let only the owner use the agent's tools

## [0.7.0] - 2026-11-03

### Changed

- The message input matches Bulle: a rounded pill with a round + button to attach files and a filled send button that dims until there is something to send
- Thread summaries read "last from Pixel" instead of "last from Pixel (@pixel:prisme.facile.studio)"

## [0.6.2] - 2026-11-03

### Changed

- Spaces use the share-circle icon

## [0.6.1] - 2026-11-03

### Fixed

- Mentions use the username: picking Pixel inserts `@pixel`, not `@Pixel (@pixel:prisme.facile.studio)`, which appeared when two members shared a display name
- Received mention pills, including older messages, show `@username` instead of the sender's text

## [0.6.0] - 2026-11-03

### Changed

- The rail entry "All channels" is now "Overview", with a home icon: it is the space's page and the phone home
- The overview page is redesigned: a space header with its channel and member counts, Invite and New channel buttons, then cards for your channels, channels you can join, and direct messages
- Space invites show as a card with a "Join the space" button at the top of the overview

## [0.5.0] - 2026-11-03

### Changed

- New logo and favicon: the Solar chat bubble with dots, in the sidebar, on the sign-in page, as the favicon and app icons, and on the social card
- The web manifest, page description and social card name Prisme instead of the template

### Fixed

- The mention menu, the user picker and the agent list show `@name` instead of the full Matrix id with `:prisme.facile.studio`

## [0.4.0] - 2026-11-03

### Added

- Screen sharing in calls, with a full-width uncropped tile for each shared screen and its audio
- Every call is end-to-end encrypted with LiveKit E2EE (media and in-call data). `POST /api/calls` returns `e2ee_key`, derived per room from `PRISME_SECRET_KEY` and given only to room members; Echo and LiveKit never see it
- A lock in the call header shows the call is encrypted; a browser that cannot encrypt is refused with an explanation

### Changed

- The call message no longer carries the Echo link in its body

## [0.3.0] - 2026-11-03

### Added

- Mention suggestions on `@` with keyboard and pointer, live highlighting in the composer, highlighted mentions in messages and a tint on rows that mention you
- Live typing indicators with avatars
- Agents drawn as Facile blobs, in the client and as their Matrix avatar (PNG rendered by the API)
- Registre shared spaces mirrored into Prisme spaces with `#general`, members joining and leaving with the directory (`POST /api/registre/sync`)
- Agents join their owner's new spaces and channels immediately
- `GET /api/users/agents`; `GET /api/users` carries `avatar_url` and `agent`
- `/settings` page with profile, appearance, notifications, agents and sign out

### Changed

- Client rebuilt on the modular chat client and muse 0.8.0 `SideBar`, `Topbar` and `MobileNav`
- Agents are named after their name in lowercase per owner (`@pixel`), `<name>-<owner>` when taken
- Message avatars and names update live from member state

### Fixed

- Legacy `@agent-<owner>-<name>` agents migrate to `@<name>` and the old account is deactivated

## [0.2.0] - 2026-11-03

### Changed

- feat(client): redesign app shell with muse SideBar, space switcher, and collapsible rail
- feat(client): adopt muse components across chat, thread, member panels and dialogs
- style(client): remove all decorative borders in accordance with Facile design charte

## [0.1.0] - 2026-11-03

### Added

- feat(client): full Svelte 5 chat web UI with muse design tokens, spaces, channels, DMs, threads, reactions, calls and agent management
- feat(api): real-time agent name availability check endpoint `GET /api/agents/check`

### Changed

- feat(api): simplified agent Matrix username to `@<name>` with global uniqueness enforcement across human accounts and active agents
