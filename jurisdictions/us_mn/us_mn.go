// Package us_mn registers Minnesota healthcare scheduling regulations:
// Minn. Stat. S 181.275 (mandatory overtime restriction for nurses),
// Minn. Stat. S 177.253-177.254 (rest breaks).
package us_mn

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USMN,
		Name:     "Minnesota",
		Type:     comply.State,
		Parent:   comply.US,
		Currency: "USD",
		TimeZone: "America/Chicago",
		Rules:    rules(),
	}
}

func rules() []*comply.RuleDef {
	return []*comply.RuleDef{
		{
			Key:         comply.RuleMandatoryOTProhibited,
			Name:        "Mandatory Overtime Restriction (Nurses)",
			Description: "Hospitals and healthcare facilities may not require nurses to work in excess of a predetermined and regularly scheduled shift.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpBool,
			StaffTypes:  []comply.Key{comply.StaffNurseRN, comply.StaffNurseLPN},
			Scope:       comply.ScopeHealthcareEmployers,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2007, time.August, 1),
					Amount: 1,
					Unit:   comply.Boolean,
					Exceptions: []string{
						"Declared national, state, or municipal emergency",
						"Unforeseen emergency that could not be prudently planned for",
						"Nurse is in a surgical or diagnostic procedure and the procedure is not yet completed",
					},
				},
			},
			Source: comply.Source{
				Title:   "Minnesota Statutes",
				Section: "S 181.275",
				URL:     "https://www.revisor.mn.gov/statutes/cite/181.275",
			},
		},
		{
			Key:         comply.RuleRestBreakDuration,
			Name:        "Rest Break",
			Description: "Adequate time to use the nearest restroom within each 4 consecutive hours of work.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1988, time.January, 1), Amount: 10, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{
				Title:   "Minnesota Statutes",
				Section: "S 177.253",
				URL:     "https://www.revisor.mn.gov/statutes/cite/177.253",
			},
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Meal Break Requirement",
			Description: "Employees working 8+ consecutive hours must have sufficient time to eat a meal.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1988, time.January, 1), Amount: 8, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{
				Title:   "Minnesota Statutes",
				Section: "S 177.254",
				URL:     "https://www.revisor.mn.gov/statutes/cite/177.254",
			},
		},
	}
}
