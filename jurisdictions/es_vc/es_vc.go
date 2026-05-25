// Package es_vc registers Valencian Community healthcare scheduling regulations.
// Conselleria de Sanitat Universal i Salut Publica manages healthcare through
// the regional health service.
package es_vc

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:      comply.ESVC,
		Name:      "Valencian Community",
		LocalName: "Comunitat Valenciana",
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
			Name:        "Annual Ordinary Hours (Valencian Health Service)",
			Description: "Annual ordinary working hours for statutory personnel in the Valencian health system.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffStatutory},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2020, time.January, 1), Amount: 1642.5, Unit: comply.Hours, Per: comply.PerYear},
			},
			Source: comply.Source{
				Title:   "Conselleria de Sanitat Universal i Salut Publica",
				Section: "Jornada ordinaria anual",
				URL:     "https://www.san.gva.es/",
			},
		},
		{
			Key:         comply.RuleMinWeeklyRest,
			Name:        "Minimum Weekly Rest (Valencian Health)",
			Description: "36 hours of uninterrupted weekly rest (24h weekly + 12h daily). Follows the national standard reinforced by Tribunal Supremo ruling.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			StaffTypes:  []comply.Key{comply.StaffStatutory},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2022, time.March, 30), Amount: 36, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{
				Title:   "Tribunal Supremo STS 280/2022 (applied to Valencian Community)",
				Section: "Sala de lo Social",
				URL:     "https://www.poderjudicial.es/",
			},
		},
	}
}
