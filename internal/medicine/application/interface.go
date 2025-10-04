package application

import "context"

type MedicineService interface {
	AddMedicine(ctx context.Context, req *AddMedicineRequest) (*AddMedicineResponse, error)
	UpdateMedicine(ctx context.Context, req *UpdateMedicineRequest) (*UpdateMedicineResponse, error)
	DeleteMedicine(ctx context.Context, req *DeleteMedicineRequest) (*DeleteMedicineResponse, error)
	GetMedicineList(ctx context.Context, req *GetMedicineListRequest) (*GetMedicineListResponse, error)
}

type AddMedicineRequest struct {
	Name         string `validate:"required"`
	CategoriesID []uint
	Items        uint   `validate:"required"`
	ItemsUnit    string `validate:"required"`
	Expires      string `validate:"required"`
}

type AddMedicineResponse struct {
	ID uint
}

type UpdateMedicineRequest struct {
	ID           uint
	Name         string `validate:"required"`
	CategoriesID []uint
	Items        uint   `validate:"required"`
	ItemsUnit    string `validate:"required"`
	Expires      string `validate:"required"`
}

type UpdateMedicineResponse struct {
	ID uint
}

type DeleteMedicineRequest struct {
	ID uint
}

type DeleteMedicineResponse struct {
}

type GetMedicineListRequest struct {
}

type MedicineListItem struct {
	ID        uint
	Name      string
	Items     uint
	ItemsUnit string
	Expires   string
}

type GetMedicineListResponse struct {
	List []*MedicineListItem
}
