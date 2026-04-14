package comply

// Key is a kebab-case identifier used for rules, staff types, and unit types.
type Key string

// Category groups related scheduling rules.
type Category string

const (
	CatWorkHours    Category = "work_hours"
	CatRest         Category = "rest"
	CatOvertime     Category = "overtime"
	CatStaffing     Category = "staffing"
	CatBreaks       Category = "breaks"
	CatOnCall       Category = "on_call"
	CatCompensation Category = "compensation"
	CatLeave        Category = "leave"
	CatNightWork    Category = "night_work"
)

// Work Hours.
const (
	RuleMaxWeeklyHours         Key = "max-weekly-hours"
	RuleMaxOrdinaryWeeklyHours Key = "max-ordinary-weekly-hours"
	RuleMaxCombinedWeeklyHours Key = "max-combined-weekly-hours"
	RuleMaxDailyHours          Key = "max-daily-hours"
	RuleMaxShiftHours          Key = "max-shift-hours"
	RuleMaxShiftTransition     Key = "max-shift-transition-hours"
	RuleMaxConsecutiveNights   Key = "max-consecutive-nights"
	RuleMaxConsecutiveDays     Key = "max-consecutive-days"
	RuleMaxAnnualHours         Key = "max-annual-hours"
)

// Rest.
const (
	RuleMinRestBetweenShifts Key = "min-rest-between-shifts"
	RuleMinRestAfterExtended Key = "min-rest-after-extended-shift"
	RuleMinWeeklyRest        Key = "min-weekly-rest"
	RuleDaysOffPerWeek       Key = "days-off-per-week"
	RuleMinDayOfRest         Key = "min-day-of-rest"
)

// Overtime.
const (
	RuleOvertimeDailyThreshold     Key = "overtime-daily-threshold"
	RuleOvertimeDailyRate          Key = "overtime-daily-rate"
	RuleDoubleTimeDailyThreshold   Key = "double-time-daily-threshold"
	RuleDoubleTimeDailyRate        Key = "double-time-daily-rate"
	RuleOvertimeWeeklyThreshold    Key = "overtime-weekly-threshold"
	RuleOvertimeWeeklyRate         Key = "overtime-weekly-rate"
	RuleOvertime7thDayRate         Key = "overtime-7th-day-rate"
	RuleDoubleTime7thDayThreshold  Key = "double-time-7th-day-threshold"
	RuleMaxOvertimeAnnual          Key = "max-overtime-annual"
	RuleMandatoryOTProhibited      Key = "mandatory-overtime-prohibited"
	RuleOvertime880Eligible        Key = "overtime-8-80-eligible"
	RuleOvertime880DailyThreshold  Key = "overtime-8-80-daily-threshold"
	RuleOvertime880PeriodThreshold Key = "overtime-8-80-period-threshold"
)

// Staffing (nurse-patient ratios use composite keys per unit).
const (
	RuleNursePatientRatioOR              Key = "nurse-patient-ratio-or"
	RuleNursePatientRatioEDTrauma        Key = "nurse-patient-ratio-ed-trauma"
	RuleNursePatientRatioICU             Key = "nurse-patient-ratio-icu"
	RuleNursePatientRatioNICU            Key = "nurse-patient-ratio-nicu"
	RuleNursePatientRatioLaborDelivery   Key = "nurse-patient-ratio-labor-delivery"
	RuleNursePatientRatioPACU            Key = "nurse-patient-ratio-pacu"
	RuleNursePatientRatioEDCritical      Key = "nurse-patient-ratio-ed-critical"
	RuleNursePatientRatioStepDown        Key = "nurse-patient-ratio-step-down"
	RuleNursePatientRatioAntepartum      Key = "nurse-patient-ratio-antepartum"
	RuleNursePatientRatioED              Key = "nurse-patient-ratio-ed"
	RuleNursePatientRatioPediatrics      Key = "nurse-patient-ratio-pediatrics"
	RuleNursePatientRatioPostpartumCplts Key = "nurse-patient-ratio-postpartum-couplets"
	RuleNursePatientRatioTelemetry       Key = "nurse-patient-ratio-telemetry"
	RuleNursePatientRatioOtherSpecialty  Key = "nurse-patient-ratio-other-specialty"
	RuleNursePatientRatioMedSurg         Key = "nurse-patient-ratio-med-surg"
	RuleNursePatientRatioPostpartum      Key = "nurse-patient-ratio-postpartum"
	RuleNursePatientRatioPsychiatric     Key = "nurse-patient-ratio-psychiatric"
	RuleNursePatientRatioICUCriticalCare Key = "nurse-patient-ratio-icu-critical-care"
)

// Breaks.
const (
	RuleMealBreakDuration        Key = "meal-break-duration"
	RuleMealBreakThreshold       Key = "meal-break-threshold"
	RuleSecondMealBreakThreshold Key = "second-meal-break-threshold"
	RuleRestBreakDuration        Key = "rest-break-duration"
	RuleRestBreakInterval        Key = "rest-break-interval"
)

// On-Call / Guards.
const (
	RuleMaxOnCallFrequency  Key = "max-on-call-frequency"
	RuleMaxGuardsMonthly    Key = "max-guards-monthly"
	RuleMinRestAfterGuard   Key = "min-rest-after-guard"
	RuleMoonlightingAllowed Key = "moonlighting-permitted"
)

// Night Work.
const (
	RuleNightPeriodStart    Key = "night-period-start"
	RuleNightPeriodEnd      Key = "night-period-end"
	RuleMaxNightShiftHours  Key = "max-night-shift-hours"
	RuleMaxNightConsecWeeks Key = "max-night-consecutive-weeks"
)

// Leave.
const (
	RuleMinAnnualLeaveDays Key = "min-annual-leave-days"
)

// Staff type constants.
const (
	StaffAll           Key = "all"
	StaffResident      Key = "resident"
	StaffResidentPGY1  Key = "resident-pgy1"
	StaffResidentPGY2P Key = "resident-pgy2-plus"
	StaffNurse         Key = "nurse"
	StaffNurseRN       Key = "nurse-rn"
	StaffNurseLPN      Key = "nurse-lpn"
	StaffNurseCNA      Key = "nurse-cna"
	StaffPhysician     Key = "physician"
	StaffAllied        Key = "allied-health"
	StaffStatutory     Key = "statutory-personnel"
	StaffVANurse       Key = "va-nurse"
)

// Hospital unit type constants.
const (
	UnitAll                Key = "all"
	UnitICU                Key = "icu"
	UnitNICU               Key = "nicu"
	UnitED                 Key = "ed"
	UnitEDTrauma           Key = "ed-trauma"
	UnitEDCritical         Key = "ed-critical"
	UnitMedSurg            Key = "med-surg"
	UnitPediatrics         Key = "pediatrics"
	UnitLaborDelivery      Key = "labor-delivery"
	UnitAntepartum         Key = "antepartum"
	UnitPostpartum         Key = "postpartum"
	UnitPostpartumCouplets Key = "postpartum-couplets"
	UnitTelemetry          Key = "telemetry"
	UnitStepDown           Key = "step-down"
	UnitPsychiatric        Key = "psychiatric"
	UnitOR                 Key = "operating-room"
	UnitPACU               Key = "pacu"
	UnitOtherSpecialty     Key = "other-specialty"
)

// Operator defines how a rule's value constrains scheduling.
type Operator string

const (
	OpLTE  Operator = "lte"  // <= maximum constraint
	OpGTE  Operator = "gte"  // >= minimum requirement
	OpEQ   Operator = "eq"   // == exact match
	OpBool Operator = "bool" // boolean flag (Amount: 1=true, 0=false)
)

// Enforcement indicates the legal force of a regulation.
type Enforcement string

const (
	Mandatory   Enforcement = "mandatory"   // Legally binding, penalties for violation.
	Recommended Enforcement = "recommended" // Strong recommendation, tracked for compliance.
	Advisory    Enforcement = "advisory"    // Guidance only, no enforcement.
)

// Unit describes what a RuleValue.Amount measures.
type Unit string

const (
	Hours            Unit = "hours"
	Minutes          Unit = "minutes"
	Days             Unit = "days"
	Count            Unit = "count"
	PatientsPerNurse Unit = "patients_per_nurse"
	Multiplier       Unit = "multiplier"
	Boolean          Unit = "boolean"
	HourOfDay        Unit = "hour_of_day"
	CalendarDays     Unit = "calendar_days"
	Weeks            Unit = "weeks"
)

// Per is the time period denominator for a rule value.
type Per string

const (
	PerShift      Per = "shift"
	PerDay        Per = "day"
	PerWeek       Per = "week"
	PerMonth      Per = "month"
	PerYear       Per = "year"
	PerPeriod     Per = "period"
	PerOccurrence Per = "occurrence"
)

// PeriodUnit is the time unit for averaging periods.
type PeriodUnit string

const (
	PeriodDays   PeriodUnit = "days"
	PeriodWeeks  PeriodUnit = "weeks"
	PeriodMonths PeriodUnit = "months"
)

// Scope defines which employers or facility types a rule applies to.
// This is critical because many healthcare regulations only cover specific
// settings (e.g., ACGME applies to accredited residency programs, not all
// hospitals; Spain's Estatuto Marco covers public health system staff only).
type Scope string

const (
	// ScopeAll means the rule applies to all employers in the jurisdiction.
	ScopeAll Scope = "all"

	// ScopePublicHealth means the rule applies only to public health system
	// employers (e.g., Spain's SNS, SERMAS, ICS, SAS).
	ScopePublicHealth Scope = "public_health"

	// ScopePrivate means the rule applies only to private sector employers.
	ScopePrivate Scope = "private"

	// ScopeHospitals means the rule applies to licensed hospitals (general
	// and special) but not clinics, nursing homes, or other facilities.
	ScopeHospitals Scope = "hospitals"

	// ScopeHealthcareEmployers is broader than hospitals: includes hospitals,
	// nursing homes, outpatient clinics, rehab facilities, residential care,
	// and similar. Used by NY Labor Law S167.
	ScopeHealthcareEmployers Scope = "healthcare_employers"

	// ScopeAccreditedPrograms means the rule applies only to institutions
	// with accredited training programs (e.g., ACGME-accredited residencies).
	// Not a law - enforced through accreditation, not statute.
	ScopeAccreditedPrograms Scope = "accredited_programs"

	// ScopeStateFacilities means the rule applies only to state-operated
	// facilities (e.g., CA Gov Code 19851.2 covers CDCR, State Hospitals, DVA, DDS).
	ScopeStateFacilities Scope = "state_facilities"

	// ScopeVA means the rule applies only to Veterans Affairs facilities.
	ScopeVA Scope = "va"

	// ScopeNursingHomes means the rule applies to nursing homes / long-term
	// care facilities but not hospitals.
	ScopeNursingHomes Scope = "nursing_homes"
)
