package analytical_test

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	insight "github.com/sibukixxx/insight-sdk-go"
	"github.com/sibukixxx/insight-sdk-go/analytical"
)

// Build, seal and submit a deterministic result. Unknown values are
// missing, never zero.
func Example() {
	at := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	a := analytical.Artifact{
		ID: "monthly-revenue-2026-01", Producer: "my-pipeline", ProducerVersion: "1.0.0", GeneratedAt: at,
		Datasets:    []analytical.DatasetRef{{ID: "sales", Version: "2026-01", Hash: analytical.SHA256([]byte("raw bytes"))}},
		Spec:        analytical.SpecRef{Kind: "sql", Reference: "queries/revenue.sql", Hash: analytical.SHA256([]byte("select ..."))},
		Period:      analytical.Period{Start: "2026-01", End: "2026-01"},
		Population:  analytical.Population{Description: "orders"},
		Metrics:     []analytical.MetricDefinition{{ID: "revenue", Name: "Revenue", Unit: "JPY"}},
		Results:     []analytical.Result{{MetricID: "revenue", Period: analytical.Period{Start: "2026-01", End: "2026-01"}, Missing: true}},
		Computation: analytical.Computation{Engine: "duckdb", EngineVersion: "1.1", Deterministic: true},
		Provenance:  []analytical.SourceProvenance{{DatasetID: "sales", Source: "urn:sales", RetrievedAt: at}},
	}
	if err := a.Seal(); err != nil {
		panic(err)
	}
	body, _ := analytical.Export(a)
	req := insight.AddEvidenceRequest{AnalyticalArtifacts: []json.RawMessage{body}}
	_ = context.Background
	fmt.Println(len(req.AnalyticalArtifacts), a.ArtifactHash.Algorithm)
	// Output: 1 sha256
}
