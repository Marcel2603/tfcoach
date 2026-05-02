# Generate Doc Pages

This tool will generate the following pages from the code:

- [Rules Overview](https://marcel2603.github.io/tfcoach/rules/) (uses information from internal api)
- [Usage](https://marcel2603.github.io/tfcoach/usage/) (uses [cobra doc](https://pkg.go.dev/github.com/spf13/cobra/doc)
  package)
- generates `nav`-section inside [zensical-config](https://github.com/Marcel2603/tfcoach/blob/main/docs/zensical.toml)
- Plain rule-site if it does not exist

## How to execute

```shell
make generate-documentation
```
