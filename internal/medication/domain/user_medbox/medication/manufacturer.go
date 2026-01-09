package medication

import (
	"errors"
	"fmt"

	"github.com/FSO-VK/final-project-vk-backend/pkg/validation"
)

const (
	maxManufacturerNameLength = 200
	maxCountryLength          = 100
)

var ErrInvalidManufacturer = errors.New("invalid manufacturer")

// ManufacturerDraft represents the raw input data
// structure for a medication manufacturer.
type ManufacturerDraft struct {
	Name    string
	Country string
}

// Manufacturer is a VO representing manufacturer information.
type Manufacturer struct {
	name    string
	country string
}

// NewManufacturer creates validated medication manufacturer.
func NewManufacturer(draft ManufacturerDraft) (Manufacturer, error) {
	err := errors.Join(
		validation.MaxLength(draft.Name, maxManufacturerNameLength),
		validation.MaxLength(draft.Country, maxCountryLength),
	)
	if err != nil {
		return Manufacturer{}, fmt.Errorf("%w: %w", ErrInvalidManufacturer, err)
	}
	return Manufacturer{
		name:    draft.Name,
		country: draft.Country,
	}, nil
}

// Name returns the manufacturer's name.
func (m Manufacturer) Name() string {
	return m.name
}

// Country returns the manufacturer's country.
func (m Manufacturer) Country() string {
	return m.country
}
