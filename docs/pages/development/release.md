# Release

We use Github Action to release using [semantic-release](https://github.com/semantic-release/semantic-release)
and [goreleaser](https://goreleaser.com/).

## ReleaseFlow

```mermaid
sequenceDiagram
    autonumber
    participant Dev as Developer
    participant DeployWF as release.yml (orchestrator)
    participant ReleaseEnvironment as Environment
    participant ReleaseGate as release-actions.yml (reusable)
    participant SR as semantic-release
    participant GH as GitHub (Tags/Releases, GHCR)
    participant GR as GoReleaser

    Dev->>DeployWF: manual dispatch Release-Gate
    DeployWF->>ReleaseEnvironment: wait for approval
    Dev->>ReleaseEnvironment: Approve Deployment
    ReleaseEnvironment->>DeployWF: Start deployment
    DeployWF->>Release: uses: ./.github/action/action.yml
    Release->>SR: analyze commits (conventional commits)
    SR-->>ReleaseGate: decision (release? yes/no)
    alt New release
        SR->>GH: generate changelog + create tag + release
        ReleaseGate-->>DeployWF: outputs {released:true, tag:vX.Y.Z, version:X.Y.Z}
        DeployWF->>GH: checkout ref = tag
        DeployWF->>GR: build artifacts (binaries)
        DeployWF->>GH: publish release assets
        DeployWF->>GH: buildx + push Docker images to GHCR
    else No release
        ReleaseGate-->>DeployWF: outputs {released:false}
        DeployWF--x GH: skip build & push
    end
```

## Artifacts

We release multiple artifacts:

| Type             | Arc                                      | Link                                                                   |
|------------------|------------------------------------------|------------------------------------------------------------------------|
| Docker Container | arm64 and amd64 for linux                | <https://github.com/Marcel2603/tfcoach/pkgs/container/tfcoach%2Ftfcoach> |
| Executable       | arm64 and amd64 for linux/darwin/windows | <https://github.com/Marcel2603/tfcoach/releases>                         |
