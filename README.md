<p align="center">
  <img src="assets/logo.png" alt="spec2scenarigo" width="400">
</p>

# spec2scenarigo

A CLI tool that generates [scenarigo](https://github.com/zoncoen/scenarigo) E2E test scenario YAML files from OpenAPI Spec files.

It makes real HTTP requests to your API, captures the responses, and writes a `scenario.yml` ready to use as scenarigo test input.

## Installation

```bash
go install github.com/o-ga09/spec2scenarigo@latest
```

## Usage

```bash
spec2scenarigo [input file] [flags]
```

**Minimal example:**

```bash
export API_KEY=your_api_key
spec2scenarigo example/input.yml
```

**With options:**

```bash
spec2scenarigo example/input.yml \
  -s https://api.example.com \
  -c example/param.csv \
  -o output/scenario.yml \
  -t 200,404
```

### Flags

| Flag | Short | Description |
|---|---|---|
| `--csv-file` | `-c` | CSV file to override path/query parameters |
| `--dry-run` | `-d` | Print parsed spec to stdout without generating a file |
| `--host` | `-s` | API base URL (defaults to `servers[0].URL` in the spec) |
| `--output-file` | `-o` | Output filename (default: `scenario.yml`) |
| `--test-case` | `-t` | Comma-separated HTTP status codes to include (e.g. `200,404`) |

## Authentication

Only `x-api-key` header authentication is supported. Set the key via the `API_KEY` environment variable:

```bash
export API_KEY=your_api_key
```

The generated scenario references it as `{{ env.API_KEY }}`.

## CSV Parameter Override

When the parameter definitions in the OpenAPI Spec are insufficient (e.g. path parameters need concrete values), you can supply them via a CSV file.

- No header row
- Column order: `path(?query_string)`, `HTTP method`, `request body (not yet implemented)`

```csv
/v1/health?userId=2222&companyId=xxxxx,GET,
/v1/user/000000?company=xxxxxx,GET,
```

Paths with parameters like `{id}` in the spec are matched against CSV paths using wildcard segment matching.

## Output Example

Scenario generated from `example/input.yml`:

```yaml
title: MH-API
vars:
  endpoint: http://localhost:8080
steps:
- title: Health check endpoint
  protocol: http
  request:
    method: GET
    url: '{{ vars.endpoint }}/v1/health'
    query:
      userId: "2222"
      companyId: xxxxx
    header:
      x-api-key: '{{ env.API_KEY }}'
  expect:
    code: 200
    body:
      message: ok
```

## Development

```bash
# Build
go build -o spec2scenarigo .

# Run tests
go test ./...

# Run tests with coverage
go test ./... -coverprofile=coverage.out

# Dry run (no file generated)
go run . example/input.yml --dry-run

# Run with CSV parameter override
go run . example/input.yml -c example/param.csv -s https://api.example.com -o output.yml
```

## Roadmap

- [ ] Add test cases
- [ ] Support Basic auth, Bearer token, and AWS Signature V4
- [ ] Support request body via CSV
