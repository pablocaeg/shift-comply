package comply

import "time"

// RuleDef defines a single scheduling regulation within a jurisdiction.
type RuleDef struct {
	// Key uniquely identifies this rule within its jurisdiction.
	// Uses kebab-case (e.g., "max-weekly-hours", "nurse-patient-ratio-icu").
	Key Key `json:"key"`

	// Name is a human-readable rule name.
	Name string `json:"name"`

	// Description explains the rule in detail.
	Description string `json:"description"`

	// Category groups related rules (work_hours, rest, overtime, staffing, etc.)
	Category Category `json:"category"`

	// Operator defines how the value constrains behavior.
	Operator Operator `json:"operator"`

	// StaffTypes limits which staff types this rule applies to.
	// Nil or empty means all staff types.
	StaffTypes []Key `json:"staff_types,omitempty"`

	// UnitTypes limits which hospital unit types this applies to.
	// Nil or empty means all units.
	UnitTypes []Key `json:"unit_types,omitempty"`

	// Scope defines which employers or facility types this rule applies to.
	// Empty or "all" means all employers in the jurisdiction.
	Scope Scope `json:"scope,omitempty"`

	// Enforcement indicates the legal force of this rule.
	Enforcement Enforcement `json:"enforcement"`

	// Values are time-versioned constraint values, ordered newest-first.
	// The first entry where Since <= target date is the effective value.
	Values []*RuleValue `json:"values"`

	// Source is the primary legal citation.
	Source Source `json:"source"`

	// Notes provides additional context not captured in other fields.
	Notes string `json:"notes,omitempty"`
}

// RuleValue is a time-versioned constraint value for a rule.
type RuleValue struct {
	// Since is when this value became effective.
	Since time.Time `json:"since"`

	// Amount is the numeric constraint value.
	// For boolean rules, 1 = true and 0 = false.
	Amount float64 `json:"amount"`

	// Unit describes what the amount measures (hours, days, patients_per_nurse, etc.)
	Unit Unit `json:"unit"`

	// Per is the time period denominator (per week, per month, per shift, etc.)
	// Empty if not applicable (e.g., a boolean rule or per-occurrence rule).
	Per Per `json:"per,omitempty"`

	// Averaged defines the averaging window, if the rule is time-averaged.
	// Nil if the rule is a hard per-period limit (no averaging).
	Averaged *AveragingPeriod `json:"averaged,omitempty"`

	// Exceptions lists conditions under which this value does not apply.
	Exceptions []string `json:"exceptions,omitempty"`
}

// AveragingPeriod defines the time window over which a limit is averaged.
type AveragingPeriod struct {
	Count int        `json:"count"` // Number of period units.
	Unit  PeriodUnit `json:"unit"`  // The time unit (days, weeks, months).
}

// Value returns the effective RuleValue for the given date.
// Returns nil if the rule was not yet in effect on that date.
func (r *RuleDef) Value(date time.Time) *RuleValue {
	for _, v := range r.Values {
		if !v.Since.After(date) {
			return v
		}
	}
	return nil
}

// Current returns the most recent RuleValue, or nil if no values exist.
func (r *RuleDef) Current() *RuleValue {
	if len(r.Values) == 0 {
		return nil
	}
	return r.Values[0]
}

// AppliesToStaff returns true if this rule applies to the given staff type.
func (r *RuleDef) AppliesToStaff(staffType Key) bool {
	if len(r.StaffTypes) == 0 {
		return true
	}
	for _, s := range r.StaffTypes {
		if s == staffType {
			return true
		}
	}
	return false
}

// AppliesToUnit returns true if this rule applies to the given unit type.
func (r *RuleDef) AppliesToUnit(unitType Key) bool {
	if len(r.UnitTypes) == 0 {
		return true
	}
	for _, u := range r.UnitTypes {
		if u == unitType {
			return true
		}
	}
	return false
}

// D creates a time.Time from year, month, day for use in RuleValue.Since.
func D(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
