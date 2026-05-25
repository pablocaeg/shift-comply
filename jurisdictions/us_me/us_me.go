// Package us_me registers Maine healthcare scheduling regulations:
// 26 MRSA S 603 (mandatory overtime restriction for nurses).
package us_me

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USME,
		Name:     "Maine",
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
			Description: "Nursing facilities may not require nurses or CNA staff to work more than 16 hours in a 24-hour period.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpBool,
			StaffTypes:  []comply.Key{comply.StaffNurseRN, comply.StaffNurseLPN, comply.StaffNurseCNA},
			Scope:       comply.ScopeHealthcareEmployers,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2001, time.September, 21),
					Amount: 1,
					Unit:   comply.Boolean,
					Exceptions: []string{
						"Declared federal or state emergency",
						"Unforeseeable emergency",
					},
				},
			},
			Source: comply.Source{
				Title:   "Maine Revised Statutes",
				Section: "26 MRSA S 603",
				URL:     "https://legislature.maine.gov/statutes/26/title26sec603.html",
			},
		},
		{
			Key:         "me-nurse-max-shift",
			Name:        "Maximum Shift Hours (Nurses)",
			Description: "Nurses and CNAs may not work more than 16 hours in any 24-hour period.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffNurseRN, comply.StaffNurseLPN, comply.StaffNurseCNA},
			Scope:       comply.ScopeHealthcareEmployers,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2001, time.September, 21), Amount: 16, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{
				Title:   "Maine Revised Statutes",
				Section: "26 MRSA S 603",
				URL:     "https://legislature.maine.gov/statutes/26/title26sec603.html",
			},
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Meal Break Duration",
			Description: "30-minute unpaid rest break after 6 consecutive hours of work.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2009, time.September, 12), Amount: 30, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{
				Title:   "Maine Revised Statutes",
				Section: "26 MRSA S 601",
				URL:     "https://legislature.maine.gov/statutes/26/title26sec601.html",
			},
		},
	}
}
