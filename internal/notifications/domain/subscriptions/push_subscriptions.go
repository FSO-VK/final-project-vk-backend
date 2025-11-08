// Package subscriptions is a domain layer for subscription.
package subscriptions

import (
	"github.com/google/uuid"
)

// PushSubscription is info about subscription.
type PushSubscription struct {
	id        uuid.UUID
	userID    uuid.UUID
	sendInfo  SendInfo
	userAgent string
	isActive  bool
}

// SendInfo is unique info for sending subscriptions.
type SendInfo struct {
	Endpoint string
	Keys     Keys
}

// Keys is unique keys for encryption.
type Keys struct {
	P256dh string
	Auth   string
}

// NewSubscription creates a new subscription.
func NewSubscription(
	userID uuid.UUID,
	endpoint string,
	p256dh string,
	auth string,
	userAgent string,
) *PushSubscription {
	return &PushSubscription{
		id:     uuid.New(),
		userID: userID,
		sendInfo: SendInfo{
			Endpoint: endpoint,
			Keys: Keys{
				P256dh: p256dh,
				Auth:   auth,
			},
		},
		userAgent: userAgent,
		isActive:  true,
	}
}

// GetID returns the unique identifier of the subscription.
func (s *PushSubscription) GetID() uuid.UUID {
	return s.id
}

// SetID sets the unique identifier of the subscription.
func (s *PushSubscription) SetID(id uuid.UUID) {
	s.id = id
}

// GetUserID returns the user identifier.
func (s *PushSubscription) GetUserID() uuid.UUID {
	return s.userID
}

// SetUserID sets the user identifier.
func (s *PushSubscription) SetUserID(userID uuid.UUID) {
	s.userID = userID
}

// GetSendInfo returns the send info.
func (s *PushSubscription) GetSendInfo() SendInfo {
	return s.sendInfo
}

// SetSendInfo sets the send info.
func (s *PushSubscription) SetSendInfo(sendInfo SendInfo) {
	s.sendInfo = sendInfo
}

// GetUserAgent returns the user agent.
func (s *PushSubscription) GetUserAgent() string {
	return s.userAgent
}

// SetUserAgent sets the user agent.
func (s *PushSubscription) SetUserAgent(userAgent string) {
	s.userAgent = userAgent
}

// GetIsActive returns the active status.
func (s *PushSubscription) GetIsActive() bool {
	return s.isActive
}

// SetIsActive sets the active status.
func (s *PushSubscription) SetIsActive(isActive bool) {
	s.isActive = isActive
}

// GetEndpoint returns the endpoint from send info.
func (s *PushSubscription) GetEndpoint() string {
	return s.sendInfo.Endpoint
}

// SetEndpoint sets the endpoint in send info.
func (s *PushSubscription) SetEndpoint(endpoint string) {
	s.sendInfo.Endpoint = endpoint
}

// GetKeys returns the keys from send info.
func (s *PushSubscription) GetKeys() Keys {
	return s.sendInfo.Keys
}

// SetKeys sets the keys in send info.
func (s *PushSubscription) SetKeys(keys Keys) {
	s.sendInfo.Keys = keys
}

// GetP256dh returns the P256dh key.
func (s *PushSubscription) GetP256dh() string {
	return s.sendInfo.Keys.P256dh
}

// SetP256dh sets the P256dh key.
func (s *PushSubscription) SetP256dh(p256dh string) {
	s.sendInfo.Keys.P256dh = p256dh
}

// GetAuth returns the auth key.
func (s *PushSubscription) GetAuth() string {
	return s.sendInfo.Keys.Auth
}

// SetAuth sets the auth key.
func (s *PushSubscription) SetAuth(auth string) {
	s.sendInfo.Keys.Auth = auth
}
