// Package at registers Austria's healthcare scheduling regulations:
// Arbeitszeitgesetz (ArbZG), Krankenanstalten-Arbeitszeitgesetz (KA-AZG)
// which is the hospital-specific working time act.
package at

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:      comply.AT,
		Name:      "Austria",
		LocalName: "Oesterreich",
		Type:      comply.Country,
		Parent:    comply.EU,
		Currency:  "EUR",
		TimeZone:  "Europe/Vienna",
		Rules:     rules(),
	}
}

func rules() []*comply.RuleDef {
	r := make([]*comply.RuleDef, 0, 15)
	r = append(r, generalRules()...)
	r = append(r, hospitalRules()...)
	return r
}

func generalRules() []*comply.RuleDef {
	azg := comply.Source{
		Title:   "Arbeitszeitgesetz (AZG, BGBl. Nr. 461/1969)",
		Section: "",
		URL:     "https://www.ris.bka.gv.at/GeltendeFassung.wxe?Abfrage=Bundesnormen&Gesetzesnummer=10008238",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxDailyHours,
			Name:        "Maximum Daily Working Hours",
			Description: "Normal daily working time may not exceed 8 hours. May be extended to 10 hours with flexible working time models.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:      comply.D(1969, time.December, 12),
					Amount:     8,
					Unit:       comply.Hours,
					Per:        comply.PerDay,
					Exceptions: []string{"May be extended to 10 hours with flexible schedules or collective agreement"},
				},
			},
			Source: comply.Source{Title: azg.Title, Section: "S 3", URL: azg.URL},
		},
		{
			Key:         comply.RuleMaxWeeklyHours,
			Name:        "Maximum Weekly Working Hours",
			Description: "Normal weekly working time may not exceed 40 hours. Including overtime: 48 hours averaged over 17 weeks.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(1969, time.December, 12),
					Amount:   48,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 17, Unit: comply.PeriodWeeks},
				},
			},
			Source: comply.Source{Title: azg.Title, Section: "S 9(4)", URL: azg.URL},
		},
		{
			Key:         comply.RuleMinRestBetweenShifts,
			Name:        "Minimum Daily Rest Period",
			Description: "At least 11 consecutive hours of uninterrupted rest per 24-hour period.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1969, time.December, 12), Amount: 11, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: azg.Title, Section: "S 12", URL: azg.URL},
		},
		{
			Key:         comply.RuleNightPeriodStart,
			Name:        "Night Period Start",
			Description: "Night time is defined as 22:00 to 05:00. Austria has an unusually early end (05:00 vs 06:00 in most EU countries).",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1969, time.December, 12), Amount: 22, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: azg.Title, Section: "S 12a", URL: azg.URL},
		},
		{
			Key:         comply.RuleNightPeriodEnd,
			Name:        "Night Period End",
			Description: "Night work period ends at 05:00 (earlier than the 06:00 standard in most EU countries).",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1969, time.December, 12), Amount: 5, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: azg.Title, Section: "S 12a", URL: azg.URL},
		},
		{
			Key:         comply.RuleMinWeeklyRest,
			Name:        "Minimum Weekly Rest (Weekend Rest)",
			Description: "At least 36 consecutive hours of uninterrupted weekly rest, including Sunday.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1969, time.December, 12), Amount: 36, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{
				Title:   "Arbeitsruhegesetz (ARG)",
				Section: "S 3",
				URL:     "https://www.ris.bka.gv.at/GeltendeFassung.wxe?Abfrage=Bundesnormen&Gesetzesnummer=10008541",
			},
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Break Requirement",
			Description: "A break is mandatory after 6 hours of continuous work.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1969, time.December, 12), Amount: 6, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{Title: azg.Title, Section: "S 11", URL: azg.URL},
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Break Duration",
			Description: "At least 30 minutes. May be split into two 15-minute periods if operational reasons require.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1969, time.December, 12), Amount: 30, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{Title: azg.Title, Section: "S 11(1)", URL: azg.URL},
		},
		{
			Key:         comply.RuleMinAnnualLeaveDays,
			Name:        "Minimum Annual Leave",
			Description: "30 working days (5 weeks) for employees with less than 25 years of service. 36 working days (6 weeks) for 25+ years.",
			Category:    comply.CatLeave,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1976, time.July, 7), Amount: 30, Unit: comply.Days, Per: comply.PerYear},
			},
			Source: comply.Source{
				Title:   "Urlaubsgesetz (UrlG)",
				Section: "S 2",
				URL:     "https://www.ris.bka.gv.at/GeltendeFassung.wxe?Abfrage=Bundesnormen&Gesetzesnummer=10008376",
			},
			Notes: "36 working days (6 weeks) for employees with 25+ years of service.",
		},
	}
}

// KA-AZG: Krankenanstalten-Arbeitszeitgesetz (Hospital Working Time Act).

func hospitalRules() []*comply.RuleDef {
	kaazg := comply.Source{
		Title:   "Krankenanstalten-Arbeitszeitgesetz (KA-AZG, BGBl. I Nr. 8/1997)",
		Section: "",
		URL:     "https://www.ris.bka.gv.at/GeltendeFassung.wxe?Abfrage=Bundesnormen&Gesetzesnummer=10009254",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxShiftHours,
			Name:        "Maximum Extended Service (Verlangerte Dienste)",
			Description: "Extended service periods (Verlangerte Dienste) in hospitals may reach 25 hours (including handover). Requires collective agreement.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:      comply.D(2015, time.January, 1),
					Amount:     25,
					Unit:       comply.Hours,
					Per:        comply.PerShift,
					Exceptions: []string{"32 hours possible for on-call with low activity if rest opportunity is guaranteed"},
				},
			},
			Source: comply.Source{Title: kaazg.Title, Section: "S 4", URL: kaazg.URL},
			Notes:  "Before 2015, 32-hour shifts were standard. The 2014 amendment (BGBl. I Nr. 94/2014) reduced the general limit to 25 hours.",
		},
		{
			Key:         comply.RuleMinRestAfterGuard,
			Name:        "Rest After Extended Service",
			Description: "After an extended service period exceeding 13 hours, the subsequent rest period must be at least 11 hours.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2015, time.January, 1), Amount: 11, Unit: comply.Hours, Per: comply.PerOccurrence},
			},
			Source: comply.Source{Title: kaazg.Title, Section: "S 4(5)", URL: kaazg.URL},
		},
		{
			Key:         "at-hospital-optout",
			Name:        "Hospital Opt-Out (Extended Working Time)",
			Description: "Austria uses the EU Article 22 opt-out specifically for hospital workers via the KA-AZG. Individual consent allows exceeding 48 hours/week when extended service periods are included.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpBool,
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2015, time.January, 1), Amount: 1, Unit: comply.Boolean},
			},
			Source: comply.Source{Title: kaazg.Title, Section: "S 4a", URL: kaazg.URL},
			Notes:  "Maximum 60 hours/week average over 17 weeks with opt-out. Worker may withdraw with 3 months notice.",
		},
		{
			Key:         "at-hospital-max-weekly-optout",
			Name:        "Maximum Weekly Hours with Opt-Out",
			Description: "With individual opt-out, hospital workers may work up to an average of 60 hours per week over 17 weeks.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2015, time.January, 1),
					Amount:   60,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 17, Unit: comply.PeriodWeeks},
				},
			},
			Source: comply.Source{Title: kaazg.Title, Section: "S 4a(2)", URL: kaazg.URL},
		},
	}
}
