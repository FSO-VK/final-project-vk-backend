package command

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/user_medbox/medication"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/apperror"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/sirupsen/logrus"
)

const (
	ValidationErrorName = "ValidationError"
)

type commandService[T, E any] interface {
	Execute(ctx context.Context, cmd T) (E, error)
}

type loggingDecorator[T, E any] struct {
	svc    commandService[T, E]
	logger *logrus.Entry
}

func (ld loggingDecorator[T, E]) Execute(ctx context.Context, cmd T) (resp E, err error) {
	serviceType := generateActionName(cmd)

	logger := ld.logger.WithFields(logrus.Fields{
		"command":      serviceType,
		"command_body": fmt.Sprintf("%#v", cmd),
	})

	logger.Debug("Executing command")
	defer func() {
		if err == nil {
			logger.Info("Command executed successfully")
		} else {
			logger.WithError(err).Error("Failed to execute command")
		}
	}()

	return ld.svc.Execute(ctx, cmd)
}

type validationDecorator[T, E any] struct {
	svc       commandService[T, E]
	validator validator.Validator
}

func (vd validationDecorator[T, E]) Execute(ctx context.Context, cmd T) (resp E, err error) {
	err = vd.validator.ValidateStruct(cmd)
	if err != nil {
		return resp, apperror.User(ValidationErrorName, err)
	}
	return vd.svc.Execute(ctx, cmd)
}

func applyCommandDecorators[T, E any](
	svc commandService[T, E],
	log *logrus.Entry,
	val validator.Validator) commandService[T, E] {
	return loggingDecorator[T, E]{
		svc: validationDecorator[T, E]{
			svc:       svc,
			validator: val,
		},
		logger: log,
	}
}

func generateActionName(service any) string {
	return strings.Split(fmt.Sprintf("%T", service), ".")[1]
}

type MedicationInfo struct {
	Name                string            `validate:"required,max=100"`
	InternationalName   string            `validate:"omitempty,max=100"`
	AmountValue         float32           `validate:"required,gte=0"`
	AmountUnit          string            `validate:"required"`
	ReleaseForm         string            `validate:"required"`
	Groups              []string          `validate:"omitempty,max=10"`
	ManufacturerName    string            `validate:"omitempty,max=100"`
	ManufacturerCountry string            `validate:"omitempty,max=100"`
	ActiveSubstances    []ActiveSubstance `validate:"omitempty,max=10"`
	ExpirationDate      string            `validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	ReleaseDate         string            `validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Commentary          string            `validate:"omitempty,max=1000"`
}

// ActiveSubstance represents active substance.
type ActiveSubstance struct {
	Name  string  `validate:"required,max=100"`
	Value float32 `validate:"required,gte=0"`
	Unit  string  `validate:"required,max=100"`
}

func adaptDomainMedicationInfo(info medication.MedicationInfo) MedicationInfo {
	var releaseDate string
	if !info.ReleaseDate().IsZero() {
		releaseDate = info.ReleaseDate().Format(time.DateOnly)
	}

	groups := make([]string, 0, len(info.Groups()))
	for _, group := range info.Groups() {
		groups = append(groups, group.String())
	}

	activeSubstances := make([]ActiveSubstance, 0, len(info.ActiveSubstances()))
	for _, activeSubstance := range info.ActiveSubstances() {
		activeSubstances = append(activeSubstances, ActiveSubstance{
			Name:  activeSubstance.Name(),
			Value: activeSubstance.Dose().Value(),
			Unit:  activeSubstance.Dose().Unit().String(),
		})
	}

	return MedicationInfo{
		Name:                string(info.Name()),
		InternationalName:   string(info.InternationalName()),
		AmountValue:         info.Amount().Value(),
		AmountUnit:          info.Amount().Unit().String(),
		ReleaseForm:         string(info.ReleaseForm()),
		Groups:              groups,
		ManufacturerName:    string(info.Manufacturer().Name()),
		ManufacturerCountry: string(info.Manufacturer().Country()),
		ActiveSubstances:    activeSubstances,
		ExpirationDate:      info.ExpirationDate().Format(time.DateOnly),
		ReleaseDate:         releaseDate,
		Commentary:          string(info.Commentary()),
	}
}
