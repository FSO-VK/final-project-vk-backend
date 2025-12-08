package http

import "github.com/FSO-VK/final-project-vk-backend/pkg/api"

const (
	// MsgFailedToGetPlan is a message for failed to get plan.
	MsgFailedToGetPlan api.ErrorType = "Failed to get plan"
	// MsgMissingSlug is a message for missing slug.
	MsgMissingSlug api.ErrorType = "Missing slug"
	// MsgFailedToAddPlan is a message for failed to add plan.
	MsgFailedToAddPlan api.ErrorType = "Failed to add plan"
)
