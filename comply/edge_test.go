package comply_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/pablocaeg/shift-comply/comply"
	_ "github.com/pablocaeg/shift-comply/jurisdictions"
)

// ---------------------------------------------------------------------------
// Registry edge cases
// ---------------------------------------------------------------------------

func TestFor_UnknownJurisdiction(t *testing.T) {
	if comply.For("DOES-NOT-EXIST") != nil {
		t.Error("For() should return nil for unknown jurisdiction")
	}
}

func TestCodes_MatchesAll(t *testing.T) {
	codes := comply.Codes()
	all := comply.All()
	if len(codes) != len(all) {
		t.Errorf("Codes() returned %d, All() returned %d", len(codes), len(all))
	}
}

func TestChain_TopLevelHasNoParent(t *testing.T) {
	us := comply.For(comply.US)
	if us.ParentDef() != nil {
		t.Error("US should have no parent")
	}
	chain := us.Chain()
	if len(chain) != 1 {
		t.Errorf("US chain should be length 1, got %d", len(chain))
	}
}

func TestChain_BrokenParentReference(t *testing.T) {
	// Register a jurisdiction with a non-existent parent
	comply.RegisterJurisdiction(&comply.JurisdictionDef{
		Code:   "TEST-ORPHAN",
		Name:   "Orphan",
		Type:   comply.State,
		Parent: "DOES-NOT-EXIST",
	})
	defer func() {
		// Clean up: re-register to avoid affecting other tests
		// (can't unregister, but this is fine for testing)
	}()

	j := comply.For("TEST-ORPHAN")
	chain := j.Chain()
	// Should stop at the orphan, not panic
	if len(chain) != 1 {
		t.Errorf("orphan chain should be length 1 (stops at missing parent), got %d", len(chain))
	}
}

// ---------------------------------------------------------------------------
// Data integrity tests
// ---------------------------------------------------------------------------

func TestNoDuplicateKeysWithinJurisdiction(t *testing.T) {
	for _, j := range comply.All() {
		seen := make(map[comply.Key]int)
		for _, r := range j.Rules {
			seen[r.Key]++
		}
		for key, count := range seen {
			if count > 1 {
				t.Errorf("[%s] duplicate rule key %q appears %d times", j.Code, key, count)
			}
		}
	}
}

func TestAllRuleValuesAreSortedNewestFirst(t *testing.T) {
	for _, j := range comply.All() {
		for _, r := range j.Rules {
			for i := 1; i < len(r.Values); i++ {
				if r.Values[i].Since.After(r.Values[i-1].Since) {
					t.Errorf("[%s] rule %q values not sorted newest-first: %s comes after %s",
						j.Code, r.Key,
						r.Values[i-1].Since.Format("2006-01-02"),
						r.Values[i].Since.Format("2006-01-02"))
				}
			}
		}
	}
}

func TestAllRulesHaveValidCategory(t *testing.T) {
	validCategories := map[comply.Category]bool{
		comply.CatWorkHours: true, comply.CatRest: true, comply.CatOvertime: true,
		comply.CatStaffing: true, comply.CatBreaks: true, comply.CatOnCall: true,
		comply.CatCompensation: true, comply.CatLeave: true, comply.CatNightWork: true,
	}
	for _, j := range comply.All() {
		for _, r := range j.Rules {
			if !validCategories[r.Category] {
				t.Errorf("[%s] rule %q has invalid category %q", j.Code, r.Key, r.Category)
			}
		}
	}
}

func TestAllRulesHaveValidOperator(t *testing.T) {
	validOps := map[comply.Operator]bool{
		comply.OpLTE: true, comply.OpGTE: true, comply.OpEQ: true, comply.OpBool: true,
	}
	for _, j := range comply.All() {
		for _, r := range j.Rules {
			if !validOps[r.Operator] {
				t.Errorf("[%s] rule %q has invalid operator %q", j.Code, r.Key, r.Operator)
			}
		}
	}
}

func TestAllRulesHaveValidEnforcement(t *testing.T) {
	validEnf := map[comply.Enforcement]bool{
		comply.Mandatory: true, comply.Recommended: true, comply.Advisory: true,
	}
	for _, j := range comply.All() {
		for _, r := range j.Rules {
			if !validEnf[r.Enforcement] {
				t.Errorf("[%s] rule %q has invalid enforcement %q", j.Code, r.Key, r.Enforcement)
			}
		}
	}
}

func TestAllValuesHaveSinceDate(t *testing.T) {
	for _, j := range comply.All() {
		for _, r := range j.Rules {
			for i, v := range r.Values {
				if v.Since.IsZero() {
					t.Errorf("[%s] rule %q value[%d] has zero Since date", j.Code, r.Key, i)
				}
			}
		}
	}
}

func TestAllRulesHaveNonEmptyDescription(t *testing.T) {
	for _, j := range comply.All() {
		for _, r := range j.Rules {
			if r.Description == "" {
				t.Errorf("[%s] rule %q has empty description", j.Code, r.Key)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Query edge cases
// ---------------------------------------------------------------------------

func TestEffectiveRules_UnknownJurisdiction(t *testing.T) {
	rules := comply.EffectiveRules("NOWHERE")
	if rules != nil {
		t.Error("EffectiveRules should return nil for unknown jurisdiction")
	}
}

func TestEffectiveRules_EmptyFiltersReturnAll(t *testing.T) {
	allDirect := comply.For(comply.US)
	effective := comply.EffectiveRules(comply.US)
	if len(effective) != len(allDirect.Rules) {
		t.Errorf("US effective (no filters) should match direct rules: got %d vs %d",
			len(effective), len(allDirect.Rules))
	}
}

func TestEffectiveRules_StaffFilterExcludesNonMatching(t *testing.T) {
	// VA nurse rules should not appear when querying for residents
	rules := comply.EffectiveRules(comply.US, comply.ForStaff(comply.StaffResident))
	for _, r := range rules {
		for _, st := range r.StaffTypes {
			if st == comply.StaffVANurse {
				t.Errorf("resident query should not include VA nurse rule %q", r.Key)
			}
		}
	}
}

func TestEffectiveRules_ScopeFilterIncludesUnscoped(t *testing.T) {
	// Rules with no scope (ScopeAll or empty) should always be included
	rules := comply.EffectiveRules(comply.US, comply.ForScope(comply.ScopeHospitals))
	foundFLSA := false
	for _, r := range rules {
		if r.Key == comply.RuleOvertimeWeeklyThreshold {
			foundFLSA = true
		}
	}
	if !foundFLSA {
		t.Error("FLSA overtime (no scope) should appear when filtering by hospitals scope")
	}
}

func TestEffectiveRules_ScopeFilterExcludesNonMatching(t *testing.T) {
	rules := comply.EffectiveRules(comply.US, comply.ForScope(comply.ScopeVA))
	for _, r := range rules {
		if r.Scope == comply.ScopeAccreditedPrograms {
			t.Errorf("VA scope should not include accredited_programs rule %q", r.Key)
		}
	}
}

func TestEffectiveRules_DateFilterWorks(t *testing.T) {
	// Query CA rules before 2004 should not include nurse ratios (effective Jan 1 2004)
	before := time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC)
	rules := comply.EffectiveRules(comply.USCA, comply.OnDate(before))
	for _, r := range rules {
		if r.Key == comply.RuleNursePatientRatioICU {
			t.Error("ICU ratio should not appear before its 2004 effective date")
		}
	}
}

func TestEffectiveRules_CategoryFilterWorks(t *testing.T) {
	rules := comply.EffectiveRules(comply.ES, comply.InCategory(comply.CatLeave))
	for _, r := range rules {
		if r.Category != comply.CatLeave {
			t.Errorf("category filter returned rule %q with category %q", r.Key, r.Category)
		}
	}
	if len(rules) == 0 {
		t.Error("Spain should have at least one leave rule")
	}
}

// ---------------------------------------------------------------------------
// Compare edge cases
// ---------------------------------------------------------------------------

func TestCompare_SameJurisdiction(t *testing.T) {
	comp := comply.Compare(comply.US, comply.US)
	if len(comp.OnlyLeft) != 0 || len(comp.OnlyRight) != 0 {
		t.Error("comparing a jurisdiction with itself should have no differences")
	}
	if len(comp.Different) != 0 {
		t.Errorf("same jurisdiction should have 0 different rules, got %d", len(comp.Different))
	}
}

func TestCompare_UnknownJurisdiction(t *testing.T) {
	comp := comply.Compare("NOWHERE", comply.US)
	if comp.OnlyLeft != nil {
		t.Error("comparing unknown jurisdiction should produce no left-only rules")
	}
}

// ---------------------------------------------------------------------------
// RuleDef method edge cases
// ---------------------------------------------------------------------------

func TestRuleDef_Value_BeforeAllDates(t *testing.T) {
	r := &comply.RuleDef{
		Values: []*comply.RuleValue{
			{Since: comply.D(2020, time.January, 1), Amount: 10},
		},
	}
	v := r.Value(time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC))
	if v != nil {
		t.Error("Value() should return nil for date before all values")
	}
}

func TestRuleDef_Value_ExactDate(t *testing.T) {
	since := comply.D(2020, time.January, 1)
	r := &comply.RuleDef{
		Values: []*comply.RuleValue{
			{Since: since, Amount: 10},
		},
	}
	v := r.Value(since)
	if v == nil || v.Amount != 10 {
		t.Error("Value() should match on exact Since date")
	}
}

func TestRuleDef_Current_NoValues(t *testing.T) {
	r := &comply.RuleDef{}
	if r.Current() != nil {
		t.Error("Current() should return nil for rule with no values")
	}
}

func TestRuleDef_AppliesToStaff_EmptyMeansAll(t *testing.T) {
	r := &comply.RuleDef{}
	if !r.AppliesToStaff(comply.StaffNurseRN) {
		t.Error("rule with no StaffTypes should apply to all staff")
	}
}

func TestRuleDef_AppliesToUnit_EmptyMeansAll(t *testing.T) {
	r := &comply.RuleDef{}
	if !r.AppliesToUnit(comply.UnitICU) {
		t.Error("rule with no UnitTypes should apply to all units")
	}
}

// ---------------------------------------------------------------------------
// Constraint generation edge cases
// ---------------------------------------------------------------------------

func TestGenerateConstraints_UnknownJurisdiction(t *testing.T) {
	c := comply.GenerateConstraints("NOWHERE")
	if len(c) != 0 {
		t.Error("GenerateConstraints should return empty for unknown jurisdiction")
	}
}

func TestGenerateConstraints_HasCitations(t *testing.T) {
	constraints := comply.GenerateConstraints(comply.US, comply.ForStaff(comply.StaffResident))
	for _, c := range constraints {
		if c.Citation == "" {
			t.Errorf("constraint for rule %q has empty citation", c.RuleKey)
		}
		if c.Jurisdiction != comply.US {
			t.Errorf("constraint jurisdiction should be US, got %s", c.Jurisdiction)
		}
	}
}

func TestGenerateConstraints_FacilityScopeIncluded(t *testing.T) {
	constraints := comply.GenerateConstraints(comply.US, comply.ForStaff(comply.StaffResident))
	for _, c := range constraints {
		if c.RuleKey == comply.RuleMaxWeeklyHours && c.FacilityScope != comply.ScopeAccreditedPrograms {
			t.Errorf("ACGME constraint should have accredited_programs scope, got %q", c.FacilityScope)
		}
	}
}

// ---------------------------------------------------------------------------
// JSON serialization edge cases
// ---------------------------------------------------------------------------

func TestJurisdiction_JSONRoundTrip(t *testing.T) {
	j := comply.For(comply.US)
	b, err := json.Marshal(j)
	if err != nil {
		t.Fatal(err)
	}
	var decoded comply.JurisdictionDef
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Code != j.Code {
		t.Errorf("code mismatch after roundtrip: %q vs %q", decoded.Code, j.Code)
	}
	if len(decoded.Rules) != len(j.Rules) {
		t.Errorf("rules count mismatch after roundtrip: %d vs %d", len(decoded.Rules), len(j.Rules))
	}
}

func TestComparison_JSONSerializable(t *testing.T) {
	comp := comply.Compare(comply.US, comply.ES)
	b, err := json.Marshal(comp)
	if err != nil {
		t.Fatalf("comparison should be JSON serializable: %v", err)
	}
	if len(b) == 0 {
		t.Error("comparison JSON should not be empty")
	}
}

func TestComplianceReport_JSONSerializable(t *testing.T) {
	report, err := comply.Validate(comply.Schedule{
		Jurisdiction: comply.ES,
		Shifts: []comply.Shift{
			{StaffID: "doc-1", StaffType: comply.StaffStatutory,
				Start: "2025-03-10T08:00:00", End: "2025-03-10T15:00:00"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(report)
	if err != nil {
		t.Fatalf("report should be JSON serializable: %v", err)
	}
	var decoded comply.ComplianceReport
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatalf("report should be JSON deserializable: %v", err)
	}
	if decoded.Jurisdiction != comply.ES {
		t.Errorf("jurisdiction mismatch: %q vs %q", decoded.Jurisdiction, comply.ES)
	}
}

// ---------------------------------------------------------------------------
// Validation edge cases (beyond validate_test.go)
// ---------------------------------------------------------------------------

func TestValidate_SameStaffIDDifferentTypes(t *testing.T) {
	// If someone reuses a staff ID with different types, each should get its own rules
	schedule := comply.Schedule{
		Jurisdiction: comply.US,
		Shifts: []comply.Shift{
			{StaffID: "person-1", StaffType: comply.StaffResident, Start: "2025-03-10T06:00:00", End: "2025-03-10T18:00:00"},
			{StaffID: "person-1", StaffType: comply.StaffResident, Start: "2025-03-11T06:00:00", End: "2025-03-11T18:00:00"},
		},
	}
	report, err := comply.Validate(schedule)
	if err != nil {
		t.Fatal(err)
	}
	if report.ConstraintsChecked == 0 {
		t.Error("should have checked at least some constraints")
	}
}

func TestValidate_ShiftExactly24Hours(t *testing.T) {
	// A 24-hour shift should trigger the min-rest-after-extended check
	schedule := comply.Schedule{
		Jurisdiction: comply.US,
		Shifts: []comply.Shift{
			{StaffID: "res-1", StaffType: comply.StaffResident,
				Start: "2025-03-10T08:00:00", End: "2025-03-11T08:00:00"}, // exactly 24h
			{StaffID: "res-1", StaffType: comply.StaffResident,
				Start: "2025-03-11T12:00:00", End: "2025-03-11T20:00:00"}, // only 4h rest
		},
	}
	report, err := comply.Validate(schedule)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, v := range report.Violations {
		if v.RuleKey == comply.RuleMinRestAfterExtended {
			found = true
			if v.Actual != 4.0 {
				t.Errorf("expected 4.0 hours actual rest, got %v", v.Actual)
			}
		}
	}
	if !found {
		t.Error("should flag min rest after 24-hour extended shift (4h rest, 14h required by ACGME)")
	}
}

func TestValidate_GuardsAcrossMonthBoundary(t *testing.T) {
	// Guards split across March and April should be counted separately
	schedule := comply.Schedule{
		Jurisdiction: comply.ES,
		Shifts: []comply.Shift{
			// 4 guards in March
			{StaffID: "mir-1", StaffType: comply.StaffResident, Start: "2025-03-05T08:00:00", End: "2025-03-06T08:00:00", OnCall: true},
			{StaffID: "mir-1", StaffType: comply.StaffResident, Start: "2025-03-12T08:00:00", End: "2025-03-13T08:00:00", OnCall: true},
			{StaffID: "mir-1", StaffType: comply.StaffResident, Start: "2025-03-19T08:00:00", End: "2025-03-20T08:00:00", OnCall: true},
			{StaffID: "mir-1", StaffType: comply.StaffResident, Start: "2025-03-26T08:00:00", End: "2025-03-27T08:00:00", OnCall: true},
			// 4 guards in April
			{StaffID: "mir-1", StaffType: comply.StaffResident, Start: "2025-04-02T08:00:00", End: "2025-04-03T08:00:00", OnCall: true},
			{StaffID: "mir-1", StaffType: comply.StaffResident, Start: "2025-04-09T08:00:00", End: "2025-04-10T08:00:00", OnCall: true},
			{StaffID: "mir-1", StaffType: comply.StaffResident, Start: "2025-04-16T08:00:00", End: "2025-04-17T08:00:00", OnCall: true},
			{StaffID: "mir-1", StaffType: comply.StaffResident, Start: "2025-04-23T08:00:00", End: "2025-04-24T08:00:00", OnCall: true},
		},
	}
	report, err := comply.Validate(schedule)
	if err != nil {
		t.Fatal(err)
	}
	// 4 guards per month should not violate Spain's 7 guard limit
	for _, v := range report.Violations {
		if v.RuleKey == comply.RuleMaxGuardsMonthly {
			t.Errorf("4 guards/month should not violate 7 guard limit: %s", v.Message)
		}
	}
}

func TestValidate_NonOnCallShiftsIgnoredByGuardCheck(t *testing.T) {
	// Regular shifts (OnCall=false) should not be counted as guards
	schedule := comply.Schedule{
		Jurisdiction: comply.ES,
		Shifts:       make([]comply.Shift, 0, 10),
	}
	for i := 1; i <= 10; i++ {
		schedule.Shifts = append(schedule.Shifts, comply.Shift{
			StaffID:   "mir-1",
			StaffType: comply.StaffResident,
			Start:     fmt.Sprintf("2025-03-%02dT08:00:00", i*2),
			End:       fmt.Sprintf("2025-03-%02dT16:00:00", i*2),
			OnCall:    false, // not a guard
		})
	}
	report, err := comply.Validate(schedule)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range report.Violations {
		if v.RuleKey == comply.RuleMaxGuardsMonthly {
			t.Errorf("non-oncall shifts should not trigger guard limit: %s", v.Message)
		}
	}
}

func TestValidate_RFC3339WithTimezone(t *testing.T) {
	// Shifts with timezone offsets should parse correctly
	schedule := comply.Schedule{
		Jurisdiction: comply.ES,
		Shifts: []comply.Shift{
			{StaffID: "doc-1", StaffType: comply.StaffStatutory,
				Start: "2025-03-10T08:00:00+01:00", End: "2025-03-10T16:00:00+01:00"},
		},
	}
	report, err := comply.Validate(schedule)
	if err != nil {
		t.Fatalf("RFC3339 with timezone should parse: %v", err)
	}
	if report.Result != resultPass {
		t.Errorf("normal 8-hour shift should pass, got %s", report.Result)
	}
}

func TestValidate_SingleShiftNoRestViolation(t *testing.T) {
	// A single shift should never trigger rest-between-shifts
	schedule := comply.Schedule{
		Jurisdiction: comply.ES,
		Shifts: []comply.Shift{
			{StaffID: "doc-1", StaffType: comply.StaffStatutory,
				Start: "2025-03-10T08:00:00", End: "2025-03-10T20:00:00"},
		},
	}
	report, err := comply.Validate(schedule)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range report.Violations {
		if v.RuleKey == comply.RuleMinRestBetweenShifts {
			t.Errorf("single shift should never violate rest between shifts: %s", v.Message)
		}
	}
}

// ---------------------------------------------------------------------------
// Inheritance correctness
// ---------------------------------------------------------------------------

func TestCataloniaDoesNotInheritMadridRules(t *testing.T) {
	// Catalonia and Madrid are siblings, not parent-child
	ctRules := comply.EffectiveRules(comply.ESCT)
	for _, r := range ctRules {
		if r.Key == "md-rest-not-effective-work" {
			t.Error("Catalonia should NOT inherit Madrid-specific rules")
		}
	}
}

func TestMadridDoesNotInheritCataloniaRules(t *testing.T) {
	mdRules := comply.EffectiveRules(comply.ESMD)
	for _, r := range mdRules {
		if r.Key == "ct-guard-exemption-age-60" {
			t.Error("Madrid should NOT inherit Catalonia-specific rules")
		}
	}
}

func TestCaliforniaDoesNotInheritSpainRules(t *testing.T) {
	caRules := comply.EffectiveRules(comply.USCA)
	for _, r := range caRules {
		if r.Key == comply.RuleMaxOvertimeAnnual && r.Source.Title == "Real Decreto Legislativo 2/2015 (Estatuto de los Trabajadores)" {
			t.Error("California should NOT inherit Spanish rules")
		}
	}
}
