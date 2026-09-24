# Pinned contract provenance

`contract/v1/schema.json` and `contract/v1/fixtures/` are a test snapshot of the upstream Public Engine Contract. They are not authoritative; `sibukixxx/insight` is.

| Field | Value |
|---|---|
| Upstream repository | https://github.com/sibukixxx/insight |
| Upstream path | `contracts/public-engine/v1/` |
| Upstream revision | `ae44542912fc5b5350c95795abac9e28162b1fcf` (branch `feat/65-per-run-model-routing`, insight PR #101) |
| Contract version | `1` |
| Synced | 2026-09-24 |

To resync, copy the upstream directory at a new revision, update this table and run `go test ./...`.

## Analytical Artifact v1

`contract/analytical-artifact/v1/` (schema and fixtures) is a test snapshot of insight `contracts/analytical-artifact/v1/` at revision `a90b5d4ae1a63c28f93adb509974dc97b75c9049` (insight `main`, 2026-09-24). Package `analytical` validates every pinned fixture.
