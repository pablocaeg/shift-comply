// Command genmap generates an SVG map of covered US jurisdictions
// from the actual jurisdiction registry. Run after adding new jurisdictions:
//
//	go run ./cmd/genmap > assets/us-coverage.svg
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/pablocaeg/shift-comply/comply"
	_ "github.com/pablocaeg/shift-comply/jurisdictions"
)

// US state positions (approximate x,y center on a 800x500 map)
var statePositions = map[string][2]int{
	"US-CA": {95, 250}, "US-NY": {700, 150}, "US-TX": {380, 380},
	"US-FL": {620, 410}, "US-IL": {500, 210}, "US-PA": {670, 180},
	"US-MA": {730, 130}, "US-WA": {120, 80}, "US-OR": {105, 140},
	"US-MN": {430, 110}, "US-OH": {590, 200}, "US-GA": {600, 340},
	"US-NC": {650, 290}, "US-MI": {540, 150}, "US-NJ": {700, 190},
	"US-VA": {650, 250}, "US-AZ": {200, 330}, "US-CO": {290, 240},
}

var stateNames = map[string]string{
	"US-CA": "California", "US-NY": "New York", "US-TX": "Texas",
	"US-FL": "Florida", "US-IL": "Illinois", "US-PA": "Pennsylvania",
	"US-MA": "Massachusetts", "US-WA": "Washington", "US-OR": "Oregon",
	"US-MN": "Minnesota", "US-OH": "Ohio", "US-GA": "Georgia",
	"US-NC": "North Carolina", "US-MI": "Michigan", "US-NJ": "New Jersey",
	"US-VA": "Virginia", "US-AZ": "Arizona", "US-CO": "Colorado",
}

func main() {
	// Find which US states are implemented
	covered := make(map[string]int) // code -> rule count
	for _, j := range comply.All() {
		code := string(j.Code)
		if strings.HasPrefix(code, "US-") {
			covered[code] = len(j.Rules)
		}
	}

	// Count totals
	us := comply.For(comply.US)
	federalRules := 0
	if us != nil {
		federalRules = len(us.Rules)
	}
	totalRules := 0
	totalJurisdictions := len(comply.All())
	for _, j := range comply.All() {
		totalRules += len(j.Rules)
	}

	var b strings.Builder
	b.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 800 520" font-family="system-ui,-apple-system,sans-serif">`)
	b.WriteString(`<rect width="800" height="520" fill="#fafafa" rx="12"/>`)

	// Title
	b.WriteString(`<text x="400" y="35" text-anchor="middle" font-size="16" font-weight="700" fill="#171717">Shift Comply Coverage</text>`)
	b.WriteString(fmt.Sprintf(`<text x="400" y="55" text-anchor="middle" font-size="12" fill="#737373">%d rules across %d jurisdictions</text>`, totalRules, totalJurisdictions))

	// US outline (simplified)
	b.WriteString(`<rect x="60" y="70" width="700" height="380" rx="8" fill="white" stroke="#e5e5e5" stroke-width="1"/>`)
	b.WriteString(fmt.Sprintf(`<text x="400" y="95" text-anchor="middle" font-size="11" fill="#a3a3a3" font-weight="600" letter-spacing="2">UNITED STATES · %d FEDERAL RULES (ALL STATES INHERIT)</text>`, federalRules))

	// Draw each state
	for code, pos := range statePositions {
		x, y := pos[0], pos[1]
		rules, isCovered := covered[code]
		name := stateNames[code]
		abbr := code[3:] // "CA", "NY", etc.

		if isCovered {
			// Covered state: green circle
			b.WriteString(fmt.Sprintf(`<circle cx="%d" cy="%d" r="28" fill="#ecfdf5" stroke="#10b981" stroke-width="2"/>`, x, y))
			b.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle" font-size="13" font-weight="700" fill="#059669">%s</text>`, x, y-4, abbr))
			b.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle" font-size="9" fill="#059669">%d rules</text>`, x, y+10, rules))
			_ = name
		} else {
			// Uncovered state: gray dashed circle
			b.WriteString(fmt.Sprintf(`<circle cx="%d" cy="%d" r="22" fill="none" stroke="#d4d4d4" stroke-width="1" stroke-dasharray="4,3"/>`, x, y))
			b.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle" font-size="11" fill="#d4d4d4">%s</text>`, x, y+4, abbr))
		}
	}

	// EU/Spain section
	b.WriteString(`<rect x="60" y="462" width="340" height="48" rx="6" fill="white" stroke="#e5e5e5"/>`)
	euRules := 0
	if eu := comply.For(comply.EU); eu != nil {
		euRules = len(eu.Rules)
	}
	esRules := 0
	if es := comply.For(comply.ES); es != nil {
		esRules = len(es.Rules)
	}
	b.WriteString(fmt.Sprintf(`<text x="80" y="490" font-size="11" fill="#525252" font-weight="600">EU</text><text x="100" y="490" font-size="10" fill="#737373">%d rules</text>`, euRules))
	b.WriteString(fmt.Sprintf(`<text x="170" y="490" font-size="11" fill="#525252" font-weight="600">Spain</text><text x="205" y="490" font-size="10" fill="#737373">%d rules</text>`, esRules))

	// Catalonia + Madrid
	ctRules := 0
	if ct := comply.For(comply.ESCT); ct != nil {
		ctRules = len(ct.Rules)
	}
	mdRules := 0
	if md := comply.For(comply.ESMD); md != nil {
		mdRules = len(md.Rules)
	}
	b.WriteString(fmt.Sprintf(`<text x="270" y="490" font-size="11" fill="#525252" font-weight="600">CAT</text><text x="295" y="490" font-size="10" fill="#737373">%d</text>`, ctRules))
	b.WriteString(fmt.Sprintf(`<text x="325" y="490" font-size="11" fill="#525252" font-weight="600">MAD</text><text x="355" y="490" font-size="10" fill="#737373">%d</text>`, mdRules))

	// Legend
	b.WriteString(`<rect x="420" y="462" width="340" height="48" rx="6" fill="white" stroke="#e5e5e5"/>`)
	b.WriteString(`<circle cx="440" cy="486" r="8" fill="#ecfdf5" stroke="#10b981" stroke-width="1.5"/>`)
	b.WriteString(`<text x="455" y="490" font-size="10" fill="#737373">Covered jurisdiction</text>`)
	b.WriteString(`<circle cx="580" cy="486" r="8" fill="none" stroke="#d4d4d4" stroke-width="1" stroke-dasharray="3,2"/>`)
	b.WriteString(`<text x="595" y="490" font-size="10" fill="#d4d4d4">Planned</text>`)

	b.WriteString(`</svg>`)

	fmt.Fprint(os.Stdout, b.String())
}
