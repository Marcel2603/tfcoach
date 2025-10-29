# Where to configure

You can configure tfcoach via file and partially via env.

## Local configuration

To configure via file, you need to create one of the following files in the same directory where you execute tfcoach:

- `.tfcoach.yml`
- `.tfcoach.yaml`
- `.tfcoach.json`
- `.tfcoach` (the syntax for `.tfcoach` is JSON)

It is possible to tell tfcoach to look for configuration in another location with the `--config` flag:

```shell
tfcoach lint --config ~/.my-config-files . # if the argument is a directory, tfcoach will look for the above-mentioned names
tfcoach lint --config ../my-tfcoach-config.yml . # if the argument is a file, it will be used as-is
```

## Global configuration

In order to avoid repeating the same config in multiple directories, you can define a global config in one of the two
following directories:

- `$HOME/.config/tfcoach/`
- `$HOME/.tfcoach/`

where `$HOME` is the home-directory of your operating system.

## Loading order

The configurations are applied in the following order, so that each potentially overrides values
from the previous ones:

- Default config [from the repository](https://github.com/Marcel2603/tfcoach/blob/main/cmd/config/.tfcoach.default.yml)
- Global config if it exists
- Local config (from the current directory or the alternative location provided via `--config`)
- Environment variables
- Command flags (see `--help`)
