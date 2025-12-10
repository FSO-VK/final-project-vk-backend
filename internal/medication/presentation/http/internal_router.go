package http

import (
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/httputil"
	"github.com/gorilla/mux"
)

// InternalRouter returns a new internal router for cross microservice communication.
func InternalRouter(
	medicationHandlers *MedicationHandlers,
) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/medication/{id}/{user_id}", medicationHandlers.GetMedicationByID).Methods("GET")
	panicMiddleware := httputil.NewPanicRecoveryMiddleware()
	r.Use(panicMiddleware.Middleware)
	return r
}
