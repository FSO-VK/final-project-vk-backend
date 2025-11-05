package http

// SubscriptionObject is info about subscription .
type SubscriptionObject struct {
	ID        string   `json:"id"`
	SendInfo  SendInfo `json:"subscription"`
	UserAgent string   `json:"ua"`
	Browser   string   `json:"browser"`
	OS        string   `json:"os"`
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
