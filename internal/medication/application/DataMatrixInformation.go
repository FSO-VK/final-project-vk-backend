// Package application is a package for application logic of the medication service.
package application

import (
	"context"
	"errors"
	"fmt"
	"strings"

	client "github.com/FSO-VK/final-project-vk-backend/internal/medication/application/api_client"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
)

// ErrCantSetCache is an error when setting cache.
var ErrCantSetCache = errors.New("error when setting cache")

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
	dataMatrixCache  client.DataMatrixCache
	validator        validator.Validator
}

// NewDataMatrixInformationService returns a new DataMatrixInformationService.
func NewDataMatrixInformationService(
	dataMatrixClient client.DataMatrixClient,
	dataMatrixCache client.DataMatrixCache,
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
	GTIN         string `validate:"required"`
	SerialNumber string `validate:"required"`
	CryptoData91 string `validate:"required"`
	CryptoData92 string `validate:"required"`
}

// DataMatrixInformationResponse is a response to get info from API.
type DataMatrixInformationResponse struct {
	// embedded struct
	CommandBase
}

// Execute executes the DataMatrixInformation command.
func (s *DataMatrixInformationService) Execute(
	ctx context.Context,
	req *DataMatrixInformationCommand,
) (*DataMatrixInformationResponse, error) {
	valErr := s.validator.ValidateStruct(req)
	if valErr != nil {
		return nil, fmt.Errorf("failed to validate request: %w", valErr)
	}

	if req.CryptoData92 != "" {
		req.CryptoData92 = strings.TrimSuffix(req.CryptoData92, "=")
	}
	dataMatrixInfo, err := s.dataMatrixCache.Get(
		ctx,
		req.GTIN+req.SerialNumber+req.CryptoData91+req.CryptoData92,
	)
	var errOut error
	fmt.Println("1111111111111111111111111111111111111")
	if err != nil {
		fmt.Println("22222222222222222222222")
		code := client.NewDataMatrixCodeInfo(
			req.GTIN,
			req.SerialNumber,
			req.CryptoData91,
			req.CryptoData92)
		dataMatrixInfo, err = s.dataMatrixClient.GetInformationByDataMatrix(code)
		if err == nil {
			fmt.Println("333333333333333333333333333")
			err = s.dataMatrixCache.Set(
				ctx,
				req.GTIN+req.SerialNumber+req.CryptoData91+req.CryptoData92, dataMatrixInfo,
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
			ActiveSubstanceName: dataMatrixInfo.ActiveSubstanceName,
			ActiveSubstanceDose: dataMatrixInfo.ActiveSubstanceDose,
			ActiveSubstanceUnit: dataMatrixInfo.ActiveSubstanceUnit,
			Expires:             dataMatrixInfo.Expires,
			Release:             dataMatrixInfo.Release,
			Commentary:          "",
		},
	}, errOut
}
