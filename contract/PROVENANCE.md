# Pinned contract provenance

`contract/v1/schema.json` and `contract/v1/fixtures/` are a test snapshot of the upstream Public Engine Contract. They are not authoritative; `sibukixxx/insight` is.

| Field | Value |
|---|---|
| Upstream repository | https://github.com/sibukixxx/insight |
| Upstream path | `contracts/public-engine/v1/` |
| Upstream revision | `9fd50185862365fbc133d1b63d2569f6a6b7a1a1` (insight `main`) |
| Contract version | `1` |
| Synced | 2026-09-25 |

To resync, copy the upstream directory at a new revision, update this table and run `go test ./...`.

## Analytical Artifact v1

`contract/analytical-artifact/v1/` (schema and fixtures) is a test snapshot of insight `contracts/analytical-artifact/v1/` at revision `9fd50185862365fbc133d1b63d2569f6a6b7a1a1` (insight `main`, 2026-09-25). Package `analytical` validates every pinned fixture.
