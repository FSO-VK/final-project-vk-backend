package http

import "github.com/FSO-VK/final-project-vk-backend/pkg/api"

const (
	// MsgFailedToReadBody is a message for failed to read body.
	MsgFailedToReadBody api.ErrorType = "Failed to read request body"
	// MsgFailedToUnmarshal is a message for failed to unmarshal body.
	MsgFailedToUnmarshal api.ErrorType = "Failed to unmarshal request body"
	// MsgFailedToAddMedication is a message for failed to add medication.
	MsgFailedToAddMedication api.ErrorType = "Failed to add medication"
	// MsgFailedToUpdateMedication is a message for failed to update medication.
	MsgFailedToUpdateMedication api.ErrorType = "Failed to update medication"
	// MsgFailedToDeleteMedication is a message for failed to delete medication.
	MsgFailedToDeleteMedication api.ErrorType = "Failed to delete medication"
	// MsgFailToParseID is a message for failed to parse id.
	MsgFailToParseID api.ErrorType = "Failed to parse id"
	// MsgFailedToGetIfoFromScan is a message for failed to get info from scan.
	MsgFailedToGetIfoFromScan api.ErrorType = "Failed to get info from scan"
	// MsgFailedToGetMedicationBox is a message for failed to get info from scan.
	MsgFailedToGetMedicationBox api.ErrorType = "Failed to get medication box"
)
