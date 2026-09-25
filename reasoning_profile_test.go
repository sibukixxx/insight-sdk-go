package insight

import (
	"encoding/json"
	"testing"
)

func TestReasoningProfileWireRoundTrip(t *testing.T) {
	req := StartAnalysisRequest{ReasoningProfile: ReasoningProfileCustomerInsight}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{"reasoningProfile":"CUSTOMER_INSIGHT"}` {
		t.Fatalf("unexpected request JSON: %s", b)
	}

	var info EngineInfo
	if err := json.Unmarshal([]byte(`{"supportedReasoningProfiles":["GENERAL_RESEARCH","CUSTOMER_INSIGHT"],"defaultReasoningProfile":"GENERAL_RESEARCH"}`), &info); err != nil {
		t.Fatal(err)
	}
	if info.DefaultReasoningProfile != ReasoningProfileGeneralResearch {
		t.Fatalf("default = %q", info.DefaultReasoningProfile)
	}
	if len(info.SupportedReasoningProfiles) != 2 || info.SupportedReasoningProfiles[1] != ReasoningProfileCustomerInsight {
		t.Fatalf("supported = %#v", info.SupportedReasoningProfiles)
	}
}
