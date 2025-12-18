//nolint:tagliatelle // this is a API's JSON format, can't change tags to camelCase.
package vidal

// Product represents the main product data.
type Product struct {
	FullForm                []string         `json:"fullForm"`
	InfoPages               []interface{}    `json:"infoPages"`
	ATCCodes                []ATCCode        `json:"atcCodes"`
	PhthGroups              []PhthGroup      `json:"phthgroups"`
	ClPhGroups              []ClPhGroup      `json:"ClPhGroups"`
	UpdatedAt               string           `json:"updatedAt"`
	IsValid                 bool             `json:"isValid"`
	ProductURL              string           `json:"productUrl"`
	Images                  []string         `json:"images"`
	MoleculeNames           []MoleculeName   `json:"moleculeNames"`
	ID                      int              `json:"id"`
	RusName                 string           `json:"rusName"`
	RusNameClear            string           `json:"rusNameClear"`
	EngName                 string           `json:"engName"`
	NonPrescriptionDrug     bool             `json:"nonPrescriptionDrug"`
	RegistrationDate        string           `json:"registrationDate"`
	DateOfCloseRegistration string           `json:"dateOfCloseRegistration"`
	RegistrationNumber      string           `json:"registrationNumber"`
	ZipInfo                 string           `json:"zipInfo"`
	Composition             string           `json:"composition"`
	ProductTypeCode         string           `json:"productTypeCode"`
	MarketStatus            MarketStatus     `json:"marketStatus"`
	DateOfReRegistration    string           `json:"dateOfReRegistration"`
	GNVLS                   bool             `json:"gnvls"`
	ListPkkn                string           `json:"listPkkn"`
	StrongMeans             bool             `json:"strongMeans"`
	Poison                  bool             `json:"poison"`
	Companies               []CompanyInfo    `json:"companies"`
	Document                Document         `json:"document"`
	Children                []interface{}    `json:"childrens"`
	Forms                   string           `json:"forms"`
	ProductPackages         []ProductPackage `json:"productPackages"`
	RegPreviousNumber       string           `json:"regPreviousNumber"`
	AnalogsUpdatedAt        string           `json:"analogsUpdatedAt"`
	BarcodeUpdatedAt        string           `json:"barcodeUpdatedAt"`
}

// ATCCode represents an ATC code entry.
type ATCCode struct {
	Code       string `json:"code"`
	RusName    string `json:"rusName"`
	EngName    string `json:"engName"`
	ParentCode string `json:"parentCode"`
}

// PhthGroup represents a pharmacotherapeutic group.
type PhthGroup struct {
	Code string `json:"code"`
}

// ClPhGroup represents a clinical-pharmacological group.
type ClPhGroup struct {
	Name string `json:"name"`
}

// MoleculeName represents the name of a molecule.
type MoleculeName struct {
	ID       int      `json:"id"`
	Molecule Molecule `json:"molecule"`
	RusName  string   `json:"rusName"`
	EngName  string   `json:"engName"`
}

// Molecule represents the details of a molecule.
type Molecule struct {
	ID                      int      `json:"id"`
	LatName                 string   `json:"latName"`
	OriginalLatName         string   `json:"originalLatName"`
	GenitiveRusName         string   `json:"genitiveRusName"`
	GenitiveOriginalLatName string   `json:"genitiveOriginalLatName"`
	RusName                 string   `json:"rusName"`
	GNParent                GNParent `json:"GNParent"`
}

// GNParent represents the parent information for a generic name.
type GNParent struct {
	GNParent    string `json:"GNParent"`
	Description string `json:"description"`
}

// MarketStatus represents the market status of the product.
type MarketStatus struct {
	ID      int    `json:"id"`
	RusName string `json:"rusName"`
}

// CompanyInfo represents information about a company related to the product.
type CompanyInfo struct {
	IsRegistrationCertificate bool    `json:"isRegistrationCertificate"`
	IsManufacturer            bool    `json:"isManufacturer"`
	Company                   Company `json:"company"`
	CompanyRusNote            string  `json:"companyRusNote"`
}

// Company represents company details.
type Company struct {
	Name     string  `json:"name"`
	GDDBName string  `json:"GDDBName"`
	Country  Country `json:"country"`
}

// Country represents country details.
type Country struct {
	Code    string `json:"code"`
	RusName string `json:"rusName"`
}

// Document represents detailed documentation about the product.
type Document struct {
	StorageCondition   string        `json:"storageCondition"`
	StorageTime        string        `json:"storageTime"`
	UpdatedAt          string        `json:"updatedAt"`
	Nozologies         []Nozology    `json:"nozologies"`
	ClphPointers       []ClphPointer `json:"clphPointers"`
	Companies          []interface{} `json:"companies"`
	DocumentID         int           `json:"documentId"`
	ArticleID          int           `json:"articleId"`
	YearEdition        string        `json:"yearEdition"`
	PhInfluence        string        `json:"phInfluence"`
	PhKinetics         string        `json:"phKinetics"`
	Dosage             string        `json:"dosage"`
	OverDosage         string        `json:"overDosage"`
	Interaction        string        `json:"interaction"`
	Lactation          string        `json:"lactation"`
	SideEffects        string        `json:"sideEffects"`
	Indication         string        `json:"indication"`
	ContraIndication   string        `json:"contraIndication"`
	SpecialInstruction string        `json:"specialInstruction"`
	PregnancyUsing     string        `json:"pregnancyUsing"`
	NursingUsing       string        `json:"nursingUsing"`
	RenalInsuf         string        `json:"renalInsuf"`
	RenalInsufUsing    string        `json:"renalInsufUsing"`
	HepatoInsuf        string        `json:"hepatoInsuf"`
	HepatoInsufUsing   string        `json:"hepatoInsufUsing"`
	PharmDelivery      string        `json:"pharmDelivery"`
	ElderlyInsuf       string        `json:"elderlyInsuf"`
	ElderlyInsufUsing  string        `json:"elderlyInsufUsing"`
	ChildInsuf         string        `json:"childInsuf"`
	ChildInsufUsing    string        `json:"childInsufUsing"`
}

// Nozology represents a nosology entry.
type Nozology struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// ClphPointer represents a clinical-pharmacological pointer.
type ClphPointer struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// ProductPackage represents a packaging variant of the product.
type ProductPackage struct {
	ProductID int    `json:"productId"`
	ID        string `json:"id"`
	BarCode   string `json:"barCode"`
}
