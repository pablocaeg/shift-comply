package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/pablocaeg/shift-comply/comply"
	_ "github.com/pablocaeg/shift-comply/jurisdictions"
)

func runCmd(args ...string) (stdout, stderr string, code int) {
	var out, err bytes.Buffer
	code = run(args, &out, &err)
	return out.String(), err.String(), code
}

func TestNoArgs(t *testing.T) {
	_, _, code := runCmd()
	if code == 0 {
		t.Error("expected non-zero exit with no args")
	}
}

func TestHelp(t *testing.T) {
	out, _, code := runCmd("help")
	if code != 0 {
		t.Errorf("help should exit 0, got %d", code)
	}
	if !strings.Contains(out, "jurisdictions") {
		t.Error("help should mention jurisdictions")
	}
}

func TestHelpFlags(t *testing.T) {
	for _, f := range []string{"-h", "--help"} {
		out, _, code := runCmd(f)
		if code != 0 {
			t.Errorf("%s should exit 0", f)
		}
		if !strings.Contains(out, "Usage") {
			t.Errorf("%s should show usage", f)
		}
	}
}

func TestUnknownCommand(t *testing.T) {
	_, stderr, code := runCmd("foobar")
	if code == 0 {
		t.Error("unknown command should exit non-zero")
	}
	if !strings.Contains(stderr, "unknown command") {
		t.Errorf("expected 'unknown command', got: %s", stderr)
	}
}

func TestJurisdictions(t *testing.T) {
	out, _, code := runCmd("jurisdictions")
	if code != 0 {
		t.Errorf("expected 0, got %d", code)
	}
	if !strings.Contains(out, "US") || !strings.Contains(out, "ES") {
		t.Error("should list US and ES")
	}
}

func TestJurisdictionsAlias(t *testing.T) {
	out, _, code := runCmd("j")
	if code != 0 {
		t.Errorf("expected 0, got %d", code)
	}
	if !strings.Contains(out, "US") {
		t.Error("alias should work")
	}
}

func TestRules(t *testing.T) {
	out, _, code := runCmd("rules", "US")
	if code != 0 {
		t.Errorf("expected 0, got %d", code)
	}
	if !strings.Contains(out, "max-weekly-hours") {
		t.Error("should show ACGME rule")
	}
}

func TestRulesFilters(t *testing.T) {
	out, _, code := runCmd("rules", "--staff", "nurse-rn", "--category", "staffing", "US-CA")
	if code != 0 {
		t.Errorf("expected 0, got %d", code)
	}
	if !strings.Contains(out, "nurse-patient-ratio") {
		t.Error("should show nurse ratios")
	}
}

func TestRulesJSON(t *testing.T) {
	out, _, code := runCmd("rules", "--json", "US")
	if code != 0 {
		t.Errorf("expected 0, got %d", code)
	}
	var rules []comply.RuleDef
	if err := json.Unmarshal([]byte(out), &rules); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(rules) == 0 {
		t.Error("should return rules")
	}
}

func TestRulesUnknown(t *testing.T) {
	_, stderr, code := runCmd("rules", "NOWHERE")
	if code == 0 {
		t.Error("should fail")
	}
	if !strings.Contains(stderr, "unknown") {
		t.Errorf("expected error, got: %s", stderr)
	}
}

func TestRulesMissing(t *testing.T) {
	_, stderr, code := runCmd("rules")
	if code == 0 {
		t.Error("should fail")
	}
	if !strings.Contains(stderr, "usage") {
		t.Errorf("expected usage, got: %s", stderr)
	}
}

func TestCompare(t *testing.T) {
	out, _, code := runCmd("compare", "US-CA", "ES")
	if code != 0 {
		t.Errorf("expected 0, got %d", code)
	}
	if !strings.Contains(out, "Comparing") {
		t.Error("should show header")
	}
	if !strings.Contains(out, "Summary") {
		t.Error("should show summary")
	}
}

func TestCompareJSON(t *testing.T) {
	out, _, code := runCmd("compare", "--json", "US", "ES")
	if code != 0 {
		t.Errorf("expected 0, got %d", code)
	}
	var comp comply.Comparison
	if err := json.Unmarshal([]byte(out), &comp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestCompareMissing(t *testing.T) {
	_, stderr, code := runCmd("compare", "US")
	if code == 0 {
		t.Error("should fail")
	}
	if !strings.Contains(stderr, "usage") {
		t.Errorf("expected usage, got: %s", stderr)
	}
}

func TestConstraints(t *testing.T) {
	out, _, code := runCmd("constraints", "--staff", "resident", "ES")
	if code != 0 {
		t.Errorf("expected 0, got %d", code)
	}
	var c []comply.Constraint
	if err := json.Unmarshal([]byte(out), &c); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(c) == 0 {
		t.Error("should return constraints")
	}
}

func TestConstraintsUnknown(t *testing.T) {
	_, stderr, code := runCmd("constraints", "NOWHERE")
	if code == 0 {
		t.Error("should fail")
	}
	if !strings.Contains(stderr, "unknown") {
		t.Errorf("expected error, got: %s", stderr)
	}
}

func TestConstraintsMissing(t *testing.T) {
	_, stderr, code := runCmd("constraints")
	if code == 0 {
		t.Error("should fail")
	}
	if !strings.Contains(stderr, "usage") {
		t.Errorf("expected usage, got: %s", stderr)
	}
}

func TestExport(t *testing.T) {
	out, _, code := runCmd("export", "US")
	if code != 0 {
		t.Errorf("expected 0, got %d", code)
	}
	var j comply.JurisdictionDef
	if err := json.Unmarshal([]byte(out), &j); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if j.Code != "US" {
		t.Errorf("expected US, got %s", j.Code)
	}
}

func TestExportUnknown(t *testing.T) {
	_, stderr, code := runCmd("export", "NOWHERE")
	if code == 0 {
		t.Error("should fail")
	}
	if !strings.Contains(stderr, "unknown") {
		t.Errorf("expected error, got: %s", stderr)
	}
}

func TestExportMissing(t *testing.T) {
	_, stderr, code := runCmd("export")
	if code == 0 {
		t.Error("should fail")
	}
	if !strings.Contains(stderr, "usage") {
		t.Errorf("expected usage, got: %s", stderr)
	}
}

func TestOpSymbol(t *testing.T) {
	tests := []struct {
		op   comply.Operator
		want string
	}{
		{comply.OpLTE, "<="}, {comply.OpGTE, ">="}, {comply.OpEQ, "=="}, {comply.OpBool, ""}, {"x", "x"},
	}
	for _, tt := range tests {
		if got := opSymbol(tt.op); got != tt.want {
			t.Errorf("opSymbol(%q) = %q, want %q", tt.op, got, tt.want)
		}
	}
}

func TestJoinKeys(t *testing.T) {
	if got := joinKeys([]comply.Key{"a", "b"}); got != "a,b" {
		t.Errorf("got %q", got)
	}
}
