// Package vidal is an implementation of instruction application service.
package vidal

import (
	"context"
	"errors"
	"fmt"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/application/medreference"
	"github.com/FSO-VK/final-project-vk-backend/pkg/validation"
)

type Service struct {
	storage Storage
	client  Client
}

func NewService(storage Storage, client Client) *Service {
	return &Service{storage: storage, client: client}
}

func (s *Service) GetProductInfo(
	ctx context.Context,
	barCode string,
) (*medreference.Product, error) {
	err := validation.EAN13(barCode)
	if err != nil {
		return nil, medreference.ErrBadBarCode
	}

	model, err := s.storage.GetProduct(ctx, barCode)
	if err != nil {
		if !errors.Is(err, ErrNoProduct) {
			return nil, fmt.Errorf("get product from storage: %w", err)
		}
	} else {
		return s.productInfoToInstruction(model, barCode)
	}

	clientResponse, err := s.client.GetInstruction(ctx, barCode)
	if err != nil {
		return nil, fmt.Errorf("external service: %w", err)
	}

	model, err = s.clientResponseToModel(clientResponse)
	if err != nil {
		return nil, fmt.Errorf("convert client response to model: %w", err)
	}

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
	nozologies := make([]medreference.Nozology, len(product.Product.Document.Nozologies))
	for i, nozology := range product.Product.Document.Nozologies {
		nozologies[i] = medreference.Nozology{
			Code: nozology.Code,
			Name: nozology.Name,
		}
	}

	clPhPointers := make([]medreference.ClPhPointer, len(product.Product.Document.ClphPointers))
	for i, clPhPointer := range product.Product.Document.ClphPointers {
		clPhPointers[i] = medreference.ClPhPointer{
			Code: clPhPointer.Code,
			Name: clPhPointer.Name,
		}
	}

	instruction := medreference.Instruction{
		Nozologies:             nozologies,
		ClPhPointers:           clPhPointers,
		PharmInfluence:         product.Product.Document.PhInfluence,
		PharmKinetics:          product.Product.Document.PhKinetics,
		Dosage:                 product.Product.Document.Dosage,
		OverDosage:             product.Product.Document.OverDosage,
		Interaction:            product.Product.Document.Interaction,
		Lactation:              product.Product.Document.Lactation,
		SideEffects:            product.Product.Document.SideEffects,
		UsingIndication:        product.Product.Document.Indication,
		UsingCounterIndication: product.Product.Document.ContraIndication,
		SpecialInstruction:     product.Product.Document.SpecialInstruction,
		RenalInsuf:             product.Product.Document.RenalInsuf,
		HepatoInsuf:            product.Product.Document.HepatoInsuf,
		ElderlyInsuf:           product.Product.Document.ElderlyInsuf,
		ChildInsuf:             product.Product.Document.ChildInsuf,
	}

	phGroups := make([]string, len(product.PhthGroups))
	for i, phGroup := range product.PhthGroups {
		phGroups[i] = phGroup.Code
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

	var manufacturer medreference.Manufacturer
	for _, company := range product.Product.Companies {
		if company.IsManufacturer {
			manufacturer = medreference.Manufacturer{
				Name:    company.Company.Name,
				Country: company.Company.Country.RusName,
			}
			break
		}
	}

	p := &medreference.Product{
		BarCode:         barCode,
		RusName:         product.Product.RusName,
		PharmGroups:     phGroups,
		ImagesLink:      images,
		ActiveSubstance: activeSubstances,
		IsPrescription:  product.NonPrescriptionDrug,
		ReleaseForm:     fullForm,
		Manufacturer:    manufacturer,
		Instruction:     instruction,
	}
	return p, nil
}

func (s *Service) clientResponseToModel(clientResponse *ClientResponse) (*StorageModel, error) {
	barCodes := make([]string, len(clientResponse.ProductPackages))
	for i, pack := range clientResponse.ProductPackages {
		barCodes[i] = pack.BarCode
	}

	model := &StorageModel{
		Product:  clientResponse.Product,
		BarCodes: barCodes,
	}
	return model, nil
}
