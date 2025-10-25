package client

import "net/http"

type AuthChecker interface {
	CheckAuth(req *Request) (*Response, error)
}

type Request struct {
	SessionID string `json:"sessionId"`
}

type Response struct {
	SessionID    string `json:"sessionId"`
	UserID       string `json:"useId"`
	IsAuthorized bool   `json:"isAuthorized"`
}

type ExpectedCheckAuthResponse struct {
	UserID string `json:"userId"`
}

// Abstraction above http.Client — any client with Do().
type HttpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}
