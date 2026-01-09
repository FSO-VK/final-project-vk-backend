package usermedbox

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var ErrFailedToGenerateEventID = errors.New("failed to generate uuid for event")

type Event interface {
	ID() uuid.UUID
	Timestamp() time.Time
	Name() string
}

type BasicEvent struct {
	id        uuid.UUID
	timestamp time.Time
}

func (be BasicEvent) ID() uuid.UUID {
	return be.id
}

func (be BasicEvent) Timestamp() time.Time {
	return be.timestamp
}

func NewBasicEvent() (BasicEvent, error) {
	eventID, err := uuid.NewV7()
	if err != nil {
		return BasicEvent{}, ErrFailedToGenerateEventID
	}
	return BasicEvent{
		id:        eventID,
		timestamp: time.Now(),
	}, nil
}
