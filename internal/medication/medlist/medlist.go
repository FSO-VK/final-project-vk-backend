// Package medlist is a package for medicine list operations
package medlist

import (
	"github.com/google/uuid"
)

// MedicineList is a list of medicines.
type MedicineList struct {
	ID          uint
	UserID      uuid.UUID
	MedicinesID []uint
}
