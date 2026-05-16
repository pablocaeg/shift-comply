// Package nl registers Netherlands healthcare scheduling regulations:
// Arbeidstijdenwet (Working Time Act, ATW) and Arbeidstijdenbesluit (ATB).
package nl

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:      comply.NL,
		Name:      "Netherlands",
		LocalName: "Nederland",
		Type:      comply.Country,
		Parent:    comply.EU,
		Currency:  "EUR",
		TimeZone:  "Europe/Amsterdam",
		Rules:     rules(),
	}
}

func rules() []*comply.RuleDef {
	atw := comply.Source{
		Title:   "Arbeidstijdenwet (Working Time Act, ATW)",
		Section: "",
		URL:     "https://wetten.overheid.nl/BWBR0007671/",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxShiftHours,
			Name:        "Maximum Shift Duration",
			Description: "A single shift may not exceed 12 hours.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1996, time.January, 1), Amount: 12, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{Title: atw.Title, Section: "Art. 5:7(2)", URL: atw.URL},
		},
		{
			Key:         comply.RuleMaxWeeklyHours,
			Name:        "Maximum Weekly Working Hours",
			Description: "Maximum 48 hours per week averaged over 16 weeks. No single week may exceed 60 hours.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(1996, time.January, 1),
					Amount:   48,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 16, Unit: comply.PeriodWeeks},
					Exceptions: []string{
						"Absolute cap of 60 hours in any single week",
					},
				},
			},
			Source: comply.Source{Title: atw.Title, Section: "Art. 5:7(4)", URL: atw.URL},
		},
		{
			Key:         comply.RuleMinRestBetweenShifts,
			Name:        "Minimum Daily Rest Period",
			Description: "At least 11 consecutive hours of rest per 24-hour period.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:      comply.D(1996, time.January, 1),
					Amount:     11,
					Unit:       comply.Hours,
					Per:        comply.PerDay,
					Exceptions: []string{"May be reduced to 8 hours once per 7-day period"},
				},
			},
			Source: comply.Source{Title: atw.Title, Section: "Art. 5:3(2)", URL: atw.URL},
		},
		{
			Key:         comply.RuleMinWeeklyRest,
			Name:        "Minimum Weekly Rest Period",
			Description: "At least 36 consecutive hours of rest per 7-day period. Alternatively, 72 hours per 14-day period, split into periods of at least 32 hours.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1996, time.January, 1), Amount: 36, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: atw.Title, Section: "Art. 5:5", URL: atw.URL},
		},
		{
			Key:         comply.RuleMaxConsecutiveDays,
			Name:        "Maximum Consecutive Working Days",
			Description: "May not work more than 7 consecutive days.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1996, time.January, 1), Amount: 7, Unit: comply.Days, Per: comply.PerPeriod},
			},
			Source: comply.Source{Title: atw.Title, Section: "Art. 5:5(4)", URL: atw.URL},
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Break Requirement",
			Description: "At least 30 minutes break after 5.5 hours of work. May be split into 2 x 15 minutes.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1996, time.January, 1), Amount: 5.5, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{Title: atw.Title, Section: "Art. 5:4(1)", URL: atw.URL},
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Break Duration",
			Description: "At least 30 minutes after 5.5 hours. At least 45 minutes after 10 hours (may be split into 3 x 15 min).",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1996, time.January, 1), Amount: 30, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{Title: atw.Title, Section: "Art. 5:4(1)", URL: atw.URL},
		},
		{
			Key:         comply.RuleMaxConsecutiveNights,
			Name:        "Maximum Consecutive Night Shifts",
			Description: "May not work more than 7 consecutive night shifts. After the last night shift, at least 46 hours of rest.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1996, time.January, 1), Amount: 7, Unit: comply.Count, Per: comply.PerPeriod},
			},
			Source: comply.Source{Title: atw.Title, Section: "Art. 5:8(3)", URL: atw.URL},
		},
		{
			Key:         comply.RuleNightPeriodStart,
			Name:        "Night Period Start",
			Description: "Night work is defined as work between 00:00 and 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1996, time.January, 1), Amount: 0, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: atw.Title, Section: "Art. 1:7", URL: atw.URL},
		},
		{
			Key:         comply.RuleNightPeriodEnd,
			Name:        "Night Period End",
			Description: "Night work period ends at 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1996, time.January, 1), Amount: 6, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: atw.Title, Section: "Art. 1:7", URL: atw.URL},
		},
		{
			Key:         comply.RuleMinAnnualLeaveDays,
			Name:        "Minimum Annual Leave",
			Description: "Minimum 20 working days (4 x weekly working days) of paid annual leave.",
			Category:    comply.CatLeave,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1996, time.January, 1), Amount: 20, Unit: comply.Days, Per: comply.PerYear},
			},
			Source: comply.Source{
				Title:   "Burgerlijk Wetboek (Civil Code)",
				Section: "Book 7, Art. 634",
				URL:     "https://wetten.overheid.nl/BWBR0005290/",
			},
		},
	}
}
