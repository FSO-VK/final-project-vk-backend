package http

// SubscriptionObject is info about subscription .
type SubscriptionObject struct {
	SendInfo  SendInfo `json:"subscription"`
	UserAgent string   `json:"ua"`
}

// SendInfo is unique info for sending notifications.
type SendInfo struct {
	Endpoint string `json:"endpoint"`
	Keys     Keys   `json:"keys"`
}

// Keys is unique keys for encryption.
type Keys struct {
	P256dh string `json:"p256dh"`
	Auth   string `json:"auth"`
}

// PushSubscription is info about subscription.
type PushSubscriptionInfo struct {
	ID        string   `json:"id"`
	UserID    string   `json:"userId"`
	SendInfo  SendInfo `json:"subscription"`
	UserAgent string   `json:"ua"`
	IsActive  bool     `json:"isActive"`
}

// PushNotificationObject is info about notification.
type PushNotificationObject struct {
	SubscriptionID string `json:"subscriptionId"`
	Title          string `json:"title"`
	Body           string `json:"body"`
	SendAt         string `json:"sendAt"`
}
