# tfcoach

A lightweight linter for Terraform code.
tfcoach helps teams enforce consistent conventions and avoid common pitfalls by running simple, composable rules against
`.tf` files.

---

## Note

This project is currently under heavy development and not yet feature-complete.
Functionality, APIs, and behavior may change without notice. Use at your own risk — contributions and feedback are
welcome!

---

[![Go Report Card](https://goreportcard.com/badge/github.com/Marcel2603/tfcoach)](https://goreportcard.com/report/github.com/Marcel2603/tfcoach)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/Marcel2603/tfcoach/blob/main/LICENSE)
[![GitHub Release](https://img.shields.io/github/v/release/Marcel2603/tfcoach)](https://github.com/Marcel2603/tfcoach/releases/latest)

---

## Features

- Command-line interface built with Cobra
- Best-Practice Rules
- Fast: parses each file once, applies multiple rules in parallel
- CI-friendly output and exit codes

---

## Installation

### From release

Download the latest binary from the [Releases](https://github.com/Marcel2603/tfcoach/releases) page and place it in your
`$PATH`.

### From docker

```bash
docker pull ghcr.io/marcel2603/tfcoach/tfcoach:latest
docker run -v $PWD:/app ghcr.io/marcel2603/tfcoach/tfcoach:latest lint .
```

### From source

```bash
go install github.com/Marcel2603/tfcoach@latest
```

This will install `tfcoach` into your `$GOPATH/bin` or `$GOBIN`.

### Pre-commit

```yaml
  - repo: https://github.com/Marcel2603/tfcoach
    rev: v0.4.0
    hooks:
      - id: tfcoach # executes tfcoach via golang
      - id: tfcoach-docker # executes tfcoach via docker
```

---

## Usage

### Lint a directory

```bash
tfcoach lint .
```

### Example output

```shell
main.tf:12:3: resource name must be "this" (core.naming.require_this)
versions.tf:1:1: terraform block must declare "required_version" (core.terraform.require_version)
```

Alternative CI-friendly output format available with `--format json` (see `tfcoach lint --help` for all options).

Exit codes:

- `0`: no issues found
- `1`: issues found
- `2`: error running the linter

### Configuration

The behavoíour of tfcoach can be configured by environment variables or by one of these files:

- .tfcoach (json syntax)
- .tfcoach.json
- .tfcoach.y[a]ml

More about the configuration [here](link to mkdocs)

---

## Development

Clone and run tests with coverage:

```bash
git clone https://github.com/Marcel2603/tfcoach.git
cd tfcoach
make test
make cover-html
```

Run the CLI locally:

```bash
go run main.go lint examples
```

---

## Roadmap

- [x] Alternative output formats (See option `--format`) → #13
- [ ] Baseline support to adopt gradually in large codebases
- [ ] Additional rule packs (AWS, GCP, Azure)
- [ ] Auto-fix for selected rules
- [ ] Pluggable rule engine
- [x] Configurable via `.tfcoach.yml` → #15

---

## License

tfcoach is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.
