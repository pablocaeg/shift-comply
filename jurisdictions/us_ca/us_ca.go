// Package us_ca registers California healthcare scheduling regulations:
// mandatory nurse-patient ratios (Title 22 CCR § 70217), overtime rules
// (Labor Code § 510), meal/rest break requirements, and mandatory overtime
// restrictions.
package us_ca

import (
	"time"

	"github.com/pablocaeg/shift-comply/comply"
)

func init() {
	comply.RegisterJurisdiction(New())
}

// New returns the California jurisdiction definition.
func New() *comply.JurisdictionDef {
	return &comply.JurisdictionDef{
		Code:     comply.USCA,
		Name:     "California",
		Type:     comply.State,
		Parent:   comply.US,
		Currency: "USD",
		TimeZone: "America/Los_Angeles",
		Rules:    rules(),
	}
}

func rules() []*comply.RuleDef {
	r := make([]*comply.RuleDef, 0, 40)
	r = append(r, nurseRatioRules()...)
	r = append(r, overtimeRules()...)
	r = append(r, breakRules()...)
	r = append(r, otherRules()...)
	return r
}

// Nurse-to-Patient Ratios - AB 394 (1999) / Title 22 CCR § 70217
// California is the ONLY US state with mandatory nurse-patient ratios.
// Ratios are minimums AT ALL TIMES - no averaging across shifts is permitted.

func nurseRatioRules() []*comply.RuleDef {
	ratioSource := comply.Source{
		Title:   "California Code of Regulations, Title 22",
		Section: "§ 70217",
		URL:     "https://www.law.cornell.edu/regulations/california/Cal-Code-Regs-Tit-22-SS-70217",
	}

	type ratio struct {
		key         comply.Key
		name        string
		unit        comply.Key
		amount      float64
		since       time.Time
		notes       string
		extraValues []*comply.RuleValue // for ratios that changed over time
	}

	ratios := []ratio{
		{comply.RuleNursePatientRatioOR, "Operating Room", comply.UnitOR, 1, comply.D(2004, time.January, 1), "One RN circulating nurse per patient; plus one scrub assistant (may be non-RN)", nil},
		{comply.RuleNursePatientRatioEDTrauma, "Emergency Dept - Trauma", comply.UnitEDTrauma, 1, comply.D(2004, time.January, 1), "RNs only", nil},
		{comply.RuleNursePatientRatioICU, "ICU / Critical Care", comply.UnitICU, 2, comply.D(2004, time.January, 1), "Includes burn, coronary care, acute respiratory care", nil},
		{comply.RuleNursePatientRatioNICU, "Neonatal ICU", comply.UnitNICU, 2, comply.D(2004, time.January, 1), "RNs only", nil},
		{comply.RuleNursePatientRatioLaborDelivery, "Labor & Delivery (active labor)", comply.UnitLaborDelivery, 2, comply.D(2004, time.January, 1), "", nil},
		{comply.RuleNursePatientRatioPACU, "Post-Anesthesia Recovery (PACU)", comply.UnitPACU, 2, comply.D(2004, time.January, 1), "", nil},
		{comply.RuleNursePatientRatioEDCritical, "Emergency Dept - Critical Care", comply.UnitEDCritical, 2, comply.D(2004, time.January, 1), "", nil},
		{
			comply.RuleNursePatientRatioStepDown, "Step-Down / Intermediate Care", comply.UnitStepDown, 3, comply.D(2008, time.January, 1), "Tightened from 1:4 in 2008",
			[]*comply.RuleValue{
				{Since: comply.D(2004, time.January, 1), Amount: 4, Unit: comply.PatientsPerNurse, Per: comply.PerShift},
			},
		},
		{comply.RuleNursePatientRatioAntepartum, "Antepartum (non-active labor)", comply.UnitAntepartum, 4, comply.D(2004, time.January, 1), "", nil},
		{comply.RuleNursePatientRatioED, "Emergency Dept - General", comply.UnitED, 4, comply.D(2004, time.January, 1), "Minimum 2 nurses present at all times when patients are receiving treatment", nil},
		{comply.RuleNursePatientRatioPediatrics, "Pediatrics", comply.UnitPediatrics, 4, comply.D(2004, time.January, 1), "", nil},
		{comply.RuleNursePatientRatioPostpartumCplts, "Postpartum Couplets (mother-baby)", comply.UnitPostpartumCouplets, 4, comply.D(2004, time.January, 1), "Max 8 patients total per nurse with multiple births", nil},
		{
			comply.RuleNursePatientRatioTelemetry, "Telemetry", comply.UnitTelemetry, 4, comply.D(2008, time.January, 1), "Tightened from 1:5 in 2008",
			[]*comply.RuleValue{
				{Since: comply.D(2004, time.January, 1), Amount: 5, Unit: comply.PatientsPerNurse, Per: comply.PerShift},
			},
		},
		{
			comply.RuleNursePatientRatioOtherSpecialty, "Other Specialty Care", comply.UnitOtherSpecialty, 4, comply.D(2008, time.January, 1), "Tightened from 1:5 in 2008",
			[]*comply.RuleValue{
				{Since: comply.D(2004, time.January, 1), Amount: 5, Unit: comply.PatientsPerNurse, Per: comply.PerShift},
			},
		},
		{
			comply.RuleNursePatientRatioMedSurg, "Medical/Surgical", comply.UnitMedSurg, 5, comply.D(2005, time.January, 1), "Was 1:6 during 2004 transition year",
			[]*comply.RuleValue{
				{Since: comply.D(2004, time.January, 1), Amount: 6, Unit: comply.PatientsPerNurse, Per: comply.PerShift},
			},
		},
		{comply.RuleNursePatientRatioPostpartum, "Postpartum (mothers only)", comply.UnitPostpartum, 6, comply.D(2004, time.January, 1), "", nil},
		{comply.RuleNursePatientRatioPsychiatric, "Psychiatric", comply.UnitPsychiatric, 6, comply.D(2004, time.January, 1), "Licensed nurses includes psychiatric technicians", nil},
	}

	rules := make([]*comply.RuleDef, 0, len(ratios))
	for _, rt := range ratios {
		values := []*comply.RuleValue{
			{Since: rt.since, Amount: rt.amount, Unit: comply.PatientsPerNurse, Per: comply.PerShift},
		}
		if rt.extraValues != nil {
			values = append(values, rt.extraValues...)
		}

		rules = append(rules, &comply.RuleDef{
			Key:         rt.key,
			Name:        "Nurse-Patient Ratio: " + rt.name,
			Description: "Maximum patients per nurse - " + rt.name + ". Enforced at all times, no averaging.",
			Category:    comply.CatStaffing,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffNurseRN},
			UnitTypes:   []comply.Key{rt.unit},
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values:      values,
			Source:      ratioSource,
			Notes:       rt.notes,
		})
	}
	return rules
}

// Overtime Rules - California Labor Code § 510

func overtimeRules() []*comply.RuleDef {
	otSource := comply.Source{
		Title:   "California Labor Code",
		Section: "§ 510(a)",
		URL:     "https://leginfo.legislature.ca.gov/faces/codes_displaySection.xhtml?sectionNum=510.&lawCode=LAB",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleOvertimeDailyThreshold,
			Name:        "Daily Overtime Threshold",
			Description: "Work exceeding 8 hours in one workday compensated at 1.5x regular rate.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2000, time.January, 1), Amount: 8, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: otSource,
		},
		{
			Key:         comply.RuleOvertimeDailyRate,
			Name:        "Daily Overtime Rate",
			Description: "1.5x regular rate for hours exceeding 8 per workday.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2000, time.January, 1), Amount: 1.5, Unit: comply.Multiplier, Per: comply.PerDay},
			},
			Source: otSource,
		},
		{
			Key:         comply.RuleDoubleTimeDailyThreshold,
			Name:        "Daily Double Time Threshold",
			Description: "Work exceeding 12 hours in one day compensated at 2x regular rate.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2000, time.January, 1), Amount: 12, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: otSource,
		},
		{
			Key:         comply.RuleDoubleTimeDailyRate,
			Name:        "Daily Double Time Rate",
			Description: "2x regular rate for hours exceeding 12 per workday.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2000, time.January, 1), Amount: 2.0, Unit: comply.Multiplier, Per: comply.PerDay},
			},
			Source: otSource,
		},
		{
			Key:         comply.RuleOvertime7thDayRate,
			Name:        "7th Consecutive Day Overtime Rate",
			Description: "First 8 hours on the 7th consecutive workday compensated at 1.5x. Hours beyond 8 on the 7th day at 2x.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2000, time.January, 1), Amount: 1.5, Unit: comply.Multiplier, Per: comply.PerDay},
			},
			Source: otSource,
		},
		{
			Key:         comply.RuleDoubleTime7thDayThreshold,
			Name:        "7th Day Double Time Threshold",
			Description: "Hours exceeding 8 on the 7th consecutive workday compensated at 2x.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2000, time.January, 1), Amount: 8, Unit: comply.Hours, Per: comply.PerDay},
			},
			Source: otSource,
		},
		{
			Key:         comply.RuleMandatoryOTProhibited,
			Name:        "Mandatory Overtime Restrictions (Private Sector Healthcare)",
			Description: "Employees assigned to 12-hour shifts shall not be required to work more than 12 hours in 24 hours unless the Chief Nursing Officer declares a healthcare emergency and all reasonable steps to provide voluntary staffing have been taken. Even in emergencies, no employee shall be required to work more than 16 hours in 24 hours unless by voluntary mutual agreement.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpBool,
			StaffTypes:  []comply.Key{comply.StaffNurseRN, comply.StaffNurseLPN, comply.StaffNurseCNA},
			Scope:       comply.ScopeHospitals,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2001, time.January, 1),
					Amount: 1,
					Unit:   comply.Boolean,
					Exceptions: []string{
						"Chief Nursing Officer declares healthcare emergency and voluntary staffing has been attempted",
						"Even in emergency, max 16 hours unless voluntary mutual agreement",
					},
				},
			},
			Source: comply.Source{
				Title:   "IWC Wage Order No. 5-2001 (Healthcare Industry)",
				Section: "Section 3(H)",
				URL:     "https://www.dir.ca.gov/iwc/wageorder5_010102.html",
			},
			Notes: "IWC Wage Order 5 covers private-sector healthcare employers since 2001. State facility employees were added later via Gov. Code 19851.2 (AB 840, effective 2016). Private sector protections predated state facility protections.",
		},
		{
			Key:         "mandatory-overtime-prohibited-state",
			Name:        "Mandatory Overtime Prohibition (State Facilities)",
			Description: "State facilities may not require nurses or CNAs to work beyond their scheduled shift. Applies to CDCR, State Hospitals, DVA, and DDS facilities.",
			Category:    comply.CatOvertime,
			Operator:    comply.OpBool,
			StaffTypes:  []comply.Key{comply.StaffNurseRN, comply.StaffNurseLPN, comply.StaffNurseCNA},
			Scope:       comply.ScopeStateFacilities,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2016, time.January, 1),
					Amount: 1,
					Unit:   comply.Boolean,
					Exceptions: []string{
						"Nurse actively engaged in a surgical procedure until completion",
						"Catastrophic event causing large influx of patients requiring immediate treatment",
						"Declared emergency (national, state, or municipal)",
					},
				},
			},
			Source: comply.Source{
				Title:   "California Government Code",
				Section: "S 19851.2 (AB 840)",
				URL:     "https://leginfo.legislature.ca.gov/faces/billTextClient.xhtml?bill_id=201520160AB840",
			},
		},
	}
}

// Meal and Rest Break Rules - Labor Code § 512, IWC Wage Order 5

func breakRules() []*comply.RuleDef {
	mealSource := comply.Source{
		Title:   "California Labor Code",
		Section: "§ 512(a)",
		URL:     "https://leginfo.legislature.ca.gov/faces/codes_displaySection.xhtml?lawCode=LAB&sectionNum=512.",
	}
	restSource := comply.Source{
		Title:   "IWC Wage Order No. 5-2001 (Healthcare Industry)",
		Section: "Section 12",
		URL:     "https://www.dir.ca.gov/iwc/wageorder5_010102.html",
	}

	return []*comply.RuleDef{
		{
			Key:         comply.RuleMealBreakThreshold,
			Name:        "Meal Break Trigger - First Break",
			Description: "A 30-minute unpaid meal break must begin before the end of the 5th hour of work.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2000, time.January, 1), Amount: 5, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: mealSource,
			Notes:  "May be waived by mutual consent if total shift is 6 hours or less.",
		},
		{
			Key:         comply.RuleMealBreakDuration,
			Name:        "Meal Break Duration",
			Description: "Meal breaks must be at least 30 minutes, unpaid, duty-free.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2000, time.January, 1), Amount: 30, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: mealSource,
		},
		{
			Key:         comply.RuleSecondMealBreakThreshold,
			Name:        "Meal Break Trigger - Second Break",
			Description: "A second 30-minute meal break required when shift exceeds 10 hours. Must begin before end of 10th hour.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2000, time.January, 1), Amount: 10, Unit: comply.Hours, Per: comply.PerShift},
			},
			Source: mealSource,
			Notes:  "Healthcare workers on shifts over 8 hours may voluntarily waive one of two meal periods per IWC Wage Order 5, Section 11(D), upheld by CA Supreme Court in Gerard v. Orange Coast Memorial (2018).",
		},
		{
			Key:         comply.RuleRestBreakDuration,
			Name:        "Rest Break Duration",
			Description: "10-minute paid, duty-free rest breaks. Cannot be waived.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2000, time.January, 1), Amount: 10, Unit: comply.Minutes, Per: comply.PerOccurrence},
			},
			Source: restSource,
		},
		{
			Key:         comply.RuleRestBreakInterval,
			Name:        "Rest Break Interval",
			Description: "One 10-minute rest break per 4-hour work period (or major fraction thereof). Under 3.5 hours: none. 3.5-6 hours: 1 break. 6-10 hours: 2 breaks. 10-14 hours: 3 breaks.",
			Category:    comply.CatBreaks,
			Operator:    comply.OpLTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{Since: comply.D(2000, time.January, 1), Amount: 4, Unit: comply.Hours, Per: comply.PerOccurrence},
			},
			Source: restSource,
			Notes:  "Penalty for missed breaks: one additional hour of pay at regular rate per workday (Labor Code § 226.7).",
		},
	}
}

// Other California Rules - Day of Rest, Alternative Workweek

func otherRules() []*comply.RuleDef {
	return []*comply.RuleDef{
		{
			Key:         comply.RuleMinDayOfRest,
			Name:        "Day of Rest",
			Description: "Every employee is entitled to one day's rest in seven. Employer may not cause employee to work more than 6 days in 7.",
			Category:    comply.CatRest,
			Operator:    comply.OpGTE,
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(1937, time.January, 1),
					Amount: 1,
					Unit:   comply.Days,
					Per:    comply.PerWeek,
					Exceptions: []string{
						"When total hours worked do not exceed 30 in any week or 6 in any one day (Labor Code § 556)",
						"When nature of employment reasonably requires 7+ consecutive days, rest may accumulate monthly (Labor Code § 554)",
					},
				},
			},
			Source: comply.Source{
				Title:   "California Labor Code",
				Section: "§§ 551, 552",
				URL:     "https://leginfo.legislature.ca.gov/faces/codes_displaySection.xhtml?sectionNum=551.&lawCode=LAB",
			},
			Notes: "Per Mendoza v. Nordstrom (2017), day of rest applies per employer-established workweek, not on a rolling basis.",
		},
		{
			Key:         "max-shift-hours-non-resident",
			Name:        "Maximum Shift Under Alternative Workweek (Healthcare)",
			Description: "Non-resident healthcare workers on an alternative workweek schedule (e.g., 3x12) may not work more than 12 hours without overtime. Cannot exceed 13 hours in any 24-hour period. Does not apply to ACGME residents who are governed by the 24-hour ACGME limit.",
			Category:    comply.CatWorkHours,
			Operator:    comply.OpLTE,
			StaffTypes:  []comply.Key{comply.StaffNurseRN, comply.StaffNurseLPN, comply.StaffNurseCNA, comply.StaffPhysician, comply.StaffAllied},
			Enforcement: comply.Mandatory,
			Values: []*comply.RuleValue{
				{
					Since:  comply.D(2000, time.January, 1),
					Amount: 12,
					Unit:   comply.Hours,
					Per:    comply.PerShift,
				},
			},
			Source: comply.Source{
				Title:   "IWC Wage Order No. 5-2001 (Healthcare Industry)",
				Section: "Section 3(B)(8)",
				URL:     "https://www.dir.ca.gov/iwc/wageorder5_010102.html",
			},
			Notes: "Requires 2/3 secret ballot approval. Must provide 14 days written disclosure. Report results to DLSE within 30 days.",
		},
	}
}
