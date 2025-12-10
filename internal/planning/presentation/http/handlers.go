// Package http is a package for http handlers
package http

import (
	"net/http"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/planning/application"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/httputil"
	"github.com/FSO-VK/final-project-vk-backend/pkg/api"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	// SlugID is a slug for id.
	SlugID = "id"
)

// PlanningHandlers is a handler for Planning.
type PlanningHandlers struct {
	app    *application.PlanningApplication
	logger *logrus.Entry
}

// NewHandlers creates a new PlanningHandlers.
func NewHandlers(
	app *application.PlanningApplication,
	logger *logrus.Entry,
) *PlanningHandlers {
	return &PlanningHandlers{
		app:    app,
		logger: logger,
	}
}

// GetAllUsersPlansItem returns user's plans.
type GetAllUsersPlansItem struct {
	// embedded struct
	PlanObject `json:",inline"`

	ID string `json:"id"`
}

// GetAllUsersPlansJSONResponse returns user's plans.
type GetAllUsersPlansJSONResponse struct {
	AllUserPlans []GetAllUsersPlansItem `json:"allUserPlans"`
}

// GetAllUsersPlans gets all users plans.
func (h *PlanningHandlers) GetAllUsersPlans(c *gin.Context) {
	auth, err := httputil.GetAuthFromCtx(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, api.Response[any]{
			StatusCode: http.StatusUnauthorized,
			Error:      api.MsgUnauthorized,
			Body:       struct{}{},
		})
		return
	}

	serviceRequest := &application.GetAllPlansCommand{
		UserID: auth.UserID,
	}
	plans, err := h.app.GetAllPlans.Execute(c.Request.Context(), serviceRequest)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get plan")
		c.JSON(http.StatusInternalServerError, api.Response[any]{
			StatusCode: http.StatusInternalServerError,
			Body:       struct{}{},
			Error:      MsgFailedToGetPlan,
		})
		return
	}

	response := &GetAllUsersPlansJSONResponse{
		AllUserPlans: make([]GetAllUsersPlansItem, 0),
	}

	for _, p := range plans.Plans {
		response.AllUserPlans = append(response.AllUserPlans, GetAllUsersPlansItem{
			PlanObject: PlanObject{
				MedicationID: p.MedicationID,
				UserID:       p.UserID,
				Amount: AmountObject{
					Value: p.AmountValue,
					Unit:  p.AmountUnit,
				},
				Condition:      p.Condition,
				StartDate:      p.StartDate,
				EndDate:        p.EndDate,
				RecurrenceRule: p.RecurrenceRule,
			},
			ID: p.ID,
		})
	}

	c.JSON(http.StatusOK, api.Response[any]{
		StatusCode: http.StatusOK,
		Body:       response,
		Error:      "",
	})
}

// AddPlanJSONRequest is a request for AddPlan.
type AddPlanJSONRequest struct {
	// embedded struct
	PlanObject `json:",inline"`
}

// AddPlanJSONResponse is a response for AddPlan.
type AddPlanJSONResponse struct {
	// embedded struct
	PlanObject `json:",inline"`

	ID string `json:"id"`
}

// AddPlan adds a plan.
func (h *PlanningHandlers) AddPlan(c *gin.Context) {
	auth, err := httputil.GetAuthFromCtx(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, api.Response[any]{
			StatusCode: http.StatusUnauthorized,
			Error:      api.MsgUnauthorized,
			Body:       struct{}{},
		})
		return
	}

	var reqJSON AddPlanJSONRequest
	if err := c.ShouldBindJSON(&reqJSON); err != nil {
		h.logger.WithError(err).Error("Failed to bind request body")
		c.JSON(http.StatusBadRequest, api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Body:       struct{}{},
			Error:      api.MsgBadBody,
		})
		return
	}

	command := &application.AddPlanCommand{
		MedicationID:   reqJSON.MedicationID,
		UserID:         auth.UserID,
		AmountValue:    reqJSON.Amount.Value,
		AmountUnit:     reqJSON.Amount.Unit,
		Condition:      reqJSON.Condition,
		StartDate:      reqJSON.StartDate,
		EndDate:        reqJSON.EndDate,
		RecurrenceRule: reqJSON.RecurrenceRule,
	}
	serviceResponse, err := h.app.AddPlan.Execute(c.Request.Context(), command)
	if err != nil {
		h.logger.WithError(err).Error("Failed to add plan")
		c.JSON(http.StatusInternalServerError, api.Response[any]{
			StatusCode: http.StatusInternalServerError,
			Body:       struct{}{},
			Error:      MsgFailedToAddPlan,
		})
		return
	}

	response := &AddPlanJSONResponse{
		PlanObject: PlanObject{
			MedicationID: serviceResponse.MedicationID,
			UserID:       serviceResponse.UserID,
			Amount: AmountObject{
				Value: serviceResponse.AmountValue,
				Unit:  serviceResponse.AmountUnit,
			},
			Condition:      serviceResponse.Condition,
			StartDate:      serviceResponse.StartDate,
			EndDate:        serviceResponse.EndDate,
			RecurrenceRule: serviceResponse.RecurrenceRule,
		},
		ID: serviceResponse.ID,
	}

	c.JSON(http.StatusOK, api.Response[any]{
		StatusCode: http.StatusOK,
		Body:       response,
		Error:      "",
	})
}

// GetPlanByIDJSONResponse is a response for GetPlanByID.
type GetPlanByIDJSONResponse struct {
	// embedded struct
	PlanObject

	ID string `json:"id"`
}

// GetPlanByID gets plan by id.
func (h *PlanningHandlers) GetPlanByID(c *gin.Context) {
	auth, err := httputil.GetAuthFromCtx(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, api.Response[any]{
			StatusCode: http.StatusUnauthorized,
			Error:      api.MsgUnauthorized,
			Body:       struct{}{},
		})
		return
	}

	slugUserID := c.Param(SlugID)
	if slugUserID == "" {
		h.logger.Error("Subscription ID not found in path params")
		c.JSON(http.StatusBadRequest, api.Response[any]{
			StatusCode: http.StatusBadRequest,
			Error:      MsgMissingSlug,
			Body:       struct{}{},
		})
		return
	}

	command := &application.GetPlanCommand{
		UserID: auth.UserID,
		ID:     slugUserID,
	}

	p, err := h.app.GetPlan.Execute(c.Request.Context(), command)
	if err != nil {
		h.logger.WithError(err).Error("Failed to send notification")
		c.JSON(http.StatusInternalServerError, api.Response[any]{
			StatusCode: http.StatusInternalServerError,
			Body:       struct{}{},
			Error:      MsgFailedToGetPlan,
		})
		return
	}

	response := &GetPlanByIDJSONResponse{
		PlanObject: PlanObject{
			MedicationID: p.MedicationID,
			UserID:       p.UserID,
			Amount: AmountObject{
				Value: p.AmountValue,
				Unit:  p.AmountUnit,
			},
			Condition:      p.Condition,
			StartDate:      p.StartDate,
			EndDate:        p.EndDate,
			RecurrenceRule: p.RecurrenceRule,
		},
		ID: p.ID,
	}

	c.JSON(http.StatusOK, api.Response[any]{
		StatusCode: http.StatusOK,
		Body:       response,
		Error:      "",
	})
}

// GetAllUsersPlansItem returns single schedule item.
type ShowScheduleItem struct {
	IntakeRecordID string       `json:"intakeRecordId"`
	MedicationID   string       `json:"medicationId"`
	MedicationName string       `json:"medicationName"`
	Amount         AmountObject `json:"amount"`
	Status         string       `json:"status"`
	PlannedAt      string       `json:"plannedAt"`
	TakenAt        string       `json:"takenAt"`
}

// ShowScheduleJSONResponse returns schedule.
type ShowScheduleJSONResponse struct {
	Schedule []ShowScheduleItem `json:"schedule"`
}

// ShowSchedule gets nearest schedule.
func (h *PlanningHandlers) ShowSchedule(c *gin.Context) {
	auth, err := httputil.GetAuthFromCtx(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, api.Response[any]{
			StatusCode: http.StatusUnauthorized,
			Error:      api.MsgUnauthorized,
			Body:       struct{}{},
		})
		return
	}

	serviceRequest := &application.ShowScheduleCommand{
		UserID:    auth.UserID,
		StartDate: c.Query("start"),
		EndDate:   c.Query("end"),
	}

	sh, err := h.app.ShowSchedule.Execute(c.Request.Context(), serviceRequest)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get schedule")
		c.JSON(http.StatusInternalServerError, api.Response[any]{
			StatusCode: http.StatusInternalServerError,
			Body:       struct{}{},
			Error:      MsgFailedToGetPlan,
		})
		return
	}

	response := &ShowScheduleJSONResponse{
		Schedule: make([]ShowScheduleItem, 0),
	}

	for _, s := range sh.Schedule {
		response.Schedule = append(response.Schedule, ShowScheduleItem{
			IntakeRecordID: s.IntakeRecordID.String(),
			MedicationID:   s.MedicationID.String(),
			MedicationName: s.MedicationName,
			Amount: AmountObject{
				Value: s.AmountValue,
				Unit:  s.AmountUnit,
			},
			Status:    s.Status,
			PlannedAt: s.PlannedAt.Format(time.RFC3339),
			TakenAt:   s.TakenAt.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, api.Response[any]{
		StatusCode: http.StatusOK,
		Body:       response,
		Error:      "",
	})
}
