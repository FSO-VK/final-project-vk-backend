// Package apiclient is a package for api client interfaces.
package apiclient

import (
	"context"
	"errors"
)

// ErrNoMedicationFound is an error when a medication is not found.
var ErrNoMedicationFound = errors.New("medication not found")

// DataMatrixCache is an interface for data matrix cache.
type DataMatrixCache interface {
	Get(ctx context.Context, key string) (*MedicationInfo, error)
	Set(ctx context.Context, key string, data *MedicationInfo) error
}

// DataMatrixClient is an interface for data matrix client.
type DataMatrixClient interface {
	GetInformationByDataMatrix(data *DataMatrixCodeInfo) (*MedicationInfo, error)
}

// MedicationInfo contains info about medication.
type MedicationInfo struct {
	Name                string
	InternationalName   string
	AmountValue         float32
	AmountUnit          string
	ReleaseForm         string
	Group               string
	ManufacturerName    string
	ManufacturerCountry string
	ActiveSubstanceName string
	ActiveSubstanceDose float32
	ActiveSubstanceUnit string
	Expires             string
	Release             string
}

// DataMatrixCodeInfo contains scanned info from data matrix.
type DataMatrixCodeInfo struct {
	GTIN         string
	SerialNumber string
	CryptoData91 string
	CryptoData92 string
}

// NewDataMatrixCodeInfo creates a new DataMatrixCodeInfo.
func NewDataMatrixCodeInfo(
	gtin string,
	serialNumber string,
	cryptoData91 string,
	cryptoData92 string,
) *DataMatrixCodeInfo {
	return &DataMatrixCodeInfo{
		GTIN:         gtin,
		SerialNumber: serialNumber,
		CryptoData91: cryptoData91,
		CryptoData92: cryptoData92,
	}
}
