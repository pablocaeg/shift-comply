// Package us_pa registers Pennsylvania healthcare scheduling regulations:
// 43 P.S. S 1602-A to 1612-A (Prohibition of Excessive Overtime in
// Healthcare Act, Act 102 of 2008).
package us_pa

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USPA,
		Name:     "Pennsylvania",
		Type:     comply.State,
		Parent:   comply.US,
		Currency: "USD",
		TimeZone: "America/New_York",
		Rules:    rules(),
	}
}

func rules() []*comply.RuleDef {
	actSource := comply.Source{
		Title:   "Pennsylvania Consolidated Statutes",
		Section: "43 P.S. S 1602-A to 1612-A (Act 102 of 2008)",
		URL:     "https://www.legis.state.pa.us/cfdocs/legis/LI/uconsCheck.cfm?txtType=HTM&yr=2008&sessInd=0&act=102",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMandatoryOTProhibited,
			Name:        "Prohibition of Excessive Overtime (Nurses)",
			Description: "Healthcare facilities may not require an employee to work in excess of an agreed-to, predetermined, and regularly scheduled daily work shift.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpBool,
			StaffTypes:  []comply.Key{comply.StaffNurseRN, comply.StaffNurseLPN, comply.StaffNurseCNA},
			Scope:       comply.ScopeHealthcareEmployers,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2009, time.July, 1),
					Amount: 1,
					Unit:   comply.Boolean,
					Exceptions: []string{
						"Declared national, state, or municipal emergency",
						"Federal or state agency determination of an emergency",
						"Employee engaged in an ongoing medical or surgical procedure",
					},
				},
			},
			Source: actSource,
		},
		{
			Key:         "pa-nurse-no-retaliation",
			Name:        "Anti-Retaliation (Overtime Refusal)",
			Description: "Healthcare facilities may not discriminate or retaliate against employees who refuse to work mandatory overtime.",
			Category:    comply.CatCompensation,
			Operator:    comply.OpBool,
			StaffTypes:  []comply.Key{comply.StaffNurseRN, comply.StaffNurseLPN, comply.StaffNurseCNA},
			Scope:       comply.ScopeHealthcareEmployers,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2009, time.July, 1), Amount: 1, Unit: comply.Boolean},
			},
			Source: actSource,
		},
	}
}
