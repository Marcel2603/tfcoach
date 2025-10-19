package format

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/Marcel2603/tfcoach/rules/core"
)

const ruleDocsFormat = "https://marcel2603.github.io/tfcoach/rules/%s"

type issueOutput struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Message  string `json:"message"`
	RuleID   string `json:"rule_id"`
	Severity string `json:"severity"`
	Category string `json:"category"`
	DocsURL  string `json:"docs_url"`
}

type jsonOutput struct {
	IssueCount int           `json:"issue_count"`
	Issues     []issueOutput `json:"issues"`
}

func WriteResults(issues []types.Issue, w io.Writer, outputFormat string) error {
	switch outputFormat {
	case "raw":
		writeTextIssues(issues, w)
		writeTextSummary(issues, w)
	case "json":
		err := writeJSON(issues, w)
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

func writeJSON(issues []types.Issue, w io.Writer) error {
	output := jsonOutput{
		IssueCount: len(issues),
		Issues:     toIssueOutputs(issues),
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

func toIssueOutputs(issues []types.Issue) []issueOutput {
	var result []issueOutput

	for _, issue := range issues {
		rule, err := core.FindByID(issue.RuleID)
		var severity, docsURL string
		if err != nil {
			severity = "UNKNOWN"
			docsURL = "about:blank"
		} else {
			rulesMeta := rule.META()
			severity = rulesMeta.Severity
			docsURL = fmt.Sprintf(ruleDocsFormat, rulesMeta.DocsURL)
		}

		result = append(result, issueOutput{
			File:     issue.File,
			Line:     issue.Range.Start.Line,
			Column:   issue.Range.Start.Column,
			Message:  issue.Message,
			RuleID:   issue.RuleID,
			Severity: severity,
			//Category: "?",  // TODO later: implement rule category
			DocsURL: docsURL,
		})
	}

	return result
}
