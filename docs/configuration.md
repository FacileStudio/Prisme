# Configuration

The full variable table lives in the README. This page covers what the table cannot.

## Where each value is read

- `apps/api/internal/env/env.go` reads the API's variables at boot. `PRISME_AS_TOKEN` and a
  `PRISME_SECRET_KEY` of at least 32 characters are required; everything else has a default or
  turns a feature off.
- `deploy/synapse/render.py` renders Synapse's config into the `prisme_synapse_data` volume at
  every start and refuses to start when a required variable is missing.
- `deploy/postgres/init-synapse.sql` creates the `synapse` database (C collation, as Synapse
  requires) the first time the Postgres volume is initialised. An existing volume is not touched.

## Registre shared spaces

The same `OIDC_ISSUER`, `OIDC_CLIENT_ID` and `OIDC_CLIENT_SECRET` that sign people in also turn on
the Registre space sync: Prisme calls Registre's directory API with HTTP Basic on its own client
id and secret. Without all three, spaces stay local and `POST /api/registre/sync` is not mounted.
For Registre to nudge Prisme, its seed entry for `prisme` needs
`directory_sync_uri: https://prisme.facile.studio/api/registre/sync`; without it the 5-minute pull
is the only trigger.

## Secrets

Production secrets live in the Dokploy compose `Prisme` (project Facile, environment
production) and, as a copy for the owner, in `tiroir` under `PRISME_*` names. None is committed.
Registre holds the matching OIDC client secret as `REGISTRE_SEED_PRISME_SECRET`.

## Rotating

- `PRISME_SECRET_KEY`: stored Matrix tokens no longer open, so every browser and agent gets a new
  device on its next call. Old devices stay valid until revoked.
- `PRISME_AS_TOKEN`: change it in Dokploy; both Synapse and the API read it from there.
- `SYNAPSE_SIGNING_KEY`: do not rotate casually; with federation off it only signs local events.
