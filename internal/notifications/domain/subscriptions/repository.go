package subscriptions

import (
	"errors"
)

// ErrNoSubscriptionsFound is an error when subscriptions is not found.
var ErrNoSubscriptionsFound = errors.New("subscriptions not found")

// Repository is a domain repository interface that defines
// data access contract for subscriptions aggregate.
type Repository interface {
}
