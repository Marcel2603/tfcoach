package constants

import "github.com/Marcel2603/tfcoach/internal/types"

var (
	SeverityHigh    = types.Severity{Str: "HIGH", Priority: 1}
	SeverityMedium  = types.Severity{Str: "MEDIUM", Priority: 2}
	SeverityLow     = types.Severity{Str: "LOW", Priority: 3}
	SeverityUnknown = types.Severity{Str: "UNKNOWN", Priority: 99}
)
