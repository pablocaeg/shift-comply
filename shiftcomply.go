// Package shiftcomply provides a structured, machine-readable database of
// healthcare scheduling regulations across jurisdictions.
//
// Shift Comply models healthcare
// labor law as structured data: each jurisdiction (country, state, or region)
// registers its scheduling rules at init time, making them queryable,
// comparable, and suitable for constraint generation in scheduling optimizers.
//
// All regulation data is compiled into the binary - no database required.
// Every rule carries its legal citation and effective date.
//
// Usage:
//
//	import (
//	    sc "github.com/pablocaeg/shift-comply/comply"
//	    _ "github.com/pablocaeg/shift-comply/jurisdictions"
//	)
//
//	// Get a specific jurisdiction
//	j := sc.For("US-CA")
//
//	// Get all effective rules including parent jurisdiction
//	rules := sc.EffectiveRules("US-CA", sc.ForStaff(sc.StaffNurseRN))
//
//	// Compare two jurisdictions
//	diff := sc.Compare("US-CA", "US-TX")
package shiftcomply
