// Package si registers Slovenia's healthcare scheduling regulations:
// Zakon o delovnih razmerjih (Employment Relationships Act, ZDR-1,
// Official Gazette RS 21/2013).
package si

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:      comply.SI,
		Name:      "Slovenia",
		LocalName: "Slovenija",
		Type:      comply.Country,
		Parent:    comply.EU,
		Currency:  "EUR",
		TimeZone:  "Europe/Ljubljana",
		Rules:     rules(),
	}
}

func rules() []*comply.RuleDef {
	zdr := comply.Source{
		Title:   "Zakon o delovnih razmerjih (ZDR-1, Ur.l. RS 21/2013)",
		Section: "",
		URL:     "http://www.pisrs.si/Pis.web/pregledPredpisa?id=ZAKO5944",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxOrdinaryWeeklyHours,
			Name:        "Maximum Ordinary Weekly Hours",
			Description: "Full-time working time is 40 hours per week. Minimum 36 hours by law or collective agreement.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2013, time.April, 12), Amount: 40, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: zdr.Title, Section: "Art. 143", URL: zdr.URL},
		},
		{
			Key:         comply.RuleMaxWeeklyHours,
			Name:        "Maximum Weekly Hours Including Overtime",
			Description: "Total working time including overtime shall not exceed 48 hours per week averaged over a 6-month reference period.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2013, time.April, 12),
					Amount:   48,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 6, Unit: comply.PeriodMonths},
				},
			},
			Source: comply.Source{Title: zdr.Title, Section: "Art. 148", URL: zdr.URL},
		},
		{
			Key:         comply.RuleMaxShiftHours,
			Name:        "Maximum Daily Working Hours",
			Description: "Working time may not exceed 10 hours per day including overtime.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2013, time.April, 12), Amount: 10, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{Title: zdr.Title, Section: "Art. 146(3)", URL: zdr.URL},
		},
		{
			Key:         comply.RuleMinRestBetweenShifts,
			Name:        "Minimum Daily Rest Period",
			Description: "At least 12 consecutive hours of rest per 24-hour period.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2013, time.April, 12), Amount: 12, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: zdr.Title, Section: "Art. 155", URL: zdr.URL},
			Notes:  "Slovenia requires 12 hours daily rest, above the EU 11-hour minimum.",
		},
		{
			Key:         comply.RuleMinWeeklyRest,
			Name:        "Minimum Weekly Rest Period",
			Description: "At least 24 consecutive hours of weekly rest plus the daily rest (36 hours total).",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2013, time.April, 12), Amount: 36, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: zdr.Title, Section: "Art. 156", URL: zdr.URL},
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Break Duration",
			Description: "At least 30 minutes during a full working day. Paid break included in working time.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2013, time.April, 12), Amount: 30, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{Title: zdr.Title, Section: "Art. 154", URL: zdr.URL},
		},
		{
			Key:         comply.RuleMaxOvertimeAnnual,
			Name:        "Maximum Annual Overtime",
			Description: "Maximum 170 hours of overtime per year. Up to 230 hours with employee written consent.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2013, time.April, 12), Amount: 170, Unit: comply.Hours, Per: comply.PerYear},
			},
			Source: comply.Source{Title: zdr.Title, Section: "Art. 146(2)", URL: zdr.URL},
		},
		{
			Key:         comply.RuleNightPeriodStart,
			Name:        "Night Period Start",
			Description: "Night work is defined as work between 23:00 and 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2013, time.April, 12), Amount: 23, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: zdr.Title, Section: "Art. 149", URL: zdr.URL},
		},
		{
			Key:         comply.RuleNightPeriodEnd,
			Name:        "Night Period End",
			Description: "Night work period ends at 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2013, time.April, 12), Amount: 6, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: zdr.Title, Section: "Art. 149", URL: zdr.URL},
		},
		{
			Key:         comply.RuleMinAnnualLeaveDays,
			Name:        "Minimum Annual Leave",
			Description: "Minimum 20 working days (4 weeks) of paid annual leave. Additional days based on age, seniority, disability, and children.",
			Category:    comply.CatLeave,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2013, time.April, 12), Amount: 20, Unit: comply.Days, Per: comply.PerYear},
			},
			Source: comply.Source{Title: zdr.Title, Section: "Art. 159", URL: zdr.URL},
		},
	}
}
