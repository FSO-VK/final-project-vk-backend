package http

import (
	httph "github.com/FSO-VK/final-project-vk-backend/internal/transport/http"
	"github.com/gorilla/mux"
)

func Router(
	medicineHandlers *MedicineHandlers,
) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/medicines", medicineHandlers.GetMedicineList).Methods("GET")
	r.HandleFunc("/medicine", medicineHandlers.AddMedicine).Methods("POST")
	r.HandleFunc("/medicine/{id}", medicineHandlers.UpdateMedicine).Methods("PUT")
	r.HandleFunc("/medicine/{id}", medicineHandlers.DeleteMedicine).Methods("DELETE")

	// r.Use(mux.CORSMethodMiddleware(r))
	panicMiddleware := httph.NewPanicRecoveryMiddleware()
	r.Use(panicMiddleware.Middleware)

	return r
}
