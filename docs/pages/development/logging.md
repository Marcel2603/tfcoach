# Logging

tfcoach uses Go's standard [`log/slog`](https://pkg.go.dev/log/slog) package with a globally configured default logger.

## Enabling logs

Set the `TF_COACH_LOG` environment variable before running tfcoach:

```shell
TF_COACH_LOG=DEBUG tfcoach lint .
```

| Value   | Effect                                          |
|---------|-------------------------------------------------|
| `DEBUG` | All messages (debug, info, warn, error)         |
| `INFO`  | Info, warn, and error messages                  |
| `WARN`  | Warn and error messages only                    |
| `ERROR` | Error messages only                             |
| `JSON`  | JSON output to stderr, INFO level               |

If `TF_COACH_LOG` is unset or has an unrecognised value, the level defaults to `ERROR`.

## Writing log statements

Import `log/slog` directly — no need to pass a logger around:

```go
import "log/slog"

slog.Debug("parsing file", "path", filePath)
slog.Info("rule applied", "rule", rule.ID(), "issues", len(issues))
slog.Warn("config not found, using defaults")
slog.Error("failed to read file", "err", err)
```

Every log entry automatically includes the timestamp and the source file + line number where the call was made.

## Output

Logs are written to **stderr**, keeping stdout clean for lint reports and JSON output.
