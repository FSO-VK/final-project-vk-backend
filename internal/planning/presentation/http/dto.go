package http

// PlanObject is info about plan.
type PlanObject struct {
	MedicationID   string       `json:"medicationId"`
	Amount         AmountObject `json:"amount"`
	Condition      string       `json:"condition"`
	Status         string       `json:"status"`
	StartDate      string       `json:"startDate"`
	EndDate        string       `json:"endDate"`
	RecurrenceRule []string     `json:"recurrenceRule"`
}

// AmountObject is a structure of JSON object of amount of medication.
type AmountObject struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}
