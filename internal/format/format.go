package format

import (
	"fmt"
	"io"

	"github.com/Marcel2603/tfcoach/internal/types"
)

func WriteIssues(issues []types.Issue, w io.Writer) {
	for _, issue := range issues {
		_, _ = fmt.Fprintf(w, "%s:%d:%d: %s (%s)\n",
			issue.File, issue.Range.Start.Line, issue.Range.Start.Column, issue.Message, issue.RuleID)
	}
}

func WriteSummary(issues []types.Issue, w io.Writer) {
	_, _ = fmt.Fprintf(w, "Summary:\n Issues: %d\n", len(issues))
}
