// Package us_ne registers Nebraska healthcare scheduling regulations:
// Neb. Rev. Stat. S 48-212 (meal breaks).
package us_ne

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USNE,
		Name:     "Nebraska",
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
			Description: "30-minute meal break for 8-hour shifts in assembly plants and workshops. Other industries follow federal FLSA only.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1989, time.January, 1), Amount: 30, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{
				Title:   "Nebraska Revised Statutes",
				Section: "S 48-212",
				URL:     "https://nebraskalegislature.gov/laws/statutes.php?statute=48-212",
			},
			Notes: "Limited applicability. Many Nebraska employers in healthcare are not covered by the state break law.",
		},
	}
}
