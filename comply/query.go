package comply

import "time"

// QueryOption configures a rules query.
type QueryOption func(*queryOptions)

type queryOptions struct {
	staffType Key
	unitType  Key
	date      time.Time
	category  Category
	scope     Scope
}

// ForStaff filters rules to those applicable to the given staff type.
func ForStaff(staffType Key) QueryOption {
	return func(o *queryOptions) { o.staffType = staffType }
}

// ForUnit filters rules to those applicable to the given hospital unit.
func ForUnit(unitType Key) QueryOption {
	return func(o *queryOptions) { o.unitType = unitType }
}

// OnDate filters rules to those effective on the given date.
func OnDate(date time.Time) QueryOption {
	return func(o *queryOptions) { o.date = date }
}

// InCategory filters rules to the given category.
func InCategory(cat Category) QueryOption {
	return func(o *queryOptions) { o.category = cat }
}

// ForScope filters rules to those applicable to the given facility scope
// (e.g., ScopePublicHealth, ScopeHospitals, ScopeAll).
func ForScope(scope Scope) QueryOption {
	return func(o *queryOptions) { o.scope = scope }
}

// EffectiveRules returns all rules that apply to a jurisdiction,
// including inherited rules from parent jurisdictions.
// When a child jurisdiction defines a rule with the same key as a parent,
// the child's rule takes precedence.
func EffectiveRules(code Code, opts ...QueryOption) []*RuleDef {
	j := For(code)
	if j == nil {
		return nil
	}

	var o queryOptions
	for _, opt := range opts {
		opt(&o)
	}

	chain := j.Chain()
	seen := make(map[Key]struct{})
	var result []*RuleDef

	for _, jur := range chain {
		for _, r := range jur.Rules {
			if _, exists := seen[r.Key]; exists {
				continue
			}
			if matchesFilter(r, &o) {
				seen[r.Key] = struct{}{}
				result = append(result, r)
			}
		}
	}

	return result
}

func matchesFilter(r *RuleDef, o *queryOptions) bool {
	if o.staffType != "" && !r.AppliesToStaff(o.staffType) {
		return false
	}
	if o.unitType != "" && !r.AppliesToUnit(o.unitType) {
		return false
	}
	if o.category != "" && r.Category != o.category {
		return false
	}
	if o.scope != "" && r.Scope != "" && r.Scope != ScopeAll && r.Scope != o.scope {
		return false
	}
	if !o.date.IsZero() && r.Value(o.date) == nil {
		return false
	}
	return true
}

// Comparison holds the differences between two jurisdictions' effective rules.
type Comparison struct {
	Left      Code        `json:"left"`
	Right     Code        `json:"right"`
	OnlyLeft  []*RuleDef  `json:"only_left,omitempty"`
	OnlyRight []*RuleDef  `json:"only_right,omitempty"`
	Different []*RulePair `json:"different,omitempty"`
	Same      []*RulePair `json:"same,omitempty"`
}

// RulePair pairs rules from two jurisdictions that share the same key.
type RulePair struct {
	Key   Key      `json:"key"`
	Left  *RuleDef `json:"left"`
	Right *RuleDef `json:"right"`
}

// Compare returns the differences between two jurisdictions' effective rules.
func Compare(left, right Code, opts ...QueryOption) *Comparison {
	leftRules := EffectiveRules(left, opts...)
	rightRules := EffectiveRules(right, opts...)

	leftMap := make(map[Key]*RuleDef, len(leftRules))
	for _, r := range leftRules {
		leftMap[r.Key] = r
	}
	rightMap := make(map[Key]*RuleDef, len(rightRules))
	for _, r := range rightRules {
		rightMap[r.Key] = r
	}

	comp := &Comparison{Left: left, Right: right}

	for _, r := range leftRules {
		if rr, ok := rightMap[r.Key]; ok {
			pair := &RulePair{Key: r.Key, Left: r, Right: rr}
			if rulesEquivalent(r, rr) {
				comp.Same = append(comp.Same, pair)
			} else {
				comp.Different = append(comp.Different, pair)
			}
		} else {
			comp.OnlyLeft = append(comp.OnlyLeft, r)
		}
	}

	for _, r := range rightRules {
		if _, ok := leftMap[r.Key]; !ok {
			comp.OnlyRight = append(comp.OnlyRight, r)
		}
	}

	return comp
}

func rulesEquivalent(a, b *RuleDef) bool {
	va, vb := a.Current(), b.Current()
	if va == nil || vb == nil {
		return va == vb
	}
	return va.Amount == vb.Amount && va.Unit == vb.Unit && a.Operator == b.Operator
}
