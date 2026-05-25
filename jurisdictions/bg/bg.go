// Package bg registers Bulgaria's healthcare scheduling regulations:
// Kodeks na truda (Labour Code, State Gazette 26/1986, last consolidated 2024).
package bg

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:      comply.BG,
		Name:      "Bulgaria",
		LocalName: "Balgariya",
		Type:      comply.Country,
		Parent:    comply.EU,
		Currency:  "BGN",
		TimeZone:  "Europe/Sofia",
		Rules:     rules(),
	}
}

func rules() []*comply.RuleDef {
	kt := comply.Source{
		Title:   "Kodeks na truda (Labour Code, SG 26/1986)",
		Section: "",
		URL:     "https://lex.bg/laws/ldoc/1594373121",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxOrdinaryWeeklyHours,
			Name:        "Maximum Ordinary Weekly Hours",
			Description: "Normal working time is 8 hours per day and 40 hours per 5-day week.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1986, time.April, 1), Amount: 40, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: kt.Title, Section: "Art. 136", URL: kt.URL},
		},
		{
			Key:         comply.RuleMaxWeeklyHours,
			Name:        "Maximum Weekly Hours Including Overtime",
			Description: "Average weekly working time including overtime shall not exceed 48 hours over a 4-month reference period.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2004, time.January, 1),
					Amount:   48,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 4, Unit: comply.PeriodMonths},
				},
			},
			Source: comply.Source{Title: kt.Title, Section: "Art. 136a", URL: kt.URL},
		},
		{
			Key:         comply.RuleMinRestBetweenShifts,
			Name:        "Minimum Daily Rest Period",
			Description: "At least 12 consecutive hours of rest between two working days.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1986, time.April, 1), Amount: 12, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: kt.Title, Section: "Art. 152", URL: kt.URL},
			Notes:  "Bulgaria requires 12 hours rest, above the EU 11-hour minimum.",
		},
		{
			Key:         comply.RuleMinWeeklyRest,
			Name:        "Minimum Weekly Rest Period",
			Description: "At least 48 consecutive hours of weekly rest, normally Saturday and Sunday.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1986, time.April, 1), Amount: 48, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: kt.Title, Section: "Art. 153", URL: kt.URL},
			Notes:  "48 hours weekly rest is among the highest in Europe, alongside Hungary and Romania.",
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Break Duration",
			Description: "At least 30 minutes meal break during a full working day. Not less than 1 hour total break time for an 8-hour day.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1986, time.April, 1), Amount: 30, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{Title: kt.Title, Section: "Art. 151", URL: kt.URL},
		},
		{
			Key:         comply.RuleMaxOvertimeAnnual,
			Name:        "Maximum Annual Overtime",
			Description: "Maximum 150 hours of overtime per year.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1986, time.April, 1), Amount: 150, Unit: comply.Hours, Per: comply.PerYear},
			},
			Source: comply.Source{Title: kt.Title, Section: "Art. 146(1)", URL: kt.URL},
		},
		{
			Key:         comply.RuleNightPeriodStart,
			Name:        "Night Period Start",
			Description: "Night work is defined as work between 22:00 and 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1986, time.April, 1), Amount: 22, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: kt.Title, Section: "Art. 140(2)", URL: kt.URL},
		},
		{
			Key:         comply.RuleNightPeriodEnd,
			Name:        "Night Period End",
			Description: "Night work period ends at 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1986, time.April, 1), Amount: 6, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: kt.Title, Section: "Art. 140(2)", URL: kt.URL},
		},
		{
			Key:         comply.RuleMinAnnualLeaveDays,
			Name:        "Minimum Annual Leave",
			Description: "Minimum 20 working days of paid annual leave.",
			Category:    comply.CatLeave,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1986, time.April, 1), Amount: 20, Unit: comply.Days, Per: comply.PerYear},
			},
			Source: comply.Source{Title: kt.Title, Section: "Art. 155(4)", URL: kt.URL},
		},
	}
}
