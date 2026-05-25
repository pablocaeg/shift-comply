// Package us_co registers Colorado healthcare scheduling regulations:
// COMPS Order #38 (meal/rest breaks), Colorado Overtime and Minimum Pay
// Standards Order (daily overtime).
package us_co

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USCO,
		Name:     "Colorado",
		Type:     comply.State,
		Parent:   comply.US,
		Currency: "USD",
		TimeZone: "America/Denver",
		Rules:    rules(),
	}
}

func rules() []*comply.RuleDef {
	comps := comply.Source{
		Title:   "Colorado Overtime and Minimum Pay Standards Order (COMPS Order #38)",
		Section: "",
		URL:     "https://cdle.colorado.gov/comps-order",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleOvertimeDailyThreshold,
			Name:        "Daily Overtime Threshold",
			Description: "Overtime pay required for hours worked in excess of 12 in a workday, or 12 consecutive hours regardless of start/end times.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2020, time.March, 16), Amount: 12, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: comps.Title, Section: "Rule 4.1.1", URL: comps.URL},
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Meal Break Requirement",
			Description: "30-minute uninterrupted meal break for shifts exceeding 5 hours.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2020, time.March, 16), Amount: 5, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{Title: comps.Title, Section: "Rule 5.2", URL: comps.URL},
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Meal Break Duration",
			Description: "At least 30 minutes. Unpaid if employee is completely relieved of duties.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2020, time.March, 16), Amount: 30, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{Title: comps.Title, Section: "Rule 5.2", URL: comps.URL},
		},
		{
			Key:         comply.RuleRestBreakDuration,
			Name:        "Rest Break",
			Description: "Paid 10-minute rest break for each 4-hour work period.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2020, time.March, 16), Amount: 10, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{Title: comps.Title, Section: "Rule 5.1", URL: comps.URL},
		},
	}
}
