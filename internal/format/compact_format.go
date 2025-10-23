package format

import (
	"fmt"
	"io"
	"slices"

	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/fatih/color"
)

func writeTextIssuesCompact(issues []types.Issue, w io.Writer) {
	preparedIssues := toIssueOutputs(issues)
	slices.SortStableFunc(preparedIssues, func(a, b issueOutput) int {
		return a.Severity.Cmp(b.Severity)
	})
	for _, issue := range preparedIssues {
		_, _ = fmt.Fprintf(
			w,
			"%s %s:%d:%d: %s %s\n",
			color.New(issue.Severity.Color(), color.Bold).Sprint(string(issue.Severity.String()[0])),
			boldFont.Sprint(issue.File),
			issue.Line,
			issue.Column,
			issue.Message,
			greyColor.Sprint("["+issue.RuleID+"]"),
		)
	}
}

func writeTextSummaryCompact(issues []types.Issue, w io.Writer) {
	_, _ = fmt.Fprintf(w, "Summary: %d issue%s\n", len(issues), condPlural(len(issues)))
}
