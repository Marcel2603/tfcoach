package runner

import (
	"fmt"
	"io"

	"github.com/Marcel2603/tfcoach/internal/engine"
	"github.com/Marcel2603/tfcoach/internal/format"
	"github.com/Marcel2603/tfcoach/internal/types"
)

func Lint(path string, src engine.Source, rules []types.Rule, w io.Writer, outputFormat string) int {
	eng := engine.New(src)
	eng.RegisterMany(rules)
	issues, err := eng.Run(path)
	if err != nil {
		_, _ = fmt.Fprintf(w, "error: %v\n", err)
		return 2
	}

	if len(issues) > 0 {
		writeErr := format.WriteResults(issues, w, outputFormat)
		if writeErr != nil {
			_, _ = fmt.Printf("error writing results: %v\n", writeErr)
			return 2
		}
		return 1
	}
	return 0
}
