// Package medication is a domain for medication
package medication

import (
	"time"
)

// Medication represents a medication entity.
type Medication struct {
	ID           uint
	Name         string
	CategoriesID []uint
	Items        uint      // количество таблеток/ампул/капсул и тп.
	ItemsUnit    string    // единица измерения
	Expires      time.Time // срок годности
}

// NewMedication creates a new medication.
func NewMedication(
	name string,
	items uint,
	categoriesID []uint,
	itemsUnit string,
	expires time.Time,
) *Medication {
	return &Medication{
		ID:           0, // could not be nill, but will be set later
		Name:         name,
		CategoriesID: categoriesID,
		Items:        items,
		ItemsUnit:    itemsUnit,
		Expires:      expires,
	}
}
