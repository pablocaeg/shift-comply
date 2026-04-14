package main

import (
	"encoding/json"
	"os/exec"
	"strings"
	"testing"

	"github.com/pablocaeg/shift-comply/comply"
)

func runCLI(t *testing.T, args ...string) (string, string, int) {
	t.Helper()
	cmd := exec.Command("go", append([]string{"run", "."}, args...)...)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("failed to run CLI: %v", err)
		}
	}
	return stdout.String(), stderr.String(), exitCode
}

// ---------------------------------------------------------------------------
// No arguments / help
// ---------------------------------------------------------------------------

func TestCLI_NoArgs(t *testing.T) {
	_, _, code := runCLI(t)
	if code == 0 {
		t.Error("expected non-zero exit code with no arguments")
	}
}

func TestCLI_Help(t *testing.T) {
	stdout, _, code := runCLI(t, "help")
	if code != 0 {
		t.Errorf("help should exit 0, got %d", code)
	}
	if !strings.Contains(stdout, "jurisdictions") {
		t.Error("help should mention jurisdictions command")
	}
}

func TestCLI_UnknownCommand(t *testing.T) {
	_, stderr, code := runCLI(t, "foobar")
	if code == 0 {
		t.Error("unknown command should exit non-zero")
	}
	if !strings.Contains(stderr, "unknown command") {
		t.Errorf("expected 'unknown command' in stderr, got: %s", stderr)
	}
}

// ---------------------------------------------------------------------------
// jurisdictions
// ---------------------------------------------------------------------------

func TestCLI_Jurisdictions(t *testing.T) {
	stdout, _, code := runCLI(t, "jurisdictions")
	if code != 0 {
		t.Errorf("expected exit 0, got %d", code)
	}
	if !strings.Contains(stdout, "US") {
		t.Error("should list US jurisdiction")
	}
	if !strings.Contains(stdout, "ES") {
		t.Error("should list ES jurisdiction")
	}
	if !strings.Contains(stdout, "jurisdictions") {
		t.Error("should show total count")
	}
}

func TestCLI_Jurisdictions_ShortAlias(t *testing.T) {
	stdout, _, code := runCLI(t, "j")
	if code != 0 {
		t.Errorf("expected exit 0, got %d", code)
	}
	if !strings.Contains(stdout, "US") {
		t.Error("short alias 'j' should work")
	}
}

// ---------------------------------------------------------------------------
// rules
// ---------------------------------------------------------------------------

func TestCLI_Rules(t *testing.T) {
	stdout, _, code := runCLI(t, "rules", "US")
	if code != 0 {
		t.Errorf("expected exit 0, got %d", code)
	}
	if !strings.Contains(stdout, "max-weekly-hours") {
		t.Error("should show ACGME max weekly hours rule")
	}
}

func TestCLI_Rules_WithFilters(t *testing.T) {
	stdout, _, code := runCLI(t, "rules", "US-CA", "--staff", "nurse-rn", "--category", "staffing")
	if code != 0 {
		t.Errorf("expected exit 0, got %d", code)
	}
	if !strings.Contains(stdout, "nurse-patient-ratio") {
		t.Error("should show nurse-patient-ratio rules")
	}
}

func TestCLI_Rules_JSON(t *testing.T) {
	stdout, _, code := runCLI(t, "rules", "--json", "US")
	if code != 0 {
		t.Errorf("expected exit 0, got %d", code)
	}
	var rules []comply.RuleDef
	if err := json.Unmarshal([]byte(stdout), &rules); err != nil {
		t.Fatalf("--json should produce valid JSON: %v", err)
	}
	if len(rules) == 0 {
		t.Error("should return rules")
	}
}

func TestCLI_Rules_UnknownJurisdiction(t *testing.T) {
	_, stderr, code := runCLI(t, "rules", "NOWHERE")
	if code == 0 {
		t.Error("unknown jurisdiction should exit non-zero")
	}
	if !strings.Contains(stderr, "unknown jurisdiction") {
		t.Errorf("expected error message, got: %s", stderr)
	}
}

func TestCLI_Rules_MissingArg(t *testing.T) {
	_, stderr, code := runCLI(t, "rules")
	if code == 0 {
		t.Error("missing jurisdiction should exit non-zero")
	}
	if !strings.Contains(stderr, "usage") {
		t.Errorf("expected usage message, got: %s", stderr)
	}
}

// ---------------------------------------------------------------------------
// compare
// ---------------------------------------------------------------------------

func TestCLI_Compare(t *testing.T) {
	stdout, _, code := runCLI(t, "compare", "US-CA", "ES")
	if code != 0 {
		t.Errorf("expected exit 0, got %d", code)
	}
	if !strings.Contains(stdout, "Comparing US-CA vs ES") {
		t.Error("should show comparison header")
	}
	if !strings.Contains(stdout, "Summary") {
		t.Error("should show summary")
	}
}

func TestCLI_Compare_JSON(t *testing.T) {
	stdout, _, code := runCLI(t, "compare", "--json", "US", "ES")
	if code != 0 {
		t.Errorf("expected exit 0, got %d", code)
	}
	var comp comply.Comparison
	if err := json.Unmarshal([]byte(stdout), &comp); err != nil {
		t.Fatalf("--json should produce valid JSON: %v", err)
	}
}

func TestCLI_Compare_MissingArg(t *testing.T) {
	_, stderr, code := runCLI(t, "compare", "US")
	if code == 0 {
		t.Error("missing second arg should exit non-zero")
	}
	if !strings.Contains(stderr, "usage") {
		t.Errorf("expected usage, got: %s", stderr)
	}
}

// ---------------------------------------------------------------------------
// constraints
// ---------------------------------------------------------------------------

func TestCLI_Constraints(t *testing.T) {
	stdout, _, code := runCLI(t, "constraints", "ES", "--staff", "resident")
	if code != 0 {
		t.Errorf("expected exit 0, got %d", code)
	}
	var constraints []comply.Constraint
	if err := json.Unmarshal([]byte(stdout), &constraints); err != nil {
		t.Fatalf("should produce valid JSON: %v", err)
	}
	if len(constraints) == 0 {
		t.Error("should return constraints")
	}
}

func TestCLI_Constraints_UnknownJurisdiction(t *testing.T) {
	_, stderr, code := runCLI(t, "constraints", "NOWHERE")
	if code == 0 {
		t.Error("unknown jurisdiction should exit non-zero")
	}
	if !strings.Contains(stderr, "unknown jurisdiction") {
		t.Errorf("expected error, got: %s", stderr)
	}
}

// ---------------------------------------------------------------------------
// export
// ---------------------------------------------------------------------------

func TestCLI_Export(t *testing.T) {
	stdout, _, code := runCLI(t, "export", "US")
	if code != 0 {
		t.Errorf("expected exit 0, got %d", code)
	}
	var j comply.JurisdictionDef
	if err := json.Unmarshal([]byte(stdout), &j); err != nil {
		t.Fatalf("should produce valid JSON: %v", err)
	}
	if j.Code != "US" {
		t.Errorf("expected US, got %s", j.Code)
	}
}

func TestCLI_Export_Unknown(t *testing.T) {
	_, stderr, code := runCLI(t, "export", "NOWHERE")
	if code == 0 {
		t.Error("unknown jurisdiction should exit non-zero")
	}
	if !strings.Contains(stderr, "unknown jurisdiction") {
		t.Errorf("expected error, got: %s", stderr)
	}
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func TestOpSymbol(t *testing.T) {
	tests := []struct {
		op   comply.Operator
		want string
	}{
		{comply.OpLTE, "<="},
		{comply.OpGTE, ">="},
		{comply.OpEQ, "=="},
		{comply.OpBool, ""},
		{"unknown", "unknown"},
	}
	for _, tt := range tests {
		got := opSymbol(tt.op)
		if got != tt.want {
			t.Errorf("opSymbol(%q) = %q, want %q", tt.op, got, tt.want)
		}
	}
}

func TestJoinKeys(t *testing.T) {
	got := joinKeys([]comply.Key{"a", "b", "c"})
	if got != "a,b,c" {
		t.Errorf("joinKeys = %q, want 'a,b,c'", got)
	}
}
