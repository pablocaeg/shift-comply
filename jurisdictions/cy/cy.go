// Package cy registers Cyprus's healthcare scheduling regulations:
// Organisation of Working Time Law (N. 63(I)/2002), transposing
// EU Directive 2003/88/EC.
package cy

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:      comply.CY,
		Name:      "Cyprus",
		LocalName: "Kypros",
		Type:      comply.Country,
		Parent:    comply.EU,
		Currency:  "EUR",
		TimeZone:  "Asia/Nicosia",
		Rules:     rules(),
	}
}

func rules() []*comply.RuleDef {
	owt := comply.Source{
		Title:   "Organisation of Working Time Law (N. 63(I)/2002)",
		Section: "",
		URL:     "http://www.cylaw.org/nomoi/enact/ind/2002_1_63/section-sc5bb7db1f-8e86-43d0-abb0-db24f891e40f.html",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxWeeklyHours,
			Name:        "Maximum Weekly Working Hours",
			Description: "Average weekly working time shall not exceed 48 hours over a 4-month reference period.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2002, time.December, 1),
					Amount:   48,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 4, Unit: comply.PeriodMonths},
				},
			},
			Source: comply.Source{Title: owt.Title, Section: "S 4", URL: owt.URL},
		},
		{
			Key:         comply.RuleMinRestBetweenShifts,
			Name:        "Minimum Daily Rest Period",
			Description: "At least 11 consecutive hours of rest per 24-hour period.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2002, time.December, 1), Amount: 11, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: owt.Title, Section: "S 5", URL: owt.URL},
		},
		{
			Key:         comply.RuleMinWeeklyRest,
			Name:        "Minimum Weekly Rest Period",
			Description: "At least 24 consecutive hours of weekly rest plus the 11-hour daily rest (35 hours total).",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2002, time.December, 1), Amount: 35, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: owt.Title, Section: "S 6", URL: owt.URL},
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Break Duration",
			Description: "At least 15 minutes when daily working time exceeds 6 hours. Collective agreements typically provide 30 minutes.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2002, time.December, 1), Amount: 15, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{Title: owt.Title, Section: "S 7", URL: owt.URL},
		},
		{
			Key:         comply.RuleNightPeriodStart,
			Name:        "Night Period Start",
			Description: "Night work is defined as work between 23:00 and 07:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2002, time.December, 1), Amount: 23, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: owt.Title, Section: "S 2", URL: owt.URL},
		},
		{
			Key:         comply.RuleNightPeriodEnd,
			Name:        "Night Period End",
			Description: "Night work period ends at 07:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2002, time.December, 1), Amount: 7, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: owt.Title, Section: "S 2", URL: owt.URL},
		},
		{
			Key:         comply.RuleMaxNightShiftHours,
			Name:        "Maximum Night Worker Hours",
			Description: "Night workers may not work more than 8 hours on average per 24-hour period.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2002, time.December, 1), Amount: 8, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: owt.Title, Section: "S 8", URL: owt.URL},
		},
		{
			Key:         comply.RuleMinAnnualLeaveDays,
			Name:        "Minimum Annual Leave",
			Description: "Minimum 20 working days (4 weeks) of paid annual leave per year.",
			Category:    comply.CatLeave,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2002, time.December, 1), Amount: 20, Unit: comply.Days, Per: comply.PerYear},
			},
			Source: comply.Source{
				Title:   "Annual Paid Leave Law (N. 8/1967 as amended)",
				Section: "S 3",
				URL:     "http://www.cylaw.org/nomoi/enact/ind/1967_1_8/index.html",
			},
		},
	}
}
