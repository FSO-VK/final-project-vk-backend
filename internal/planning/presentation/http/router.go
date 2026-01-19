package http

import (
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/httputil"
	"github.com/gin-gonic/gin"
)

// Router returns a new Gin engine with routes and Gin-native middleware.
func Router(
	planningHandlers *PlanningHandlers,
	authMw *httputil.AuthMiddleware,
) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(httputil.NewPanicRecoveryMiddleware().Handler())
	authGroup := r.Group("/")
	authGroup.Use(authMw.Middleware())
	{
		authGroup.POST(
			"/intake/:id/take",
			planningHandlers.TakeMedication,
		)
		authGroup.POST(
			"/intake/:id/change",
			planningHandlers.ChangeTakeMedication,
		)
		authGroup.DELETE(
			"/intake/:id/cancel",
			planningHandlers.CancelMedicationTake,
		)
		authGroup.GET("/plan/all", planningHandlers.GetAllUsersPlans)
		authGroup.GET("/plan/:id", planningHandlers.GetPlanByID)
		authGroup.POST("/plan", planningHandlers.AddPlan)
		authGroup.GET("/plan/schedule", planningHandlers.ShowSchedule)
		authGroup.DELETE("/plan/:id", planningHandlers.FinishPlan)
	}

	return r
}
