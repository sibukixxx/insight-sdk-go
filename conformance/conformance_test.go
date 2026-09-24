package conformance

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	insight "github.com/sibukixxx/insight-sdk-go"
)

const fixtureDir = "../contract/v1/fixtures"

// Every operation used by the pinned fixtures must be callable by this SDK.
func TestPinnedFixturesUseOnlyOperationsTheSDKImplements(t *testing.T) {
	fixtures, err := Load(fixtureDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(fixtures) == 0 {
		t.Fatal("no pinned fixtures")
	}
	for _, f := range fixtures {
		for _, s := range f.Steps {
			_, err := call(context.Background(), insight.NewClient("http://127.0.0.1:1"), s.Op, map[string]string{}, nil)
			if err != nil && strings.HasPrefix(err.Error(), "unknown op") {
				t.Errorf("%s uses operation %q the SDK does not implement", f.Name, s.Op)
			}
		}
	}
}

// Live conformance against a running engine. Set INSIGHT_DETERMINISTIC_URL
// (and optionally INSIGHT_MODEL_BACKED_URL for model-backed fixtures).
func TestPinnedFixturesPassAgainstLiveEngine(t *testing.T) {
	urls := map[string]string{"deterministic": os.Getenv("INSIGHT_DETERMINISTIC_URL"), "model_backed": os.Getenv("INSIGHT_MODEL_BACKED_URL")}
	if urls["deterministic"] == "" {
		t.Skip("INSIGHT_DETERMINISTIC_URL not set; live conformance skipped")
	}
	fixtures, err := Load(fixtureDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range fixtures {
		t.Run(f.Name, func(t *testing.T) {
			url := urls[f.Engine]
			if url == "" {
				t.Skipf("no engine URL for %s", f.Engine)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()
			if err := Run(ctx, insight.NewClient(url, insight.WithPollInterval(50*time.Millisecond)), f); err != nil {
				t.Fatal(err)
			}
		})
	}
}
