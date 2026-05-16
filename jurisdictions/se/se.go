// Package se registers Sweden's healthcare scheduling regulations:
// Arbetstidslag (Working Time Act, SFS 1982:673). Sweden does NOT use
// the Article 22 opt-out, making its 48-hour limit absolute like Italy and France.
package se

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:      comply.SE,
		Name:      "Sweden",
		LocalName: "Sverige",
		Type:      comply.Country,
		Parent:    comply.EU,
		Currency:  "SEK",
		TimeZone:  "Europe/Stockholm",
		Rules:     rules(),
	}
}

func rules() []*comply.RuleDef {
	atl := comply.Source{
		Title:   "Arbetstidslag (SFS 1982:673)",
		Section: "",
		URL:     "https://www.riksdagen.se/sv/dokument-och-lagar/dokument/svensk-forfattningssamling/arbetstidslag-1982673_sfs-1982-673/",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxWeeklyHours,
			Name:        "Maximum Weekly Working Hours",
			Description: "Regular working time may not exceed 40 hours per week on average. Including overtime: 48 hours per week averaged over 4 months.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(1982, time.July, 1),
					Amount:   48,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 4, Unit: comply.PeriodMonths},
				},
			},
			Source: comply.Source{Title: atl.Title, Section: "S 5, S 10a", URL: atl.URL},
		},
		{
			Key:         comply.RuleMaxOrdinaryWeeklyHours,
			Name:        "Maximum Ordinary Weekly Hours",
			Description: "Ordinary working time may not exceed 40 hours per week.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1982, time.July, 1), Amount: 40, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: atl.Title, Section: "S 5", URL: atl.URL},
		},
		{
			Key:         comply.RuleMinRestBetweenShifts,
			Name:        "Minimum Daily Rest Period",
			Description: "At least 11 consecutive hours of rest per 24-hour period.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1982, time.July, 1), Amount: 11, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: atl.Title, Section: "S 13", URL: atl.URL},
		},
		{
			Key:         comply.RuleMinWeeklyRest,
			Name:        "Minimum Weekly Rest Period",
			Description: "At least 36 consecutive hours of uninterrupted rest per 7-day period. Should preferably be on weekends.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1982, time.July, 1), Amount: 36, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: atl.Title, Section: "S 14", URL: atl.URL},
		},
		{
			Key:         comply.RuleMaxOvertimeAnnual,
			Name:        "Maximum Annual Overtime",
			Description: "Maximum 200 hours of overtime per calendar year. Additional 150 hours of extra overtime (mertid) possible for part-time workers.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1982, time.July, 1), Amount: 200, Unit: comply.Hours, Per: comply.PerYear},
			},
			Source: comply.Source{Title: atl.Title, Section: "S 8", URL: atl.URL},
			Notes:  "Also limited to 48 hours per 4-week period and 50 hours per calendar month.",
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Break Requirement",
			Description: "Workers may not work more than 5 consecutive hours without a break (rast).",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1982, time.July, 1), Amount: 5, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{Title: atl.Title, Section: "S 15", URL: atl.URL},
		},
		{
			Key:         comply.RuleNightPeriodStart,
			Name:        "Night Period Start",
			Description: "Night time is defined as 22:00 to 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1982, time.July, 1), Amount: 22, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: atl.Title, Section: "S 13a", URL: atl.URL},
		},
		{
			Key:         comply.RuleNightPeriodEnd,
			Name:        "Night Period End",
			Description: "Night work period ends at 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1982, time.July, 1), Amount: 6, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: atl.Title, Section: "S 13a", URL: atl.URL},
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
					Since:    comply.D(1982, time.July, 1),
					Amount:   8,
					Unit:     comply.Hours,
					Per:      comply.PerDay,
					Averaged: &comply.AveragingPeriod{Count: 4, Unit: comply.PeriodMonths},
				},
			},
			Source: comply.Source{Title: atl.Title, Section: "S 13a", URL: atl.URL},
		},
		{
			Key:         comply.RuleMinAnnualLeaveDays,
			Name:        "Minimum Annual Leave",
			Description: "Minimum 25 working days of paid annual leave per year.",
			Category:    comply.CatLeave,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1977, time.April, 1), Amount: 25, Unit: comply.Days, Per: comply.PerYear},
			},
			Source: comply.Source{
				Title:   "Semesterlag (SFS 1977:480)",
				Section: "S 4",
				URL:     "https://www.riksdagen.se/sv/dokument-och-lagar/dokument/svensk-forfattningssamling/semesterlag-1977480_sfs-1977-480/",
			},
		},
		{
			Key:         "se-no-optout",
			Name:        "No Individual Opt-Out (Article 22 Not Adopted)",
			Description: "Sweden does NOT use the EU Article 22 individual opt-out. The 48-hour weekly maximum is absolute. Sweden relies on collective agreements for flexibility but cannot exceed the 48-hour average.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpBool,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1982, time.July, 1), Amount: 0, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "Arbetstidslag (SFS 1982:673)",
				Section: "S 10a (no Art. 22 transposition)",
				URL:     "https://www.riksdagen.se/sv/dokument-och-lagar/dokument/svensk-forfattningssamling/arbetstidslag-1982673_sfs-1982-673/",
			},
			Notes: "Sweden, France, and Italy are the three major EU economies that refuse the Article 22 opt-out.",
		},
	}
}
