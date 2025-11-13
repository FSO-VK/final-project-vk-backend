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
