// Package pl registers Poland's healthcare scheduling regulations:
// Kodeks Pracy (Labour Code), and Ustawa o dzialalnosci leczniczej
// (Act on Healthcare Entities) which governs medical on-call duty
// (dyzur medyczny) and Poland's extensive use of the EU opt-out.
package pl

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

// New returns the Poland jurisdiction definition.
func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:      comply.PL,
		Name:      "Poland",
		LocalName: "Polska",
		Type:      comply.Country,
		Parent:    comply.EU,
		Currency:  "PLN",
		TimeZone:  "Europe/Warsaw",
		Rules:     rules(),
	}
}

func rules() []*comply.RuleDef {
	r := make([]*comply.RuleDef, 0, 15)
	r = append(r, labourCodeRules()...)
	r = append(r, healthcareRules()...)
	return r
}

// Kodeks Pracy (Labour Code) - general provisions.

func labourCodeRules() []*comply.RuleDef {
	kp := comply.Source{
		Title:   "Kodeks Pracy (Labour Code, Ustawa z dnia 26 czerwca 1974 r.)",
		Section: "",
		URL:     "https://isap.sejm.gov.pl/isap.nsf/DocDetails.xsp?id=wdu19740240141",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxDailyHours,
			Name:        "Maximum Daily Working Hours",
			Description: "Standard daily working time is 8 hours in a 5-day work week.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1974, time.June, 26), Amount: 8, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: kp.Title, Section: "Art. 129 S 1", URL: kp.URL},
		},
		{
			Key:         comply.RuleMaxOrdinaryWeeklyHours,
			Name:        "Maximum Ordinary Weekly Working Hours",
			Description: "Standard weekly working time is 40 hours in a 5-day work week, averaged over the reference period.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(1974, time.June, 26),
					Amount:   40,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 4, Unit: comply.PeriodMonths},
				},
			},
			Source: comply.Source{Title: kp.Title, Section: "Art. 129 S 1", URL: kp.URL},
		},
		{
			Key:         comply.RuleMinRestBetweenShifts,
			Name:        "Minimum Daily Rest Period",
			Description: "At least 11 consecutive hours of uninterrupted rest per 24-hour period.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1974, time.June, 26), Amount: 11, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: kp.Title, Section: "Art. 132 S 1", URL: kp.URL},
		},
		{
			Key:         comply.RuleMinWeeklyRest,
			Name:        "Minimum Weekly Rest Period",
			Description: "At least 35 consecutive hours of uninterrupted weekly rest, including at least 11 hours of daily rest.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1974, time.June, 26), Amount: 35, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: kp.Title, Section: "Art. 133 S 1", URL: kp.URL},
		},
		{
			Key:         comply.RuleMaxOvertimeAnnual,
			Name:        "Maximum Annual Overtime",
			Description: "Maximum 150 hours of overtime per calendar year per employee. May be increased by collective agreement but must not cause average weekly hours to exceed 48.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1974, time.June, 26), Amount: 150, Unit: comply.Hours, Per: comply.PerYear},
			},
			Source: comply.Source{Title: kp.Title, Section: "Art. 151 S 3", URL: kp.URL},
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Break Requirement",
			Description: "Employees working at least 6 hours are entitled to a 15-minute break included in working time.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1974, time.June, 26), Amount: 6, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{Title: kp.Title, Section: "Art. 134", URL: kp.URL},
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Break Duration",
			Description: "At least 15 minutes, included in working time (paid break).",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1974, time.June, 26), Amount: 15, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{Title: kp.Title, Section: "Art. 134", URL: kp.URL},
		},
		{
			Key:         comply.RuleMinAnnualLeaveDays,
			Name:        "Minimum Annual Leave",
			Description: "20 working days for employees with less than 10 years of service; 26 working days for 10+ years (education counts toward seniority).",
			Category:    comply.CatLeave,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1974, time.June, 26), Amount: 20, Unit: comply.Days, Per: comply.PerYear},
			},
			Source: comply.Source{Title: kp.Title, Section: "Art. 154 S 1", URL: kp.URL},
			Notes:  "26 days for 10+ years seniority (higher education counts as 8 years).",
		},
	}
}

// Healthcare-specific rules - Ustawa o dzialalnosci leczniczej (u.d.l.).

func healthcareRules() []*comply.RuleDef {
	udl := comply.Source{
		Title:   "Ustawa z dnia 15 kwietnia 2011 r. o dzialalnosci leczniczej (Act on Healthcare Entities)",
		Section: "",
		URL:     "https://isap.sejm.gov.pl/isap.nsf/DocDetails.xsp?id=wdu20111120654",
	}

	return []*comply.RuleDef{
		{
			Key:         "pl-healthcare-base-hours",
			Name:        "Healthcare Worker Basic Working Time",
			Description: "Basic working time for healthcare workers providing round-the-clock services is 7 hours 35 minutes per day and 37 hours 55 minutes per week, averaged over a reference period not exceeding 3 months.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2011, time.July, 1),
					Amount:   37.92,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 3, Unit: comply.PeriodMonths},
				},
			},
			Source: comply.Source{Title: udl.Title, Section: "Art. 93", URL: udl.URL},
			Notes:  "7h 35min/day = 7.583h. 37h 55min/week = 37.917h. Lower than the general 40h/week.",
		},
		{
			Key:         "pl-oncall-is-working-time",
			Name:        "Medical On-Call (Dyzur Medyczny) is Working Time",
			Description: "Medical on-call duty (dyzur medyczny) requiring presence at the healthcare facility counts as working time. Time worked during on-call is compensated as overtime unless included in normal working time.",
			Category:    comply.CatOnCall,
			Operator:    comply.OpBool,
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2008, time.January, 1), Amount: 1, Unit: comply.Boolean},
			},
			Source: comply.Source{Title: udl.Title, Section: "Art. 95(2)", URL: udl.URL},
		},
		{
			Key:         "pl-healthcare-optout",
			Name:        "Healthcare Opt-Out (Klauzula opt-out)",
			Description: "Healthcare workers may individually consent in writing to exceed the 48-hour weekly limit when medical on-call duty (dyzur medyczny) is included. Poland uses this extensively. Worker may withdraw consent with 1-month notice.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpBool,
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2008, time.January, 1), Amount: 1, Unit: comply.Boolean},
			},
			Source: comply.Source{Title: udl.Title, Section: "Art. 96", URL: udl.URL},
			Notes:  "Poland is one of the most extensive users of the opt-out in EU healthcare. Without opt-out, most Polish hospitals could not maintain 24-hour on-call staffing. Worker may withdraw with 1-month written notice.",
		},
		{
			Key:         comply.RuleMaxWeeklyHours,
			Name:        "Maximum Weekly Hours Including On-Call",
			Description: "Weekly working time including on-call duty shall not exceed 48 hours averaged over the reference period. With individual opt-out, this limit may be exceeded.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2008, time.January, 1),
					Amount:   48,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 3, Unit: comply.PeriodMonths},
					Exceptions: []string{
						"May be exceeded with individual written opt-out consent (Art. 96 u.d.l.)",
					},
				},
			},
			Source: comply.Source{Title: udl.Title, Section: "Art. 95-96", URL: udl.URL},
		},
		{
			Key:         comply.RuleMinRestAfterGuard,
			Name:        "Rest After On-Call Duty (Dyzur Medyczny)",
			Description: "After medical on-call duty, the worker must receive the equivalent of the daily rest period (11 hours) immediately following the end of duty.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2008, time.January, 1), Amount: 11, Unit: comply.Hours, Per: comply.PerOccurrence},
			},
			Source: comply.Source{Title: udl.Title, Section: "Art. 97", URL: udl.URL},
		},
		{
			Key:         comply.RuleMaxShiftHours,
			Name:        "Maximum On-Call Shift Duration",
			Description: "A single medical on-call duty period (dyzur medyczny) may not exceed 24 hours.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2008, time.January, 1), Amount: 24, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{Title: udl.Title, Section: "Art. 95(3)", URL: udl.URL},
		},
	}
}
