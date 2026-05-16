// Package hr registers Croatia's healthcare scheduling regulations:
// Zakon o radu (Labour Act, NN 93/14, 127/17, 98/19, 151/22).
package hr

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:      comply.HR,
		Name:      "Croatia",
		LocalName: "Hrvatska",
		Type:      comply.Country,
		Parent:    comply.EU,
		Currency:  "EUR",
		TimeZone:  "Europe/Zagreb",
		Rules:     rules(),
	}
}

func rules() []*comply.RuleDef {
	zr := comply.Source{
		Title:   "Zakon o radu (Labour Act, NN 93/14, 127/17, 98/19, 151/22)",
		Section: "",
		URL:     "https://www.zakon.hr/z/307/Zakon-o-radu",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxOrdinaryWeeklyHours,
			Name:        "Maximum Ordinary Weekly Hours",
			Description: "Full-time working hours may not exceed 40 hours per week.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2014, time.August, 7), Amount: 40, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: zr.Title, Section: "Art. 61", URL: zr.URL},
		},
		{
			Key:         comply.RuleMaxWeeklyHours,
			Name:        "Maximum Weekly Hours Including Overtime",
			Description: "Average weekly working time including overtime shall not exceed 48 hours per 4-month period (6 months by collective agreement).",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2014, time.August, 7),
					Amount:   48,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 4, Unit: comply.PeriodMonths},
				},
			},
			Source: comply.Source{Title: zr.Title, Section: "Art. 66", URL: zr.URL},
		},
		{
			Key:         comply.RuleMinRestBetweenShifts,
			Name:        "Minimum Daily Rest Period",
			Description: "At least 12 consecutive hours of uninterrupted rest per 24-hour period.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2014, time.August, 7), Amount: 12, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: zr.Title, Section: "Art. 74", URL: zr.URL},
			Notes:  "Croatia requires 12 hours (above EU minimum of 11). This is the same as Romania and Spain.",
		},
		{
			Key:         comply.RuleMinWeeklyRest,
			Name:        "Minimum Weekly Rest Period",
			Description: "At least 24 consecutive hours of uninterrupted weekly rest, plus the daily rest period. Total minimum 36 hours.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2014, time.August, 7), Amount: 36, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: zr.Title, Section: "Art. 75", URL: zr.URL},
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Break Requirement",
			Description: "Workers are entitled to a break when daily working time is at least 6 hours.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2014, time.August, 7), Amount: 6, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{Title: zr.Title, Section: "Art. 73", URL: zr.URL},
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Break Duration",
			Description: "At least 30 minutes. Counted as working time.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2014, time.August, 7), Amount: 30, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{Title: zr.Title, Section: "Art. 73(3)", URL: zr.URL},
		},
		{
			Key:         comply.RuleMaxOvertimeAnnual,
			Name:        "Maximum Annual Overtime",
			Description: "Maximum 180 hours of overtime per year. May be extended to 250 hours by collective agreement.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:      comply.D(2014, time.August, 7),
					Amount:     180,
					Unit:       comply.Hours,
					Per:        comply.PerYear,
					Exceptions: []string{"Up to 250 hours by collective agreement"},
				},
			},
			Source: comply.Source{Title: zr.Title, Section: "Art. 65", URL: zr.URL},
		},
		{
			Key:         comply.RuleNightPeriodStart,
			Name:        "Night Period Start",
			Description: "Night work is defined as work between 22:00 and 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2014, time.August, 7), Amount: 22, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: zr.Title, Section: "Art. 69", URL: zr.URL},
		},
		{
			Key:         comply.RuleNightPeriodEnd,
			Name:        "Night Period End",
			Description: "Night work period ends at 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2014, time.August, 7), Amount: 6, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: zr.Title, Section: "Art. 69", URL: zr.URL},
		},
		{
			Key:         comply.RuleMaxNightShiftHours,
			Name:        "Maximum Night Worker Hours",
			Description: "Night workers may not work more than 8 hours on average per 24-hour period over 4 months.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2014, time.August, 7),
					Amount:   8,
					Unit:     comply.Hours,
					Per:      comply.PerDay,
					Averaged: &comply.AveragingPeriod{Count: 4, Unit: comply.PeriodMonths},
				},
			},
			Source: comply.Source{Title: zr.Title, Section: "Art. 71", URL: zr.URL},
		},
		{
			Key:         comply.RuleMinAnnualLeaveDays,
			Name:        "Minimum Annual Leave",
			Description: "Minimum 20 working days (4 weeks) of paid annual leave per year.",
			Category:    comply.CatLeave,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2014, time.August, 7), Amount: 20, Unit: comply.Days, Per: comply.PerYear},
			},
			Source: comply.Source{Title: zr.Title, Section: "Art. 77", URL: zr.URL},
		},
	}
}
