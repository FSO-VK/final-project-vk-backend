// Package plan is subdomain for planning domain.
package plan

import (
	"errors"
	"fmt"
	"time"

	intake "github.com/FSO-VK/final-project-vk-backend/internal/planning/domain/record"
	"github.com/google/uuid"
)

var (
	// ErrGenIntakeRecord tells that something bad happen while generating records.
	ErrGenIntakeRecord = errors.New("cannot generate intake records")
	// ErrCourseRange tells that course end time is before course start time.
	ErrCourseRange = errors.New("course ends before it starts")
	// ErrFinishedPlan tells that plan is already finished and can;t be mutated.
	ErrFinishedPlan = errors.New("can't modify finished plan")
)

// Plan is an aggregate that represents a plan for medication intake.
type Plan struct {
	id           uuid.UUID
	medicationID uuid.UUID
	userID       uuid.UUID
	// dosage is an amount of medication intake per one take.
	dosage dosage
	status Status
	// schedule contains the schedule of the plan.
	schedule schedule
	// condition is a description of the condition
	// under which the medication should be taken.
	condition string
	createdAt time.Time
	updatedAt time.Time
}

// NewPlan creates validated plan.
func NewPlan(
	id uuid.UUID,
	medicationID uuid.UUID,
	userID uuid.UUID,
	dosage dosage,
	schedule schedule,
	condition string,
	createdAt time.Time,
	updatedAt time.Time,
) (*Plan, error) {
	return &Plan{
		id:           id,
		medicationID: medicationID,
		userID:       userID,
		dosage:       dosage,
		schedule:     schedule,
		status:       StatusActive,
		condition:    condition,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}, nil
}

// ChangeDosage executes business logic for changing the dosage of the plan.
func (p *Plan) ChangeDosage(d dosage) (*Plan, error) {
	if p.status != StatusActive {
		return nil, ErrFinishedPlan
	}

	p.dosage = d
	return p, nil
}

// ChangeSchedule executes business logic for changing the schedule of the plan.
func (p *Plan) ChangeSchedule(
	newSchedule schedule,
) (*Plan, error) {
	if p.status != StatusActive {
		return nil, ErrFinishedPlan
	}

	p.schedule = newSchedule
	return p, nil
}

// Deactivate executes business logic for finishing the plan (soft deletion).
func (p *Plan) Deactivate() (*Plan, error) {
	p.status = StatusFinished
	return p, nil
}

// Schedule returns the schedule of the plan in range [from, to].
// If there is no records in the range, it returns nil.
func (p *Plan) Schedule(from, to time.Time) []time.Time {
	if from.After(to) {
		return nil
	}

	var schedule []time.Time
	for t := p.schedule.Next(from); t.Before(to) && !t.IsZero(); t = p.schedule.Next(t) {
		schedule = append(schedule, t)
	}
	return schedule
}

// GenerateIntakeRecords is a factory for intake records related to the plan.
func (p *Plan) GenerateIntakeRecords(from, to time.Time) ([]*intake.IntakeRecord, error) {
	records := make([]*intake.IntakeRecord, 0)

	times := p.Schedule(from, to)
	for _, t := range times {
		record, err := intake.NewIntakeRecord(
			uuid.New(),
			p.id,
			t,
			time.Now(),
			time.Now(),
		)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrGenIntakeRecord, err)
		}
		records = append(records, record)
	}

	return records, nil
}

// ID returns the ID of the plan.
func (p *Plan) ID() uuid.UUID {
	return p.id
}

// MedicationID returns the medication ID of the plan.
func (p *Plan) MedicationID() uuid.UUID {
	return p.medicationID
}

// UserID returns the id of user (owner of the plan).
func (p *Plan) UserID() uuid.UUID {
	return p.userID
}

// Dosage returns the dosage per intake.
// It returns amount and unit.
func (p *Plan) Dosage() (float64, string) {
	return p.dosage.value, p.dosage.unit
}

// IsActive tells whether the plan is active.
func (p *Plan) IsActive() bool {
	return p.status == StatusActive
}

// Condition returns the intake condition of the plan.
func (p *Plan) Condition() string {
	return p.condition
}

// CourseStart returns the start of the plan.
func (p *Plan) CourseStart() time.Time {
	return p.schedule.start
}

// CourseEnd returns the end of the plan.
func (p *Plan) CourseEnd() time.Time {
	return p.schedule.end
}

// ScheduleIcal returns the recurrence rules in iCalendar RFC 5545 format.
// Each rule string defines a recurrence pattern for the schedule using the RRULE property.
// The format specifies how the event repeats over time (frequency, interval, by day, etc.).
// Multiple rules can be returned for complex schedules with different patterns.
// RFC 5545 Specification: https://tools.ietf.org/html/rfc5545#section-3.3.10
func (p *Plan) ScheduleIcal() []string {
	rules := make([]string, 0, len(p.schedule.rules))
	for _, rule := range p.schedule.rules {
		rules = append(rules, rule.String())
	}
	return rules
}

// Status returns the status of the plan.
func (r *Plan) Status() Status {
	return r.status
}
