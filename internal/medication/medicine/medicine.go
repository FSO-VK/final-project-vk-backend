// Package medicine is a package for medicines
package medicine

import (
	"time"
)

// Medicine represents a medicine.
type Medicine struct {
	ID           uint
	Name         string
	CategoriesID []uint
	Items        uint      // количество таблеток/ампул/капсул и тп.
	ItemsUnit    string    // единица измерения
	Expires      time.Time // срок годности
}

// NewMedicine creates a new medicine.
func NewMedicine(
	name string,
	items uint,
	categoriesID []uint,
	itemsUnit string,
	expires time.Time,
) *Medicine {
	return &Medicine{
		Name:         name,
		CategoriesID: categoriesID,
		Items:        items,
		ItemsUnit:    itemsUnit,
		Expires:      expires,
	}
}
