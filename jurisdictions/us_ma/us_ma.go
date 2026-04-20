// Package us_ma registers Massachusetts healthcare scheduling regulations:
// mandatory overtime prohibition for nurses (MGL Ch. 111 S226), maximum shift
// hours, ICU nurse-patient ratio (MGL Ch. 111 S231), meal break requirements
// (MGL Ch. 149 S100), day of rest (MGL Ch. 149 S48), and overtime rules
// (MGL Ch. 151 S1A).
package us_ma

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

// New returns the Massachusetts jurisdiction definition.
func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USMA,
		Name:     "Massachusetts",
		Type:     comply.State,
		Parent:   comply.US,
		Currency: "USD",
		TimeZone: "America/New_York",
		Rules:    rules(),
	}
}

func rules() []*comply.RuleDef {
	r := make([]*comply.RuleDef, 0, 10)
	r = append(r, nurseRules()...)
	r = append(r, staffingRules()...)
	r = append(r, breakRules()...)
	r = append(r, laborRules()...)
	return r
}

// Nurse Rules - MGL Part I, Title XVI, Chapter 111, Section 226

func nurseRules() []*comply.RuleDef {
	nurseSource := comply.Source{
		Title:   "Massachusetts General Laws",
		Section: "Part I, Title XVI, Chapter 111, Section 226",
		URL:     "https://malegislature.gov/Laws/GeneralLaws/PartI/TitleXVI/Chapter111/Section226",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMandatoryOTProhibited,
			Name:        "Nurse Mandatory Overtime Prohibition",
			Description: "Hospitals may not require nurses to work mandatory overtime.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpBool,
			StaffTypes:  []comply.Key{comply.StaffNurseRN, comply.StaffNurseLPN},
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2012, time.November, 5),
					Amount: 1,
					Unit:   comply.Boolean,
					Exceptions: []string{
						"emergency situations only",
					},
				},
			},
			Source: nurseSource,
			Notes:  "Part of Chapter 224 of the Acts of 2012. Hospitals must report all mandatory OT to DPH.",
		},
		{
			Key:         comply.RuleMaxShiftHours,
			Name:        "Maximum Shift Hours (Nurses)",
			Description: "Nurses may not work more than 16 hours in a shift.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffNurseRN, comply.StaffNurseLPN},
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2012, time.November, 5), Amount: 16, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: nurseSource,
			Notes:  "Hard cap regardless of voluntary/mandatory.",
		},
		{
			Key:         comply.RuleMinRestAfterExtended,
			Name:        "Minimum Rest After Extended Shift (Nurses)",
			Description: "Nurses must receive at least 8 hours of rest after a 16-hour shift.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			StaffTypes:  []comply.Key{comply.StaffNurseRN, comply.StaffNurseLPN},
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2012, time.November, 5), Amount: 8, Unit: comply.Hours, Per: comply.PerOccurrence},
			},
			Source: nurseSource,
		},
	}
}

// Staffing Rules - MGL Part I, Title XVI, Chapter 111, Section 231

func staffingRules() []*comply.RuleDef {
	return []*comply.RuleDef{
		{
			Key:         "ma-icu-nurse-ratio",
			Name:        "ICU Nurse-Patient Ratio",
			Description: "ICU nurse-to-patient ratio must not exceed 1:2. Acuity-based 1:1 or 1:2.",
			Category:    comply.CatStaffing,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffNurseRN},
			UnitTypes:   []comply.Key{comply.UnitICU},
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2014, time.June, 30), Amount: 2, Unit: comply.PatientsPerNurse, Per: comply.PerShift},
			},
			Source: comply.Source{
				Title:   "Massachusetts General Laws",
				Section: "Part I, Title XVI, Chapter 111, Section 231",
				URL:     "https://malegislature.gov/Laws/GeneralLaws/PartI/TitleXVI/Chapter111/Section231",
			},
			Notes: "Acuity-based 1:1 or 1:2. Only state besides CA with mandated ICU ratio.",
		},
		{
			Key:         "ma-nurse-ratios-rejected",
			Name:        "Statewide Nurse Ratios Rejected",
			Description: "Ballot Question 1 (2018) to establish statewide nurse-patient ratios was rejected by voters.",
			Category:    comply.CatStaffing,
			Operator:    comply.OpBool,
			Enforcement: comply.Advisory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2018, time.November, 6), Amount: 0, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "Massachusetts General Laws",
				Section: "Part I, Title XVI, Chapter 111, Section 231",
				URL:     "https://malegislature.gov/Laws/GeneralLaws/PartI/TitleXVI/Chapter111/Section231",
			},
			Notes: "Ballot Question 1 (2018) rejected 70.4%-29.6%.",
		},
	}
}

// Break Rules - MGL Chapter 149, Section 100

func breakRules() []*comply.RuleDef {
	mealSource := comply.Source{
		Title:   "Massachusetts General Laws",
		Section: "Chapter 149, Section 100",
		URL:     "https://malegislature.gov/Laws/GeneralLaws/PartI/TitleXXI/Chapter149/Section100",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Meal Break Trigger",
			Description: "Employees working more than 6 hours must receive a 30-minute meal break.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1994, time.January, 1), Amount: 6, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: mealSource,
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Meal Break Duration",
			Description: "Meal breaks must be at least 30 minutes.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1994, time.January, 1), Amount: 30, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: mealSource,
		},
	}
}

// Labor Rules - Day of Rest, Overtime

func laborRules() []*comply.RuleDef {
	return []*comply.RuleDef{
		{
			Key:         comply.RuleMinDayOfRest,
			Name:        "Day of Rest",
			Description: "Every employee is entitled to one day of rest in seven.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1940, time.January, 1), Amount: 1, Unit: comply.Days, Per: comply.PerWeek},
			},
			Source: comply.Source{
				Title:   "Massachusetts General Laws",
				Section: "Chapter 149, Section 48",
				URL:     "https://malegislature.gov/Laws/GeneralLaws/PartI/TitleXXI/Chapter149/Section48",
			},
		},
		{
			Key:         comply.RuleOvertimeWeeklyThreshold,
			Name:        "Weekly Overtime Threshold",
			Description: "Work exceeding 40 hours per week must be compensated at 1.5x regular rate.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1998, time.January, 1), Amount: 40, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{
				Title:   "Massachusetts General Laws",
				Section: "Chapter 151, Section 1A",
				URL:     "https://malegislature.gov/Laws/GeneralLaws/PartI/TitleXXI/Chapter151/Section1A",
			},
		},
		{
			Key:         comply.RuleOvertimeWeeklyRate,
			Name:        "Weekly Overtime Rate",
			Description: "1.5x regular rate for hours exceeding 40 per week.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1998, time.January, 1), Amount: 1.5, Unit: comply.Multiplier, Per: comply.PerWeek},
			},
			Source: comply.Source{
				Title:   "Massachusetts General Laws",
				Section: "Chapter 151, Section 1A",
				URL:     "https://malegislature.gov/Laws/GeneralLaws/PartI/TitleXXI/Chapter151/Section1A",
			},
		},
	}
}
