// Package us_ak registers Alaska healthcare scheduling regulations:
// AS 23.10.060 (daily overtime after 8 hours), AS 23.10.070 (weekly overtime).
// Alaska is one of few states with daily overtime outside of California.
package us_ak

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USAK,
		Name:     "Alaska",
		Type:     comply.State,
		Parent:   comply.US,
		Currency: "USD",
		TimeZone: "America/Anchorage",
		Rules:    rules(),
	}
}

func rules() []*comply.RuleDef {
	return []*comply.RuleDef{
		{
			Key:         comply.RuleOvertimeDailyThreshold,
			Name:        "Daily Overtime Threshold",
			Description: "Overtime pay (1.5x) required for hours worked in excess of 8 in a workday.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1959, time.January, 3), Amount: 8, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{
				Title:   "Alaska Statutes",
				Section: "AS 23.10.060",
				URL:     "https://www.akleg.gov/basis/statutes.asp#23.10.060",
			},
		},
		{
			Key:         comply.RuleMandatoryOTProhibited,
			Name:        "Mandatory Overtime Prohibition -- NOT ENACTED",
			Description: "Alaska does not prohibit mandatory overtime for nurses or other healthcare workers.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpBool,
			Enforcement: comply.Advisory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1900, time.January, 1), Amount: 0, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "Alaska Department of Labor",
				Section: "No healthcare-specific scheduling statute",
				URL:     "https://labor.alaska.gov/",
			},
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Meal/Rest Break Requirement -- NOT ENACTED",
			Description: "Alaska does not require meal or rest breaks for adult workers.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpBool,
			Enforcement: comply.Advisory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1900, time.January, 1), Amount: 0, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "Alaska Department of Labor",
				Section: "No state break requirement",
				URL:     "https://labor.alaska.gov/",
			},
		},
	}
}
