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

type MedicationName string

func NewMedicationName(name string) (MedicationName, error) {
	err := errors.Join(
		validation.Required(name),
		validation.MaxLength(name, 200),
	)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrInvalidName, err)
	}
	return MedicationName(name), nil
}

func (n MedicationName) GetName() string {
	return string(n)
}

type MedicationInternationalName string

func NewMedicationInternationalName(name string) (MedicationInternationalName, error) {
	err := errors.Join(
		validation.MaxLength(name, maxNameLength),
	)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrInvalidInternationalName, err)
	}
	return MedicationInternationalName(name), nil
}

func (n MedicationInternationalName) GetInternationalName() string {
	return string(n)
}

type MedicationGroup string

func NewMedicationGroup(group string) (MedicationGroup, error) {
	err := errors.Join(
		validation.MaxLength(group, maxGroupNameLength),
	)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrInvalidGroup, err)
	}

	return MedicationGroup(group), nil
}

func (g MedicationGroup) GetGroup() string {
	return string(g)
}

type MedicationManufacturerDraft struct {
	Name    string
	Country string
}

type MedicationManufacturer struct {
	name    string
	country string
}

func NewMedicationManufacturer(name string, country string) (MedicationManufacturer, error) {
	err := errors.Join(
		validation.MaxLength(name, maxNameLength),
		validation.MaxLength(country, maxCommentaryLength),
	)
	if err != nil {
		return MedicationManufacturer{}, fmt.Errorf("%w: %w", ErrInvalidManufacturer, err)
	}
	return MedicationManufacturer{
		name:    name,
		country: country,
	}, nil
}

func (m MedicationManufacturer) GetName() string {
	return m.name
}

func (m MedicationManufacturer) GetCountry() string {
	return m.country
}

type MedicationReleaseForm int

const (
	UnknownForm MedicationReleaseForm = iota
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
var releaseFormToString = map[string]MedicationReleaseForm{
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
var stringToReleaseForm = map[MedicationReleaseForm]string{
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

func NewMedicationReleaseForm(form string) (MedicationReleaseForm, error) {
	err := errors.Join(
		validation.Required(form),
	)
	if err != nil {
		return UnknownForm, fmt.Errorf("%w: %w", ErrInvalidReleaseForm, err)
	}
	return releaseFormToString[form], nil
}

func (f MedicationReleaseForm) String() string {
	return stringToReleaseForm[f]
}

type MedicationUnit int

const (
	UnknownUnit MedicationUnit = iota
	Piece
	Gram
	Milligram
	Milliliter
)

// Unexported global variable.
//
//nolint:gochecknoglobals
var unitToString = map[string]MedicationUnit{
	"piece":      Piece,
	"gram":       Gram,
	"milligram":  Milligram,
	"milliliter": Milliliter,
	"pcs":        Piece,
	"g":          Gram,
	"mg":         Milligram,
	"ml":         Milliliter,
}

// Unexported global variable.
//
//nolint:gochecknoglobals
var stringToUnit = map[MedicationUnit]string{
	UnknownUnit: "unknown unit",
	Piece:       "piece",
	Gram:        "gram",
	Milligram:   "milligram",
	Milliliter:  "milliliter",
}

func NewMedicationUnit(unit string) (MedicationUnit, error) {
	err := errors.Join(
		validation.Required(unit),
	)
	if err != nil {
		return UnknownUnit, fmt.Errorf("%w: %w", ErrInvalidUnit, err)
	}
	return unitToString[unit], nil
}

func (u MedicationUnit) String() string {
	return stringToUnit[u]
}

type MedicationAmount struct {
	value float32
	unit  MedicationUnit
}

func NewMedicationAmount(value float32, unit string) (MedicationAmount, error) {
	medicationUnit, err := NewMedicationUnit(unit)

	err = errors.Join(
		err,
		validation.Positive(value),
	)
	if err != nil {
		return MedicationAmount{}, fmt.Errorf("%w: %w", ErrInvalidAmount, err)
	}

	return MedicationAmount{
		value: value,
		unit:  medicationUnit,
	}, nil
}

func (a MedicationAmount) GetValue() float32 {
	return a.value
}

func (a MedicationAmount) GetUnit() MedicationUnit {
	return a.unit
}

type MedicationCommentary string

func NewMedicationCommentary(commentary string) (MedicationCommentary, error) {
	err := errors.Join(
		validation.MaxLength(commentary, maxCommentaryLength),
	)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrInvalidCommentary, err)
	}
	return MedicationCommentary(commentary), nil
}

func (c MedicationCommentary) GetCommentary() string {
	return string(c)
}

type MedicationActiveSubstance struct {
	name string
	dose MedicationAmount
}

func NewMedicationActiveSubstance(
	name string,
	doseValue float32,
	doseUnit string,
) (MedicationActiveSubstance, error) {
	dose, err := NewMedicationAmount(
		doseValue,
		doseUnit,
	)
	err = errors.Join(
		validation.MaxLength(name, 200),
		err,
	)
	if err != nil {
		return MedicationActiveSubstance{}, fmt.Errorf("%w: %w", ErrInvalidActiveSubstance, err)
	}

	return MedicationActiveSubstance{
		name: name,
		dose: dose,
	}, nil
}

func (a MedicationActiveSubstance) GetName() string { return a.name }

func (a MedicationActiveSubstance) GetDose() MedicationAmount { return a.dose }
