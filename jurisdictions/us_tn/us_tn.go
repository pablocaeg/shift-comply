// Package us_tn registers Tennessee healthcare scheduling regulations:
// TCA 50-2-103(h) (30-minute meal break for 6+ hour shifts).
package us_tn

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USTN,
		Name:     "Tennessee",
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
			Key:         comply.RuleMealBreakDuration,
			Name:        "Meal Break",
			Description: "30-minute unpaid meal or rest period for employees scheduled to work 6 consecutive hours. Cannot be scheduled during or before the first hour of work.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1999, time.July, 1), Amount: 30, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{
				Title:   "Tennessee Code Annotated",
				Section: "TCA S 50-2-103(h)",
				URL:     "https://www.tn.gov/workforce/employees/labor-laws/labor-laws-redirect/meal-and-rest-breaks.html",
			},
			Notes: "Penalty: Class B misdemeanor, fine $100-$500 per violation. Exemption: workplaces that by nature provide ample opportunity for breaks.",
		},
		{
			Key:         comply.RuleMandatoryOTProhibited,
			Name:        "Mandatory Overtime Prohibition -- NOT ENACTED",
			Description: "Tennessee does not prohibit mandatory overtime for nurses or other healthcare workers.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpBool,
			Enforcement: comply.Advisory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1900, time.January, 1), Amount: 0, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   "Tennessee Department of Labor and Workforce Development",
				Section: "No healthcare-specific scheduling statute",
				URL:     "https://www.tn.gov/workforce/",
			},
		},
	}
}
