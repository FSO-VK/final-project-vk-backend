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
	// dosage is an amount of medication intake per one take.
	dosage dosage
	status Status
	// schedule contains the schedule of the plan.
	schedule schedule
	// courseStart is a date and time when the planned course starts.
	courseStart courseStart
	// courseEnd is a date and time when the planned course ends.
	courseEnd courseEnd
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
	dosage dosage,
	schedule schedule,
	start courseStart,
	end courseEnd,
	condition string,
	createdAt time.Time,
	updatedAt time.Time,
) (*Plan, error) {
	if end.ToTime().Before(start.ToTime()) {
		return nil, ErrCourseRange
	}

	return &Plan{
		id:           id,
		medicationID: medicationID,
		dosage:       dosage,
		schedule:     schedule,
		status:       StatusActive,
		courseStart:  start,
		courseEnd:    end,
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
	newStart courseStart,
) (*Plan, error) {
	if p.status != StatusActive {
		return nil, ErrFinishedPlan
	}

	p.schedule = newSchedule
	p.courseStart = newStart
	return p, nil
}

// Schedule returns the schedule of the plan in range [from, to].
// If there is no records in the range, it returns nil.
func (p *Plan) Schedule(from, to time.Time) []time.Time {
	if from.After(to) {
		return nil
	}

	if to.After(p.courseEnd.ToTime()) {
		to = p.courseEnd.ToTime()
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
