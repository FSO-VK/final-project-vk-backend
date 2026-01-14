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

type UpdateMedicationCommand struct {
	MedicationInfo

	UserID string `validate:"required,uuid"`
	ID     string `validate:"required,uuid"`
}

type UpdateMedicationResponse struct {
	MedicationInfo

	ID string
}

type UpdateMedicationService commandService[*UpdateMedicationCommand, *UpdateMedicationResponse]

type updateMedicationService struct {
	repo usermedbox.Repository
}

func NewUpdateMedicationService(
	repo usermedbox.Repository,
	log *logrus.Entry,
	val validator.Validator,
) UpdateMedicationService {
	if repo == nil {
		panic(fmt.Sprintf("%T is nil", repo))
	}

	return applyCommandDecorators(
		&updateMedicationService{
			repo: repo,
		},
		log,
		val,
	)
}

func (ums *updateMedicationService) Execute(ctx context.Context, cmd *UpdateMedicationCommand) (*UpdateMedicationResponse, error) {
	// these errors is impossible because uuid is checked by middleware
	userID, _ := uuid.Parse(cmd.UserID)
	medicationID, _ := uuid.Parse(cmd.ID)
	expirationDate, _ := time.Parse(time.RFC3339, cmd.ExpirationDate)
	var releaseDate time.Time
	if cmd.ReleaseDate != "" {
		releaseDate, _ = time.Parse(time.RFC3339, cmd.ReleaseDate)
	}

	medbox, err := ums.repo.GetOrCreate(ctx, userID)
	if err != nil {
		return nil, apperror.Internal(RepositoryErrorName, err)
	}

	changedInfoDraft := medication.MedicationInfoDraft{
		Name:              cmd.Name,
		InternationalName: cmd.InternationalName,
		Groups:            cmd.Groups,
		Manufacturer: medication.ManufacturerDraft{
			Name:    cmd.ManufacturerName,
			Country: cmd.ManufacturerCountry,
		},
		ReleaseForm: cmd.ReleaseDate,
		Amount: medication.AmountDraft{
			Value: cmd.AmountValue,
			Unit:  cmd.AmountUnit,
		},
		Commentary:       cmd.Commentary,
		ActiveSubstances: adaptActiveSubstancesToDomain(cmd.ActiveSubstances),
		ReleaseDate:      releaseDate,
		ExpirationDate:   expirationDate,
	}

	changedInfo, err := medication.NewMedicationInfo(changedInfoDraft)
	if err != nil {
		return nil, apperror.User(DomainErrorName, err)
	}

	err = medbox.UpdateMedicationInfo(medicationID, changedInfo)
	if err != nil {
		if errors.Is(err, usermedbox.ErrNoMedication) {
			return nil, apperror.User(MedicationDoesNotExistErrorName, err)
		}
		return nil, apperror.Internal(DomainErrorName, err)
	}

	err = ums.repo.Save(ctx, medbox)
	if err != nil {
		return nil, apperror.Internal(RepositoryErrorName, err)
	}

	changedMedication, err := medbox.GetMedication(medicationID)
	// situation when medication does not exist is impossible
	// because it checked by UpdateMedicationInfo
	if err != nil {
		return nil, apperror.Internal(DomainErrorName, err)
	}

	return &UpdateMedicationResponse{
		MedicationInfo: adaptDomainMedicationInfo(changedMedication.Info()),
		ID:             cmd.ID,
	}, nil
}
