// Package fr registers France's healthcare scheduling regulations:
// Code du travail (Labour Code), 35-hour work week, specific provisions
// for healthcare workers in public hospitals (fonction publique hospitaliere),
// and resident (interne) work hour limits.
package fr

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.FR,
		Name:     "France",
		Type:     comply.Country,
		Parent:   comply.EU,
		Currency: "EUR",
		TimeZone: "Europe/Paris",
		Rules:    rules(),
	}
}

func rules() []*comply.RuleDef {
	r := make([]*comply.RuleDef, 0, 20)
	r = append(r, codeRules()...)
	r = append(r, healthcareRules()...)
	return r
}

func codeRules() []*comply.RuleDef {
	ct := comply.Source{
		Title:   "Code du travail",
		Section: "",
		URL:     "https://www.legifrance.gouv.fr/codes/id/LEGITEXT000006072050/",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxOrdinaryWeeklyHours,
			Name:        "Legal Weekly Working Time (35 Hours)",
			Description: "The legal duration of work is 35 hours per week. Hours beyond 35 are overtime.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2000, time.February, 1), Amount: 35, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. L3121-27", URL: ct.URL},
			Notes:  "35 hours is the threshold for overtime calculation, not an absolute maximum.",
		},
		{
			Key:         comply.RuleMaxWeeklyHours,
			Name:        "Absolute Maximum Weekly Hours",
			Description: "Working time may not exceed 48 hours in any single week, nor 44 hours averaged over 12 consecutive weeks.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2000, time.February, 1),
					Amount:   44,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 12, Unit: comply.PeriodWeeks},
					Exceptions: []string{
						"Absolute cap of 48 hours in any single week",
						"Derogation to 46 hours averaged over 12 weeks possible by decree or collective agreement",
					},
				},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. L3121-20 to L3121-22", URL: ct.URL},
		},
		{
			Key:         comply.RuleMaxDailyHours,
			Name:        "Maximum Daily Working Hours",
			Description: "Daily working time may not exceed 10 hours. May be extended to 12 hours by collective agreement for specific activities.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:      comply.D(2000, time.February, 1),
					Amount:     10,
					Unit:       comply.Hours,
					Per:        comply.PerDay,
					Exceptions: []string{"May be extended to 12 hours by collective agreement or authorization"},
				},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. L3121-18", URL: ct.URL},
		},
		{
			Key:         comply.RuleMinRestBetweenShifts,
			Name:        "Minimum Daily Rest Period",
			Description: "Every worker is entitled to at least 11 consecutive hours of daily rest.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:      comply.D(2000, time.February, 1),
					Amount:     11,
					Unit:       comply.Hours,
					Per:        comply.PerDay,
					Exceptions: []string{"May be reduced to 9 hours by collective agreement for urgent activities or continuous-service sectors"},
				},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. L3131-1", URL: ct.URL},
		},
		{
			Key:         comply.RuleMinWeeklyRest,
			Name:        "Minimum Weekly Rest Period",
			Description: "At least 35 consecutive hours of weekly rest (24 hours + 11 hours daily rest). Must include Sunday in principle.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2000, time.February, 1), Amount: 35, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. L3132-1 to L3132-2", URL: ct.URL},
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Break Duration",
			Description: "At least 20 consecutive minutes after 6 hours of continuous work.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2000, time.February, 1), Amount: 20, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. L3121-16", URL: ct.URL},
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Break Trigger",
			Description: "A break is mandatory after 6 hours of continuous work.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2000, time.February, 1), Amount: 6, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. L3121-16", URL: ct.URL},
		},
		{
			Key:         comply.RuleNightPeriodStart,
			Name:        "Night Period Start",
			Description: "Night work period runs from 21:00 to 06:00 (or 21:00-07:00 by collective agreement).",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2000, time.February, 1), Amount: 21, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. L3122-2", URL: ct.URL},
		},
		{
			Key:         comply.RuleNightPeriodEnd,
			Name:        "Night Period End",
			Description: "Night work period ends at 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2000, time.February, 1), Amount: 6, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. L3122-2", URL: ct.URL},
		},
		{
			Key:         comply.RuleMaxNightShiftHours,
			Name:        "Maximum Night Worker Daily Hours",
			Description: "Night workers may not work more than 8 hours per day. May be extended to 12 hours in healthcare by collective agreement.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2000, time.February, 1), Amount: 8, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. L3122-6", URL: ct.URL},
		},
		{
			Key:         comply.RuleMaxOvertimeAnnual,
			Name:        "Annual Overtime Quota",
			Description: "Maximum 220 hours of overtime per year per employee (default contingent). May be modified by collective agreement.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2000, time.February, 1), Amount: 220, Unit: comply.Hours, Per: comply.PerYear},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. D3121-24", URL: ct.URL},
		},
		{
			Key:         comply.RuleMinAnnualLeaveDays,
			Name:        "Minimum Annual Leave",
			Description: "2.5 working days per month of work, totaling 30 working days (5 weeks) per year.",
			Category:    comply.CatLeave,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1982, time.January, 16), Amount: 30, Unit: comply.Days, Per: comply.PerYear},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. L3141-3", URL: ct.URL},
			Notes:  "30 working days = 5 weeks. France has the most generous statutory leave in the EU.",
		},
	}
}

func healthcareRules() []*comply.RuleDef {
	return []*comply.RuleDef{
		{
			Key:         "fr-no-optout",
			Name:        "No Individual Opt-Out (Article 22 Not Adopted)",
			Description: "France does NOT use the EU Article 22 individual opt-out. The 48-hour absolute weekly maximum applies to all workers without exception.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpBool,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2000, time.February, 1), Amount: 0, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "Code du travail",
				Section: "Art. L3121-20 (no Art. 22 transposition)",
				URL:     "https://www.legifrance.gouv.fr/codes/id/LEGITEXT000006072050/",
			},
		},
		{
			Key:         comply.RuleMaxShiftHours,
			Name:        "Maximum Guard Duration (Public Hospitals)",
			Description: "Guard duty (garde) in public hospitals is limited to 24 consecutive hours. Followed by mandatory rest.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffPhysician},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2002, time.January, 1), Amount: 24, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{
				Title:   "Arrete du 14 septembre 2001 relatif a l'organisation et a l'indemnisation des gardes",
				Section: "Art. 2",
				URL:     "https://www.legifrance.gouv.fr/loda/id/JORFTEXT000000406025",
			},
		},
		{
			Key:         comply.RuleMinRestAfterGuard,
			Name:        "Rest After Guard Duty",
			Description: "After a 24-hour guard, the practitioner must have at least 11 hours of immediate rest. The post-guard rest is a demi-journee (half day) off.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			StaffTypes:  []comply.Key{comply.StaffPhysician},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2002, time.January, 1), Amount: 11, Unit: comply.Hours, Per: comply.PerOccurrence},
			},
			Source: comply.Source{
				Title:   "Code de la sante publique",
				Section: "Art. R6152-26",
				URL:     "https://www.legifrance.gouv.fr/codes/id/LEGITEXT000006072665/",
			},
		},
		{
			Key:         "fr-intern-max-weekly",
			Name:        "Intern (Interne) Maximum Weekly Hours",
			Description: "Medical interns (internes) in public hospitals may not exceed 48 hours per week averaged over a trimester (3 months), including guard duties.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffResident},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2015, time.May, 1),
					Amount:   48,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 3, Unit: comply.PeriodMonths},
				},
			},
			Source: comply.Source{
				Title:   "Arrete du 8 juillet 2010 relatif aux gardes des internes",
				Section: "Art. 3, modified 2015",
				URL:     "https://www.legifrance.gouv.fr/loda/id/JORFTEXT000022507414",
			},
			Notes: "Obligations: 10 half-days of work per week (including guards). Max 2 guards per week (Saturday 14h to Monday 8h counts as 1 guard).",
		},
		{
			Key:         comply.RuleMaxGuardsMonthly,
			Name:        "Intern Maximum Monthly Guards",
			Description: "Medical interns may not perform more than 8 guards per month.",
			Category:    comply.CatOnCall,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffResident},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2010, time.July, 8), Amount: 8, Unit: comply.Count, Per: comply.PerMonth},
			},
			Source: comply.Source{
				Title:   "Arrete du 8 juillet 2010 relatif aux gardes des internes",
				Section: "Art. 2",
				URL:     "https://www.legifrance.gouv.fr/loda/id/JORFTEXT000022507414",
			},
		},
	}
}
