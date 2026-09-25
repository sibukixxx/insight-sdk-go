// Package analytical is the producer-side Go implementation of the Insight
// Analytical Artifact v1 wire contract (insight contracts/analytical-artifact/v1,
// pinned in contract/analytical-artifact/v1). Producers build, seal (hash) and
// validate artifacts here, then submit them with AddEvidenceRequest
// .AnalyticalArtifacts. It never imports Insight internals.
//
// An Analytical Artifact records deterministic calculation output only. It is
// never a cause, hypothesis, claim or insight; meaning is assigned by the
// Research Engine.
package analytical

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

const (
	Schema  = "insight-lab.analytical-artifact"
	Version = "1"
)

var (
	ErrInvalidArtifact  = errors.New("invalid analytical artifact")
	ErrIdentityConflict = errors.New("analytical artifact identity conflict")
)

type Hash struct {
	Algorithm string `json:"algorithm"`
	Value     string `json:"value"`
}

// SHA256 returns the contract hash of data.
func SHA256(data []byte) Hash {
	sum := sha256.Sum256(data)
	return Hash{Algorithm: "sha256", Value: hex.EncodeToString(sum[:])}
}

type DatasetRef struct {
	ID      string `json:"id"`
	URI     string `json:"uri,omitempty"`
	Version string `json:"version"`
	Hash    Hash   `json:"hash"`
}

type SpecRef struct {
	Kind      string `json:"kind"`
	Reference string `json:"reference"`
	Hash      Hash   `json:"hash"`
}

// ExternalSubjectRef is an opaque producer-owned identity. Insight never
// interprets it as domain state.
type ExternalSubjectRef struct {
	Namespace string `json:"namespace"`
	ID        string `json:"id"`
}

type Period struct {
	Start string `json:"start"`
	End   string `json:"end"`
	Basis string `json:"basis,omitempty"`
}

type Population struct {
	Description string `json:"description"`
	Unit        string `json:"unit,omitempty"`
}

type MetricDefinition struct {
	Version     string `json:"version,omitempty"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Unit        string `json:"unit"`
	Aggregation string `json:"aggregation,omitempty"`
}

type QualityFlag struct {
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
}

type TemporalMetadata struct {
	ObservedAt time.Time `json:"observedAt"`
	Origin     string    `json:"origin"`
	Geography  string    `json:"geography"`
	ValueBasis string    `json:"valueBasis"`
}

type Result struct {
	MetricID   string            `json:"metricId"`
	Temporal   *TemporalMetadata `json:"temporal,omitempty"`
	Dimensions map[string]string `json:"dimensions,omitempty"`
	Period     Period            `json:"period"`
	Value      json.RawMessage   `json:"value,omitempty"`
	Missing    bool              `json:"missing,omitempty"`
	Quality    []QualityFlag     `json:"qualityFlags,omitempty"`
}

type Computation struct {
	Engine        string `json:"engine"`
	EngineVersion string `json:"engineVersion"`
	Deterministic bool   `json:"deterministic"`
	Timezone      string `json:"timezone,omitempty"`
}

type SourceProvenance struct {
	DatasetID          string    `json:"datasetId"`
	Source             string    `json:"source"`
	RetrievedAt        time.Time `json:"retrievedAt"`
	License            string    `json:"license,omitempty"`
	TransformationRefs []string  `json:"transformationRefs,omitempty"`
}

type Artifact struct {
	ArtifactSchema  string              `json:"artifactSchema"`
	SchemaVersion   string              `json:"schemaVersion"`
	ID              string              `json:"id"`
	ArtifactHash    Hash                `json:"artifactHash"`
	Producer        string              `json:"producer"`
	ProducerVersion string              `json:"producerVersion"`
	GeneratedAt     time.Time           `json:"generatedAt"`
	ExternalSubject *ExternalSubjectRef `json:"externalSubject,omitempty"`

	Datasets    []DatasetRef       `json:"datasets"`
	Spec        SpecRef            `json:"spec"`
	Parameters  map[string]any     `json:"parameters,omitempty"`
	Dimensions  []string           `json:"dimensions,omitempty"`
	Filters     map[string]any     `json:"filters,omitempty"`
	Period      Period             `json:"period"`
	Population  Population         `json:"population"`
	Metrics     []MetricDefinition `json:"metrics"`
	Results     []Result           `json:"results"`
	Quality     []QualityFlag      `json:"qualityFlags,omitempty"`
	Computation Computation        `json:"computation"`
	Provenance  []SourceProvenance `json:"provenance"`
}

// NumberValue encodes a numeric result value.
func NumberValue(v float64) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}

// Import decodes and validates exactly one artifact.
func Import(data []byte) (*Artifact, error) {
	var a Artifact
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := dec.Decode(&a); err != nil {
		return nil, fmt.Errorf("%w: decode: %v", ErrInvalidArtifact, err)
	}
	var trailing any
	if err := dec.Decode(&trailing); err != io.EOF {
		return nil, fmt.Errorf("%w: input must contain exactly one JSON value", ErrInvalidArtifact)
	}
	if err := a.Validate(); err != nil {
		return nil, err
	}
	return &a, nil
}

// Export validates and serializes an artifact.
func Export(a Artifact) ([]byte, error) {
	if err := a.Validate(); err != nil {
		return nil, err
	}
	return json.MarshalIndent(a, "", "  ")
}

// Seal fills schema identifiers and computes artifactHash from the canonical
// serialization of the artifact with an empty artifactHash. The artifact
// producer owns this hash; the SDK does not attach consumer-specific policy.
func (a *Artifact) Seal() error {
	a.ArtifactSchema, a.SchemaVersion = Schema, Version
	a.ArtifactHash = Hash{}
	canonical, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf("%w: canonical encode: %v", ErrInvalidArtifact, err)
	}
	a.ArtifactHash = SHA256(canonical)
	return a.Validate()
}

// Validate mirrors the public v1 contract rules.
func (a Artifact) Validate() error {
	fail := func(format string, args ...any) error {
		return fmt.Errorf("%w: %s", ErrInvalidArtifact, fmt.Sprintf(format, args...))
	}
	if a.ArtifactSchema != Schema || a.SchemaVersion != Version {
		return fail("unsupported schema %q version %q", a.ArtifactSchema, a.SchemaVersion)
	}
	if blank(a.ID) || blank(a.Producer) || blank(a.ProducerVersion) || a.GeneratedAt.IsZero() {
		return fail("id, producer, producerVersion and generatedAt are required")
	}
	if a.ExternalSubject != nil && (blank(a.ExternalSubject.Namespace) || blank(a.ExternalSubject.ID)) {
		return fail("externalSubject requires namespace and id")
	}
	if err := validateHash("artifactHash", a.ArtifactHash); err != nil {
		return fail("%v", err)
	}
	if len(a.Datasets) == 0 {
		return fail("at least one dataset is required")
	}
	datasetIDs := map[string]struct{}{}
	for i, d := range a.Datasets {
		if blank(d.ID) || blank(d.Version) {
			return fail("datasets[%d] requires id and version", i)
		}
		if _, dup := datasetIDs[d.ID]; dup {
			return fail("duplicate dataset id %q", d.ID)
		}
		datasetIDs[d.ID] = struct{}{}
		if err := validateHash(fmt.Sprintf("datasets[%d].hash", i), d.Hash); err != nil {
			return fail("%v", err)
		}
	}
	if blank(a.Spec.Kind) || blank(a.Spec.Reference) {
		return fail("spec kind and reference are required")
	}
	if err := validateHash("spec.hash", a.Spec.Hash); err != nil {
		return fail("%v", err)
	}
	if err := validatePeriod("period", a.Period); err != nil {
		return fail("%v", err)
	}
	if blank(a.Population.Description) {
		return fail("population.description is required")
	}
	if len(a.Metrics) == 0 || len(a.Results) == 0 {
		return fail("metrics and results must not be empty")
	}
	metricIDs := map[string]struct{}{}
	for i, m := range a.Metrics {
		if blank(m.ID) || blank(m.Name) || blank(m.Unit) {
			return fail("metrics[%d] requires id, name and unit", i)
		}
		if _, dup := metricIDs[m.ID]; dup {
			return fail("duplicate metric id %q", m.ID)
		}
		metricIDs[m.ID] = struct{}{}
	}
	for i, r := range a.Results {
		if _, ok := metricIDs[r.MetricID]; !ok {
			return fail("results[%d] references unknown metric %q", i, r.MetricID)
		}
		if err := validatePeriod(fmt.Sprintf("results[%d].period", i), r.Period); err != nil {
			return fail("%v", err)
		}
		if r.Temporal != nil {
			if err := validateTemporal(r.Temporal); err != nil {
				return fail("results[%d].temporal: %v", i, err)
			}
		}
		if r.Missing && len(r.Value) != 0 {
			return fail("results[%d] cannot have value when missing", i)
		}
		if !r.Missing && !validScalar(r.Value) {
			return fail("results[%d].value must be a JSON scalar", i)
		}
	}
	if blank(a.Computation.Engine) || blank(a.Computation.EngineVersion) || !a.Computation.Deterministic {
		return fail("computation requires engine, engineVersion and deterministic=true")
	}
	if len(a.Provenance) == 0 {
		return fail("source provenance is required")
	}
	covered := map[string]struct{}{}
	for i, p := range a.Provenance {
		if _, ok := datasetIDs[p.DatasetID]; !ok {
			return fail("provenance[%d] references unknown dataset %q", i, p.DatasetID)
		}
		if blank(p.Source) || p.RetrievedAt.IsZero() {
			return fail("provenance[%d] requires source and retrievedAt", i)
		}
		covered[p.DatasetID] = struct{}{}
	}
	for id := range datasetIDs {
		if _, ok := covered[id]; !ok {
			return fail("dataset %q has no source provenance", id)
		}
	}
	return nil
}

// CheckDuplicate: same id + same hash is a no-op duplicate; same id with a
// different hash is an identity conflict.
func CheckDuplicate(existing, incoming Artifact) (bool, error) {
	if existing.ID != incoming.ID {
		return false, nil
	}
	if strings.EqualFold(existing.ArtifactHash.Value, incoming.ArtifactHash.Value) {
		return true, nil
	}
	return false, fmt.Errorf("%w: id %q has different artifactHash", ErrIdentityConflict, incoming.ID)
}

// ReproducibilityKey follows the public definition: dataset ids/versions/hashes,
// spec hash, parameters and filters.
func (a Artifact) ReproducibilityKey() string {
	parts := []string{strings.ToLower(a.Spec.Hash.Algorithm) + ":" + strings.ToLower(a.Spec.Hash.Value)}
	for _, d := range a.Datasets {
		parts = append(parts, d.ID+"@"+d.Version+":"+strings.ToLower(d.Hash.Algorithm)+":"+strings.ToLower(d.Hash.Value))
	}
	params, _ := json.Marshal(a.Parameters)
	filters, _ := json.Marshal(a.Filters)
	parts = append(parts, "parameters:"+string(params), "filters:"+string(filters))
	sort.Strings(parts)
	sum := sha256.Sum256([]byte(strings.Join(parts, "|")))
	return "sha256:" + hex.EncodeToString(sum[:])
}

func validateTemporal(t *TemporalMetadata) error {
	if t.ObservedAt.IsZero() || blank(t.Geography) {
		return errors.New("observedAt and geography are required")
	}
	switch t.Origin {
	case "observed", "derived":
	default:
		return fmt.Errorf("origin %q must be observed or derived", t.Origin)
	}
	switch t.ValueBasis {
	case "nominal", "real", "not_applicable":
	default:
		return fmt.Errorf("valueBasis %q is not supported", t.ValueBasis)
	}
	return nil
}

func validateHash(name string, h Hash) error {
	if strings.TrimSpace(h.Algorithm) != "sha256" {
		return fmt.Errorf("%s.algorithm must be sha256", name)
	}
	v := strings.TrimSpace(h.Value)
	if len(v) != 64 {
		return fmt.Errorf("%s.value must contain 64 hexadecimal characters", name)
	}
	if _, err := hex.DecodeString(v); err != nil {
		return fmt.Errorf("%s.value must be hexadecimal", name)
	}
	return nil
}

func validatePeriod(name string, p Period) error {
	if blank(p.Start) || blank(p.End) {
		return fmt.Errorf("%s requires start and end", name)
	}
	if p.Start > p.End {
		return fmt.Errorf("%s.start must not be after end", name)
	}
	return nil
}

func validScalar(raw json.RawMessage) bool {
	if len(raw) == 0 {
		return false
	}
	var v any
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if dec.Decode(&v) != nil {
		return false
	}
	switch v.(type) {
	case nil, bool, string, json.Number:
		return true
	}
	return false
}

func blank(s string) bool { return strings.TrimSpace(s) == "" }
