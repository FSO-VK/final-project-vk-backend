package http

import (
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/httputil"
	"github.com/gorilla/mux"
)

// InternalRouter returns a new internal router for cross microservice communication.
func InternalRouter(
	medicationHandlers *MedicationHandlers,
	loggingMw *httputil.LoggingMiddleware,
) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc(
		"/internal/medication/{id}/{user_id}",
		medicationHandlers.InternalGetMedicationByID,
	).Methods("GET")
	panicMiddleware := httputil.NewPanicRecoveryMiddleware()
	r.Use(panicMiddleware.Middleware)
	r.Use(loggingMw.MiddlewareNetHTTP)
	return r
}
