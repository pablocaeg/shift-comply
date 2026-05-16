// Package dk registers Denmark's healthcare scheduling regulations.
// Denmark is unique: working time is regulated primarily through collective
// agreements (overenskomster) rather than statute. The Working Environment Act
// (Arbejdsmiljoeloven) provides the safety framework, while EU WTD is
// implemented via Bekendtgoerelse nr. 324/2020.
package dk

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:      comply.DK,
		Name:      "Denmark",
		LocalName: "Danmark",
		Type:      comply.Country,
		Parent:    comply.EU,
		Currency:  "DKK",
		TimeZone:  "Europe/Copenhagen",
		Rules:     rules(),
	}
}

func rules() []*comply.RuleDef {
	wtd := comply.Source{
		Title:   "Bekendtgoerelse om hvileperiode og fridoegn (BEK nr. 324/2020)",
		Section: "",
		URL:     "https://www.retsinformation.dk/eli/lta/2020/324",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxWeeklyHours,
			Name:        "Maximum Weekly Working Hours",
			Description: "Average weekly working time shall not exceed 48 hours over a 4-month reference period. Denmark implements this via BEK 324/2020.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2002, time.August, 1),
					Amount:   48,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 4, Unit: comply.PeriodMonths},
				},
			},
			Source: comply.Source{Title: wtd.Title, Section: "S 4", URL: wtd.URL},
			Notes:  "The standard work week (37 hours) is set by collective agreement, not statute.",
		},
		{
			Key:         comply.RuleMinRestBetweenShifts,
			Name:        "Minimum Daily Rest Period (11 Hours)",
			Description: "Within each 24-hour period, at least 11 consecutive hours of rest.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:      comply.D(2002, time.August, 1),
					Amount:     11,
					Unit:       comply.Hours,
					Per:        comply.PerDay,
					Exceptions: []string{"May be reduced to 8 hours by collective agreement for shift work and healthcare, max twice per 7 days"},
				},
			},
			Source: comply.Source{Title: wtd.Title, Section: "S 2", URL: wtd.URL},
		},
		{
			Key:         comply.RuleDaysOffPerWeek,
			Name:        "Weekly Day Off (Fridoegn)",
			Description: "Within each 7-day period, at least one uninterrupted rest period of 24 hours (fridoegn) plus the 11 hours daily rest = 35 hours.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2002, time.August, 1), Amount: 1, Unit: comply.Days, Per: comply.PerWeek},
			},
			Source: comply.Source{Title: wtd.Title, Section: "S 3", URL: wtd.URL},
		},
		{
			Key:         comply.RuleNightPeriodStart,
			Name:        "Night Period Start",
			Description: "Night time is defined as 23:00 to 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2002, time.August, 1), Amount: 23, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: wtd.Title, Section: "S 5", URL: wtd.URL},
		},
		{
			Key:         comply.RuleNightPeriodEnd,
			Name:        "Night Period End",
			Description: "Night time ends at 06:00.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpEQ,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2002, time.August, 1), Amount: 6, Unit: comply.HourOfDay},
			},
			Source: comply.Source{Title: wtd.Title, Section: "S 5", URL: wtd.URL},
		},
		{
			Key:         comply.RuleMaxNightShiftHours,
			Name:        "Maximum Night Worker Hours",
			Description: "Night workers may not work more than 8 hours on average per 24-hour period.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2002, time.August, 1), Amount: 8, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{Title: wtd.Title, Section: "S 5", URL: wtd.URL},
		},
		{
			Key:         comply.RuleMinAnnualLeaveDays,
			Name:        "Minimum Annual Leave",
			Description: "25 days (5 weeks) of paid annual leave per year.",
			Category:    comply.CatLeave,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2000, time.September, 1), Amount: 25, Unit: comply.Days, Per: comply.PerYear},
			},
			Source: comply.Source{
				Title:   "Ferieloven (Holiday Act, LOV nr. 60/2018)",
				Section: "S 6",
				URL:     "https://www.retsinformation.dk/eli/lta/2018/60",
			},
		},
		{
			Key:         "dk-collective-agreement-model",
			Name:        "Collective Agreement Working Time Model",
			Description: "Denmark regulates working time primarily through collective agreements (overenskomster), not statute. The standard 37-hour week, overtime rules, and break provisions are all set by agreement between employer organizations and unions.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpBool,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1900, time.January, 1), Amount: 1, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "Danish Industrial Relations Model",
				Section: "Hovedaftalen (Basic Agreement)",
				URL:     "https://www.da.dk/",
			},
			Notes: "No statutory minimum break duration, overtime cap, or ordinary weekly hours. All set by overenskomst. This makes Denmark unique in Europe.",
		},
	}
}
