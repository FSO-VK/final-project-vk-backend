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
	panicMw := httputil.NewPanicRecoveryMiddleware()

	// PUBLIC ROUTER
	public := mux.NewRouter()
	public.HandleFunc("/medication/all", medicationHandlers.GetMedicationBox).Methods("GET")
	public.HandleFunc("/medication/{id}", medicationHandlers.GetMedicationByID).Methods("GET")
	public.HandleFunc("/medication", medicationHandlers.AddMedication).Methods("POST")
	public.HandleFunc("/medication/{id}", medicationHandlers.UpdateMedication).Methods("PUT")
	public.HandleFunc("/medication/{id}", medicationHandlers.DeleteMedication).Methods("DELETE")
	public.HandleFunc("/scan", medicationHandlers.DataMatrixInformation).Methods("GET")
	public.HandleFunc(
		"/medication/{id}/assistant",
		medicationHandlers.InstructionAssistant,
	).Methods("GET")
	public.HandleFunc(
		"/medication/{id}/instruction",
		medicationHandlers.GetInstruction,
	).Methods("GET")

	// middlewares for public API
	public.Use(panicMw.Middleware)
	public.Use(loggingMw.MiddlewareNetHTTP)
	public.Use(authMw.AuthMiddlewareWrapper)

	// INTERNAL ROUTER
	internal := mux.NewRouter()
	internal.HandleFunc(
		"/internal/medication/{id}",
		medicationHandlers.InternalGetMedicationByID,
	).Methods("GET")

	// ROOT ROUTER
	root := mux.NewRouter()
	root.PathPrefix("/internal/").Handler(internal)
	root.PathPrefix("/").Handler(public)

	return root
}
