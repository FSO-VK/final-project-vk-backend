package http

// PlanObject is info about plan.
type PlanObject struct {
	MedicationID   string   `json:"medicationId"`
	UserID         string   `json:"userId"`
	AmountValue    float64  `json:"amountValue"`
	AmountUnit     string   `json:"amountUnit"`
	Condition      string   `json:"condition"`
	StartDate      string   `json:"startDate"`
	EndDate        string   `json:"endDate"`
	RecurrenceRule []string `json:"recurrenceRule"`
}
