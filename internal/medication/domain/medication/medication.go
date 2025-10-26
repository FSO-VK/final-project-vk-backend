// Package medication is a domain for medication
package medication

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Medication represents a medication entity.
type Medication struct {
	ID                uuid.UUID
	Name              MedicationName
	InternationalName MedicationInternationalName
	Group             MedicationGroup
	Manufacturer      MedicationManufacturer
	ReleaseForm       MedicationReleaseForm
	Amount            MedicationAmount

	// Точно ли комментарий относится к лекарству?
	// Скорее это относится к записи в аптечке человека.
	Commentary MedicationCommentary

	ActiveSubstance MedicationActiveSubstance
	ReleaseDate     time.Time // дата выпуска
	ExpirationDate  time.Time // срок годности
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// NewMedication creates a new medication.
func NewMedication(
	id uuid.UUID,
	name MedicationName,
	internationalName MedicationInternationalName,
	group MedicationGroup,
	manufacturer MedicationManufacturer,
	releaseForm MedicationReleaseForm,
	amount MedicationAmount,
	activeSubstance MedicationActiveSubstance,
	releaseDate time.Time,
	expirationDate time.Time,
	commentary MedicationCommentary,
	createdAt time.Time,
	updatedAt time.Time,
) *Medication {
	return &Medication{
		ID:                id,
		Name:              name,
		InternationalName: internationalName,
		Group:             group,
		Manufacturer:      manufacturer,
		ReleaseForm:       releaseForm,
		Amount:            amount,
		ActiveSubstance:   activeSubstance,
		ReleaseDate:       releaseDate,
		ExpirationDate:    expirationDate,
		Commentary:        commentary,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}
}

// MedicationDraft represents a medication draft entity
// that uses built-in types.
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
	Manufacturer      MedicationManufacturerDraft

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
	name        MedicationName
	releaseForm MedicationReleaseForm
	amount      MedicationAmount
}

type optionalFields struct {
	internationalName MedicationInternationalName
	group             MedicationGroup
	manufacturer      MedicationManufacturer
	activeSubstance   MedicationActiveSubstance
	commentary        MedicationCommentary
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

func NewMedicationParse(draft MedicationDraft) (*Medication, error) {
	required, err := validateRequired(draft)
	if err != nil {
		return nil, err
	}

	optional, err := validateOptional(draft)
	if err != nil {
		return nil, err
	}

	return &Medication{
		ID:                required.id,
		Name:              required.name,
		ReleaseForm:       required.releaseForm,
		Amount:            required.amount,
		InternationalName: optional.internationalName,
		Group:             optional.group,
		Manufacturer:      optional.manufacturer,
		ActiveSubstance:   optional.activeSubstance,
		Commentary:        optional.commentary,
		ReleaseDate:       draft.ReleaseDate,
		ExpirationDate:    draft.ExpirationDate,
		CreatedAt:         draft.CreatedAt,
		UpdatedAt:         draft.UpdatedAt,
	}, nil
}

func (m *Medication) GetID() uuid.UUID { return m.ID }

func (m *Medication) SetName(name MedicationName) {
	m.Name = name
}

func (m *Medication) SetInternationalName(name MedicationInternationalName) {
	m.InternationalName = name
}

func (m *Medication) SetGroup(group MedicationGroup) {
	m.Group = group
}

func (m *Medication) SetManufacturer(manufacturer MedicationManufacturer) {
	m.Manufacturer = manufacturer
}

func (m *Medication) SetReleaseForm(form MedicationReleaseForm) {
	m.ReleaseForm = form
}

func (m *Medication) SetAmount(amount MedicationAmount) {
	m.Amount = amount
}

func (m *Medication) SetCommentary(commentary MedicationCommentary) {
	m.Commentary = commentary
}

func (m *Medication) SetActiveSubstance(substance MedicationActiveSubstance) {
	m.ActiveSubstance = substance
}

func (m *Medication) SetReleaseDate(date time.Time) {
	m.ReleaseDate = date
}

func (m *Medication) SetExpirationDate(date time.Time) {
	m.ExpirationDate = date
}

func (m *Medication) SetUpdatedAt(date time.Time) {
	m.UpdatedAt = date
}

func (m *Medication) GetName() MedicationName {
	return m.Name
}

func (m *Medication) GetInternationalName() MedicationInternationalName {
	return m.InternationalName
}

func (m *Medication) GetGroup() MedicationGroup {
	return m.Group
}

func (m *Medication) GetManufacturer() MedicationManufacturer {
	return m.Manufacturer
}

func (m *Medication) GetReleaseForm() MedicationReleaseForm {
	return m.ReleaseForm
}

func (m *Medication) GetAmount() MedicationAmount {
	return m.Amount
}

func (m *Medication) GetCommentary() MedicationCommentary {
	return m.Commentary
}

func (m *Medication) GetActiveSubstance() MedicationActiveSubstance {
	return m.ActiveSubstance
}

func (m *Medication) GetReleaseDate() time.Time {
	return m.ReleaseDate
}

func (m *Medication) GetExpirationDate() time.Time {
	return m.ExpirationDate
}

func (m *Medication) GetCreatedAt() time.Time {
	return m.CreatedAt
}

func (m *Medication) GetUpdatedAt() time.Time {
	return m.UpdatedAt
}
