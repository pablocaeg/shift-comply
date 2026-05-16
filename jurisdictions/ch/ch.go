// Package ch registers Switzerland's healthcare scheduling regulations:
// Arbeitsgesetz (ArG, Labour Act, SR 822.11) and Verordnung 1 zum
// Arbeitsgesetz (ArGV 1). Switzerland is not in the EU but applies
// equivalent protections through its own legislation.
package ch

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:      comply.CH,
		Name:      "Switzerland",
		LocalName: "Schweiz",
		Type:      comply.Country,
		Currency:  "CHF",
		TimeZone:  "Europe/Zurich",
		Rules:     rules(),
	}
}

func rules() []*comply.RuleDef {
	arg := comply.Source{
		Title:   "Arbeitsgesetz (ArG, SR 822.11)",
		Section: "",
		URL:     "https://www.fedlex.admin.ch/eli/cc/1966/57_57_57/de",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxWeeklyHours,
			Name:        "Maximum Weekly Working Hours",
			Description: "Maximum 45 hours per week for industrial workers, office staff, technical employees, and retail. 50 hours for all other workers.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:      comply.D(1966, time.February, 1),
					Amount:     50,
					Unit:       comply.Hours,
					Per:        comply.PerWeek,
					Exceptions: []string{"45 hours for industrial workers, office staff, technical employees, and retail staff"},
				},
			},
			Source: comply.Source{Title: arg.Title, Section: "Art. 9(1)", URL: arg.URL},
			Notes:  "Switzerland has no EU WTD 48-hour limit. Its own limits are 45/50 hours depending on category.",
		},
		{
			Key:         comply.RuleMaxDailyHours,
			Name:        "Maximum Daily Working Hours",
			Description: "Maximum 14 hours from start of work to end of work including breaks. Effective max work time is approximately 12.5 hours.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1966, time.February, 1), Amount: 14, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: arg.Title, Section: "Art. 10(3)", URL: arg.URL},
			Notes:  "14 hours is the maximum span (Tagesrahmen) from first to last hour of work, including breaks.",
		},
		{
			Key:         comply.RuleMinRestBetweenShifts,
			Name:        "Minimum Daily Rest Period",
			Description: "At least 11 consecutive hours of rest between two working days. May be reduced to 8 hours once per week if averaged to 11 over 2 weeks.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:      comply.D(1966, time.February, 1),
					Amount:     11,
					Unit:       comply.Hours,
					Per:        comply.PerDay,
					Exceptions: []string{"May be reduced to 8 hours once per 7 days if 11-hour average maintained over 2 weeks"},
				},
			},
			Source: comply.Source{Title: arg.Title, Section: "Art. 15a(1)", URL: arg.URL},
		},
		{
			Key:         comply.RuleMinWeeklyRest,
			Name:        "Minimum Weekly Rest Period",
			Description: "At least 35 consecutive hours of weekly rest (Sunday + preceding or following half-day).",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1966, time.February, 1), Amount: 35, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: arg.Title, Section: "Art. 18-21", URL: arg.URL},
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Break Requirement",
			Description: "15 minutes if work exceeds 5.5 hours; 30 minutes if work exceeds 7 hours; 60 minutes if work exceeds 9 hours.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1966, time.February, 1), Amount: 5.5, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{Title: arg.Title, Section: "Art. 15", URL: arg.URL},
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Break Duration",
			Description: "15 minutes after 5.5 hours; 30 minutes after 7 hours; 60 minutes after 9 hours.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1966, time.February, 1), Amount: 15, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{Title: arg.Title, Section: "Art. 15", URL: arg.URL},
			Notes:  "Scaled: 15 min (5.5h), 30 min (7h), 60 min (9h).",
		},
		{
			Key:         comply.RuleNightPeriodStart,
			Name:        "Night Period Start",
			Description: "Night work is defined as work between 23:00 and 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1966, time.February, 1), Amount: 23, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: arg.Title, Section: "Art. 16", URL: arg.URL},
		},
		{
			Key:         comply.RuleNightPeriodEnd,
			Name:        "Night Period End",
			Description: "Night work period ends at 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1966, time.February, 1), Amount: 6, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: arg.Title, Section: "Art. 16", URL: arg.URL},
		},
		{
			Key:         comply.RuleMaxOvertimeAnnual,
			Name:        "Maximum Annual Overtime",
			Description: "Maximum 170 hours per year (for 45-hour maximum workers) or 140 hours per year (for 50-hour maximum workers).",
			Category:    comply.CatOvertime,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1966, time.February, 1), Amount: 170, Unit: comply.Hours, Per: comply.PerYear},
			},
			Source: comply.Source{Title: arg.Title, Section: "Art. 12(2)", URL: arg.URL},
			Notes:  "170 hours for the 45h/week category; 140 hours for the 50h/week category.",
		},
		{
			Key:         comply.RuleMinAnnualLeaveDays,
			Name:        "Minimum Annual Leave",
			Description: "Minimum 20 working days (4 weeks) per year. 25 days for workers under 20.",
			Category:    comply.CatLeave,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1984, time.January, 1), Amount: 20, Unit: comply.Days, Per: comply.PerYear},
			},
			Source: comply.Source{
				Title:   "Obligationenrecht (OR, Code of Obligations)",
				Section: "Art. 329a",
				URL:     "https://www.fedlex.admin.ch/eli/cc/27/317_321_377/de",
			},
			Notes: "25 days for workers under 20. Many collective agreements provide 5 weeks (25 days) for all.",
		},
	}
}
