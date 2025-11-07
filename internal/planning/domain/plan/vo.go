package plan

import (
	"errors"
	"fmt"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/pkg/validation"
	"github.com/robfig/cron/v3"
)

var (
	// ErrInvalidDosage means that the dosage is invalid.
	ErrInvalidDosage = errors.New("invalid dosage")
	// ErrInvalidSchedule means that the schedule expression format is invalid.
	ErrInvalidSchedule = errors.New("invalid schedule expression format")
	// ErrInvalidCourseStart means that something is wrong with the course start time.
	ErrInvalidCourseStart = errors.New("invalid course start time")
	// ErrInvalidCourseEnd means that course end time is in the past.
	ErrInvalidCourseEnd = errors.New("course ends in the past")
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

// schedule is a VO that describes the schedule in unix cron format.
type schedule struct {
	s cron.Schedule
}

// NewSchedule returns valid schedule. Accepts cronExpr is a cron expression.
// It requires 5 entries representing:
// minute, hour, day of month, month and day of week, in that
// order. It returns a descriptive error if the spec is not valid.
// For more info see https://en.wikipedia.org/wiki/Cron.
//
// It accepts
//   - Standard crontab specs, e.g. "* * * * ?"
//   - Descriptors, e.g. "@midnight", "@every 1h30m"
func NewSchedule(cronExpr string) (schedule, error) {
	s, err := cron.ParseStandard(cronExpr)
	if err != nil {
		return schedule{}, fmt.Errorf("%w: %w", ErrInvalidSchedule, err)
	}

	return schedule{
		s: s,
	}, nil
}

// Next returns the next scheduled time after the given time.
// If there is no next time, it returns the zero time.
func (s *schedule) Next(from time.Time) time.Time {
	return s.s.Next(from)
}

// courseStart is a VO that describes the start time of the course.
type courseStart time.Time

// NewCourseStart creates validated courseStart.
func NewCourseStart(t time.Time) (courseStart, error) {
	if t.Before(time.Now()) {
		return courseStart(time.Time{}), ErrInvalidCourseStart
	}
	return courseStart(t), nil
}

// ToTime returns the time.Time representation of the course start.
func (c courseStart) ToTime() time.Time {
	return time.Time(c)
}

// courseEnd is a VO that describes the end time of the course.
// It shouldn't be in the past.
type courseEnd time.Time

// NewCourseEnd creates validated courseEnd.
func NewCourseEnd(t time.Time) (courseEnd, error) {
	if t.Before(time.Now()) {
		return courseEnd(time.Time{}), ErrInvalidCourseEnd
	}
	return courseEnd(t), nil
}

// ToTime returns the time.Time representation of the course end.
func (c courseEnd) ToTime() time.Time {
	return time.Time(c)
}
