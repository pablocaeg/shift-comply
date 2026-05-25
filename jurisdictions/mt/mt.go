// Package mt registers Malta's healthcare scheduling regulations:
// Employment and Industrial Relations Act (EIRA, Cap. 452) and
// Organisation of Working Time Regulations (S.L. 452.87).
package mt

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.MT,
		Name:     "Malta",
		Type:     comply.Country,
		Parent:   comply.EU,
		Currency: "EUR",
		TimeZone: "Europe/Malta",
		Rules:    rules(),
	}
}

func rules() []*comply.RuleDef {
	owtr := comply.Source{
		Title:   "Organisation of Working Time Regulations (S.L. 452.87)",
		Section: "",
		URL:     "https://legislation.mt/eli/sl/452.87/eng",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxWeeklyHours,
			Name:        "Maximum Weekly Working Hours",
			Description: "Average weekly working time shall not exceed 48 hours over a 17-week reference period.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2004, time.May, 1),
					Amount:   48,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 17, Unit: comply.PeriodWeeks},
				},
			},
			Source: comply.Source{Title: owtr.Title, Section: "Reg. 14", URL: owtr.URL},
		},
		{
			Key:         comply.RuleMaxOrdinaryWeeklyHours,
			Name:        "Maximum Ordinary Weekly Hours",
			Description: "Normal working time is 40 hours per week.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2004, time.May, 1), Amount: 40, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{
				Title:   "Employment and Industrial Relations Act (EIRA, Cap. 452)",
				Section: "S 13",
				URL:     "https://legislation.mt/eli/cap/452/eng",
			},
		},
		{
			Key:         comply.RuleMinRestBetweenShifts,
			Name:        "Minimum Daily Rest Period",
			Description: "At least 11 consecutive hours of rest per 24-hour period.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2004, time.May, 1), Amount: 11, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: owtr.Title, Section: "Reg. 5", URL: owtr.URL},
		},
		{
			Key:         comply.RuleMinWeeklyRest,
			Name:        "Minimum Weekly Rest Period",
			Description: "At least 24 consecutive hours of weekly rest plus the 11-hour daily rest (35 hours total) per 7-day period.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2004, time.May, 1), Amount: 35, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: owtr.Title, Section: "Reg. 6", URL: owtr.URL},
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Break Duration",
			Description: "At least 15 minutes when daily working time exceeds 6 hours.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2004, time.May, 1), Amount: 15, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{Title: owtr.Title, Section: "Reg. 7", URL: owtr.URL},
		},
		{
			Key:         comply.RuleNightPeriodStart,
			Name:        "Night Period Start",
			Description: "Night time is defined as 22:00 to 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2004, time.May, 1), Amount: 22, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: owtr.Title, Section: "Reg. 2", URL: owtr.URL},
		},
		{
			Key:         comply.RuleNightPeriodEnd,
			Name:        "Night Period End",
			Description: "Night work period ends at 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2004, time.May, 1), Amount: 6, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: owtr.Title, Section: "Reg. 2", URL: owtr.URL},
		},
		{
			Key:         comply.RuleMaxNightShiftHours,
			Name:        "Maximum Night Worker Hours",
			Description: "Night workers may not work more than 8 hours on average per 24-hour period.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2004, time.May, 1), Amount: 8, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: owtr.Title, Section: "Reg. 8", URL: owtr.URL},
		},
		{
			Key:         comply.RuleMinAnnualLeaveDays,
			Name:        "Minimum Annual Leave",
			Description: "Minimum 24 working days (192 hours based on 8-hour days) of paid annual leave.",
			Category:    comply.CatLeave,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2004, time.May, 1), Amount: 24, Unit: comply.Days, Per: comply.PerYear},
			},
			Source: comply.Source{Title: owtr.Title, Section: "Reg. 9", URL: owtr.URL},
			Notes:  "Plus 14 public holidays. Total: 38 days off per year, among the highest in Europe.",
		},
	}
}
