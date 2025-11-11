// Package application is a package for application logic of the planning service.
package application

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/planning/domain/plan"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/google/uuid"
)

const (
	dateLayout = "2006-01-02T15:04:05.000Z"
	timeLayout = "15:04"
)

// AddPlan is an interface for adding a notification.
type AddPlan interface {
	Execute(
		ctx context.Context,
		cmd *AddPlanCommand,
	) (*AddPlanResponse, error)
}

// AddPlanService is a service for creating a subscription.
type AddPlanService struct {
	planningRepo plan.Repository
	validator    validator.Validator
}

// NewAddPlanService returns a new AddPlanService.
func NewAddPlanService(
	planningRepo plan.Repository,
	valid validator.Validator,
) *AddPlanService {
	return &AddPlanService{
		planningRepo: planningRepo,
		validator:    valid,
	}
}

// AddPlanCommand is a request to add a plan.
type AddPlanCommand struct {
	MedicationID string
	UserID       string
	AmountValue  float64
	AmountUnit   string

	TakingTime string // "15:04"
	TimeZone   string // +03:00 from UTC
	Frequency  string // "daily", "weekly", "monthly", "custom"
	WeekDays   []int  // [1,3,5] для пн,ср,пт (0-6, 0=воскресенье)
	MonthDay   int    // 15 для 15-го числа каждого месяца

	StartDate       string
	EndDate         string
	IntakeCondition string
}

// AddPlanResponse is a response to add a plan.
type AddPlanResponse struct {
}

// Execute executes the AddPlan command.
func (s *AddPlanService) Execute(
	ctx context.Context,
	req *AddPlanCommand,
) (*AddPlanResponse, error) {
	valErr := s.validator.ValidateStruct(req)
	if valErr != nil {
		return nil, fmt.Errorf("request is not valid: %w", valErr)
	}
	parsedUser, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid format: %w", err)
	}
	parsedMedicationID, err := uuid.Parse(req.MedicationID)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid format: %w", err)
	}

	dosage, err := plan.NewDosage(
		req.AmountValue,
		req.AmountUnit,
	)
	if err != nil {
		return nil, fmt.Errorf("invalid dosage: %w", err)
	}
	cronString, err := generateCronExpression(req)

	if err != nil {
		return nil, fmt.Errorf("invalid cron expression: %w", err)
	}

	schedule, err := plan.NewSchedule(cronString)

	parsedStart, err := time.Parse(dateLayout, req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid course end: %w", err)
	}
	start, err := plan.NewCourseStart(parsedStart)
	if err != nil {
		return nil, fmt.Errorf("invalid course start: %w", err)
	}
	parsedEnd, err := time.Parse(dateLayout, req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid course end: %w", err)
	}
	end, err := plan.NewCourseEnd(parsedEnd)
	if err != nil {
		return nil, fmt.Errorf("invalid course end: %w", err)
	}
	newPlan, err := plan.NewPlan(
		uuid.New(),
		parsedMedicationID,
		parsedUser,
		dosage,
		schedule,
		start,
		end,
		req.IntakeCondition,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create plan: %w", err)
	}

	err = s.planningRepo.Save(ctx, *newPlan)

	if err != nil {
		return nil, fmt.Errorf("failed to save plan: %w", err)
	}

	response := &AddPlanResponse{}
	return response, nil
}

func generateCronExpression(cmd *AddPlanCommand) (string, error) {

	takingTime, err := time.Parse(timeLayout, cmd.TakingTime)
	if err != nil {
		return "", err
	}

	minute := takingTime.Minute()
	hour := takingTime.Hour()

	switch cmd.Frequency {
	case "daily":
		return fmt.Sprintf("%d %d * * *", minute, hour), nil

	case "weekly":
		if len(cmd.WeekDays) == 0 {
			return "", errors.New("week days required for weekly frequency")
		}
		days := make([]string, len(cmd.WeekDays))
		for i, day := range cmd.WeekDays {
			days[i] = strconv.Itoa(day)
		}
		daysStr := strings.Join(days, ",")
		return fmt.Sprintf("%d %d * * %s", minute, hour, daysStr), nil

	case "monthly":
		if cmd.MonthDay == 0 {
			cmd.MonthDay = 1
		}
		return fmt.Sprintf("%d %d %d * *", minute, hour, cmd.MonthDay), nil

	default:
		return "", fmt.Errorf("unsupported frequency: %s", cmd.Frequency)
	}
}
