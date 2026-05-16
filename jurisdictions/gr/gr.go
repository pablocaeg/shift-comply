// Package gr registers Greece's healthcare scheduling regulations:
// Labour law (Nomos 3850/2010 on health and safety, Presidential Decree
// 88/1999 implementing EU WTD), and healthcare-specific provisions.
package gr

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:      comply.GR,
		Name:      "Greece",
		LocalName: "Ellada",
		Type:      comply.Country,
		Parent:    comply.EU,
		Currency:  "EUR",
		TimeZone:  "Europe/Athens",
		Rules:     rules(),
	}
}

func rules() []*comply.RuleDef {
	pd := comply.Source{
		Title:   "Presidential Decree 88/1999 (EU WTD implementation)",
		Section: "",
		URL:     "https://www.e-nomothesia.gr/kat-ergasia-koinonike-asphalise/pd-88-1999.html",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxWeeklyHours,
			Name:        "Maximum Weekly Working Hours",
			Description: "Average weekly working time shall not exceed 48 hours over a 4-month reference period (6 months by collective agreement).",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(1999, time.March, 29),
					Amount:   48,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 4, Unit: comply.PeriodMonths},
				},
			},
			Source: comply.Source{Title: pd.Title, Section: "Art. 3", URL: pd.URL},
		},
		{
			Key:         comply.RuleMaxOrdinaryWeeklyHours,
			Name:        "Maximum Ordinary Weekly Hours",
			Description: "Standard legal working time is 40 hours per week (5 days x 8 hours) for private sector.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1984, time.January, 1), Amount: 40, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{
				Title:   "Nomos 1876/1990",
				Section: "Art. 1",
				URL:     "https://www.e-nomothesia.gr/kat-ergasia-koinonike-asphalise/n-1876-1990.html",
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
				{Since: comply.D(1999, time.March, 29), Amount: 11, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: pd.Title, Section: "Art. 3", URL: pd.URL},
		},
		{
			Key:         comply.RuleMinWeeklyRest,
			Name:        "Minimum Weekly Rest Period",
			Description: "At least 24 consecutive hours of uninterrupted weekly rest plus the 11-hour daily rest (35 hours). Sunday is the mandatory rest day.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1999, time.March, 29), Amount: 35, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: pd.Title, Section: "Art. 5", URL: pd.URL},
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Break Requirement",
			Description: "A break when daily working time exceeds 6 hours.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1999, time.March, 29), Amount: 6, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{Title: pd.Title, Section: "Art. 4", URL: pd.URL},
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Break Duration",
			Description: "At least 15 minutes (30 minutes for shifts exceeding 8 hours by most collective agreements).",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1999, time.March, 29), Amount: 15, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{Title: pd.Title, Section: "Art. 4", URL: pd.URL},
		},
		{
			Key:         comply.RuleNightPeriodStart,
			Name:        "Night Period Start",
			Description: "Night time is defined as 22:00 to 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1999, time.March, 29), Amount: 22, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: pd.Title, Section: "Art. 6", URL: pd.URL},
		},
		{
			Key:         comply.RuleNightPeriodEnd,
			Name:        "Night Period End",
			Description: "Night time ends at 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1999, time.March, 29), Amount: 6, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: pd.Title, Section: "Art. 6", URL: pd.URL},
		},
		{
			Key:         comply.RuleMaxNightShiftHours,
			Name:        "Maximum Night Worker Hours",
			Description: "Night workers may not work more than 8 hours on average per 24-hour period.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1999, time.March, 29), Amount: 8, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: pd.Title, Section: "Art. 6", URL: pd.URL},
		},
		{
			Key:         comply.RuleMinAnnualLeaveDays,
			Name:        "Minimum Annual Leave",
			Description: "20 working days for 5-day week; 24 working days for 6-day week. Increases by 1 day after 10 years with same employer (to 25/26 days).",
			Category:    comply.CatLeave,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2001, time.January, 1), Amount: 20, Unit: comply.Days, Per: comply.PerYear},
			},
			Source: comply.Source{
				Title:   "Nomos 539/1945 as amended",
				Section: "Art. 1-2",
				URL:     "https://www.e-nomothesia.gr/kat-ergasia-koinonike-asphalise/n-539-1945.html",
			},
		},
		{
			Key:         "gr-healthcare-optout",
			Name:        "Healthcare Opt-Out (Article 22)",
			Description: "Greece uses the EU Article 22 opt-out for healthcare workers. Hospital doctors regularly exceed 48 hours/week with on-call (efimerias). Individual written consent required.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpBool,
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2005, time.January, 1), Amount: 1, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "Nomos 3329/2005",
				Section: "Art. 45",
				URL:     "https://www.e-nomothesia.gr/kat-ygeia/n-3329-2005.html",
			},
			Notes: "Greek hospital doctors routinely work 32-hour shifts (efimeria). Compliance with WTD rest periods remains a persistent challenge.",
		},
		{
			Key:         comply.RuleMaxShiftHours,
			Name:        "Maximum Guard Shift (Efimeria)",
			Description: "Hospital guard duty (efimeria) in public hospitals. Standard efimeria is 24 hours added after a normal work day, resulting in shifts of up to 32 hours in practice.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffPhysician},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:      comply.D(2005, time.January, 1),
					Amount:     24,
					Unit:       comply.Hours,
					Per:        comply.PerShift,
					Exceptions: []string{"May extend to 32 hours when combined with preceding regular shift (7h + 24h efimeria + 1h handover)"},
				},
			},
			Source: comply.Source{
				Title:   "Nomos 3329/2005",
				Section: "Art. 45",
				URL:     "https://www.e-nomothesia.gr/kat-ygeia/n-3329-2005.html",
			},
		},
	}
}
