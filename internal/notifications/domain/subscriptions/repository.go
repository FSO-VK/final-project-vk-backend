package subscriptions

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// ErrNoSubscriptionsFound is an error when subscriptions is not found.
var ErrNoSubscriptionsFound = errors.New("subscriptions not found")

// Repository is a domain repository interface that defines
// data access contract for subscriptions aggregate.
type Repository interface {
	GetSubscriptionsByUserID(ctx context.Context, userID uuid.UUID) ([]*PushSubscription, error)
	GetSubscriptionByID(ctx context.Context, subscriptionID uuid.UUID) (*PushSubscription, error)
	SetSubscription(ctx context.Context, subscription *PushSubscription) error
}
