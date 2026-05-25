// Package us_wv registers West Virginia healthcare scheduling regulations:
// WV Code S 21-5F (Nurse Overtime Standards Act), meal break requirements.
package us_wv

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USWV,
		Name:     "West Virginia",
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
			Name:        "Mandatory Overtime Restriction (Nurses)",
			Description: "Healthcare facilities may not require a nurse to work overtime. Refusal is not patient abandonment and may not be used as grounds for dismissal or discipline.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpBool,
			StaffTypes:  []comply.Key{comply.StaffNurseRN, comply.StaffNurseLPN},
			Scope:       comply.ScopeHealthcareEmployers,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2007, time.June, 8),
					Amount: 1,
					Unit:   comply.Boolean,
					Exceptions: []string{
						"Declared federal, state, or local emergency",
						"Unforeseen emergency that could not be prudently planned for",
						"Nurse engaged in an ongoing medical or surgical procedure until completed",
					},
				},
			},
			Source: comply.Source{
				Title:   "West Virginia Code",
				Section: "S 21-5F-1 to 21-5F-7 (Nurse Overtime Standards Act)",
				URL:     "https://code.wvlegislature.gov/21-5F/",
			},
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Meal Break",
			Description: "At least 20-minute meal break for employees working 6+ hours. Must be provided within a reasonable time after starting work.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1979, time.July, 1), Amount: 20, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{
				Title:   "West Virginia Code",
				Section: "S 21-3-10a",
				URL:     "https://code.wvlegislature.gov/21-3-10A/",
			},
		},
	}
}
