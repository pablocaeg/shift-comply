// Package us_md registers Maryland healthcare scheduling regulations:
// Health-General Art. S 19-311.1 (mandatory overtime restriction for nurses).
package us_md

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USMD,
		Name:     "Maryland",
		Type:     comply.State,
		Parent:   comply.US,
		Currency: "USD",
		TimeZone: "America/New_York",
		Rules:    rules(),
	}
}

func rules() []*comply.RuleDef {
	return []*comply.RuleDef{
		{
			Key:         comply.RuleMandatoryOTProhibited,
			Name:        "Mandatory Overtime Restriction (Nurses)",
			Description: "A hospital or related institution may not require a nurse to work more than the predetermined work shift. Refusal to work mandatory overtime is not patient abandonment.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpBool,
			StaffTypes:  []comply.Key{comply.StaffNurseRN, comply.StaffNurseLPN},
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2007, time.October, 1),
					Amount: 1,
					Unit:   comply.Boolean,
					Exceptions: []string{
						"Declaration of emergency by federal, state, or local government",
						"Unforeseen emergency that could not be prudently planned for",
						"Nurse actively engaged in an ongoing medical or surgical procedure",
					},
				},
			},
			Source: comply.Source{
				Title:   "Maryland Code, Health-General",
				Section: "S 19-311.1 (SB 555, 2007)",
				URL:     "https://mgaleg.maryland.gov/mgawebsite/Laws/StatuteText?article=ghg&section=19-311.1",
			},
		},
	}
}
