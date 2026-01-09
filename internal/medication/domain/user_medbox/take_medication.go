package usermedbox

import (
	"time"

	"github.com/google/uuid"
)

func (mb *UserMedbox) TakeMedication(id uuid.UUID, amount float32) error {
	medication := mb.getMedicationByID(id)
	if medication == nil {
		return ErrNoMedication
	}

	updatedBasicEvent, err := NewBasicEvent()
	if err != nil {
		return err
	}
	updatedEvent := MedicationUpdatedEvent{
		BasicEvent: updatedBasicEvent,
		medication: *medication,
	}

	ranOutBasicEvent, err := NewBasicEvent()
	if err != nil {
		return err
	}
	ranOutEvent := MedicationRanOutEvent{
		BasicEvent:   ranOutBasicEvent,
		medicationID: medication.ID(),
	}

	err = medication.Take(amount)
	if err != nil {
		return err
	}

	mb.events = append(mb.events, updatedEvent)

	if medication.IsRanOut() {
		mb.events = append(mb.events, ranOutEvent)
	}

	mb.updatedAt = time.Now()
	return nil
}
