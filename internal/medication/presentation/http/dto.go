package http

// ProducerObject represents JSON object of producer of medication.
type ProducerObject struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

// ActiveSubstanceObject represents JSON object of active substance.
type ActiveSubstanceObject struct {
	Name  string  `json:"name,omitempty"`
	Value float32 `json:"value,omitzero"`
	Unit  string  `json:"unit,omitempty"`
}

// AmountObject is a structure of JSON object of amount of medication.
type AmountObject struct {
	Value float32 `json:"value"`
	Unit  string  `json:"unit"`
}

// BodyCommonObject is a common structure of JSON object.
type BodyCommonObject struct {
	Name              string                  `json:"name"`
	InternationalName string                  `json:"internationalName"`
	Amount            AmountObject            `json:"amount"`
	ReleaseForm       string                  `json:"releaseForm"`
	Group             []string                `json:"group"`
	Producer          ProducerObject          `json:"producer"`
	ActiveSubstance   []ActiveSubstanceObject `json:"activeSubstance"`
	Expiration        string                  `json:"expirationDate"`
	Release           string                  `json:"releaseDate,omitempty"`
	Commentary        string                  `json:"commentary"`
}

// BodyAPIObject is a api structure of JSON object.
type BodyAPIObject struct {
	Name              string         `json:"name"`
	InternationalName string         `json:"internationalName"`
	Amount            AmountObject   `json:"amount"`
	ReleaseForm       string         `json:"releaseForm"`
	Group             []string       `json:"group"`
	Producer          ProducerObject `json:"producer"`
	Expiration        string         `json:"expirationDate"`
	Release           string         `json:"releaseDate"`
}

// Nosology is an illness.
type Nosology struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// ClPhPointer is a clinical-pharmacological pointer.
type ClPhPointer struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// InstructionCommonObject is a common structure of JSON instruction.
type InstructionCommonObject struct {
	Nosologies             []Nosology    `json:"diseases"`
	ClPhPointers           []ClPhPointer `json:"clPhPointers"`
	PharmInfluence         string        `json:"pharmInfluence"`
	PharmKinetics          string        `json:"pharmKinetics"`
	Dosage                 string        `json:"dosage"`
	OverDosage             string        `json:"overDosage"`
	Interaction            string        `json:"interaction"`
	Lactation              string        `json:"lactation"`
	SideEffects            string        `json:"sideEffects"`
	UsingIndication        string        `json:"usingIndication"`
	UsingCounterIndication string        `json:"usingCounterIndication"`
	SpecialInstruction     string        `json:"specialInstruction"`
	RenalInsuf             string        `json:"renalInfluence"`
	HepatoInsuf            string        `json:"hepaticInfluence"`
	ElderlyInsuf           string        `json:"elderlyUsage"`
	ChildInsuf             string        `json:"childUsage"`
}
