package http

import "github.com/FSO-VK/final-project-vk-backend/pkg/api"

const (
	// MsgFailedToGetVapidPublicKey is a message for failed to get vapid public key.
	MsgFailedToGetVapidPublicKey api.ErrorType = "Failed to get vapid public key"
	// MsgFailedToCreateSubscription is a message for failed to create subscription.
	MsgFailedToCreateSubscription api.ErrorType = "Failed to create subscription"
	// MsgFailedToDeleteSubscription is a message for failed to delete subscription.
	MsgFailedToDeleteSubscription api.ErrorType = "Failed to delete subscription"
	// MsgFailedToSendNotification is a message for failed to send notification.
	MsgFailedToSendNotification api.ErrorType = "Failed to send notification"
	// MsgMissingSlug is a message for missing slug.
	MsgMissingSlug api.ErrorType = "Missing slug"
)
