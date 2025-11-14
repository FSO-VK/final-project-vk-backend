// Package instruction is an application service interface for instruction.
package medreference

import (
	"context"
	"errors"
)

var ErrBadBarCode = errors.New("invalid bar code")

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

type ClPhPointer struct {
	Code string
	Name string
}

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
