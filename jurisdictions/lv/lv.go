// Package lv registers Latvia's healthcare scheduling regulations:
// Darba likums (Labour Law, adopted June 20, 2001).
package lv

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:      comply.LV,
		Name:      "Latvia",
		LocalName: "Latvija",
		Type:      comply.Country,
		Parent:    comply.EU,
		Currency:  "EUR",
		TimeZone:  "Europe/Riga",
		Rules:     rules(),
	}
}

func rules() []*comply.RuleDef {
	dl := comply.Source{
		Title:   "Darba likums (Labour Law, 2001)",
		Section: "",
		URL:     "https://likumi.lv/ta/id/26019",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxOrdinaryWeeklyHours,
			Name:        "Maximum Ordinary Weekly Hours",
			Description: "Normal working time may not exceed 40 hours per week (8 hours per day).",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2002, time.June, 1), Amount: 40, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: dl.Title, Section: "S 131", URL: dl.URL},
		},
		{
			Key:         comply.RuleMaxWeeklyHours,
			Name:        "Maximum Weekly Hours Including Overtime",
			Description: "Average weekly working time including overtime shall not exceed 48 hours over a reference period up to 4 months.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2002, time.June, 1),
					Amount:   48,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 4, Unit: comply.PeriodMonths},
				},
			},
			Source: comply.Source{Title: dl.Title, Section: "S 136(4)", URL: dl.URL},
		},
		{
			Key:         comply.RuleMinRestBetweenShifts,
			Name:        "Minimum Daily Rest Period",
			Description: "At least 12 consecutive hours of rest per 24-hour period.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2002, time.June, 1), Amount: 12, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: dl.Title, Section: "S 139(1)", URL: dl.URL},
			Notes:  "Latvia requires 12 hours rest, above the EU 11-hour minimum.",
		},
		{
			Key:         comply.RuleMinWeeklyRest,
			Name:        "Minimum Weekly Rest Period",
			Description: "At least 42 consecutive hours of uninterrupted weekly rest.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2002, time.June, 1), Amount: 42, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: dl.Title, Section: "S 140(1)", URL: dl.URL},
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Break Duration",
			Description: "At least 30 minutes if daily working time exceeds 6 hours. Not included in working time unless specified in employment contract.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2002, time.June, 1), Amount: 30, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{Title: dl.Title, Section: "S 145", URL: dl.URL},
		},
		{
			Key:         comply.RuleMaxOvertimeAnnual,
			Name:        "Maximum Annual Overtime",
			Description: "Maximum 144 hours of overtime per 6-month period (approximately 288 hours per year).",
			Category:    comply.CatOvertime,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2002, time.June, 1), Amount: 144, Unit: comply.Hours, Per: comply.PerPeriod},
			},
			Source: comply.Source{Title: dl.Title, Section: "S 136(2)", URL: dl.URL},
			Notes:  "144 hours per 6-month period. Also limited to 8 hours per week.",
		},
		{
			Key:         comply.RuleNightPeriodStart,
			Name:        "Night Period Start",
			Description: "Night work is defined as work between 22:00 and 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2002, time.June, 1), Amount: 22, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: dl.Title, Section: "S 138(1)", URL: dl.URL},
		},
		{
			Key:         comply.RuleNightPeriodEnd,
			Name:        "Night Period End",
			Description: "Night work period ends at 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2002, time.June, 1), Amount: 6, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: dl.Title, Section: "S 138(1)", URL: dl.URL},
		},
		{
			Key:         comply.RuleMinAnnualLeaveDays,
			Name:        "Minimum Annual Leave",
			Description: "Minimum 20 working days (4 calendar weeks) of paid annual leave.",
			Category:    comply.CatLeave,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2002, time.June, 1), Amount: 20, Unit: comply.Days, Per: comply.PerYear},
			},
			Source: comply.Source{Title: dl.Title, Section: "S 149(1)", URL: dl.URL},
		},
	}
}
