package application

import (
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medication"
)

// CommandBase contains common fields for commands.
type CommandBase struct {
	Name                string `validate:"required"`
	InternationalName   string
	AmountValue         float32 `validate:"required,gte=0"`
	AmountUnit          string  `validate:"required"`
	ReleaseForm         string  `validate:"required"`
	Group               string
	ManufacturerName    string
	ManufacturerCountry string
	ActiveSubstanceName string
	ActiveSubstanceDose float32
	ActiveSubstanceUnit string
	Expires             string `validate:"required"`
	Release             string
	Commentary          string
}

// ResponseBase contains common fields for responses.
type ResponseBase struct {
	ID                  string
	Name                string
	InternationalName   string
	AmountValue         float32
	AmountUnit          string
	ReleaseForm         string
	Group               string
	ManufacturerName    string
	ManufacturerCountry string
	ActiveSubstanceName string
	ActiveSubstanceDose float32
	ActiveSubstanceUnit string
	Expires             string
	Release             string
	Commentary          string
}

// responseBaseMapper maps medication.Medication to ResponseBase.
func responseBaseMapper(m *medication.Medication) ResponseBase {
	var release string
	if !m.GetReleaseDate().IsZero() {
		release = m.GetReleaseDate().Format(time.DateOnly)
	}
	return ResponseBase{
		ID:                  m.GetID().String(),
		Name:                m.GetName().GetName(),
		InternationalName:   m.GetInternationalName().GetInternationalName(),
		AmountValue:         m.GetAmount().GetValue(),
		AmountUnit:          m.GetAmount().GetUnit().String(),
		ReleaseForm:         m.GetReleaseForm().String(),
		Group:               m.GetGroup().GetGroup(),
		ManufacturerName:    m.GetManufacturer().GetName(),
		ManufacturerCountry: m.GetManufacturer().GetCountry(),
		ActiveSubstanceName: m.GetActiveSubstance().GetName(),
		ActiveSubstanceDose: m.GetActiveSubstance().GetDose().GetValue(),
		ActiveSubstanceUnit: m.GetActiveSubstance().GetDose().GetUnit().String(),
		Expires:             m.GetExpirationDate().Format(time.DateOnly),
		Release:             release,
		Commentary:          m.GetCommentary().GetCommentary(),
	}
}
