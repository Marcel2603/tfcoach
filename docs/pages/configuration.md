# Configuration

You can configure tfcoach via file and partially via env

To configure via file, you need to create one of the following files in the same directory where you execute tfcoach:

- .tfcoach
- .tfcoach.json
- .tfcoach.y[a]ml

Syntax for ".tfcoach" is JSON.

## What can you configure

The default values can be found [in the repository](https://github.com/marcel2603/tfcoach/cmd/config/.tfcoach.default.yml).

```yaml
rules: # map to restrict rule configurations
  rule_id: # rule_id of the rule you want to configure
    enabled: false # decide to enable or disable the rule (enabled by default)
    spec: {} # config specific configuration, see rule documentation 
```
