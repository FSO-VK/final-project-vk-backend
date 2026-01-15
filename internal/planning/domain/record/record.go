// Package record is subdomain for planning domain.
package record

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ErrRecordOutdated tells that record is outdated and can't be rescheduled.
var ErrRecordOutdated = errors.New("cannot reschedule outdated record")

// IntakeRecord is an aggregate that represents a record for medication intake.
type IntakeRecord struct {
	id        uuid.UUID
	planID    uuid.UUID
	status    Status
	plannedAt time.Time
	takenAt   time.Time
	createdAt time.Time
	updatedAt time.Time
}

// NewIntakeRecord creates validated IntakeRecord.
func NewIntakeRecord(
	id uuid.UUID,
	planID uuid.UUID,
	plannedAt time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) (*IntakeRecord, error) {
	return &IntakeRecord{
		id:        id,
		planID:    planID,
		status:    StatusDraft,
		plannedAt: plannedAt,
		takenAt:   time.Time{}, // zero value
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

// MarkTaken executes business logic for marking the record as taken.
func (r *IntakeRecord) MarkTaken(t time.Time) *IntakeRecord {
	r.status = StatusTaken
	r.takenAt = t
	return r
}

// MarkMissed executes business logic for marking the record as missed.
func (r *IntakeRecord) MarkMissed() *IntakeRecord {
	r.status = StatusMissed
	return r
}

// Cancel executes business logic for marking the record as draft if user want to cancel taking it.
func (r *IntakeRecord) Cancel() *IntakeRecord {
	r.status = StatusDraft
	return r
}

// Reschedule executes business logic for rescheduling the future record.
func (r *IntakeRecord) Reschedule(newPlannedTime time.Time) (*IntakeRecord, error) {
	if r.status != StatusDraft {
		return nil, ErrRecordOutdated
	}
	r.plannedAt = newPlannedTime
	return r, nil
}

// IsTaken returns the status of the record.
func (r *IntakeRecord) IsTaken() bool {
	return r.status == StatusTaken
}

// PlannedTime returns the planned time of the intake.
func (r *IntakeRecord) PlannedTime() time.Time {
	return r.plannedAt
}

// ID returns the ID of the intake record.
func (r *IntakeRecord) ID() uuid.UUID {
	return r.id
}

// PlanID returns the plan ID associated with the intake record.
func (r *IntakeRecord) PlanID() uuid.UUID {
	return r.planID
}

// TakenAt returns the time the record was taken.
func (r *IntakeRecord) TakenAt() time.Time {
	return r.takenAt
}

// Status returns the status of the record.
func (r *IntakeRecord) Status() Status {
	return r.status
}
