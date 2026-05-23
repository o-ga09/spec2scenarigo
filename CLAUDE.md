# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Project Does

`spec2scenarigo` is a CLI tool that generates [scenarigo](https://github.com/zoncoen/scenarigo) E2E test scenario YAML files from OpenAPI Spec files. It works by:
1. Parsing an OpenAPI spec file
2. Making actual HTTP requests to the target API to capture real responses
3. Writing a `scenario.yml` file usable as scenarigo test input

Authentication is currently limited to `x-api-key` header only (set via `export API_KEY=xxx`).

## Commands

```bash
# Build
go build -o spec2scenarigo .

# Run tests
go test ./...

# Run tests with coverage
go test ./... -coverprofile=coverage.out

# Run a single test
go test ./internal/... -run TestGenItem

# Run the CLI (dry run, no file generated)
go run . example/input.yml --dry-run

# Run the CLI with CSV parameter override
go run . example/input.yml -c example/param.csv -s https://api.example.com -o output.yml
```

## Architecture

```
main.go                  → Entry point; calls internal.Execute()
internal/
  root.go                → Cobra CLI definition and command handler
  util.go                → Core logic: GenItem, GenScenario, GetResponse, AddParam
  types.go               → All structs: APISpec, Scenario, step, requestInfo, etc.
  util_test.go           → Table-driven tests (test cases are TODO stubs)
pkg/
  util.go                → Shared helpers: CompPath (path matching), InArray
example/
  input.yml              → Sample OpenAPI spec
  param.csv              → Sample CSV for path/query parameter overrides
```

### Key Data Flow

`GenItem(inputFile, cases)` → `*APISpec`  
Reads an OpenAPI YAML file using `kin-openapi`, extracts paths/methods/params/responses into `APISpec`.

`GenScenario(apiSpec, outputFile, [param])` → `scenario.yml`  
Iterates over `APISpec`, calls `GetResponse()` for each path+method to get real API responses, then marshals a `Scenario` struct to YAML.

`AddParam(csvFile)` → `*map[string]addParam`  
Parses a headerless CSV (columns: `path?query=value`, `METHOD`, `body`) to override path/query params when the OpenAPI spec values are insufficient.

### CSV Parameter Override Logic

When `--csv-file` is provided, `CompPath` in `pkg/util.go` matches OpenAPI spec paths (e.g., `/v1/user/{id}`) against CSV paths (e.g., `/v1/user/000000`) by treating `{…}` segments as wildcards. The matched CSV entry's URL and query params replace the spec defaults.

### Scenarigo Output Format

The generated `scenario.yml` uses `{{ vars.endpoint }}` (set to the base URL) and `{{ env.API_KEY }}` (from environment) as template variables. Each response status code from the spec becomes a separate test step with the actual API response as the expected body.

## Dependencies

- `github.com/getkin/kin-openapi` — OpenAPI 3.0 parsing
- `github.com/spf13/cobra` — CLI framework
- `gopkg.in/yaml.v2` — YAML marshaling for output

## Known Limitations / Planned Work

- Only `x-api-key` authentication is supported
- Tests in `internal/util_test.go` have no test cases yet (all TODO)
- Request body support via CSV is not yet implemented
- Basic auth and Bearer token support are planned
