// Package no registers Norway's healthcare scheduling regulations:
// Arbeidsmiljoeloven (Working Environment Act, LOV-2005-06-17-62).
// Norway is not in the EU but applies the WTD via the EEA Agreement.
// Norway does NOT use the Article 22 opt-out.
package no

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:      comply.NO,
		Name:      "Norway",
		LocalName: "Norge",
		Type:      comply.Country,
		Parent:    comply.EU,
		Currency:  "NOK",
		TimeZone:  "Europe/Oslo",
		Rules:     rules(),
	}
}

func rules() []*comply.RuleDef {
	aml := comply.Source{
		Title:   "Arbeidsmiljoeloven (LOV-2005-06-17-62)",
		Section: "",
		URL:     "https://lovdata.no/dokument/NL/lov/2005-06-17-62",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxOrdinaryWeeklyHours,
			Name:        "Maximum Ordinary Weekly Hours",
			Description: "Normal working time may not exceed 9 hours per day and 40 hours per week. Shift workers and continuous operations: 36 or 38 hours per week.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2006, time.January, 1), Amount: 40, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: aml.Title, Section: "S 10-4(1)", URL: aml.URL},
		},
		{
			Key:         comply.RuleMaxDailyHours,
			Name:        "Maximum Daily Working Hours",
			Description: "Normal daily working time may not exceed 9 hours.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2006, time.January, 1), Amount: 9, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: aml.Title, Section: "S 10-4(1)", URL: aml.URL},
		},
		{
			Key:         comply.RuleMaxWeeklyHours,
			Name:        "Maximum Weekly Hours Including Overtime",
			Description: "Total working time including overtime may not exceed 48 hours per week averaged over 8 weeks. Absolute single-week cap of 69 hours (with collective agreement).",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2006, time.January, 1),
					Amount:   48,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 8, Unit: comply.PeriodWeeks},
				},
			},
			Source: comply.Source{Title: aml.Title, Section: "S 10-6(5)", URL: aml.URL},
		},
		{
			Key:         comply.RuleMinRestBetweenShifts,
			Name:        "Minimum Daily Rest Period",
			Description: "At least 11 consecutive hours of uninterrupted rest per 24-hour period.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:      comply.D(2006, time.January, 1),
					Amount:     11,
					Unit:       comply.Hours,
					Per:        comply.PerDay,
					Exceptions: []string{"May be reduced to 8 hours by collective agreement with compensatory rest"},
				},
			},
			Source: comply.Source{Title: aml.Title, Section: "S 10-8(1)", URL: aml.URL},
		},
		{
			Key:         comply.RuleMinWeeklyRest,
			Name:        "Minimum Weekly Rest Period",
			Description: "At least 35 consecutive hours of rest per 7-day period. Must normally include Sunday.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2006, time.January, 1), Amount: 35, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: aml.Title, Section: "S 10-8(2)", URL: aml.URL},
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Break Requirement",
			Description: "At least one break when daily working time exceeds 5.5 hours. If 8+ hours, break must be at least 30 minutes.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2006, time.January, 1), Amount: 5.5, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{Title: aml.Title, Section: "S 10-9(1)", URL: aml.URL},
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Break Duration",
			Description: "At least 30 minutes if daily working time is 8 hours or more.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2006, time.January, 1), Amount: 30, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{Title: aml.Title, Section: "S 10-9(1)", URL: aml.URL},
		},
		{
			Key:         comply.RuleMaxOvertimeAnnual,
			Name:        "Maximum Annual Overtime",
			Description: "Maximum 200 hours of overtime per year. Up to 300 hours by collective agreement, 400 hours with Labour Inspection consent.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2006, time.January, 1), Amount: 200, Unit: comply.Hours, Per: comply.PerYear},
			},
			Source: comply.Source{Title: aml.Title, Section: "S 10-6(4)", URL: aml.URL},
		},
		{
			Key:         comply.RuleNightPeriodStart,
			Name:        "Night Period Start",
			Description: "Night work is defined as work between 21:00 and 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2006, time.January, 1), Amount: 21, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: aml.Title, Section: "S 10-11(1)", URL: aml.URL},
		},
		{
			Key:         comply.RuleNightPeriodEnd,
			Name:        "Night Period End",
			Description: "Night work period ends at 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2006, time.January, 1), Amount: 6, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: aml.Title, Section: "S 10-11(1)", URL: aml.URL},
		},
		{
			Key:         comply.RuleMinAnnualLeaveDays,
			Name:        "Minimum Annual Leave",
			Description: "25 working days (21 days by statute + 4 days by agreement). Workers 60+ receive 6 additional days.",
			Category:    comply.CatLeave,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1988, time.January, 1), Amount: 25, Unit: comply.Days, Per: comply.PerYear},
			},
			Source: comply.Source{
				Title:   "Ferieloven (LOV-1988-04-29-21)",
				Section: "S 5(1)",
				URL:     "https://lovdata.no/dokument/NL/lov/1988-04-29-21",
			},
			Notes: "Statutory minimum is 21 days, but virtually all workers receive 25 days through collective agreement (avtalefestet ferie).",
		},
		{
			Key:         "no-no-optout",
			Name:        "No Individual Opt-Out",
			Description: "Norway does NOT use the Article 22 opt-out. The 48-hour weekly maximum averaged over 8 weeks is absolute.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpBool,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2006, time.January, 1), Amount: 0, Unit: comply.Boolean},
			},
			Source: comply.Source{Title: aml.Title, Section: "S 10-6 (no opt-out provision)", URL: aml.URL},
		},
	}
}
