// Package us_sc registers South Carolina healthcare scheduling regulations.
// South Carolina has no state-specific healthcare scheduling laws beyond
// federal FLSA and ACGME. Documented regulatory absences are included.
package us_sc

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USSC,
		Name:     "South Carolina",
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
			Name:        "Mandatory Overtime Prohibition -- NOT ENACTED",
			Description: "South Carolina does not prohibit mandatory overtime for nurses or other healthcare workers.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpBool,
			Enforcement: comply.Advisory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1900, time.January, 1), Amount: 0, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "South Carolina Department of Labor",
				Section: "No healthcare-specific scheduling statute",
				URL:     "https://www.llr.sc.gov/",
			},
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Meal/Rest Break Requirement -- NOT ENACTED",
			Description: "South Carolina does not require meal or rest breaks for adult workers.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpBool,
			Enforcement: comply.Advisory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1900, time.January, 1), Amount: 0, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "South Carolina Department of Labor",
				Section: "No state break requirement",
				URL:     "https://www.llr.sc.gov/",
			},
		},
	}
}
