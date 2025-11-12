// Package application is a package for application logic of the medication service.
package application

import (
	"context"
	"errors"
	"fmt"
	"strings"

	client "github.com/FSO-VK/final-project-vk-backend/internal/medication/application/api_client"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/FSO-VK/final-project-vk-backend/pkg/validation"
)

var (
	// ErrCantSetCache is an error when setting cache.
	ErrCantSetCache = errors.New("fail to set cache")
	// ErrEmptyInput is an error when input is empty.
	ErrEmptyInput = errors.New("empty input")
)

// DataMatrixInformation is an interface for scanned info from data matrix.
type DataMatrixInformation interface {
	Execute(
		ctx context.Context,
		cmd *DataMatrixInformationCommand,
	) (*DataMatrixInformationResponse, error)
}

// DataMatrixInformationService is a service for get info from API.
type DataMatrixInformationService struct {
	dataMatrixClient client.DataMatrixClient
	dataMatrixCache  client.DataMatrixCacher
	validator        validator.Validator
}

// NewDataMatrixInformationService returns a new DataMatrixInformationService.
func NewDataMatrixInformationService(
	dataMatrixClient client.DataMatrixClient,
	dataMatrixCache client.DataMatrixCacher,
	valid validator.Validator,
) *DataMatrixInformationService {
	return &DataMatrixInformationService{
		dataMatrixClient: dataMatrixClient,
		dataMatrixCache:  dataMatrixCache,
		validator:        valid,
	}
}

// DataMatrixInformationCommand is a request to get info from API.
type DataMatrixInformationCommand struct {
	Data string `validate:"required"`
}

// DataMatrixInformationResponse is a response to get info from API.
type DataMatrixInformationResponse struct {
	// embedded struct
	CommandBase
}

// ParsedInformation is a struct for parsed data matrix.
type ParsedInformation struct {
	GTIN         string
	SerialNumber string
	CryptoData91 string
	CryptoData92 string
}

// Execute executes the DataMatrixInformation command.
func (s *DataMatrixInformationService) Execute(
	ctx context.Context,
	req *DataMatrixInformationCommand,
) (*DataMatrixInformationResponse, error) {
	valErr := s.validator.ValidateStruct(req)
	if valErr != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidationFail, valErr)
	}
	parsedData, err := ParseDataMatrix(req.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse data matrix: %w", err)
	}
	concatenatedInfo := parsedData.GTIN + parsedData.SerialNumber + parsedData.CryptoData91 + parsedData.CryptoData92
	dataMatrixInfo, err := s.dataMatrixCache.Get(
		ctx,
		concatenatedInfo,
	)
	var errOut error
	if err != nil {
		code := client.NewDataMatrixCodeInfo(
			parsedData.GTIN,
			parsedData.SerialNumber,
			parsedData.CryptoData91,
			parsedData.CryptoData92)
		dataMatrixInfo, err = s.dataMatrixClient.GetInformationByDataMatrix(code)
		if err == nil {
			err = s.dataMatrixCache.Set(
				ctx,
				concatenatedInfo,
				dataMatrixInfo,
			)
			if err != nil {
				errOut = ErrCantSetCache
			}
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get medication: %w", err)
	}

	return &DataMatrixInformationResponse{
		CommandBase: CommandBase{
			Name:                dataMatrixInfo.Name,
			InternationalName:   dataMatrixInfo.InternationalName,
			AmountValue:         dataMatrixInfo.AmountValue,
			AmountUnit:          dataMatrixInfo.AmountUnit,
			ReleaseForm:         dataMatrixInfo.ReleaseForm,
			Group:               dataMatrixInfo.Group,
			ManufacturerName:    dataMatrixInfo.ManufacturerName,
			ManufacturerCountry: dataMatrixInfo.ManufacturerCountry,
			ActiveSubstance:     MapAPIActiveSubstanceToLocal(dataMatrixInfo.ActiveSubstance),
			Expires:             dataMatrixInfo.Expires,
			Release:             dataMatrixInfo.Release,
			Commentary:          "",
		},
	}, errOut
}

// ParseDataMatrix creates validated data matrix string.
func ParseDataMatrix(data string) (*ParsedInformation, error) {
	if data == "" {
		return nil, ErrEmptyInput
	}

	var gtin, serial, crypto91, crypto92 string
	const fmtPattern = "01%14s21%13s91%4s92%44s"

	_, scanErr := fmt.Sscanf(data, fmtPattern, &gtin, &serial, &crypto91, &crypto92)

	if scanErr != nil {
		return nil, fmt.Errorf("failed to parse datamatrix string: %w", scanErr)
	}

	err := errors.Join(
		validation.Required(gtin),
		validation.Required(serial),
		validation.Required(crypto91),
		validation.Required(crypto92),
		validation.GTIN(gtin),
		validation.Serial(serial),
		validation.Crypto91(crypto91),
		validation.Crypto92(crypto92),
	)
	if err != nil {
		return nil, err
	}
	crypto92 = strings.TrimSuffix(crypto92, "=")
	return &ParsedInformation{
		GTIN:         gtin,
		SerialNumber: serial,
		CryptoData91: crypto91,
		CryptoData92: crypto92,
	}, nil
}

// MapAPIActiveSubstanceToLocal maps api active substance to local active substance.
func MapAPIActiveSubstanceToLocal(apiSubstances []client.ActiveSubstance) []ActiveSubstance {
	result := make([]ActiveSubstance, len(apiSubstances))
	for i, substance := range apiSubstances {
		result[i] = ActiveSubstance{
			Name:  substance.Name,
			Value: substance.Value,
			Unit:  substance.Unit,
		}
	}
	return result
}
