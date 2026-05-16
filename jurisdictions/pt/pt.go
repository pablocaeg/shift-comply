// Package pt registers Portugal's healthcare scheduling regulations:
// Codigo do Trabalho (Labour Code, Lei 7/2009), specific healthcare
// provisions for NHS (SNS) workers, and resident (interno) rules.
package pt

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.PT,
		Name:     "Portugal",
		Type:     comply.Country,
		Parent:   comply.EU,
		Currency: "EUR",
		TimeZone: "Europe/Lisbon",
		Rules:    rules(),
	}
}

func rules() []*comply.RuleDef {
	ct := comply.Source{
		Title:   "Codigo do Trabalho (Lei n. 7/2009)",
		Section: "",
		URL:     "https://diariodarepublica.pt/dr/legislacao-consolidada/lei/2009-34546475",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxOrdinaryWeeklyHours,
			Name:        "Maximum Ordinary Weekly Hours",
			Description: "Normal weekly working time is 40 hours for private sector. Public sector (SNS) is 35 hours since 2016.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2009, time.February, 12), Amount: 40, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. 203", URL: ct.URL},
		},
		{
			Key:         comply.RuleMaxWeeklyHours,
			Name:        "Maximum Weekly Hours Including Overtime",
			Description: "Total weekly working time including overtime shall not exceed 48 hours, averaged over a reference period of 4 months (extendable to 6 by collective agreement).",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2009, time.February, 12),
					Amount:   48,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 4, Unit: comply.PeriodMonths},
				},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. 211", URL: ct.URL},
		},
		{
			Key:         comply.RuleMaxDailyHours,
			Name:        "Maximum Daily Working Hours",
			Description: "Normal working time may not exceed 8 hours per day.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2009, time.February, 12), Amount: 8, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. 203(1)", URL: ct.URL},
		},
		{
			Key:         comply.RuleMinRestBetweenShifts,
			Name:        "Minimum Daily Rest Period",
			Description: "At least 11 consecutive hours of rest between two consecutive working periods.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2009, time.February, 12), Amount: 11, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. 214(1)", URL: ct.URL},
		},
		{
			Key:         comply.RuleMinWeeklyRest,
			Name:        "Minimum Weekly Rest Period",
			Description: "At least one day (24 consecutive hours) of weekly rest, plus the 11-hour daily rest (35 hours total).",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2009, time.February, 12), Amount: 35, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. 232", URL: ct.URL},
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Break Requirement",
			Description: "Workers may not work more than 5 consecutive hours without a break.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2009, time.February, 12), Amount: 5, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. 213(1)", URL: ct.URL},
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Break Duration",
			Description: "Break must be at least 1 hour and no more than 2 hours, unless otherwise agreed.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2009, time.February, 12), Amount: 60, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. 213(1)", URL: ct.URL},
		},
		{
			Key:         comply.RuleMaxOvertimeAnnual,
			Name:        "Maximum Annual Overtime",
			Description: "Maximum 150 hours of overtime per year for medium and large companies. 175 hours for micro/small companies.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2009, time.February, 12), Amount: 150, Unit: comply.Hours, Per: comply.PerYear},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. 228", URL: ct.URL},
		},
		{
			Key:         comply.RuleNightPeriodStart,
			Name:        "Night Period Start",
			Description: "Night work period: 22:00 to 07:00 (or defined by collective agreement within the 00:00-05:00 window).",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2009, time.February, 12), Amount: 22, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. 223(1)", URL: ct.URL},
		},
		{
			Key:         comply.RuleNightPeriodEnd,
			Name:        "Night Period End",
			Description: "Night work period ends at 07:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2009, time.February, 12), Amount: 7, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. 223(1)", URL: ct.URL},
		},
		{
			Key:         comply.RuleMaxNightShiftHours,
			Name:        "Maximum Night Worker Hours",
			Description: "Night workers may not work more than 8 hours per day on average over a reference period.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2009, time.February, 12), Amount: 8, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. 224(1)", URL: ct.URL},
		},
		{
			Key:         comply.RuleMinAnnualLeaveDays,
			Name:        "Minimum Annual Leave",
			Description: "Minimum 22 working days of paid annual leave per year.",
			Category:    comply.CatLeave,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2009, time.February, 12), Amount: 22, Unit: comply.Days, Per: comply.PerYear},
			},
			Source: comply.Source{Title: ct.Title, Section: "Art. 238(1)", URL: ct.URL},
		},
		{
			Key:         "pt-sns-weekly-hours",
			Name:        "SNS Public Health Worker Weekly Hours",
			Description: "Public health workers in the SNS (Servico Nacional de Saude) have a 35-hour work week since 2016.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2016, time.July, 1), Amount: 35, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{
				Title:   "Lei n. 18/2016 (reposicao das 35 horas)",
				Section: "Art. 2",
				URL:     "https://diariodarepublica.pt/dr/detalhe/lei/18-2016-74739812",
			},
		},
		{
			Key:         comply.RuleMaxShiftHours,
			Name:        "Maximum Guard Shift Duration (SNS)",
			Description: "Guard shifts in SNS hospitals are limited to 24 consecutive hours maximum.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffPhysician},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2009, time.February, 12), Amount: 24, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{
				Title:   "Decreto-Lei n. 73/90 (carreiras medicas)",
				Section: "Art. 31",
				URL:     "https://diariodarepublica.pt/dr/detalhe/decreto-lei/73-1990-332831",
			},
		},
	}
}
