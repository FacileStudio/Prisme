# Prisme API

Go API and Matrix application service. See the root README for configuration and
[docs/api.md](../../docs/api.md) for the routes.

```sh
go test ./...
go run .   # needs DATABASE_URL, PRISME_AS_TOKEN, PRISME_SECRET_KEY and a reachable Synapse
```
