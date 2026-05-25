// Package lu registers Luxembourg's healthcare scheduling regulations:
// Code du travail (Labour Code), specifically Book III, Title I.
package lu

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.LU,
		Name:     "Luxembourg",
		Type:     comply.Country,
		Parent:   comply.EU,
		Currency: "EUR",
		TimeZone: "Europe/Luxembourg",
		Rules:    rules(),
	}
}

func rules() []*comply.RuleDef {
	ct := comply.Source{
		Title:   "Code du travail (Luxembourg)",
		Section: "",
		URL:     "https://legilux.public.lu/eli/etat/leg/code/travail/",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxOrdinaryWeeklyHours,
			Name:        "Maximum Ordinary Weekly Hours",
			Description: "Normal working time may not exceed 8 hours per day and 40 hours per week.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2006, time.September, 1), Amount: 40, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. L.211-5", URL: ct.URL},
		},
		{
			Key:         comply.RuleMaxDailyHours,
			Name:        "Maximum Daily Working Hours",
			Description: "Daily working time may not exceed 8 hours. May be extended to 10 hours by collective agreement.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2006, time.September, 1), Amount: 8, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. L.211-5", URL: ct.URL},
		},
		{
			Key:         comply.RuleMaxWeeklyHours,
			Name:        "Maximum Weekly Hours Including Overtime",
			Description: "Total working time including overtime may not exceed 48 hours per week averaged over a 4-week reference period.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2006, time.September, 1),
					Amount:   48,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 4, Unit: comply.PeriodWeeks},
				},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. L.211-12", URL: ct.URL},
		},
		{
			Key:         comply.RuleMinRestBetweenShifts,
			Name:        "Minimum Daily Rest Period",
			Description: "At least 11 consecutive hours of rest per 24-hour period.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2006, time.September, 1), Amount: 11, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. L.211-16", URL: ct.URL},
		},
		{
			Key:         comply.RuleMinWeeklyRest,
			Name:        "Minimum Weekly Rest Period",
			Description: "At least 44 consecutive hours of weekly rest, normally including Sunday.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2006, time.September, 1), Amount: 44, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. L.231-1", URL: ct.URL},
			Notes:  "44 hours weekly rest is one of the highest in Europe.",
		},
		{
			Key:         comply.RuleNightPeriodStart,
			Name:        "Night Period Start",
			Description: "Night work is defined as work between 22:00 and 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2006, time.September, 1), Amount: 22, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. L.211-22", URL: ct.URL},
		},
		{
			Key:         comply.RuleNightPeriodEnd,
			Name:        "Night Period End",
			Description: "Night work period ends at 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2006, time.September, 1), Amount: 6, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. L.211-22", URL: ct.URL},
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Break Requirement",
			Description: "A break is mandatory when daily working time exceeds 6 hours.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2006, time.September, 1), Amount: 6, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. L.211-16", URL: ct.URL},
		},
		{
			Key:         comply.RuleMaxOvertimeAnnual,
			Name:        "Maximum Annual Overtime",
			Description: "Maximum 150 hours of overtime per year per employee. Additional overtime requires authorization from the Inspectorate of Labour.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2006, time.September, 1), Amount: 150, Unit: comply.Hours, Per: comply.PerYear},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. L.211-27", URL: ct.URL},
		},
		{
			Key:         comply.RuleMaxNightShiftHours,
			Name:        "Maximum Night Worker Hours",
			Description: "Night workers may not work more than 8 hours on average per 24-hour period.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2006, time.September, 1), Amount: 8, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. L.211-22", URL: ct.URL},
		},
		{
			Key:         comply.RuleMinAnnualLeaveDays,
			Name:        "Minimum Annual Leave",
			Description: "Minimum 26 working days of paid annual leave per year.",
			Category:    comply.CatLeave,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2019, time.January, 1), Amount: 26, Unit: comply.Days, Per: comply.PerYear},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. L.233-4", URL: ct.URL},
			Notes:  "Increased from 25 to 26 days in 2019. One of the highest statutory minimums in Europe.",
		},
	}
}
