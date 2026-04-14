# Shift Comply

[![Test](https://github.com/pablocaeg/shift-comply/actions/workflows/test.yaml/badge.svg)](https://github.com/pablocaeg/shift-comply/actions/workflows/test.yaml)
[![Lint](https://github.com/pablocaeg/shift-comply/actions/workflows/lint.yaml/badge.svg)](https://github.com/pablocaeg/shift-comply/actions/workflows/lint.yaml)
[![codecov](https://codecov.io/gh/pablocaeg/shift-comply/graph/badge.svg)](https://codecov.io/gh/pablocaeg/shift-comply)
[![Go Reference](https://pkg.go.dev/badge/github.com/pablocaeg/shift-comply.svg)](https://pkg.go.dev/github.com/pablocaeg/shift-comply/comply)

Open-source, machine-readable database of healthcare scheduling regulations across jurisdictions.

**[Live Demo](https://pablocaeg.github.io/shift-comply/)** | Every regulation carries a real legal citation, an effective date, and is queryable by jurisdiction, staff type, and hospital unit. All data compiles into the binary. No database required.

## The problem

Healthcare labor laws vary enormously across jurisdictions:

- The US has ACGME rules limiting residents to 80 hours/week, but states layer on top: California mandates nurse-patient ratios for 17 unit types while Texas has almost no state-level rules.
- The EU Working Time Directive caps working time at 48 hours/week, but member states implement it differently, with varying opt-out provisions and on-call definitions shaped by CJEU case law (SIMAP, Jaeger, Matzak).
- Spain has national rules (Estatuto de los Trabajadores, Estatuto Marco, MIR residency regulations) plus autonomous community variations: Catalonia limits guards to 4/month, Andalusia is reforming guard duration from 24 to 17 hours, Madrid mandates 36-hour weekly rest after the 2022 Supreme Court ruling.

Scheduling systems like [YouShift](https://www.you-shift.com/) (YC W25) already have powerful constraint engines: max hours, mandatory rest, consecutive day limits, shift sequence rules, team composition requirements. The system can enforce all of these. But the constraint *values* have to come from somewhere. Today, hospital admins configure them manually, which assumes they know their own regulations, can translate legal language into scheduling parameters, and will keep up with changes across jurisdictions.

That works for a single hospital. It breaks when you scale internationally.

Shift Comply is the layer that answers: **what should those values be, and what law requires them?** When a hospital selects its jurisdiction, Shift Comply provides the legally correct constraint values with citations attached. The scheduling system consumes them directly.

## How it connects to scheduling systems

A scheduling optimizer like YouShift already supports constraints like:

| YouShift constraint | Shift Comply provides |
|---|---|
| Maximum hours per week | 80 for US residents (ACGME CPR VI.F.1), 48 for Spanish public health (Ley 55/2003 Art. 48) |
| Mandatory rest between shifts | 12h in Spain (Estatuto de los Trabajadores Art. 34.3), 8h recommended by ACGME |
| Maximum consecutive night shifts | 6 for US residents (ACGME), 2 weeks in Spain (Art. 36.3) |
| Monthly on-call limits | 7 guards nationally in Spain (RD 1146/2006), 4 in Catalonia (ICS III Acord) |
| Nurse staffing minimums | 1:2 in California ICU (Title 22 CCR S70217), no state mandate in Texas |
| Overtime thresholds | 8h daily in California (Labor Code S510), no daily threshold federally |
| Shift duration caps | 12h for CA non-resident healthcare (IWC Wage Order 5), 24h for ACGME residents |

The integration: when a hospital onboards and selects "Catalonia, public hospital," Shift Comply auto-populates every constraint value with the legally correct number from the ICS collective agreement, Spanish national law, and the EU Working Time Directive. No manual research, no guessing, every value traceable to a specific statute.

## What it does

**Query rules by jurisdiction, staff type, and unit**. Get all regulations that apply to an ICU nurse in California, including inherited federal rules.

**Generate optimizer-ready constraints**. Turn legal rules into structured JSON that scheduling engines can consume directly, with constraint types, limits, averaging periods, and citations.

**Compare jurisdictions**. "What changes if you expand from Madrid to New York?" Get a structured diff showing rules that exist in one jurisdiction but not the other, and rules where the values differ.

**Track regulation history**. California's step-down nurse ratio was 1:4 in 2004, tightened to 1:3 in 2008. Query any historical date.

**Document regulatory absences**. Texas explicitly registers that it has *no* daily overtime law and *no* break requirements. Florida registers that it does *not* prohibit mandatory nurse overtime. These documented gaps are critical for jurisdiction comparison.

## Coverage

| Jurisdiction | Rules | Key Sources |
|---|---|---|
| **US Federal** | 17 | ACGME CPR Section VI (accreditation, not statute), FLSA 29 U.S.C. S207, VA Personnel Enhancement Act |
| **California** | 32 | Title 22 CCR S70217 (nurse ratios), Labor Code S510, IWC Wage Order 5-2001 S3(H), Gov. Code S19851.2 |
| **EU** | 8 | Directive 2003/88/EC, CJEU case law (SIMAP, Jaeger, Matzak) |
| **Spain** | 21 | RDL 2/2015 (Estatuto de los Trabajadores), Ley 55/2003 (Estatuto Marco, public health only), RD 1146/2006 (MIR) |
| **Catalonia** | 6 | III Acord ICS (DOGC Jan 23, 2024, Resolucio EMT/74/2024). 4 guards/month, age exemptions, post-guard rest |
| **Madrid** | 3 | SERMAS Resolution Feb 26, 2021 + STS 280/2022 (March 30, 2022). 36-hour weekly rest, 1,642.5 annual hours |

**87 rules across 6 jurisdictions.** Every rule has a verified legal citation and correct facility scope. Spanish region rules are verified against DOGC/BOE publications, Supreme Court rulings, and official union/health service sources.

## Install

```
go get github.com/pablocaeg/shift-comply
```

## Usage

### As a Go library

```go
import (
    "github.com/pablocaeg/shift-comply/comply"
    _ "github.com/pablocaeg/shift-comply/jurisdictions"
)

// Get a jurisdiction
ca := comply.For("US-CA")

// Query effective rules (includes inherited federal rules)
rules := comply.EffectiveRules("US-CA",
    comply.ForStaff(comply.StaffNurseRN),
    comply.ForUnit(comply.UnitICU),
)

// Generate optimizer-ready constraints
constraints := comply.GenerateConstraints("US-CA",
    comply.ForStaff(comply.StaffNurseRN),
)

// Compare two jurisdictions
diff := comply.Compare("US-CA", "ES")

// Validate a schedule
report, err := comply.Validate(comply.Schedule{
    Jurisdiction:  "US-CA",
    FacilityScope: comply.ScopeHospitals,
    Shifts: []comply.Shift{
        {StaffID: "nurse-1", StaffType: comply.StaffNurseRN, UnitType: comply.UnitICU,
         Start: "2025-03-10T07:00:00", End: "2025-03-10T20:30:00"},
    },
})
// report.Result = "fail"
// report.Violations[0].Citation = "IWC Wage Order No. 5-2001..."
```

### CLI

```
go install github.com/pablocaeg/shift-comply/cmd/shiftcomply@latest
```

```sh
# List all jurisdictions
shiftcomply jurisdictions

# Query rules with filters
shiftcomply rules US-CA --staff nurse-rn --category staffing

# Compare two jurisdictions
shiftcomply compare US-CA ES

# Generate optimizer-ready JSON constraints
shiftcomply constraints ES --staff resident

# Full JSON export of a jurisdiction
shiftcomply export US-CA
```

### REST API

```sh
go install github.com/pablocaeg/shift-comply/cmd/shiftcomply-api@latest
shiftcomply-api              # listens on :8080
shiftcomply-api -addr :3000  # custom port
PORT=8080 shiftcomply-api    # or via PORT env var
```

Endpoints:

| Method | Path | Description |
|---|---|---|
| GET | `/jurisdictions` | List all jurisdictions |
| GET | `/rules?jurisdiction=US-CA&staff=nurse-rn&scope=hospitals` | Query rules with filters |
| GET | `/constraints?jurisdiction=ES&staff=resident` | Generate optimizer-ready constraints |
| GET | `/compare?left=US-CA&right=ES` | Compare two jurisdictions |
| POST | `/validate` | Validate a schedule (JSON body) |
| GET | `/export/US-CA` | Full JSON export of a jurisdiction |
| GET | `/health` | Health check |

CORS is enabled by default. All responses are JSON.

### WebAssembly

Build the entire regulation database + validation engine as a single .wasm file:

```sh
GOOS=js GOARCH=wasm go build -o shiftcomply.wasm ./cmd/wasm
# Copy the Go WASM support file
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" .
```

Use from JavaScript:

```html
<script src="wasm_exec.js"></script>
<script>
const go = new Go();
WebAssembly.instantiateStreaming(fetch("shiftcomply.wasm"), go.importObject)
  .then(result => {
    go.run(result.instance);

    // All 87 rules available client-side, zero API calls
    const rules = JSON.parse(shiftcomply.rules("US-CA", "nurse-rn"));
    const constraints = JSON.parse(shiftcomply.constraints("ES", "resident"));
    const diff = JSON.parse(shiftcomply.compare("US-CA", "ES"));

    // Validate a schedule in the browser
    const report = JSON.parse(shiftcomply.validate(JSON.stringify({
      jurisdiction: "US-CA",
      facility_scope: "hospitals",
      shifts: [
        {staff_id: "nurse-1", staff_type: "nurse-rn",
         start: "2025-03-10T07:00:00", end: "2025-03-10T20:30:00"}
      ]
    })));
  });
</script>
```

The .wasm file is ~3.7MB and contains the full regulation database, validation engine, and constraint generator. No server round-trips needed.

### Example: constraints output

```sh
shiftcomply constraints US-CA --staff nurse-rn --unit icu
```

```json
[
  {
    "type": "staffing_ratio",
    "scope": "shift",
    "limit": 2,
    "limit_unit": "patients_per_nurse",
    "operator": "lte",
    "enforcement": "mandatory",
    "citation": "California Code of Regulations, Title 22, S 70217",
    "jurisdiction": "US-CA",
    "rule_key": "nurse-patient-ratio-icu"
  },
  {
    "type": "max_shift",
    "scope": "shift",
    "limit": 12,
    "limit_unit": "hours",
    "operator": "lte",
    "enforcement": "mandatory",
    "citation": "IWC Wage Order No. 5-2001 (Healthcare Industry), Section 3(B)(8)",
    "jurisdiction": "US-CA",
    "rule_key": "max-shift-hours"
  }
]
```

### Example: jurisdiction comparison

```sh
shiftcomply compare US-CA ES --staff nurse-rn
```

```
Comparing US-CA vs ES

--- Only in US-CA (44 rules) ---
  nurse-patient-ratio-icu: <=2 patients_per_nurse
  nurse-patient-ratio-med-surg: <=5 patients_per_nurse
  overtime-daily-threshold: >=8 hours
  mandatory-overtime-prohibited: 1 boolean
  ...

--- Only in ES (19 rules) ---
  max-ordinary-weekly-hours: <=40 hours
  min-weekly-rest: >=36 hours
  max-overtime-annual: <=80 hours
  min-annual-leave-days: >=30 calendar_days
  ...

--- Different values (4 rules) ---
  meal-break-threshold: >=5 hours vs >=6 hours
  max-weekly-hours: <=80 hours vs <=48 hours
  min-rest-between-shifts: >=8 hours vs >=12 hours

Summary: 44 only-US-CA, 19 only-ES, 4 different, 1 same
```

## Architecture

Each jurisdiction is a self-contained Go package that registers its rules at init time:

```
comply/                  Core types, registry, query API, constraint generation
jurisdictions/           One package per jurisdiction (registered at init time)
  us/                    US Federal (ACGME, FLSA, VA)
  us_ca/                 California (nurse ratios, overtime, breaks, mandatory OT restrictions)
  eu/                    EU Working Time Directive 2003/88/EC
  es/                    Spain (Estatuto de los Trabajadores, Estatuto Marco, MIR residency)
  es_ct/                 Catalonia / ICS (III Acord, guard limits, age exemptions)
  es_md/                 Community of Madrid / SERMAS (36-hour rest, Supreme Court ruling)
cmd/shiftcomply/         CLI tool
```

Each jurisdiction registers via `init()` + `comply.RegisterJurisdiction()`. Blank-importing `jurisdictions` loads all of them. The parent chain (US-CA -> US, ES-MD -> ES -> EU) means child rules inherit from and can override parent rules with the same key.

### Core types

- **JurisdictionDef**: a jurisdiction with its code, name, type, parent, and rules
- **RuleDef**: a single regulation with key, category, operator, staff/unit scopes, enforcement level, source citation, and time-versioned values
- **RuleValue**: a value effective from a given date, with amount, unit, per-period, optional averaging window, and exceptions
- **Constraint**: optimizer-ready output generated from rules
- **Source**: legal citation (title, section, URL)

## Design decisions

**Compiled data, no database.** Regulations change slowly (yearly at most). Compiling into the binary means zero infrastructure, instant lookups, and the data is version-controlled alongside the code. When a regulation changes, it's a pull request with a diff, not a database migration.

**Time-versioned values.** Rules change. Rather than replacing old values, each rule stores a time-ordered list of values. This supports historical queries and makes regulation changes auditable through git history.

**Documented absences.** Knowing what a jurisdiction does *not* regulate is as important as knowing what it does. Texas and Florida explicitly register their regulatory gaps so comparison tools surface them.

**Hierarchical inheritance.** Query California and get both state rules and inherited federal rules. Query the Community of Madrid and get Madrid rules, Spanish national rules, and EU directive rules. Child overrides parent when both define the same rule key.

**Every rule has a citation.** No placeholder data. Every value traces back to a specific statute, regulation, or court ruling.

## Adding a new jurisdiction

Create a new package under `jurisdictions/`:

```go
package us_il

import (
    "time"
    "github.com/pablocaeg/shift-comply/comply"
)

func init() {
    comply.RegisterJurisdiction(New())
}

func New() *comply.JurisdictionDef {
    return &comply.JurisdictionDef{
        Code:     "US-IL",
        Name:     "Illinois",
        Type:     comply.State,
        Parent:   comply.US,
        Currency: "USD",
        TimeZone: "America/Chicago",
        Rules: []*comply.RuleDef{
            {
                Key:         comply.RuleMandatoryOTProhibited,
                Name:        "Mandatory Overtime Prohibition (Nurses)",
                // ... real legal citation required
                Source: comply.Source{
                    Title:   "Illinois Nurse Staffing by Patient Acuity Act",
                    Section: "210 ILCS 85/10.10",
                    URL:     "https://www.ilga.gov/...",
                },
            },
        },
    }
}
```

Then add a blank import in `jurisdictions/jurisdictions.go`.

Requirements for contributed jurisdictions:
- Every rule must have a real legal citation (statute, regulation, or court ruling)
- Values must be accurate and current
- Include effective dates
- Document regulatory absences where relevant for comparison
- Include exceptions where they exist

## Use cases

### Schedule validation API

The most direct application: an API that receives a schedule and returns every regulation it violates, with legal citations.

A hospital uploads their monthly shift calendar as JSON. The API checks every shift assignment against the jurisdiction's rules and returns a compliance report: which staff members exceed the weekly hour limit, where rest periods are too short, which units are understaffed relative to patient ratios.

```
POST /validate
{
  "jurisdiction": "US-CA",
  "facility_scope": "hospitals",
  "shifts": [
    {"staff_id": "nurse-42", "staff_type": "nurse-rn", "unit": "icu",
     "start": "2025-03-10T07:00:00", "end": "2025-03-10T19:30:00"},
    {"staff_id": "nurse-42", "staff_type": "nurse-rn", "unit": "icu",
     "start": "2025-03-11T07:00:00", "end": "2025-03-11T19:30:00"}
  ]
}

Response:
{
  "result": "fail",
  "violations": [
    {
      "rule_key": "nurse-patient-ratio-icu",
      "rule_name": "Nurse-Patient Ratio: ICU / Critical Care",
      "severity": "mandatory",
      "staff_id": "nurse-42",
      "message": "ICU shift on 2025-03-10 has 3 patients assigned, maximum is 2",
      "citation": "California Code of Regulations, Title 22, § 70217"
    }
  ]
}
```

This turns a legal research problem into an API call. The hospital doesn't need to know that California Title 22 CCR § 70217 mandates a 1:2 ICU ratio. The system knows.

### Pre-scheduling constraint feed

Before the optimizer runs, it queries Shift Comply for the constraints that apply to this hospital's jurisdiction, staff types, and units. The optimizer treats these as hard constraints (mandatory rules) or soft constraints (recommended rules) in its solver model.

This is the integration pattern for scheduling systems like YouShift: instead of hospital admins manually configuring "max 80 hours per week for residents, averaged over 4 weeks," the system pulls it automatically from the ACGME rules and feeds it to the optimizer with the legal citation attached.

### Expansion planning

A hospital network operating in Spain is evaluating expansion into California. They query the comparison API to see exactly what changes:

```sh
shiftcomply compare ES US-CA --staff nurse-rn
```

California has 17 mandatory nurse-patient ratios that Spain doesn't. California has daily overtime at 8 hours that Spain doesn't. Spain has a 30-day annual leave minimum that California doesn't. The output is a concrete, rule-by-rule diff that operations teams can plan around.

### Onboarding automation

When a new hospital signs up for a scheduling platform, the onboarding flow asks for jurisdiction and facility type. The system automatically loads the applicable constraints instead of requiring manual configuration. A public hospital in Catalonia gets the ICS collective agreement rules (4 guards/month, age exemptions) layered on top of Spanish national law and EU directives. A private hospital in Madrid gets only the Estatuto de los Trabajadores baseline, not the SERMAS rules.

### Compliance auditing

After a scheduling period ends, run the full schedule through the validation engine to produce an audit report. Every shift that touched a regulatory boundary is flagged with the specific statute or regulation it relates to. The report serves as documentation that the organization evaluated compliance, which matters for accreditation reviews (ACGME) and labor inspections.

### Regulatory change monitoring

When a regulation changes (a new collective agreement is signed, a court ruling modifies rest requirements), the change is a pull request to this repository. Any system consuming the data can diff the old and new versions to see exactly which constraints changed, for which staff types, in which jurisdictions. This is version-controlled regulatory intelligence.

## Roadmap

- [x] Core types, registry, and query API
- [x] 6 verified jurisdictions (US Federal, California, EU, Spain, Catalonia, Madrid) with 87 rules
- [x] Facility scope on every rule (public health, hospitals, accredited programs, VA, etc.)
- [x] Jurisdiction comparison and hierarchical inheritance
- [x] Constraint generation (optimizer-ready JSON output)
- [x] CLI tool with filtering and JSON export
- [x] JSON serialization on all types
- [x] Schedule validation engine (max shift hours, weekly hours with averaging, rest between shifts, days off, guard limits)
- [x] REST API (`shiftcomply-api`) with /jurisdictions, /rules, /constraints, /compare, /validate, /export endpoints
- [x] WASM build (3.7MB) exposing full API to JavaScript
- [x] CI: GitHub Actions (lint with golangci-lint, test with race detector, WASM build verification)
- [ ] More US states (NY, TX, FL, MA, IL)
- [ ] Spanish autonomous communities (Andalusia/SAS)
- [ ] More EU countries (France, Germany, Italy)
- [ ] Codecov integration
- [ ] GoReleaser for versioned releases

## Why Go

Shift Comply is a data library, not a web service. The language choice is driven by how the data gets consumed:

**Single binary, zero dependencies.** All 78 rules compile into the binary. No database, no config files, no runtime. Deploy as a sidecar, a CLI tool, or a Lambda function with nothing to install.

**Cross-compilation.** `GOOS=linux GOARCH=amd64 go build` produces a Linux binary from a Mac. Same for ARM, Windows, or any target. One `go build` command, no Docker required.

**WebAssembly.** `GOOS=js GOARCH=wasm go build` compiles the entire regulation database into a .wasm file that runs in a browser. A frontend can query regulations client-side with zero API calls.

**Fast startup, low memory.** The binary starts in milliseconds and uses minimal memory. This matters for serverless functions and sidecar containers that spin up per-request.

**init() registration pattern.** Go's init functions run at import time, which means jurisdiction packages register themselves automatically. Adding a new jurisdiction is a new package + one blank import line. No registry files to maintain, no dependency injection.

**Type safety for structured data.** Every rule key, staff type, unit type, scope, operator, and enforcement level is a typed constant. Typos are compile errors, not runtime bugs.

## Integration with scheduling systems

Shift Comply is designed to feed into scheduling optimizers. There are several integration patterns depending on the consuming system's architecture.

### Pattern 1: JSON at build time (simplest)

Generate static JSON files during CI/CD and bundle them with your application. Works with any language.

```sh
# Generate constraint files for each jurisdiction your system supports
shiftcomply constraints US-CA --staff nurse-rn > constraints/us-ca-nurse.json
shiftcomply constraints ES --staff resident > constraints/es-resident.json
shiftcomply export ES > jurisdictions/es.json
```

Your Python optimizer, Java backend, or TypeScript frontend reads these JSON files. No Go dependency at runtime. Update by re-running the CLI when regulations change.

### Pattern 2: HTTP sidecar

Wrap shift-comply in a thin HTTP server and run it alongside your backend. Any service can query it.

```go
// 15 lines to expose the full API over HTTP
http.HandleFunc("/rules", func(w http.ResponseWriter, r *http.Request) {
    code := comply.Code(r.URL.Query().Get("jurisdiction"))
    opts := []comply.QueryOption{}
    if s := r.URL.Query().Get("staff"); s != "" {
        opts = append(opts, comply.ForStaff(comply.Key(s)))
    }
    if s := r.URL.Query().Get("scope"); s != "" {
        opts = append(opts, comply.ForScope(comply.Scope(s)))
    }
    rules := comply.EffectiveRules(code, opts...)
    json.NewEncoder(w).Encode(rules)
})
```

Deploy as a container next to your main service. Sub-millisecond responses, no database.

### Pattern 3: Go library import (if your backend is in Go)

```go
import (
    "github.com/pablocaeg/shift-comply/comply"
    _ "github.com/pablocaeg/shift-comply/jurisdictions"
)

// In your scheduling optimizer
func buildConstraints(jurisdiction string, staffType string) []comply.Constraint {
    return comply.GenerateConstraints(
        comply.Code(jurisdiction),
        comply.ForStaff(comply.Key(staffType)),
        comply.ForScope(comply.ScopeHospitals),
    )
}
```

### Pattern 4: WebAssembly (browser)

Compile to WASM and call from JavaScript. The entire regulation database runs client-side.

```sh
GOOS=js GOARCH=wasm go build -o comply.wasm ./cmd/wasm
```

Useful for: rule explorer UIs, jurisdiction comparison tools, onboarding flows where a hospital selects its location and sees applicable regulations instantly.

### How constraints map to a scheduling optimizer

The `Constraint` struct is designed to be directly consumable by constraint-based scheduling solvers (OR-Tools, OptaPlanner, custom solvers). Each constraint tells the optimizer:

| Field | What the optimizer does with it |
|---|---|
| `type` | Maps to a constraint class (max_hours, min_rest, staffing_ratio, etc.) |
| `limit` + `operator` | The bound: `<=80 hours`, `>=11 hours`, `<=2 patients_per_nurse` |
| `time_scope` | What period the constraint applies to (per shift, per week, per month) |
| `averaged_over_days` | Whether the constraint is hard per-period or averaged (28 days = 4 weeks) |
| `staff_types` | Which worker roles this constraint binds |
| `unit_types` | Which hospital units this constraint applies to |
| `facility_scope` | Whether this applies to this specific facility type |
| `enforcement` | Whether the solver must satisfy it (mandatory) or optimize toward it (recommended) |
| `citation` | Attached to any violation report for legal traceability |

A typical integration: your optimizer loads constraints for the hospital's jurisdiction and staff types, adds them as hard/soft constraints to the solver model, and attaches citations to any violations in the output.

### Scope filtering for your facility

Not all rules apply to every facility. A private hospital in California queries differently than a VA hospital:

```go
// Private hospital in California
rules := comply.EffectiveRules("US-CA",
    comply.ForStaff(comply.StaffNurseRN),
    comply.ForScope(comply.ScopeHospitals),
)

// VA hospital in California
rules := comply.EffectiveRules("US-CA",
    comply.ForStaff(comply.StaffNurseRN),
    comply.ForScope(comply.ScopeVA),
)

// Public hospital in Spain
rules := comply.EffectiveRules("ES",
    comply.ForStaff(comply.StaffStatutory),
    comply.ForScope(comply.ScopePublicHealth),
)
```

## Context

This project is designed to complement [YouShift](https://www.you-shift.com/) (YC W25), which automates hospital shift scheduling with AI. As YouShift expands internationally, the regulatory complexity of healthcare labor law across jurisdictions becomes a core infrastructure problem. Shift Comply aims to be the open-source foundation for that layer, usable by YouShift and by anyone else building healthcare scheduling tools.

For full context on YouShift's product, market position, and gaps, see the [youshift-brain](https://github.com/pablocaeg/youshift-brain) knowledge base.

## License

Apache 2.0. See [LICENSE](LICENSE).
