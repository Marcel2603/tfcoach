package format

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/Marcel2603/tfcoach/internal/types"
)

type jsonOutput struct {
	IssueCount int
	Issues     []types.Issue
}

func WriteResults(issues []types.Issue, w io.Writer, outputFormat string) error {
	switch outputFormat {
	case "raw":
		writeTextIssues(issues, w)
		writeTextSummary(issues, w)
	case "json":
		err := writeJson(issues, w)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown output format: %s", outputFormat)
	}
	return nil
}

func writeTextIssues(issues []types.Issue, w io.Writer) {
	for _, issue := range issues {
		_, _ = fmt.Fprintf(w, "%s:%d:%d: %s (%s)\n",
			issue.File, issue.Range.Start.Line, issue.Range.Start.Column, issue.Message, issue.RuleID)
	}
}

func writeTextSummary(issues []types.Issue, w io.Writer) {
	_, _ = fmt.Fprintf(w, "Summary:\n Issues: %d\n", len(issues))
}

func writeJson(issues []types.Issue, w io.Writer) error {
	output := jsonOutput{
		IssueCount: len(issues),
		Issues:     issues,
	}
	outputAsStr, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}
	_, err = w.Write(outputAsStr)
	if err != nil {
		return err
	}
	return nil
}
