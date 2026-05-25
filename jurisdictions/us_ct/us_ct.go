// Package us_ct registers Connecticut healthcare scheduling regulations:
// Public Act 16-2 (nurse mandatory overtime ban), CGS S 31-76b (meal breaks).
package us_ct

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USCT,
		Name:     "Connecticut",
		Type:     comply.State,
		Parent:   comply.US,
		Currency: "USD",
		TimeZone: "America/New_York",
		Rules:    rules(),
	}
}

func rules() []*comply.RuleDef {
	return []*comply.RuleDef{
		{
			Key:         comply.RuleMandatoryOTProhibited,
			Name:        "Mandatory Overtime Prohibition (Nurses)",
			Description: "Hospitals may not require nurses to work overtime. Does not apply to voluntary overtime.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpBool,
			StaffTypes:  []comply.Key{comply.StaffNurseRN, comply.StaffNurseLPN},
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2016, time.October, 1),
					Amount: 1,
					Unit:   comply.Boolean,
					Exceptions: []string{
						"Declaration of emergency by federal, state, or municipal government",
						"Unforeseen emergent circumstance",
					},
				},
			},
			Source: comply.Source{
				Title:   "Connecticut General Statutes",
				Section: "CGS S 19a-490o (Public Act 16-2)",
				URL:     "https://www.cga.ct.gov/current/pub/chap_368v.htm#sec_19a-490o",
			},
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Meal Break Requirement",
			Description: "Employees working 7.5 or more consecutive hours must receive a 30-minute meal break within the first 7.5 hours.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1967, time.October, 1), Amount: 7.5, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{
				Title:   "Connecticut General Statutes",
				Section: "CGS S 31-76b",
				URL:     "https://www.cga.ct.gov/current/pub/chap_558.htm#sec_31-76b",
			},
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Meal Break Duration",
			Description: "At least 30 consecutive minutes.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1967, time.October, 1), Amount: 30, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{
				Title:   "Connecticut General Statutes",
				Section: "CGS S 31-76b",
				URL:     "https://www.cga.ct.gov/current/pub/chap_558.htm#sec_31-76b",
			},
		},
		{
			Key:         comply.RuleMinDayOfRest,
			Name:        "Day of Rest",
			Description: "Employers must provide at least 24 consecutive hours of rest in each 7-day period.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1967, time.October, 1), Amount: 24, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{
				Title:   "Connecticut General Statutes",
				Section: "CGS S 53-303e",
				URL:     "https://www.cga.ct.gov/current/pub/chap_945.htm#sec_53-303e",
			},
		},
	}
}
