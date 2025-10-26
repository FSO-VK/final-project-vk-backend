package http

type ProducerObject struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

type ActiveSubstanceObject struct {
	Name  string  `json:"name"`
	Value float32 `json:"value"`
	Unit  string  `json:"unit"`
}

type AmountObject struct {
	Value float32 `json:"value"`
	Unit  string  `json:"unit"`
}

type BodyCommonObject struct {
	Name              string                `json:"name"`
	InternationalName string                `json:"internationalName"`
	Amount            AmountObject          `json:"amount"`
	ReleaseForm       string                `json:"releaseForm"`
	Group             string                `json:"group"`
	Producer          ProducerObject        `json:"producer"`
	ActiveSubstance   ActiveSubstanceObject `json:"activeSubstance"`
	Expiration        string                `json:"expirationDate"`
	Release           string                `json:"releaseDate"`
	Commentary        string                `json:"commentary"`
}
