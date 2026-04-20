// Package us_or registers Oregon healthcare scheduling regulations:
// nurse mandatory overtime prohibition and shift limits (ORS 441.166,
// strengthened by HB 2697 2023), meal and rest break requirements
// (ORS 653.261), and overtime rules.
package us_or

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

// New returns the Oregon jurisdiction definition.
func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USOR,
		Name:     "Oregon",
		Type:     comply.State,
		Parent:   comply.US,
		Currency: "USD",
		TimeZone: "America/Los_Angeles",
		Rules:    rules(),
	}
}

func rules() []*comply.RuleDef {
	r := make([]*comply.RuleDef, 0, 8)
	r = append(r, nurseRules()...)
	r = append(r, breakRules()...)
	r = append(r, overtimeRules()...)
	return r
}

// Nurse Rules - ORS 441.166 (HB 2697, 2023)

func nurseRules() []*comply.RuleDef {
	nurseSource := comply.Source{
		Title:   "Oregon Revised Statutes",
		Section: "ORS 441.166",
		URL:     "https://www.oregonlaws.org/ors/441.166",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMandatoryOTProhibited,
			Name:        "Nurse Mandatory Overtime Prohibition",
			Description: "Hospitals may not require nurses or CNAs to work mandatory overtime.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpBool,
			StaffTypes:  []comply.Key{comply.StaffNurseRN, comply.StaffNurseLPN, comply.StaffNurseCNA},
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2023, time.September, 1),
					Amount: 1,
					Unit:   comply.Boolean,
					Exceptions: []string{
						"one additional hour for shift vacancy or patient safety",
					},
				},
			},
			Source: nurseSource,
			Notes:  "Significantly strengthened by HB 2697 (2023). Covers RNs, LPNs, CNAs.",
		},
		{
			Key:         comply.RuleMaxShiftHours,
			Name:        "Maximum Shift Hours (Nurses)",
			Description: "Nurses and CNAs may not work more than 12 consecutive hours in a 24-hour period.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffNurseRN, comply.StaffNurseLPN, comply.StaffNurseCNA},
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2023, time.September, 1), Amount: 12, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: nurseSource,
			Notes:  "Cannot exceed 12 consecutive hours in 24-hour period.",
		},
		{
			Key:         "or-max-weekly-hours-nurse",
			Name:        "Maximum Weekly Hours (Nurses)",
			Description: "Nurses and CNAs may not work more than 48 hours per week.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffNurseRN, comply.StaffNurseLPN, comply.StaffNurseCNA},
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2023, time.September, 1), Amount: 48, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: nurseSource,
		},
		{
			Key:         comply.RuleMinRestBetweenShifts,
			Name:        "Minimum Rest Between Shifts (Nurses)",
			Description: "Nurses and CNAs must receive at least 10 hours of rest after working 12 hours or after any shift exceeding 12 hours.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			StaffTypes:  []comply.Key{comply.StaffNurseRN, comply.StaffNurseLPN, comply.StaffNurseCNA},
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2023, time.September, 1), Amount: 10, Unit: comply.Hours, Per: comply.PerOccurrence},
			},
			Source: nurseSource,
			Notes:  "10-hour rest required after 12th hour worked or after any shift exceeding 12 hours.",
		},
	}
}

// Break Rules - ORS 653.261

func breakRules() []*comply.RuleDef {
	breakSource := comply.Source{
		Title:   "Oregon Revised Statutes",
		Section: "ORS 653.261",
		URL:     "https://www.oregon.gov/boli/workers/pages/meals-and-breaks.aspx",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Meal Break Trigger",
			Description: "Employees working more than 6 hours must receive a meal break.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2000, time.January, 1), Amount: 6, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: breakSource,
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Meal Break Duration",
			Description: "Meal breaks must be at least 30 minutes.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2000, time.January, 1), Amount: 30, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: breakSource,
		},
		{
			Key:         comply.RuleRestBreakDuration,
			Name:        "Rest Break Duration",
			Description: "10-minute paid rest break per 4-hour work period.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2000, time.January, 1), Amount: 10, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: breakSource,
		},
	}
}

// Overtime Rules - ORS 653.261

func overtimeRules() []*comply.RuleDef {
	return []*comply.RuleDef{
		{
			Key:         comply.RuleOvertimeWeeklyThreshold,
			Name:        "Weekly Overtime Threshold",
			Description: "Work exceeding 40 hours per week must be compensated at 1.5x regular rate.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2000, time.January, 1), Amount: 40, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{
				Title:   "Oregon Revised Statutes",
				Section: "ORS 653.261",
				URL:     "https://www.oregon.gov/boli/employers/pages/overtime.aspx",
			},
		},
	}
}
