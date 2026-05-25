// Package us_nd registers North Dakota healthcare scheduling regulations:
// NDCC S 34-06-03 (meal breaks).
package us_nd

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USND,
		Name:     "North Dakota",
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
			Description: "30-minute meal break for shifts of 5+ hours if two or more employees are on duty.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1979, time.July, 1), Amount: 30, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{
				Title:   "North Dakota Century Code",
				Section: "NDCC S 34-06-03",
				URL:     "https://www.ndlegis.gov/cencode/t34c06.pdf",
			},
		},
	}
}
