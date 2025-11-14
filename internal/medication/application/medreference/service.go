// Package medreference contains service for drug reference.
package medreference

import (
	"context"
	"errors"
)

// Common errors for application layer.
var (
	ErrBadBarCode = errors.New("invalid bar code")
	ErrNoProduct  = errors.New("can't get product with such bar code")
)

// Manufacturer represents a medication manufacturer.
type Manufacturer struct {
	Name    string
	Country string
}

// Nozology is an illness.
type Nozology struct {
	// Code is illness ICD-10 code
	Code string
	Name string
}

// ClPhPointer is a clinical-pharmacological pointer.
type ClPhPointer struct {
	Code string
	Name string
}

// Instruction is a medication instruction.
type Instruction struct {
	Nozologies             []Nozology
	ClPhPointers           []ClPhPointer
	PharmInfluence         string
	PharmKinetics          string
	Dosage                 string
	OverDosage             string
	Interaction            string
	Lactation              string
	SideEffects            string
	UsingIndication        string
	UsingCounterIndication string
	SpecialInstruction     string
	RenalInsuf             string
	HepatoInsuf            string
	ElderlyInsuf           string
	ChildInsuf             string
}

// Product is a data structure for instruction.
type Product struct {
	BarCode         string
	RusName         string
	PharmGroups     []string
	ImagesLink      []string
	ActiveSubstance []string
	IsPrescription  bool
	ReleaseForm     string
	Manufacturer    Manufacturer
	Instruction     Instruction
}

// MedicationReferenceProvider is an interface for instruction application service.
type MedicationReferenceProvider interface {
	GetProductInfo(ctx context.Context, barCode string) (*Product, error)
}
