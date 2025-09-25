package format

import (
	"fmt"
	"io"

	"github.com/Marcel2603/tfcoach/internal/engine"
)

func WriteIssues(issues []engine.Issue, w io.Writer) {
	for _, issue := range issues {
		fmt.Fprintf(w, "%s:%d:%d: %s (%s)\n",
			issue.File, issue.Range.Start.Line, issue.Range.Start.Column, issue.Message, issue.RuleID)
	}
}
