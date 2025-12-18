package http

import "github.com/FSO-VK/final-project-vk-backend/pkg/api"

const (
	MsgWrongCredentials api.ErrorType = "Wrong credentials"
	MsgLogoutFailed     api.ErrorType = "Logout failed"
	MsgUnauthorized     api.ErrorType = "Session does not exist"
	MsgSetCookieFail    api.ErrorType = "Unable to set cookie"
	MsgInvalidEmail     api.ErrorType = "Invalid email"
	MsgInvalidPassword  api.ErrorType = "Invalid password"
	MsgUserAlreadyExist api.ErrorType = "User with this email already exist"
)
