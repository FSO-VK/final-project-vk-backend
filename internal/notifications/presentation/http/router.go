package http

import (
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/httputil"
	"github.com/gorilla/mux"
)

// Router returns a new HTTP router.
func Router(
	notificationHandlers *NotificationsHandlers,
	authMw *httputil.AuthMiddleware,
) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/vapidPublicKey", notificationHandlers.GetVapidPublicKey).Methods("GET")
	r.HandleFunc("/subscribe", notificationHandlers.CreateSubscription).Methods("POST")
	r.HandleFunc("/subscribe", notificationHandlers.DeleteSubscription).Methods("DELETE")
	r.HandleFunc("/send", notificationHandlers.SendNotification).Methods("POST")
	r.HandleFunc("/action/{id}", notificationHandlers.InteractWithNotification).Methods("POST")

	panicMiddleware := httputil.NewPanicRecoveryMiddleware()
	r.Use(panicMiddleware.Middleware)
	r.Use(authMw.AuthMiddlewareWrapper)

	return r
}