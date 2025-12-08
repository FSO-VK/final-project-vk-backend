package http

// PlanObject is info about plan.
type PlanObject struct {
	MedicationID   string       `json:"medicationId"`
	UserID         string       `json:"userId"`
	Amount         AmountObject `json:"amount"`
	Condition      string       `json:"condition"`
	StartDate      string       `json:"startDate"`
	EndDate        string       `json:"endDate"`
	RecurrenceRule []string     `json:"recurrenceRule"`
}

// AmountObject is a structure of JSON object of amount of medication.
type AmountObject struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}
