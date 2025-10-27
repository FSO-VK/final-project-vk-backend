// Package application implements dependency injection container for medication service use cases.
package application

// MedicationApplication is a dependency injection container that aggregates all use cases
// for the medication domain to be injected from main.go.
type MedicationApplication struct {
	GetMedicationList     GetMedicationList
	AddMedication         AddMedication
	UpdateMedication      UpdateMedication
	DeleteMedication      DeleteMedication
	DataMatrixInformation DataMatrixInformation
}
