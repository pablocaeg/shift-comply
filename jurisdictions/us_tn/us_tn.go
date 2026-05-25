// Package us_tn registers Tennessee healthcare scheduling regulations.
// Tennessee has no state-specific healthcare scheduling laws beyond
// federal FLSA and ACGME. Documented regulatory absences are included.
package us_tn

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USTN,
		Name:     "Tennessee",
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
			Name:        "Mandatory Overtime Prohibition -- NOT ENACTED",
			Description: "Tennessee does not prohibit mandatory overtime for nurses or other healthcare workers.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpBool,
			Enforcement: comply.Advisory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1900, time.January, 1), Amount: 0, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "Tennessee Department of Labor",
				Section: "No healthcare-specific scheduling statute",
				URL:     "https://www.tn.gov/workforce/",
			},
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Meal/Rest Break Requirement -- NOT ENACTED",
			Description: "Tennessee does not require meal or rest breaks for adult workers.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpBool,
			Enforcement: comply.Advisory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1900, time.January, 1), Amount: 0, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "Tennessee Department of Labor",
				Section: "No state break requirement",
				URL:     "https://www.tn.gov/workforce/",
			},
		},
	}
}
