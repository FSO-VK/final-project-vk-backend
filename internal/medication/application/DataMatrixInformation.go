// Package application is a package for application logic of the medication service.
package application

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
)

var (
	ErrCantSetCache = errors.New("error when setting cache")
)

// DataMatrixInformation is an interface for scanned info from data matrix.
type DataMatrixInformation interface {
	Execute(
		ctx context.Context,
		cmd *DataMatrixInformationCommand,
	) (*DataMatrixInformationResponse, error)
}

// DataMatrixInformationService is a service for adding a medication.
type DataMatrixInformationService struct {
	dataMatrixClient DataMatrixClient
	dataMatrixCache  DataMatrixCache
	validator        validator.Validator
}

// NewDataMatrixInformationService returns a new DataMatrixInformationService.
func NewDataMatrixInformationService(
	dataMatrixClient DataMatrixClient,
	dataMatrixCache DataMatrixCache,
	valid validator.Validator,
) *DataMatrixInformationService {
	return &DataMatrixInformationService{
		dataMatrixClient: dataMatrixClient,
		validator:        valid,
	}
}

// DataMatrixInformationCommand is a request to add a medication.
type DataMatrixInformationCommand struct {
	GTIN         string `validate:"required"`
	SerialNumber string `validate:"required"`
	CryptoData91 string `validate:"required"`
	CryptoData92 string `validate:"required"`
}

// DataMatrixInformationResponse is a response to add a medication.
type DataMatrixInformationResponse struct {
	ExpDate string
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
	dataMatrixInfo, err := s.dataMatrixCache.Get(req.GTIN + req.SerialNumber + req.CryptoData91 + req.CryptoData92)
	var errOut error = nil
	if err != nil {
		dataMatrixInfo, err = s.dataMatrixClient.GetInformationByDataMatrix(req.GTIN, req.SerialNumber, req.CryptoData91, req.CryptoData92)
		if err == nil {
			err = s.dataMatrixCache.Set(req.GTIN+req.SerialNumber+req.CryptoData91+req.CryptoData92, dataMatrixInfo)
			if err != nil {
				errOut = ErrCantSetCache
			}
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get medication: %w", err)
	}

	return &DataMatrixInformationResponse{
		ExpDate: dataMatrixInfo.ExpDate,
	}, errOut
}
