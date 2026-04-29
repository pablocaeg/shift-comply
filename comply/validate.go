package comply

import (
	"fmt"
	"math"
	"sort"
	"time"
)

// Validate checks a schedule against the regulations of its jurisdiction
// and returns a compliance report listing all violations with legal citations.
//
// It validates time-based rules that can be derived from shift start/end times:
// max shift duration, max weekly hours (with averaging), min rest between shifts,
// min rest after extended shifts, days off per week, max guards per month, and
// max consecutive night shifts.
//
// Rules requiring data not present in the schedule (nurse-patient ratios, break
// compliance, overtime pay) are skipped.
func Validate(schedule Schedule) (*ComplianceReport, error) {
	if For(schedule.Jurisdiction) == nil {
		return nil, fmt.Errorf("unknown jurisdiction: %s", schedule.Jurisdiction)
	}

	parsed, err := parseShifts(schedule.Shifts)
	if err != nil {
		return nil, fmt.Errorf("invalid shift data: %w", err)
	}

	byStaff := groupByStaff(parsed)

	var baseOpts []QueryOption
	if schedule.FacilityScope != "" {
		baseOpts = append(baseOpts, ForScope(schedule.FacilityScope))
	}

	// Resolve night period for this jurisdiction.
	nightStart, nightEnd := resolveNightPeriod(schedule.Jurisdiction)

	report := &ComplianceReport{
		Jurisdiction: schedule.Jurisdiction,
	}

	for staffID, shifts := range byStaff {
		sort.Slice(shifts, func(i, j int) bool {
			return shifts[i].start.Before(shifts[j].start)
		})

		staffType := shifts[0].StaffType

		staffOpts := make([]QueryOption, len(baseOpts), len(baseOpts)+1)
		copy(staffOpts, baseOpts)
		staffOpts = append(staffOpts, ForStaff(staffType))

		rules := EffectiveRules(schedule.Jurisdiction, staffOpts...)

		for _, r := range rules {
			v := r.Current()
			if v == nil {
				continue
			}
			report.ConstraintsChecked++

			violations := dispatchRule(r, v, staffID, shifts, nightStart, nightEnd)
			report.Violations = append(report.Violations, violations...)
		}
	}

	if len(report.Violations) > 0 {
		report.Result = "fail"
	} else {
		report.Result = "pass"
	}

	return report, nil
}

// SwapRequest describes a proposed shift swap for compliance validation.
type SwapRequest struct {
	// Jurisdiction to validate against.
	Jurisdiction Code `json:"jurisdiction"`

	// FacilityScope filters rules by facility type.
	FacilityScope Scope `json:"facility_scope,omitempty"`

	// StaffID is the worker whose schedule will change.
	StaffID string `json:"staff_id"`

	// StaffType is the worker's role.
	StaffType Key `json:"staff_type"`

	// CurrentShifts are all of the worker's existing shift assignments.
	CurrentShifts []Shift `json:"current_shifts"`

	// Remove is the shift being given away (nil if just taking a new shift).
	Remove *Shift `json:"remove,omitempty"`

	// Add is the shift being received (nil if just giving away a shift).
	Add *Shift `json:"add,omitempty"`
}

// ValidateSwap checks whether a proposed shift swap would leave the worker's
// schedule compliant with jurisdiction rules. It simulates the swap (removing
// the old shift, adding the new one) and validates the resulting schedule.
//
// This is the primary integration point for exchange systems: call this before
// accepting a swap to ensure it won't create a compliance violation.
func ValidateSwap(req SwapRequest) (*ComplianceReport, error) {
	if For(req.Jurisdiction) == nil {
		return nil, fmt.Errorf("unknown jurisdiction: %s", req.Jurisdiction)
	}

	// Build the post-swap shift list.
	shifts := make([]Shift, 0, len(req.CurrentShifts)+1)
	for _, s := range req.CurrentShifts {
		// Skip the shift being removed (match by start+end time).
		if req.Remove != nil && s.Start == req.Remove.Start && s.End == req.Remove.End {
			continue
		}
		shifts = append(shifts, s)
	}
	if req.Add != nil {
		add := *req.Add
		add.StaffID = req.StaffID
		add.StaffType = req.StaffType
		shifts = append(shifts, add)
	}

	// Stamp all shifts with the worker's identity.
	for i := range shifts {
		shifts[i].StaffID = req.StaffID
		shifts[i].StaffType = req.StaffType
	}

	return Validate(Schedule{
		Jurisdiction:  req.Jurisdiction,
		FacilityScope: req.FacilityScope,
		Shifts:        shifts,
	})
}

// ---------------------------------------------------------------------------
// Rule dispatch — uses Category + Operator instead of hardcoded key maps.
// New jurisdiction keys are automatically handled by their category.
// ---------------------------------------------------------------------------

func dispatchRule(r *RuleDef, v *RuleValue, staffID string, shifts []parsedShift, nightStart, nightEnd int) []Violation {
	// Boolean/policy rules can't be validated against shift times.
	if r.Operator == OpBool {
		return nil
	}

	switch r.Category {
	case CatWorkHours:
		return dispatchWorkHours(r, v, staffID, shifts, nightStart, nightEnd)
	case CatRest:
		return dispatchRest(r, v, staffID, shifts)
	case CatOnCall:
		return dispatchOnCall(r, v, staffID, shifts)
	case CatNightWork:
		return dispatchNightWork(r, v, staffID, shifts, nightStart, nightEnd)
	default:
		// Categories we can't validate from shift times alone:
		// CatOvertime (pay rates), CatStaffing (ratios), CatBreaks (not in shift data),
		// CatCompensation, CatLeave.
		return nil
	}
}

func dispatchWorkHours(r *RuleDef, v *RuleValue, staffID string, shifts []parsedShift, nightStart, nightEnd int) []Violation {
	switch r.Operator {
	case OpLTE:
		switch v.Per {
		case PerShift:
			return checkMaxShiftHours(r, v, staffID, shifts)
		case PerWeek:
			return checkMaxWeeklyHours(r, v, staffID, shifts)
		case PerDay:
			return checkMaxDailyHours(r, v, staffID, shifts)
		case PerYear:
			// Annual hour limits need a full year of data; skip for now.
			return nil
		}
	}
	return nil
}

func dispatchRest(r *RuleDef, v *RuleValue, staffID string, shifts []parsedShift) []Violation {
	switch r.Operator {
	case OpGTE:
		switch v.Per {
		case PerShift, PerDay:
			return checkMinRestBetweenShifts(r, v, staffID, shifts)
		case PerWeek:
			return checkDaysOffPerWeek(r, v, staffID, shifts)
		case PerOccurrence:
			// Min rest after extended shift or guard.
			return checkMinRestAfterExtended(r, v, staffID, shifts)
		}
	}
	return nil
}

func dispatchOnCall(r *RuleDef, v *RuleValue, staffID string, shifts []parsedShift) []Violation {
	if r.Operator == OpLTE && v.Per == PerMonth {
		return checkMaxGuardsMonthly(r, v, staffID, shifts)
	}
	return nil
}

func dispatchNightWork(r *RuleDef, v *RuleValue, staffID string, shifts []parsedShift, nightStart, nightEnd int) []Violation {
	if r.Operator == OpLTE && v.Unit == Count {
		return checkMaxConsecutiveNights(r, v, staffID, shifts, nightStart, nightEnd)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Night period resolution — reads from jurisdiction rules, not hardcoded.
// ---------------------------------------------------------------------------

// resolveNightPeriod finds the night period start/end hours for a jurisdiction
// by walking the jurisdiction chain. Falls back to 19:00-07:00 if not defined.
func resolveNightPeriod(code Code) (start, end int) {
	start, end = 19, 7 // default fallback
	rules := EffectiveRules(code)
	for _, r := range rules {
		if r.Key == RuleNightPeriodStart {
			if v := r.Current(); v != nil {
				start = int(v.Amount)
			}
		}
		if r.Key == RuleNightPeriodEnd {
			if v := r.Current(); v != nil {
				end = int(v.Amount)
			}
		}
	}
	return start, end
}

// isNightShift returns true if the shift overlaps with the night period
// as defined by the jurisdiction (not hardcoded).
func isNightShift(s parsedShift, nightStart, nightEnd int) bool {
	if s.start.Hour() >= nightStart {
		return true
	}
	// Shift spans midnight and ends in the morning.
	startDay := s.start.Truncate(24 * time.Hour)
	endDay := s.end.Truncate(24 * time.Hour)
	if endDay.After(startDay) && s.end.Hour() <= nightEnd {
		return true
	}
	return false
}

// ---------------------------------------------------------------------------
// Shift parsing and grouping (unchanged).
// ---------------------------------------------------------------------------

type parsedShift struct {
	Shift
	start    time.Time
	end      time.Time
	duration time.Duration
}

func parseShifts(shifts []Shift) ([]parsedShift, error) {
	result := make([]parsedShift, 0, len(shifts))
	for _, s := range shifts {
		start, err := time.Parse(time.RFC3339, s.Start)
		if err != nil {
			start, err = time.Parse("2006-01-02T15:04:05", s.Start)
			if err != nil {
				return nil, fmt.Errorf("cannot parse start time %q for staff %s: %w", s.Start, s.StaffID, err)
			}
		}
		end, err := time.Parse(time.RFC3339, s.End)
		if err != nil {
			end, err = time.Parse("2006-01-02T15:04:05", s.End)
			if err != nil {
				return nil, fmt.Errorf("cannot parse end time %q for staff %s: %w", s.End, s.StaffID, err)
			}
		}
		if !end.After(start) {
			return nil, fmt.Errorf("shift end %q is not after start %q for staff %s", s.End, s.Start, s.StaffID)
		}
		result = append(result, parsedShift{
			Shift:    s,
			start:    start,
			end:      end,
			duration: end.Sub(start),
		})
	}
	return result, nil
}

func groupByStaff(shifts []parsedShift) map[string][]parsedShift {
	groups := make(map[string][]parsedShift)
	for _, s := range shifts {
		groups[s.StaffID] = append(groups[s.StaffID], s)
	}
	return groups
}

// ---------------------------------------------------------------------------
// Check functions.
// ---------------------------------------------------------------------------

func checkMaxShiftHours(r *RuleDef, v *RuleValue, staffID string, shifts []parsedShift) []Violation {
	limit := v.Amount
	var violations []Violation
	for _, s := range shifts {
		hours := s.duration.Hours()
		if hours > limit {
			violations = append(violations, Violation{
				RuleKey:  r.Key,
				RuleName: r.Name,
				Severity: r.Enforcement,
				StaffID:  staffID,
				Message:  fmt.Sprintf("Shift on %s is %.1f hours, maximum is %.0f hours", s.start.Format("2006-01-02"), hours, limit),
				Citation: r.Source.Citation(),
				Actual:   math.Round(hours*10) / 10,
				Limit:    limit,
			})
		}
	}
	return violations
}

func checkMaxDailyHours(r *RuleDef, v *RuleValue, staffID string, shifts []parsedShift) []Violation {
	limit := v.Amount
	var violations []Violation

	// Group hours by calendar date.
	daily := make(map[string]float64)
	for _, s := range shifts {
		day := s.start.Format("2006-01-02")
		daily[day] += s.duration.Hours()
	}

	for day, hours := range daily {
		if hours > limit {
			violations = append(violations, Violation{
				RuleKey:  r.Key,
				RuleName: r.Name,
				Severity: r.Enforcement,
				StaffID:  staffID,
				Message:  fmt.Sprintf("%.1f hours on %s, maximum is %.0f hours/day", hours, day, limit),
				Citation: r.Source.Citation(),
				Actual:   math.Round(hours*10) / 10,
				Limit:    limit,
			})
		}
	}
	return violations
}

func checkMaxWeeklyHours(r *RuleDef, v *RuleValue, staffID string, shifts []parsedShift) []Violation {
	if len(shifts) == 0 {
		return nil
	}

	limit := v.Amount

	windowDays := 7
	if v.Averaged != nil {
		windowDays = averagingDays(v.Averaged)
	}

	minDate := shifts[0].start
	maxDate := shifts[0].end
	for _, s := range shifts {
		if s.start.Before(minDate) {
			minDate = s.start
		}
		if s.end.After(maxDate) {
			maxDate = s.end
		}
	}

	windowDur := time.Duration(windowDays) * 24 * time.Hour
	windowStart := minDate.Truncate(24 * time.Hour)
	windowEnd := windowStart.Add(windowDur)

	var violations []Violation

	for windowStart.Before(maxDate) {
		var totalHours float64
		for _, s := range shifts {
			overlapStart := s.start
			if overlapStart.Before(windowStart) {
				overlapStart = windowStart
			}
			overlapEnd := s.end
			if overlapEnd.After(windowEnd) {
				overlapEnd = windowEnd
			}
			if overlapEnd.After(overlapStart) {
				totalHours += overlapEnd.Sub(overlapStart).Hours()
			}
		}

		weeks := float64(windowDays) / 7.0
		if weeks == 0 {
			break
		}
		avgWeeklyHours := totalHours / weeks

		if avgWeeklyHours > limit {
			violations = append(violations, Violation{
				RuleKey:  r.Key,
				RuleName: r.Name,
				Severity: r.Enforcement,
				StaffID:  staffID,
				Message: fmt.Sprintf("%.1f hours/week (averaged over %d days starting %s), maximum is %.0f hours/week",
					avgWeeklyHours, windowDays, windowStart.Format("2006-01-02"), limit),
				Citation: r.Source.Citation(),
				Actual:   math.Round(avgWeeklyHours*10) / 10,
				Limit:    limit,
			})
		}

		windowStart = windowStart.Add(7 * 24 * time.Hour)
		windowEnd = windowStart.Add(windowDur)
	}

	return violations
}

func checkMinRestBetweenShifts(r *RuleDef, v *RuleValue, staffID string, shifts []parsedShift) []Violation {
	minRest := v.Amount
	var violations []Violation

	for i := 1; i < len(shifts); i++ {
		gap := shifts[i].start.Sub(shifts[i-1].end).Hours()
		if gap < 0 {
			violations = append(violations, Violation{
				RuleKey:  r.Key,
				RuleName: r.Name,
				Severity: r.Enforcement,
				StaffID:  staffID,
				Message: fmt.Sprintf("Overlapping shifts: shift ending %s overlaps with shift starting %s",
					shifts[i-1].end.Format("2006-01-02 15:04"), shifts[i].start.Format("2006-01-02 15:04")),
				Citation: r.Source.Citation(),
				Actual:   0,
				Limit:    minRest,
			})
		} else if gap < minRest {
			violations = append(violations, Violation{
				RuleKey:  r.Key,
				RuleName: r.Name,
				Severity: r.Enforcement,
				StaffID:  staffID,
				Message: fmt.Sprintf("%.1f hours rest between shifts on %s and %s, minimum is %.0f hours",
					gap, shifts[i-1].end.Format("2006-01-02 15:04"), shifts[i].start.Format("2006-01-02 15:04"), minRest),
				Citation: r.Source.Citation(),
				Actual:   math.Round(gap*10) / 10,
				Limit:    minRest,
			})
		}
	}

	return violations
}

func checkMinRestAfterExtended(r *RuleDef, v *RuleValue, staffID string, shifts []parsedShift) []Violation {
	minRest := v.Amount
	var violations []Violation

	for i, s := range shifts {
		if s.duration.Hours() < 24 {
			continue
		}
		if i+1 < len(shifts) {
			gap := shifts[i+1].start.Sub(s.end).Hours()
			if gap < minRest {
				actual := gap
				if actual < 0 {
					actual = 0
				}
				violations = append(violations, Violation{
					RuleKey:  r.Key,
					RuleName: r.Name,
					Severity: r.Enforcement,
					StaffID:  staffID,
					Message: fmt.Sprintf("%.1f hours rest after extended shift (%.0fh) ending %s, minimum is %.0f hours",
						actual, s.duration.Hours(), s.end.Format("2006-01-02 15:04"), minRest),
					Citation: r.Source.Citation(),
					Actual:   math.Round(actual*10) / 10,
					Limit:    minRest,
				})
			}
		}
	}

	return violations
}

func checkDaysOffPerWeek(r *RuleDef, v *RuleValue, staffID string, shifts []parsedShift) []Violation {
	if len(shifts) == 0 {
		return nil
	}

	// If the value is in hours (e.g., EU "35 hours weekly rest"), convert to days.
	// Otherwise treat as number of days off required.
	requiredDaysOff := v.Amount
	if v.Unit == Hours {
		// Hours-based weekly rest: check the longest continuous gap per week.
		return checkMinWeeklyRestHours(r, v, staffID, shifts)
	}

	var violations []Violation

	minDate := shifts[0].start.Truncate(24 * time.Hour)
	maxDate := shifts[len(shifts)-1].end

	for weekStart := minDate; weekStart.Before(maxDate); weekStart = weekStart.Add(7 * 24 * time.Hour) {
		weekEnd := weekStart.Add(7 * 24 * time.Hour)

		daysWorked := make(map[string]bool)
		for _, s := range shifts {
			if s.end.After(weekStart) && s.start.Before(weekEnd) {
				day := s.start.Truncate(24 * time.Hour)
				endDay := s.end.Truncate(24 * time.Hour)
				if s.end.Equal(endDay) && !endDay.Equal(s.start.Truncate(24*time.Hour)) {
					endDay = endDay.Add(-24 * time.Hour)
				}
				for !day.After(endDay) {
					if !day.Before(weekStart) && day.Before(weekEnd) {
						daysWorked[day.Format("2006-01-02")] = true
					}
					day = day.Add(24 * time.Hour)
				}
			}
		}

		daysOff := 7 - len(daysWorked)
		if float64(daysOff) < requiredDaysOff {
			violations = append(violations, Violation{
				RuleKey:  r.Key,
				RuleName: r.Name,
				Severity: r.Enforcement,
				StaffID:  staffID,
				Message: fmt.Sprintf("%d days off in week starting %s, minimum is %.0f",
					daysOff, weekStart.Format("2006-01-02"), requiredDaysOff),
				Citation: r.Source.Citation(),
				Actual:   float64(daysOff),
				Limit:    requiredDaysOff,
			})
		}
	}

	return violations
}

// checkMinWeeklyRestHours validates hour-based weekly rest (e.g., EU 35 hours,
// Spain 36 hours). Finds the longest continuous gap between shifts per week.
func checkMinWeeklyRestHours(r *RuleDef, v *RuleValue, staffID string, shifts []parsedShift) []Violation {
	if len(shifts) == 0 {
		return nil
	}

	requiredHours := v.Amount
	var violations []Violation

	windowDays := 7
	if v.Averaged != nil {
		windowDays = averagingDays(v.Averaged)
	}

	minDate := shifts[0].start.Truncate(24 * time.Hour)
	maxDate := shifts[len(shifts)-1].end

	for weekStart := minDate; weekStart.Before(maxDate); weekStart = weekStart.Add(7 * 24 * time.Hour) {
		weekEnd := weekStart.Add(time.Duration(windowDays) * 24 * time.Hour)

		// Collect shifts in this window, sorted.
		var windowShifts []parsedShift
		for _, s := range shifts {
			if s.end.After(weekStart) && s.start.Before(weekEnd) {
				windowShifts = append(windowShifts, s)
			}
		}

		// Find the longest gap.
		longestGap := 0.0
		cursor := weekStart
		for _, s := range windowShifts {
			effectiveStart := s.start
			if effectiveStart.Before(weekStart) {
				effectiveStart = weekStart
			}
			gap := effectiveStart.Sub(cursor).Hours()
			if gap > longestGap {
				longestGap = gap
			}
			effectiveEnd := s.end
			if effectiveEnd.After(cursor) {
				cursor = effectiveEnd
			}
		}
		// Gap after last shift to end of window.
		finalGap := weekEnd.Sub(cursor).Hours()
		if finalGap > longestGap {
			longestGap = finalGap
		}

		if longestGap < requiredHours {
			violations = append(violations, Violation{
				RuleKey:  r.Key,
				RuleName: r.Name,
				Severity: r.Enforcement,
				StaffID:  staffID,
				Message: fmt.Sprintf("Longest continuous rest in week starting %s is %.1f hours, minimum is %.0f hours",
					weekStart.Format("2006-01-02"), longestGap, requiredHours),
				Citation: r.Source.Citation(),
				Actual:   math.Round(longestGap*10) / 10,
				Limit:    requiredHours,
			})
		}
	}

	return violations
}

func checkMaxGuardsMonthly(r *RuleDef, v *RuleValue, staffID string, shifts []parsedShift) []Violation {
	limit := int(v.Amount)
	var violations []Violation

	monthly := make(map[string]int)
	for _, s := range shifts {
		if s.OnCall {
			key := s.start.Format("2006-01")
			monthly[key]++
		}
	}

	for month, count := range monthly {
		if count > limit {
			violations = append(violations, Violation{
				RuleKey:  r.Key,
				RuleName: r.Name,
				Severity: r.Enforcement,
				StaffID:  staffID,
				Message:  fmt.Sprintf("%d guard duties in %s, maximum is %d", count, month, limit),
				Citation: r.Source.Citation(),
				Actual:   float64(count),
				Limit:    float64(limit),
			})
		}
	}

	return violations
}

func checkMaxConsecutiveNights(r *RuleDef, v *RuleValue, staffID string, shifts []parsedShift, nightStart, nightEnd int) []Violation {
	limit := int(v.Amount)
	var violations []Violation

	consecutive := 0
	var streakStart time.Time

	for _, s := range shifts {
		if isNightShift(s, nightStart, nightEnd) {
			if consecutive == 0 {
				streakStart = s.start
			}
			consecutive++
			if consecutive > limit {
				violations = append(violations, Violation{
					RuleKey:  r.Key,
					RuleName: r.Name,
					Severity: r.Enforcement,
					StaffID:  staffID,
					Message: fmt.Sprintf("%d consecutive night shifts starting %s, maximum is %d",
						consecutive, streakStart.Format("2006-01-02"), limit),
					Citation: r.Source.Citation(),
					Actual:   float64(consecutive),
					Limit:    float64(limit),
				})
			}
		} else {
			consecutive = 0
		}
	}

	return violations
}

// ---------------------------------------------------------------------------
// Helpers.
// ---------------------------------------------------------------------------

// averagingDays converts an AveragingPeriod to days using proper calendar math.
func averagingDays(a *AveragingPeriod) int {
	switch a.Unit {
	case PeriodDays:
		return a.Count
	case PeriodWeeks:
		return a.Count * 7
	case PeriodMonths:
		// Use 30.44 days/month (365.25/12) for better accuracy than flat 30.
		return int(math.Round(float64(a.Count) * 30.44))
	}
	return 7
}
