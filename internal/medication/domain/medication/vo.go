package medication

// This file contains value objects (VO) for medication.

import (
	"errors"
	"fmt"

	"github.com/FSO-VK/final-project-vk-backend/pkg/validation"
)

const (
	maxNameLength       = 100
	maxGroupNameLength  = 100
	maxCountryLength    = 100
	maxCommentaryLength = 1000
)

// Errors of the medication domain VO's.
var (
	ErrInvalidID                = errors.New("invalid id")
	ErrInvalidName              = errors.New("invalid name")
	ErrInvalidInternationalName = errors.New("invalid international name")
	ErrInvalidGroup             = errors.New("invalid group")
	ErrInvalidManufacturer      = errors.New("invalid manufacturer")
	ErrInvalidReleaseForm       = errors.New("invalid release form")
	ErrInvalidUnit              = errors.New("invalid unit")
	ErrInvalidAmount            = errors.New("invalid amount")
	ErrInvalidActiveSubstance   = errors.New("invalid active substance")
	ErrInvalidCommentary        = errors.New("invalid commentary")
	ErrInvalidExpirationTime    = errors.New("invalid expiration time")
	ErrInvalidDateRange         = errors.New("expiration date must be greater than release date")
)

// Name is a Value Object representing the name of a medication.
type Name string

// NewMedicationName creates validated medication name.
func NewMedicationName(name string) (Name, error) {
	err := errors.Join(
		validation.Required(name),
		validation.MaxLength(name, 200),
	)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrInvalidName, err)
	}
	return Name(name), nil
}

// GetName returns the string value of the medication name.
func (n Name) GetName() string {
	return string(n)
}

// InternationalName is a Value Object representing the international non-proprietary name of a medication.
type InternationalName string

// NewMedicationInternationalName creates validated medication international name.
func NewMedicationInternationalName(name string) (InternationalName, error) {
	err := errors.Join(
		validation.MaxLength(name, maxNameLength),
	)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrInvalidInternationalName, err)
	}
	return InternationalName(name), nil
}

// GetInternationalName returns the international name of the medication.
func (n InternationalName) GetInternationalName() string {
	return string(n)
}

// Group is a Value Object representing
// the therapeutic or pharmacological group of a medication.
type Group string

// NewMedicationGroup creates validated medication group.
func NewMedicationGroup(groups []string) ([]Group, error) {
	result := []Group{}
	for _, group := range groups {
		err := errors.Join(
			validation.MaxLength(group, maxGroupNameLength),
		)
		if err != nil {
			return []Group{}, fmt.Errorf("%w: %w", ErrInvalidGroup, err)
		}
		result = append(result, Group(group))
	}
	return result, nil
}

// GetGroup returns the medication group name.
func (g Group) GetGroup() string {
	return string(g)
}

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

// NewMedicationManufacturer creates validated medication manufacturer.
func NewMedicationManufacturer(name string, country string) (Manufacturer, error) {
	err := errors.Join(
		validation.MaxLength(name, maxNameLength),
		validation.MaxLength(country, maxCountryLength),
	)
	if err != nil {
		return Manufacturer{}, fmt.Errorf("%w: %w", ErrInvalidManufacturer, err)
	}
	return Manufacturer{
		name:    name,
		country: country,
	}, nil
}

// GetName returns the manufacturer's name.
func (m Manufacturer) GetName() string {
	return m.name
}

// GetCountry returns the manufacturer's country.
func (m Manufacturer) GetCountry() string {
	return m.country
}

// ReleaseForm is a Value Object representing the physical form
// in which the medication is released.
type ReleaseForm int

// Enum values of possible medication release form.
const (
	UnknownForm ReleaseForm = iota
	Tablet
	Capsule
	Injection
	Ointment
	Syrup
	Drops
	Inhalation
	Patch
)

// Unexported global variable.
//
//nolint:gochecknoglobals
var releaseFormToString = map[string]ReleaseForm{
	"tablet":     Tablet,
	"capsule":    Capsule,
	"injection":  Injection,
	"ointment":   Ointment,
	"syrup":      Syrup,
	"drops":      Drops,
	"inhalation": Inhalation,
	"patch":      Patch,
}

// Unexported global variable.
//
//nolint:gochecknoglobals
var stringToReleaseForm = map[ReleaseForm]string{
	UnknownForm: "unknown form",
	Tablet:      "tablet",
	Capsule:     "capsule",
	Injection:   "injection",
	Ointment:    "ointment",
	Syrup:       "syrup",
	Drops:       "drops",
	Inhalation:  "inhalation",
	Patch:       "patch",
}

// NewMedicationReleaseForm creates validated medication release form.
func NewMedicationReleaseForm(form string) (ReleaseForm, error) {
	err := errors.Join(
		validation.Required(form),
	)
	if err != nil {
		return UnknownForm, fmt.Errorf("%w: %w", ErrInvalidReleaseForm, err)
	}
	return releaseFormToString[form], nil
}

// String returns the string representation of the release form.
func (f ReleaseForm) String() string {
	return stringToReleaseForm[f]
}

// Unit is a VO representing the unit of measurement
// for medication quantities.
type Unit int

// Enum values of possible medication unit.
const (
	UnsetUnit Unit = iota
	UnknownUnit
	Piece
	Gram
	Milligram
	Milliliter
)

// Unexported global variable.
//
//nolint:gochecknoglobals
var stringToUnit = map[string]Unit{
	"piece":        Piece,
	"gram":         Gram,
	"milligram":    Milligram,
	"milliliter":   Milliliter,
	"pcs":          Piece,
	"g":            Gram,
	"mg":           Milligram,
	"ml":           Milliliter,
	"":             UnsetUnit,
	"unknown unit": UnknownUnit,
}

// Unexported global variable.
//
//nolint:gochecknoglobals
var unitToString = map[Unit]string{
	UnsetUnit:   "",
	UnknownUnit: "unknown unit",
	Piece:       "piece",
	Gram:        "gram",
	Milligram:   "milligram",
	Milliliter:  "milliliter",
}

// NewMedicationUnit creates validated medication unit.
func NewMedicationUnit(unit string) (Unit, error) {
	if unit == "" {
		return UnsetUnit, nil
	}
	_, ok := stringToUnit[unit]
	if !ok {
		return UnknownUnit, fmt.Errorf("%w: unknown unit", ErrInvalidUnit)
	}
	return stringToUnit[unit], nil
}

// String returns the string representation of the medication unit.
func (u Unit) String() string {
	return unitToString[u]
}

// Amount is a VO representing the quantity of a medication
// with its unit of measurement.
type Amount struct {
	value float32
	unit  Unit
}

// NewMedicationAmount creates validated medication amount.
func NewMedicationAmount(value float32, unit string) (Amount, error) {
	medicationUnit, err := NewMedicationUnit(unit)

	err = errors.Join(
		err,
		validation.Positive(value),
	)
	if err != nil {
		return Amount{}, fmt.Errorf("%w: %w", ErrInvalidAmount, err)
	}

	return Amount{
		value: value,
		unit:  medicationUnit,
	}, nil
}

// GetValue returns the numeric value of the amount.
func (a Amount) GetValue() float32 {
	return a.value
}

// GetUnit returns the unit of measurement for the amount.
func (a Amount) GetUnit() Unit {
	return a.unit
}

// Commentary is a Value Object representing additional notes or comments about a medication.
type Commentary string

// NewMedicationCommentary creates validated medication commentary.
func NewMedicationCommentary(commentary string) (Commentary, error) {
	err := errors.Join(
		validation.MaxLength(commentary, maxCommentaryLength),
	)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrInvalidCommentary, err)
	}
	return Commentary(commentary), nil
}

// GetCommentary returns the text of the medication commentary.
func (c Commentary) GetCommentary() string {
	return string(c)
}

// ActiveSubstance is a VO representing the active
// pharmaceutical substance and its dosage.
type ActiveSubstance struct {
	name string
	dose Amount
}

// NewMedicationActiveSubstance creates validated medication active substance.
func NewMedicationActiveSubstance(
	activeSubstances []ActiveSubstanceDraft,
) ([]ActiveSubstance, error) {
	result := []ActiveSubstance{}
	for _, activeSubstance := range activeSubstances {
		dose, err := NewMedicationAmount(
			activeSubstance.Value,
			activeSubstance.Unit,
		)
		err = errors.Join(
			validation.MaxLength(activeSubstance.Name, 200),
			err,
		)
		if err != nil {
			return []ActiveSubstance{}, fmt.Errorf("%w: %w", ErrInvalidActiveSubstance, err)
		}
		result = append(result, ActiveSubstance{
			name: activeSubstance.Name,
			dose: dose,
		})
	}
	return result, nil
}

// GetName returns the name of the active substance.
func (a ActiveSubstance) GetName() string { return a.name }

// GetDose returns the dose amount of the active substance.
func (a ActiveSubstance) GetDose() Amount { return a.dose }
