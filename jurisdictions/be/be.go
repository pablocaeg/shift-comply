// Package be registers Belgium's healthcare scheduling regulations:
// Loi sur le travail (Labour Act of March 16, 1971) and Arrete royal
// provisions for healthcare workers.
package be

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:      comply.BE,
		Name:      "Belgium",
		LocalName: "Belgique",
		Type:      comply.Country,
		Parent:    comply.EU,
		Currency:  "EUR",
		TimeZone:  "Europe/Brussels",
		Rules:     rules(),
	}
}

func rules() []*comply.RuleDef {
	lt := comply.Source{
		Title:   "Loi sur le travail du 16 mars 1971",
		Section: "",
		URL:     "https://www.ejustice.just.fgov.be/cgi_loi/change_lg.pl?language=fr&la=F&cn=1971031602&table_name=loi",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxDailyHours,
			Name:        "Maximum Daily Working Hours",
			Description: "Working time may not exceed 8 hours per day (9 hours in a 5-day week system, 10 hours for certain shift patterns).",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1971, time.March, 16), Amount: 8, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: lt.Title, Section: "Art. 19", URL: lt.URL},
		},
		{
			Key:         comply.RuleMaxOrdinaryWeeklyHours,
			Name:        "Maximum Ordinary Weekly Hours",
			Description: "Standard weekly working time is 38 hours (since January 1, 2003).",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2003, time.January, 1), Amount: 38, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: lt.Title, Section: "Art. 19(2)", URL: lt.URL},
		},
		{
			Key:         comply.RuleMaxWeeklyHours,
			Name:        "Absolute Maximum Weekly Hours",
			Description: "Working time including overtime may not exceed 48 hours per week averaged over a reference period. Internal limit of 50 hours in any single week (11 hours overtime).",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2003, time.January, 1),
					Amount:   48,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 4, Unit: comply.PeriodMonths},
				},
			},
			Source: comply.Source{Title: lt.Title, Section: "Art. 26bis", URL: lt.URL},
		},
		{
			Key:         comply.RuleMinRestBetweenShifts,
			Name:        "Minimum Daily Rest Period",
			Description: "At least 11 consecutive hours of rest per 24-hour period.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1971, time.March, 16), Amount: 11, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: lt.Title, Section: "Art. 38ter", URL: lt.URL},
		},
		{
			Key:         comply.RuleMinWeeklyRest,
			Name:        "Minimum Weekly Rest Period",
			Description: "At least 35 consecutive hours of weekly rest (24 hours + 11 hours daily rest). Sunday rest is mandatory in principle.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1971, time.March, 16), Amount: 35, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: lt.Title, Section: "Art. 11", URL: lt.URL},
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Break Requirement",
			Description: "Workers may not work more than 6 consecutive hours without a break.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1971, time.March, 16), Amount: 6, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{Title: lt.Title, Section: "Art. 38quater", URL: lt.URL},
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Break Duration",
			Description: "At least 15 minutes. Collective agreements typically grant 30-60 minutes.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1971, time.March, 16), Amount: 15, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{Title: lt.Title, Section: "Art. 38quater", URL: lt.URL},
		},
		{
			Key:         comply.RuleNightPeriodStart,
			Name:        "Night Period Start",
			Description: "Night work is defined as work between 20:00 and 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1971, time.March, 16), Amount: 20, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: lt.Title, Section: "Art. 35(1)", URL: lt.URL},
			Notes:  "Belgium has one of the earliest night period starts in Europe (20:00 vs 22:00-23:00 elsewhere).",
		},
		{
			Key:         comply.RuleNightPeriodEnd,
			Name:        "Night Period End",
			Description: "Night work period ends at 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1971, time.March, 16), Amount: 6, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: lt.Title, Section: "Art. 35(1)", URL: lt.URL},
		},
		{
			Key:         comply.RuleMinAnnualLeaveDays,
			Name:        "Minimum Annual Leave",
			Description: "Minimum 20 working days (4 weeks) of paid annual leave per year.",
			Category:    comply.CatLeave,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1971, time.March, 16), Amount: 20, Unit: comply.Days, Per: comply.PerYear},
			},
			Source: comply.Source{
				Title:   "Arrete royal du 30 mars 1967 (conges annuels)",
				Section: "Art. 3",
				URL:     "https://www.ejustice.just.fgov.be/cgi_loi/change_lg.pl?language=fr&la=F&cn=1967033004&table_name=loi",
			},
		},
		{
			Key:         comply.RuleMaxShiftHours,
			Name:        "Maximum Shift Duration (Healthcare)",
			Description: "In healthcare, shifts may be up to 12 hours with appropriate rest periods. 24-hour on-call possible under collective agreement.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:      comply.D(2003, time.January, 1),
					Amount:     12,
					Unit:       comply.Hours,
					Per:        comply.PerShift,
					Exceptions: []string{"24-hour on-call possible under sectoral collective agreement for healthcare"},
				},
			},
			Source: comply.Source{Title: lt.Title, Section: "Art. 27", URL: lt.URL},
		},
		{
			Key:         comply.RuleMaxOvertimeAnnual,
			Name:        "Maximum Annual Overtime",
			Description: "Maximum 120 hours of overtime per quarter. Annual cap varies by sector (typically 130-143 hours, extendable to 220 by agreement).",
			Category:    comply.CatOvertime,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2003, time.January, 1), Amount: 143, Unit: comply.Hours, Per: comply.PerYear},
			},
			Source: comply.Source{Title: lt.Title, Section: "Art. 26bis S 1", URL: lt.URL},
			Notes:  "Base quota is 91 hours per year, raised to 143 by most sector agreements. Can reach 220 by individual agreement.",
		},
	}
}
