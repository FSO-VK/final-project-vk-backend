package plan

import (
	"errors"
	"fmt"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/pkg/validation"
	"github.com/teambition/rrule-go"
)

var (
	// ErrInvalidDosage means that the dosage is invalid.
	ErrInvalidDosage = errors.New("invalid dosage")
	// ErrInvalidSchedule means that the schedule expression format is invalid.
	ErrInvalidSchedule = errors.New("invalid schedule expression format")
)

// dosage is a VO representing the dosage of planned medication.
type dosage struct {
	value float64
	unit  string
}

// Units is a list of possible units for dosage.
//
//nolint:gochecknoglobals
var Units = []string{"мг", "шт", "мл"}

// NewDosage creates validated dosage.
func NewDosage(value float64, unit string) (dosage, error) {
	err := errors.Join(
		validation.Positive(value),
		validation.Required(unit),
	)
	if err != nil {
		return dosage{}, fmt.Errorf("%w: %w", ErrInvalidDosage, err)
	}

	var u string
	for _, u = range Units {
		if u == unit {
			break
		}
	}
	if u == "" {
		return dosage{}, ErrInvalidDosage
	}

	return dosage{
		value: value,
		unit:  u,
	}, nil
}

// Status is a VO representing the status of the plan.
type Status uint

// Enum of statuses.
const (
	StatusDraft Status = iota
	StatusActive
	StatusFinished
)

type schedule struct {
	start time.Time
	end   time.Time

	// specific RFC5545 fields to contain complex schedule
	rules []*rrule.RRule
}

func NewSchedule(start, end time.Time, rules []*rrule.RRule) (schedule, error) {
	if end.Before(start) {
		return schedule{}, ErrInvalidSchedule
	}

	r := make([]*rrule.RRule, 0, len(rules))
	for _, rule := range rules {
		if rule == nil {
			continue
		}
		// rule is limited by range
		rule.DTStart(start)
		rule.Until(end)
		r = append(r, rule)
	}

	return schedule{
		start: start,
		end:   end,
		rules: r,
	}, nil
}

// Next returns the next scheduled time after the given time.
// If there is no next time, it returns the zero time.
func (s *schedule) Next(from time.Time) time.Time {
	var t time.Time
	for _, rules := range s.rules {
		next := rules.After(from, false)
		t = time.Unix(min(t.Unix(), next.Unix()), 0)
	}
	return t
}