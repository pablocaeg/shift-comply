// Package de registers Germany's healthcare scheduling regulations:
// Arbeitszeitgesetz (ArbZG), collective agreements for hospital doctors
// (TV-Aerzte/VKA), and nurse staffing minimums (PpUGV).
package de

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

// New returns the Germany jurisdiction definition.
func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:      comply.DE,
		Name:      "Germany",
		LocalName: "Deutschland",
		Type:      comply.Country,
		Parent:    comply.EU,
		Currency:  "EUR",
		TimeZone:  "Europe/Berlin",
		Rules:     rules(),
	}
}

func rules() []*comply.RuleDef {
	r := make([]*comply.RuleDef, 0, 20)
	r = append(r, arbzgRules()...)
	r = append(r, healthcareRules()...)
	return r
}

// Arbeitszeitgesetz (Working Time Act) - general labor law.
// BGBl. I S. 1170, last amended 2020.

func arbzgRules() []*comply.RuleDef {
	arbzg := comply.Source{
		Title:   "Arbeitszeitgesetz (ArbZG)",
		Section: "",
		URL:     "https://www.gesetze-im-internet.de/arbzg/",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxDailyHours,
			Name:        "Maximum Daily Working Hours",
			Description: "Working time may not exceed 8 hours per working day. May be extended to 10 hours if within a 6-month or 24-week reference period the average does not exceed 8 hours per working day.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(1994, time.July, 1),
					Amount: 8,
					Unit:   comply.Hours,
					Per:    comply.PerDay,
					Exceptions: []string{
						"May extend to 10 hours/day if averaged to 8 hours over 6 months or 24 weeks",
					},
				},
			},
			Source: comply.Source{Title: arbzg.Title, Section: "S 3", URL: arbzg.URL},
		},
		{
			Key:         comply.RuleMaxWeeklyHours,
			Name:        "Maximum Weekly Working Hours",
			Description: "Derived from the 8-hour daily limit across 6 working days (Mon-Sat): 48 hours per week. With extension: 60 hours if averaged to 48 over the reference period.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(1994, time.July, 1),
					Amount:   48,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 6, Unit: comply.PeriodMonths},
					Exceptions: []string{
						"May reach 60 hours/week if averaged to 48 over 6 months or 24 weeks",
					},
				},
			},
			Source: comply.Source{Title: arbzg.Title, Section: "S 3", URL: arbzg.URL},
			Notes:  "Germany counts Saturday as a working day, so 6 x 8 = 48 hours baseline.",
		},
		{
			Key:         comply.RuleMinRestBetweenShifts,
			Name:        "Minimum Daily Rest Period",
			Description: "Workers must have at least 11 uninterrupted hours of rest after the end of the daily working time.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(1994, time.July, 1),
					Amount: 11,
					Unit:   comply.Hours,
					Per:    comply.PerDay,
					Exceptions: []string{
						"In hospitals and care facilities: may be reduced to 10 hours if compensated by 12 hours rest within 4 weeks (S 5(2) ArbZG)",
					},
				},
			},
			Source: comply.Source{Title: arbzg.Title, Section: "S 5(1)", URL: arbzg.URL},
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Break Requirement - First Break",
			Description: "Work may not continue for more than 6 hours without a break.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1994, time.July, 1), Amount: 6, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{Title: arbzg.Title, Section: "S 4", URL: arbzg.URL},
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Break Duration",
			Description: "At least 30 minutes for shifts of 6-9 hours; at least 45 minutes for shifts over 9 hours. Breaks may be split into 15-minute blocks.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1994, time.July, 1), Amount: 30, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{Title: arbzg.Title, Section: "S 4", URL: arbzg.URL},
			Notes:  "45 minutes required if shift exceeds 9 hours.",
		},
		{
			Key:         comply.RuleNightPeriodStart,
			Name:        "Night Period Start",
			Description: "Night time is defined as 11:00 PM to 6:00 AM (23:00-06:00).",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1994, time.July, 1), Amount: 23, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: arbzg.Title, Section: "S 2(3)", URL: arbzg.URL},
		},
		{
			Key:         comply.RuleNightPeriodEnd,
			Name:        "Night Period End",
			Description: "Night time ends at 6:00 AM.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1994, time.July, 1), Amount: 6, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: arbzg.Title, Section: "S 2(3)", URL: arbzg.URL},
		},
		{
			Key:         comply.RuleMaxNightShiftHours,
			Name:        "Maximum Night Worker Daily Hours",
			Description: "Night workers may not work more than 8 hours per working day, averaged over a calendar month or 4 weeks.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(1994, time.July, 1),
					Amount:   8,
					Unit:     comply.Hours,
					Per:      comply.PerDay,
					Averaged: &comply.AveragingPeriod{Count: 4, Unit: comply.PeriodWeeks},
				},
			},
			Source: comply.Source{Title: arbzg.Title, Section: "S 6(2)", URL: arbzg.URL},
		},
		{
			Key:         comply.RuleMinAnnualLeaveDays,
			Name:        "Minimum Annual Leave",
			Description: "Minimum 24 working days of paid annual leave (based on a 6-day work week). Equivalent to 20 days for a 5-day week.",
			Category:    comply.CatLeave,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1963, time.January, 1), Amount: 24, Unit: comply.Days, Per: comply.PerYear},
			},
			Source: comply.Source{
				Title:   "Bundesurlaubsgesetz (BUrlG)",
				Section: "S 3(1)",
				URL:     "https://www.gesetze-im-internet.de/burlg/",
			},
			Notes: "24 working days based on 6-day week = 20 days for 5-day week. Most collective agreements provide 30 days.",
		},
	}
}

// Healthcare-specific rules.

func healthcareRules() []*comply.RuleDef {
	return []*comply.RuleDef{
		{
			Key:         "de-oncall-is-working-time",
			Name:        "On-Call at Workplace is Working Time",
			Description: "Bereitschaftsdienst (on-call duty requiring physical presence at the workplace) counts as working time in its entirety, per CJEU SIMAP/Jaeger rulings as implemented in German law.",
			Category:    comply.CatOnCall,
			Operator:    comply.OpBool,
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2004, time.January, 1), Amount: 1, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "Arbeitszeitgesetz (ArbZG) as interpreted per CJEU C-303/98 (SIMAP) and C-151/02 (Jaeger)",
				Section: "S 7(1) No. 1a",
				URL:     "https://www.gesetze-im-internet.de/arbzg/__7.html",
			},
			Notes: "S 7(1) No. 1a ArbZG allows collective agreements (Tarifvertrag) to extend daily working time beyond 10 hours if the extension includes on-call time (Bereitschaftsdienst). This is Germany's mechanism for accommodating 24-hour on-call in hospitals while formally complying with CJEU rulings.",
		},
		{
			Key:         "de-healthcare-optout",
			Name:        "Opt-Out for Healthcare (Collective Agreement)",
			Description: "Via collective agreements (Tarifvertrag), healthcare workers may individually consent in writing to exceed the 48-hour weekly limit when shifts include on-call time. Germany uses this for hospital doctors under TV-Aerzte/VKA and TV-Aerzte/TdL.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpBool,
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2004, time.January, 1), Amount: 1, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "Arbeitszeitgesetz (ArbZG)",
				Section: "S 7(2a) (opt-out via Tarifvertrag, implementing EU Art. 22)",
				URL:     "https://www.gesetze-im-internet.de/arbzg/__7.html",
			},
			Notes: "Individual written consent required. Worker may withdraw with 6-month notice. Average weekly hours including on-call must not exceed 58 hours under most collective agreements (TV-Aerzte).",
		},
		{
			Key:         "de-rest-reduction-hospitals",
			Name:        "Reduced Rest Period in Hospitals",
			Description: "In hospitals and care facilities, the 11-hour daily rest may be reduced to 10 hours by collective agreement if a compensatory rest of 12 hours is granted within 4 weeks.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1994, time.July, 1), Amount: 10, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{
				Title:   "Arbeitszeitgesetz (ArbZG)",
				Section: "S 5(2)",
				URL:     "https://www.gesetze-im-internet.de/arbzg/__5.html",
			},
		},
		{
			Key:         comply.RuleMaxShiftHours,
			Name:        "Maximum Shift Including On-Call (TV-Aerzte)",
			Description: "Under the TV-Aerzte collective agreement for hospital doctors, a single duty period including on-call (Bereitschaftsdienst) may not exceed 24 hours. With opt-out consent, transitional activities of up to 3 hours may follow.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffPhysician},
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2006, time.August, 1),
					Amount: 24,
					Unit:   comply.Hours,
					Per:    comply.PerShift,
				},
			},
			Source: comply.Source{
				Title:   "Tarifvertrag fuer Aerztinnen und Aerzte an kommunalen Krankenhaeusern (TV-Aerzte/VKA)",
				Section: "S 7(1)",
				URL:     "https://www.marburger-bund.de/tarifvertraege",
			},
		},
		{
			Key:         comply.RuleMinRestAfterGuard,
			Name:        "Rest After On-Call Duty (TV-Aerzte)",
			Description: "After on-call duty (Bereitschaftsdienst) of 24 hours, the worker must have the remaining day free as compensatory rest.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			StaffTypes:  []comply.Key{comply.StaffPhysician},
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2006, time.August, 1), Amount: 11, Unit: comply.Hours, Per: comply.PerOccurrence},
			},
			Source: comply.Source{
				Title:   "TV-Aerzte/VKA",
				Section: "S 7(3)",
				URL:     "https://www.marburger-bund.de/tarifvertraege",
			},
			Notes: "In practice, after a 24-hour on-call shift ending in the morning, the doctor is free for the rest of that day.",
		},
	}
}
