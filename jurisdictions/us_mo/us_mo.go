// Package us_mo registers Missouri healthcare scheduling regulations:
// RSMo S 335.300 (nurse mandatory overtime restriction in hospitals).
package us_mo

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USMO,
		Name:     "Missouri",
		Type:     comply.State,
		Parent:   comply.US,
		Currency: "USD",
		TimeZone: "America/Chicago",
		Rules:    rules(),
	}
}

func rules() []*comply.RuleDef {
	return []*comply.RuleDef{
		{
			Key:         comply.RuleMandatoryOTProhibited,
			Name:        "Mandatory Overtime Prohibition (Nurses - Hospitals)",
			Description: "Hospitals may not require a registered nurse to work in excess of the regularly scheduled hours for the nurse. Covers only RNs, not LPNs or CNAs.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpBool,
			StaffTypes:  []comply.Key{comply.StaffNurseRN, comply.StaffNurseLPN},
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2009, time.August, 28),
					Amount: 1,
					Unit:   comply.Boolean,
					Exceptions: []string{
						"Declared federal, state, or local emergency",
						"Unforeseen emergency that could not be prudently planned for",
					},
				},
			},
			Source: comply.Source{
				Title:   "Missouri Revised Statutes",
				Section: "RSMo S 335.300 (HB 904, 2009)",
				URL:     "https://revisor.mo.gov/main/OneSection.aspx?section=335.300",
			},
			Notes: "Narrower than many states: covers only RNs (not LPNs/CNAs) and only hospitals (not nursing homes or clinics).",
		},
	}
}
