// Package us_sd registers South Dakota healthcare scheduling regulations.
// South Dakota has no state-specific healthcare scheduling laws beyond
// federal FLSA and ACGME. Documented regulatory absences are included.
package us_sd

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USSD,
		Name:     "South Dakota",
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
			Key:         "no-mandatory-overtime-ban",
			Name:        "Mandatory Overtime Prohibition -- NOT ENACTED",
			Description: "South Dakota does not prohibit mandatory overtime for nurses or other healthcare workers.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpBool,
			Enforcement: comply.Advisory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1900, time.January, 1), Amount: 0, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "South Dakota Department of Labor",
				Section: "No healthcare-specific scheduling statute",
				URL:     "https://dlr.sd.gov/",
			},
		},
		{
			Key:         "no-meal-break-requirement",
			Name:        "Meal/Rest Break Requirement -- NOT ENACTED",
			Description: "South Dakota does not require meal or rest breaks for adult workers.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpBool,
			Enforcement: comply.Advisory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1900, time.January, 1), Amount: 0, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "South Dakota Department of Labor",
				Section: "No state break requirement",
				URL:     "https://dlr.sd.gov/",
			},
		},
	}
}
