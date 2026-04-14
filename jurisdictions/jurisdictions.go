// Package jurisdictions imports all jurisdiction packages, registering their
// regulatory data at init time.
//
// Import this package with a blank identifier to load all jurisdictions:
//
//	import _ "github.com/pablocaeg/shift-comply/jurisdictions"
package jurisdictions

import (
	_ "github.com/pablocaeg/shift-comply/jurisdictions/es"    // ES jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/es_ct" // ES-CT jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/es_md" // ES-MD jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/eu"    // EU jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us"    // US jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_ca" // US-CA jurisdiction
)
