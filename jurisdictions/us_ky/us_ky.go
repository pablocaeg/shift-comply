// Package us_ky registers Kentucky healthcare scheduling regulations:
// KRS 337.355 (meal breaks), KRS 337.365 (rest breaks).
package us_ky

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USKY,
		Name:     "Kentucky",
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
			Key:         comply.RuleMealBreakDuration,
			Name:        "Meal Break",
			Description: "Reasonable unpaid meal break between the 3rd and 5th hour of work.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1974, time.January, 1), Amount: 30, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{
				Title:   "Kentucky Revised Statutes",
				Section: "KRS 337.355",
				URL:     "https://apps.legislature.ky.gov/law/statutes/statute.aspx?id=32092",
			},
		},
		{
			Key:         comply.RuleRestBreakDuration,
			Name:        "Rest Break",
			Description: "Paid 10-minute rest break per 4 hours of work.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1974, time.January, 1), Amount: 10, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{
				Title:   "Kentucky Revised Statutes",
				Section: "KRS 337.365",
				URL:     "https://apps.legislature.ky.gov/law/statutes/statute.aspx?id=32094",
			},
		},
	}
}
