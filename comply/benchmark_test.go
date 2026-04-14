package comply_test

import (
	"fmt"
	"testing"

	"github.com/pablocaeg/shift-comply/comply"
	_ "github.com/pablocaeg/shift-comply/jurisdictions"
)

func BenchmarkFor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		comply.For(comply.USCA)
	}
}

func BenchmarkEffectiveRules_NoFilters(b *testing.B) {
	for i := 0; i < b.N; i++ {
		comply.EffectiveRules(comply.USCA)
	}
}

func BenchmarkEffectiveRules_StaffFilter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		comply.EffectiveRules(comply.USCA, comply.ForStaff(comply.StaffNurseRN))
	}
}

func BenchmarkEffectiveRules_ThreeLevelChain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		comply.EffectiveRules(comply.ESCT, comply.ForStaff(comply.StaffStatutory))
	}
}

func BenchmarkGenerateConstraints(b *testing.B) {
	for i := 0; i < b.N; i++ {
		comply.GenerateConstraints(comply.USCA, comply.ForStaff(comply.StaffNurseRN))
	}
}

func BenchmarkCompare(b *testing.B) {
	for i := 0; i < b.N; i++ {
		comply.Compare(comply.USCA, comply.ES)
	}
}

func BenchmarkValidate_SmallSchedule(b *testing.B) {
	schedule := comply.Schedule{
		Jurisdiction: comply.ES,
		Shifts: []comply.Shift{
			{StaffID: "doc-1", StaffType: comply.StaffStatutory, Start: "2025-03-10T08:00:00", End: "2025-03-10T15:00:00"},
			{StaffID: "doc-1", StaffType: comply.StaffStatutory, Start: "2025-03-11T08:00:00", End: "2025-03-11T15:00:00"},
			{StaffID: "doc-1", StaffType: comply.StaffStatutory, Start: "2025-03-12T08:00:00", End: "2025-03-12T15:00:00"},
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = comply.Validate(schedule)
	}
}

func BenchmarkValidate_LargeSchedule(b *testing.B) {
	shifts := make([]comply.Shift, 0, 100)
	for i := 0; i < 20; i++ {
		for d := 1; d <= 5; d++ {
			day := 1 + (i/4)*7 + d
			if day > 28 {
				continue
			}
			shifts = append(shifts, comply.Shift{
				StaffID:   fmt.Sprintf("staff-%d", i),
				StaffType: comply.StaffStatutory,
				Start:     fmt.Sprintf("2025-03-%02dT08:00:00", day),
				End:       fmt.Sprintf("2025-03-%02dT16:00:00", day),
			})
		}
	}
	schedule := comply.Schedule{
		Jurisdiction: comply.ES,
		Shifts:       shifts,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = comply.Validate(schedule)
	}
}
