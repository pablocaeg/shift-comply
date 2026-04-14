// Package es registers Spain's national healthcare scheduling regulations:
// Estatuto de los Trabajadores (RDL 2/2015), Estatuto Marco for public health
// personnel (Ley 55/2003), and MIR residency regulations (RD 1146/2006).
package es

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

// New returns the Spain national jurisdiction definition.
func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:      comply.ES,
		Name:      "Spain",
		LocalName: "Espana",
		Type:      comply.Country,
		Parent:    comply.EU,
		Currency:  "EUR",
		TimeZone:  "Europe/Madrid",
		Rules:     rules(),
	}
}

func rules() []*comply.RuleDef {
	r := make([]*comply.RuleDef, 0, 30)
	r = append(r, estatutoRules()...)
	r = append(r, estatutoMarcoRules()...)
	r = append(r, mirRules()...)
	return r
}

// Estatuto de los Trabajadores - Real Decreto Legislativo 2/2015
// Baseline labor law for private sector and supplementary for public sector.

func estatutoRules() []*comply.RuleDef {
	etSource := comply.Source{
		Title:   "Real Decreto Legislativo 2/2015 (Estatuto de los Trabajadores)",
		Section: "",
		URL:     "https://www.boe.es/buscar/act.php?id=BOE-A-2015-11430",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxOrdinaryWeeklyHours,
			Name:        "Maximum Ordinary Weekly Hours",
			Description: "Maximum 40 hours of effective work per week, calculated as an annual average. Allows irregular distribution throughout the year.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(1980, time.March, 14),
					Amount:   40,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 12, Unit: comply.PeriodMonths},
				},
			},
			Source: comply.Source{Title: etSource.Title, Section: "Articulo 34.1", URL: etSource.URL},
		},
		{
			Key:         comply.RuleMaxDailyHours,
			Name:        "Maximum Ordinary Daily Hours",
			Description: "Maximum 9 hours of ordinary work per day, unless collective agreement establishes different daily distribution (must respect rest between shifts).",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:      comply.D(1980, time.March, 14),
					Amount:     9,
					Unit:       comply.Hours,
					Per:        comply.PerDay,
					Exceptions: []string{"Collective agreement may establish different daily distribution"},
				},
			},
			Source: comply.Source{Title: etSource.Title, Section: "Articulo 34.3", URL: etSource.URL},
		},
		{
			Key:         comply.RuleMinRestBetweenShifts,
			Name:        "Minimum Rest Between Shifts",
			Description: "Minimum 12 hours between the end of one working day and the beginning of the next.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1980, time.March, 14), Amount: 12, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: etSource.Title, Section: "Articulo 34.3", URL: etSource.URL},
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Break After 6 Continuous Hours",
			Description: "When continuous daily working hours exceed 6, a rest break of at least 15 minutes must be provided. Workers under 18: 30 minutes after 4.5 continuous hours.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1980, time.March, 14), Amount: 6, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{Title: etSource.Title, Section: "Articulo 34.4", URL: etSource.URL},
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Minimum Break Duration",
			Description: "Break must be at least 15 minutes.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1980, time.March, 14), Amount: 15, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{Title: etSource.Title, Section: "Articulo 34.4", URL: etSource.URL},
		},
		{
			Key:         comply.RuleMaxOvertimeAnnual,
			Name:        "Maximum Annual Overtime",
			Description: "Maximum 80 overtime hours per year. Hours compensated with equivalent rest within 4 months do not count toward limit. Hours for emergency/damage prevention also excluded.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(1980, time.March, 14),
					Amount: 80,
					Unit:   comply.Hours,
					Per:    comply.PerYear,
					Exceptions: []string{
						"Overtime compensated with equivalent rest within 4 months is excluded",
						"Overtime for preventing or repairing extraordinary/urgent damage is excluded",
					},
				},
			},
			Source: comply.Source{Title: etSource.Title, Section: "Articulo 35.2", URL: etSource.URL},
		},
		{
			Key:         comply.RuleMinWeeklyRest,
			Name:        "Minimum Weekly Rest",
			Description: "One and a half consecutive days of weekly rest (generally Saturday afternoon/Monday morning plus Sunday). Accumulable over 14 days. Workers under 18: 2 consecutive days.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(1980, time.March, 14),
					Amount:   36,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 14, Unit: comply.PeriodDays},
				},
			},
			Source: comply.Source{Title: etSource.Title, Section: "Articulo 37.1", URL: etSource.URL},
			Notes:  "1.5 days = 36 hours (combining 24h weekly rest + 12h daily rest).",
		},
		// --- Night Work ---
		{
			Key:         comply.RuleNightPeriodStart,
			Name:        "Night Period Start",
			Description: "Night period runs from 10:00 PM to 6:00 AM.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1980, time.March, 14), Amount: 22, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: etSource.Title, Section: "Articulo 36.1", URL: etSource.URL},
		},
		{
			Key:         comply.RuleNightPeriodEnd,
			Name:        "Night Period End",
			Description: "Night period runs from 10:00 PM to 6:00 AM.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1980, time.March, 14), Amount: 6, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: etSource.Title, Section: "Articulo 36.1", URL: etSource.URL},
		},
		{
			Key:         comply.RuleMaxNightShiftHours,
			Name:        "Maximum Night Worker Daily Hours",
			Description: "Night workers' daily working time cannot exceed 8 hours on average over a 15-day reference period. Night workers cannot perform overtime hours (horas extraordinarias).",
			Category:    comply.CatNightWork,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(1980, time.March, 14),
					Amount:   8,
					Unit:     comply.Hours,
					Per:      comply.PerDay,
					Averaged: &comply.AveragingPeriod{Count: 15, Unit: comply.PeriodDays},
				},
			},
			Source: comply.Source{Title: etSource.Title, Section: "Articulo 36.1", URL: etSource.URL},
		},
		{
			Key:         comply.RuleMaxNightConsecWeeks,
			Name:        "Maximum Consecutive Night Shift Weeks",
			Description: "In continuous 24-hour operations, no worker may work the night shift for more than 2 consecutive weeks, except by voluntary assignment.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1980, time.March, 14), Amount: 2, Unit: comply.Weeks, Per: comply.PerPeriod},
			},
			Source: comply.Source{Title: etSource.Title, Section: "Articulo 36.3", URL: etSource.URL},
		},
		// --- Annual Leave ---
		{
			Key:         comply.RuleMinAnnualLeaveDays,
			Name:        "Minimum Annual Leave",
			Description: "Minimum 30 calendar days of paid annual leave per year. Not reducible by collective agreement.",
			Category:    comply.CatLeave,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1980, time.March, 14), Amount: 30, Unit: comply.CalendarDays, Per: comply.PerYear},
			},
			Source: comply.Source{Title: etSource.Title, Section: "Articulo 38.1", URL: etSource.URL},
		},
	}
}

// Estatuto Marco - Ley 55/2003 (Public Health System Personnel)

func estatutoMarcoRules() []*comply.RuleDef {
	emSource := comply.Source{
		Title:   "Ley 55/2003, del Estatuto Marco del personal estatutario de los servicios de salud",
		Section: "",
		URL:     "https://www.boe.es/buscar/act.php?id=BOE-A-2003-23101",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxCombinedWeeklyHours,
			Name:        "Maximum Combined Weekly Hours (Ordinary + Complementary)",
			Description: "Combined ordinary and complementary (guard duty) working time shall not exceed 48 hours per week of effective work, averaged over a semester (6 months). On-call by contactability not counted unless called in.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffStatutory},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2003, time.December, 17),
					Amount:   48,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 6, Unit: comply.PeriodMonths},
				},
			},
			Source: comply.Source{Title: emSource.Title, Section: "Articulo 48", URL: emSource.URL},
			Notes:  "Complementary hours (jornada complementaria) apply when continuous care services require permanent attention. Only applies to staff categories that were already covering guard duty before the law's entry into force.",
		},
		{
			Key:         "es-max-ordinary-shift",
			Name:        "Maximum Ordinary Shift Duration",
			Description: "Ordinary working time shall not exceed 12 uninterrupted hours.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffStatutory},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2003, time.December, 17), Amount: 12, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{Title: emSource.Title, Section: "Articulo 51", URL: emSource.URL},
		},
		{
			Key:         comply.RuleMaxShiftHours,
			Name:        "Maximum Exceptional Shift Duration (24-Hour Guards)",
			Description: "Working days of up to 24 hours may be established for certain units/services when organizational or care reasons warrant it.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffStatutory},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:      comply.D(2003, time.December, 17),
					Amount:     24,
					Unit:       comply.Hours,
					Per:        comply.PerShift,
					Exceptions: []string{"Only through functional programming for specific units/services with organizational or care justification"},
				},
			},
			Source: comply.Source{Title: emSource.Title, Section: "Articulo 51", URL: emSource.URL},
		},
		{
			Key:         "es-statutory-min-weekly-rest",
			Name:        "Minimum Weekly Rest (Public Health Personnel)",
			Description: "Average 24 hours uninterrupted weekly rest. Combined with 12-hour daily rest yields 36 hours practical minimum. Reference period: 14 days.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			StaffTypes:  []comply.Key{comply.StaffStatutory},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2003, time.December, 17),
					Amount:   36,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 14, Unit: comply.PeriodDays},
				},
			},
			Source: comply.Source{Title: emSource.Title, Section: "Articulo 52", URL: emSource.URL},
			Notes:  "36 hours cannot be monetized or substituted - must be taken as rest. Quarterly minimum: 96 hours of weekly rest including rest actually taken.",
		},
		{
			Key:         "es-special-regime-optout",
			Name:        "Special Working Time Regime (48-Hour Opt-Out)",
			Description: "Workers may individually, freely, and in writing consent to exceed the 48-hour combined limit. Worker cannot suffer detriment for refusing. Spain's healthcare-specific implementation of EU Article 22 opt-out.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpBool,
			StaffTypes:  []comply.Key{comply.StaffStatutory},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2003, time.December, 17), Amount: 1, Unit: comply.Boolean},
			},
			Source: comply.Source{Title: emSource.Title, Section: "Articulo 49", URL: emSource.URL},
		},
	}
}

// MIR Residency Regulations - Real Decreto 1146/2006

func mirRules() []*comply.RuleDef {
	mirSource := comply.Source{
		Title:   "Real Decreto 1146/2006 (relacion laboral especial de residencia)",
		Section: "",
		URL:     "https://www.boe.es/buscar/act.php?id=BOE-A-2006-17498",
	}

	return []*comply.RuleDef{
		{
			Key:         "es-mir-max-ordinary-weekly",
			Name:        "MIR Maximum Ordinary Weekly Hours",
			Description: "Maximum ordinary working hours for medical residents: 37.5 hours per week, averaged over a semester.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffResident},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2006, time.October, 7),
					Amount:   37.5,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 6, Unit: comply.PeriodMonths},
				},
			},
			Source: comply.Source{Title: mirSource.Title, Section: "Articulo 5.1.a", URL: mirSource.URL},
		},
		{
			Key:         comply.RuleMaxGuardsMonthly,
			Name:        "MIR Maximum Monthly Guards",
			Description: "Residents may not perform more than 7 guard duties (guardias) per month.",
			Category:    comply.CatOnCall,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffResident},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2006, time.October, 7), Amount: 7, Unit: comply.Count, Per: comply.PerMonth},
			},
			Source: comply.Source{Title: mirSource.Title, Section: "Articulo 5.1.c", URL: mirSource.URL},
			Notes:  "Residents are obligated to perform only the complementary hours established by their training program for the corresponding year.",
		},
		{
			Key:         comply.RuleMinRestAfterGuard,
			Name:        "MIR Minimum Rest After 24-Hour Guard",
			Description: "After 24 hours of uninterrupted work, the resident is entitled to 12 continuous hours of rest.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			StaffTypes:  []comply.Key{comply.StaffResident},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2006, time.October, 7),
					Amount: 12,
					Unit:   comply.Hours,
					Per:    comply.PerOccurrence,
					Exceptions: []string{
						"Emergency care situations (emergencia asistencial)",
					},
				},
			},
			Source: comply.Source{Title: mirSource.Title, Section: "Articulo 5.1.b", URL: mirSource.URL},
			Notes:  "Tribunal Supremo ruling (April 21, 2022, Amyts v. SERMAS) established 36 hours uninterrupted rest after Saturday/pre-holiday 24-hour guards, applying the Workers' Statute supplementarily.",
		},
		{
			Key:         "es-mir-min-rest-between-shifts",
			Name:        "MIR Minimum Rest Between Shifts",
			Description: "Minimum 12 continuous hours between the end of one working day and the beginning of the next.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			StaffTypes:  []comply.Key{comply.StaffResident},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2006, time.October, 7), Amount: 12, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{Title: mirSource.Title, Section: "Articulo 5.1.b", URL: mirSource.URL},
		},
	}
}
