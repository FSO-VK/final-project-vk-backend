package usermedbox

import (
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/user_medbox/medication"
	"github.com/google/uuid"
)

type MedicationDeletedEvent struct {
	BasicEvent
	medication medication.Medication
}

func (mde MedicationDeletedEvent) Name() string {
	return "user_medbox.medication_deleted"
}

func (mb *UserMedbox) DeleteMedication(id uuid.UUID) error {
	medication := mb.getMedicationByID(id)
	if medication == nil {
		return ErrNoMedication
	}

	basicEvent, err := NewBasicEvent()
	if err != nil {
		return err
	}

	event := MedicationDeletedEvent{
		BasicEvent: basicEvent,
		medication: *medication,
	}

	delete(mb.medications, medication.ID())

	mb.events = append(mb.events, event)
	mb.updatedAt = time.Now()

	return nil
}
