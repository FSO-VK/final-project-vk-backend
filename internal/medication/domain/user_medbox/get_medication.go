package usermedbox

import (
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/user_medbox/medication"
	"github.com/google/uuid"
)

func (mb *UserMedbox) GetMedication(id uuid.UUID) (*medication.Medication, error) {
	medication := mb.getMedicationByID(id)
	if medication == nil {
		return nil, ErrNoMedication
	}
	return medication, nil
}
