# tfcoach

⚠️ Disclaimer

> This project is currently under heavy development and not yet feature-complete.
> Functionality, APIs, and behavior may change without notice. Use at your own risk — contributions and feedback are
> welcome!

A lightweight linter for Terraform code.  
tfcoach helps teams enforce consistent conventions and avoid common pitfalls by running simple, composable rules against
`.tf` files.

---

## Features

- Command-line interface built with Cobra
- Best-Practice Rules
- Fast: parses each file once, applies multiple rules
- CI-friendly exit codes

---

## Installation

### From source

```bash
go install github.com/Marcel2603/tfcoach@latest
```

This will install `tfcoach` into your `$GOPATH/bin` or `$GOBIN`.

### From release

Download the latest binary from the [Releases](https://github.com/Marcel2603/tfcoach/releases) page and place it in your
`$PATH`.

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

Exit codes:

- `0`: no issues found
- `1`: issues found
- `2`: error running the linter

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
go run ./cmd/tfcoach lint examples
```

---

## Roadmap

- [ ] JSON output (`--format json`) for machine parsing
- [ ] Baseline support to adopt gradually in large codebases
- [ ] Additional rule packs (AWS, GCP, Azure)
- [ ] Auto-fix for selected rules
- [ ] Pluggable rule engine
- [ ] Configurable via `.tfcoach.yml`

---

## License

tfcoach is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.
