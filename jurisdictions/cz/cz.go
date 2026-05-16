// Package cz registers Czech Republic's healthcare scheduling regulations:
// Zakonik prace (Labour Code, Act 262/2006 Coll.) as amended.
package cz

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:      comply.CZ,
		Name:      "Czech Republic",
		LocalName: "Cesko",
		Type:      comply.Country,
		Parent:    comply.EU,
		Currency:  "CZK",
		TimeZone:  "Europe/Prague",
		Rules:     rules(),
	}
}

func rules() []*comply.RuleDef {
	zp := comply.Source{
		Title:   "Zakonik prace (Act No. 262/2006 Coll., Labour Code)",
		Section: "",
		URL:     "https://www.zakonyprolidi.cz/cs/2006-262",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxOrdinaryWeeklyHours,
			Name:        "Maximum Ordinary Weekly Hours",
			Description: "Standard weekly working time is 40 hours. Reduced for specific conditions: 37.5 hours for multi-shift, underground, or healthcare workers.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2007, time.January, 1), Amount: 40, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: zp.Title, Section: "S 79", URL: zp.URL},
			Notes:  "37.5 hours for 2-shift operations; 37.5 hours for underground work; 38.75 hours for 3-shift and continuous operations.",
		},
		{
			Key:         comply.RuleMaxWeeklyHours,
			Name:        "Maximum Weekly Hours Including Overtime",
			Description: "Working time including overtime shall not exceed 48 hours per week averaged over 26 consecutive weeks (52 weeks by collective agreement).",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2007, time.January, 1),
					Amount:   48,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 6, Unit: comply.PeriodMonths},
				},
			},
			Source: comply.Source{Title: zp.Title, Section: "S 93a", URL: zp.URL},
		},
		{
			Key:         comply.RuleMaxShiftHours,
			Name:        "Maximum Shift Duration",
			Description: "A single shift may not exceed 12 hours.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2007, time.January, 1), Amount: 12, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{Title: zp.Title, Section: "S 83", URL: zp.URL},
		},
		{
			Key:         comply.RuleMinRestBetweenShifts,
			Name:        "Minimum Rest Between Shifts",
			Description: "At least 11 consecutive hours of rest between the end of one shift and the beginning of the next within 24 hours.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:      comply.D(2007, time.January, 1),
					Amount:     11,
					Unit:       comply.Hours,
					Per:        comply.PerDay,
					Exceptions: []string{"May be reduced to 8 hours for healthcare, agriculture, or continuous operations with compensatory rest within 2 weeks"},
				},
			},
			Source: comply.Source{Title: zp.Title, Section: "S 90", URL: zp.URL},
		},
		{
			Key:         comply.RuleMinWeeklyRest,
			Name:        "Minimum Weekly Rest Period",
			Description: "At least 35 consecutive hours of uninterrupted rest per 7-day period.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2007, time.January, 1), Amount: 35, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: zp.Title, Section: "S 92", URL: zp.URL},
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Break Requirement",
			Description: "A break of at least 30 minutes after 6 consecutive hours of work. Meal break must be granted after 6 hours maximum.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2007, time.January, 1), Amount: 6, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{Title: zp.Title, Section: "S 88", URL: zp.URL},
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Break Duration",
			Description: "At least 30 minutes. May be split into parts of at least 15 minutes each.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2007, time.January, 1), Amount: 30, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{Title: zp.Title, Section: "S 88", URL: zp.URL},
		},
		{
			Key:         comply.RuleMaxOvertimeAnnual,
			Name:        "Maximum Annual Overtime",
			Description: "Maximum 150 hours of ordered overtime per year. Up to 416 hours with employee consent (8 hours/week average).",
			Category:    comply.CatOvertime,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:      comply.D(2007, time.January, 1),
					Amount:     150,
					Unit:       comply.Hours,
					Per:        comply.PerYear,
					Exceptions: []string{"Up to 416 hours per year with individual employee consent"},
				},
			},
			Source: comply.Source{Title: zp.Title, Section: "S 93", URL: zp.URL},
		},
		{
			Key:         comply.RuleNightPeriodStart,
			Name:        "Night Period Start",
			Description: "Night work is defined as work between 22:00 and 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2007, time.January, 1), Amount: 22, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: zp.Title, Section: "S 78(1)(j)", URL: zp.URL},
		},
		{
			Key:         comply.RuleNightPeriodEnd,
			Name:        "Night Period End",
			Description: "Night work period ends at 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2007, time.January, 1), Amount: 6, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: zp.Title, Section: "S 78(1)(j)", URL: zp.URL},
		},
		{
			Key:         comply.RuleMaxNightShiftHours,
			Name:        "Maximum Night Worker Hours",
			Description: "Night workers may not work more than 8 hours on average per 24-hour period over 26 weeks.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2007, time.January, 1),
					Amount:   8,
					Unit:     comply.Hours,
					Per:      comply.PerDay,
					Averaged: &comply.AveragingPeriod{Count: 6, Unit: comply.PeriodMonths},
				},
			},
			Source: comply.Source{Title: zp.Title, Section: "S 94(1)", URL: zp.URL},
		},
		{
			Key:         comply.RuleMinAnnualLeaveDays,
			Name:        "Minimum Annual Leave",
			Description: "Minimum 20 working days (4 weeks) per year. Public sector employees: 25 working days (5 weeks).",
			Category:    comply.CatLeave,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2007, time.January, 1), Amount: 20, Unit: comply.Days, Per: comply.PerYear},
			},
			Source: comply.Source{Title: zp.Title, Section: "S 213", URL: zp.URL},
			Notes:  "Public sector (including hospitals): 5 weeks (25 days). Private sector: 4 weeks (20 days).",
		},
	}
}
