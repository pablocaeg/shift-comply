// Package fi registers Finland's healthcare scheduling regulations:
// Tyoaikalaki (Working Time Act 872/2019, effective January 1, 2020).
package fi

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:      comply.FI,
		Name:      "Finland",
		LocalName: "Suomi",
		Type:      comply.Country,
		Parent:    comply.EU,
		Currency:  "EUR",
		TimeZone:  "Europe/Helsinki",
		Rules:     rules(),
	}
}

func rules() []*comply.RuleDef {
	tal := comply.Source{
		Title:   "Tyoaikalaki (Working Time Act 872/2019)",
		Section: "",
		URL:     "https://www.finlex.fi/fi/laki/ajantasa/2019/20190872",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxOrdinaryWeeklyHours,
			Name:        "Maximum Ordinary Weekly Hours",
			Description: "Regular working time may not exceed 8 hours per day and 40 hours per week.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2020, time.January, 1), Amount: 40, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: tal.Title, Section: "S 5", URL: tal.URL},
		},
		{
			Key:         comply.RuleMaxWeeklyHours,
			Name:        "Maximum Weekly Hours Including Overtime",
			Description: "Total working time including overtime shall not exceed 48 hours per week averaged over 4 months.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2020, time.January, 1),
					Amount:   48,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 4, Unit: comply.PeriodMonths},
				},
			},
			Source: comply.Source{Title: tal.Title, Section: "S 18", URL: tal.URL},
		},
		{
			Key:         comply.RuleMaxDailyHours,
			Name:        "Maximum Daily Working Hours",
			Description: "Regular daily working time may not exceed 8 hours.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2020, time.January, 1), Amount: 8, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: tal.Title, Section: "S 5", URL: tal.URL},
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
					Since:      comply.D(2020, time.January, 1),
					Amount:     11,
					Unit:       comply.Hours,
					Per:        comply.PerDay,
					Exceptions: []string{"May be reduced to 9 hours temporarily by collective agreement, with compensatory rest"},
				},
			},
			Source: comply.Source{Title: tal.Title, Section: "S 25", URL: tal.URL},
		},
		{
			Key:         comply.RuleMinWeeklyRest,
			Name:        "Minimum Weekly Rest Period",
			Description: "At least 35 consecutive hours of uninterrupted weekly rest. If shift work scheduling makes this impossible, minimum 24 hours.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2020, time.January, 1), Amount: 35, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: tal.Title, Section: "S 27", URL: tal.URL},
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Break Requirement",
			Description: "A break of at least 30 minutes when daily working time exceeds 6 hours. Worker must be able to leave the workplace.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2020, time.January, 1), Amount: 6, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{Title: tal.Title, Section: "S 24", URL: tal.URL},
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Break Duration",
			Description: "At least 30 minutes. Not included in working time if worker may leave workplace.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2020, time.January, 1), Amount: 30, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{Title: tal.Title, Section: "S 24", URL: tal.URL},
		},
		{
			Key:         comply.RuleMaxOvertimeAnnual,
			Name:        "Maximum Annual Overtime",
			Description: "Maximum 250 hours of overtime per calendar year.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2020, time.January, 1), Amount: 250, Unit: comply.Hours, Per: comply.PerYear},
			},
			Source: comply.Source{Title: tal.Title, Section: "S 19", URL: tal.URL},
		},
		{
			Key:         comply.RuleNightPeriodStart,
			Name:        "Night Period Start",
			Description: "Night work is defined as work performed between 23:00 and 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2020, time.January, 1), Amount: 23, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: tal.Title, Section: "S 8", URL: tal.URL},
		},
		{
			Key:         comply.RuleNightPeriodEnd,
			Name:        "Night Period End",
			Description: "Night work period ends at 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2020, time.January, 1), Amount: 6, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: tal.Title, Section: "S 8", URL: tal.URL},
		},
		{
			Key:         comply.RuleMinAnnualLeaveDays,
			Name:        "Minimum Annual Leave",
			Description: "24 working days for less than 1 year of service; 30 working days for 1+ year.",
			Category:    comply.CatLeave,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2005, time.April, 1), Amount: 24, Unit: comply.Days, Per: comply.PerYear},
			},
			Source: comply.Source{
				Title:   "Vuosilomalaki (Annual Holidays Act 162/2005)",
				Section: "S 5",
				URL:     "https://www.finlex.fi/fi/laki/ajantasa/2005/20050162",
			},
			Notes: "30 days for employees with 1+ year of service by March 31.",
		},
	}
}
