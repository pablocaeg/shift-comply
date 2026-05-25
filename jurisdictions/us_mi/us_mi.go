// Package us_mi registers Michigan healthcare scheduling regulations:
// MCL 409.112 (meal breaks for certain industries).
// Michigan also has a Youth Employment Standards Act with break
// requirements for minors, but no general adult break mandate.
package us_mi

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USMI,
		Name:     "Michigan",
		Type:     comply.State,
		Parent:   comply.US,
		Currency: "USD",
		TimeZone: "America/Detroit",
		Rules:    rules(),
	}
}

func rules() []*comply.RuleDef {
	return []*comply.RuleDef{
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Meal Break",
			Description: "30-minute uncompensated meal break for shifts of 5+ consecutive hours. Applies to most employers.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Recommended,
			Values: []*comply.RuleValue{
				{Since: comply.D(2014, time.January, 1), Amount: 30, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{
				Title:   "Michigan Compiled Laws",
				Section: "MCL 409.112",
				URL:     "https://www.legislature.mi.gov/Laws/MCL?objectName=mcl-409-112",
			},
			Notes: "Michigan does not have a general mandatory break law for adults. The statute provides for breaks in specific industries. Federal FLSA applies for compensable time determination.",
		},
		{
			Key:         "no-mandatory-overtime-ban",
			Name:        "Mandatory Overtime Prohibition -- NOT ENACTED",
			Description: "Michigan does not prohibit mandatory overtime for nurses or other healthcare workers.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpBool,
			Enforcement: comply.Advisory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1900, time.January, 1), Amount: 0, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "Michigan Department of Labor and Economic Opportunity",
				Section: "No healthcare-specific scheduling statute",
				URL:     "https://www.michigan.gov/leo/",
			},
		},
	}
}
