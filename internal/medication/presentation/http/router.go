package http

import (
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/httputil"
	"github.com/gorilla/mux"
)

// Router returns a new HTTP router.
func Router(
	medicationHandlers *MedicationHandlers,
	authMw *httputil.AuthMiddleware,
	loggingMw *httputil.LoggingMiddleware,
) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/medication/all", medicationHandlers.GetMedicationBox).Methods("GET")
	r.HandleFunc("/medication/{id}", medicationHandlers.GetMedicationByID).Methods("GET")
	r.HandleFunc("/medication", medicationHandlers.AddMedication).Methods("POST")
	r.HandleFunc("/medication/{id}", medicationHandlers.UpdateMedication).Methods("PUT")
	r.HandleFunc("/medication/{id}", medicationHandlers.DeleteMedication).Methods("DELETE")
	r.HandleFunc("/scan", medicationHandlers.DataMatrixInformation).Methods("GET")
	r.HandleFunc("/medication/{id}/assistant", medicationHandlers.InstructionAssistant).
		Methods("GET")
	r.HandleFunc("/medication/{id}/instruction", medicationHandlers.GetInstruction).Methods("GET")
	r.HandleFunc("/medication/{id}/take", medicationHandlers.TakeMedication).Methods("POST")

	panicMiddleware := httputil.NewPanicRecoveryMiddleware()
	r.Use(panicMiddleware.Middleware)
	r.Use(loggingMw.MiddlewareNetHTTP)
	r.Use(authMw.AuthMiddlewareWrapper)

	return r
}
