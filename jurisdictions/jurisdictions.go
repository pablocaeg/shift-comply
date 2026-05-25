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
	_ "github.com/pablocaeg/shift-comply/jurisdictions/bg"    // BG jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/ch"    // CH jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/cy"    // CY jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/cz"    // CZ jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/de"    // DE jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/dk"    // DK jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/ee"    // EE jurisdiction
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
	_ "github.com/pablocaeg/shift-comply/jurisdictions/is"    // IS jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/it"    // IT jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/lt"    // LT jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/lu"    // LU jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/lv"    // LV jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/mt"    // MT jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/nl"    // NL jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/no"    // NO jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/pl"    // PL jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/pt"    // PT jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/ro"    // RO jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/se"    // SE jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/si"    // SI jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/sk"    // SK jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us"    // US jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_ak" // US-AK jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_al" // US-AL jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_ar" // US-AR jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_az" // US-AZ jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_ca" // US-CA jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_co" // US-CO jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_ct" // US-CT jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_de" // US-DE jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_fl" // US-FL jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_ga" // US-GA jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_hi" // US-HI jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_ia" // US-IA jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_id" // US-ID jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_il" // US-IL jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_in" // US-IN jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_ks" // US-KS jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_ky" // US-KY jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_la" // US-LA jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_ma" // US-MA jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_md" // US-MD jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_me" // US-ME jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_mi" // US-MI jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_mn" // US-MN jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_mo" // US-MO jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_ms" // US-MS jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_mt" // US-MT jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_nc" // US-NC jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_nd" // US-ND jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_ne" // US-NE jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_nh" // US-NH jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_nj" // US-NJ jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_nm" // US-NM jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_nv" // US-NV jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_ny" // US-NY jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_oh" // US-OH jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_ok" // US-OK jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_or" // US-OR jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_pa" // US-PA jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_ri" // US-RI jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_sc" // US-SC jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_sd" // US-SD jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_tn" // US-TN jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_tx" // US-TX jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_ut" // US-UT jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_va" // US-VA jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_vt" // US-VT jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_wa" // US-WA jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_wi" // US-WI jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_wv" // US-WV jurisdiction
	_ "github.com/pablocaeg/shift-comply/jurisdictions/us_wy" // US-WY jurisdiction
)
