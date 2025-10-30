package auth

// AuthChecker is the interface for the CheckAuth method.
type AuthChecker interface {
	CheckAuth(req *Request) (*Response, error)
}

// Request is the request for the CheckAuth method.
type Request struct {
	SessionID string `json:"sessionId"`
}

// Response is the response for the CheckAuth method.
type Response struct {
	SessionID    string `json:"sessionId"`
	UserID       string `json:"useId"`
	IsAuthorized bool   `json:"isAuthorized"`
}

// ExpectedCheckAuthResponse is the expected response for the CheckAuth method.
type ExpectedCheckAuthResponse struct {
	UserID string `json:"userId"`
}
