package comply

import "sync"

// Code is a jurisdiction identifier following ISO 3166-1 alpha-2 for countries,
// with hyphenated subdivision codes (e.g., "US-CA", "ES-MD").
type Code string

// Jurisdiction codes for implemented jurisdictions.
const (
	US   Code = "US"
	USAL Code = "US-AL"
	USAK Code = "US-AK"
	USAZ Code = "US-AZ"
	USAR Code = "US-AR"
	USCA Code = "US-CA"
	USCO Code = "US-CO"
	USCT Code = "US-CT"
	USDE Code = "US-DE"
	USFL Code = "US-FL"
	USGA Code = "US-GA"
	USHI Code = "US-HI"
	USID Code = "US-ID"
	USIL Code = "US-IL"
	USIN Code = "US-IN"
	USIA Code = "US-IA"
	USKS Code = "US-KS"
	USKY Code = "US-KY"
	USLA Code = "US-LA"
	USME Code = "US-ME"
	USMD Code = "US-MD"
	USMA Code = "US-MA"
	USMI Code = "US-MI"
	USMN Code = "US-MN"
	USMS Code = "US-MS"
	USMO Code = "US-MO"
	USMT Code = "US-MT"
	USNE Code = "US-NE"
	USNV Code = "US-NV"
	USNH Code = "US-NH"
	USNJ Code = "US-NJ"
	USNM Code = "US-NM"
	USNY Code = "US-NY"
	USNC Code = "US-NC"
	USND Code = "US-ND"
	USOH Code = "US-OH"
	USOK Code = "US-OK"
	USOR Code = "US-OR"
	USPA Code = "US-PA"
	USRI Code = "US-RI"
	USSC Code = "US-SC"
	USSD Code = "US-SD"
	USTN Code = "US-TN"
	USTX Code = "US-TX"
	USUT Code = "US-UT"
	USVT Code = "US-VT"
	USVA Code = "US-VA"
	USWA Code = "US-WA"
	USWV Code = "US-WV"
	USWI Code = "US-WI"
	USWY Code = "US-WY"
	EU   Code = "EU"
	AT   Code = "AT"
	BE   Code = "BE"
	BG   Code = "BG"
	CH   Code = "CH"
	CY   Code = "CY"
	CZ   Code = "CZ"
	DE   Code = "DE"
	DK   Code = "DK"
	EE   Code = "EE"
	FI   Code = "FI"
	FR   Code = "FR"
	GR   Code = "GR"
	HR   Code = "HR"
	HU   Code = "HU"
	IE   Code = "IE"
	IS   Code = "IS"
	IT   Code = "IT"
	LT   Code = "LT"
	LU   Code = "LU"
	LV   Code = "LV"
	MT   Code = "MT"
	NL   Code = "NL"
	NO   Code = "NO"
	PL   Code = "PL"
	PT   Code = "PT"
	RO   Code = "RO"
	SE   Code = "SE"
	SI   Code = "SI"
	SK   Code = "SK"
	ES   Code = "ES"
	ESMD Code = "ES-MD"
	ESCT Code = "ES-CT"
)

// JurisdictionType classifies the hierarchical level of a jurisdiction.
type JurisdictionType string

const (
	// Jurisdiction type constants.
	Supranational JurisdictionType = "supranational"
	Country       JurisdictionType = "country"
	State         JurisdictionType = "state"
	Region        JurisdictionType = "region"
)

// JurisdictionDef defines a jurisdiction and all its scheduling regulations.
type JurisdictionDef struct {
	// Code uniquely identifies this jurisdiction.
	Code Code `json:"code"`

	// Name is the English name.
	Name string `json:"name"`

	// LocalName is the name in the local language, if different.
	LocalName string `json:"local_name,omitempty"`

	// Type classifies this jurisdiction (country, state, region, supranational).
	Type JurisdictionType `json:"type"`

	// Parent is the parent jurisdiction code. Empty for top-level.
	Parent Code `json:"parent,omitempty"`

	// Currency is the ISO 4217 currency code.
	Currency string `json:"currency"`

	// TimeZone is the primary IANA timezone.
	TimeZone string `json:"timezone"`

	// Rules contains all scheduling regulations for this jurisdiction.
	Rules []*RuleDef `json:"rules"`
}

var (
	registryMu    sync.RWMutex
	jurisdictions = make(map[Code]*JurisdictionDef)
)

// RegisterJurisdiction adds a jurisdiction to the global registry.
// It is called by jurisdiction packages in their init() functions.
func RegisterJurisdiction(j *JurisdictionDef) {
	registryMu.Lock()
	defer registryMu.Unlock()
	jurisdictions[j.Code] = j
}

// For returns the jurisdiction definition for the given code, or nil.
func For(code Code) *JurisdictionDef {
	registryMu.RLock()
	defer registryMu.RUnlock()
	return jurisdictions[code]
}

// All returns all registered jurisdictions in no guaranteed order.
func All() []*JurisdictionDef {
	registryMu.RLock()
	defer registryMu.RUnlock()
	result := make([]*JurisdictionDef, 0, len(jurisdictions))
	for _, j := range jurisdictions {
		result = append(result, j)
	}
	return result
}

// Codes returns all registered jurisdiction codes.
func Codes() []Code {
	registryMu.RLock()
	defer registryMu.RUnlock()
	codes := make([]Code, 0, len(jurisdictions))
	for c := range jurisdictions {
		codes = append(codes, c)
	}
	return codes
}

// ParentDef returns the parent jurisdiction definition, or nil.
func (j *JurisdictionDef) ParentDef() *JurisdictionDef {
	if j.Parent == "" {
		return nil
	}
	return For(j.Parent)
}

// Chain returns this jurisdiction and all its ancestors,
// ordered from most specific (self) to least specific (root).
func (j *JurisdictionDef) Chain() []*JurisdictionDef {
	var chain []*JurisdictionDef
	for cur := j; cur != nil; cur = cur.ParentDef() {
		chain = append(chain, cur)
	}
	return chain
}
