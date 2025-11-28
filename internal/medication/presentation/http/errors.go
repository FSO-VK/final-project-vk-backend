package http

import "github.com/FSO-VK/final-project-vk-backend/pkg/api"

const (
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
	// MsgFailedToGetMedicationBox is a message for failed to get medication box.
	MsgFailedToGetMedicationBox api.ErrorType = "Failed to get medication box"
	// MsgNoMedication indicates that no medication found.
	MsgNoMedication api.ErrorType = "No medication found by such id"
	// MsgFailedToGetMedication is shown when service failed to process get by id request.
	MsgFailedToGetMedication api.ErrorType = "Failed to get medication"
	// MsgFailedToGetInstructions is a message for failed to get instructions.
	MsgFailedToGetInstructions api.ErrorType = "Failed to get instructions"
	// MsgFailedToGetInfoFromLLM is a message for failed to get info from LLM.
	MsgFailedToGetInfoFromLLM api.ErrorType = "Failed to get info from LLM"
)
