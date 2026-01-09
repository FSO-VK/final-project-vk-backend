package medication

import (
	"errors"
	"fmt"
	"strings"

	"github.com/FSO-VK/final-project-vk-backend/pkg/validation"
)

type DataMatrix struct {
	// gtin represents a class of medication
	gtin string
	// unique identifier of a medication
	serialNumber string

	cryptoData91 string
	cryptoData92 string
}

// BarCode extracts medication barcode from DataMatrix.
func (dm *DataMatrix) BarCode() string {
	return dm.gtin[1:]
}

// String is a Stringer implementation for DataMatrix.
func (dm *DataMatrix) String() string {
	var b strings.Builder
	b.WriteString("01")
	b.WriteString(dm.gtin)
	b.WriteString("21")
	b.WriteString(dm.serialNumber)
	b.WriteString("%1D91")
	b.WriteString(dm.cryptoData91)
	b.WriteString("%1D92")
	b.WriteString(dm.cryptoData92)
	return b.String()
}

var emptyDataMatrix = DataMatrix{}

// IsEmpty checks is DataMatrix contains no data.
func (dm *DataMatrix) IsEmpty() bool {
	return *dm == emptyDataMatrix
}

// NewEmptyDataMatrix makes an empty DataMatrix.
func NewEmptyDataMatrix() DataMatrix {
	return emptyDataMatrix
}

// ParseDataMatrix creates DataMatrix from string.
func ParseDataMatrix(data string) (DataMatrix, error) {
	var gtin, serial, crypto91, crypto92 string
	const fmtPattern = "01%14s21%13s91%4s92%44s"

	_, scanErr := fmt.Sscanf(data, fmtPattern, &gtin, &serial, &crypto91, &crypto92)

	if scanErr != nil {
		return DataMatrix{}, fmt.Errorf("failed to parse datamatrix string: %w", scanErr)
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
		return DataMatrix{}, err
	}
	crypto92 = strings.TrimSuffix(crypto92, "=")
	return DataMatrix{
		gtin:         gtin,
		serialNumber: serial,
		cryptoData91: crypto91,
		cryptoData92: crypto92,
	}, nil
}

