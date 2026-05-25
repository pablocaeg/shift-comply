// Package us_nj registers New Jersey healthcare scheduling regulations:
// NJSA 34:11-56a46-56a52 (mandatory overtime restriction for nurses, P.L. 2002 c.83),
// safe patient handling.
package us_nj

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USNJ,
		Name:     "New Jersey",
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
			Name:        "Mandatory Overtime Prohibition (Nurses)",
			Description: "Healthcare facilities may not require nurses to work more than their regularly scheduled hours. Refusal cannot be grounds for discrimination, dismissal, or retaliation.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpBool,
			StaffTypes:  []comply.Key{comply.StaffNurseRN, comply.StaffNurseLPN, comply.StaffNurseCNA},
			Scope:       comply.ScopeHealthcareEmployers,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2002, time.December, 3),
					Amount: 1,
					Unit:   comply.Boolean,
					Exceptions: []string{
						"Declared national, state, or municipal emergency",
						"Nurse engaged in an ongoing medical or surgical procedure",
					},
				},
			},
			Source: comply.Source{
				Title:   "New Jersey Statutes Annotated",
				Section: "NJSA 34:11-56a46 to 56a52 (P.L. 2002, c.83)",
				URL:     "https://lis.njleg.state.nj.us/nxt/gateway.dll?f=templates&fn=default.htm",
			},
			Notes: "One of the earliest state nurse OT bans (2002). Covers hospitals, nursing homes, and other healthcare facilities.",
		},
	}
}
