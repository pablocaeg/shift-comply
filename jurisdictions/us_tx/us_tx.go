// Package us_tx registers Texas healthcare scheduling regulations:
// nurse mandatory overtime prohibition for hospitals (Health & Safety Code
// Ch. 258), nurse safe harbor protections (Occ. Code Ch. 303), staffing
// committee requirements (Health & Safety Code Ch. 257), and documented
// absences of daily overtime, meal break, and rest break requirements.
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
	r := make([]*comply.RuleDef, 0, 6)
	r = append(r, nurseRules()...)
	r = append(r, documentedAbsences()...)
	return r
}

// Nurse Rules - TX Health & Safety Code Ch. 258, Occ. Code Ch. 303

func nurseRules() []*comply.RuleDef {
	return []*comply.RuleDef{
		{
			Key:         comply.RuleMandatoryOTProhibited,
			Name:        "Nurse Mandatory Overtime Prohibition (Hospitals)",
			Description: "Hospitals may not require nurses to work mandatory overtime. Refusal to work mandatory overtime does not constitute patient abandonment per Occ. Code 301.356.",
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
						"health care disaster",
						"government emergency",
						"unforeseen emergency",
						"ongoing procedure",
					},
				},
			},
			Source: comply.Source{
				Title:   "Texas Health and Safety Code",
				Section: "Chapter 258, Sections 258.001-258.005 (S.B. 476)",
				URL:     "https://statutes.capitol.texas.gov/Docs/HS/htm/HS.258.htm",
			},
			Notes: "Hospitals only, not clinics. No civil penalties specified. Refusal not abandonment per Occ. Code 301.356.",
		},
		{
			Key:         "tx-nurse-safe-harbor",
			Name:        "Nurse Safe Harbor Protection",
			Description: "Nurses may invoke safe harbor when asked to accept an assignment that the nurse believes would expose a patient to risk of harm. Protects the nurse from retaliation.",
			Category:    comply.CatCompensation,
			Operator:    comply.OpBool,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2009, time.September, 1), Amount: 1, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "Texas Occupations Code",
				Section: "Chapter 303, Section 303.005",
				URL:     "https://statutes.capitol.texas.gov/Docs/OC/htm/OC.303.htm",
			},
		},
		{
			Key:         "tx-nurse-staffing-committee",
			Name:        "Nurse Staffing Committee Requirement",
			Description: "Hospitals must establish nurse staffing committees with at least 60% nurse composition. Committees must meet quarterly. No mandatory nurse-patient ratios are imposed.",
			Category:    comply.CatStaffing,
			Operator:    comply.OpBool,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2009, time.September, 1), Amount: 1, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "Texas Health and Safety Code",
				Section: "Chapter 257 (S.B. 476)",
				URL:     "https://statutes.capitol.texas.gov/Docs/HS/htm/HS.257.htm",
			},
			Notes: "60% nurse composition, quarterly meetings. No mandatory ratios.",
		},
	}
}

// Documented Absences - Texas has minimal state-level labor protections.

func documentedAbsences() []*comply.RuleDef {
	return []*comply.RuleDef{
		{
			Key:         comply.RuleOvertimeDailyThreshold,
			Name:        "No Daily Overtime Threshold",
			Description: "Texas does not impose a daily overtime threshold. Only FLSA weekly overtime (40 hours) applies.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpGTE,
			Enforcement: comply.Advisory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1970, time.January, 1), Amount: 0, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{
				Title: "Texas Workforce Commission",
				URL:   "https://efte.twc.texas.gov/flsa_does_and_doesnt_do.html",
			},
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "No Meal Break Requirement",
			Description: "Texas does not require employers to provide meal breaks. Only FLSA rules apply.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Advisory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1970, time.January, 1), Amount: 0, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{
				Title:   "Texas Workforce Commission",
				Section: "Breaks",
				URL:     "https://efte.twc.texas.gov/d_breaks.html",
			},
		},
		{
			Key:         comply.RuleRestBreakDuration,
			Name:        "No Rest Break Requirement",
			Description: "Texas does not require employers to provide rest breaks. Only FLSA rules apply.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Advisory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1970, time.January, 1), Amount: 0, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{
				Title:   "Texas Workforce Commission",
				Section: "Breaks",
				URL:     "https://efte.twc.texas.gov/d_breaks.html",
			},
		},
	}
}
