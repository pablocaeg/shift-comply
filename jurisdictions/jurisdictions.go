// Package jurisdictions imports all jurisdiction packages, registering their
// regulatory data at init time.
//
// Import this package with a blank identifier to load all jurisdictions:
//
//	import _ "github.com/pablocaeg/shift-comply/jurisdictions"
package jurisdictions

import (
	_ "github.com/pablocaeg/shift-comply/jurisdictions/at"    // AT jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/be"    // BE jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/ch"    // CH jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/cz"    // CZ jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/de"    // DE jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/dk"    // DK jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/es"    // ES jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/es_ct" // ES-CT jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/es_md" // ES-MD jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/eu"    // EU jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/fi"    // FI jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/fr"    // FR jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/gr"    // GR jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/hr"    // HR jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/hu"    // HU jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/ie"    // IE jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/it"    // IT jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/nl"    // NL jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/pl"    // PL jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/pt"    // PT jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/ro"    // RO jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/se"    // SE jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us"    // US jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_ca" // US-CA jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_fl" // US-FL jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_il" // US-IL jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_ma" // US-MA jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_ny" // US-NY jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_or" // US-OR jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_tx" // US-TX jurisdiction
)
