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

### Homebrew

```shell
brew tap Marcel2603/tap
brew install tfcoach
```

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
    rev: v0.5.0
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

<!-- markdownlint-disable MD013 -->
```shell
Summary: 2 rules broken (3 issues total)

─── Naming Convention (Severity HIGH) ─────────────────────────────

💡  Terraform names should only contain lowercase alphanumeric characters and underscores.

🆔  [core.naming_convention]
📑  https://marcel2603.github.io/tfcoach/rules/core/naming_convention

⚠️  Broken at:
🔹 examples/non_compliant/main.tf:5:1 ➡️  Block "tEst" violates naming convention, it should only contain lowercase alphanumeric characters and underscores.
🔹 examples/non_compliant/main.tf:7:1 ➡️  Block "is-not-compliant" violates naming convention, it should only contain lowercase alphanumeric characters and underscores.


─── Avoid using hashicorp/null provider (Severity MEDIUM) ─────────

💡  With newer Terraform version, use locals and terraform_data as native replacement for hashicorp/null

🆔  [core.avoid_null_provider]
📑  https://marcel2603.github.io/tfcoach/rules/core/avoid_null_provider

⚠️  Broken at:
🔹 examples/non_compliant/data.tf:1:1 ➡️  Use locals instead of null_data_source
```

Alternative output formats (e.g. CI-friendly JSON) available with the `--format` flags (see `tfcoach lint --help`
for all options).

Exit codes:

- `0`: no issues found
- `1`: issues found
- `2`: error running the linter

### Convert a JSON-report into a human-friendly format

To avoid re-running the analysis in your CI-pipeline, run `tfcoach lint` with `--format json` and perform
the analysis needed, then reformat it into any human-friendly format with `tfcoach print`:

```shell
tfcoach lint . --format json > report.json  # run with machine-friendly output
cat report.json | jq '.issue_count'  # do something with report
tfcoach print report.json --format pretty  # output result in human friendly output
```

### Configuration

The behavíour of tfcoach can be configured by environment variables or by one of these files:

- .tfcoach (json syntax)
- .tfcoach.json
- .tfcoach.y[a]ml

More about the configuration [in the docs](https://marcel2603.github.io/tfcoach/configuration/)

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
- [x] Configurable via `.tfcoach.yml` → #15
- [ ] Baseline support to adopt gradually in large codebases → <https://marcel2603.github.io/tfcoach/rule-ideas/>
- [ ] Auto-fix for selected rules
- [ ] Third party ruleset support
- [ ] Additional rule packs (AWS, GCP, Azure)

---

## License

tfcoach is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.
