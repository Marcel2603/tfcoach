package runner

import (
	"fmt"
	"io"

	"github.com/Marcel2603/tfcoach/internal/engine"
	"github.com/Marcel2603/tfcoach/internal/format"
)

func Lint(path string, src engine.Source, rules []engine.Rule, w io.Writer) int {
	eng := engine.New(src)
	eng.RegisterMany(rules)
	issues, err := eng.Run(path)
	if err != nil {
		_, _ = fmt.Fprintf(w, "error: %v\n", err)
		return 1
	}

	if len(issues) > 0 {
		format.WriteIssues(issues, w)
		format.WriteSummary(issues, w)
		return 2
	}
	return 0
}
