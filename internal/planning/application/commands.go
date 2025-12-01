package application

import (
	"fmt"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/planning/domain/plan"
	"github.com/google/uuid"
	"github.com/teambition/rrule-go"
)

// PlanDraft is all the information needed to create a plan.
type PlanDraft struct {
	ID             string
	MedicationID   string
	UserID         string
	AmountValue    float64
	AmountUnit     string
	Condition      string
	StartDate      string
	EndDate        string
	Duration       string
	RecurrenceRule []string
}

// parsePlanID parses or generates a plan UUID.
func parsePlanID(id string) (uuid.UUID, error) {
	if id == "" {
		return uuid.NewV7()
	}
	return uuid.Parse(id)
}

// parseUUIDs parses UserID and MedicationID.
func parseUUIDs(userID, medicationID string) (uuid.UUID, uuid.UUID, error) {
	parsedUser, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	parsedMedicationID, err := uuid.Parse(medicationID)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	return parsedUser, parsedMedicationID, nil
}

// parseDates parses start and end dates.
func parseDates(startDate, endDate string) (time.Time, time.Time, error) {
	parsedStart, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	parsedEnd, err := time.Parse(time.RFC3339, endDate)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return parsedStart, parsedEnd, nil
}

// parseRRules parses recurrence rules.
func parseRRules(rules []string) ([]*rrule.RRule, error) {
	if len(rules) == 0 {
		return nil, ErrUnsupportedRrule
	}

	result := make([]*rrule.RRule, 0, len(rules))
	for _, ruleStr := range rules {
		rule, err := rrule.StrToRRule(ruleStr)
		if err != nil {
			return nil, ErrUnsupportedRrule
		}
		result = append(result, rule)
	}
	return result, nil
}

// createPlan creates a new plan from draft.
func createPlan(req *PlanDraft) (*plan.Plan, error) {
	planID, err := parsePlanID(req.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid plan id: %w", err)
	}

	userID, medicationID, err := parseUUIDs(req.UserID, req.MedicationID)
	if err != nil {
		return nil, ErrValidationFail
	}

	dosage, err := plan.NewDosage(req.AmountValue, req.AmountUnit)
	if err != nil {
		return nil, fmt.Errorf("invalid dosage: %w", err)
	}

	startDate, endDate, err := parseDates(req.StartDate, req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid dates: %w", err)
	}

	rules, err := parseRRules(req.RecurrenceRule)
	if err != nil {
		return nil, ErrUnsupportedRrule
	}

	schedule, err := plan.NewSchedule(startDate, endDate, rules)
	if err != nil {
		return nil, fmt.Errorf("invalid schedule: %w", err)
	}

	return plan.NewPlan(
		planID,
		medicationID,
		userID,
		dosage,
		schedule,
		req.Condition,
		time.Now(),
		time.Now(),
	)
}
