// Package us registers US federal healthcare scheduling regulations:
// ACGME duty hour requirements, FLSA overtime rules, and OSHA guidelines.
package us

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

// New returns the US federal jurisdiction definition.
func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.US,
		Name:     "United States",
		Type:     comply.Country,
		Currency: "USD",
		TimeZone: "America/New_York",
		Rules:    rules(),
	}
}

func rules() []*comply.RuleDef {
	return append(append(acgmeRules(), flsaRules()...), oshaRules()...)
}

// ACGME Common Program Requirements - Duty Hour Standards
// Effective: July 1, 2017 (current revision unifying PGY-1 and PGY-2+)
// Source: ACGME Common Program Requirements (Residency), Section VI
//
// IMPORTANT: ACGME is a private nonprofit, NOT a government body.
// These rules are enforced through ACCREDITATION, not statute or regulation.
// Non-compliance leads to accreditation actions (warning, probation,
// withdrawal), which causes loss of Medicare GME funding. This is distinct
// from state law (only NY and IL have codified resident hours into statute).
// Scope is set to ScopeAccreditedPrograms accordingly.

func acgmeRules() []*comply.RuleDef {
	acgmeSource := comply.Source{
		Title:   "ACGME Common Program Requirements (Residency)",
		Section: "Section VI - The Learning and Working Environment",
		URL:     "https://www.acgme.org/programs-and-institutions/programs/common-program-requirements/",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMaxWeeklyHours,
			Name:        "Maximum Weekly Clinical and Educational Work Hours",
			Description: "Inclusive of all in-house clinical and educational activities, clinical work done from home, and all moonlighting.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffResident},
			Scope:       comply.ScopeAccreditedPrograms,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2003, time.July, 1),
					Amount:   80,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 4, Unit: comply.PeriodWeeks},
					Exceptions: []string{
						"Specialty Review Committee may approve up to 88 hours/week for specific rotations with sound educational rationale",
					},
				},
			},
			Source: comply.Source{
				Title:   acgmeSource.Title,
				Section: "Section VI.F.1",
				URL:     acgmeSource.URL,
			},
		},
		{
			Key:         comply.RuleMaxShiftHours,
			Name:        "Maximum Continuous Scheduled Clinical Assignment",
			Description: "Maximum hours of continuous duty. Applies to all PGY levels since July 2017. Prior to 2017, PGY-1 was limited to 16 hours.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffResident},
			Scope:       comply.ScopeAccreditedPrograms,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2017, time.July, 1),
					Amount: 24,
					Unit:   comply.Hours,
					Per:    comply.PerShift,
				},
			},
			Source: comply.Source{
				Title:   acgmeSource.Title,
				Section: "Section VI.F",
				URL:     acgmeSource.URL,
			},
			Notes: "The 2017 revision eliminated the separate 16-hour PGY-1 limit based on the FIRST trial and iCOMPARE study results.",
		},
		{
			Key:         comply.RuleMaxShiftTransition,
			Name:        "Additional Hours After 24-Hour Period for Transitions",
			Description: "Up to 4 additional hours after a 24-hour duty period for activities related to patient safety, transitions of care, resident education, and continuity of care. No new patients may be accepted.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffResident},
			Scope:       comply.ScopeAccreditedPrograms,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2003, time.July, 1),
					Amount: 4,
					Unit:   comply.Hours,
					Per:    comply.PerShift,
				},
			},
			Source: comply.Source{
				Title:   acgmeSource.Title,
				Section: "Section VI.F",
				URL:     acgmeSource.URL,
			},
		},
		{
			Key:         comply.RuleMinRestBetweenShifts,
			Name:        "Minimum Time Off Between Scheduled Clinical Work",
			Description: "Residents should have 8 hours free of clinical work between scheduled periods. This uses 'should have' language - a strong recommendation, not an absolute mandate.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			StaffTypes:  []comply.Key{comply.StaffResident},
			Scope:       comply.ScopeAccreditedPrograms,
			Enforcement: comply.Recommended,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2003, time.July, 1),
					Amount: 8,
					Unit:   comply.Hours,
					Per:    comply.PerShift,
				},
			},
			Source: comply.Source{
				Title:   acgmeSource.Title,
				Section: "Section VI.F",
				URL:     acgmeSource.URL,
			},
		},
		{
			Key:         comply.RuleMinRestAfterExtended,
			Name:        "Minimum Time Off After 24 Hours of In-House Call",
			Description: "Mandatory 14 hours free of clinical work and education after 24 hours of in-house call. Uses 'must have' language.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			StaffTypes:  []comply.Key{comply.StaffResident},
			Scope:       comply.ScopeAccreditedPrograms,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2003, time.July, 1),
					Amount: 14,
					Unit:   comply.Hours,
					Per:    comply.PerOccurrence,
				},
			},
			Source: comply.Source{
				Title:   acgmeSource.Title,
				Section: "Section VI.F",
				URL:     acgmeSource.URL,
			},
		},
		{
			Key:         comply.RuleDaysOffPerWeek,
			Name:        "Minimum Days Off Per Week",
			Description: "One continuous 24-hour period free of all clinical work and education per week. At-home call cannot be assigned on free days.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			StaffTypes:  []comply.Key{comply.StaffResident},
			Scope:       comply.ScopeAccreditedPrograms,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2003, time.July, 1),
					Amount:   1,
					Unit:     comply.Days,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 4, Unit: comply.PeriodWeeks},
				},
			},
			Source: comply.Source{
				Title:   acgmeSource.Title,
				Section: "Section VI.F",
				URL:     acgmeSource.URL,
			},
		},
		{
			Key:         comply.RuleMaxOnCallFrequency,
			Name:        "Maximum In-House Call Frequency",
			Description: "In-house call no more frequently than every 3rd night, averaged over 4 weeks.",
			Category:    comply.CatOnCall,
			Operator:    comply.OpGTE,
			StaffTypes:  []comply.Key{comply.StaffResident},
			Scope:       comply.ScopeAccreditedPrograms,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(2003, time.July, 1),
					Amount:   3,
					Unit:     comply.Days,
					Per:      comply.PerOccurrence,
					Averaged: &comply.AveragingPeriod{Count: 4, Unit: comply.PeriodWeeks},
				},
			},
			Source: comply.Source{
				Title:   acgmeSource.Title,
				Section: "Section VI.F",
				URL:     acgmeSource.URL,
			},
			Notes: "Means no more than ~1 in 3 nights on in-house call. At-home call is not subject to this limit.",
		},
		{
			Key:         comply.RuleMaxConsecutiveNights,
			Name:        "Maximum Consecutive Nights of Night Float",
			Description: "Night float rotations limited to 6 consecutive nights. Specialty Review Committees may impose further restrictions.",
			Category:    comply.CatNightWork,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffResident},
			Scope:       comply.ScopeAccreditedPrograms,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2003, time.July, 1),
					Amount: 6,
					Unit:   comply.Count,
					Per:    comply.PerPeriod,
				},
			},
			Source: comply.Source{
				Title:   acgmeSource.Title,
				Section: "Section VI.F",
				URL:     acgmeSource.URL,
			},
		},
		{
			Key:         "moonlighting-prohibited-pgy1",
			Name:        "PGY-1 Moonlighting Prohibition",
			Description: "First-year residents (interns) are not permitted to moonlight, internally or externally.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpBool,
			StaffTypes:  []comply.Key{comply.StaffResidentPGY1},
			Scope:       comply.ScopeAccreditedPrograms,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2003, time.July, 1),
					Amount: 1,
					Unit:   comply.Boolean,
				},
			},
			Source: comply.Source{
				Title:   acgmeSource.Title,
				Section: "Section VI.F.5",
				URL:     acgmeSource.URL,
			},
		},
		{
			Key:         comply.RuleMoonlightingAllowed,
			Name:        "PGY-2+ Moonlighting Permitted with Approval",
			Description: "Requires written advance approval from the program director. All moonlighting hours count toward the 80-hour weekly maximum.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpBool,
			StaffTypes:  []comply.Key{comply.StaffResidentPGY2P},
			Scope:       comply.ScopeAccreditedPrograms,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2003, time.July, 1),
					Amount: 1,
					Unit:   comply.Boolean,
				},
			},
			Source: comply.Source{
				Title:   acgmeSource.Title,
				Section: "Section VI.F.5",
				URL:     acgmeSource.URL,
			},
			Notes: "Internal and external moonlighting hours must be counted toward the 80-hour weekly maximum. Must not interfere with patient safety or educational goals.",
		},
	}
}

// FLSA (Fair Labor Standards Act) - Overtime Rules for Healthcare

func flsaRules() []*comply.RuleDef {
	flsaSource := comply.Source{
		Title:   "Fair Labor Standards Act",
		Section: "29 U.S.C. § 207",
		URL:     "https://www.dol.gov/agencies/whd/flsa",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleOvertimeWeeklyThreshold,
			Name:        "Standard Overtime Threshold",
			Description: "Overtime pay required for non-exempt employees working more than 40 hours in a 7-day workweek.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(1938, time.June, 25),
					Amount: 40,
					Unit:   comply.Hours,
					Per:    comply.PerWeek,
				},
			},
			Source: comply.Source{
				Title:   flsaSource.Title,
				Section: "29 U.S.C. § 207(a)(1)",
				URL:     flsaSource.URL,
			},
		},
		{
			Key:         comply.RuleOvertimeWeeklyRate,
			Name:        "Standard Overtime Rate",
			Description: "Non-exempt employees receive 1.5x regular rate for hours exceeding 40 per workweek.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(1938, time.June, 25),
					Amount: 1.5,
					Unit:   comply.Multiplier,
					Per:    comply.PerWeek,
				},
			},
			Source: comply.Source{
				Title:   flsaSource.Title,
				Section: "29 U.S.C. § 207(a)(1)",
				URL:     flsaSource.URL,
			},
		},
		{
			Key:         comply.RuleOvertime880Eligible,
			Name:        "8/80 Overtime System Eligibility",
			Description: "Hospitals and residential care establishments may use a 14-day work period with overtime thresholds of 8 hours/day or 80 hours/14-day period. Requires prior agreement with employees.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpBool,
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(1966, time.September, 23),
					Amount: 1,
					Unit:   comply.Boolean,
				},
			},
			Source: comply.Source{
				Title:   flsaSource.Title,
				Section: "29 U.S.C. § 207(j); 29 CFR 778.601",
				URL:     "https://www.dol.gov/agencies/whd/fact-sheets/54-healthcare-overtime",
			},
			Notes: "Eligible establishments: hospitals, skilled nursing facilities, nursing facilities, assisted living, residential care. Requires fixed, regularly recurring 14-day period and prior employee agreement.",
		},
		{
			Key:         comply.RuleOvertime880DailyThreshold,
			Name:        "8/80 System - Daily Overtime Threshold",
			Description: "Under the 8/80 system, overtime is owed for hours exceeding 8 in a single workday.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpGTE,
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(1966, time.September, 23),
					Amount: 8,
					Unit:   comply.Hours,
					Per:    comply.PerDay,
				},
			},
			Source: comply.Source{
				Title:   flsaSource.Title,
				Section: "29 U.S.C. § 207(j)",
				URL:     "https://www.dol.gov/agencies/whd/fact-sheets/54-healthcare-overtime",
			},
		},
		{
			Key:         comply.RuleOvertime880PeriodThreshold,
			Name:        "8/80 System - 14-Day Period Overtime Threshold",
			Description: "Under the 8/80 system, overtime is owed for hours exceeding 80 in the 14-day work period.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpGTE,
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(1966, time.September, 23),
					Amount:   80,
					Unit:     comply.Hours,
					Per:      comply.PerPeriod,
					Averaged: &comply.AveragingPeriod{Count: 14, Unit: comply.PeriodDays},
				},
			},
			Source: comply.Source{
				Title:   flsaSource.Title,
				Section: "29 U.S.C. § 207(j)",
				URL:     "https://www.dol.gov/agencies/whd/fact-sheets/54-healthcare-overtime",
			},
		},
	}
}

// OSHA / VA - Fatigue and Nurse Hour Limits

func oshaRules() []*comply.RuleDef {
	return []*comply.RuleDef{
		{
			Key:         "va-nurse-max-consecutive-hours",
			Name:        "VA Nurse Maximum Consecutive Hours",
			Description: "Veterans Affairs nurses may not work more than 12 consecutive hours, except in emergency care situations.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffVANurse},
			Scope:       comply.ScopeVA,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2004, time.December, 3),
					Amount: 12,
					Unit:   comply.Hours,
					Per:    comply.PerShift,
					Exceptions: []string{
						"Emergency care situations",
					},
				},
			},
			Source: comply.Source{
				Title:   "VA Healthcare Personnel Enhancement Act of 2004",
				Section: "P.L. 108-445",
				URL:     "https://www.congress.gov/bill/108th-congress/senate-bill/2484",
			},
		},
		{
			Key:         "va-nurse-max-weekly-hours",
			Name:        "VA Nurse Maximum Weekly Hours",
			Description: "Veterans Affairs nurses may not work more than 60 hours in any 7-day period, except in emergencies.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffVANurse},
			Scope:       comply.ScopeVA,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2004, time.December, 3),
					Amount: 60,
					Unit:   comply.Hours,
					Per:    comply.PerWeek,
					Exceptions: []string{
						"Emergency care situations",
					},
				},
			},
			Source: comply.Source{
				Title:   "VA Healthcare Personnel Enhancement Act of 2004",
				Section: "P.L. 108-445",
				URL:     "https://www.congress.gov/bill/108th-congress/senate-bill/2484",
			},
			Notes: "This is the only federal statutory limit on nurse work hours. Applies only to VA-employed nurses, not private sector.",
		},
	}
}
