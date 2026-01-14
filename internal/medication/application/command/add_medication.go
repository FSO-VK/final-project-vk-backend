package command

import (
	"context"
	"errors"
	"fmt"
	"time"

	usermedbox "github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/user_medbox"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/user_medbox/medication"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/apperror"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	IDGenerationErrorName    = "IDGenerationError"
	MedicationExistErrorName = "MedicationExistError"
)

type AddMedicationCommand struct {
	MedicationInfo

	DataMatrix string `validate:"omitempty,max=1000"`
	UserID     string `validate:"required,uuid"`
}

type AddMedicationResponse struct {
	MedicationInfo

	ID string
}

type AddMedicationService commandService[*AddMedicationCommand, *AddMedicationResponse]

type addMedicationService struct {
	repo usermedbox.Repository
}

func NewAddMedicationService(repo usermedbox.Repository, log *logrus.Entry, val validator.Validator) AddMedicationService {
	if repo == nil {
		panic(fmt.Sprintf("%T is nil", repo))
	}

	return applyCommandDecorators(
		&addMedicationService{
			repo: repo,
		},
		log,
		val,
	)
}

func (ams *addMedicationService) Execute(ctx context.Context, cmd *AddMedicationCommand) (*AddMedicationResponse, error) {
	// these errors is impossible because uuid is checked by middleware
	userID, _ := uuid.Parse(cmd.UserID)
	expirationDate, _ := time.Parse(time.RFC3339, cmd.ExpirationDate)
	var releaseDate time.Time
	if cmd.ReleaseDate != "" {
		releaseDate, _ = time.Parse(time.RFC3339, cmd.ReleaseDate)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, apperror.Internal(IDGenerationErrorName, err)
	}

	draft := medication.MedicationDraft{
		ID: id,
		Info: medication.MedicationInfoDraft{
			Name:              cmd.Name,
			InternationalName: cmd.InternationalName,
			Groups:            cmd.Groups,
			Manufacturer: medication.ManufacturerDraft{
				Name:    cmd.ManufacturerName,
				Country: cmd.ManufacturerCountry,
			},
			ReleaseForm: cmd.ReleaseForm,
			Amount: medication.AmountDraft{
				Value: cmd.AmountValue,
				Unit:  cmd.AmountUnit,
			},
			Commentary:     cmd.Commentary,
			ActiveSubstances: adaptActiveSubstancesToDomain(cmd.ActiveSubstances),
			ReleaseDate:    releaseDate,
			ExpirationDate: expirationDate,
		},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		DataMatrix: cmd.DataMatrix,
	}
	medbox, err := ams.repo.GetOrCreate(ctx, userID)
	if err != nil {
		return nil, apperror.Internal(RepositoryErrorName, err)
	}

	newMedication, err := medication.NewMedication(draft)
	if err != nil {
		return nil, apperror.User(ValidationErrorName, err)
	}

	err = medbox.AddMedication(newMedication)
	if err != nil {
		if errors.Is(err, usermedbox.ErrMedicationExist) {
			return nil, apperror.User(MedicationExistErrorName, err)
		}
		return nil, apperror.Internal(DomainErrorName, err)
	}

	err = ams.repo.Save(ctx, medbox)
	if err != nil {
		return nil, apperror.Internal(RepositoryErrorName, err)
	}

	return &AddMedicationResponse{
		MedicationInfo: adaptDomainMedicationInfo(newMedication.Info()),
		ID:             id.String(),
	}, nil
}
