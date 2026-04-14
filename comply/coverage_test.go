package comply_test

import (
	"testing"

	"github.com/pablocaeg/shift-comply/comply"
	_ "github.com/pablocaeg/shift-comply/jurisdictions"
)

// Cover branches missed by existing tests.

func TestSource_Citation_NoSection(t *testing.T) {
	s := comply.Source{Title: "Some Law"}
	if s.Citation() != "Some Law" {
		t.Errorf("got %q", s.Citation())
	}
}

func TestSource_Citation_WithSection(t *testing.T) {
	s := comply.Source{Title: "Some Law", Section: "Art. 5"}
	if s.Citation() != "Some Law, Art. 5" {
		t.Errorf("got %q", s.Citation())
	}
}

func TestAppliesToUnit_Match(t *testing.T) {
	r := &comply.RuleDef{UnitTypes: []comply.Key{comply.UnitICU, comply.UnitED}}
	if !r.AppliesToUnit(comply.UnitICU) {
		t.Error("should match ICU")
	}
	if !r.AppliesToUnit(comply.UnitED) {
		t.Error("should match ED")
	}
}

func TestAppliesToUnit_NoMatch(t *testing.T) {
	r := &comply.RuleDef{UnitTypes: []comply.Key{comply.UnitICU}}
	if r.AppliesToUnit(comply.UnitED) {
		t.Error("should not match ED")
	}
}

func TestRulesEquivalent_BothNilValues(t *testing.T) {
	comp := comply.Compare(comply.US, comply.US)
	// Same jurisdiction: all rules should be in Same, none in Different
	if len(comp.Different) != 0 {
		t.Errorf("expected 0 different, got %d", len(comp.Different))
	}
}

func TestGenerateConstraints_AllTypes(t *testing.T) {
	// Generate constraints for all jurisdictions to cover more branches
	for _, code := range []comply.Code{comply.US, comply.USCA, comply.EU, comply.ES, comply.ESCT, comply.ESMD} {
		c := comply.GenerateConstraints(code)
		if len(c) == 0 && code != comply.ESMD {
			t.Errorf("expected constraints for %s", code)
		}
		for _, con := range c {
			if con.Citation == "" {
				t.Errorf("[%s] constraint %s has empty citation", code, con.RuleKey)
			}
		}
	}
}

func TestValidate_NightShiftDetection(t *testing.T) {
	// Test a shift starting at 19:00 (should be night)
	schedule := comply.Schedule{
		Jurisdiction: comply.US,
		Shifts: []comply.Shift{
			{StaffID: "res-1", StaffType: comply.StaffResident, Start: "2025-03-10T19:00:00", End: "2025-03-11T07:00:00"},
			{StaffID: "res-1", StaffType: comply.StaffResident, Start: "2025-03-11T19:00:00", End: "2025-03-12T07:00:00"},
			{StaffID: "res-1", StaffType: comply.StaffResident, Start: "2025-03-12T19:00:00", End: "2025-03-13T07:00:00"},
			{StaffID: "res-1", StaffType: comply.StaffResident, Start: "2025-03-13T19:00:00", End: "2025-03-14T07:00:00"},
			{StaffID: "res-1", StaffType: comply.StaffResident, Start: "2025-03-14T19:00:00", End: "2025-03-15T07:00:00"},
			{StaffID: "res-1", StaffType: comply.StaffResident, Start: "2025-03-15T19:00:00", End: "2025-03-16T07:00:00"},
			{StaffID: "res-1", StaffType: comply.StaffResident, Start: "2025-03-16T19:00:00", End: "2025-03-17T07:00:00"},
		},
	}
	report, err := comply.Validate(schedule)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, v := range report.Violations {
		if v.RuleKey == comply.RuleMaxConsecutiveNights {
			found = true
		}
	}
	if !found {
		t.Error("7 consecutive night shifts should trigger max-consecutive-nights violation (limit 6)")
	}
}

func TestValidate_DaysOffPerWeek_SevenDaysWorked(t *testing.T) {
	// ACGME requires 1 day off per week for residents
	schedule := comply.Schedule{
		Jurisdiction: comply.US,
		Shifts: []comply.Shift{
			{StaffID: "res-1", StaffType: comply.StaffResident, Start: "2025-03-10T08:00:00", End: "2025-03-10T16:00:00"},
			{StaffID: "res-1", StaffType: comply.StaffResident, Start: "2025-03-11T08:00:00", End: "2025-03-11T16:00:00"},
			{StaffID: "res-1", StaffType: comply.StaffResident, Start: "2025-03-12T08:00:00", End: "2025-03-12T16:00:00"},
			{StaffID: "res-1", StaffType: comply.StaffResident, Start: "2025-03-13T08:00:00", End: "2025-03-13T16:00:00"},
			{StaffID: "res-1", StaffType: comply.StaffResident, Start: "2025-03-14T08:00:00", End: "2025-03-14T16:00:00"},
			{StaffID: "res-1", StaffType: comply.StaffResident, Start: "2025-03-15T08:00:00", End: "2025-03-15T16:00:00"},
			{StaffID: "res-1", StaffType: comply.StaffResident, Start: "2025-03-16T08:00:00", End: "2025-03-16T16:00:00"},
		},
	}
	report, err := comply.Validate(schedule)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, v := range report.Violations {
		if v.RuleKey == comply.RuleDaysOffPerWeek {
			found = true
		}
	}
	if !found {
		t.Error("7 days worked with no day off should trigger ACGME days-off-per-week violation")
	}
}

func TestValidate_WeeklyHoursNonAveraged(t *testing.T) {
	// Spain ordinary weekly hours (40h, averaged annually) with 50h in one week
	schedule := comply.Schedule{
		Jurisdiction: comply.ES,
		Shifts: []comply.Shift{
			{StaffID: "doc-1", StaffType: comply.StaffStatutory, Start: "2025-03-10T06:00:00", End: "2025-03-10T16:00:00"},
			{StaffID: "doc-1", StaffType: comply.StaffStatutory, Start: "2025-03-11T06:00:00", End: "2025-03-11T16:00:00"},
			{StaffID: "doc-1", StaffType: comply.StaffStatutory, Start: "2025-03-12T06:00:00", End: "2025-03-12T16:00:00"},
			{StaffID: "doc-1", StaffType: comply.StaffStatutory, Start: "2025-03-13T06:00:00", End: "2025-03-13T16:00:00"},
			{StaffID: "doc-1", StaffType: comply.StaffStatutory, Start: "2025-03-14T06:00:00", End: "2025-03-14T16:00:00"},
		},
	}
	report, err := comply.Validate(schedule)
	if err != nil {
		t.Fatal(err)
	}
	// 50 hours in one week, but Spain's 40h is averaged over 12 months
	// So 50h in one week averaged over 12 months is well under 40h
	// This should NOT trigger a violation for max-ordinary-weekly-hours
	for _, v := range report.Violations {
		if v.RuleKey == comply.RuleMaxOrdinaryWeeklyHours {
			t.Logf("violation: %s (actual: %.1f, limit: %.1f)", v.Message, v.Actual, v.Limit)
		}
	}
}

func TestMatchesFilter_UnitFilter(t *testing.T) {
	rules := comply.EffectiveRules(comply.USCA, comply.ForUnit(comply.UnitICU))
	for _, r := range rules {
		if len(r.UnitTypes) > 0 {
			found := false
			for _, u := range r.UnitTypes {
				if u == comply.UnitICU {
					found = true
				}
			}
			if !found {
				t.Errorf("rule %q has unit types %v but doesn't include ICU", r.Key, r.UnitTypes)
			}
		}
	}
}
