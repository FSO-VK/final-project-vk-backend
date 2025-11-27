package plan

// Testing as a white box because VO are unexported.

import (
	"testing"
	"time"

	"github.com/teambition/rrule-go"
)


func Test_schedule_Next(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		start time.Time
		end   time.Time
		rules []*rrule.RRule
		// Named input parameters for target function.
		from time.Time
		want time.Time
	}{
		{
			name: "Should return next occurrence",
			start: time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
			end: time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC),
			rules: func() []*rrule.RRule {
				var rule1, rule2 *rrule.RRule
				rule1, err := rrule.NewRRule(rrule.ROption{
					Freq: rrule.DAILY,
					Byhour: []int{9, 19},
				})
				rule2, err = rrule.NewRRule(rrule.ROption{
					// friday is a January 5, 2024
					Byweekday: []rrule.Weekday{rrule.MO, rrule.WE, rrule.FR},
					Freq:      rrule.WEEKLY,
					Byhour:    []int{15},
				})
				if err != nil {
					t.Fatalf("arrange failed: %v", err)
				}
				return []*rrule.RRule{rule1, rule2}
			}(),
			from: time.Date(2024, 1, 5, 10, 0, 0, 0, time.UTC),
			want: time.Date(2024, 1, 5, 15, 0, 0, 0, time.UTC),
		},
		{
			name: "Should return zero time if there is no next occurrence",
			start: time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
			end: time.Date(2024, 1, 3, 9, 0, 0, 0, time.UTC),
			rules: func() []*rrule.RRule {
				rule, err := rrule.NewRRule(rrule.ROption{
					Freq:   rrule.DAILY,
					Byhour: []int{9},
				})
				if err != nil {
					t.Fatalf("arrange failed: %v", err)
				}
				return []*rrule.RRule{rule}
			}(),
			from: time.Date(2024, 1, 3, 10, 0, 0, 0, time.UTC),
			want: time.Time{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s, err := NewSchedule(tt.start, tt.end, tt.rules)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got := s.Next(tt.from)
			if !tt.want.Equal(got) {
				t.Errorf("Next() = %v, want %v", got, tt.want)
			}
		})
	}
}