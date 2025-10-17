// Package application is a package for application logic of the medication service.
package application

// MedicineApplication is a struct for application logic of the medication service.
type MedicineApplication struct {
	GetMedicineList GetMedicineList
	AddMedicine     AddMedicine
	UpdateMedicine  UpdateMedicine
	DeleteMedicine  DeleteMedicine
}
