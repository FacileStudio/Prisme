# Changelog

## [0.1.0] - 2026-11-03

### Added

- feat(client): full Svelte 5 chat web UI with muse design tokens, spaces, channels, DMs, threads, reactions, calls and agent management
- feat(api): real-time agent name availability check endpoint `GET /api/agents/check`

### Changed

- feat(api): simplified agent Matrix username to `@<name>` with global uniqueness enforcement across human accounts and active agents
