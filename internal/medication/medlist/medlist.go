package medlist

import (
	"github.com/google/uuid"
)

type MedicineList struct {
	ID          uint
	UserID      uuid.UUID
	MedicinesID []uint
}
