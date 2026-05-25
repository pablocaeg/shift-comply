// Package us_in registers Indiana healthcare scheduling regulations:
// IC 22-2-7 (one day of rest in seven).
package us_in

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USIN,
		Name:     "Indiana",
		Type:     comply.State,
		Parent:   comply.US,
		Currency: "USD",
		TimeZone: "America/Indiana/Indianapolis",
		Rules:    rules(),
	}
}

func rules() []*comply.RuleDef {
	return []*comply.RuleDef{
		{
			Key:         comply.RuleMinDayOfRest,
			Name:        "Day of Rest (One Day in Seven)",
			Description: "Every employer must allow every employee at least 24 consecutive hours of rest in each 7-day period.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1971, time.January, 1), Amount: 24, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{
				Title:   "Indiana Code",
				Section: "IC 22-2-7-1",
				URL:     "http://iga.in.gov/laws/2024/ic/titles/22#22-2-7",
			},
		},
		{
			Key:         comply.RuleMandatoryOTProhibited,
			Name:        "Mandatory Overtime Prohibition -- NOT ENACTED",
			Description: "Indiana does not prohibit mandatory overtime for nurses or other healthcare workers.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpBool,
			Enforcement: comply.Advisory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1900, time.January, 1), Amount: 0, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "Indiana Department of Labor",
				Section: "No healthcare-specific scheduling statute",
				URL:     "https://www.in.gov/dol/",
			},
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Meal/Rest Break Requirement -- NOT ENACTED",
			Description: "Indiana does not require meal or rest breaks for adult workers.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpBool,
			Enforcement: comply.Advisory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1900, time.January, 1), Amount: 0, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "Indiana Department of Labor",
				Section: "No state break requirement",
				URL:     "https://www.in.gov/dol/",
			},
		},
	}
}
