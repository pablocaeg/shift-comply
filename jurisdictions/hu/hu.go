// Package hu registers Hungary's healthcare scheduling regulations:
// Labour Code (Munka Torvenykonyve, Act I of 2012), and healthcare-specific
// on-call and opt-out provisions from Act CLIV of 1997 on Health (Eutv.).
//
// Hungary makes extensive use of the EU Article 22 opt-out for healthcare
// workers, allowing combined regular + on-call hours well beyond 48/week.
package hu

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

// New returns the Hungary jurisdiction definition.
func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:      comply.HU,
		Name:      "Hungary",
		LocalName: "Magyarorszag",
		Type:      comply.Country,
		Parent:    comply.EU,
		Currency:  "HUF",
		TimeZone:  "Europe/Budapest",
		Rules:     rules(),
	}
}

func rules() []*comply.RuleDef {
	r := make([]*comply.RuleDef, 0, 15)
	r = append(r, labourCodeRules()...)
	r = append(r, healthcareRules()...)
	return r
}

// Munka Torvenykonyve (Labour Code) - Act I of 2012.

func labourCodeRules() []*comply.RuleDef {
	mt := comply.Source{
		Title:   "Munka Torvenykonyve (Labour Code, Act I of 2012)",
		Section: "",
		URL:     "https://njt.hu/jogszabaly/2012-1-00-00",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxDailyHours,
			Name:        "Maximum Daily Working Hours",
			Description: "Standard daily working time is 8 hours. May be extended to 12 hours in special schedules.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2012, time.July, 1), Amount: 8, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: mt.Title, Section: "S 92", URL: mt.URL},
		},
		{
			Key:         comply.RuleMaxWeeklyHours,
			Name:        "Maximum Weekly Working Hours",
			Description: "Weekly working time including overtime shall not exceed 48 hours, averaged over a reference period. Standard reference period is 4 months, extendable to 6 months by collective agreement.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2012, time.July, 1),
					Amount:   48,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 4, Unit: comply.PeriodMonths},
				},
			},
			Source: comply.Source{Title: mt.Title, Section: "S 99", URL: mt.URL},
		},
		{
			Key:         comply.RuleMinRestBetweenShifts,
			Name:        "Minimum Daily Rest Period",
			Description: "At least 11 consecutive hours of rest between working days.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2012, time.July, 1),
					Amount: 11,
					Unit:   comply.Hours,
					Per:    comply.PerDay,
					Exceptions: []string{
						"May be reduced to 8 hours for healthcare workers performing on-call duty",
					},
				},
			},
			Source: comply.Source{Title: mt.Title, Section: "S 104", URL: mt.URL},
		},
		{
			Key:         comply.RuleMinWeeklyRest,
			Name:        "Minimum Weekly Rest Period",
			Description: "At least 48 consecutive hours of weekly rest, or 2 x 40 hours in each 14-day period.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2012, time.July, 1),
					Amount:   48,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 14, Unit: comply.PeriodDays},
				},
			},
			Source: comply.Source{Title: mt.Title, Section: "S 105", URL: mt.URL},
		},
		{
			Key:         comply.RuleMaxOvertimeAnnual,
			Name:        "Maximum Annual Overtime",
			Description: "Maximum 250 hours of overtime per year. May be extended to 300 hours by collective agreement.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:      comply.D(2012, time.July, 1),
					Amount:     250,
					Unit:       comply.Hours,
					Per:        comply.PerYear,
					Exceptions: []string{"May be extended to 300 hours by collective agreement"},
				},
			},
			Source: comply.Source{Title: mt.Title, Section: "S 109", URL: mt.URL},
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Break Requirement",
			Description: "At least 20 minutes break after 6 hours of continuous work.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2012, time.July, 1), Amount: 6, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{Title: mt.Title, Section: "S 103", URL: mt.URL},
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Break Duration",
			Description: "At least 20 minutes after 6 hours of work. Additional 25 minutes after 9 hours.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2012, time.July, 1), Amount: 20, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{Title: mt.Title, Section: "S 103", URL: mt.URL},
		},
	}
}

// Healthcare-specific rules - Act CLIV of 1997 on Health (Eutv.).

func healthcareRules() []*comply.RuleDef {
	eutv := comply.Source{
		Title:   "Act CLIV of 1997 on Health (Egeszsegugyi torveny, Eutv.)",
		Section: "",
		URL:     "https://njt.hu/jogszabaly/1997-154-00-00",
	}

	return []*comply.RuleDef{
		{
			Key:         "hu-oncall-is-working-time",
			Name:        "On-Call Duty (Ugyelet) is Working Time",
			Description: "On-call duty requiring physical presence at the healthcare facility (ugyelet) counts as working time in its entirety, in compliance with CJEU SIMAP/Jaeger rulings.",
			Category:    comply.CatOnCall,
			Operator:    comply.OpBool,
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2004, time.May, 1), Amount: 1, Unit: comply.Boolean},
			},
			Source: comply.Source{Title: eutv.Title, Section: "S 95", URL: eutv.URL},
		},
		{
			Key:         "hu-healthcare-optout",
			Name:        "Healthcare Opt-Out (Extended Working Time)",
			Description: "Healthcare workers may individually consent in writing to work beyond 48 hours/week when on-call duty is included. This is Hungary's implementation of the EU Article 22 opt-out for the healthcare sector.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpBool,
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2004, time.May, 1), Amount: 1, Unit: comply.Boolean},
			},
			Source: comply.Source{Title: eutv.Title, Section: "S 96", URL: eutv.URL},
			Notes:  "Hungary is one of the most extensive users of the opt-out in the EU healthcare sector. Without opt-out, the 24-hour on-call model used in most Hungarian hospitals would be impossible under the 48-hour WTD limit.",
		},
		{
			Key:         comply.RuleMaxShiftHours,
			Name:        "Maximum On-Call Duty Duration",
			Description: "A single on-call duty period (ugyelet) may not exceed 24 hours. Combined with regular working time immediately preceding it, the total continuous work period is limited.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2004, time.May, 1), Amount: 24, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{Title: eutv.Title, Section: "S 95(2)", URL: eutv.URL},
		},
		{
			Key:         comply.RuleMinRestAfterGuard,
			Name:        "Rest After On-Call Duty",
			Description: "After on-call duty (ugyelet), the healthcare worker must receive the daily rest period (11 hours) before the next working period.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2004, time.May, 1), Amount: 11, Unit: comply.Hours, Per: comply.PerOccurrence},
			},
			Source: comply.Source{Title: eutv.Title, Section: "S 97", URL: eutv.URL},
		},
	}
}
