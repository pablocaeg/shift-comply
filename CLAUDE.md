# Shift Comply

Healthcare scheduling regulation engine. Go library with CLI, REST API, WASM build, and Next.js demo website.

## Project structure

- `comply/` - Core Go package: types, registry, query, validation, constraints
- `jurisdictions/` - One Go package per jurisdiction (registered at init time)
- `cmd/shiftcomply/` - CLI tool
- `cmd/shiftcomply-api/` - REST API server
- `cmd/wasm/` - WebAssembly build
- `website/` - Next.js demo (TypeScript, Tailwind, shadcn/ui)
- `.claude/agents/` - Custom agents for adding jurisdictions

## Key rules

- Every regulation must have a real legal citation (Source.Title and Source.Section are never empty)
- Every rule must have the correct Scope (public_health, hospitals, accredited_programs, etc.)
- No em dashes or double dashes in Go strings
- Use `comply.D(year, time.Month, day)` for dates
- ACGME is not law: use ScopeAccreditedPrograms
- Run `go test ./...` and `golangci-lint run ./...` before committing
- Website WASM files (shiftcomply.wasm, wasm_exec.js) are build artifacts, not committed

## Adding a jurisdiction

Use the new-jurisdiction agent: `/agents/new-jurisdiction Add France`

Or manually: create a package under `jurisdictions/`, register via `init()`, add blank import in `jurisdictions/jurisdictions.go`, add tests.

## Testing

```sh
go test ./...              # all Go tests
go test -race ./...        # with race detector
cd website && npm run build  # website build check
```
