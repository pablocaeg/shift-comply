// Package us_mt registers Montana healthcare scheduling regulations:
// MCA 39-3-405 (daily overtime after 8 hours for some employers),
// MCA 39-3-601 (meal break for miners, limited applicability).
package us_mt

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USMT,
		Name:     "Montana",
		Type:     comply.State,
		Parent:   comply.US,
		Currency: "USD",
		TimeZone: "America/Denver",
		Rules:    rules(),
	}
}

func rules() []*comply.RuleDef {
	return []*comply.RuleDef{
		{
			Key:         comply.RuleOvertimeDailyThreshold,
			Name:        "Daily Overtime Threshold",
			Description: "Overtime pay (1.5x) required for hours worked in excess of 8 in a workday for employers not covered by FLSA. FLSA-covered employers follow the federal 40-hour weekly threshold.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:      comply.D(1973, time.January, 1),
					Amount:     8,
					Unit:       comply.Hours,
					Per:        comply.PerDay,
					Exceptions: []string{"Applies only to employers not covered by federal FLSA"},
				},
			},
			Source: comply.Source{
				Title:   "Montana Code Annotated",
				Section: "MCA 39-3-405",
				URL:     "https://leg.mt.gov/bills/mca/title_0390/chapter_0030/part_0040/section_0050/0390-0030-0040-0050.html",
			},
		},
		{
			Key:         "no-mandatory-overtime-ban",
			Name:        "Mandatory Overtime Prohibition -- NOT ENACTED",
			Description: "Montana does not prohibit mandatory overtime for nurses or other healthcare workers.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpBool,
			Enforcement: comply.Advisory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1900, time.January, 1), Amount: 0, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "Montana Department of Labor and Industry",
				Section: "No healthcare-specific scheduling statute",
				URL:     "https://erd.dli.mt.gov/",
			},
		},
		{
			Key:         "no-meal-break-requirement",
			Name:        "Meal/Rest Break Requirement -- NOT ENACTED",
			Description: "Montana does not require meal or rest breaks for adult workers in general. Mining-specific breaks exist under MCA 39-3-601.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpBool,
			Enforcement: comply.Advisory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1900, time.January, 1), Amount: 0, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "Montana Department of Labor and Industry",
				Section: "No general state break requirement",
				URL:     "https://erd.dli.mt.gov/",
			},
		},
	}
}
