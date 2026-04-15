// Package eu registers EU-level healthcare scheduling regulations:
// Working Time Directive 2003/88/EC and key CJEU case law on on-call time.
package eu

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

// New returns the EU jurisdiction definition.
func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.EU,
		Name:     "European Union",
		Type:     comply.Supranational,
		Currency: "EUR",
		TimeZone: "Europe/Brussels",
		Rules:    rules(),
	}
}

func rules() []*comply.RuleDef {
	wtdSource := comply.Source{
		Title:   "Directive 2003/88/EC (Working Time Directive)",
		Section: "",
		URL:     "https://eur-lex.europa.eu/legal-content/EN/TXT/?uri=celex:32003L0088",
	}

	return []*comply.RuleDef{
		// Article 6: Maximum Weekly Working Time
		{
			Key:         comply.RuleMaxWeeklyHours,
			Name:        "Maximum Weekly Working Time",
			Description: "Average working time including overtime shall not exceed 48 hours per 7-day period. Cannot be derogated under Article 17 - only Article 22 individual opt-out can override.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2003, time.November, 4),
					Amount:   48,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 4, Unit: comply.PeriodMonths},
					Exceptions: []string{
						"Member States may extend reference period to 6 months by law or collective agreement",
						"Reference period may extend to 12 months by collective agreement only, for objective/technical reasons",
						"Article 22 opt-out: individual written consent allows exceeding 48 hours (approx. 16 Member States use this)",
					},
				},
			},
			Source: comply.Source{
				Title:   wtdSource.Title,
				Section: "Article 6",
				URL:     wtdSource.URL,
			},
			Notes: "The 48-hour limit is absolute unless the individual opt-out (Article 22) is activated by the Member State. Approximately 16 Member States have some form of opt-out, many limited to healthcare.",
		},
		// Article 3: Daily Rest
		{
			Key:         comply.RuleMinRestBetweenShifts,
			Name:        "Minimum Daily Rest Period",
			Description: "Every worker entitled to a minimum of 11 consecutive hours of rest per 24-hour period.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2003, time.November, 4),
					Amount: 11,
					Unit:   comply.Hours,
					Per:    comply.PerDay,
					Exceptions: []string{
						"Article 17(3)(c)(i): hospitals and similar establishments may derogate if equivalent compensatory rest is granted",
						"Derogation requires compensatory rest immediately following the working period (per CJEU Jaeger C-151/02)",
					},
				},
			},
			Source: comply.Source{
				Title:   wtdSource.Title,
				Section: "Article 3",
				URL:     wtdSource.URL,
			},
		},
		// Article 5: Weekly Rest
		{
			Key:         comply.RuleMinWeeklyRest,
			Name:        "Minimum Weekly Rest Period",
			Description: "Minimum 24 hours uninterrupted rest per 7-day period, plus the 11 hours daily rest (35 hours total). May be averaged over 14 days.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2003, time.November, 4),
					Amount:   35,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 14, Unit: comply.PeriodDays},
				},
			},
			Source: comply.Source{
				Title:   wtdSource.Title,
				Section: "Article 5",
				URL:     wtdSource.URL,
			},
		},
		// Article 4: Breaks
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Break Entitlement Trigger",
			Description: "Workers entitled to a rest break when working day exceeds 6 hours. Duration and terms set by Member State law or collective agreement.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2003, time.November, 4), Amount: 6, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{
				Title:   wtdSource.Title,
				Section: "Article 4",
				URL:     wtdSource.URL,
			},
			Notes: "The Directive does not specify break duration - this is left to Member States.",
		},
		// Article 8: Night Work
		{
			Key:         comply.RuleMaxNightShiftHours,
			Name:        "Maximum Night Work Duration",
			Description: "Normal hours of night workers shall not exceed an average of 8 hours per 24-hour period.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2003, time.November, 4), Amount: 8, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{
				Title:   wtdSource.Title,
				Section: "Article 8",
				URL:     wtdSource.URL,
			},
			Notes: "Night time = any period of not less than 7 hours, defined by national law, which must include midnight to 5:00 AM. Night worker = regularly works at least 3 hours during night time.",
		},
		// Article 7: Annual Leave
		{
			Key:         comply.RuleMinAnnualLeaveDays,
			Name:        "Minimum Annual Leave",
			Description: "Every worker entitled to paid annual leave of at least 4 weeks (20 working days). May not be replaced by financial compensation except on termination.",
			Category:    comply.CatLeave,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2003, time.November, 4), Amount: 20, Unit: comply.Days, Per: comply.PerYear},
			},
			Source: comply.Source{
				Title:   wtdSource.Title,
				Section: "Article 7",
				URL:     wtdSource.URL,
			},
		},
		// Article 17(3)(c)(i): Healthcare Derogation
		{
			Key:         "eu-healthcare-derogation",
			Name:        "Healthcare Sector Derogation (Article 17)",
			Description: "Hospitals and similar establishments may derogate from daily rest (Art. 3), breaks (Art. 4), weekly rest (Art. 5), night work (Art. 8), and reference periods (Art. 16). CANNOT derogate from the 48-hour maximum (Art. 6). Equivalent compensatory rest must be provided.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpBool,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2003, time.November, 4), Amount: 1, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   wtdSource.Title,
				Section: "Article 17(3)(c)(i)",
				URL:     wtdSource.URL,
			},
			Notes: "CJEU case law: SIMAP (C-303/98, 2000) - on-call at workplace = working time in entirety. Jaeger (C-151/02, 2003) - compensatory rest must be immediate and consecutive. Matzak (C-518/15, 2018) - home standby with severe response-time constraints may be working time.",
		},
		// Article 22: Individual Opt-Out
		{
			Key:         "eu-article-22-optout",
			Name:        "Individual Opt-Out from 48-Hour Maximum",
			Description: "Member States may allow workers to individually consent in writing to exceed the 48-hour weekly limit. Worker cannot suffer detriment for refusing. May withdraw consent with up to 3 months notice.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpBool,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2003, time.November, 4), Amount: 1, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   wtdSource.Title,
				Section: "Article 22(1)",
				URL:     wtdSource.URL,
			},
			Notes: "Countries with broad opt-out: Bulgaria, Cyprus, Estonia, Malta. Healthcare-specific opt-outs: Belgium, Croatia, France, Germany, Hungary, Netherlands, Poland, Slovakia, Slovenia, Spain (via Ley 55/2003 Art. 49). Denmark added limited opt-out in 2024. Countries WITHOUT opt-out: Austria, Finland, Greece, Ireland, Italy, Lithuania, Luxembourg, Portugal, Romania, Sweden.",
		},
	}
}
