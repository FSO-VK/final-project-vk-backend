package api

// ErrorType represents a specific category or type of error within the API.
type ErrorType string

const (
	MsgNoBody           ErrorType = "No request body"
	MsgBadBody          ErrorType = "Bad request body"
	MsgServerError      ErrorType = "Server error"
	MsgNotFound         ErrorType = "No such endpoint"
	MsgMethodNotAllowed ErrorType = "Method not allowed"
)
