// Package application implements dependency injection container for notifications service use cases.
package application

// NotificationsApplication is a dependency injection container that aggregates all use cases
// for the notification domain to be injected from main.go.
type NotificationsApplication struct {
	GetVapidPublicKey  GetVapidPublicKey
	CreateSubscription CreateSubscription
	DeleteSubscription DeleteSubscription
	SendNotification   SendNotification
}
