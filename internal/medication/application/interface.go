package application

import "context"

// MedicineService is a service for medicine.
type MedicineService interface {
	AddMedicine(ctx context.Context, req *AddMedicineRequest) (*AddMedicineResponse, error)
	UpdateMedicine(ctx context.Context, req *UpdateMedicineRequest) (*UpdateMedicineResponse, error)
	DeleteMedicine(ctx context.Context, req *DeleteMedicineRequest) (*DeleteMedicineResponse, error)
	GetMedicineList(
		ctx context.Context,
		req *GetMedicineListRequest,
	) (*GetMedicineListResponse, error)
}

// AddMedicineRequest is a request to add a medicine.
type AddMedicineRequest struct {
	Name         string `validate:"required"`
	CategoriesID []uint
	Items        uint   `validate:"required"`
	ItemsUnit    string `validate:"required"`
	Expires      string `validate:"required"`
}

// AddMedicineResponse is a response to add a medicine.
type AddMedicineResponse struct {
	ID uint
}

// UpdateMedicineRequest is a request to update a medicine.
type UpdateMedicineRequest struct {
	ID           uint
	Name         string `validate:"required"`
	CategoriesID []uint
	Items        uint   `validate:"required"`
	ItemsUnit    string `validate:"required"`
	Expires      string `validate:"required"`
}

// UpdateMedicineResponse is a response to update a medicine.
type UpdateMedicineResponse struct {
	ID           uint
	Name         string
	CategoriesID []uint
	Items        uint
	ItemsUnit    string
	Expires      string
}

// DeleteMedicineRequest is a request to delete a medicine.
type DeleteMedicineRequest struct {
	ID uint
}

// DeleteMedicineResponse is a response to delete a medicine.
type DeleteMedicineResponse struct{}

// GetMedicineListRequest is a request to get a list of medicines.
type GetMedicineListRequest struct{}

// MedicineListItem contains information about one medicine in the list.
type MedicineListItem struct {
	ID        uint
	Name      string
	Items     uint
	ItemsUnit string
	Expires   string
}

// GetMedicineListResponse contains a list of medicines.
type GetMedicineListResponse struct {
	List []*MedicineListItem
}
