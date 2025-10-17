package http

import "github.com/FSO-VK/final-project-vk-backend/pkg/api"

const (
	// MsgFailedToReadBody is a message for failed to read body.
	MsgFailedToReadBody api.ErrorType = "Failed to read request body"
	// MsgFailedToUnmarshal is a message for failed to unmarshal body.
	MsgFailedToUnmarshal api.ErrorType = "Failed to unmarshal request body"
	// MsgFailedToAddMedicine is a message for failed to add medicine.
	MsgFailedToAddMedicine api.ErrorType = "Failed to add medicine"
	// MsgFailedToUpdateMedicine is a message for failed to update medicine.
	MsgFailedToUpdateMedicine api.ErrorType = "Failed to update medicine"
	// MsgFailedToDeleteMedicine is a message for failed to delete medicine.
	MsgFailedToDeleteMedicine api.ErrorType = "Failed to delete medicine"
	// MsgFailedToGetMedicineList is a message for failed to get medicine list.
	MsgFailedToGetMedicineList api.ErrorType = "Failed to get medicine list"
	// MsgFailToParseID is a message for failed to parse id.
	MsgFailToParseID api.ErrorType = "Failed to parse id"
)
