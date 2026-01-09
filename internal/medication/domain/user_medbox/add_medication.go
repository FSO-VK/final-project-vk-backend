package usermedbox

import (
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/user_medbox/medication"
	"time"
)

type MedicationAddedEvent struct {
	BasicEvent

	medication medication.Medication
}

func (mae MedicationAddedEvent) Name() string {
	return "user_medbox.medication_added"
}

func (mb *UserMedbox) AddMedication(m medication.Medication) error {
	basicEvent, err := NewBasicEvent()
	if err != nil {
		return err
	}

	addedEvent := MedicationAddedEvent{
		BasicEvent: basicEvent,
		medication: m,
	}

	mb.medications[m.ID()] = &m
	mb.events = append(mb.events, addedEvent)
	mb.updatedAt = time.Now()
	return nil
}
