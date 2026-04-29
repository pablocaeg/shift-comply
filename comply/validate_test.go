package comply_test

import (
	"fmt"
	"testing"

	"github.com/pablocaeg/shift-comply/comply"
	_ "github.com/pablocaeg/shift-comply/jurisdictions"
)

const resultPass = "pass"
const resultFail = "fail"

func TestValidate_MaxShiftHours(t *testing.T) {
	schedule := comply.Schedule{
		Jurisdiction: comply.USCA,
		Shifts: []comply.Shift{
			{StaffID: "nurse-1", StaffType: comply.StaffNurseRN, Start: "2025-03-10T07:00:00", End: "2025-03-10T20:30:00"}, // 13.5 hours, exceeds 12
		},
	}
	report, err := comply.Validate(schedule)
	if err != nil {
		t.Fatal(err)
	}
	if report.Result != resultFail {
		t.Error("expected fail, 13.5 hour shift exceeds CA 12-hour max")
	}
	found := false
	for _, v := range report.Violations {
		if v.StaffID == "nurse-1" && v.Limit == 12 && v.Actual >= 13 {
			found = true
		}
	}
	if !found {
		t.Error("expected max shift hours violation for nurse-1")
	}
}

func TestValidate_ShiftWithinLimits(t *testing.T) {
	schedule := comply.Schedule{
		Jurisdiction: comply.USCA,
		Shifts: []comply.Shift{
			{StaffID: "nurse-1", StaffType: comply.StaffNurseRN, Start: "2025-03-10T07:00:00", End: "2025-03-10T19:00:00"}, // 12 hours, exactly at limit
		},
	}
	report, err := comply.Validate(schedule)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range report.Violations {
		if v.Limit == 12 && v.Actual > 12 {
			t.Errorf("12-hour shift should not violate 12-hour max, got violation: %s", v.Message)
		}
	}
}

func TestValidate_MinRestBetweenShifts(t *testing.T) {
	schedule := comply.Schedule{
		Jurisdiction:  comply.ES,
		FacilityScope: comply.ScopePublicHealth,
		Shifts: []comply.Shift{
			{StaffID: "doc-1", StaffType: comply.StaffStatutory, Start: "2025-03-10T08:00:00", End: "2025-03-10T20:00:00"},
			{StaffID: "doc-1", StaffType: comply.StaffStatutory, Start: "2025-03-11T04:00:00", End: "2025-03-11T16:00:00"}, // only 8 hours rest, Spain requires 12
		},
	}
	report, err := comply.Validate(schedule)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, v := range report.Violations {
		if v.RuleKey == comply.RuleMinRestBetweenShifts {
			found = true
			if v.Actual != 8.0 {
				t.Errorf("expected 8.0 hours actual rest, got %v", v.Actual)
			}
		}
	}
	if !found {
		t.Error("expected min-rest-between-shifts violation (8h rest, 12h required)")
	}
}

func TestValidate_SpainResidentGuards(t *testing.T) {
	// Spain MIR: max 7 guards per month
	shifts := make([]comply.Shift, 0, 8)
	for i := 1; i <= 8; i++ {
		shifts = append(shifts, comply.Shift{
			StaffID:   "mir-1",
			StaffType: comply.StaffResident,
			Start:     "2025-03-" + pad(i*3) + "T08:00:00",
			End:       "2025-03-" + pad(i*3+1) + "T08:00:00",
			OnCall:    true,
		})
	}
	schedule := comply.Schedule{
		Jurisdiction: comply.ES,
		Shifts:       shifts,
	}
	report, err := comply.Validate(schedule)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, v := range report.Violations {
		if v.RuleKey == comply.RuleMaxGuardsMonthly {
			found = true
			if v.Actual != 8 {
				t.Errorf("expected 8 guards, got %v", v.Actual)
			}
		}
	}
	if !found {
		t.Error("expected max-guards-monthly violation (8 guards, max 7)")
	}
}

func TestValidate_PassingSchedule(t *testing.T) {
	schedule := comply.Schedule{
		Jurisdiction: comply.ES,
		Shifts: []comply.Shift{
			{StaffID: "doc-1", StaffType: comply.StaffStatutory, Start: "2025-03-10T08:00:00", End: "2025-03-10T15:00:00"},
			{StaffID: "doc-1", StaffType: comply.StaffStatutory, Start: "2025-03-11T08:00:00", End: "2025-03-11T15:00:00"},
			{StaffID: "doc-1", StaffType: comply.StaffStatutory, Start: "2025-03-12T08:00:00", End: "2025-03-12T15:00:00"},
		},
	}
	report, err := comply.Validate(schedule)
	if err != nil {
		t.Fatal(err)
	}
	if report.Result != resultPass {
		t.Errorf("expected pass for normal 7-hour shifts with 17-hour rest, got %s with %d violations",
			report.Result, len(report.Violations))
		for _, v := range report.Violations {
			t.Logf("  violation: %s - %s", v.RuleKey, v.Message)
		}
	}
}

func TestValidate_UnknownJurisdiction(t *testing.T) {
	schedule := comply.Schedule{
		Jurisdiction: "XX",
		Shifts:       []comply.Shift{},
	}
	_, err := comply.Validate(schedule)
	if err == nil {
		t.Error("expected error for unknown jurisdiction")
	}
}

func TestValidate_CataloniaOverridesGuardLimit(t *testing.T) {
	// Catalonia limits guards to 4/month, overriding Spain's 7
	shifts := make([]comply.Shift, 0, 5)
	for i := 1; i <= 5; i++ {
		shifts = append(shifts, comply.Shift{
			StaffID:   "doc-ct-1",
			StaffType: comply.StaffStatutory,
			Start:     "2025-03-" + pad(i*5) + "T08:00:00",
			End:       "2025-03-" + pad(i*5+1) + "T08:00:00",
			OnCall:    true,
		})
	}
	schedule := comply.Schedule{
		Jurisdiction: comply.ESCT,
		Shifts:       shifts,
	}
	report, err := comply.Validate(schedule)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, v := range report.Violations {
		if v.RuleKey == comply.RuleMaxGuardsMonthly {
			found = true
			if v.Limit != 4 {
				t.Errorf("Catalonia guard limit should be 4, got %v", v.Limit)
			}
		}
	}
	if !found {
		t.Error("expected max-guards-monthly violation in Catalonia (5 guards, max 4)")
	}
}

func TestValidate_ACGMEWeeklyHours(t *testing.T) {
	// ACGME: max 80 hours/week averaged over 4 weeks
	// Schedule 6 days/week * 14h/day = 84h/week for 4 weeks = 336h / 4 = 84 avg > 80
	var shifts []comply.Shift
	for week := 0; week < 4; week++ {
		baseDay := 3 + (week * 7) // start March 3
		for d := 0; d < 6; d++ {
			day := baseDay + d
			shifts = append(shifts, comply.Shift{
				StaffID:   "res-1",
				StaffType: comply.StaffResident,
				Start:     fmt.Sprintf("2025-03-%02dT06:00:00", day),
				End:       fmt.Sprintf("2025-03-%02dT20:00:00", day), // 14h * 6 days = 84h/week
			})
		}
	}

	schedule := comply.Schedule{
		Jurisdiction: comply.US,
		Shifts:       shifts,
	}
	report, err := comply.Validate(schedule)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, v := range report.Violations {
		if v.RuleKey == comply.RuleMaxWeeklyHours {
			found = true
		}
	}
	if !found {
		t.Error("expected max-weekly-hours violation (84 hours/week averaged over 4 weeks, ACGME max 80)")
	}
}

func TestValidate_OverlappingShifts(t *testing.T) {
	schedule := comply.Schedule{
		Jurisdiction: comply.ES,
		Shifts: []comply.Shift{
			{StaffID: "doc-1", StaffType: comply.StaffStatutory, Start: "2025-03-10T08:00:00", End: "2025-03-10T20:00:00"},
			{StaffID: "doc-1", StaffType: comply.StaffStatutory, Start: "2025-03-10T18:00:00", End: "2025-03-11T06:00:00"}, // overlaps by 2 hours
		},
	}
	report, err := comply.Validate(schedule)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, v := range report.Violations {
		if v.RuleKey == comply.RuleMinRestBetweenShifts && v.Actual == 0 {
			found = true
		}
	}
	if !found {
		t.Error("expected overlapping shifts to be flagged as zero rest violation")
	}
}

func TestValidate_ShiftEndBeforeStart(t *testing.T) {
	schedule := comply.Schedule{
		Jurisdiction: comply.US,
		Shifts: []comply.Shift{
			{StaffID: "doc-1", StaffType: comply.StaffResident, Start: "2025-03-10T20:00:00", End: "2025-03-10T08:00:00"},
		},
	}
	_, err := comply.Validate(schedule)
	if err == nil {
		t.Error("expected error for shift where end is before start")
	}
}

func TestValidate_EmptySchedule(t *testing.T) {
	schedule := comply.Schedule{
		Jurisdiction: comply.US,
		Shifts:       []comply.Shift{},
	}
	report, err := comply.Validate(schedule)
	if err != nil {
		t.Fatal(err)
	}
	if report.Result != resultPass {
		t.Errorf("empty schedule should pass, got %s", report.Result)
	}
}

func TestValidate_MidnightBoundaryShift(t *testing.T) {
	// Shift ending exactly at midnight should not count the next day as worked
	schedule := comply.Schedule{
		Jurisdiction: comply.ES,
		Shifts: []comply.Shift{
			{StaffID: "doc-1", StaffType: comply.StaffStatutory, Start: "2025-03-10T16:00:00", End: "2025-03-11T00:00:00"}, // 8 hours, ends at midnight
			{StaffID: "doc-1", StaffType: comply.StaffStatutory, Start: "2025-03-11T12:00:00", End: "2025-03-11T20:00:00"}, // next day
		},
	}
	report, err := comply.Validate(schedule)
	if err != nil {
		t.Fatal(err)
	}
	// Rest between shifts: midnight to 12:00 = 12 hours. Spain requires 12, so this should pass.
	for _, v := range report.Violations {
		if v.RuleKey == comply.RuleMinRestBetweenShifts {
			t.Errorf("12 hours rest (midnight to noon) should satisfy Spain's 12-hour minimum, got violation: %s", v.Message)
		}
	}
}

func TestValidate_MultipleStaffTypes(t *testing.T) {
	// Two staff members with different types should get different rules
	schedule := comply.Schedule{
		Jurisdiction: comply.USCA,
		Shifts: []comply.Shift{
			// Nurse: 13-hour shift violates CA 12-hour max
			{StaffID: "nurse-1", StaffType: comply.StaffNurseRN, Start: "2025-03-10T06:00:00", End: "2025-03-10T19:30:00"},
			// Resident: 13-hour shift does NOT violate (ACGME allows 24)
			{StaffID: "res-1", StaffType: comply.StaffResident, Start: "2025-03-10T06:00:00", End: "2025-03-10T19:30:00"},
		},
	}
	report, err := comply.Validate(schedule)
	if err != nil {
		t.Fatal(err)
	}
	nurseViolation := false
	residentViolation := false
	for _, v := range report.Violations {
		// Check for shift-duration violations specifically (limit 12 or 24)
		if v.StaffID == "nurse-1" && v.Limit == 12 && v.Actual >= 13 {
			nurseViolation = true
		}
		if v.StaffID == "res-1" && v.Limit == 12 {
			residentViolation = true
		}
	}
	if !nurseViolation {
		t.Error("nurse 13-hour shift should violate CA 12-hour max")
	}
	if residentViolation {
		t.Error("resident 13-hour shift should NOT violate (ACGME allows 24 hours)")
	}
}

func TestValidate_InvalidTimeFormat(t *testing.T) {
	schedule := comply.Schedule{
		Jurisdiction: comply.US,
		Shifts: []comply.Shift{
			{StaffID: "doc-1", StaffType: comply.StaffResident, Start: "March 10 2025", End: "March 10 2025"},
		},
	}
	_, err := comply.Validate(schedule)
	if err == nil {
		t.Error("expected error for invalid time format")
	}
}

func TestValidate_ExactlyAtLimit(t *testing.T) {
	// Spain: exactly 12 hours rest should NOT violate the 12-hour minimum
	schedule := comply.Schedule{
		Jurisdiction: comply.ES,
		Shifts: []comply.Shift{
			{StaffID: "doc-1", StaffType: comply.StaffStatutory, Start: "2025-03-10T08:00:00", End: "2025-03-10T16:00:00"},
			{StaffID: "doc-1", StaffType: comply.StaffStatutory, Start: "2025-03-11T04:00:00", End: "2025-03-11T12:00:00"}, // exactly 12h rest
		},
	}
	report, err := comply.Validate(schedule)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range report.Violations {
		if v.RuleKey == comply.RuleMinRestBetweenShifts {
			t.Errorf("exactly 12 hours rest should not violate 12-hour minimum, got: %s", v.Message)
		}
	}
}

// ---------------------------------------------------------------------------
// ValidateSwap tests
// ---------------------------------------------------------------------------

func TestValidateSwap_TakeShiftViolatesRest(t *testing.T) {
	// A nurse already has a shift ending at 20:00. Taking a shift starting at
	// 04:00 the next day gives only 8 hours rest — Spain requires 12.
	req := comply.SwapRequest{
		Jurisdiction:  comply.ES,
		FacilityScope: comply.ScopePublicHealth,
		StaffID:       "nurse-1",
		StaffType:     comply.StaffStatutory,
		CurrentShifts: []comply.Shift{
			{Start: "2025-03-10T08:00:00", End: "2025-03-10T20:00:00"},
		},
		Add: &comply.Shift{Start: "2025-03-11T04:00:00", End: "2025-03-11T12:00:00"},
	}
	report, err := comply.ValidateSwap(req)
	if err != nil {
		t.Fatal(err)
	}
	if report.Result != resultFail {
		t.Error("expected fail: taking shift creates 8h rest gap, Spain requires 12h")
	}
	found := false
	for _, v := range report.Violations {
		if v.RuleKey == comply.RuleMinRestBetweenShifts {
			found = true
		}
	}
	if !found {
		t.Error("expected min-rest-between-shifts violation")
	}
}

func TestValidateSwap_SwapIsCompliant(t *testing.T) {
	// A nurse swaps a morning shift for an afternoon shift the next day.
	// Removing the 08-16 shift and adding 16-00 the next day is compliant.
	req := comply.SwapRequest{
		Jurisdiction: comply.ES,
		StaffID:      "nurse-1",
		StaffType:    comply.StaffStatutory,
		CurrentShifts: []comply.Shift{
			{Start: "2025-03-10T08:00:00", End: "2025-03-10T16:00:00"},
			{Start: "2025-03-11T08:00:00", End: "2025-03-11T16:00:00"},
			{Start: "2025-03-12T08:00:00", End: "2025-03-12T16:00:00"},
		},
		Remove: &comply.Shift{Start: "2025-03-11T08:00:00", End: "2025-03-11T16:00:00"},
		Add:    &comply.Shift{Start: "2025-03-13T08:00:00", End: "2025-03-13T16:00:00"},
	}
	report, err := comply.ValidateSwap(req)
	if err != nil {
		t.Fatal(err)
	}
	if report.Result != resultPass {
		t.Errorf("expected pass, got %s with violations:", report.Result)
		for _, v := range report.Violations {
			t.Logf("  %s: %s", v.RuleKey, v.Message)
		}
	}
}

func TestValidateSwap_TakeShiftExceedsWeeklyHours(t *testing.T) {
	// A resident working 84h/week for 4 weeks (above ACGME 80h avg).
	// They swap out a short shift and take a longer one, pushing avg over 80.
	// Build 4 weeks of 6 days * 14h = 84h/week each.
	var shifts []comply.Shift
	for week := 0; week < 4; week++ {
		baseDay := 3 + (week * 7) // March 3, 10, 17, 24
		for d := 0; d < 6; d++ {
			day := baseDay + d
			shifts = append(shifts, comply.Shift{
				Start: fmt.Sprintf("2025-03-%02dT06:00:00", day),
				End:   fmt.Sprintf("2025-03-%02dT19:30:00", day), // 13.5h * 6 = 81h/week
			})
		}
	}
	// Take an additional 4h shift in week 4, pushing avg even higher.
	req := comply.SwapRequest{
		Jurisdiction:  comply.US,
		StaffID:       "res-1",
		StaffType:     comply.StaffResident,
		CurrentShifts: shifts,
		Add:           &comply.Shift{Start: "2025-03-30T08:00:00", End: "2025-03-30T12:00:00"},
	}
	report, err := comply.ValidateSwap(req)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, v := range report.Violations {
		if v.RuleKey == comply.RuleMaxWeeklyHours {
			found = true
		}
	}
	if !found {
		t.Error("expected max-weekly-hours violation after taking extra shift (>80h/week avg over 4 weeks)")
	}
}

func TestValidateSwap_RemoveOnly(t *testing.T) {
	// Giving away a shift (no replacement) should always be compliant.
	req := comply.SwapRequest{
		Jurisdiction: comply.USCA,
		StaffID:      "nurse-1",
		StaffType:    comply.StaffNurseRN,
		CurrentShifts: []comply.Shift{
			{Start: "2025-03-10T07:00:00", End: "2025-03-10T19:00:00"},
			{Start: "2025-03-11T07:00:00", End: "2025-03-11T19:00:00"},
		},
		Remove: &comply.Shift{Start: "2025-03-11T07:00:00", End: "2025-03-11T19:00:00"},
	}
	report, err := comply.ValidateSwap(req)
	if err != nil {
		t.Fatal(err)
	}
	if report.Result != resultPass {
		t.Errorf("giving away a shift should always pass, got %s", report.Result)
	}
}

func TestValidateSwap_UnknownJurisdiction(t *testing.T) {
	req := comply.SwapRequest{
		Jurisdiction: "XX",
		StaffID:      "doc-1",
		StaffType:    comply.StaffResident,
	}
	_, err := comply.ValidateSwap(req)
	if err == nil {
		t.Error("expected error for unknown jurisdiction")
	}
}

// ---------------------------------------------------------------------------
// Derived constraint tests
// ---------------------------------------------------------------------------

func TestGenerateConstraints_DerivedMaxConsecutiveDays(t *testing.T) {
	// ACGME has "1 day off per week" → should derive "max 6 consecutive days"
	constraints := comply.GenerateConstraints(comply.US, comply.ForStaff(comply.StaffResident))
	found := false
	for _, c := range constraints {
		if c.RuleKey == "derived-max-consecutive-days" {
			found = true
			if c.Limit != 6 {
				t.Errorf("expected derived max 6 consecutive days, got %.0f", c.Limit)
			}
			if c.Type != comply.ConstraintMaxConsecutive {
				t.Errorf("expected type max_consecutive, got %s", c.Type)
			}
		}
	}
	if !found {
		t.Error("expected derived-max-consecutive-days constraint from ACGME days-off rule")
	}
}

func TestGenerateConstraints_NoPolicyDropped(t *testing.T) {
	// Verify boolean/policy rules are included, not silently dropped.
	constraints := comply.GenerateConstraints(comply.US, comply.ForStaff(comply.StaffResidentPGY1))
	foundPolicy := false
	for _, c := range constraints {
		if c.Type == comply.ConstraintPolicy {
			foundPolicy = true
			break
		}
	}
	if !foundPolicy {
		t.Error("expected boolean/policy constraints to be included (e.g., moonlighting-prohibited-pgy1)")
	}
}

func pad(d int) string {
	return fmt.Sprintf("%02d", d)
}
