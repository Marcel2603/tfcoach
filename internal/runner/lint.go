package runner

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/Marcel2603/tfcoach/internal/engine"
	"github.com/Marcel2603/tfcoach/internal/formatter"
	"github.com/Marcel2603/tfcoach/internal/types"
)

func Lint(path string, src engine.Source, rules []types.Rule, w io.Writer, outputFormat string, allowEmojis bool) int {
	eng := engine.New(src)
	eng.RegisterMany(rules)
	issues, err := eng.Run(path)
	if err != nil {
		_, _ = fmt.Fprintf(w, "error: %v\n", err)
		return 2
	}

	if len(issues) > 0 {
		writeErr := formatter.WriteResults(issues, w, outputFormat, allowEmojis)
		if writeErr != nil {
			slog.Error("error writing results", "err", writeErr)
			return 2
		}
		return 1
	}
	return 0
}
