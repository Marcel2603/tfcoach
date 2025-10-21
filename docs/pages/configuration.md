# Configuration

You can configure tfcoach via file and partially via env.

To configure via file, you need to create one of the following files in the same directory where you execute tfcoach:

- .tfcoach
- .tfcoach.json
- .tfcoach.y[a]ml

Syntax for ".tfcoach" is JSON.

The configurations are applied in the following order, so that each potentially overrides values
from the previous ones:

- Default config [from the repository](https://github.com/Marcel2603/tfcoach/blob/main/cmd/config/.tfcoach.default.yml)
- User config using one of the above-mentioned files
- Environment variables
- Command flags

## What can you configure

Example `.tfcoach.yml`:

```yaml
rules:  # map to restrict rule configurations
  rule_id:  # rule_id of the rule you want to configure
    enabled: false  # decide to enable or disable the rule (enabled by default)
    spec: {}  # config specific configuration, see rule documentation
output:
  format: pretty  # see "--help" for supported output formats
  color: true  # enable or disable color; if set to false, equivalent to the "--no-color" flag
  emojis: true  # enable or disable emojis; if set to false, equivalent to the "--no-emojis" flag
```
