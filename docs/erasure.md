# Erasure

When Registre erases a person, Prisme deletes their account, their agents and their Matrix
data. This page says how, when Prisme answers the notice, and exactly what remains afterwards,
so the output of the suite's read-only leftover check (`scripts/erasure-check.sh` in Mycelium)
can be read correctly.

## How Registre tells Prisme

Registre queues one deletion notice per suite app when it erases a person. A person who deletes
their own account waits 30 days first; an operator's erasure is immediate. Prisme hears it two
ways, both handled by `porte/directory`:

- a signed `user.deleted` push to `POST /api/registre/deletions`, retried by Registre up to 12
  times when Prisme does not answer 2xx;
- a poll of `GET /api/directory/deletions` at boot and every five minutes, which lists erasures
  for 60 days.

Registre marks a notice delivered on any 2xx and never asks again. So Prisme answers success
only when the person was found and erased, or provably never held an account here.

## What Prisme does

`internal/erasure` runs the erasure. It finds the account through the identity row the subject
names, under either spelling of the issuer (`https://sso.facile.studio` and the same with a
trailing slash). Then, in this order:

1. **Synapse.** For every Matrix account the person owns (their own and each agent's), through
   Synapse's admin API:
   1. settle each room the account is joined to (see "The room rule");
   2. delete every file the account uploaded (`DELETE /_synapse/admin/v1/users/<id>/media`);
   3. deactivate the account with erasure
      (`POST /_synapse/admin/v1/deactivate/<id>` with `{"erase": true}`).
2. **Prisme's database**, in one transaction through `porte/kit`'s erase hook: record the
   erased subject as a hash, clear the person from a running call and from the spaces they
   created, delete their Matrix sessions, agents and account rows, then delete the `porte_users`
   row. Its cascade deletes the identities, sessions, login codes, magic links and push devices.

Every step is idempotent. If Synapse fails, nothing is deleted from the database, the notice is
answered with an error, and the retry finds the person again. If the transaction fails after
Synapse succeeded, the rows still name the person, and the retry repeats the Synapse calls
harmlessly (a deactivated account is in no room, has no media, and deactivates again) before
the transaction runs. Tests cover both failures (`internal/erasure/erasure_test.go`).

## When a notice is answered as done

| Case | Answer |
|---|---|
| The subject maps to an account | Done after both halves succeed |
| Prisme erased this subject before | Done (the hash in `erased_subjects` proves it) |
| The subject maps to nothing, and every account holds a Registre subject | Done: it was never here |
| The subject maps to nothing, and some account holds no subject | Error, so the notice stays pending |

An unlinked account may belong to the erased person, and the notice carries only the subject,
so Prisme cannot tell. It answers with an error instead of guessing.

Three things link an account to its subject: the person's first Prisme sign-in, the Registre
space sync (it links every member it sees), and, since v0.14.0, a service call that sends
`X-Prisme-Acting-Subject` next to `X-Prisme-Acting-For` (see [the API](api.md)). Before v0.14.0
Bulle named people by email only, so an account Bulle created for someone who never signed in to
Prisme and is in no Registre space held no subject. Such an account gains its subject at Bulle's
next call that sends the header, and from then on a deletion notice for that subject erases it
like any other. A notice that was pending meanwhile succeeds at the next delivery or poll.

The header never re-points an identity. When the subject already belongs to another account,
for example after the person changed their email in Registre, the call is refused with
`409 acting_subject_taken`, nothing is written, and the old account stays the one the notice
erases. An account a service key created under the new address before v0.14.0 then stays
unlinked and needs `erase-user` below.

## An erased subject never gets an account again

Since v0.14.2, every path that creates or links an account from a Registre subject first looks
the subject up in `erased_subjects` and refuses it when Prisme erased it:

| Path | Answer for an erased subject |
|---|---|
| Service call with `X-Prisme-Acting-Subject` | `409 acting_subject_erased`, nothing written (see [the API](api.md)) |
| Sign-in through Registre (porte's OIDC callback, and porte's bearer provisioning) | Refused with `permission_denied`, nothing written |
| Registre space sync | The member is skipped and logged, nothing written |

Before v0.14.2 a service call that arrived just after an erasure created and linked a new account
for the erased subject (2026-11-11). Two things now clean up such an account:

- The next delivery of the notice, by push or poll, erases it. porte/directory skips subjects it
  already handled, so Prisme's state store hands porte back as unhandled every handled subject an
  identity holds again. It reads the store once per process, so this takes effect from the next
  start.
- `erase-user <email>` erases it: an identity whose subject Prisme erased does not count as a
  link.

A service call that names a person by email alone (no `X-Prisme-Acting-Subject`) carries no
subject to check, so it still creates an unlinked account; `erase-user` settles that one.

At boot Prisme logs how many accounts hold no subject (a count, no names). To settle one, run
inside the API container:

```sh
/api erase-user <email>
```

It erases that account exactly as a notice would, Synapse first, and refuses an account that
holds a Registre subject Prisme has not erased. Once no unlinked account is left, pending notices for unknown
subjects succeed on the next poll.

## The admin credential

The admin API needs a server admin. Prisme uses the appservice token (`PRISME_AS_TOKEN`) for
it, which works because the appservice sender `@prisme` is a server admin. The Synapse image's
entrypoint (`deploy/synapse/render.py`) sets `users.admin = 1` for `@prisme` while Synapse is
stopped, which is the SQL Synapse's own admin API documentation gives for a first admin; Synapse
caches the flag per process, so it cannot be set while Synapse runs. On a first start, when
`@prisme` has no row yet, the entrypoint runs Synapse once to register it, stops it, and sets
the flag. A failure there is printed and Synapse still starts.

At boot the API makes one read-only admin call (`GET /_synapse/admin/v1/background_updates/enabled`,
which returns nothing about any person) and logs whether the token works. If it does not,
every erasure fails and is retried, and the log line says so.

This makes the appservice token more powerful than before: it could already act as every
Prisme user, and now it can also purge rooms and media. It still lives only in Dokploy's env.

## The room rule

For each room an erased account is joined to, before deactivating:

| Room | Action |
|---|---|
| Nobody else is in it (the bot and the person's own agents do not count), and it is not a Prisme space or channel | Purged: the room and its history are deleted |
| Nobody else is in it, but it is a Prisme space or channel | Kept with the bot, which administers it; the space or Registre owns it |
| Others are in it and one of them, or the bot, can still change its power levels | Kept |
| Others are in it, only an erased account could administer it | The most powerful remaining member (lowest id on a tie) is raised to the erased account's power through `make_room_admin`, then the room is kept |
| Others are in it and nobody could administer it before the erasure | Kept as it was, logged at warn |

From room version 12 a creator holds unlimited power; Prisme reads the create event for that.

## Matrix ids are never reused

A deactivated Matrix id is never handed to anyone else. When a new person signs in with the
email of an erased one, Prisme tries the same localpart, Synapse answers that the id is taken,
the admin API confirms it is deactivated, and Prisme picks the next candidate: `alice2` instead
of `alice`. Agent names follow the same rule. If the admin API cannot confirm the state of a
taken id, Prisme also moves on, so a new person is never bound to an erased account. Tests:
`TestANewPersonWithAnErasedEmailGetsANewMatrixAccount` (fake Synapse) and
`TestLiveErasureAgainstSynapse` (real Synapse).

Erasures before v0.13.0 kept a row in `accounts` (under a negated user id) and in `agents`
(with no owner) to reserve the id. Boot deletes those rows: Synapse answers instead.

## What remains, and why

### Prisme's database

The leftover check should report **0** for `prisme`. What remains holds no email, subject or
Matrix id:

| Table | What stays | Why |
|---|---|---|
| `erased_subjects` | `sha256(issuer + " " + subject)` and the time | Proves a repeated notice was handled |
| `spaces` | Rows the person created, with `created_by = 0` | The space belongs to its other members |
| `calls` | A running call they started, with `started_by = ''` | The call belongs to its room |
| `directory_deletion_state` | porte's cursor and handled subjects | porte/directory's own state; subjects are Registre's random ids |

### Synapse's database

The leftover check should report a non-zero `synapse` total. Measured on Synapse v1.164.1 after
erasing a person with one agent, three rooms, three messages and one upload (one room purged):

| Table and column | Rows | What it is |
|---|---|---|
| `users.name` | 1 per account | The deactivated account. Synapse keeps it so the id is never reissued |
| `erased_users.user_id` | 1 per account | Synapse's erasure flag, which it reads to prune their events |
| `events.sender`, `events.state_key` | 1 per event | Event metadata: id, type, room, time. No content |
| `event_json` (matched on content, not by the check) | 1 per message in a kept room | The message body stays in the database. Synapse serves it pruned (no content) to anyone who was not in the room when it was sent; people who were there keep seeing it, as the Matrix specification describes for erasure |
| `room_memberships`, `local_current_membership`, `current_state_events`, `current_state_delta_stream`, `state_events`, `state_groups_state`, `sliding_sync_membership_snapshots` | per membership change | Membership history, now ending in `leave`. Synapse redacts the content of the person's membership events (display name, avatar) when it parts them |
| `rooms.creator` | 1 per kept room they created | The room belongs to its other members |
| `event_push_actions`, `event_push_summary` | a few | Notification counters for the person; no content |
| `presence_stream`, `user_stats_current` | 1 each | Last presence update and room counters |
| `scheduled_tasks.resource_id` | 1 | Synapse's completed `update_join_states` task for the deactivation |

Gone from Synapse: the profile (display name and avatar), devices and access tokens, account
data (including direct-chat lists and push rules), the user directory entry, every uploaded
file, every purged room with all its events, and every pending invite.

The leftover check finds Synapse rows by `@<part of the email before the @>:prisme.facile.studio`.
That is the person's Matrix id only for a simple address. Prisme turns other characters into dots,
caps the localpart at 40 characters and numbers it when taken (`alice2`); agents have their own
ids. A zero for `synapse` can therefore mean the check looked under the wrong id.

### Outside both databases

- **Federation**: none. `homeserver.yaml` sets `federation_domain_whitelist: []` and the
  listener serves the client API only, so no other server holds a copy.
- **Devices**: messages and files the person's contacts already downloaded stay on their
  devices, and Android push notifications already shown stay in their notification history.
- **Logs**: Synapse's container log records the deactivation with the Matrix id. Prisme's own
  log lines carry counts and room ids only.
- **Bulle** keeps its own data about the person's agents and erases it on its own notice.

Back to the [documentation index](README.md).
