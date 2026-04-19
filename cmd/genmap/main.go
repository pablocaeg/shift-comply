// Command genmap generates an SVG map of covered US jurisdictions.
// Run after adding jurisdictions:
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

// Simplified US state outlines as SVG paths (viewBox 0 0 960 600)
// Source: public domain US state boundaries, heavily simplified
var statePaths = map[string]string{
	"CA": "M122,270 L118,380 L133,400 L133,410 L100,500 L60,490 L50,400 L55,350 L65,280 L80,230 L100,210 L122,270Z",
	"NY": "M810,120 L830,115 L850,130 L845,155 L830,170 L800,175 L785,160 L780,140 L790,125 L810,120Z",
	"TX": "M350,340 L430,340 L460,350 L480,380 L470,430 L450,470 L410,490 L370,480 L340,460 L310,430 L310,390 L320,360 L350,340Z",
	"FL": "M640,380 L680,370 L710,385 L720,420 L700,470 L670,490 L650,470 L640,440 L635,400 L640,380Z",
	"WA": "M100,60 L155,55 L160,110 L100,115 L90,90 L100,60Z",
	"OR": "M65,115 L160,110 L155,175 L115,185 L55,165 L50,135 L65,115Z",
	"NV": "M130,180 L165,175 L175,300 L120,270 L100,210 L130,180Z",
	"AZ": "M140,340 L210,330 L225,420 L190,450 L130,440 L110,400 L115,360 L140,340Z",
	"UT": "M185,190 L230,185 L235,280 L190,285 L175,260 L185,190Z",
	"CO": "M240,230 L330,225 L335,295 L240,300 L240,230Z",
	"NM": "M225,335 L310,330 L315,430 L225,435 L225,335Z",
	"ID": "M170,75 L210,70 L215,180 L185,190 L160,170 L165,110 L170,75Z",
	"MT": "M215,55 L335,50 L340,115 L215,120 L210,70 L215,55Z",
	"WY": "M230,120 L330,115 L335,195 L235,200 L230,120Z",
	"ND": "M340,55 L445,55 L445,110 L340,110 L340,55Z",
	"SD": "M340,115 L445,115 L445,175 L340,175 L340,115Z",
	"NE": "M340,180 L445,180 L450,230 L335,235 L340,180Z",
	"KS": "M340,240 L450,235 L455,300 L340,305 L340,240Z",
	"OK": "M340,310 L455,305 L465,345 L430,340 L350,340 L340,310Z",
	"MN": "M450,60 L520,55 L525,140 L450,145 L450,60Z",
	"IA": "M450,150 L530,145 L535,210 L455,215 L450,150Z",
	"MO": "M460,220 L545,215 L555,300 L470,305 L460,220Z",
	"AR": "M470,310 L545,305 L550,370 L475,375 L470,310Z",
	"LA": "M480,380 L540,375 L560,430 L520,450 L480,430 L480,380Z",
	"WI": "M520,65 L575,60 L585,145 L530,150 L520,65Z",
	"IL": "M545,155 L580,150 L590,260 L555,265 L545,155Z",
	"IN": "M585,155 L620,150 L625,250 L590,255 L585,155Z",
	"OH": "M625,155 L670,150 L675,235 L630,240 L625,155Z",
	"MI": "M560,65 L620,55 L625,140 L585,145 L575,100 L560,65Z",
	"KY": "M580,265 L670,245 L690,275 L600,290 L580,265Z",
	"TN": "M570,290 L695,280 L700,310 L575,320 L570,290Z",
	"MS": "M555,325 L580,320 L590,400 L560,410 L555,325Z",
	"AL": "M595,320 L640,315 L650,400 L630,410 L595,400 L595,320Z",
	"GA": "M645,310 L700,305 L710,380 L680,395 L650,395 L645,310Z",
	"SC": "M700,300 L745,285 L750,320 L715,340 L700,300Z",
	"NC": "M690,270 L790,250 L800,280 L710,295 L690,270Z",
	"VA": "M680,235 L780,210 L790,245 L700,265 L680,235Z",
	"WV": "M670,220 L700,215 L705,260 L680,265 L670,220Z",
	"PA": "M700,155 L790,140 L795,185 L710,195 L700,155Z",
	"NJ": "M795,170 L810,165 L815,205 L800,210 L795,170Z",
	"CT": "M820,145 L845,140 L848,162 L823,165 L820,145Z",
	"MA": "M825,125 L870,118 L875,140 L830,145 L825,125Z",
	"VT": "M810,80 L825,78 L828,120 L813,122 L810,80Z",
	"NH": "M830,75 L843,73 L847,118 L832,120 L830,75Z",
	"ME": "M850,40 L880,35 L885,100 L855,105 L850,40Z",
	"MD": "M740,210 L790,200 L795,225 L745,232 L740,210Z",
	"DE": "M800,200 L812,198 L814,225 L802,227 L800,200Z",
	"RI": "M855,140 L867,138 L868,155 L856,157 L855,140Z",
}

func main() {
	covered := make(map[string]int)
	for _, j := range comply.All() {
		code := string(j.Code)
		if strings.HasPrefix(code, "US-") {
			covered[code[3:]] = len(j.Rules)
		}
	}

	us := comply.For(comply.US)
	federalRules := 0
	if us != nil {
		federalRules = len(us.Rules)
	}
	totalRules := 0
	for _, j := range comply.All() {
		totalRules += len(j.Rules)
	}

	var b strings.Builder

	// SVG header
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 960 620" font-family="system-ui,-apple-system,sans-serif">`)
	fmt.Fprintf(&b, `<defs><style>`)
	fmt.Fprintf(&b, `.state-covered{fill:#dcfce7;stroke:#16a34a;stroke-width:2}`)
	fmt.Fprintf(&b, `.state-default{fill:#f5f5f5;stroke:#d4d4d4;stroke-width:1}`)
	fmt.Fprintf(&b, `.label{font-size:11px;font-weight:700;text-anchor:middle;pointer-events:none}`)
	fmt.Fprintf(&b, `.label-covered{fill:#15803d}`)
	fmt.Fprintf(&b, `.label-default{fill:#a3a3a3}`)
	fmt.Fprintf(&b, `.rules-count{font-size:9px;font-weight:400;fill:#16a34a;text-anchor:middle}`)
	fmt.Fprintf(&b, `</style></defs>`)

	// Background
	fmt.Fprintf(&b, `<rect width="960" height="620" fill="white" rx="12"/>`)

	// Title
	fmt.Fprintf(&b, `<text x="480" y="30" text-anchor="middle" font-size="15" font-weight="700" fill="#171717">Shift Comply: US Coverage</text>`)
	fmt.Fprintf(&b, `<text x="480" y="48" text-anchor="middle" font-size="11" fill="#737373">%d federal rules inherited by all states | %d total rules across %d jurisdictions</text>`, federalRules, totalRules, len(comply.All()))

	// Draw states
	for abbr, path := range statePaths {
		rules, isCovered := covered[abbr]
		cls := "state-default"
		lblCls := "label-default"
		if isCovered {
			cls = "state-covered"
			lblCls = "label-covered"
		}

		fmt.Fprintf(&b, `<path d="%s" class="%s"/>`, path, cls)

		// Find center of path for label (approximate from first point)
		cx, cy := pathCenter(path)
		fmt.Fprintf(&b, `<text x="%d" y="%d" class="label %s">%s</text>`, cx, cy, lblCls, abbr)
		if isCovered {
			fmt.Fprintf(&b, `<text x="%d" y="%d" class="rules-count">%d rules</text>`, cx, cy+12, rules)
		}
	}

	// Legend
	fmt.Fprintf(&b, `<rect x="30" y="560" width="900" height="45" rx="6" fill="#fafafa" stroke="#e5e5e5"/>`)
	fmt.Fprintf(&b, `<rect x="50" y="576" width="16" height="16" rx="3" fill="#dcfce7" stroke="#16a34a" stroke-width="1.5"/>`)
	fmt.Fprintf(&b, `<text x="72" y="589" font-size="11" fill="#525252">Covered (state-specific rules)</text>`)
	fmt.Fprintf(&b, `<rect x="250" y="576" width="16" height="16" rx="3" fill="#f5f5f5" stroke="#d4d4d4"/>`)
	fmt.Fprintf(&b, `<text x="272" y="589" font-size="11" fill="#a3a3a3">Federal rules only (no state-specific)</text>`)
	fmt.Fprintf(&b, `<text x="900" y="589" text-anchor="end" font-size="10" fill="#a3a3a3">Regenerate: go run ./cmd/genmap > assets/us-coverage.svg</text>`)

	fmt.Fprintf(&b, `</svg>`)

	_, _ = os.Stdout.WriteString(b.String())
}

func pathCenter(path string) (int, int) {
	var x, y int
	path = strings.TrimPrefix(path, "M")
	parts := strings.Split(path, " ")
	if len(parts) > 0 {
		coords := strings.Split(parts[0], ",")
		if len(coords) >= 2 {
			_, _ = fmt.Sscanf(coords[0], "%d", &x)
			_, _ = fmt.Sscanf(coords[1], "%d", &y)
		}
	}
	if len(parts) > 2 {
		var x2, y2 int
		c2 := strings.Split(strings.TrimRight(parts[2], "LlZz"), ",")
		if len(c2) >= 2 {
			_, _ = fmt.Sscanf(c2[0], "%d", &x2)
			_, _ = fmt.Sscanf(c2[1], "%d", &y2)
			x = (x + x2) / 2
			y = (y + y2) / 2
		}
	}
	return x, y
}
