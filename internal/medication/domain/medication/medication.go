// Package medication is a domain for medication
package medication

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Medication represents a medication entity.
type Medication struct {
	id                uuid.UUID
	name              Name
	internationalName InternationalName
	group             Group
	manufacturer      Manufacturer
	releaseForm       ReleaseForm
	amount            Amount

	// Точно ли комментарий относится к лекарству?
	// Скорее это относится к записи в аптечке человека.
	commentary Commentary

	activeSubstance ActiveSubstance
	releaseDate     time.Time // дата выпуска
	expirationDate  time.Time // срок годности
	createdAt       time.Time
	updatedAt       time.Time
}

// NewMedication creates a new medication.
func NewMedication(
	id uuid.UUID,
	name Name,
	internationalName InternationalName,
	group Group,
	manufacturer Manufacturer,
	releaseForm ReleaseForm,
	amount Amount,
	activeSubstance ActiveSubstance,
	releaseDate time.Time,
	expirationDate time.Time,
	commentary Commentary,
	createdAt time.Time,
	updatedAt time.Time,
) *Medication {
	return &Medication{
		id:                id,
		name:              name,
		internationalName: internationalName,
		group:             group,
		manufacturer:      manufacturer,
		releaseForm:       releaseForm,
		amount:            amount,
		activeSubstance:   activeSubstance,
		releaseDate:       releaseDate,
		expirationDate:    expirationDate,
		commentary:        commentary,
		createdAt:         createdAt,
		updatedAt:         updatedAt,
	}
}

// MedicationDraft represents a medication draft entity
// that uses built-in types.
//
// Revive consider to name this struct as Draft
// but it's more clear to name the struct as-is.
//
//nolint:revive
type MedicationDraft struct {
	// required fields
	ID             uuid.UUID
	Name           string
	ReleaseForm    string
	AmountValue    float32
	AmountUnit     string
	ExpirationDate time.Time

	// optional fields, nullable
	InternationalName string
	Group             string
	Manufacturer      ManufacturerDraft

	ActiveSubstanceName      string
	ActiveSubstanceDoseValue float32
	ActiveSubstanceDoseUnit  string

	Commentary string

	ReleaseDate time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

type requiredFields struct {
	id          uuid.UUID
	name        Name
	releaseForm ReleaseForm
	amount      Amount
}

type optionalFields struct {
	internationalName InternationalName
	group             Group
	manufacturer      Manufacturer
	activeSubstance   ActiveSubstance
	commentary        Commentary
}

func validateRequired(draft MedicationDraft) (requiredFields, error) {
	var allErrors error

	if draft.ID == uuid.Nil {
		return requiredFields{}, ErrInvalidID
	}

	name, err := NewMedicationName(draft.Name)
	allErrors = errors.Join(allErrors, err)

	releaseForm, err := NewMedicationReleaseForm(draft.ReleaseForm)
	allErrors = errors.Join(allErrors, err)

	amount, err := NewMedicationAmount(draft.AmountValue, draft.AmountUnit)
	allErrors = errors.Join(allErrors, err)

	if draft.ExpirationDate.IsZero() {
		allErrors = errors.Join(allErrors, ErrInvalidExpirationTime)
	}

	if allErrors != nil {
		return requiredFields{}, allErrors
	}

	return requiredFields{
		id:          draft.ID,
		name:        name,
		releaseForm: releaseForm,
		amount:      amount,
	}, nil
}

func validateOptional(draft MedicationDraft) (optionalFields, error) {
	var allErrors error

	internationalName, err := NewMedicationInternationalName(draft.InternationalName)
	allErrors = errors.Join(allErrors, err)

	group, err := NewMedicationGroup(draft.Group)
	allErrors = errors.Join(allErrors, err)

	manufacturer, err := NewMedicationManufacturer(
		draft.Manufacturer.Name,
		draft.Manufacturer.Country,
	)
	allErrors = errors.Join(allErrors, err)

	activeSubstance, err := NewMedicationActiveSubstance(
		draft.ActiveSubstanceName,
		draft.ActiveSubstanceDoseValue,
		draft.ActiveSubstanceDoseUnit,
	)
	allErrors = errors.Join(allErrors, err)

	commentary, err := NewMedicationCommentary(draft.Commentary)
	allErrors = errors.Join(allErrors, err)

	if draft.ReleaseDate.After(draft.ExpirationDate) {
		allErrors = errors.Join(allErrors, ErrInvalidDateRange)
	}

	if allErrors != nil {
		return optionalFields{}, allErrors
	}

	return optionalFields{
		internationalName: internationalName,
		group:             group,
		manufacturer:      manufacturer,
		activeSubstance:   activeSubstance,
		commentary:        commentary,
	}, nil
}

// Parse creates a new medication from a draft.
func Parse(draft MedicationDraft) (*Medication, error) {
	required, err := validateRequired(draft)
	if err != nil {
		return nil, err
	}

	optional, err := validateOptional(draft)
	if err != nil {
		return nil, err
	}

	return &Medication{
		id:                required.id,
		name:              required.name,
		releaseForm:       required.releaseForm,
		amount:            required.amount,
		internationalName: optional.internationalName,
		group:             optional.group,
		manufacturer:      optional.manufacturer,
		activeSubstance:   optional.activeSubstance,
		commentary:        optional.commentary,
		releaseDate:       draft.ReleaseDate,
		expirationDate:    draft.ExpirationDate,
		createdAt:         draft.CreatedAt,
		updatedAt:         draft.UpdatedAt,
	}, nil
}

// GetID returns the unique identifier of the medication.
func (m *Medication) GetID() uuid.UUID { return m.id }

// SetName updates the name Value Object of the medication.
func (m *Medication) SetName(name Name) {
	m.name = name
}

// SetInternationalName updates the international name Value Object of the medication.
func (m *Medication) SetInternationalName(name InternationalName) {
	m.internationalName = name
}

// SetGroup updates the group Value Object of the medication.
func (m *Medication) SetGroup(group Group) {
	m.group = group
}

// SetManufacturer updates the manufacturer Value Object of the medication.
func (m *Medication) SetManufacturer(manufacturer Manufacturer) {
	m.manufacturer = manufacturer
}

// SetReleaseForm updates the release form Value Object of the medication.
func (m *Medication) SetReleaseForm(form ReleaseForm) {
	m.releaseForm = form
}

// SetAmount updates the amount Value Object of the medication.
func (m *Medication) SetAmount(amount Amount) {
	m.amount = amount
}

// SetCommentary updates the commentary Value Object of the medication.
func (m *Medication) SetCommentary(commentary Commentary) {
	m.commentary = commentary
}

// SetActiveSubstance updates the active substance Value Object of the medication.
func (m *Medication) SetActiveSubstance(substance ActiveSubstance) {
	m.activeSubstance = substance
}

// SetReleaseDate updates the release date of the medication.
func (m *Medication) SetReleaseDate(date time.Time) {
	m.releaseDate = date
}

// SetExpirationDate updates the expiration date of the medication.
func (m *Medication) SetExpirationDate(date time.Time) {
	m.expirationDate = date
}

// SetUpdatedAt updates the last modification timestamp of the medication.
func (m *Medication) SetUpdatedAt(date time.Time) {
	m.updatedAt = date
}

// GetName returns the name Value Object of the medication.
func (m *Medication) GetName() Name {
	return m.name
}

// GetInternationalName returns the international name Value Object of the medication.
func (m *Medication) GetInternationalName() InternationalName {
	return m.internationalName
}

// GetGroup returns the group Value Object of the medication.
func (m *Medication) GetGroup() Group {
	return m.group
}

// GetManufacturer returns the manufacturer Value Object of the medication.
func (m *Medication) GetManufacturer() Manufacturer {
	return m.manufacturer
}

// GetReleaseForm returns the release form Value Object of the medication.
func (m *Medication) GetReleaseForm() ReleaseForm {
	return m.releaseForm
}

// GetAmount returns the amount Value Object of the medication.
func (m *Medication) GetAmount() Amount {
	return m.amount
}

// GetCommentary returns the commentary Value Object of the medication.
func (m *Medication) GetCommentary() Commentary {
	return m.commentary
}

// GetActiveSubstance returns the active substance Value Object of the medication.
func (m *Medication) GetActiveSubstance() ActiveSubstance {
	return m.activeSubstance
}

// GetReleaseDate returns the release date of the medication.
func (m *Medication) GetReleaseDate() time.Time {
	return m.releaseDate
}

// GetExpirationDate returns the expiration date of the medication.
func (m *Medication) GetExpirationDate() time.Time {
	return m.expirationDate
}

// GetCreatedAt returns the creation timestamp of the medication.
func (m *Medication) GetCreatedAt() time.Time {
	return m.createdAt
}

// GetUpdatedAt returns the last modification timestamp of the medication.
func (m *Medication) GetUpdatedAt() time.Time {
	return m.updatedAt
}
