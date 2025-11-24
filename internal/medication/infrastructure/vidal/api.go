// Package vidal is an implementation of instruction application service.
package vidal

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/application/medreference"
	"github.com/FSO-VK/final-project-vk-backend/pkg/validation"
)

// Service is an implementation of medreference service.
type Service struct {
	storage Storage
	client  Client
}

// NewService creates a new Service.
func NewService(storage Storage, client Client) *Service {
	return &Service{storage: storage, client: client}
}

// GetProductInfo returns a product info by bar code.
func (s *Service) GetProductInfo(
	ctx context.Context,
	barCode string,
) (*medreference.Product, error) {
	err := validation.EAN13(barCode)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", medreference.ErrBadBarCode, err)
	}

	model, err := s.storage.GetProduct(ctx, barCode)
	if err != nil {
		if !errors.Is(err, ErrStorageNoProduct) {
			return nil, fmt.Errorf("get product from storage: %w", err)
		}
	} else {
		return s.productInfoToInstruction(model, barCode)
	}

	clientResponse, err := s.client.GetInstruction(ctx, barCode)
	if err != nil {
		if errors.Is(err, ErrClientNoProduct) {
			return nil, errors.Join(medreference.ErrNoProduct, err)
		}
		return nil, fmt.Errorf("external service: %w", err)
	}

	model = s.clientResponseToModel(clientResponse)

	err = s.storage.SaveProduct(ctx, model)
	if err != nil {
		return nil, fmt.Errorf("save product to s: %w", err)
	}

	return s.productInfoToInstruction(model, barCode)
}

func (s *Service) productInfoToInstruction(
	product *StorageModel,
	barCode string,
) (*medreference.Product, error) {
	nozologies := make([]medreference.Nozology, len(product.Document.Nozologies))
	for i, nozology := range product.Document.Nozologies {
		nozologies[i] = medreference.Nozology{
			Code: nozology.Code,
			Name: nozology.Name,
		}
	}

	clPhPointers := make([]medreference.ClPhPointer, len(product.Document.ClphPointers))
	for i, clPhPointer := range product.Document.ClphPointers {
		clPhPointers[i] = medreference.ClPhPointer{
			Code: clPhPointer.Code,
			Name: clPhPointer.Name,
		}
	}

	instruction := medreference.Instruction{
		Nozologies:             nozologies,
		ClPhPointers:           clPhPointers,
		PharmInfluence:         product.Document.PhInfluence,
		PharmKinetics:          product.Document.PhKinetics,
		Dosage:                 product.Document.Dosage,
		OverDosage:             product.Document.OverDosage,
		Interaction:            product.Document.Interaction,
		Lactation:              product.Document.Lactation,
		SideEffects:            product.Document.SideEffects,
		UsingIndication:        product.Document.Indication,
		UsingCounterIndication: product.Document.ContraIndication,
		SpecialInstruction:     product.Document.SpecialInstruction,
		RenalInsuf:             product.Document.RenalInsuf,
		HepatoInsuf:            product.Document.HepatoInsuf,
		ElderlyInsuf:           product.Document.ElderlyInsuf,
		ChildInsuf:             product.Document.ChildInsuf,
	}

	phGroups := extractPharmGroups(product.PhthGroups)

	clPhGroups := make([]string, len(product.ClPhGroups))
	for i, clPhGroup := range product.ClPhGroups {
		clPhGroups[i] = clPhGroup.Name
	}

	images := make([]string, len(product.Images))
	copy(images, product.Images)

	activeSubstances := make([]string, len(product.MoleculeNames))
	for i, substance := range product.MoleculeNames {
		activeSubstances[i] = substance.RusName
	}

	fullForm := ""
	if len(product.FullForm) > 0 {
		fullForm = product.FullForm[0]
	}

	manufacturer := extractManufacturer(product.Companies)

	p := &medreference.Product{
		BarCode:           barCode,
		RusName:           product.RusName,
		PharmGroups:       phGroups,
		ClinicPharmGroups: clPhGroups,
		ImagesLink:        images,
		ActiveSubstance:   activeSubstances,
		IsPrescription:    product.NonPrescriptionDrug,
		ReleaseForm:       fullForm,
		Manufacturer:      manufacturer,
		Instruction:       instruction,
	}
	return p, nil
}

func (s *Service) clientResponseToModel(clientResponse *ClientResponse) *StorageModel {
	barCodes := make([]string, len(clientResponse.ProductPackages))
	for i, pack := range clientResponse.ProductPackages {
		barCodes[i] = pack.BarCode
	}

	model := &StorageModel{
		Product:   clientResponse.Product,
		BarCodes:  barCodes,
		CreatedAt: time.Now(),
	}
	return model
}

func extractPharmGroups(phThGroups []PhthGroup) []string {
	phGroups := make([]string, 0, len(phThGroups))
	for _, phGroup := range phThGroups {
		// sometimes pharm groups contained in one string separated by semicolon.
		groups := strings.SplitSeq(phGroup.Code, ";")
		for group := range groups {
			phGroups = append(phGroups, strings.TrimSpace(group))
		}
	}
	return phGroups
}

func extractManufacturer(companies []CompanyInfo) medreference.Manufacturer {
	var manufacturer medreference.Manufacturer
	for _, company := range companies {
		if len(companies) == 1 {
			manufacturer = medreference.Manufacturer{
				Name:    company.Company.Name,
				Country: company.Company.Country.RusName,
			}
		} else if company.IsManufacturer {
			manufacturer = medreference.Manufacturer{
				Name:    company.Company.Name,
				Country: company.Company.Country.RusName,
			}
			break
		}
	}
	return manufacturer
}
