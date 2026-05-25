// Package us_ak registers Alaska healthcare scheduling regulations:
// AS 23.10.060 (daily overtime after 8 hours), AS 18.20.400-499 (nurse
// mandatory overtime ban with 10h rest requirement).
package us_ak

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USAK,
		Name:     "Alaska",
		Type:     comply.State,
		Parent:   comply.US,
		Currency: "USD",
		TimeZone: "America/Anchorage",
		Rules:    rules(),
	}
}

func rules() []*comply.RuleDef {
	return []*comply.RuleDef{
		{
			Key:         comply.RuleOvertimeDailyThreshold,
			Name:        "Daily Overtime Threshold",
			Description: "Overtime pay (1.5x) required for hours worked in excess of 8 in a workday.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1959, time.January, 3), Amount: 8, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{
				Title:   "Alaska Statutes",
				Section: "AS 23.10.060",
				URL:     "https://www.akleg.gov/basis/statutes.asp#23.10.060",
			},
		},
		{
			Key:         comply.RuleMandatoryOTProhibited,
			Name:        "Mandatory Overtime Prohibition (Nurses and CNAs)",
			Description: "Healthcare facilities may not require or coerce nurses or CNAs to work beyond a predetermined, regularly scheduled shift. After completing a scheduled shift, a nurse must receive at least 10 consecutive hours of off-duty time.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpBool,
			StaffTypes:  []comply.Key{comply.StaffNurseRN, comply.StaffNurseLPN, comply.StaffNurseCNA},
			Scope:       comply.ScopeHealthcareEmployers,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2018, time.January, 1),
					Amount: 1,
					Unit:   comply.Boolean,
					Exceptions: []string{
						"Unforeseen emergency jeopardizing patient safety",
						"Nurse actively engaged in a surgical or medical procedure",
						"Voluntary overtime",
					},
				},
			},
			Source: comply.Source{
				Title:   "Alaska Statutes",
				Section: "AS 18.20.400 to 18.20.499",
				URL:     "https://www.akleg.gov/basis/statutes.asp#18.20.400",
			},
			Notes: "Also requires 10 consecutive hours off-duty after completing a scheduled shift.",
		},
		{
			Key:         "no-meal-break-requirement",
			Name:        "Meal/Rest Break Requirement -- NOT ENACTED",
			Description: "Alaska does not require meal or rest breaks for adult workers.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpBool,
			Enforcement: comply.Advisory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1900, time.January, 1), Amount: 0, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "Alaska Department of Labor",
				Section: "No state break requirement",
				URL:     "https://labor.alaska.gov/",
			},
		},
	}
}
