// Package sk registers Slovakia's healthcare scheduling regulations:
// Zakonnik prace (Labour Code, Act 311/2001 Coll.) as amended.
package sk

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:      comply.SK,
		Name:      "Slovakia",
		LocalName: "Slovensko",
		Type:      comply.Country,
		Parent:    comply.EU,
		Currency:  "EUR",
		TimeZone:  "Europe/Bratislava",
		Rules:     rules(),
	}
}

func rules() []*comply.RuleDef {
	zp := comply.Source{
		Title:   "Zakonnik prace (Act No. 311/2001 Coll., Labour Code)",
		Section: "",
		URL:     "https://www.slov-lex.sk/pravne-predpisy/SK/ZZ/2001/311/",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxOrdinaryWeeklyHours,
			Name:        "Maximum Ordinary Weekly Hours",
			Description: "Weekly working time may not exceed 40 hours. Reduced schedules: 38.75 hours for 2-shift, 37.5 hours for 3-shift and continuous operations.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2002, time.April, 1), Amount: 40, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: zp.Title, Section: "S 85", URL: zp.URL},
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
					Since:    comply.D(2002, time.April, 1),
					Amount:   48,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 4, Unit: comply.PeriodMonths},
				},
			},
			Source: comply.Source{Title: zp.Title, Section: "S 85a", URL: zp.URL},
		},
		{
			Key:         comply.RuleMaxShiftHours,
			Name:        "Maximum Shift Duration",
			Description: "A single shift may not exceed 12 hours.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2002, time.April, 1), Amount: 12, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{Title: zp.Title, Section: "S 85(5)", URL: zp.URL},
		},
		{
			Key:         comply.RuleMinRestBetweenShifts,
			Name:        "Minimum Rest Between Shifts",
			Description: "At least 12 consecutive hours of rest between shifts. May be reduced to 8 hours for healthcare workers with compensatory rest.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:      comply.D(2002, time.April, 1),
					Amount:     12,
					Unit:       comply.Hours,
					Per:        comply.PerDay,
					Exceptions: []string{"May be reduced to 8 hours in healthcare with compensatory rest within 30 days"},
				},
			},
			Source: comply.Source{Title: zp.Title, Section: "S 92", URL: zp.URL},
			Notes:  "Slovakia requires 12 hours rest, above the EU 11-hour minimum.",
		},
		{
			Key:         comply.RuleMinWeeklyRest,
			Name:        "Minimum Weekly Rest Period",
			Description: "At least 35 consecutive hours of uninterrupted rest per 7-day period, normally including Sunday.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2002, time.April, 1), Amount: 35, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: zp.Title, Section: "S 93", URL: zp.URL},
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Break Requirement",
			Description: "A break of at least 30 minutes after 6 hours of continuous work.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2002, time.April, 1), Amount: 6, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{Title: zp.Title, Section: "S 91", URL: zp.URL},
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Break Duration",
			Description: "At least 30 minutes. May be split into 15-minute periods.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2002, time.April, 1), Amount: 30, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{Title: zp.Title, Section: "S 91", URL: zp.URL},
		},
		{
			Key:         comply.RuleMaxOvertimeAnnual,
			Name:        "Maximum Annual Overtime",
			Description: "Maximum 150 hours of ordered overtime per year. Up to 400 hours with employee agreement.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2002, time.April, 1), Amount: 150, Unit: comply.Hours, Per: comply.PerYear},
			},
			Source: comply.Source{Title: zp.Title, Section: "S 97(6)", URL: zp.URL},
		},
		{
			Key:         comply.RuleNightPeriodStart,
			Name:        "Night Period Start",
			Description: "Night work is defined as work between 22:00 and 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2002, time.April, 1), Amount: 22, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: zp.Title, Section: "S 98(1)", URL: zp.URL},
		},
		{
			Key:         comply.RuleNightPeriodEnd,
			Name:        "Night Period End",
			Description: "Night work period ends at 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2002, time.April, 1), Amount: 6, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: zp.Title, Section: "S 98(1)", URL: zp.URL},
		},
		{
			Key:         comply.RuleMinAnnualLeaveDays,
			Name:        "Minimum Annual Leave",
			Description: "Minimum 20 working days (4 weeks). 25 working days (5 weeks) for employees aged 33+.",
			Category:    comply.CatLeave,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2002, time.April, 1), Amount: 20, Unit: comply.Days, Per: comply.PerYear},
			},
			Source: comply.Source{Title: zp.Title, Section: "S 103", URL: zp.URL},
			Notes:  "25 days for employees aged 33 and above.",
		},
	}
}
