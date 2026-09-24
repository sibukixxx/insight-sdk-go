// Operations added after v0.1.0 (insight #66 #71 #73 #74 #83 #92). One method per
// contract operation; no research logic.

package insight

import (
	"context"
	"net/http"
	"net/url"
)

// CreateDatasetProfile calls POST /subjects/{subjectId}/dataset-profiles.
func (c *Client) CreateDatasetProfile(ctx context.Context, subjectID string, req CreateDatasetProfileRequest) (*DatasetProfile, error) {
	fillEnvelope(&req.ContractVersion, &req.IdempotencyKey)
	var out DatasetProfile
	return &out, c.do(ctx, http.MethodPost, "/subjects/"+url.PathEscape(subjectID)+"/dataset-profiles", req, &out)
}

// GetDatasetProfile calls GET /subjects/{subjectId}/dataset-profiles/{profileId}.
func (c *Client) GetDatasetProfile(ctx context.Context, subjectID, profileID string) (*DatasetProfile, error) {
	var out DatasetProfile
	return &out, c.do(ctx, http.MethodGet, "/subjects/"+url.PathEscape(subjectID)+"/dataset-profiles/"+url.PathEscape(profileID), nil, &out)
}

// Triage calls POST /subjects/{subjectId}/dataset-profiles/{profileId}/triage.
func (c *Client) Triage(ctx context.Context, subjectID, profileID string, req TriageRequest) (*SelectionPlan, error) {
	fillEnvelope(&req.ContractVersion, &req.IdempotencyKey)
	var out SelectionPlan
	return &out, c.do(ctx, http.MethodPost, "/subjects/"+url.PathEscape(subjectID)+"/dataset-profiles/"+url.PathEscape(profileID)+"/triage", req, &out)
}

// ListSelectionPlans calls GET /subjects/{subjectId}/dataset-profiles/{profileId}/selection-plans.
func (c *Client) ListSelectionPlans(ctx context.Context, subjectID, profileID string) (*SelectionPlanList, error) {
	var out SelectionPlanList
	return &out, c.do(ctx, http.MethodGet, "/subjects/"+url.PathEscape(subjectID)+"/dataset-profiles/"+url.PathEscape(profileID)+"/selection-plans", nil, &out)
}

// GetSelectionPlan calls GET /selection-plans/{planId}.
func (c *Client) GetSelectionPlan(ctx context.Context, planID string) (*SelectionPlan, error) {
	var out SelectionPlan
	return &out, c.do(ctx, http.MethodGet, "/selection-plans/"+url.PathEscape(planID), nil, &out)
}

// ReviseSelectionPlan calls POST /selection-plans/{planId}/revisions.
func (c *Client) ReviseSelectionPlan(ctx context.Context, planID string, req ReviseSelectionPlanRequest) (*SelectionPlan, error) {
	fillEnvelope(&req.ContractVersion, &req.IdempotencyKey)
	var out SelectionPlan
	return &out, c.do(ctx, http.MethodPost, "/selection-plans/"+url.PathEscape(planID)+"/revisions", req, &out)
}

// ListAnalyses calls GET /subjects/{subjectId}/analyses.
func (c *Client) ListAnalyses(ctx context.Context, subjectID string) (*AnalysisList, error) {
	var out AnalysisList
	return &out, c.do(ctx, http.MethodGet, "/subjects/"+url.PathEscape(subjectID)+"/analyses", nil, &out)
}

// CompareAnalyses calls GET /subjects/{subjectId}/analyses/{analysisId}/compare/{otherAnalysisId}.
func (c *Client) CompareAnalyses(ctx context.Context, subjectID, analysisID, otherAnalysisID string) (*RunComparisonResult, error) {
	var out RunComparisonResult
	return &out, c.do(ctx, http.MethodGet, "/subjects/"+url.PathEscape(subjectID)+"/analyses/"+url.PathEscape(analysisID)+"/compare/"+url.PathEscape(otherAnalysisID), nil, &out)
}

// ListResearchRuns calls GET /subjects/{subjectId}/research-runs.
func (c *Client) ListResearchRuns(ctx context.Context, subjectID string) (*ResearchRunList, error) {
	var out ResearchRunList
	return &out, c.do(ctx, http.MethodGet, "/subjects/"+url.PathEscape(subjectID)+"/research-runs", nil, &out)
}

// ReEvaluate calls POST /research-runs/{researchRunId}/re-evaluations.
func (c *Client) ReEvaluate(ctx context.Context, researchRunID string, req ReEvaluationRequest) (*ReEvaluationResult, error) {
	fillEnvelope(&req.ContractVersion, &req.IdempotencyKey)
	var out ReEvaluationResult
	return &out, c.do(ctx, http.MethodPost, "/research-runs/"+url.PathEscape(researchRunID)+"/re-evaluations", req, &out)
}

// GetScenarios calls GET /research-runs/{researchRunId}/scenarios.
func (c *Client) GetScenarios(ctx context.Context, researchRunID string) (*ScenarioAnalysis, error) {
	var out ScenarioAnalysis
	return &out, c.do(ctx, http.MethodGet, "/research-runs/"+url.PathEscape(researchRunID)+"/scenarios", nil, &out)
}

// CreateScenarioSet calls POST /research-runs/{researchRunId}/scenario-sets.
func (c *Client) CreateScenarioSet(ctx context.Context, researchRunID string, req CreateScenarioSetRequest) (*ScenarioSetResult, error) {
	fillEnvelope(&req.ContractVersion, &req.IdempotencyKey)
	var out ScenarioSetResult
	return &out, c.do(ctx, http.MethodPost, "/research-runs/"+url.PathEscape(researchRunID)+"/scenario-sets", req, &out)
}

// ScaffoldScenarioSet calls POST /research-runs/{researchRunId}/scenario-sets/scaffold.
func (c *Client) ScaffoldScenarioSet(ctx context.Context, researchRunID string, req ScaffoldScenarioSetRequest) (*ScenarioSetResult, error) {
	fillEnvelope(&req.ContractVersion, &req.IdempotencyKey)
	var out ScenarioSetResult
	return &out, c.do(ctx, http.MethodPost, "/research-runs/"+url.PathEscape(researchRunID)+"/scenario-sets/scaffold", req, &out)
}

// EvaluateScenarios calls POST /research-runs/{researchRunId}/scenario-sets/{scenarioSetId}/evaluations.
func (c *Client) EvaluateScenarios(ctx context.Context, researchRunID, scenarioSetID string, req EvaluateScenariosRequest) (*ScenarioEvaluationResult, error) {
	fillEnvelope(&req.ContractVersion, &req.IdempotencyKey)
	var out ScenarioEvaluationResult
	return &out, c.do(ctx, http.MethodPost, "/research-runs/"+url.PathEscape(researchRunID)+"/scenario-sets/"+url.PathEscape(scenarioSetID)+"/evaluations", req, &out)
}

// GetResearchTimeline calls GET /research-runs/{researchRunId}/timeline.
func (c *Client) GetResearchTimeline(ctx context.Context, researchRunID string) (*ResearchTimeline, error) {
	var out ResearchTimeline
	return &out, c.do(ctx, http.MethodGet, "/research-runs/"+url.PathEscape(researchRunID)+"/timeline", nil, &out)
}

// ApplyTemporalOperation calls POST /temporal-operations.
func (c *Client) ApplyTemporalOperation(ctx context.Context, req TemporalOperationRequest) (*TemporalOperationResult, error) {
	if req.ContractVersion == "" {
		req.ContractVersion = ContractVersion
	}
	var out TemporalOperationResult
	return &out, c.do(ctx, http.MethodPost, "/temporal-operations", req, &out)
}
