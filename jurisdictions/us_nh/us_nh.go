// Package us_nh registers New Hampshire healthcare scheduling regulations:
// RSA 326-B:18-a (mandatory overtime prohibition for nurses),
// RSA 275:30-a (meal breaks).
package us_nh

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USNH,
		Name:     "New Hampshire",
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
			Description: "Nursing facilities may not require nurses to work mandatory overtime beyond their regularly scheduled hours.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpBool,
			StaffTypes:  []comply.Key{comply.StaffNurseRN, comply.StaffNurseLPN},
			Scope:       comply.ScopeHealthcareEmployers,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2007, time.September, 16),
					Amount: 1,
					Unit:   comply.Boolean,
					Exceptions: []string{
						"Declared federal or state emergency",
						"Unforeseen emergency",
					},
				},
			},
			Source: comply.Source{
				Title:   "New Hampshire Revised Statutes Annotated",
				Section: "RSA 326-B:18-a",
				URL:     "https://www.gencourt.state.nh.us/rsa/html/XXX/326-B/326-B-18-a.htm",
			},
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Meal Break",
			Description: "30-minute meal break after 5 consecutive hours of work.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1985, time.January, 1), Amount: 30, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{
				Title:   "New Hampshire Revised Statutes Annotated",
				Section: "RSA 275:30-a",
				URL:     "https://www.gencourt.state.nh.us/rsa/html/XXIII/275/275-30-a.htm",
			},
		},
	}
}
