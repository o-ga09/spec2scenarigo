<p align="center">
  <img src="assets/logo.png" alt="spec2scenario" width="400">
</p>

# spec2scenario

[![Test](https://github.com/o-ga09/spec2scenario/actions/workflows/test.yaml/badge.svg)](https://github.com/o-ga09/spec2scenario/actions/workflows/test.yaml)
[![Lint](https://github.com/o-ga09/spec2scenario/actions/workflows/lint.yaml/badge.svg)](https://github.com/o-ga09/spec2scenario/actions/workflows/lint.yaml)
[![Security](https://github.com/o-ga09/spec2scenario/actions/workflows/security.yml/badge.svg)](https://github.com/o-ga09/spec2scenario/actions/workflows/security.yml)
[![Go version](https://img.shields.io/github/go-mod/go-version/o-ga09/spec2scenario)](https://github.com/o-ga09/spec2scenario)
[![Latest release](https://img.shields.io/github/v/release/o-ga09/spec2scenario)](https://github.com/o-ga09/spec2scenario/releases)
[![License](https://img.shields.io/github/license/o-ga09/spec2scenario)](https://github.com/o-ga09/spec2scenario/blob/main/LICENSE)

A CLI tool that generates [scenarigo](https://github.com/zoncoen/scenarigo) E2E test scenario YAML files from OpenAPI Spec files.

It makes real HTTP requests to your API, captures the responses, and writes a `scenario.yml` ready to use as scenarigo test input.

## Installation

```bash
go install github.com/o-ga09/spec2scenario@latest
```

## Usage

```bash
spec2scenario [input file] [flags]
```

**Minimal example:**

```bash
export API_KEY=your_api_key
spec2scenario example/input.yml
```

**With options:**

```bash
spec2scenario example/input.yml \
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
go build -o spec2scenario .

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
