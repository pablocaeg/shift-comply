// Package it registers Italy's healthcare scheduling regulations:
// D.Lgs. 66/2003 (Working Time Directive implementation), Legge 161/2014
// (which finally enforced WTD rest periods for healthcare), and CCNL
// Sanita collective agreement provisions.
//
// Italy does NOT use the Article 22 individual opt-out. Instead, it relies
// on Article 17 derogations for healthcare, which permit deviations from
// daily rest and break rules but NOT from the 48-hour weekly maximum.
package it

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

// New returns the Italy jurisdiction definition.
func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:      comply.IT,
		Name:      "Italy",
		LocalName: "Italia",
		Type:      comply.Country,
		Parent:    comply.EU,
		Currency:  "EUR",
		TimeZone:  "Europe/Rome",
		Rules:     rules(),
	}
}

func rules() []*comply.RuleDef {
	r := make([]*comply.RuleDef, 0, 15)
	r = append(r, dlgs66Rules()...)
	r = append(r, healthcareRules()...)
	return r
}

// D.Lgs. 66/2003 - Implementation of EU Working Time Directive.

func dlgs66Rules() []*comply.RuleDef {
	dlgs := comply.Source{
		Title:   "Decreto Legislativo 8 aprile 2003, n. 66",
		Section: "",
		URL:     "https://www.normattiva.it/uri-res/N2Ls?urn:nir:stato:decreto.legislativo:2003-04-08;66",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxWeeklyHours,
			Name:        "Maximum Weekly Working Hours",
			Description: "Average weekly working time including overtime shall not exceed 48 hours, averaged over a 4-month reference period. May be extended to 6 or 12 months by collective agreement.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2003, time.April, 14),
					Amount:   48,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 4, Unit: comply.PeriodMonths},
				},
			},
			Source: comply.Source{Title: dlgs.Title, Section: "Art. 4, comma 2", URL: dlgs.URL},
		},
		{
			Key:         comply.RuleMinRestBetweenShifts,
			Name:        "Minimum Daily Rest Period",
			Description: "Every worker is entitled to 11 consecutive hours of rest per 24-hour period.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2003, time.April, 14),
					Amount: 11,
					Unit:   comply.Hours,
					Per:    comply.PerDay,
				},
			},
			Source: comply.Source{Title: dlgs.Title, Section: "Art. 7", URL: dlgs.URL},
			Notes:  "Until Legge 161/2014, healthcare was effectively exempt from this requirement. Since November 25, 2015, this applies to all healthcare workers without exception.",
		},
		{
			Key:         comply.RuleMinWeeklyRest,
			Name:        "Minimum Weekly Rest Period",
			Description: "Every worker is entitled to at least 24 consecutive hours of rest per 7-day period, plus the 11-hour daily rest (35 hours total). Averaged over a 14-day reference period.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2003, time.April, 14),
					Amount:   35,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 14, Unit: comply.PeriodDays},
				},
			},
			Source: comply.Source{Title: dlgs.Title, Section: "Art. 9", URL: dlgs.URL},
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Break Requirement",
			Description: "Workers are entitled to a break when daily working time exceeds 6 hours. Duration and terms set by collective agreement, minimum 10 minutes.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2003, time.April, 14), Amount: 6, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{Title: dlgs.Title, Section: "Art. 8", URL: dlgs.URL},
		},
		{
			Key:         comply.RuleMaxOvertimeAnnual,
			Name:        "Maximum Annual Overtime",
			Description: "Maximum 250 hours of overtime per year, unless collective agreements set a different limit.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2003, time.April, 14), Amount: 250, Unit: comply.Hours, Per: comply.PerYear},
			},
			Source: comply.Source{Title: dlgs.Title, Section: "Art. 5, comma 3", URL: dlgs.URL},
		},
		{
			Key:         comply.RuleMaxNightShiftHours,
			Name:        "Maximum Night Worker Hours",
			Description: "Night workers may not work more than 8 hours on average per 24-hour period.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2003, time.April, 14), Amount: 8, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: dlgs.Title, Section: "Art. 13", URL: dlgs.URL},
		},
		{
			Key:         comply.RuleMinAnnualLeaveDays,
			Name:        "Minimum Annual Leave",
			Description: "Every worker is entitled to at least 4 weeks (20 working days) of paid annual leave. Cannot be replaced by financial compensation except on termination.",
			Category:    comply.CatLeave,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2003, time.April, 14), Amount: 20, Unit: comply.Days, Per: comply.PerYear},
			},
			Source: comply.Source{Title: dlgs.Title, Section: "Art. 10", URL: dlgs.URL},
		},
	}
}

// Healthcare-specific rules.

func healthcareRules() []*comply.RuleDef {
	return []*comply.RuleDef{
		{
			Key:         "it-no-optout",
			Name:        "No Individual Opt-Out (Article 22 Not Adopted)",
			Description: "Italy does NOT use the EU Article 22 individual opt-out. The 48-hour weekly maximum applies absolutely to all workers including healthcare. This is a key difference from Germany, Hungary, Poland, and Spain.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpBool,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2003, time.April, 14), Amount: 0, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "Decreto Legislativo 66/2003",
				Section: "Art. 17 (derogation only, no Art. 22 transposition)",
				URL:     "https://www.normattiva.it/uri-res/N2Ls?urn:nir:stato:decreto.legislativo:2003-04-08;66",
			},
			Notes: "Italy relies solely on Article 17 derogations for healthcare (rest periods, breaks), not Article 22 opt-out (weekly maximum). This has created significant staffing challenges in Italian hospitals.",
		},
		{
			Key:         "it-legge161-enforcement",
			Name:        "Enforcement of WTD Rest Periods in Healthcare (Legge 161/2014)",
			Description: "Legge 161/2014 repealed the healthcare exemption from daily and weekly rest requirements that had existed since 2008. Since November 25, 2015, all healthcare workers including doctors performing guard duty (guardia) must receive the full 11-hour daily rest.",
			Category:    comply.CatRest,
			Operator:    comply.OpBool,
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2015, time.November, 25), Amount: 1, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "Legge 30 ottobre 2014, n. 161 (Legge europea 2013-bis)",
				Section: "Art. 14",
				URL:     "https://www.normattiva.it/uri-res/N2Ls?urn:nir:stato:legge:2014-10-30;161",
			},
			Notes: "Before this law, D.L. 112/2008 Art. 41(13) had exempted healthcare from the rest requirements of D.Lgs. 66/2003. The 'smonto guardia' (post-guard rest) became legally mandatory only from November 25, 2015. Compliance remains inconsistent across Italian hospitals.",
		},
		{
			Key:         comply.RuleMinRestAfterGuard,
			Name:        "Rest After Guard Duty (Smonto Guardia)",
			Description: "After a guard shift (guardia) in a hospital, the doctor must receive the full 11-hour daily rest before the next working period. This became mandatory after Legge 161/2014.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			StaffTypes:  []comply.Key{comply.StaffPhysician},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2015, time.November, 25), Amount: 11, Unit: comply.Hours, Per: comply.PerOccurrence},
			},
			Source: comply.Source{
				Title:   "D.Lgs. 66/2003 Art. 7, as enforced by Legge 161/2014 Art. 14",
				Section: "Art. 7",
				URL:     "https://www.normattiva.it/uri-res/N2Ls?urn:nir:stato:decreto.legislativo:2003-04-08;66",
			},
		},
		{
			Key:         comply.RuleMaxShiftHours,
			Name:        "Maximum Guard Shift Duration (CCNL Sanita)",
			Description: "Under the CCNL Comparto Sanita, guard shifts (guardia) for SSN doctors may not exceed 12 hours for regular shifts. Extended 24-hour guard duty is permitted under specific collective agreement provisions.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffPhysician},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2018, time.May, 21),
					Amount: 12,
					Unit:   comply.Hours,
					Per:    comply.PerShift,
					Exceptions: []string{
						"24-hour guard shifts permitted under specific collective agreement provisions for continuous care services",
					},
				},
			},
			Source: comply.Source{
				Title:   "CCNL Comparto Sanita 2016-2018",
				Section: "Art. 27 (Orario di lavoro dei dirigenti)",
				URL:     "",
			},
		},
	}
}
