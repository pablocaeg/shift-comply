// Package us_nc registers North Carolina healthcare scheduling regulations.
// North Carolina has no state-specific healthcare scheduling laws, nurse
// mandatory overtime bans, or break requirements beyond federal FLSA.
package us_nc

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USNC,
		Name:     "North Carolina",
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
			Description: "North Carolina does not prohibit mandatory overtime for nurses or any healthcare workers.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpBool,
			Enforcement: comply.Advisory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1900, time.January, 1), Amount: 0, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "North Carolina Department of Labor",
				Section: "Wage and Hour Act (NCGS 95-25)",
				URL:     "https://www.labor.nc.gov/workplace-rights/employee-rights-regarding-time-worked-and-wages-earned",
			},
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Break Requirement -- NOT ENACTED",
			Description: "North Carolina does not require meal or rest breaks for workers aged 16 and over.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpBool,
			Enforcement: comply.Advisory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1900, time.January, 1), Amount: 0, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "North Carolina Department of Labor",
				Section: "Wage and Hour Act (NCGS 95-25.14)",
				URL:     "https://www.labor.nc.gov/workplace-rights/employee-rights-regarding-time-worked-and-wages-earned/breaks-meals",
			},
		},
	}
}
