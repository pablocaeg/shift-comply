// Package ro registers Romania's healthcare scheduling regulations:
// Codul Muncii (Labour Code, Law 53/2003) as amended.
package ro

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:      comply.RO,
		Name:      "Romania",
		LocalName: "Romania",
		Type:      comply.Country,
		Parent:    comply.EU,
		Currency:  "RON",
		TimeZone:  "Europe/Bucharest",
		Rules:     rules(),
	}
}

func rules() []*comply.RuleDef {
	cm := comply.Source{
		Title:   "Codul Muncii (Legea nr. 53/2003, republicata)",
		Section: "",
		URL:     "https://legislatie.just.ro/Public/DetaliiDocument/41625",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxOrdinaryWeeklyHours,
			Name:        "Maximum Ordinary Weekly Hours",
			Description: "Normal working time is 8 hours per day and 40 hours per week.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2003, time.March, 1), Amount: 40, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: cm.Title, Section: "Art. 112", URL: cm.URL},
		},
		{
			Key:         comply.RuleMaxWeeklyHours,
			Name:        "Maximum Weekly Hours Including Overtime",
			Description: "Total working time including overtime shall not exceed 48 hours per week averaged over a 4-month reference period.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2003, time.March, 1),
					Amount:   48,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 4, Unit: comply.PeriodMonths},
				},
			},
			Source: comply.Source{Title: cm.Title, Section: "Art. 114", URL: cm.URL},
		},
		{
			Key:         comply.RuleMaxDailyHours,
			Name:        "Maximum Daily Working Hours",
			Description: "Daily working time may not exceed 12 hours per day followed by at least 24 hours rest.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2003, time.March, 1), Amount: 12, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: cm.Title, Section: "Art. 115", URL: cm.URL},
		},
		{
			Key:         comply.RuleMinRestBetweenShifts,
			Name:        "Minimum Daily Rest Period",
			Description: "At least 12 consecutive hours of rest between two working days. In shifts: at least 8 hours between shifts.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:      comply.D(2003, time.March, 1),
					Amount:     12,
					Unit:       comply.Hours,
					Per:        comply.PerDay,
					Exceptions: []string{"May be reduced to 8 hours for shift work"},
				},
			},
			Source: comply.Source{Title: cm.Title, Section: "Art. 135", URL: cm.URL},
			Notes:  "Romania requires 12 hours rest (above EU 11-hour minimum) but allows 8 hours for shift workers.",
		},
		{
			Key:         comply.RuleMinWeeklyRest,
			Name:        "Minimum Weekly Rest Period",
			Description: "At least 48 consecutive hours of weekly rest, normally Saturday and Sunday.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2003, time.March, 1), Amount: 48, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: cm.Title, Section: "Art. 137", URL: cm.URL},
			Notes:  "48 hours weekly rest is above the EU 35-hour minimum. One of the most generous in Europe alongside Hungary.",
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Break Requirement",
			Description: "Workers are entitled to a break when daily working time exceeds 6 hours.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2003, time.March, 1), Amount: 6, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{Title: cm.Title, Section: "Art. 134", URL: cm.URL},
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Break Duration",
			Description: "At least 15 minutes. Specific duration set by internal rules or collective agreement (typically 30 minutes).",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2003, time.March, 1), Amount: 15, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{Title: cm.Title, Section: "Art. 134", URL: cm.URL},
		},
		{
			Key:         comply.RuleMaxOvertimeAnnual,
			Name:        "Maximum Annual Overtime",
			Description: "Maximum overtime is 8 hours per week. Annualized: approximately 360 hours (but must stay within the 48-hour weekly average).",
			Category:    comply.CatOvertime,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2003, time.March, 1), Amount: 360, Unit: comply.Hours, Per: comply.PerYear},
			},
			Source: comply.Source{Title: cm.Title, Section: "Art. 114, Art. 120", URL: cm.URL},
			Notes:  "No explicit annual cap; the 48h/week average over 4 months acts as the constraint. 360h is the theoretical max (8h/week x 45 weeks).",
		},
		{
			Key:         comply.RuleNightPeriodStart,
			Name:        "Night Period Start",
			Description: "Night work is defined as work performed between 22:00 and 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2003, time.March, 1), Amount: 22, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: cm.Title, Section: "Art. 125", URL: cm.URL},
		},
		{
			Key:         comply.RuleNightPeriodEnd,
			Name:        "Night Period End",
			Description: "Night work period ends at 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2003, time.March, 1), Amount: 6, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: cm.Title, Section: "Art. 125", URL: cm.URL},
		},
		{
			Key:         comply.RuleMaxNightShiftHours,
			Name:        "Maximum Night Worker Hours",
			Description: "Night workers may not work more than 8 hours per 24-hour period.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2003, time.March, 1), Amount: 8, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: cm.Title, Section: "Art. 125(2)", URL: cm.URL},
		},
		{
			Key:         comply.RuleMinAnnualLeaveDays,
			Name:        "Minimum Annual Leave",
			Description: "Minimum 20 working days of paid annual leave per year.",
			Category:    comply.CatLeave,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2003, time.March, 1), Amount: 20, Unit: comply.Days, Per: comply.PerYear},
			},
			Source: comply.Source{Title: cm.Title, Section: "Art. 145", URL: cm.URL},
		},
	}
}
