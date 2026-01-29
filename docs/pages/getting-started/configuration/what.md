# What to configure

## Disable or fine-tune rules

Rules can be disabled at will with `rules.<rule_id>.enabled: false`.

Some rules may allow for further configuration using the `rules.<rule_id>.spec` map, please refer to the rule-specific
documentation for more information.

## Output format

Several output formats are supported under `output.format`:

- `educational` (default): issues grouped by rule in a colorful human-friendly format including explanation
- `pretty`: issues grouped by file in a colorful human-friendly format
- `compact`: one line per issue with location information for quick scanning, sorted by severity
- `json`: detailed output optimized for machine-parsing

Color can be disabled with `output.color: false`.

Emojis can be disabled with `output.emojis: false` (they get replaced with equivalent text).

Issues found in `.terragrunt-cache` directories are usually not wanted so disabled by default. These directories can be
scanned by setting the property `output.include_terragrunt_cache: true`.

## Example

Example `.tfcoach.yml` (same options available with the JSON format):

<!-- markdownlint-disable MD013 -->
```yaml
rules: # map to restrict rule configurations
  core.example_rule: # rule_id of the rule you want to configure
    enabled: false  # decide to enable or disable the rule (enabled by default)
    spec: { }  # config specific configuration, see rule documentation
output:
  format: pretty  # see "--help" for supported output formats
  color: true  # enable or disable color; if set to false, equivalent to the "--no-color" flag
  emojis: true  # enable or disable emojis; if set to false, equivalent to the "--no-emojis" flag
  include_terragrunt_cache: false  # enable or disable terragrunt-cache scanning; if set to true, equivalent to the "--include-terragrunt-cache" flag
```
