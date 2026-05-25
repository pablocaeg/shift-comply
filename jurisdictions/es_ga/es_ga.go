// Package es_ga registers Galicia (SERGAS) healthcare scheduling regulations.
// SERGAS (Servizo Galego de Saude) has its own collective agreement with
// specific provisions on guard limits and post-guard rest.
package es_ga

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:      comply.ESGA,
		Name:      "Galicia",
		LocalName: "Galicia",
		Type:      comply.Region,
		Parent:    comply.ES,
		Currency:  "EUR",
		TimeZone:  "Europe/Madrid",
		Rules:     rules(),
	}
}

func rules() []*comply.RuleDef {
	sergas := comply.Source{
		Title:   "SERGAS Pacto sobre condiciones de trabajo",
		Section: "",
		URL:     "https://www.sergas.es/",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxGuardsMonthly,
			Name:        "Maximum Guards Per Month (SERGAS)",
			Description: "Maximum of 6 guard duties per month for SERGAS statutory personnel.",
			Category:    comply.CatOnCall,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffStatutory},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2020, time.January, 1), Amount: 6, Unit: comply.Count, Per: comply.PerMonth},
			},
			Source: comply.Source{Title: sergas.Title, Section: "Atencion continuada", URL: sergas.URL},
		},
		{
			Key:         comply.RuleMaxAnnualHours,
			Name:        "SERGAS Annual Ordinary Hours",
			Description: "Annual ordinary working hours for SERGAS statutory personnel.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffStatutory},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2020, time.January, 1), Amount: 1642.5, Unit: comply.Hours, Per: comply.PerYear},
			},
			Source: comply.Source{Title: sergas.Title, Section: "Jornada ordinaria anual", URL: sergas.URL},
		},
	}
}
