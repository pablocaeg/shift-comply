// Package es_pv registers Basque Country (Osakidetza) healthcare scheduling
// regulations. The Basque health service has its own collective agreement
// with specific guard limits and rest provisions.
package es_pv

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:      comply.ESPV,
		Name:      "Basque Country",
		LocalName: "Euskadi",
		Type:      comply.Region,
		Parent:    comply.ES,
		Currency:  "EUR",
		TimeZone:  "Europe/Madrid",
		Rules:     rules(),
	}
}

func rules() []*comply.RuleDef {
	osakidetza := comply.Source{
		Title:   "Acuerdo regulador de condiciones de trabajo de Osakidetza",
		Section: "",
		URL:     "https://www.osakidetza.euskadi.eus/",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxGuardsMonthly,
			Name:        "Maximum Guards Per Month (Osakidetza)",
			Description: "Maximum of 5 guard duties (guardias) per month for Osakidetza statutory personnel.",
			Category:    comply.CatOnCall,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffStatutory},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2019, time.January, 1), Amount: 5, Unit: comply.Count, Per: comply.PerMonth},
			},
			Source: comply.Source{Title: osakidetza.Title, Section: "Atencion continuada", URL: osakidetza.URL},
			Notes:  "Stricter than the national 7 guards/month (RD 1146/2006) but more permissive than Catalonia's 4.",
		},
		{
			Key:         comply.RuleMinRestAfterGuard,
			Name:        "Post-Guard Rest (Osakidetza)",
			Description: "After a 24-hour guard, the worker receives the following day as rest. Post-guard rest is counted as effective working time.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			StaffTypes:  []comply.Key{comply.StaffStatutory},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2019, time.January, 1), Amount: 24, Unit: comply.Hours, Per: comply.PerOccurrence},
			},
			Source: comply.Source{Title: osakidetza.Title, Section: "Descanso post-guardia", URL: osakidetza.URL},
		},
		{
			Key:         comply.RuleMaxAnnualHours,
			Name:        "Osakidetza Annual Ordinary Hours",
			Description: "Annual ordinary working hours for Osakidetza statutory personnel.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffStatutory},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2019, time.January, 1), Amount: 1592, Unit: comply.Hours, Per: comply.PerYear},
			},
			Source: comply.Source{Title: osakidetza.Title, Section: "Jornada ordinaria anual", URL: osakidetza.URL},
			Notes:  "Lower than the national 1,642.5 hours. One of the shortest annual schedules in Spain.",
		},
	}
}
