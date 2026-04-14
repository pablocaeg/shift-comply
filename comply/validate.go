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
// min rest after extended shifts, days off per week, and max guards per month.
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

	report := &ComplianceReport{
		Jurisdiction: schedule.Jurisdiction,
	}

	for staffID, shifts := range byStaff {
		sort.Slice(shifts, func(i, j int) bool {
			return shifts[i].start.Before(shifts[j].start)
		})

		staffType := shifts[0].StaffType

		// Build a new slice for each staff member to avoid mutating baseOpts.
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

			violations := checkRule(r, v, staffID, shifts)
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

// maxShiftKeys contains all rule keys that represent max shift hour limits.
var maxShiftKeys = map[Key]bool{
	RuleMaxShiftHours:              true,
	"max-shift-hours-non-resident": true,
}

var maxWeeklyKeys = map[Key]bool{
	RuleMaxWeeklyHours:         true,
	RuleMaxCombinedWeeklyHours: true,
	RuleMaxOrdinaryWeeklyHours: true,
}

var minRestKeys = map[Key]bool{
	RuleMinRestBetweenShifts:         true,
	"es-mir-min-rest-between-shifts": true,
}

var daysOffKeys = map[Key]bool{
	RuleDaysOffPerWeek: true,
	RuleMinDayOfRest:   true,
}

func checkRule(r *RuleDef, v *RuleValue, staffID string, shifts []parsedShift) []Violation {
	if maxShiftKeys[r.Key] {
		return checkMaxShiftHours(r, v, staffID, shifts)
	}
	if maxWeeklyKeys[r.Key] {
		return checkMaxWeeklyHours(r, v, staffID, shifts)
	}
	if minRestKeys[r.Key] {
		return checkMinRestBetweenShifts(r, v, staffID, shifts)
	}
	if daysOffKeys[r.Key] {
		return checkDaysOffPerWeek(r, v, staffID, shifts)
	}
	switch r.Key {
	case RuleMinRestAfterExtended:
		return checkMinRestAfterExtended(r, v, staffID, shifts)
	case RuleMaxGuardsMonthly:
		return checkMaxGuardsMonthly(r, v, staffID, shifts)
	case RuleMaxConsecutiveNights:
		return checkMaxConsecutiveNights(r, v, staffID, shifts)
	case RuleMinRestAfterGuard:
		return checkMinRestAfterExtended(r, v, staffID, shifts)
	default:
		return nil
	}
}

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

func checkMaxWeeklyHours(r *RuleDef, v *RuleValue, staffID string, shifts []parsedShift) []Violation {
	if len(shifts) == 0 {
		return nil
	}

	limit := v.Amount

	// Determine averaging window
	windowDays := 7
	if v.Averaged != nil {
		switch v.Averaged.Unit {
		case PeriodWeeks:
			windowDays = v.Averaged.Count * 7
		case PeriodDays:
			windowDays = v.Averaged.Count
		case PeriodMonths:
			windowDays = v.Averaged.Count * 30
		}
	}

	// Find the date range across all shifts
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
			// Overlapping shifts: flag as zero rest
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

	requiredDaysOff := v.Amount
	var violations []Violation

	minDate := shifts[0].start.Truncate(24 * time.Hour)
	maxDate := shifts[len(shifts)-1].end

	for weekStart := minDate; weekStart.Before(maxDate); weekStart = weekStart.Add(7 * 24 * time.Hour) {
		weekEnd := weekStart.Add(7 * 24 * time.Hour)

		daysWorked := make(map[string]bool)
		for _, s := range shifts {
			if s.end.After(weekStart) && s.start.Before(weekEnd) {
				// Use shift start date and end date, but only count a day
				// as worked if the shift was active during working hours of that day.
				// Simpler approach: count each calendar date where the shift has
				// at least 1 hour of overlap.
				day := s.start.Truncate(24 * time.Hour)
				endDay := s.end.Truncate(24 * time.Hour)
				// If shift ends exactly at midnight, don't count that next day
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

func checkMaxConsecutiveNights(r *RuleDef, v *RuleValue, staffID string, shifts []parsedShift) []Violation {
	limit := int(v.Amount)
	var violations []Violation

	consecutive := 0
	var streakStart time.Time

	for _, s := range shifts {
		if isNightShift(s) {
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

// isNightShift returns true if the shift overlaps with the night period.
// A shift is considered a night shift if it starts at 19:00 or later,
// or if it spans past midnight (end date is after start date).
func isNightShift(s parsedShift) bool {
	if s.start.Hour() >= 19 {
		return true
	}
	// Shift spans midnight: start and end are on different calendar days
	startDay := s.start.Truncate(24 * time.Hour)
	endDay := s.end.Truncate(24 * time.Hour)
	if endDay.After(startDay) && s.end.Hour() <= 12 {
		return true
	}
	return false
}
