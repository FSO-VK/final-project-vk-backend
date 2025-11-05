package notifications

import (
	"errors"
)

// ErrNoNotificationsFound is an error when notifications is not found.
var ErrNoNotificationsFound = errors.New("notifications not found")

// Repository is a domain repository interface that defines
// data access contract for notifications box aggregate.
type Repository interface {
}
