// Package es_an registers Andalusia (SAS) healthcare scheduling regulations.
//
// Andalusia has been actively reforming guard duty duration, moving from
// 24-hour to 17-hour maximum guards. The SAS (Servicio Andaluz de Salud)
// is the largest regional health service in Spain by population.
//
// Key source: Acuerdo de 18 de julio de 2023 de la Mesa Sectorial de
// Sanidad (guard reform), BOJA publication.
package es_an

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:      comply.ESAN,
		Name:      "Andalusia",
		LocalName: "Andalucia",
		Type:      comply.Region,
		Parent:    comply.ES,
		Currency:  "EUR",
		TimeZone:  "Europe/Madrid",
		Rules:     rules(),
	}
}

func rules() []*comply.RuleDef {
	sasSource := comply.Source{
		Title:   "Acuerdo de la Mesa Sectorial de Sanidad de Andalucia (2023)",
		Section: "Jornada y guardias",
		URL:     "https://www.sspa.juntadeandalucia.es/servicioandaluzdesalud/",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxShiftHours,
			Name:        "Maximum Guard Duration (SAS Reform)",
			Description: "Andalusia is progressively reducing maximum guard duration from 24 to 17 hours. The reform applies to SAS statutory personnel in primary care and hospitals.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffStatutory},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2023, time.July, 18),
					Amount: 17,
					Unit:   comply.Hours,
					Per:    comply.PerShift,
					Exceptions: []string{
						"24-hour guards still permitted in specific services during transition period",
						"Full implementation timeline depends on staffing availability",
					},
				},
			},
			Source: sasSource,
			Notes:  "Andalusia is the first Spanish region to legislate a sub-24h guard maximum. The transition from 24h to 17h is being implemented progressively.",
		},
		{
			Key:         comply.RuleMinRestAfterGuard,
			Name:        "Post-Guard Rest (SAS)",
			Description: "After guard duty, the worker must receive rest equivalent to the guard duration. Post-guard rest counts as effective working time for annual hours computation.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			StaffTypes:  []comply.Key{comply.StaffStatutory},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2023, time.July, 18), Amount: 17, Unit: comply.Hours, Per: comply.PerOccurrence},
			},
			Source: sasSource,
		},
		{
			Key:         comply.RuleMaxAnnualHours,
			Name:        "SAS Annual Ordinary Hours",
			Description: "Annual ordinary working hours for SAS statutory personnel: 1,642.5 hours (same as national SERMAS standard).",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffStatutory},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2020, time.January, 1), Amount: 1642.5, Unit: comply.Hours, Per: comply.PerYear},
			},
			Source: comply.Source{
				Title:   "SAS Pacto de Condiciones de Trabajo",
				Section: "Jornada ordinaria",
				URL:     "https://www.sspa.juntadeandalucia.es/servicioandaluzdesalud/",
			},
		},
	}
}
