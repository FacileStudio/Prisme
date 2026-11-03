# Changelog

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
