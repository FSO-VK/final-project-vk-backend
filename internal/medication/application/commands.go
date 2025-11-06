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
	Group               []string
	ManufacturerName    string
	ManufacturerCountry string
	ActiveSubstance     []ActiveSubstance
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
	Group               []string
	ManufacturerName    string
	ManufacturerCountry string
	ActiveSubstance     []ActiveSubstance
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
		Group:               GroupsToStrings(m.GetGroup()),
		ManufacturerName:    m.GetManufacturer().GetName(),
		ManufacturerCountry: m.GetManufacturer().GetCountry(),
		ActiveSubstance:     convertToActiveSubstance(m.GetActiveSubstance()),
		Expires:             m.GetExpirationDate().Format(time.DateOnly),
		Release:             release,
		Commentary:          m.GetCommentary().GetCommentary(),
	}
}

func convertToActiveSubstance(substances []medication.ActiveSubstance) []ActiveSubstance {
	result := make([]ActiveSubstance, len(substances))
	for i, v := range substances {
		result[i] = ActiveSubstance{
			Name:  v.GetName(),
			Value: v.GetDose().GetValue(),
			Unit:  v.GetDose().GetUnit().String(),
		}
	}
	return result
}

// MapActiveSubstanceToDraft maps ActiveSubstance to ActiveSubstanceDraft.
func MapActiveSubstanceToDraft(substances []ActiveSubstance) []medication.ActiveSubstanceDraft {
	result := make([]medication.ActiveSubstanceDraft, len(substances))
	for i, substance := range substances {
		result[i] = medication.ActiveSubstanceDraft{
			Name:  substance.Name,
			Value: substance.Value,
			Unit:  substance.Unit,
		}
	}
	return result
}

// GroupsToStrings converts []medication.Group to []string.
func GroupsToStrings(groups []medication.Group) []string {
	result := make([]string, len(groups))
	for i, v := range groups {
		result[i] = string(v)
	}
	return result
}
