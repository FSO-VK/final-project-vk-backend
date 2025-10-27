// Package application is a package for application logic of the medication service.
package application

import (
	"context"
	"errors"
	"fmt"
	"regexp"
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
	Data string `validate:"required"`
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
	gtin, serialNumber, cryptoData91, cryptoData92, err := parseDataMatrixString(req.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse data matrix: %w", err)
	}
	cryptoData92 = strings.TrimSuffix(cryptoData92, "=")
	dataMatrixInfo, err := s.dataMatrixCache.Get(
		ctx,
		gtin+serialNumber+cryptoData91+cryptoData92,
	)
	var errOut error
	if err != nil {
		code := client.NewDataMatrixCodeInfo(
			gtin,
			serialNumber,
			cryptoData91,
			cryptoData92)
		dataMatrixInfo, err = s.dataMatrixClient.GetInformationByDataMatrix(code)
		if err == nil {
			err = s.dataMatrixCache.Set(
				ctx,
				gtin+serialNumber+cryptoData91+cryptoData92, dataMatrixInfo,
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

func parseDataMatrixString(data string) (gtin, serialNumber, crypto91, crypto92 string, err error) {
	re, err := regexp.Compile(`\(01\)([^\(]+)|\(21\)([^\(]+)|\(91\)([^\(]+)|\(92\)([^\(]+)`)
	if err != nil {
		
	}
	matches := re.FindAllStringSubmatch(data, -1)

	if len(matches) < 4 {
		return "", "", "", "", fmt.Errorf("invalid data matrix format")
	}

	for _, match := range matches {
		switch {
		case match[1] != "":
			gtin = match[1]
		case match[2] != "":
			serialNumber = match[2]
		case match[3] != "":
			crypto91 = match[3]
		case match[4] != "":
			crypto92 = match[4]
		}
	}

	if gtin == "" || serialNumber == "" || crypto91 == "" || crypto92 == "" {
		return "", "", "", "", fmt.Errorf("missing required data matrix fields")
	}

	return gtin, serialNumber, crypto91, crypto92, nil
}
