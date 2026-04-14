package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pablocaeg/shift-comply/comply"
	_ "github.com/pablocaeg/shift-comply/jurisdictions"
)

// ---------------------------------------------------------------------------
// /health
// ---------------------------------------------------------------------------

func TestHealthEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	handleHealth(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("expected status ok, got %v", body["status"])
	}
}

// ---------------------------------------------------------------------------
// /jurisdictions
// ---------------------------------------------------------------------------

func TestJurisdictionsEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/jurisdictions", nil)
	w := httptest.NewRecorder()
	handleJurisdictions(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Error("expected application/json content type")
	}
	var jurisdictions []comply.JurisdictionDef
	if err := json.Unmarshal(w.Body.Bytes(), &jurisdictions); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(jurisdictions) == 0 {
		t.Error("expected at least one jurisdiction")
	}
}

func TestJurisdictionsEndpoint_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/jurisdictions", nil)
	w := httptest.NewRecorder()
	handleJurisdictions(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

// ---------------------------------------------------------------------------
// /rules
// ---------------------------------------------------------------------------

func TestRulesEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/rules?jurisdiction=US&staff=resident", nil)
	w := httptest.NewRecorder()
	handleRules(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var rules []comply.RuleDef
	if err := json.Unmarshal(w.Body.Bytes(), &rules); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(rules) == 0 {
		t.Error("expected rules for US residents")
	}
}

func TestRulesEndpoint_AllQueryParams(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/rules?jurisdiction=US-CA&staff=nurse-rn&unit=icu&scope=hospitals&category=staffing", nil)
	w := httptest.NewRecorder()
	handleRules(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var rules []comply.RuleDef
	if err := json.Unmarshal(w.Body.Bytes(), &rules); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	for _, r := range rules {
		if r.Category != comply.CatStaffing {
			t.Errorf("expected only staffing rules, got %s", r.Category)
		}
	}
}

func TestRulesEndpoint_MissingJurisdiction(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/rules", nil)
	w := httptest.NewRecorder()
	handleRules(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRulesEndpoint_UnknownJurisdiction(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/rules?jurisdiction=NOWHERE", nil)
	w := httptest.NewRecorder()
	handleRules(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestRulesEndpoint_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/rules?jurisdiction=US", nil)
	w := httptest.NewRecorder()
	handleRules(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

// ---------------------------------------------------------------------------
// /constraints
// ---------------------------------------------------------------------------

func TestConstraintsEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/constraints?jurisdiction=ES&staff=resident", nil)
	w := httptest.NewRecorder()
	handleConstraints(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var constraints []comply.Constraint
	if err := json.Unmarshal(w.Body.Bytes(), &constraints); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(constraints) == 0 {
		t.Error("expected constraints for ES residents")
	}
}

func TestConstraintsEndpoint_MissingJurisdiction(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/constraints", nil)
	w := httptest.NewRecorder()
	handleConstraints(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestConstraintsEndpoint_UnknownJurisdiction(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/constraints?jurisdiction=NOWHERE", nil)
	w := httptest.NewRecorder()
	handleConstraints(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestConstraintsEndpoint_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/constraints?jurisdiction=US", nil)
	w := httptest.NewRecorder()
	handleConstraints(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

// ---------------------------------------------------------------------------
// /compare
// ---------------------------------------------------------------------------

func TestCompareEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/compare?left=US-CA&right=ES", nil)
	w := httptest.NewRecorder()
	handleCompare(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var comp comply.Comparison
	if err := json.Unmarshal(w.Body.Bytes(), &comp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if comp.Left != "US-CA" || comp.Right != "ES" {
		t.Errorf("wrong comparison codes: %s vs %s", comp.Left, comp.Right)
	}
}

func TestCompareEndpoint_WithStaffFilter(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/compare?left=US&right=ES&staff=resident", nil)
	w := httptest.NewRecorder()
	handleCompare(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestCompareEndpoint_MissingParams(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/compare?left=US", nil)
	w := httptest.NewRecorder()
	handleCompare(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCompareEndpoint_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/compare?left=US&right=ES", nil)
	w := httptest.NewRecorder()
	handleCompare(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

// ---------------------------------------------------------------------------
// /validate
// ---------------------------------------------------------------------------

func TestValidateEndpoint(t *testing.T) {
	schedule := comply.Schedule{
		Jurisdiction: "US-CA",
		Shifts: []comply.Shift{
			{StaffID: "nurse-1", StaffType: comply.StaffNurseRN,
				Start: "2025-03-10T07:00:00", End: "2025-03-10T20:30:00"},
		},
	}
	body, _ := json.Marshal(schedule)
	req := httptest.NewRequest(http.MethodPost, "/validate", bytes.NewReader(body))
	w := httptest.NewRecorder()
	handleValidate(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var report comply.ComplianceReport
	if err := json.Unmarshal(w.Body.Bytes(), &report); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if report.Result != "fail" {
		t.Error("13.5h nurse shift in CA should fail")
	}
}

func TestValidateEndpoint_UnknownJurisdiction(t *testing.T) {
	schedule := comply.Schedule{
		Jurisdiction: "NOWHERE",
		Shifts: []comply.Shift{
			{StaffID: "doc-1", StaffType: comply.StaffResident,
				Start: "2025-03-10T08:00:00", End: "2025-03-10T16:00:00"},
		},
	}
	body, _ := json.Marshal(schedule)
	req := httptest.NewRequest(http.MethodPost, "/validate", bytes.NewReader(body))
	w := httptest.NewRecorder()
	handleValidate(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for unknown jurisdiction, got %d", w.Code)
	}
}

func TestValidateEndpoint_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/validate", bytes.NewReader([]byte("not json")))
	w := httptest.NewRecorder()
	handleValidate(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestValidateEndpoint_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/validate", nil)
	w := httptest.NewRecorder()
	handleValidate(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

// ---------------------------------------------------------------------------
// /export
// ---------------------------------------------------------------------------

func TestExportEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/export/US", nil)
	w := httptest.NewRecorder()
	handleExport(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var j comply.JurisdictionDef
	if err := json.Unmarshal(w.Body.Bytes(), &j); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if j.Code != "US" {
		t.Errorf("expected US, got %s", j.Code)
	}
}

func TestExportEndpoint_Unknown(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/export/NOWHERE", nil)
	w := httptest.NewRecorder()
	handleExport(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestExportEndpoint_EmptyCode(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/export/", nil)
	w := httptest.NewRecorder()
	handleExport(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestExportEndpoint_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/export/US", nil)
	w := httptest.NewRecorder()
	handleExport(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

// ---------------------------------------------------------------------------
// CORS
// ---------------------------------------------------------------------------

func TestCORS_Options(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	handler := withCORS(mux)

	req := httptest.NewRequest(http.MethodOptions, "/health", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("missing CORS Allow-Origin header")
	}
	if w.Header().Get("Access-Control-Allow-Methods") != "GET, POST, OPTIONS" {
		t.Error("missing CORS Allow-Methods header")
	}
	if w.Header().Get("Access-Control-Allow-Headers") != "Content-Type" {
		t.Error("missing CORS Allow-Headers header")
	}
	if w.Code != http.StatusNoContent {
		t.Errorf("OPTIONS should return 204, got %d", w.Code)
	}
}

func TestCORS_PassthroughGet(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	handler := withCORS(mux)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("CORS headers should be set on non-OPTIONS requests too")
	}
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
