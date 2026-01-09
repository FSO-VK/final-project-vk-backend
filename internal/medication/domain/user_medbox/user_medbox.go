package usermedbox

import (
	"errors"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/user_medbox/medication"
	"github.com/google/uuid"
)

var ErrInvalidUpdationTime = errors.New("updatedAt must be after createdAt")

var (
	ErrNoMedication = errors.New("no such medication")
)

type UserMedboxVersion uint

type UserMedbox struct {
	id          uuid.UUID
	userID      uuid.UUID
	medications map[uuid.UUID]*medication.Medication

	createdAt time.Time
	updatedAt time.Time
	version   UserMedboxVersion
	events    []Event
}

func New(id, userID uuid.UUID, medications []medication.Medication, createdAt, updatedAt time.Time) (UserMedbox, error) {
	if updatedAt.Before(createdAt) {
		return UserMedbox{}, nil
	}

	var medicationsMap = make(map[uuid.UUID]*medication.Medication)
	for _, med := range medications {
		medicationsMap[med.ID()] = &med
	}

	return UserMedbox{
		id:          id,
		userID:      userID,
		medications: medicationsMap,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		version:     0,
		events:      make([]Event, 0),
	}, nil
}

func (mb *UserMedbox) getMedicationByID(id uuid.UUID) *medication.Medication {
	return mb.medications[id]
}

func (mb *UserMedbox) AlreadyHas(m *medication.Medication) bool {
	dataMatrix := m.DataMatrix()
	for k, v := range mb.medications {
		if k == m.ID() {
			return false
		}
		if !dataMatrix.IsEmpty() && dataMatrix == v.DataMatrix() {
			return false
		}
	}
	return true
}
