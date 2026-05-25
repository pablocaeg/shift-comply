// Package us_oh registers Ohio healthcare scheduling regulations:
// ORC S 4723.09 (nurse overtime restriction).
package us_oh

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USOH,
		Name:     "Ohio",
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
			Description: "Hospitals may not schedule or require a nurse to work in excess of an agreed-upon and predetermined regular work schedule. Includes RNs, LPNs, and nursing assistants.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpBool,
			StaffTypes:  []comply.Key{comply.StaffNurseRN, comply.StaffNurseLPN, comply.StaffNurseCNA},
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2022, time.April, 7),
					Amount: 1,
					Unit:   comply.Boolean,
					Exceptions: []string{
						"Federal, state, or local declaration of emergency",
						"Unforeseen emergency posing immediate threat to patient health and safety",
					},
				},
			},
			Source: comply.Source{
				Title:   "Ohio Revised Code",
				Section: "ORC S 4723.09 (HB 279, 134th General Assembly)",
				URL:     "https://codes.ohio.gov/ohio-revised-code/section-4723.09",
			},
			Notes: "Effective April 7, 2022. Relatively recent compared to other states. Applies only to hospitals.",
		},
	}
}
