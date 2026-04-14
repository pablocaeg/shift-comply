package comply

// Constraint is an optimizer-ready scheduling constraint generated from
// one or more RuleDefs. This is the format that scheduling engines
// (like YouShift's optimizer) consume directly.
type Constraint struct {
	// Type identifies the constraint kind.
	Type ConstraintType `json:"type"`

	// TimeScope is the time period of the constraint.
	TimeScope Per `json:"time_scope"`

	// FacilityScope identifies which employers/facilities this applies to.
	FacilityScope Scope `json:"facility_scope,omitempty"`

	// Limit is the numeric boundary.
	Limit float64 `json:"limit"`

	// LimitUnit describes what the limit measures.
	LimitUnit Unit `json:"limit_unit"`

	// Operator is how the limit constrains (lte, gte, eq, bool).
	Operator Operator `json:"operator"`

	// AveragedOver is the averaging window in days. 0 means hard limit.
	AveragedOverDays int `json:"averaged_over_days,omitempty"`

	// StaffTypes this constraint applies to. Empty means all.
	StaffTypes []Key `json:"staff_types,omitempty"`

	// UnitTypes (hospital units) this constraint applies to. Empty means all.
	UnitTypes []Key `json:"unit_types,omitempty"`

	// Enforcement level.
	Enforcement Enforcement `json:"enforcement"`

	// Citation is the legal source for this constraint.
	Citation string `json:"citation"`

	// Jurisdiction is where this constraint originates.
	Jurisdiction Code `json:"jurisdiction"`

	// RuleKey links back to the originating rule.
	RuleKey Key `json:"rule_key"`
}

// ConstraintType categorizes constraints for optimizer consumption.
type ConstraintType string

const (
	ConstraintMaxHours       ConstraintType = "max_hours"
	ConstraintMinRest        ConstraintType = "min_rest"
	ConstraintMaxShift       ConstraintType = "max_shift"
	ConstraintMaxConsecutive ConstraintType = "max_consecutive"
	ConstraintMinDaysOff     ConstraintType = "min_days_off"
	ConstraintStaffingRatio  ConstraintType = "staffing_ratio"
	ConstraintMaxOvertime    ConstraintType = "max_overtime"
	ConstraintBreakRequired  ConstraintType = "break_required"
	ConstraintMaxGuards      ConstraintType = "max_guards"
	ConstraintPolicy         ConstraintType = "policy" // boolean policies
)

// ruleToConstraintType maps rule keys to constraint types.
var ruleToConstraintType = map[Key]ConstraintType{
	RuleMaxWeeklyHours:         ConstraintMaxHours,
	RuleMaxOrdinaryWeeklyHours: ConstraintMaxHours,
	RuleMaxCombinedWeeklyHours: ConstraintMaxHours,
	RuleMaxDailyHours:          ConstraintMaxHours,
	RuleMaxAnnualHours:         ConstraintMaxHours,
	RuleMaxShiftHours:          ConstraintMaxShift,
	RuleMaxShiftTransition:     ConstraintMaxShift,
	RuleMaxConsecutiveNights:   ConstraintMaxConsecutive,
	RuleMaxConsecutiveDays:     ConstraintMaxConsecutive,
	RuleMinRestBetweenShifts:   ConstraintMinRest,
	RuleMinRestAfterExtended:   ConstraintMinRest,
	RuleMinWeeklyRest:          ConstraintMinRest,
	RuleDaysOffPerWeek:         ConstraintMinDaysOff,
	RuleMinDayOfRest:           ConstraintMinDaysOff,
	RuleMaxOvertimeAnnual:      ConstraintMaxOvertime,
	RuleMealBreakThreshold:     ConstraintBreakRequired,
	RuleMealBreakDuration:      ConstraintBreakRequired,
	RuleRestBreakDuration:      ConstraintBreakRequired,
	RuleRestBreakInterval:      ConstraintBreakRequired,
	RuleMaxGuardsMonthly:       ConstraintMaxGuards,
	RuleMinRestAfterGuard:      ConstraintMinRest,
	RuleMaxOnCallFrequency:     ConstraintMaxGuards,
}

// GenerateConstraints produces optimizer-ready constraints from the effective
// rules for a jurisdiction. Each rule is translated into a Constraint struct
// with normalized fields that scheduling engines can consume directly.
func GenerateConstraints(code Code, opts ...QueryOption) []Constraint {
	rules := EffectiveRules(code, opts...)
	constraints := make([]Constraint, 0, len(rules))

	for _, r := range rules {
		v := r.Current()
		if v == nil {
			continue
		}

		ct, ok := ruleToConstraintType[r.Key]
		if !ok {
			// Staffing ratios use composite keys
			if r.Category == CatStaffing {
				ct = ConstraintStaffingRatio
			} else if r.Operator == OpBool {
				ct = ConstraintPolicy
			} else {
				continue // rule type not yet mapped
			}
		}

		c := Constraint{
			Type:          ct,
			TimeScope:     v.Per,
			FacilityScope: r.Scope,
			Limit:         v.Amount,
			LimitUnit:     v.Unit,
			Operator:      r.Operator,
			StaffTypes:    r.StaffTypes,
			UnitTypes:     r.UnitTypes,
			Enforcement:   r.Enforcement,
			Jurisdiction:  code,
			RuleKey:       r.Key,
		}

		c.Citation = r.Source.Citation()

		// Convert averaging period to days
		if v.Averaged != nil {
			switch v.Averaged.Unit {
			case PeriodDays:
				c.AveragedOverDays = v.Averaged.Count
			case PeriodWeeks:
				c.AveragedOverDays = v.Averaged.Count * 7
			case PeriodMonths:
				c.AveragedOverDays = v.Averaged.Count * 30
			}
		}

		constraints = append(constraints, c)
	}

	return constraints
}

// Schedule represents a work schedule to be validated against jurisdiction rules.
type Schedule struct {
	// Jurisdiction is the jurisdiction code to validate against.
	Jurisdiction Code `json:"jurisdiction"`

	// FacilityScope filters rules by facility type (hospitals, public_health, etc.)
	FacilityScope Scope `json:"facility_scope,omitempty"`

	// Shifts contains the individual shift assignments.
	Shifts []Shift `json:"shifts"`
}

// Shift represents a single shift assignment.
type Shift struct {
	// StaffID identifies the worker.
	StaffID string `json:"staff_id"`

	// StaffType is the worker's role.
	StaffType Key `json:"staff_type"`

	// UnitType is the hospital unit for this shift.
	UnitType Key `json:"unit_type,omitempty"`

	// Start is the shift start time (RFC 3339).
	Start string `json:"start"`

	// End is the shift end time (RFC 3339).
	End string `json:"end"`

	// OnCall indicates this is an on-call shift.
	OnCall bool `json:"on_call,omitempty"`
}

// Violation represents a single compliance violation found in a schedule.
type Violation struct {
	// RuleKey identifies the violated rule.
	RuleKey Key `json:"rule_key"`

	// RuleName is the human-readable rule name.
	RuleName string `json:"rule_name"`

	// Severity indicates the legal force of the violated rule.
	Severity Enforcement `json:"severity"`

	// StaffID identifies the affected worker.
	StaffID string `json:"staff_id"`

	// Message describes the violation.
	Message string `json:"message"`

	// Citation is the legal reference.
	Citation string `json:"citation"`

	// Actual is the actual value found.
	Actual float64 `json:"actual"`

	// Limit is the constraint that was violated.
	Limit float64 `json:"limit"`
}

// ComplianceReport is the result of validating a schedule.
type ComplianceReport struct {
	// Jurisdiction that was checked against.
	Jurisdiction Code `json:"jurisdiction"`

	// Result is pass, fail, or warnings.
	Result string `json:"result"`

	// Violations lists all found violations.
	Violations []Violation `json:"violations"`

	// ConstraintsChecked is the number of constraints evaluated.
	ConstraintsChecked int `json:"constraints_checked"`
}
