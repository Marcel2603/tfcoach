# Getting started

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

## First run

```bash
tfcoach lint .
```
