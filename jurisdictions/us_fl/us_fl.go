// Package us_fl registers Florida healthcare scheduling regulations:
// nursing home shift limits (FAC 59A-4.108) and documented absences
// of mandatory overtime prohibition, daily overtime, meal/rest break,
// and day of rest requirements.
package us_fl

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

// New returns the Florida jurisdiction definition.
func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USFL,
		Name:     "Florida",
		Type:     comply.State,
		Parent:   comply.US,
		Currency: "USD",
		TimeZone: "America/New_York",
		Rules:    rules(),
	}
}

func rules() []*comply.RuleDef {
	r := make([]*comply.RuleDef, 0, 6)
	r = append(r, nursingHomeRules()...)
	r = append(r, documentedAbsences()...)
	return r
}

// Nursing Home Shift Rules - FAC 59A-4.108

func nursingHomeRules() []*comply.RuleDef {
	return []*comply.RuleDef{
		{
			Key:         "fl-nursing-home-max-shift",
			Name:        "Nursing Home Maximum Shift Hours",
			Description: "Nursing home staff may not work more than 16 hours in a 24-hour period, for a maximum of 3 consecutive days.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffNurseRN, comply.StaffNurseLPN, comply.StaffNurseCNA},
			Scope:       comply.ScopeNursingHomes,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2015, time.December, 21), Amount: 16, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{
				Title:   "Florida Administrative Code",
				Section: "59A-4.108 (Nursing Home Staffing)",
				URL:     "https://www.flrules.org/gateway/ruleNo.asp?id=59A-4.108",
			},
			Notes: "Only state-level hour limit for healthcare workers in Florida. Nursing homes only.",
		},
	}
}

// Documented Absences - Florida has minimal state-level labor protections.

func documentedAbsences() []*comply.RuleDef {
	return []*comply.RuleDef{
		{
			Key:         comply.RuleMandatoryOTProhibited,
			Name:        "No Mandatory Overtime Prohibition",
			Description: "Florida does not prohibit mandatory overtime for nurses. Florida is not among the 18 states with nurse mandatory overtime restrictions.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpBool,
			Enforcement: comply.Advisory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1970, time.January, 1), Amount: 0, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title: "Florida Statutes, Chapter 464 (Nurse Practice Act)",
				URL:   "https://www.flsenate.gov/Laws/Statutes/2023/Chapter464/All",
			},
			Notes: "Florida is NOT among the 18 states with nurse mandatory overtime restrictions.",
		},
		{
			Key:         comply.RuleOvertimeDailyThreshold,
			Name:        "No Daily Overtime Threshold",
			Description: "Florida does not impose a daily overtime threshold. Follows FLSA only.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpGTE,
			Enforcement: comply.Advisory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1970, time.January, 1), Amount: 0, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{
				Title: "Florida Statutes",
			},
			Notes: "Follows FLSA only.",
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "No Meal Break Requirement",
			Description: "Florida does not require employers to provide meal breaks for adult employees.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Advisory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1970, time.January, 1), Amount: 0, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{
				Title:   "Florida Statutes",
				Section: "450.081(4) (minors only)",
			},
		},
		{
			Key:         comply.RuleRestBreakDuration,
			Name:        "No Rest Break Requirement",
			Description: "Florida does not require employers to provide rest breaks.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Advisory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1970, time.January, 1), Amount: 0, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{
				Title: "Florida Statutes",
			},
		},
		{
			Key:         comply.RuleMinDayOfRest,
			Name:        "No Day of Rest Requirement",
			Description: "Florida does not require employers to provide a mandatory day of rest.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Advisory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1970, time.January, 1), Amount: 0, Unit: comply.Days, Per: comply.PerWeek},
			},
			Source: comply.Source{
				Title: "Florida Statutes",
			},
		},
	}
}
