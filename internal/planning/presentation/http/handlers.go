// Package http is a package for http handlers
package http

import (
	"net/http"

	"github.com/FSO-VK/final-project-vk-backend/internal/planning/application"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/httputil"
	"github.com/FSO-VK/final-project-vk-backend/pkg/api"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
				MedicationID:   p.MedicationID,
				UserID:         p.UserID,
				AmountValue:    p.AmountValue,
				AmountUnit:     p.AmountUnit,
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
