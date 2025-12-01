// Package application is a package for application logic of the planning service.
package application

import (
	"context"
	"time"

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
	planningRepo plan.Repository
	recordsRepo  record.Repository
	validator    validator.Validator
	// createdShift is the offset from 00:00 when records are generated.
	// At 00:00 + createdShift, all records for that day are created. (basically 24h - today creating for the next day)
	createdShift time.Duration
}

// NewShowScheduleService returns a new ShowScheduleService.
func NewShowScheduleService(
	planningRepo plan.Repository,
	recordsRepo record.Repository,
	valid validator.Validator,
	createdShift time.Duration,
) *ShowScheduleService {
	return &ShowScheduleService{
		planningRepo: planningRepo,
		recordsRepo:  recordsRepo,
		validator:    valid,
		createdShift: createdShift,
	}
}

// ShowScheduleCommand is a request to get a plan.
type ShowScheduleCommand struct {
	UserID    string `validate:"required,uuid"`
	StartDate string `validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	EndDate   string `validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
}

// IntakeRecord is an aggregate that represents a record for medication intake.
type ScheduleTime struct {
	intakeRecordID uuid.UUID
	medicationID   uuid.UUID
	medicationName string
	AmountValue    float64
	AmountUnit     string
	status         bool // is taken
	plannedAt      time.Time
	takenAt        time.Time
}

// ShowScheduleResponse is a response to get a plan.
type ShowScheduleResponse struct {
	schedule []*ScheduleTime
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

	parsedUser, parsedStart, parsedEnd, err := ParseInfo(req.UserID, req.StartDate, req.EndDate)
	if err != nil {
		return nil, ErrValidationFail
	}

	userPlans, err := s.planningRepo.UserPlans(ctx, parsedUser)
	if err != nil {
		return nil, ErrNoPlan
	}

	pastScheduleList := make([]*ScheduleTime, 0, len(*userPlans))
	futureScheduleList := make([]*ScheduleTime, 0, len(*userPlans))

	for _, onePlan := range *userPlans {
		amountValue, amountUnit := onePlan.Dosage()

		records, err := s.recordsRepo.GetByPlanID(ctx, onePlan.ID())
		// if there is an error, it means that the plan has no records
		// because all of them are in the future and we will calculate them after a while
		if err == nil {
			for _, oneRecord := range records {
				if !oneRecord.PlannedTime().After(parsedEnd) &&
					!oneRecord.PlannedTime().Before(parsedStart) {
					pastScheduleList = append(pastScheduleList, &ScheduleTime{
						intakeRecordID: oneRecord.ID(),
						medicationID:   onePlan.MedicationID(),
						medicationName: "Medication Name", // need client
						AmountValue:    amountValue,
						AmountUnit:     amountUnit,
						status:         oneRecord.IsTaken(),
						plannedAt:      oneRecord.PlannedTime(),
						takenAt:        oneRecord.TakenAt(),
					})
				}
			}
		}
		// we are calculating all future records that are not created in db
		futureTimes := onePlan.Schedule(
			time.Now().Truncate(24*time.Hour).Add(s.createdShift),
			parsedEnd,
		)

		for _, oneTime := range futureTimes {
			futureScheduleList = append(futureScheduleList, &ScheduleTime{
				intakeRecordID: uuid.Nil,
				medicationID:   onePlan.MedicationID(),
				medicationName: "Medication Name", // need client
				AmountValue:    amountValue,
				AmountUnit:     amountUnit,
				status:         false,
				plannedAt:      oneTime,
				takenAt:        time.Time{},
			})
		}
	}

	response := &ShowScheduleResponse{
		schedule: append(futureScheduleList, pastScheduleList...),
	}
	return response, nil
}

func ParseInfo(
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
