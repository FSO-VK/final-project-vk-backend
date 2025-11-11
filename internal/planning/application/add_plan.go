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

var (
	// ErrWeekDaysRequired is an error when week days are required for weekly frequency.
	ErrWeekDaysRequired = errors.New("week days required for weekly frequency")
	// ErrUnsupportedFrequency is an error when frequency is unsupported.
	ErrUnsupportedFrequency = errors.New("unsupported frequency")
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

	CustomFrequency CustomFrequency

	StartDate       string
	EndDate         string
	IntakeCondition string
}

type CustomFrequency struct {
	TakingTime string // "15:04"
	Frequency  string // "every x days", "daily", "weekly", "monthly", "custom"
	FrequencyX int
	WeekDays   []int // [1,3,5] для пн,ср,пт (0-6, 0=воскресенье)
	MonthDay   int   // 15 для 15-го числа каждого месяца
}

// AddPlanResponse is a response to add a plan.
type AddPlanResponse struct{}

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

	newPlan, err := createPlan(req, parsedUser, parsedMedicationID)
	if err != nil {
		return nil, fmt.Errorf("failed to create plan: %w", err)
	}

	err = s.planningRepo.Save(ctx, newPlan)
	if err != nil {
		return nil, fmt.Errorf("failed to save plan: %w", err)
	}

	response := &AddPlanResponse{}
	return response, nil
}

func convertToUTC(
	takingTimeStr string,
) (int, int, int, error) {
	var t time.Time
	var err error
	t, err = time.Parse(time.RFC3339Nano, takingTimeStr)
	if err != nil {
		t, err = time.Parse(time.RFC3339, takingTimeStr)
		if err != nil {
			t, err = time.Parse(dateLayout, takingTimeStr)
			if err != nil {
				return 0, 0, 0, fmt.Errorf("invalid taking time format: %w", err)
			}
		}
	}

	localHour := t.Hour()
	localMinute := t.Minute()
	_, offsetSeconds := t.Zone()
	offsetMinutes := offsetSeconds / 60

	totalMinutes := localHour*60 + localMinute - offsetMinutes

	dayOffset := 0
	if totalMinutes < 0 {
		totalMinutes += 24 * 60
		dayOffset = -1
	} else if totalMinutes >= 24*60 {
		totalMinutes -= 24 * 60
		dayOffset = 1
	}

	hour := totalMinutes / 60
	minute := totalMinutes % 60

	return hour, minute, dayOffset, nil
}

func generateCronExpression(cmd *CustomFrequency, hour int, minute int) (string, error) {
	switch cmd.Frequency {
	case "every x days":
		return fmt.Sprintf("%d %d */%d * *", minute, hour, cmd.FrequencyX), nil

	case "daily":
		return fmt.Sprintf("%d %d * * *", minute, hour), nil

	case "weekly":
		if len(cmd.WeekDays) == 0 {
			return "", ErrWeekDaysRequired
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
		return "", ErrUnsupportedFrequency
	}
}

func createPlan(req *AddPlanCommand,
	parsedUser uuid.UUID,
	parsedMedicationID uuid.UUID,
) (plan.Plan, error) {
	dosage, err := plan.NewDosage(
		req.AmountValue,
		req.AmountUnit,
	)
	if err != nil {
		return plan.Plan{}, fmt.Errorf("invalid dosage: %w", err)
	}

	hour, minute, dayOffset, err := convertToUTC(
		req.CustomFrequency.TakingTime,
	)
	if err != nil {
		return plan.Plan{}, err
	}

	cronString, err := generateCronExpression(
		&req.CustomFrequency,
		hour,
		minute,
	)
	if err != nil {
		return plan.Plan{}, fmt.Errorf("invalid cron expression: %w", err)
	}

	schedule, err := plan.NewSchedule(cronString)
	if err != nil {
		return plan.Plan{}, fmt.Errorf("invalid schedule: %w", err)
	}

	parsedStart, err := parseDate(req.StartDate, dayOffset)
	if err != nil {
		return plan.Plan{}, fmt.Errorf("invalid course start: %w", err)
	}
	start, err := plan.NewCourseStart(parsedStart)
	if err != nil {
		return plan.Plan{}, fmt.Errorf("invalid course start: %w", err)
	}

	parsedEnd, err := parseDate(req.EndDate, dayOffset)
	if err != nil {
		return plan.Plan{}, fmt.Errorf("invalid course end: %w", err)
	}
	end, err := plan.NewCourseEnd(parsedEnd)
	if err != nil {
		return plan.Plan{}, fmt.Errorf("invalid course end: %w", err)
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
	return *newPlan, err
}

func parseDate(date string, dayOffset int) (time.Time, error) {
	parsedStart, err := time.Parse(dateLayout, date)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid course start: %w", err)
	}
	if dayOffset > 0 {
		parsedStart = parsedStart.AddDate(0, 0, -1)
	}
	return parsedStart, nil
}
