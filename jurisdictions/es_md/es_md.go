// Package es_md registers Community of Madrid (SERMAS) healthcare scheduling
// regulations.
//
// Key sources:
//   - Resolucion de 26 de febrero de 2021 de la Direccion General de RRHH y
//     Relaciones Laborales del SERMAS (post-guard rest rules).
//   - STS 280/2022, Sala de lo Social, Rec 63/2020, dated March 30, 2022
//     (Tribunal Supremo ruling confirming 36-hour weekly rest for MIR residents).
//
// IMPORTANT: The SERMAS resolution (Feb 26, 2021) PREDATES the Supreme Court
// ruling (March 30, 2022). The instructions were issued during pending litigation
// and a MIR strike, not as a response to the final ruling. The Supreme Court
// later confirmed and reinforced the same 36-hour rest principle.
//
// Scope: Applies only to SERMAS statutory employees (personal estatutario del
// Servicio Madrileno de Salud). Does not apply to private hospitals in Madrid.
package es_md

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

// New returns the Community of Madrid jurisdiction definition.
func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:      comply.ESMD,
		Name:      "Community of Madrid",
		LocalName: "Comunidad de Madrid",
		Type:      comply.Region,
		Parent:    comply.ES,
		Currency:  "EUR",
		TimeZone:  "Europe/Madrid",
		Rules:     rules(),
	}
}

func rules() []*comply.RuleDef {
	sermasResolution := comply.Source{
		Title:   "Resolucion de 26 de febrero de 2021 de la Direccion General de RRHH y Relaciones Laborales del SERMAS",
		Section: "",
		URL:     "https://www.redaccionmedica.com/autonomias/madrid/sermas-aprueba-nueva-orden-descansos-obligatorios-guardias-1935",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMinWeeklyRest,
			Name:        "Minimum Weekly Rest (SERMAS - 36 Hours)",
			Description: "Minimum 36 hours of uninterrupted weekly rest (24 hours weekly rest + 12 hours daily rest). After Saturday guard duty, rest must be scheduled on the following Monday. If Monday rest is not possible due to care needs, 72 hours of uninterrupted rest must be guaranteed within 14 days from the guard.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			StaffTypes:  []comply.Key{comply.StaffStatutory},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2021, time.February, 26),
					Amount: 36,
					Unit:   comply.Hours,
					Per:    comply.PerWeek,
					Exceptions: []string{
						"If Monday rest not possible due to care needs: 72 hours uninterrupted rest within 14 days from the guard",
					},
				},
			},
			Source: sermasResolution,
			Notes:  "SERMAS resolution issued Feb 26, 2021, during pending litigation. Confirmed by Tribunal Supremo STS 280/2022 (Sala de lo Social, Rec 63/2020, March 30, 2022, ECLI:ES:TS:2022:1543, Amyts v. SERMAS). The CGPJ press release was April 21, 2022 - this is often confused with the ruling date. Legal basis: Arts. 37.1 and 37.2 Estatuto de los Trabajadores, applied supplementarily per Arts. 3, 5, 16, 17 EU Directive 2003/88/CE. Full resolution text: https://madrid.ccoo.es/a3a2c8aa00f0a56b9d6b702db0eab6da000045.pdf",
		},
		{
			Key:         "md-rest-not-effective-work",
			Name:        "Post-Guard Rest Does Not Count as Effective Working Time",
			Description: "Rest periods after guard duty do not have the character of effective work, nor are they taken into consideration for compliance with the ordinary working day. Staff must still complete the full 1,642.5 annual ordinary hours.",
			Category:    comply.CatRest,
			Operator:    comply.OpBool,
			StaffTypes:  []comply.Key{comply.StaffStatutory},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2021, time.February, 26), Amount: 1, Unit: comply.Boolean},
			},
			Source: sermasResolution,
			Notes:  "Exact text: 'Los periodos de descanso no tendran el caracter ni la consideracion de trabajo efectivo, ni en ningun caso tomados en consideracion para el cumplimiento de la jornada ordinaria.' This differs from Catalonia ICS, where post-guard rest IS counted as a full working day.",
		},
		{
			Key:         comply.RuleMaxAnnualHours,
			Name:        "SERMAS Annual Ordinary Hours",
			Description: "Annual ordinary working hours for SERMAS statutory hospital staff: 1,642.5 hours, based on 219 working days and 146 free days. Guard hours (jornada complementaria) are separate and do not count toward this total.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffStatutory},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2020, time.January, 1), Amount: 1642.5, Unit: comply.Hours, Per: comply.PerYear},
			},
			Source: comply.Source{
				Title:   "SERMAS Working Conditions (CSIF / CCOO documentation)",
				Section: "Jornada ordinaria anual",
				URL:     "https://feccoo-madrid.org/3867f2fddb62b34d9a69fac3512fe2ed000063.pdf",
			},
			Notes: "Weekly equivalent: 37.5 hours average. Applies to Monday-Friday daytime shift staff.",
		},
	}
}
