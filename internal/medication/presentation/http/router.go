package http

import (
	httph "github.com/FSO-VK/final-project-vk-backend/internal/transport/http"
	"github.com/gorilla/mux"
)

// Router returns a new HTTP router.
func Router(
	medicineHandlers *MedicineHandlers,
) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/medication/all", medicineHandlers.GetMedicineList).Methods("GET")
	r.HandleFunc("/medication", medicineHandlers.AddMedicine).Methods("POST")
	r.HandleFunc("/medication/{id}", medicineHandlers.UpdateMedicine).Methods("PUT")
	r.HandleFunc("/medication/{id}", medicineHandlers.DeleteMedicine).Methods("DELETE")

	// r.Use(mux.CORSMethodMiddleware(r))
	panicMiddleware := httph.NewPanicRecoveryMiddleware()
	r.Use(panicMiddleware.Middleware)

	return r
}
