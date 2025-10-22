# Usage 
## tfcoach

Tiny Terraform coach

```
tfcoach [flags]
```

### Options

```
  -h, --help   help for tfcoach
```



### Exit Codes

| Code | Meaning|
|------|--------|
| 0 | OK |

## tfcoach lint

Lint Terraform files

```
tfcoach lint [path] [flags]
```

### Options

```
  -f, --format string   Output format. Supported: json|compact|pretty|educational (default "pretty")
  -h, --help            help for lint
      --no-color        Disable color output
      --no-emojis       Prevent emojis in output
```



### Exit Codes

| Code | Meaning|
|------|--------|
| 0 | No issues found |
| 1 | Issues found |
| 2 | Runtime error |

## tfcoach version

Print the version number

```
tfcoach version [flags]
```

### Options

```
  -h, --help   help for version
```



### Exit Codes

| Code | Meaning|
|------|--------|
| 0 | OK |

