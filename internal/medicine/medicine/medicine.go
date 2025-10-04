package medicine

import (
	"time"
)

type Medicine struct {
	ID           uint
	Name         string
	CategoriesID []uint
	Items        uint      // количество таблеток/ампул/капсул и тп.
	ItemsUnit    string    // единица измерения
	Expires      time.Time // срок годности
}

func NewMedicine(
	name string,
	items uint,
	itemsUnit string,
	expires time.Time,
) *Medicine {
	return &Medicine{
		Name:         name,
		Items:        items,
		ItemsUnit:    itemsUnit,
		Expires:      expires,
	}
}
