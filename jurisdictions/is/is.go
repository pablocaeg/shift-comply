// Package is_ registers Iceland's healthcare scheduling regulations:
// Vinnuverndarloegin (Act on Working Environment, Health and Safety
// No. 46/1980) and Loegin um 40 stunda vinnuviku (Act on 40-Hour
// Working Week). Iceland applies the WTD via the EEA Agreement.
package is

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:      comply.IS,
		Name:      "Iceland",
		LocalName: "Island",
		Type:      comply.Country,
		Parent:    comply.EU,
		Currency:  "ISK",
		TimeZone:  "Atlantic/Reykjavik",
		Rules:     rules(),
	}
}

func rules() []*comply.RuleDef {
	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxOrdinaryWeeklyHours,
			Name:        "Maximum Ordinary Weekly Hours",
			Description: "Standard working time is 40 hours per week. Iceland moved to a 36-hour week for many public sector workers in 2021 after successful trials.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1971, time.January, 1), Amount: 40, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{
				Title:   "Loegin um 40 stunda vinnuviku (Act No. 88/1971)",
				Section: "Art. 1",
				URL:     "https://www.althingi.is/lagas/nuna/1971088.html",
			},
			Notes: "Many public sector workers moved to 36-hour week in 2021 after the 2015-2019 reduced-hours trials.",
		},
		{
			Key:         comply.RuleMaxWeeklyHours,
			Name:        "Maximum Weekly Hours Including Overtime",
			Description: "Average weekly working time shall not exceed 48 hours over a 4-month reference period, per EEA WTD implementation.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2002, time.January, 1),
					Amount:   48,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 4, Unit: comply.PeriodMonths},
				},
			},
			Source: comply.Source{
				Title:   "Reglugerdin um vinnuverndarstarfsemi a vinnustoedum (Regulation No. 1000/2004)",
				Section: "Art. 52",
				URL:     "https://www.reglugerd.is/reglugerdir/eftir-raduneytum/felagsmalaraduneyti/nr/20773",
			},
		},
		{
			Key:         comply.RuleMinRestBetweenShifts,
			Name:        "Minimum Daily Rest Period",
			Description: "At least 11 consecutive hours of rest per 24-hour period.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2002, time.January, 1), Amount: 11, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{
				Title:   "Reglugerdin um vinnuverndarstarfsemi a vinnustoedum (Regulation No. 1000/2004)",
				Section: "Art. 53",
				URL:     "https://www.reglugerd.is/reglugerdir/eftir-raduneytum/felagsmalaraduneyti/nr/20773",
			},
		},
		{
			Key:         comply.RuleDaysOffPerWeek,
			Name:        "Weekly Day Off",
			Description: "At least one full day of rest per 7-day period, normally Sunday.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2002, time.January, 1), Amount: 1, Unit: comply.Days, Per: comply.PerWeek},
			},
			Source: comply.Source{
				Title:   "Reglugerdin um vinnuverndarstarfsemi a vinnustoedum (Regulation No. 1000/2004)",
				Section: "Art. 54",
				URL:     "https://www.reglugerd.is/reglugerdir/eftir-raduneytum/felagsmalaraduneyti/nr/20773",
			},
		},
		{
			Key:         comply.RuleMinAnnualLeaveDays,
			Name:        "Minimum Annual Leave",
			Description: "Minimum 24 working days of paid annual leave per year. Increases to 25 days after 5 years and 30 days after 10 years with same employer.",
			Category:    comply.CatLeave,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1987, time.June, 30), Amount: 24, Unit: comply.Days, Per: comply.PerYear},
			},
			Source: comply.Source{
				Title:   "Orlofsloegin (Holiday Act No. 30/1987)",
				Section: "Art. 3",
				URL:     "https://www.althingi.is/lagas/nuna/1987030.html",
			},
		},
	}
}
