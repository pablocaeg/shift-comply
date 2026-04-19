// Package us_ny registers New York healthcare scheduling regulations:
// resident work hour limits (10 NYCRR 405.4), nurse mandatory overtime
// prohibition (Labor Law S167), meal/rest break requirements, staffing
// committee mandates, and compensation rules.
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
	r := make([]*comply.RuleDef, 0, 12)
	r = append(r, residentRules()...)
	r = append(r, nurseRules()...)
	r = append(r, breakRules()...)
	r = append(r, compensationRules()...)
	r = append(r, staffingRules()...)
	return r
}

// Resident Work Hour Rules - 10 NYCRR 405.4(b)(6)

func residentRules() []*comply.RuleDef {
	dohSource := comply.Source{
		Title:   "New York Codes, Rules and Regulations, Title 10",
		Section: "405.4(b)(6)",
		URL:     "https://www.health.ny.gov/facilities/hospital/reports/resident_work_hours/",
	}

	return []*comply.RuleDef{
		{
			Key:         "ny-resident-max-consecutive-hours",
			Name:        "Resident Maximum Consecutive Hours",
			Description: "Residents may not work more than 24 consecutive hours. DOH allows up to 3 additional hours for transition activities (handoff, education, continuity), not 4 as ACGME permits.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffResident},
			Scope:       comply.ScopeAccreditedPrograms,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1989, time.July, 1), Amount: 24, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: dohSource,
			Notes:  "Predates ACGME by 14 years. DOH allows up to 3-hour transition (not 4 as ACGME). Regulation text says 24 hours.",
		},
		{
			Key:         comply.RuleMaxWeeklyHours,
			Name:        "Resident Maximum Weekly Hours",
			Description: "Residents may not work more than 80 hours per week, averaged over 4 weeks.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffResident},
			Scope:       comply.ScopeAccreditedPrograms,
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
			Source: dohSource,
		},
		{
			Key:         comply.RuleMinRestBetweenShifts,
			Name:        "Resident Minimum Rest Between Shifts",
			Description: "Residents must have at least 8 hours of rest between scheduled shifts. Unlike ACGME's recommended 8 hours, New York's 8 hours is mandatory.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			StaffTypes:  []comply.Key{comply.StaffResident},
			Scope:       comply.ScopeAccreditedPrograms,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1989, time.July, 1), Amount: 8, Unit: comply.Hours, Per: comply.PerOccurrence},
			},
			Source: dohSource,
			Notes:  "Unlike ACGME recommended 8h, NY's 8h is mandatory.",
		},
		{
			Key:         "ny-resident-ed-max-hours",
			Name:        "Resident ED Maximum Shift Hours",
			Description: "Residents assigned to emergency departments may not work shifts exceeding 12 hours.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffResident},
			UnitTypes:   []comply.Key{comply.UnitED},
			Scope:       comply.ScopeAccreditedPrograms,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1989, time.July, 1), Amount: 12, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: dohSource,
		},
		{
			Key:         comply.RuleDaysOffPerWeek,
			Name:        "Resident Days Off Per Week",
			Description: "Residents must have at least 1 day (24 consecutive hours) off per week.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			StaffTypes:  []comply.Key{comply.StaffResident},
			Scope:       comply.ScopeAccreditedPrograms,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1989, time.July, 1), Amount: 1, Unit: comply.Days, Per: comply.PerWeek},
			},
			Source: dohSource,
		},
	}
}

// Nurse Mandatory Overtime - NY Labor Law S167

func nurseRules() []*comply.RuleDef {
	return []*comply.RuleDef{
		{
			Key:         comply.RuleMandatoryOTProhibited,
			Name:        "Nurse Mandatory Overtime Prohibition",
			Description: "Employers may not require nurses to work overtime. Covers hospitals, nursing homes, clinics, rehab facilities, residential care, and drug/alcohol treatment facilities.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpBool,
			Scope:       comply.ScopeHealthcareEmployers,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2009, time.July, 1),
					Amount: 1,
					Unit:   comply.Boolean,
					Exceptions: []string{
						"health care disaster",
						"government emergency declaration",
						"unforeseen patient care emergency",
						"ongoing procedure",
					},
				},
			},
			Source: comply.Source{
				Title:   "New York Labor Law",
				Section: "Section 167; 12 NYCRR Part 177",
				URL:     "https://dol.ny.gov/mandatory-overtime-nurses",
			},
			Notes: "Covers hospitals, nursing homes, clinics, rehab, residential care, drug/alcohol treatment. Penalties: $1,000 first, $2,000 second (12 months), $3,000 third+ per 2023 amendments.",
		},
	}
}

// Meal and Rest Break Rules - NY Labor Law S162, S161

func breakRules() []*comply.RuleDef {
	mealSource := comply.Source{
		Title:   "New York Labor Law",
		Section: "Section 162",
		URL:     "https://dol.ny.gov/day-rest-and-meal-periods",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Meal Break Trigger",
			Description: "Employees working a shift of more than 6 hours must receive a 30-minute meal break.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1940, time.January, 1), Amount: 6, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: mealSource,
			Notes:  "Shift workers (1pm-6am start): 45 minutes. Factory workers: 60 minutes.",
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Meal Break Duration",
			Description: "Meal breaks must be at least 30 minutes.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1940, time.January, 1), Amount: 30, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: mealSource,
		},
		{
			Key:         comply.RuleMinDayOfRest,
			Name:        "Day of Rest",
			Description: "Every employee is entitled to 24 consecutive hours of rest per calendar week.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(1940, time.January, 1), Amount: 1, Unit: comply.Days, Per: comply.PerWeek},
			},
			Source: comply.Source{
				Title:   "New York Labor Law",
				Section: "Section 161",
				URL:     "https://dol.ny.gov/day-rest-and-meal-periods",
			},
		},
	}
}

// Compensation Rules - 12 NYCRR

func compensationRules() []*comply.RuleDef {
	return []*comply.RuleDef{
		{
			Key:         "ny-spread-of-hours",
			Name:        "Spread of Hours Pay",
			Description: "When the spread of hours (time between the start of the first work period and end of the last work period in a day) exceeds 10 hours, the employee must be paid 1 extra hour at minimum wage.",
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
			Name:        "Call-In Pay",
			Description: "Employees who report to work must be paid a minimum of 4 hours at minimum wage, regardless of hours actually worked.",
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
	}
}

// Staffing Rules - NY Public Health Law

func staffingRules() []*comply.RuleDef {
	return []*comply.RuleDef{
		{
			Key:         "ny-icu-nurse-ratio",
			Name:        "ICU Nurse-Patient Ratio",
			Description: "ICU nurse-to-patient ratio must not exceed 1:2. Staffing committees required.",
			Category:    comply.CatStaffing,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffNurseRN},
			UnitTypes:   []comply.Key{comply.UnitICU},
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2023, time.January, 1), Amount: 2, Unit: comply.PatientsPerNurse, Per: comply.PerShift},
			},
			Source: comply.Source{
				Title:   "New York Public Health Law",
				Section: "Section 2805-t",
				URL:     "https://www.nysenate.gov/legislation/laws/PBH/2805-T",
			},
			Notes: "Staffing committees required. ICU 1:2 is the only mandated ratio.",
		},
	}
}
