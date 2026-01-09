package usermedbox

import (
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/user_medbox/medication"
	"github.com/google/uuid"
)

type MedicationUpdatedEvent struct {
	BasicEvent
	medication medication.Medication
}

func (mue MedicationUpdatedEvent) Name() string {
	return "user_medbox.medication_updated"
}

type MedicationRanOutEvent struct {
	BasicEvent
	medicationID uuid.UUID
}

func (mroe MedicationRanOutEvent) Name() string {
	return "user_medbox.medication_ran_out"
}

func (mb *UserMedbox) UpdateMedicationInfo(id uuid.UUID, info medication.MedicationInfo) error {
	medication := mb.getMedicationByID(id)
	if medication == nil {
		return ErrNoMedication
	}

	medication.UpdateInfo(info)

	updatedBasicEvent, err := NewBasicEvent()
	if err != nil {
		return err
	}
	updatedEvent := MedicationUpdatedEvent{
		BasicEvent: updatedBasicEvent,
		medication: *medication,
	}
	mb.events = append(mb.events, updatedEvent)

	if medication.IsRanOut() {

		ranOutBasicEvent, err := NewBasicEvent()
		if err != nil {
			return err
		}

		ranOutEvent := MedicationRanOutEvent{
			BasicEvent:   ranOutBasicEvent,
			medicationID: medication.ID(),
		}
		mb.events = append(mb.events, ranOutEvent)
	}

	mb.updatedAt = time.Now()
	return nil
}
