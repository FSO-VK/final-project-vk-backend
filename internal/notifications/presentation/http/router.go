package http

import (
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/httputil"
	"github.com/gin-gonic/gin"
)

// Router returns a new Gin engine with routes and Gin-native middleware.
func Router(
	notificationHandlers *NotificationsHandlers,
	authMw *httputil.GinAuthMiddleware,
) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(httputil.NewPanicRecoveryMiddleware().Handler())
	r.Use(authMw.Middleware())

	r.GET("/vapidPublicKey", notificationHandlers.GetVapidPublicKeyGin)
	r.POST("/pushSubscription", notificationHandlers.CreateSubscriptionGin)
	r.DELETE("/pushSubscription/:id", notificationHandlers.DeleteSubscriptionGin)
	r.POST("/send", notificationHandlers.SendNotificationGin)

	return r
}
