// Package notificationsclient is a package for interface for notifications client.
package notificationsclient

// PushNotificationClient is an interface for data matrix client.
type PushNotificationClient interface {
	PushNotification(data *PushNotificationInfo) (*NotificationResponse, error)
}

// PushNotificationInfo contains required info for push notification.
type PushNotificationInfo struct {
	// TODO
}

// NotificationResponse contains info about pushed notification.
type NotificationResponse struct {
	// TODO
}

// NewPushNotificationInfo creates a new PushNotificationInfo.
func NewPushNotificationInfo(
// TODO
) *PushNotificationInfo {
	return &PushNotificationInfo{
		// TODO
	}
}
