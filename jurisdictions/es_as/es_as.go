// Package es_as registers Asturias healthcare scheduling regulations.
// SESPA (Servicio de Salud del Principado de Asturias) is the regional health service.
package es_as

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:      comply.ESAS,
		Name:      "Asturias",
		LocalName: "Principado de Asturias",
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
			Name:        "Annual Ordinary Hours (SESPA)",
			Description: "Annual ordinary working hours for SESPA statutory personnel.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffStatutory},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2020, time.January, 1), Amount: 1642.5, Unit: comply.Hours, Per: comply.PerYear},
			},
			Source: comply.Source{
				Title:   "SESPA (Servicio de Salud del Principado de Asturias)",
				Section: "Jornada ordinaria anual",
				URL:     "https://www.astursalud.es/",
			},
		},
		{
			Key:         comply.RuleMinWeeklyRest,
			Name:        "Minimum Weekly Rest (SESPA)",
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
