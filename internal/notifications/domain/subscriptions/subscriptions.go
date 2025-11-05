// Package subscriptions is a domain layer for subscription.
package subscriptions

import (
	"github.com/google/uuid"
)

// Subscription is info about subscription.
type Subscription struct {
	id        uuid.UUID
	userID    uuid.UUID
	sendInfo  SendInfo
	userAgent string
	browser   string
	os        string
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
func NewSubscription(userID uuid.UUID) *Subscription {
	return &Subscription{
		id:       uuid.New(),
		userID:   userID,
		sendInfo: SendInfo{},
		isActive: true,
	}
}

// GetID returns the unique identifier of the subscription.
func (s *Subscription) GetID() uuid.UUID {
	return s.id
}

// SetID sets the unique identifier of the subscription.
func (s *Subscription) SetID(id uuid.UUID) {
	s.id = id
}

// GetUserID returns the user identifier.
func (s *Subscription) GetUserID() uuid.UUID {
	return s.userID
}

// SetUserID sets the user identifier.
func (s *Subscription) SetUserID(userID uuid.UUID) {
	s.userID = userID
}

// GetSendInfo returns the send info.
func (s *Subscription) GetSendInfo() SendInfo {
	return s.sendInfo
}

// SetSendInfo sets the send info.
func (s *Subscription) SetSendInfo(sendInfo SendInfo) {
	s.sendInfo = sendInfo
}

// GetUserAgent returns the user agent.
func (s *Subscription) GetUserAgent() string {
	return s.userAgent
}

// SetUserAgent sets the user agent.
func (s *Subscription) SetUserAgent(userAgent string) {
	s.userAgent = userAgent
}

// GetBrowser returns the browser.
func (s *Subscription) GetBrowser() string {
	return s.browser
}

// SetBrowser sets the browser.
func (s *Subscription) SetBrowser(browser string) {
	s.browser = browser
}

// GetOS returns the operating system.
func (s *Subscription) GetOS() string {
	return s.os
}

// SetOS sets the operating system.
func (s *Subscription) SetOS(os string) {
	s.os = os
}

// GetIsActive returns the active status.
func (s *Subscription) GetIsActive() bool {
	return s.isActive
}

// SetIsActive sets the active status.
func (s *Subscription) SetIsActive(isActive bool) {
	s.isActive = isActive
}

// GetEndpoint returns the endpoint from send info.
func (s *Subscription) GetEndpoint() string {
	return s.sendInfo.Endpoint
}

// SetEndpoint sets the endpoint in send info.
func (s *Subscription) SetEndpoint(endpoint string) {
	s.sendInfo.Endpoint = endpoint
}

// GetKeys returns the keys from send info.
func (s *Subscription) GetKeys() Keys {
	return s.sendInfo.Keys
}

// SetKeys sets the keys in send info.
func (s *Subscription) SetKeys(keys Keys) {
	s.sendInfo.Keys = keys
}

// GetP256dh returns the P256dh key.
func (s *Subscription) GetP256dh() string {
	return s.sendInfo.Keys.P256dh
}

// SetP256dh sets the P256dh key.
func (s *Subscription) SetP256dh(p256dh string) {
	s.sendInfo.Keys.P256dh = p256dh
}

// GetAuth returns the auth key.
func (s *Subscription) GetAuth() string {
	return s.sendInfo.Keys.Auth
}

// SetAuth sets the auth key.
func (s *Subscription) SetAuth(auth string) {
	s.sendInfo.Keys.Auth = auth
}
