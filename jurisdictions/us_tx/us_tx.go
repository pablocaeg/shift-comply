// Package us_tx registers Texas healthcare scheduling regulations.
//
// Texas is deliberately employer-friendly: no state overtime law, no meal/rest
// break requirements, no mandatory nurse-patient ratios. Key provisions are the
// nurse mandatory overtime prohibition (H&S Code Ch. 258) and nurse staffing
// committee requirement (Ch. 257), both effective September 1, 2009.
//
// Documented regulatory ABSENCES are included as advisory rules with Amount=0
// so comparison tools can show what Texas explicitly lacks vs other states.
package us_tx

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

// New returns the Texas jurisdiction definition.
func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USTX,
		Name:     "Texas",
		Type:     comply.State,
		Parent:   comply.US,
		Currency: "USD",
		TimeZone: "America/Chicago",
		Rules:    rules(),
	}
}

func rules() []*comply.RuleDef {
	return []*comply.RuleDef{
		// === Enacted regulations ===
		{
			Key:         comply.RuleMandatoryOTProhibited,
			Name:        "Mandatory Overtime Prohibition (Nurses — Hospitals Only)",
			Description: "Hospitals may not require nurses (RN or LVN) to work mandatory overtime. A nurse may refuse without penalty. On-call time cannot substitute for mandatory overtime.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpBool,
			StaffTypes:  []comply.Key{comply.StaffNurseRN, comply.StaffNurseLPN},
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2009, time.September, 1),
					Amount: 1,
					Unit:   comply.Boolean,
					Exceptions: []string{
						"Health care disaster (natural or other) in county or contiguous county",
						"Federal, state, or county declaration of emergency",
						"Unforeseen emergency that does not regularly occur and could not be prudently anticipated",
						"Nurse actively engaged in ongoing medical or surgical procedure",
					},
				},
			},
			Source: comply.Source{
				Title:   "Texas Health and Safety Code",
				Section: "Chapter 258, SS 258.001-258.005 (S.B. 476, 81st Legislature)",
				URL:     "https://statutes.capitol.texas.gov/Docs/HS/htm/HS.258.htm",
			},
			Notes: "Narrower than NY: applies ONLY to hospitals (general and special), not nursing homes or clinics. No civil penalties — enforcement is through anti-retaliation (S 258.005). Refusal cannot constitute patient abandonment per Occ. Code S 301.356.",
		},
		{
			Key:         "tx-nurse-safe-harbor",
			Name:        "Nurse Safe Harbor — Peer Review",
			Description: "A nurse asked to engage in conduct they believe violates their duty to a patient may request peer review. Requesting peer review is protected activity.",
			Category:    comply.CatCompensation,
			Operator:    comply.OpBool,
			StaffTypes:  []comply.Key{comply.StaffNurseRN, comply.StaffNurseLPN},
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2009, time.September, 1), Amount: 1, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "Texas Occupations Code",
				Section: "Chapter 303, S 303.005",
				URL:     "https://statutes.capitol.texas.gov/Docs/OC/htm/OC.303.htm",
			},
		},
		{
			Key:         "tx-nurse-staffing-committee",
			Name:        "Nurse Staffing Committee Requirement",
			Description: "Every hospital must establish a nurse staffing committee. At least 60% must be RNs providing direct patient care at least 50% of work time. Must meet at least quarterly.",
			Category:    comply.CatStaffing,
			Operator:    comply.OpBool,
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2009, time.September, 1), Amount: 1, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "Texas Health and Safety Code",
				Section: "Chapter 257 (S.B. 476, 81st Legislature)",
				URL:     "https://statutes.capitol.texas.gov/Docs/HS/htm/HS.257.htm",
			},
			Notes: "Committee develops staffing policies but does NOT mandate specific ratios. Only California has mandatory ratio laws.",
		},
		// === Documented absences ===
		// These are explicit "not enacted" markers so comparison tools can show
		// what Texas lacks versus states like CA and NY.
		{
			Key:         comply.RuleOvertimeDailyThreshold,
			Name:        "Daily Overtime Threshold — NOT ENACTED",
			Description: "Texas has no state daily overtime law. Only federal FLSA weekly overtime (40 hrs) applies.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpBool,
			Enforcement: comply.Advisory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1900, time.January, 1), Amount: 0, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "Texas Workforce Commission",
				Section: "FLSA Does and Does Not Do",
				URL:     "https://efte.twc.texas.gov/flsa_does_and_doesnt_do.html",
			},
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Meal Break Requirement — NOT ENACTED",
			Description: "Texas has no state meal break requirement for any workers, including healthcare.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpBool,
			Enforcement: comply.Advisory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1900, time.January, 1), Amount: 0, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "Texas Workforce Commission",
				Section: "Breaks",
				URL:     "https://efte.twc.texas.gov/d_breaks.html",
			},
		},
		{
			Key:         comply.RuleRestBreakDuration,
			Name:        "Rest Break Requirement — NOT ENACTED",
			Description: "Texas has no state rest break requirement.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpBool,
			Enforcement: comply.Advisory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1900, time.January, 1), Amount: 0, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "Texas Workforce Commission",
				Section: "Breaks",
				URL:     "https://efte.twc.texas.gov/d_breaks.html",
			},
		},
	}
}
