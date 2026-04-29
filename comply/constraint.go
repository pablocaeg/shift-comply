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

// categoryToConstraintType maps rule categories to constraint types.
// This replaces the old key-based map so new jurisdiction keys are handled
// automatically without code changes.
var categoryToConstraintType = map[Category]ConstraintType{
	CatWorkHours: ConstraintMaxHours,
	CatRest:      ConstraintMinRest,
	CatOnCall:    ConstraintMaxGuards,
	CatOvertime:  ConstraintMaxOvertime,
	CatBreaks:    ConstraintBreakRequired,
	CatStaffing:  ConstraintStaffingRatio,
	CatNightWork: ConstraintMaxConsecutive,
	CatLeave:     ConstraintMinDaysOff,
}

// refineConstraintType narrows the constraint type based on rule semantics.
// The category gives a rough mapping; this function refines it using the
// key, operator, and value unit.
func refineConstraintType(r *RuleDef, v *RuleValue, base ConstraintType) ConstraintType {
	switch r.Category {
	case CatWorkHours:
		if v.Per == PerShift {
			return ConstraintMaxShift
		}
		if v.Unit == Count {
			return ConstraintMaxConsecutive
		}
		return ConstraintMaxHours
	case CatRest:
		if v.Per == PerWeek && v.Unit == Days {
			return ConstraintMinDaysOff
		}
		if v.Per == PerWeek && v.Unit == Hours {
			return ConstraintMinRest
		}
		return ConstraintMinRest
	}
	return base
}

// GenerateConstraints produces optimizer-ready constraints from the effective
// rules for a jurisdiction. Each rule is translated into a Constraint struct
// with normalized fields that scheduling engines can consume directly.
//
// All rules are included. Boolean/policy rules use ConstraintPolicy. Rules
// whose category isn't mapped get ConstraintPolicy as a fallback — nothing
// is silently dropped.
func GenerateConstraints(code Code, opts ...QueryOption) []Constraint {
	rules := EffectiveRules(code, opts...)
	constraints := make([]Constraint, 0, len(rules))

	// Track days-off rules so we can derive max-consecutive-days.
	var daysOffRules []*RuleDef

	for _, r := range rules {
		v := r.Current()
		if v == nil {
			continue
		}

		var ct ConstraintType
		if r.Operator == OpBool {
			ct = ConstraintPolicy
		} else if base, ok := categoryToConstraintType[r.Category]; ok {
			ct = refineConstraintType(r, v, base)
		} else {
			ct = ConstraintPolicy
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
			Citation:      r.Source.Citation(),
		}

		if v.Averaged != nil {
			c.AveragedOverDays = averagingDays(v.Averaged)
		}

		constraints = append(constraints, c)

		// Collect days-off rules for deriving max consecutive days.
		if ct == ConstraintMinDaysOff && v.Unit == Days && v.Amount > 0 {
			daysOffRules = append(daysOffRules, r)
		}
	}

	// Derive max-consecutive-days from days-off rules.
	// "1 day off per 7" implies max 6 consecutive working days.
	for _, r := range daysOffRules {
		v := r.Current()
		if v == nil || v.Per != PerWeek {
			continue
		}
		maxConsecutive := 7 - v.Amount
		if maxConsecutive <= 0 {
			continue
		}
		constraints = append(constraints, Constraint{
			Type:          ConstraintMaxConsecutive,
			TimeScope:     PerWeek,
			FacilityScope: r.Scope,
			Limit:         maxConsecutive,
			LimitUnit:     Days,
			Operator:      OpLTE,
			StaffTypes:    r.StaffTypes,
			UnitTypes:     r.UnitTypes,
			Enforcement:   r.Enforcement,
			Jurisdiction:  code,
			RuleKey:       Key("derived-max-consecutive-days"),
			Citation:      "Derived from: " + r.Source.Citation(),
		})
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
