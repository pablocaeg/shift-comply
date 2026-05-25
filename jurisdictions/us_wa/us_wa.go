// Package us_wa registers Washington State healthcare scheduling regulations:
// RCW 49.28.130-49.28.150 (nurse overtime restrictions and mandatory rest),
// RCW 49.12.480 (meal/rest breaks).
package us_wa

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USWA,
		Name:     "Washington",
		Type:     comply.State,
		Parent:   comply.US,
		Currency: "USD",
		TimeZone: "America/Los_Angeles",
		Rules:    rules(),
	}
}

func rules() []*comply.RuleDef {
	return []*comply.RuleDef{
		{
			Key:         comply.RuleMandatoryOTProhibited,
			Name:        "Mandatory Overtime Prohibition (Nurses)",
			Description: "Healthcare facilities may not require nurses to work overtime. Nurses may volunteer. On-call time that is activated counts toward hours worked.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpBool,
			StaffTypes:  []comply.Key{comply.StaffNurseRN, comply.StaffNurseLPN, comply.StaffNurseCNA},
			Scope:       comply.ScopeHealthcareEmployers,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2002, time.June, 13),
					Amount: 1,
					Unit:   comply.Boolean,
					Exceptions: []string{
						"Unforeseen emergent circumstance",
						"National or state emergency",
						"Nurse in a surgical procedure until completed",
					},
				},
			},
			Source: comply.Source{
				Title:   "Revised Code of Washington",
				Section: "RCW 49.28.140",
				URL:     "https://app.leg.wa.gov/RCW/default.aspx?cite=49.28.140",
			},
		},
		{
			Key:         "wa-nurse-rest-between-shifts",
			Name:        "Mandatory Rest Between Shifts (Nurses)",
			Description: "Nurses must have at least 10 consecutive hours of uninterrupted rest between shifts. 12 hours if preceding shift was 12+ hours.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			StaffTypes:  []comply.Key{comply.StaffNurseRN, comply.StaffNurseLPN, comply.StaffNurseCNA},
			Scope:       comply.ScopeHealthcareEmployers,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2019, time.July, 28), Amount: 10, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{
				Title:   "Revised Code of Washington",
				Section: "RCW 49.28.130(2)(b)",
				URL:     "https://app.leg.wa.gov/RCW/default.aspx?cite=49.28.130",
			},
			Notes: "HB 1155 (2019). 12 hours rest required if preceding shift was 12+ hours.",
		},
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Meal Break Requirement",
			Description: "Employees must receive a 30-minute meal break when working more than 5 hours.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1961, time.January, 1), Amount: 5, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{
				Title:   "Washington Administrative Code",
				Section: "WAC 296-126-092(1)",
				URL:     "https://app.leg.wa.gov/WAC/default.aspx?cite=296-126-092",
			},
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Meal Break Duration",
			Description: "At least 30 minutes, not less than 2 hours nor more than 5 hours from beginning of shift.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1961, time.January, 1), Amount: 30, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{
				Title:   "Washington Administrative Code",
				Section: "WAC 296-126-092(1)",
				URL:     "https://app.leg.wa.gov/WAC/default.aspx?cite=296-126-092",
			},
		},
		{
			Key:         comply.RuleRestBreakDuration,
			Name:        "Rest Break Duration",
			Description: "Paid 10-minute rest break for each 4-hour work period.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1961, time.January, 1), Amount: 10, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{
				Title:   "Washington Administrative Code",
				Section: "WAC 296-126-092(4)",
				URL:     "https://app.leg.wa.gov/WAC/default.aspx?cite=296-126-092",
			},
		},
	}
}
