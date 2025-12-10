// Package application is a package for application logic of the planning service.
package application

import (
	"context"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/planning/application/medication"
	"github.com/FSO-VK/final-project-vk-backend/internal/planning/domain/plan"
	"github.com/FSO-VK/final-project-vk-backend/internal/planning/domain/record"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/google/uuid"
)

// ShowSchedule is an interface for getting a notification.
type ShowSchedule interface {
	Execute(
		ctx context.Context,
		cmd *ShowScheduleCommand,
	) (*ShowScheduleResponse, error)
}

// ShowScheduleService is a service for creating a subscription.
type ShowScheduleService struct {
	planningRepo       plan.Repository
	recordsRepo        record.Repository
	validator          validator.Validator
	medicationProvider medication.MedicationService
	// createdShift is the offset from 00:00 when records are generated.
	// At 00:00 + createdShift, all records for that day are created. (basically 24h - today creating for the next day)
	createdShift time.Duration
}

// NewShowScheduleService returns a new ShowScheduleService.
func NewShowScheduleService(
	planningRepo plan.Repository,
	recordsRepo record.Repository,
	valid validator.Validator,
	medicationProvider medication.MedicationService,
	createdShift time.Duration,
) *ShowScheduleService {
	return &ShowScheduleService{
		planningRepo:       planningRepo,
		recordsRepo:        recordsRepo,
		validator:          valid,
		medicationProvider: medicationProvider,
		createdShift:       createdShift,
	}
}

// ShowScheduleCommand is a request to get a plan.
type ShowScheduleCommand struct {
	UserID    string `validate:"required,uuid"`
	StartDate string `validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	EndDate   string `validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
}

// ScheduleTime is an item of the schedules list to be returned.
type ScheduleTime struct {
	IntakeRecordID uuid.UUID
	MedicationID   uuid.UUID
	MedicationName string
	AmountValue    float64
	AmountUnit     string
	Status         bool // is taken
	PlannedAt      time.Time
	TakenAt        time.Time
}

// ShowScheduleResponse is a response to get a plan.
type ShowScheduleResponse struct {
	Schedule []*ScheduleTime
}

// Execute executes the ShowSchedule command.
func (s *ShowScheduleService) Execute(
	ctx context.Context,
	req *ShowScheduleCommand,
) (*ShowScheduleResponse, error) {
	valErr := s.validator.ValidateStruct(req)
	if valErr != nil {
		return nil, ErrValidationFail
	}

	parsedUser, parsedStart, parsedEnd, err := parseInfo(req.UserID, req.StartDate, req.EndDate)
	if err != nil {
		return nil, ErrValidationFail
	}

	userPlans, err := s.planningRepo.UserPlans(ctx, parsedUser)
	if err != nil {
		return nil, ErrNoPlan
	}

	response := &ShowScheduleResponse{
		Schedule: s.scheduleList(ctx, userPlans, parsedStart, parsedEnd),
	}
	return response, nil
}

func parseInfo(
	userID string,
	startTime string,
	endTime string,
) (uuid.UUID, time.Time, time.Time, error) {
	parsedUser, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, time.Time{}, time.Time{}, ErrValidationFail
	}

	parsedStart, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		return uuid.Nil, time.Time{}, time.Time{}, ErrValidationFail
	}

	parsedEnd, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		return uuid.Nil, time.Time{}, time.Time{}, ErrValidationFail
	}

	return parsedUser, parsedStart, parsedEnd, nil
}

func (s *ShowScheduleService) scheduleList(
	ctx context.Context,
	userPlans []*plan.Plan,
	parsedStart time.Time,
	parsedEnd time.Time,
) []*ScheduleTime {
	pastScheduleList := make([]*ScheduleTime, 0, len(userPlans))
	futureScheduleList := make([]*ScheduleTime, 0, len(userPlans))
	for _, p := range userPlans {
		amountValue, amountUnit := p.Dosage()

		records, err := s.recordsRepo.GetByPlanID(ctx, p.ID())
		// if there is an error, it means that the plan has no records
		// because all of them are in the future and we will calculate them after a while
		medicationName, nameErr := s.medicationProvider.MedicationName(p.MedicationID(), p.UserID())
		if err == nil && nameErr == nil {
			for _, record := range records {
				if !record.PlannedTime().After(parsedEnd) &&
					!record.PlannedTime().Before(parsedStart) {
					pastScheduleList = append(pastScheduleList, &ScheduleTime{
						IntakeRecordID: record.ID(),
						MedicationID:   p.MedicationID(),
						MedicationName: medicationName,
						AmountValue:    amountValue,
						AmountUnit:     amountUnit,
						Status:         record.IsTaken(),
						PlannedAt:      record.PlannedTime(),
						TakenAt:        record.TakenAt(),
					})
				}
			}
		}
		// we are calculating all future records that are not created in db
		now := time.Now()
		futureTimes := p.Schedule(
			time.Date(
				now.Year(), now.Month(), now.Day(),
				0, 0, 0, 0, now.Location(),
			).Add(s.createdShift),
			parsedEnd,
		)

		for _, t := range futureTimes {
			futureScheduleList = append(futureScheduleList, &ScheduleTime{
				IntakeRecordID: uuid.Nil,
				MedicationID:   p.MedicationID(),
				MedicationName: medicationName,
				AmountValue:    amountValue,
				AmountUnit:     amountUnit,
				Status:         false,
				PlannedAt:      t,
				TakenAt:        time.Time{},
			})
		}
	}
	return append(futureScheduleList, pastScheduleList...)
}
