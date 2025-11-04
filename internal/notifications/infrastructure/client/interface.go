package client

// ScannedInfoFromDataMatrix is a request to find medication info in API.
type ScannedInfoFromDataMatrix struct {
	GTIN         string
	SerialNumber string
	CryptoData91 string
	CryptoData92 string
}

// MedicationInfoFromAPI is the response for the CheckAuth method.
type MedicationInfoFromAPI struct {
	ExpDate string
}

// ExpectedDataMatrixAPIResponse is the expected response for dataMatrix API.
type ExpectedDataMatrixAPIResponse struct {
	CodeFound     bool      `json:"codeFounded"`
	CheckDate     int64     `json:"checkDate"`
	Category      string    `json:"category"`
	ProductName   string    `json:"productName"`
	ProducerName  string    `json:"producerName"`
	Status        string    `json:"status"`
	ExpDate       string    `json:"expDate"`
	OperationDate int64     `json:"operationDate"`
	Batch         string    `json:"batch"`
	Warning       string    `json:"warning"`
	IsPerishable  bool      `json:"isPerishable"`
	DrugsData     DrugsData `json:"drugsData"`
}

// DrugsData represents Drugs data from API.
type DrugsData struct {
	ProductDescLabel      string    `json:"prodDescLabel"`
	ProducerName          string    `json:"packingName"`
	Batch                 string    `json:"batch"`
	Status                string    `json:"status"`
	LastOperationDate     int64     `json:"lastOperationDate"`
	ReceiptDate           int64     `json:"receiptDate"`
	SourceType            int       `json:"sourceType"`
	UtilizationOpDate     int64     `json:"utilizationOpDate"`
	EmissionOperationDate int64     `json:"emissionOperationDate"`
	ContainsVZN           bool      `json:"containsVzn"`
	ReleaseDate           int64     `json:"releaseDate"`
	IsPerishable          bool      `json:"isPerishable"`
	FOIV                  FOIVData  `json:"foiv"`
	VidalData             VidalData `json:"vidalData"`
}

// FOIVData represents FOIV data from API.
type FOIVData struct {
	ProductFormName     string `json:"prodFormNormName"`
	Dosage              string `json:"prodDNormName"`
	PackageType         string `json:"prodPack1Name"`
	PackageQuantity     string `json:"prodPack12"`
	PackageSize         string `json:"prodPack1Size"`
	Manufacturer        string `json:"glfName"`
	ManufacturerCountry string `json:"glfCountry"`
	RegHolder           string `json:"regHolder"`
	RegNumber           string `json:"regNumber"`
	RegDate             int64  `json:"regDate"`
	ProductName         string `json:"prodNormName"`
}

// VidalData represents Vidal data from API.
type VidalData struct {
	Pharmacology string `json:"phKinetics"`
}
