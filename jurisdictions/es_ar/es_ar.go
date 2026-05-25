// Package es_ar registers Aragon healthcare scheduling regulations.
// Servicio Aragones de Salud (SALUD) is the regional health service.
package es_ar

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:      comply.ESAR,
		Name:      "Aragon",
		LocalName: "Aragon",
		Type:      comply.Region,
		Parent:    comply.ES,
		Currency:  "EUR",
		TimeZone:  "Europe/Madrid",
		Rules:     rules(),
	}
}

func rules() []*comply.RuleDef {
	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxAnnualHours,
			Name:        "Annual Ordinary Hours (Servicio Aragones de Salud)",
			Description: "Annual ordinary working hours for Servicio Aragones de Salud statutory personnel.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffStatutory},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2020, time.January, 1), Amount: 1642.5, Unit: comply.Hours, Per: comply.PerYear},
			},
			Source: comply.Source{
				Title:   "Servicio Aragones de Salud (SALUD)",
				Section: "Jornada ordinaria anual",
				URL:     "https://www.saludinforma.es/",
			},
		},
		{
			Key:         comply.RuleMinWeeklyRest,
			Name:        "Minimum Weekly Rest (Servicio Aragones de Salud)",
			Description: "36 hours of uninterrupted weekly rest (24h weekly + 12h daily), as established by national law and reinforced by Tribunal Supremo STS 280/2022.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			StaffTypes:  []comply.Key{comply.StaffStatutory},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2022, time.March, 30), Amount: 36, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{
				Title:   "Tribunal Supremo STS 280/2022 (applied nationally)",
				Section: "Sala de lo Social, Rec 63/2020",
				URL:     "https://www.poderjudicial.es/",
			},
		},
	}
}
