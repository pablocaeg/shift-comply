// Package us_ri registers Rhode Island healthcare scheduling regulations:
// RIGL S 23-17.17 (Safe Staffing and Mandatory Overtime Act),
// RIGL S 28-3-14 (meal breaks).
package us_ri

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USRI,
		Name:     "Rhode Island",
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
			Description: "Healthcare facilities may not require nurses to work mandatory overtime beyond their regularly scheduled shift. Covers hospitals, nursing facilities, and home health agencies.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpBool,
			StaffTypes:  []comply.Key{comply.StaffNurseRN, comply.StaffNurseLPN, comply.StaffNurseCNA},
			Scope:       comply.ScopeHealthcareEmployers,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2008, time.January, 1),
					Amount: 1,
					Unit:   comply.Boolean,
					Exceptions: []string{
						"Declared federal, state, or local emergency",
						"Unforeseen emergency circumstance",
					},
				},
			},
			Source: comply.Source{
				Title:   "Rhode Island General Laws",
				Section: "RIGL S 23-17.17",
				URL:     "http://webserver.rilin.state.ri.us/Statutes/TITLE23/23-17.17/INDEX.htm",
			},
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Meal Break Duration",
			Description: "30-minute meal break within 6 hours of starting work for shifts of 6+ hours. Unpaid if employee is completely relieved of duties.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1985, time.January, 1), Amount: 30, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{
				Title:   "Rhode Island General Laws",
				Section: "RIGL S 28-3-14",
				URL:     "http://webserver.rilin.state.ri.us/Statutes/TITLE28/28-3/28-3-14.htm",
			},
		},
	}
}
