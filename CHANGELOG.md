# Changelog

## [0.14.2] - 2026-11-11

### Fixed

- An erased Registre subject never gets a Prisme account again. On 2026-11-11 a Bulle service
  call that arrived one second after an erasure created and linked a new account for the erased
  subject. Every path that creates or links an account from a subject now checks
  `erased_subjects` first: a service call with `X-Prisme-Acting-Subject` answers
  `409 acting_subject_erased` (final, nothing written), a sign-in through Registre is refused,
  and the Registre space sync skips the member. See `docs/api.md` and `docs/erasure.md`.
- An account re-created for an erased subject is erased instead of skipped. The next delivery of
  the deletion notice erases it (porte/directory's handled subjects that an identity holds again
  are handed back as unhandled, from the next start), and `erase-user <email>` no longer refuses
  it, because an identity whose subject Prisme erased does not count as a link.

### Security

- build(client): vite 7.3.8, and overrides for devalue 5.9.4, nanoid 3.3.20 and source-map-js
  1.2.2, which clear the four high advisories douane reported in `apps/client/bun.lock`.

## [0.14.1] - 2026-11-11

### Security

- build(api): Go 1.26.11 (`go` directive, `golang:1.26.11-alpine` builder, mise pins). It fixes GO-2026-6848 to GO-2026-6859 in the standard library and toolchain, which douane reported as unscored gaps and which failed the gate.

## [0.14.0] - 2026-11-10

### Added

- Service-key calls accept `X-Prisme-Acting-Subject: <Registre subject>` beside
  `X-Prisme-Acting-For`. Prisme links the subject to the account for that email through porte's
  identity store, the way a sign-in does, so a Registre deletion notice can erase someone Bulle
  provisioned who never signed in to Prisme. Before, such a notice stayed pending until an
  operator ran `erase-user`. A call without the header behaves exactly as before.
- The header never moves an identity. A subject that belongs to another account answers
  `409 acting_subject_taken`, an account holding another subject answers
  `409 acting_subject_mismatch`, a malformed value answers `400 acting_subject_invalid`, and a
  Prisme with no `OIDC_ISSUER` answers `412 acting_subject_unsupported`. Each refusal writes
  nothing. Session and delegated-token requests ignore the header. See `docs/api.md`.

### Changed

- CI pins douane v0.7.2 and filet v0.19.2.

## [0.13.0] - 2026-11-10

### Fixed

- A Registre erasure now erases the person everywhere. Before, a notice for an account that
  held no Registre subject (one Bulle created through the service key) erased nothing and was
  still answered as done, so the account, its agents and its Matrix data stayed. Prisme now
  answers a notice as done only when it erased the person or can prove they were never here,
  and leaves the notice pending otherwise. See `docs/erasure.md`.
- The Matrix half of an erasure goes through Synapse's admin API: rooms only the person was in
  are purged, a shared room they alone administered gets a new admin first, their uploads are
  deleted, and their accounts and their agents' accounts are deactivated with erasure.
- The database half runs in one transaction through `porte/kit`'s erase hook and deletes the
  account, identities, sessions, push devices, Matrix sessions and agents. A failure on either
  side leaves the notice pending, and the retry converges.
- A new person who signs in with the email of an erased one gets a new Matrix id (`alice2`),
  never the erased account. Prisme asks Synapse whether a taken id is deactivated, so the
  reservation rows older erasures kept in `accounts` and `agents` are deleted at boot.

### Added

- Pending Registre space invitations in the user menu, from `GET /api/me/invitations`
  (`porte/kit`'s handler, cached 60 seconds per person, behind the session).
- `/api erase-user <email>` erases an account that holds no Registre subject, which no
  deletion notice can reach.
- The Synapse image makes the appservice sender a server admin at start, and the API checks
  the admin credential with one read-only call at boot.

### Changed

- tronc v0.21.0, porte v0.18.2 (`porte/directory` v0.4.0) and `@facile/muse` v0.18.0.
- Go 1.26.10 (the standard library's critical advisories), `golang.org/x/text` v0.41.0 and
  `go-jose` v4.1.6. CI scans dependencies with douane v0.7.1 and fails on a critical finding.
- `filet.yml` is `preset: facile` (identical findings), CI runs filet v0.19.1, and the pre-push
  hook runs `filet gate .`.

- Licensed under the GNU Affero General Public License v3.0 only (`AGPL-3.0-only`). See `LICENSE`.

## [0.12.1] - 2026-11-07

### Changed

- `@facile/muse` v0.14.4 to v0.16.1. Secondary buttons use the `secondary` variant in place of the deprecated `outline`, and the call panel's icon buttons use `secondary` and `ghost-danger` in place of `default` and `danger`.
- `GET /version` comes from the version in `apps/client/package.json`, which the release commit bumps (still `v` plus the number). The `APP_VERSION` build argument and the Dokploy env no longer take part, so a stale env value cannot make the server report an old release.

## [0.12.0] - 2026-11-07

### Added

- An Android app (`Prisme-android.apk`) and a desktop app for macOS, Windows and Linux, attached to every GitHub release. Both sign in with your Facile account in the browser and come back through `studio.facile.prisme://auth`.
- A first-run intro in both apps that ends on the sign-in slide, then a welcome step that asks for notifications.
- Desktop notifications for direct messages and mentions while the window is in the background, and self-update from the latest release.
- Android push for direct messages and mentions, skipping muted rooms, even when the app is closed. Synapse now sends its events to the API (`PRISME_AS_URL`, `PRISME_HS_TOKEN`), and `FCM_SERVICE_ACCOUNT` turns sending on.
- `GET /api/push` and `POST /api/push/devices`; device tokens are deleted when their session signs out.
- A downloads page at `/downloads`, linked from the user menu on the web.

### Changed

- CORS always allows the apps' origins (`https://localhost`, `tauri://localhost`, `http://tauri.localhost`) while OIDC is on.
- The notification setting reads "Notifications on this device" in the apps.

## [0.11.3] - 2026-11-07

### Changed

- `@facile/muse` v0.14.2 to v0.14.4: the language chosen in the user menu re-renders the whole UI without a reload.
- The user menu's Subscription row drops its empty `plan` placeholder and shows only its label.

## [0.11.2] - 2026-11-07

### Security

- Registre tokens are sealed at rest. porte to v0.13.0: `porte_identities.access_token` and `refresh_token`
  are encrypted with AES-256-GCM under `PORTE_TOKEN_KEY`. Existing rows are sealed once, in the background,
  after boot, and the count is logged. An invalid `PORTE_TOKEN_KEY` stops the boot. With the key unset nothing changes.
- Losing `PORTE_TOKEN_KEY` makes the stored tokens unreadable: people sign in again, nothing else is lost.
  Do not roll back to a build older than porte v0.13.0 after the first reseal.

## [0.11.1] - 2026-11-07

### Changed

- `@facile/muse` v0.14.1 to v0.14.2: the user menu clamps itself inside the viewport at the right edge of a header.
- The user menu has a Subscription row that opens `https://sso.facile.studio/spaces`, where plans are managed.

## [0.11.0] - 2026-11-07
### Changed

- The sidebar user card and the mobile avatar open muse's user menu (account, settings, language, theme, help, sign out). The cog and the in-settings sign-out are gone
- Settings use muse's `SettingsLayout` with the sections Profile, Preferences, Notifications and Agents, at `/settings/<section>`. Profile is read-only and links to Facile. `/settings/appearance` redirects to `/settings/preferences`
- A language picked in the app is remembered in the `fc_locale` cookie and wins over the Registre profile locale
- muse v0.13.0 to v0.14.1

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
