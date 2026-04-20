// Package us_il registers Illinois healthcare scheduling regulations:
// ACGME compliance codified into law (210 ILCS 85/6.14), nurse mandatory
// overtime prohibition (210 ILCS 85/10.9), nurse staffing committees
// (210 ILCS 85/10.10), One Day Rest in Seven Act (820 ILCS 140),
// and overtime rules (820 ILCS 105/4a).
package us_il

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

// New returns the Illinois jurisdiction definition.
func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USIL,
		Name:     "Illinois",
		Type:     comply.State,
		Parent:   comply.US,
		Currency: "USD",
		TimeZone: "America/Chicago",
		Rules:    rules(),
	}
}

func rules() []*comply.RuleDef {
	r := make([]*comply.RuleDef, 0, 9)
	r = append(r, residentRules()...)
	r = append(r, nurseRules()...)
	r = append(r, staffingRules()...)
	r = append(r, laborRules()...)
	return r
}

// Resident Rules - 210 ILCS 85/6.14

func residentRules() []*comply.RuleDef {
	return []*comply.RuleDef{
		{
			Key:         "il-acgme-codified",
			Name:        "ACGME Compliance Codified",
			Description: "Illinois codifies ACGME duty hour standards into state law, making them legally enforceable.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpBool,
			StaffTypes:  []comply.Key{comply.StaffResident},
			Scope:       comply.ScopeAccreditedPrograms,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1992, time.January, 1), Amount: 1, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "Illinois Hospital Licensing Act",
				Section: "210 ILCS 85/6.14",
				URL:     "https://ilga.gov/documents/legislation/ilcs/documents/021000850K6.14.htm",
			},
			Notes: "Makes ACGME duty hour standards legally enforceable in Illinois. One of only two states (with NY) to codify resident hours into law.",
		},
	}
}

// Nurse Rules - 210 ILCS 85/10.9

func nurseRules() []*comply.RuleDef {
	nurseSource := comply.Source{
		Title:   "Illinois Hospital Licensing Act",
		Section: "210 ILCS 85/10.9",
		URL:     "https://www.ilga.gov/legislation/ilcs/fulltext.asp?DocName=021000850K10.9",
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
					Since:  comply.D(2005, time.January, 1),
					Amount: 1,
					Unit:   comply.Boolean,
					Exceptions: []string{
						"unforeseen emergent circumstance (declared disasters, hospital disaster plan, specialized nursing skills for procedure completion)",
						"does NOT include routine staffing shortages",
					},
				},
			},
			Source: nurseSource,
			Notes:  "Mandated OT cannot exceed 4 hours beyond predetermined shift. After 12 consecutive hours, nurse must receive 8 hours off.",
		},
		{
			Key:         "il-max-mandatory-ot-hours",
			Name:        "Maximum Mandatory Overtime Hours Beyond Shift",
			Description: "During permitted emergency overtime, mandatory OT cannot exceed 4 hours beyond the predetermined shift.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffNurseRN, comply.StaffNurseLPN},
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2005, time.January, 1), Amount: 4, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: nurseSource,
		},
		{
			Key:         "il-min-rest-after-12h",
			Name:        "Minimum Rest After 12 Consecutive Hours (Nurses)",
			Description: "After 12 consecutive hours of work, nurses must receive at least 8 hours off.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			StaffTypes:  []comply.Key{comply.StaffNurseRN, comply.StaffNurseLPN},
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2005, time.January, 1), Amount: 8, Unit: comply.Hours, Per: comply.PerOccurrence},
			},
			Source: nurseSource,
		},
	}
}

// Staffing Rules - 210 ILCS 85/10.10

func staffingRules() []*comply.RuleDef {
	return []*comply.RuleDef{
		{
			Key:         "il-nurse-staffing-committee",
			Name:        "Nurse Staffing Committee Requirement",
			Description: "Hospitals must establish nurse staffing committees with 50% RNs providing direct care. Must develop hospital-wide staffing plan.",
			Category:    comply.CatStaffing,
			Operator:    comply.OpBool,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2005, time.January, 1), Amount: 1, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "Illinois Hospital Licensing Act",
				Section: "210 ILCS 85/10.10",
				URL:     "https://www.ilga.gov/legislation/ilcs/fulltext.asp?DocName=021000850K10.10",
			},
			Notes: "50% RNs providing direct care. Must develop hospital-wide staffing plan. No mandated ratios.",
		},
	}
}

// Labor Rules - One Day Rest in Seven Act, Overtime

func laborRules() []*comply.RuleDef {
	odrisaSource := comply.Source{
		Title:   "One Day Rest in Seven Act",
		Section: "820 ILCS 140",
		URL:     "https://labor.illinois.gov/laws-rules/fls/odrisa.html",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMinDayOfRest,
			Name:        "Day of Rest",
			Description: "Employees must receive at least 24 consecutive hours of rest per 7-day period.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2023, time.January, 1), Amount: 1, Unit: comply.Days, Per: comply.PerWeek},
			},
			Source: odrisaSource,
			Notes:  "2023 amendment: employees cannot work 7+ consecutive days. Penalties $250-$500.",
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Meal Break Trigger",
			Description: "Employees working 7.5 or more hours must receive a meal break.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2023, time.January, 1), Amount: 7.5, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: odrisaSource,
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Meal Break Duration",
			Description: "Meal breaks must be at least 20 minutes.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2023, time.January, 1), Amount: 20, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: odrisaSource,
			Notes:  "Must begin no later than 5 hours after shift start. Additional 20min for every 4.5h beyond 7.5h.",
		},
		{
			Key:         comply.RuleOvertimeWeeklyThreshold,
			Name:        "Weekly Overtime Threshold",
			Description: "Work exceeding 40 hours per week must be compensated at 1.5x regular rate.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2004, time.January, 1), Amount: 40, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{
				Title:   "Illinois Minimum Wage Law",
				Section: "820 ILCS 105/4a",
				URL:     "https://www.ilga.gov/legislation/ilcs/fulltext.asp?DocName=082001050K4a",
			},
		},
	}
}
