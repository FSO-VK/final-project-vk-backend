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
		ReleaseDate:       releaseDate,
		ExpirationDate:    expirationDate,
		Commentary:        commentary,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}
}

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

	releaseForm, ok := ReleaseFormString[draft.ReleaseForm]
	if !ok {
		allErrors = errors.Join(allErrors, ErrInvalidReleaseForm)
	}

	amountUnit, ok := UnitString[draft.AmountUnit]
	if !ok {
		allErrors = errors.Join(allErrors, ErrInvalidUnit)
	}

	amount, err := NewMedicationAmount(draft.AmountValue, amountUnit)
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

	activeSubstanceUnit, ok := UnitString[draft.ActiveSubstanceDoseUnit]
	if !ok {
		allErrors = errors.Join(allErrors, ErrInvalidUnit)
	}

	activeSubstance, err := NewMedicationActiveSubstance(
		draft.ActiveSubstanceName,
		draft.ActiveSubstanceDoseValue,
		activeSubstanceUnit,
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
