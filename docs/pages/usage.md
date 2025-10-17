# Usage

```bash
tfcoach lint [path]
```

Available output formats using the `--format` flag:

- `pretty` (default): colorful output with explanation in a human-friendly format and link to docs
- `compact`: one line per issue with location information and rule explanation in a human-friendly format
- `json`: more detailed output optimized for machine-parsing

Additional options:

- `--no-color`: Disable color in output

Exit codes:

- 0: no issues
- 1: issues found
- 2: runtime error
