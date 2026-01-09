package medication

import (
	"errors"
	"time"
)

// ErrInvalidExpirationDate is an error occured when expiration date is before release date.
var ErrInvalidExpirationDate = errors.New("invalid expiration date")

// MedicationInfoDraft is a helper struct for constructing MedicationInfo.
type MedicationInfoDraft struct {
	Name              string
	InternationalName string
	Groups            []string
	Manufacturer      ManufacturerDraft
	ReleaseForm       string
	Amount            AmountDraft
	Commentary        string
	ActiveSubstances  []ActiveSubstanceDraft
	ReleaseDate       time.Time
	ExpirationDate    time.Time
}

// MedicationInfo is a user editable medication common info.
type MedicationInfo struct {
	name              Name
	internationalName InternationalName
	groups            []Group
	manufacturer      Manufacturer
	releaseForm       ReleaseForm
	amount            Amount
	commentary        Commentary
	activeSubstances  []ActiveSubstance
	releaseDate       time.Time
	expirationDate    time.Time
}

// NewMedicationInfo creates a new instance of MedicationInfo.
func NewMedicationInfo(d MedicationInfoDraft) (MedicationInfo, error) {
	var allErrors error

	name, err := NewName(d.Name)
	allErrors = errors.Join(allErrors, err)

	internationalName, err := NewInternationalName(d.InternationalName)
	allErrors = errors.Join(allErrors, err)

	groups := make([]Group, 0, len(d.Groups))
	for _, groupDraft := range d.Groups {
		group, err := NewGroup(groupDraft)
		allErrors = errors.Join(allErrors, err)
		groups = append(groups, group)
	}

	manufacturer, err := NewManufacturer(d.Manufacturer)
	allErrors = errors.Join(allErrors, err)

	releaseForm, err := NewReleaseForm(d.ReleaseForm)
	allErrors = errors.Join(allErrors, err)

	amount, err := NewAmount(d.Amount)
	allErrors = errors.Join(allErrors, err)

	commentary, err := NewCommentary(d.Commentary)
	allErrors = errors.Join(allErrors, err)

	activeSubstances := make([]ActiveSubstance, 0, len(d.ActiveSubstances))
	for _, activeSubstanceDraft := range d.ActiveSubstances {
		activeSubstance, err := NewActiveSubstance(activeSubstanceDraft)
		allErrors = errors.Join(allErrors, err)
		activeSubstances = append(activeSubstances, activeSubstance)
	}

	if d.ExpirationDate.Before(d.ReleaseDate) {
		allErrors = errors.Join(allErrors, ErrInvalidExpirationDate)
	}

	if allErrors != nil {
		return MedicationInfo{}, allErrors
	}

	return MedicationInfo{
		name:              name,
		internationalName: internationalName,
		groups:            groups,
		manufacturer:      manufacturer,
		releaseForm:       releaseForm,
		amount:            amount,
		commentary:        commentary,
		activeSubstances:  activeSubstances,
		releaseDate:       d.ReleaseDate,
		expirationDate:    d.ExpirationDate,
	}, nil
}

// Name returns the name Value Object of the medication.
func (m *MedicationInfo) Name() Name {
	return m.name
}

// InternationalName returns the international name Value Object of the medication.
func (m *MedicationInfo) InternationalName() InternationalName {
	return m.internationalName
}

// Group returns the groups of the medication.
func (m *MedicationInfo) Groups() []Group {
	return m.groups
}

// Manufacturer returns the manufacturer Value Object of the medication.
func (m *MedicationInfo) Manufacturer() Manufacturer {
	return m.manufacturer
}

// ReleaseForm returns the release form Value Object of the medication.
func (m *MedicationInfo) ReleaseForm() ReleaseForm {
	return m.releaseForm
}

// Amount returns the amount Value Object of the medication.
func (m *MedicationInfo) Amount() Amount {
	return m.amount
}

// Commentary returns the commentary Value Object of the medication.
func (m *MedicationInfo) Commentary() Commentary {
	return m.commentary
}

// ActiveSubstance returns the active substance of the medication.
func (m *MedicationInfo) ActiveSubstances() []ActiveSubstance {
	return m.activeSubstances
}

// ReleaseDate returns the release date of the medication.
func (m *MedicationInfo) ReleaseDate() time.Time {
	return m.releaseDate
}

// ExpirationDate returns the expiration date of the medication.
func (m *MedicationInfo) ExpirationDate() time.Time {
	return m.expirationDate
}
