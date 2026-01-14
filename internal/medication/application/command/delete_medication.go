package command

import (
	"context"
	"errors"
	"fmt"

	usermedbox "github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/user_medbox"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/apperror"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	MedicationDoesNotExistErrorName = "MedicationDoesNotExistError"
)

type DeleteMedicationCommand struct {
	UserID string `validate:"required,uuid"`
	ID     string `validate:"required,uuid"`
}

type DeleteMedicationResponse struct{}

type DeleteMedicationService commandService[*DeleteMedicationCommand, *DeleteMedicationResponse]

type deleteMedicationService struct {
	repo usermedbox.Repository
}

func NewDeleteMedicationService(
	repo usermedbox.Repository,
	log *logrus.Entry,
	val validator.Validator,
) DeleteMedicationService {
	if repo == nil {
		panic(fmt.Sprintf("%T is nil", repo))
	}

	return applyCommandDecorators(
		&deleteMedicationService{
			repo: repo,
		},
		log,
		val,
	)
}

func (dms *deleteMedicationService) Execute(ctx context.Context, cmd *DeleteMedicationCommand) (*DeleteMedicationResponse, error) {
	// these errors is impossible because uuid is checked by middleware
	userID, _ := uuid.Parse(cmd.UserID)
	medicationID, _ := uuid.Parse(cmd.ID)

	medbox, err := dms.repo.GetOrCreate(ctx, userID)
	if err != nil {
		return nil, apperror.Internal(RepositoryErrorName, err)
	}

	err = medbox.DeleteMedication(medicationID)
	if err != nil {
		if errors.Is(err, usermedbox.ErrNoMedication) {
			return nil, apperror.User(MedicationDoesNotExistErrorName, err)
		}
		return nil, apperror.Internal(DomainErrorName, err)
	}

	err = dms.repo.Save(ctx, medbox)
	if err != nil {
		return nil, apperror.Internal(RepositoryErrorName, err)
	}

	return &DeleteMedicationResponse{}, nil
}
