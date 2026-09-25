# insight-sdk-go

Thin, dependency-free Go client for the [Insight](https://github.com/sibukixxx/insight) **Public Engine Contract v1**.

`sibukixxx/insight` owns Research semantics and the Public Contract. This SDK holds no research logic, scoring, storage, queue or domain policy, and never imports Insight's internal packages. The SDK is optional: the engine is fully usable through its HTTP contract alone.

```text
consumer
   ↓
insight-sdk-go        (this repository)
   ↓
Public Engine Contract (insight: contracts/public-engine/v1)
   ↓
insight OSS
```


<!-- role-boundary:v1 -->
## Role and boundaries

**Role:** Go consumersからInsight Public Engine Contractを安全に利用するための **thin client / compatibility layer**。

### Owns

- public request/response types
- transport, cancellation, timeout, typed errors
- contract-version compatibility and conformance
- Analytical Artifact helper exposed by the public contract

### Does not own

- Research algorithms or semantics
- storage, queue, scheduler, managed infrastructure
- consumer/domain policy
- private business workflow
- a second copy of the Insight engine

### Integration

```text
Go consumer → insight-sdk-go → Public Engine Contract → Insight
```

Upstream Insight contract is authoritative. SDK-specific convenience must not become new Research semantics.
## Install

```sh
go get github.com/sibukixxx/insight-sdk-go@v0.5.0
```

## Quickstart

Start an engine (`insight-lab`, default `http://127.0.0.1:8787`), then:

```go
client := insight.NewClient("http://127.0.0.1:8787")
subject, err := client.CreateSubject(ctx, insight.CreateSubjectRequest{
	IdempotencyKey: "my-subject",
	Subject:        insight.SubjectRef{Namespace: "my-app", ID: "item-42"},
})
// AddEvidence -> StartAnalysis(researchQuestion) -> WaitForAnalysis -> GetAnalysisResults
// CreateResearchRun -> AppendIteration -> GetResearchRun
```

Runnable example: `go run ./example/minimal -engine http://127.0.0.1:8787`.

- Errors are `*insight.Error` with a contract `Code`; `errors.Is(err, &insight.Error{Code: insight.CodeNotFound})` matches on the code. `CodeUnavailable` means the engine could not be reached or did not answer with a contract message.
- The client fills `ContractVersion` and a random `IdempotencyKey` when empty. Reuse your own key on retry to get the first result back.
- `context.Context` cancellation/timeouts are honoured on every call.
- `WithTransport` replaces HTTP (in-process engine, recorded responses, queues) without changing request/result types. `WithHTTPClient` configures the default transport.
- `ResearchResult.View()` decodes the research artifact fields the contract promises.

## Analytical Artifacts

Package `github.com/sibukixxx/insight-sdk-go/analytical` builds, seals (`artifactHash`) and validates Analytical Artifact v1 — deterministic results produced outside Insight (SQL, dataframes, BI exports). Submit the exported JSON via `AddEvidenceRequest.AnalyticalArtifacts`. A result is never a cause or insight; unknown values are `missing`, never zero.

## Compatibility

| SDK version | Contract versions | Pinned contract source |
|---|---|---|
| v0.5.x | `1` (adds domain-neutral `researchQuestion`, question input diff, generic evidence sources; insight #108) | `contract/v1` + `contract/analytical-artifact/v1` — see [contract/PROVENANCE.md](contract/PROVENANCE.md) |
| v0.4.x | `1` (adds `EngineInfo.ModelBacked`, insight #104) | `contract/v1` + `contract/analytical-artifact/v1` — see [contract/PROVENANCE.md](contract/PROVENANCE.md) |
| v0.3.x | `1` (adds modelBindings / modelRouting, timeline scenarioEvents, comparison informational diff; `analytical` package) | `contract/v1` + `contract/analytical-artifact/v1` — see [contract/PROVENANCE.md](contract/PROVENANCE.md) |
| v0.2.x | `1` (adds InputSource, ExecutionProfile, run comparison, re-evaluation, timeline, temporal operations, scenarios, data triage) | `contract/v1` — see [contract/PROVENANCE.md](contract/PROVENANCE.md) |
| v0.1.x | `1` (original v0 surface) | insight `e6e402d` |

- Pre-1.0 semver: minor versions may add operations/fields; patch versions never change behavior.
- v0.5.0 makes the analysis question first-class semantic input. Set `StartAnalysisRequest.ResearchQuestion` when the later ResearchRun is answering a specific question; leave it empty for open-ended discovery. Same evidence under a different question is reported as an input change.
- Unknown response fields are ignored (additive contract evolution). A response with a different `contractVersion` fails with `UNSUPPORTED_CONTRACT_VERSION`.
- v0.2.0 added InputSource / RawArtifact (insight #90) and ExecutionProfile (insight #91) additively; every v0.1 call keeps working.
- Types for operations added after v0.1 are generated from the pinned schema: `go run ./internal/typegen` (a test fails when `types_gen.go` is stale).

## Verification

```sh
go vet ./... && go test ./...
# live conformance against a running engine (deterministic fixtures; 08/09 need
# -input-root pointing at contract/v1/fixtures/data and -heavy-dir):
INSIGHT_DETERMINISTIC_URL=http://127.0.0.1:8787 go test ./conformance -run Live -v
```

`drift_test.go` fails when SDK types and the pinned `contract/v1/schema.json` diverge, and when `go.mod` gains any dependency. The pinned files are a test snapshot; upstream `insight` is authoritative.

### Model-backed fixtures without a real LLM

Fixtures whose `engine` is `model_backed` need an engine with a model. The insight repository ships a deterministic test model for this (insight PR #100, `cmd/insight-scripted-llm`):

```sh
# in a checkout of sibukixxx/insight
go run ./cmd/insight-scripted-llm -addr 127.0.0.1:8788 &
# model-backed engine: fixture 17 asserts allowedModels[0] == "scripted-model-large" and binds "scripted-model"
go run ./cmd/insight-lab -port 8787 -no-browser -db /tmp/insight-model.db -base-url http://127.0.0.1:8788 \
  -model scripted-model -allowed-models scripted-model-large -api-key scripted &
# deterministic engine: fixtures 08/09 need -input-root and -heavy-dir
go run ./cmd/insight-lab -port 8789 -no-browser -db /tmp/insight-det.db \
  -input-root <this repo>/contract/v1/fixtures/data -heavy-dir /tmp/heavy &
# here
INSIGHT_DETERMINISTIC_URL=http://127.0.0.1:8789 INSIGHT_MODEL_BACKED_URL=http://127.0.0.1:8787 go test ./conformance -run Live -v
```

Give each engine its own `-db`; without it both would share the default database in the OS data directory. The pinned fixture set remains 17 cases; v0.5.0 additionally relies on SDK drift/unit tests for the additive question-conditioned fields. The canonical contract is insight `main` at `9407a9f`. The canonical description of this setup is insight `docs/public-engine-contract.md` ("Running model-backed fixtures outside this repository"); if the two disagree, insight wins.

## License

Apache-2.0
