// Package us_vt registers Vermont healthcare scheduling regulations:
// 21 V.S.A. S 304 (meal breaks).
package us_vt

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USVT,
		Name:     "Vermont",
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
			Description: "Reasonable opportunity to eat and use toilet facilities during work shifts. Employees must be given a meal break if working more than 6 hours.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1997, time.July, 1), Amount: 30, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{
				Title:   "Vermont Statutes Annotated",
				Section: "21 V.S.A. S 304",
				URL:     "https://legislature.vermont.gov/statutes/section/21/005/00304",
			},
		},
	}
}
