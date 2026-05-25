// Package us_wi registers Wisconsin healthcare scheduling regulations:
// Wis. Admin. Code DWD 274.02 (meal breaks), Wis. Stat. S 103.85 (day of rest).
package us_wi

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USWI,
		Name:     "Wisconsin",
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
			Key:         comply.RuleMealBreakDuration,
			Name:        "Meal Break",
			Description: "30-minute meal break for shifts of 6+ hours. Must be provided reasonably close to the usual meal time.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Recommended,
			Values: []*comply.RuleValue{
				{Since: comply.D(1980, time.January, 1), Amount: 30, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{
				Title:   "Wisconsin Administrative Code",
				Section: "DWD 274.02",
				URL:     "https://docs.legis.wisconsin.gov/code/admin_code/dwd/272_to_280/274/02",
			},
		},
		{
			Key:         comply.RuleMinDayOfRest,
			Name:        "Day of Rest (One Day in Seven)",
			Description: "Each employer must provide at least 24 consecutive hours of rest in each 7-day period.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1931, time.January, 1), Amount: 24, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{
				Title:   "Wisconsin Statutes",
				Section: "S 103.85",
				URL:     "https://docs.legis.wisconsin.gov/statutes/statutes/103/85",
			},
		},
	}
}
