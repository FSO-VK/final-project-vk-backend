package http

import (
	httph "github.com/FSO-VK/final-project-vk-backend/internal/transport/http"
	"github.com/gorilla/mux"
)

// Router returns a new HTTP router.
func Router(
	medicationHandlers *MedicationHandlers,
) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/all", medicationHandlers.GetMedicationList).Methods("GET")
	r.HandleFunc("", medicationHandlers.AddMedication).Methods("POST")
	r.HandleFunc("/{id}", medicationHandlers.UpdateMedication).Methods("PUT")
	r.HandleFunc("/{id}", medicationHandlers.DeleteMedication).Methods("DELETE")

	panicMiddleware := httph.NewPanicRecoveryMiddleware()
	r.Use(panicMiddleware.Middleware)

	return r
}
