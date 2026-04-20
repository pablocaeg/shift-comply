# Contributing to Shift Comply

## Adding a new jurisdiction

The fastest way: use the [new-jurisdiction agent](.claude/agents/README.md). It researches the legislation, presents findings for your review, then generates the Go package with tests. Run: `/agents/new-jurisdiction Add France`

Manual process:
1. Create a new package under `jurisdictions/` (e.g., `jurisdictions/fr/`).
2. Implement a `New()` function that returns `*comply.JurisdictionDef`.
3. Register it in `init()` with `comply.RegisterJurisdiction(New())`.
4. Add a blank import in `jurisdictions/jurisdictions.go`.
5. Add tests in `comply/comply_test.go` for the new jurisdiction.
6. Add healthcare stats to `website/src/lib/jurisdiction-data.ts` (`JURISDICTION_STATS`).
7. Update the coverage table in `README.md`.

### Requirements for jurisdiction data

Every rule must have:
- A real legal citation (statute, regulation, collective agreement, or court ruling)
- A source URL where the citation can be verified (official gazette, legislature website, legal database)
- An accurate numeric value with correct unit and operator
- An effective date
- A correct `Scope` field (hospitals, public_health, accredited_programs, etc.)
- Documented exceptions where they exist

Do not use placeholder data. Do not guess values. If you cannot verify a rule from a primary source, do not include it.

### Documenting regulatory absences

If a jurisdiction notably lacks a regulation that other jurisdictions have (e.g., no mandatory nurse overtime prohibition), register it as an advisory rule with `Amount: 0` and a source confirming the absence. This is important for the comparison feature.

## Development

### Prerequisites

- Go 1.23+
- golangci-lint (for linting)

### Running tests

```sh
go test ./...                    # all tests
go test -race ./...              # with race detector
go test -v ./comply/ -run Test   # verbose, specific package
```

### Running the linter

```sh
golangci-lint run ./...
```

### Building

```sh
go build ./...                                           # all packages
go build ./cmd/shiftcomply                               # CLI
go build ./cmd/shiftcomply-api                           # REST API
GOOS=js GOARCH=wasm go build -o shiftcomply.wasm ./cmd/wasm  # WASM
```

## Code style

- No em dashes or double dashes in prose.
- Use colons, commas, or periods instead.
- Go standard formatting (gofmt/goimports enforced by CI).
- Keep descriptions factual. Cite the law, not your interpretation of it.

## Pull request process

1. Fork the repository and create a feature branch.
2. Make your changes with tests.
3. Run `go test -race ./...` and `golangci-lint run ./...` locally.
4. Open a pull request with a description of what changed and why.
5. CI must pass (lint + test + WASM build).
