// Package ie registers Ireland's healthcare scheduling regulations:
// Organisation of Working Time Act 1997 (S.I. No. 20/1998) and
// European Communities (Organisation of Working Time) Regulations.
package ie

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.IE,
		Name:     "Ireland",
		Type:     comply.Country,
		Parent:   comply.EU,
		Currency: "EUR",
		TimeZone: "Europe/Dublin",
		Rules:    rules(),
	}
}

func rules() []*comply.RuleDef {
	owta := comply.Source{
		Title:   "Organisation of Working Time Act 1997",
		Section: "",
		URL:     "https://www.irishstatutebook.ie/eli/1997/act/20/enacted/en/html",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxWeeklyHours,
			Name:        "Maximum Weekly Working Hours",
			Description: "Average weekly working time shall not exceed 48 hours, calculated over a reference period of 4 months (6 months for healthcare by collective agreement).",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(1997, time.September, 30),
					Amount:   48,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 4, Unit: comply.PeriodMonths},
				},
			},
			Source: comply.Source{Title: owta.Title, Section: "S. 15", URL: owta.URL},
		},
		{
			Key:         comply.RuleMinRestBetweenShifts,
			Name:        "Minimum Daily Rest Period",
			Description: "At least 11 consecutive hours of rest per 24-hour period.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1997, time.September, 30), Amount: 11, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: owta.Title, Section: "S. 11(1)", URL: owta.URL},
		},
		{
			Key:         comply.RuleMinWeeklyRest,
			Name:        "Minimum Weekly Rest Period",
			Description: "At least 24 consecutive hours of rest per 7-day period, preceded by the 11-hour daily rest (35 hours total). Should include Sunday where practicable.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1997, time.September, 30), Amount: 35, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: owta.Title, Section: "S. 13", URL: owta.URL},
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Break Requirement",
			Description: "Entitled to a break after 4.5 hours of work. A second break after 6 hours.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1997, time.September, 30), Amount: 4.5, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{Title: owta.Title, Section: "S. 12(1)", URL: owta.URL},
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Break Duration",
			Description: "At least 15 minutes after 4.5 hours; at least 30 minutes total after 6 hours (may include the first 15 minutes).",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1997, time.September, 30), Amount: 15, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{Title: owta.Title, Section: "S. 12(2)", URL: owta.URL},
			Notes:  "30 minutes total required after 6 hours.",
		},
		{
			Key:         comply.RuleNightPeriodStart,
			Name:        "Night Period Start",
			Description: "Night time is midnight to 7 AM.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1997, time.September, 30), Amount: 0, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: owta.Title, Section: "S. 16(1)", URL: owta.URL},
		},
		{
			Key:         comply.RuleNightPeriodEnd,
			Name:        "Night Period End",
			Description: "Night work period ends at 07:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1997, time.September, 30), Amount: 7, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: owta.Title, Section: "S. 16(1)", URL: owta.URL},
		},
		{
			Key:         comply.RuleMaxNightShiftHours,
			Name:        "Maximum Night Worker Hours",
			Description: "Night workers may not work more than 8 hours on average per 24-hour period, averaged over 2 months.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(1997, time.September, 30),
					Amount:   8,
					Unit:     comply.Hours,
					Per:      comply.PerDay,
					Averaged: &comply.AveragingPeriod{Count: 2, Unit: comply.PeriodMonths},
				},
			},
			Source: comply.Source{Title: owta.Title, Section: "S. 16(2)", URL: owta.URL},
		},
		{
			Key:         comply.RuleMinAnnualLeaveDays,
			Name:        "Minimum Annual Leave",
			Description: "Minimum 20 working days (4 working weeks) of paid annual leave per year.",
			Category:    comply.CatLeave,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1997, time.September, 30), Amount: 20, Unit: comply.Days, Per: comply.PerYear},
			},
			Source: comply.Source{Title: owta.Title, Section: "S. 19(1)", URL: owta.URL},
		},
		{
			Key:         "ie-healthcare-optout",
			Name:        "Healthcare Opt-Out (Extended Working Time)",
			Description: "Ireland uses the EU Article 22 opt-out for healthcare workers. Individual written consent allows exceeding the 48-hour weekly average. Used extensively for NCHDs (Non-Consultant Hospital Doctors).",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpBool,
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2004, time.August, 1), Amount: 1, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "European Communities (Organisation of Working Time) (Activities of Doctors in Training) Regulations 2004",
				Section: "S.I. No. 494/2004, Reg. 4",
				URL:     "https://www.irishstatutebook.ie/eli/2004/si/494/made/en/print",
			},
			Notes: "NCHDs (Non-Consultant Hospital Doctors) frequently work 60+ hours/week with opt-out. Subject to ongoing reform efforts.",
		},
		{
			Key:         comply.RuleMaxShiftHours,
			Name:        "Maximum Shift Duration (NCHD)",
			Description: "Non-Consultant Hospital Doctors: shifts should not exceed 24 hours including on-call. European Working Time Directive compliance has been a persistent challenge.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffResident},
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2004, time.August, 1), Amount: 24, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{
				Title:   "HSE EWTD Implementation Plan for NCHDs",
				Section: "Schedule 3",
				URL:     "https://www.hse.ie/eng/staff/resources/hr-publications/",
			},
		},
	}
}
