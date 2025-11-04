// Package application implements dependency injection container for medication service use cases.
package application

// MedicationsApplication is a dependency injection container that aggregates all use cases
// for the medication domain to be injected from main.go.
type MedicationsApplication struct {
	GetVapidPublicKey        GetVapidPublicKey
	CreateSubscription       CreateSubscription
	DeleteSubscription       DeleteSubscription
	InteractWithNotification InteractWithNotification
}
