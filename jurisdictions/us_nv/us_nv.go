// Package us_nv registers Nevada healthcare scheduling regulations:
// NRS 608.019 (meal/rest breaks), NRS 608.018 (daily overtime).
package us_nv

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USNV,
		Name:     "Nevada",
		Type:     comply.State,
		Parent:   comply.US,
		Currency: "USD",
		TimeZone: "America/Los_Angeles",
		Rules:    rules(),
	}
}

func rules() []*comply.RuleDef {
	return []*comply.RuleDef{
		{
			Key:         comply.RuleOvertimeDailyThreshold,
			Name:        "Daily Overtime Threshold",
			Description: "Overtime pay required for hours worked in excess of 8 in a 24-hour period.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1973, time.January, 1), Amount: 8, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{
				Title:   "Nevada Revised Statutes",
				Section: "NRS 608.018",
				URL:     "https://www.leg.state.nv.us/nrs/nrs-608.html#NRS608Sec018",
			},
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Meal Break Duration",
			Description: "30-minute meal break for 8-hour continuous shifts.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1973, time.January, 1), Amount: 30, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{
				Title:   "Nevada Revised Statutes",
				Section: "NRS 608.019",
				URL:     "https://www.leg.state.nv.us/nrs/nrs-608.html#NRS608Sec019",
			},
		},
		{
			Key:         comply.RuleRestBreakDuration,
			Name:        "Rest Break",
			Description: "Paid 10-minute rest break for each 3.5 hours of continuous work.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1973, time.January, 1), Amount: 10, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{
				Title:   "Nevada Revised Statutes",
				Section: "NRS 608.019",
				URL:     "https://www.leg.state.nv.us/nrs/nrs-608.html#NRS608Sec019",
			},
		},
	}
}
