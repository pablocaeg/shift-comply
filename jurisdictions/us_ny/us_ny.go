// Package us_ny registers New York healthcare scheduling regulations:
// Bell Regulations / "Libby Zion Law" (10 NYCRR S 405.4, resident work hours),
// nurse mandatory overtime prohibition (NY Labor Law S 167), meal break
// requirements (S 162), and staffing committee requirements (PHL S 2805-t).
package us_ny

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

// New returns the New York jurisdiction definition.
func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USNY,
		Name:     "New York",
		Type:     comply.State,
		Parent:   comply.US,
		Currency: "USD",
		TimeZone: "America/New_York",
		Rules:    rules(),
	}
}

func rules() []*comply.RuleDef {
	r := make([]*comply.RuleDef, 0, 15)
	r = append(r, residentRules()...)
	r = append(r, nurseRules()...)
	r = append(r, generalRules()...)
	return r
}

// Bell Regulations — 10 NYCRR S 405.4(b)(6) ("Libby Zion Law")
// First state-level resident work hour regulations in the US (effective July 1, 1989).

func residentRules() []*comply.RuleDef {
	bellSource := comply.Source{
		Title:   "New York Codes, Rules and Regulations, Title 10",
		Section: "S 405.4(b)(6)",
		URL:     "https://www.law.cornell.edu/regulations/new-york/10-NYCRR-405.4",
	}

	return []*comply.RuleDef{
		{
			Key:         "ny-resident-max-consecutive-hours",
			Name:        "Maximum Consecutive Work Hours (Bell Regulations)",
			Description: "Residents may not work more than 24 consecutive hours. Stricter than ACGME 24+4 — New York does not allow the additional 4 transition hours.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffResident},
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1989, time.July, 1), Amount: 24, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: bellSource,
			Notes: "Adopted following the 1984 death of Libby Zion at New York Hospital. Preceded ACGME national standards by 14 years.",
		},
		{
			Key:         comply.RuleMaxWeeklyHours,
			Name:        "Maximum Weekly Hours (Bell Regulations)",
			Description: "Residents may not work more than 80 hours per week, averaged over 4 weeks.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffResident},
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:    comply.D(1989, time.July, 1),
					Amount:   80,
					Unit:     comply.Hours,
					Per:      comply.PerWeek,
					Averaged: &comply.AveragingPeriod{Count: 4, Unit: comply.PeriodWeeks},
				},
			},
			Source: bellSource,
		},
		{
			Key:         comply.RuleMinRestBetweenShifts,
			Name:        "Minimum Rest Between Assignments (Bell Regulations)",
			Description: "At least 8 nonworking hours between scheduled on-duty assignments.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			StaffTypes:  []comply.Key{comply.StaffResident},
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1989, time.July, 1), Amount: 8, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: bellSource,
			Notes: "Unlike ACGME's 'should have' (recommended) 8-hour rest, New York's is mandatory.",
		},
		{
			Key:         "ny-resident-ed-max-hours",
			Name:        "ED Maximum Shift (Residents)",
			Description: "Emergency department assignments limited to 12 consecutive hours.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffResident},
			UnitTypes:   []comply.Key{comply.UnitED},
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:      comply.D(1989, time.July, 1),
					Amount:     12,
					Unit:       comply.Hours,
					Per:        comply.PerShift,
					Exceptions: []string{"Commissioner may approve extensions to 15 hours for attendings under specific conditions"},
				},
			},
			Source: bellSource,
		},
		{
			Key:         comply.RuleDaysOffPerWeek,
			Name:        "Weekly Time Off (Bell Regulations)",
			Description: "At least one period of 24 consecutive nonworking hours per week.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			StaffTypes:  []comply.Key{comply.StaffResident},
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1989, time.July, 1), Amount: 1, Unit: comply.Days, Per: comply.PerWeek},
			},
			Source: bellSource,
		},
	}
}

// Nurse mandatory overtime — NY Labor Law S 167 / 12 NYCRR Part 177.

func nurseRules() []*comply.RuleDef {
	return []*comply.RuleDef{
		{
			Key:         comply.RuleMandatoryOTProhibited,
			Name:        "Mandatory Overtime Prohibition (Nurses)",
			Description: "Healthcare employers may not require nurses (RN or LPN) to work overtime. On-call time counts as work time for this determination.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpBool,
			StaffTypes:  []comply.Key{comply.StaffNurseRN, comply.StaffNurseLPN},
			Scope:       comply.ScopeHealthcareEmployers,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2009, time.July, 1),
					Amount: 1,
					Unit:   comply.Boolean,
					Exceptions: []string{
						"Health care disaster (natural or other) in county or contiguous county",
						"Federal, state, or local government declaration of emergency",
						"Unforeseen patient care emergency that could not be prudently planned for",
						"Nurse actively engaged in ongoing medical or surgical procedure",
					},
				},
			},
			Source: comply.Source{
				Title:   "New York Labor Law",
				Section: "S 167; 12 NYCRR Part 177",
				URL:     "https://dol.ny.gov/mandatory-overtime-nurses",
			},
			Notes: "Covers hospitals, nursing homes, outpatient clinics, rehab hospitals, residential care, drug/alcohol treatment, adult day health care, diagnostic centers. Penalties (2023 amendments): $1,000 first, $2,000 second, $3,000 third+ within 12 months.",
		},
	}
}

// General NY labor rules.

func generalRules() []*comply.RuleDef {
	return []*comply.RuleDef{
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Meal Break Requirement",
			Description: "Employees working more than 6 hours spanning 11 AM to 2 PM get at least 30 minutes. Shift workers (start 1 PM-6 AM) get 45 minutes midway.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1940, time.January, 1), Amount: 6, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: comply.Source{
				Title:   "New York Labor Law",
				Section: "S 162",
				URL:     "https://dol.ny.gov/day-rest-and-meal-periods",
			},
			Notes: "Employees starting before 11 AM continuing past 7 PM also get an additional 20-minute break between 5-7 PM.",
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Meal Break Duration",
			Description: "Standard: 30 minutes. Shift workers (1 PM-6 AM start): 45 minutes.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1940, time.January, 1), Amount: 30, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: comply.Source{
				Title:   "New York Labor Law",
				Section: "S 162",
				URL:     "https://dol.ny.gov/day-rest-and-meal-periods",
			},
		},
		{
			Key:         comply.RuleMinDayOfRest,
			Name:        "Day of Rest",
			Description: "Employers must provide at least 24 consecutive hours of rest per calendar week.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1940, time.January, 1), Amount: 24, Unit: comply.Hours, Per: comply.PerWeek},
			},
			Source: comply.Source{
				Title:   "New York Labor Law",
				Section: "S 161",
				URL:     "https://dol.ny.gov/day-rest-and-meal-periods",
			},
		},
		{
			Key:         "ny-spread-of-hours",
			Name:        "Spread of Hours Pay",
			Description: "When interval between start and end of workday exceeds 10 hours, employee is owed 1 additional hour at minimum wage.",
			Category:    comply.CatCompensation,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1960, time.January, 1), Amount: 10, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: comply.Source{
				Title:   "New York Codes, Rules and Regulations, Title 12",
				Section: "12 NYCRR 142-2.4",
				URL:     "https://www.law.cornell.edu/regulations/new-york/12-NYCRR-142-2.4",
			},
		},
		{
			Key:         "ny-call-in-pay",
			Name:        "Call-In / Reporting Time Pay",
			Description: "Employee who reports to work must be paid for at least 4 hours (or scheduled shift length, whichever is less).",
			Category:    comply.CatCompensation,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1960, time.January, 1), Amount: 4, Unit: comply.Hours, Per: comply.PerOccurrence},
			},
			Source: comply.Source{
				Title:   "New York Codes, Rules and Regulations, Title 12",
				Section: "12 NYCRR 142-2.3",
				URL:     "https://www.law.cornell.edu/regulations/new-york/12-NYCRR-142-2.3",
			},
		},
		{
			Key:         comply.RuleNursePatientRatioICUCriticalCare,
			Name:        "Nurse-Patient Ratio: ICU/Critical Care",
			Description: "New York requires a minimum 1:2 nurse-to-patient ratio for critical care / ICU units.",
			Category:    comply.CatStaffing,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffNurseRN},
			UnitTypes:   []comply.Key{comply.UnitICU},
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2021, time.June, 22), Amount: 2, Unit: comply.PatientsPerNurse, Per: comply.PerShift},
			},
			Source: comply.Source{
				Title:   "New York Public Health Law",
				Section: "S 2805-t (Clinical Staffing Committees)",
				URL:     "https://www.nysenate.gov/legislation/laws/PBH/2805-T",
			},
			Notes: "Signed June 22, 2021. Only specific ratio mandated statewide; other units rely on committee-developed plans.",
		},
	}
}
