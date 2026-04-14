// Package es_ct registers Catalonia (ICS) healthcare scheduling regulations
// from the III Acord de condicions de treball del personal estatutari de l'ICS.
//
// Signed: November 22, 2023 at Mesa Sectorial de Negociacio de Sanitat.
// Ratified by Govern: December 12, 2023.
// Published in DOGC: January 23, 2024 (Resolucio EMT/74/2024, DOGC no. 9085).
// Conveni code: 79100162132024.
// Signatories: CCOO, UGT, Metges de Catalunya, SATSE (workers); ICS (employer).
//
// Scope: Over 55,000 statutory workers across 8 hospitals and 289 primary care
// teams. Only applies to ICS public health service employees, not private hospitals.
//
// Sources verified against:
// - DOGC publication via CIDO registry
// - ICS official news (ics.gencat.cat)
// - Infermeres de Catalunya (agreement text by topic)
// - Metges de Catalunya press release
// - Redaccion Medica coverage
package es_ct

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

// New returns the Catalonia jurisdiction definition.
func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:      comply.ESCT,
		Name:      "Catalonia",
		LocalName: "Catalunya",
		Type:      comply.Region,
		Parent:    comply.ES,
		Currency:  "EUR",
		TimeZone:  "Europe/Madrid",
		Rules:     rules(),
	}
}

func rules() []*comply.RuleDef {
	// Primary source: III Acord de condicions de treball del personal estatutari
	// de l'Institut Catala de la Salut (ICS), signed November 22, 2023.
	// Published in DOGC January 23, 2024 via Resolucio EMT/74/2024.
	icsSource := comply.Source{
		Title:   "III Acord de condicions de treball de l'ICS",
		Section: "Signed November 22, 2023; DOGC January 23, 2024 (Resolucio EMT/74/2024)",
		URL:     "https://ics.gencat.cat/ca/detall/noticia/III-acord-ics",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxGuardsMonthly,
			Name:        "Maximum Guard Modules Per Month (ICS)",
			Description: "Maximum of 4 modules of continuous care (guardias) per month on average, with a maximum of 1 falling on a weekend or holiday. The 'on average' qualifier means a given month may exceed 4 if compensated in other months.",
			Category:    comply.CatOnCall,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffStatutory},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2024, time.January, 23),
					Amount: 4,
					Unit:   comply.Count,
					Per:    comply.PerMonth,
					Exceptions: []string{
						"Averaged: a given month may exceed 4 if compensated by fewer in other months",
						"Maximum of 1 guard on weekend/holiday (protective cap, not a minimum obligation)",
					},
				},
			},
			Source: comply.Source{
				Title:   icsSource.Title,
				Section: "Complements d'atencio continuada",
				URL:     "https://llibreria.infermeresdecatalunya.cat/books/ics/page/complements-datencio-continuada",
			},
			Notes: "Overrides Spain national limit of 7 guards/month (RD 1146/2006, Art. 5.1.c). Applies to all statutory personnel in categories providing continuous care, not only doctors. Module compensation rates cover medical specialists, emergency doctors, primary care doctors, nursing, and TCAI staff.",
		},
		{
			Key:         comply.RuleMaxAnnualHours,
			Name:        "Maximum Combined Annual Hours (ICS)",
			Description: "Combined ordinary and complementary working hours shall not exceed 2,187 hours per year (equivalent to 48 hours/week), as the sum of ordinary working hours (jornada ordinaria) and complementary hours (jornada complementaria).",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffStatutory},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2024, time.January, 23),
					Amount: 2187,
					Unit:   comply.Hours,
					Per:    comply.PerYear,
					Exceptions: []string{
						"Up to 150 additional voluntary hours annually under basic sectoral regulations fall outside this maximum",
					},
				},
			},
			Source: comply.Source{
				Title:   icsSource.Title,
				Section: "Jornada complementaria",
				URL:     "https://llibreria.infermeresdecatalunya.cat/books/ics/page/jornada-complementaria",
			},
			Notes: "Ordinary hours within this total: 1,642 hrs/year for doctors; 1,599 hrs/year for day-shift non-medical staff; 1,445 hrs/year for fixed night-shift staff.",
		},
		{
			Key:         comply.RuleMinRestAfterGuard,
			Name:        "Post-Guard Rest Protection (ICS)",
			Description: "Rest after a 24-hour guard is guaranteed and counts within the working day (blindatge del descans postguardia). The rest period is treated as effective working time for salary and annual hours computation - the worker does not lose a day from their ordinary schedule.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			StaffTypes:  []comply.Key{comply.StaffStatutory},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2024, time.January, 23),
					Amount: 24,
					Unit:   comply.Hours,
					Per:    comply.PerOccurrence,
				},
			},
			Source: comply.Source{
				Title:   icsSource.Title,
				Section: "Blindatge del descans postguardia",
				URL:     "https://metgesdecatalunya.cat/ca/premsa/metges-de-catalunya-signa-el-iii-acord-de-lics-i-engega-una-nova-etapa-daccio-sindical-sense-renunciar-a-res",
			},
			Notes: "Differs from national rule (RD 1146/2006 Art. 5.1.b: 12 hours rest after guard). ICS counts the post-guard rest day as a full working day for annual hours purposes.",
		},
		// Age-based guard exemptions - applies to specialist physicians (metges especialistes) only.
		{
			Key:         "ct-guard-exemption-age-60",
			Name:        "Guard Exemption - Age 60+ (Specialist Physicians)",
			Description: "From January 1, 2024, specialist physicians turning 60 may request exemption from guard duty. Compensated at EUR 3,000/year (complement de metge senior). Requires minimum average of 2 monthly guards in 5 years preceding request.",
			Category:    comply.CatOnCall,
			Operator:    comply.OpBool,
			StaffTypes:  []comply.Key{comply.StaffPhysician},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2024, time.January, 1),
					Amount: 1,
					Unit:   comply.Boolean,
					Exceptions: []string{
						"ICS management may temporarily suspend exemptions for up to 1 year when service needs require it, with 10% premium above standard guard rates",
						"Requires 3 months notice to employer",
						"Incompatible with continuous care supplements for actual guard performance",
					},
				},
			},
			Source: comply.Source{
				Title:   icsSource.Title,
				Section: "Exempció de guardies / Complement per exoneració per edat",
				URL:     "https://llibreria.infermeresdecatalunya.cat/books/ics/page/exempció-de-guàrdies",
			},
			Notes: "Professionals with fewer guards but at least 1 monthly 24-hour guard on average receive 50% (EUR 1,500/year). Reference: average of 5 years immediately preceding request. Applies retroactively to specialists who requested exemptions up to 5 years before agreement signing.",
		},
		{
			Key:         "ct-guard-exemption-age-55",
			Name:        "Guard Exemption - Age 55+ (Specialist Physicians)",
			Description: "From January 1, 2025, specialist physicians turning 55 may request exemption from guard duty. Compensated at EUR 2,500/year. Requires minimum average of 2 monthly guards (1 weekday + 1 holiday) in 5 years preceding request.",
			Category:    comply.CatOnCall,
			Operator:    comply.OpBool,
			StaffTypes:  []comply.Key{comply.StaffPhysician},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2025, time.January, 1), Amount: 1, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   icsSource.Title,
				Section: "Exempció de guardies / Complement per exoneració per edat",
				URL:     "https://llibreria.infermeresdecatalunya.cat/books/ics/page/nou-tram-datencio-continuada-com-a-complement-per-exoneracio-de-guardies-per-edat-del-metge-especialista",
			},
			Notes: "50% rate (EUR 1,250/year) for those with fewer guards but at least 1 monthly 24-hour guard on average.",
		},
		{
			Key:         "ct-guard-exemption-age-50",
			Name:        "Guard Exemption - Age 50+ (Specialist Physicians)",
			Description: "From January 1, 2027, specialist physicians turning 50 may request exemption from guard duty. Compensated at EUR 2,000/year. Requires minimum average of 3 monthly guards (2 weekday + 1 holiday) in 5 years preceding request.",
			Category:    comply.CatOnCall,
			Operator:    comply.OpBool,
			StaffTypes:  []comply.Key{comply.StaffPhysician},
			Scope:       comply.ScopePublicHealth,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2027, time.January, 1), Amount: 1, Unit: comply.Boolean},
			},
			Source: comply.Source{
				Title:   icsSource.Title,
				Section: "Exempció de guardies / Complement per exoneració per edat",
				URL:     "https://redaccionmedica.com/autonomias/cataluna/el-tercer-convenio-ics-pacta-la-exencion-de-guardias-a-los-50-y-mas-sueldos-1057",
			},
			Notes: "50% rate (EUR 1,000/year) for those with fewer guards but at least 1 monthly 24-hour guard on average. Localized guards count at 50%.",
		},
	}
}
