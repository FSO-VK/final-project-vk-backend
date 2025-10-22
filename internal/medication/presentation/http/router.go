package http

import (
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/httputil"
	"github.com/gorilla/mux"
)

// Router returns a new HTTP router.
func Router(
	medicationHandlers *MedicationHandlers,
) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/medication/all", medicationHandlers.GetMedicationList).Methods("GET")
	r.HandleFunc("/medication", medicationHandlers.AddMedication).Methods("POST")
	r.HandleFunc("/medication/{id}", medicationHandlers.UpdateMedication).Methods("PUT")
	r.HandleFunc("/medication/{id}", medicationHandlers.DeleteMedication).Methods("DELETE")

	panicMiddleware := httputil.NewPanicRecoveryMiddleware()
	r.Use(panicMiddleware.Middleware)

	return r
}
