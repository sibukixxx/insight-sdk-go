package analytical

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestImportAcceptsVendoredPublicFixtures(t *testing.T) {
	paths, err := filepath.Glob("../contract/analytical-artifact/v1/fixtures/*.json")
	if err != nil || len(paths) == 0 {
		t.Fatalf("no vendored fixtures: %v", err)
	}
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := Import(data); err != nil {
			t.Errorf("%s: public fixture rejected: %v", p, err)
		}
	}
}

func TestImportRejectsValueOnMissingResult(t *testing.T) {
	a := validArtifact(t)
	a.Results[0].Missing = true
	if err := a.Validate(); !errors.Is(err, ErrInvalidArtifact) {
		t.Fatalf("want ErrInvalidArtifact, got %v", err)
	}
}

func TestValidateRejectsDatasetWithoutProvenance(t *testing.T) {
	a := validArtifact(t)
	a.Datasets = append(a.Datasets, DatasetRef{ID: "other", Version: "1", Hash: SHA256([]byte("x"))})
	if err := a.Validate(); !errors.Is(err, ErrInvalidArtifact) {
		t.Fatalf("want ErrInvalidArtifact, got %v", err)
	}
}

func TestSealProducesStableHashForSameContent(t *testing.T) {
	a, b := validArtifact(t), validArtifact(t)
	if a.ArtifactHash != b.ArtifactHash {
		t.Fatalf("same content produced different hashes: %v vs %v", a.ArtifactHash, b.ArtifactHash)
	}
	b.Population.Description = "changed"
	if err := b.Seal(); err != nil {
		t.Fatal(err)
	}
	if a.ArtifactHash == b.ArtifactHash {
		t.Fatal("changed content kept the same hash")
	}
}

func TestCheckDuplicateReturnsConflictWhenSameIDHasDifferentHash(t *testing.T) {
	a, b := validArtifact(t), validArtifact(t)
	dup, err := CheckDuplicate(a, b)
	if !dup || err != nil {
		t.Fatalf("identical artifacts: dup=%v err=%v", dup, err)
	}
	b.Population.Description = "changed"
	_ = b.Seal()
	if _, err := CheckDuplicate(a, b); !errors.Is(err, ErrIdentityConflict) {
		t.Fatalf("want ErrIdentityConflict, got %v", err)
	}
}

func TestExportRoundTripsThroughImport(t *testing.T) {
	a := validArtifact(t)
	data, err := Export(a)
	if err != nil {
		t.Fatal(err)
	}
	got, err := Import(data)
	if err != nil {
		t.Fatal(err)
	}
	if got.ReproducibilityKey() != a.ReproducibilityKey() {
		t.Fatal("reproducibility key changed across round trip")
	}
	var generic map[string]any
	if err := json.Unmarshal(data, &generic); err != nil {
		t.Fatal(err)
	}
	if generic["artifactSchema"] != Schema {
		t.Fatalf("artifactSchema = %v", generic["artifactSchema"])
	}
}

func validArtifact(t *testing.T) Artifact {
	t.Helper()
	ts := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	a := Artifact{
		ID: "a-1", Producer: "techvit-test", ProducerVersion: "0.1.0", GeneratedAt: ts,
		Datasets:    []DatasetRef{{ID: "ds", Version: "v1", Hash: SHA256([]byte("rows"))}},
		Spec:        SpecRef{Kind: "declarative", Reference: "spec.json", Hash: SHA256([]byte("spec"))},
		Period:      Period{Start: "2026-01", End: "2026-06"},
		Population:  Population{Description: "rows"},
		Metrics:     []MetricDefinition{{ID: "m", Name: "M", Unit: "count"}},
		Results:     []Result{{MetricID: "m", Period: Period{Start: "2026-01", End: "2026-06"}, Value: NumberValue(3)}},
		Computation: Computation{Engine: "techvit", EngineVersion: "0.1.0", Deterministic: true},
		Provenance:  []SourceProvenance{{DatasetID: "ds", Source: "urn:test", RetrievedAt: ts}},
	}
	if err := a.Seal(); err != nil {
		t.Fatal(err)
	}
	return a
}
