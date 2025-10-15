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
	RuleId   string `json:"rule_id"`
	Severity string `json:"severity"`
	Category string `json:"category"`
	DocsUrl  string `json:"docs_url"`
}

type jsonOutput struct {
	IssueCount int           `json:"issue_count"`
	Issues     []issueOutput `json:"issues"`
}

func WriteResults(issues []types.Issue, w io.Writer, outputFormat string) error {
	switch outputFormat {
	case "compact":
		writeTextIssuesCompact(issues, w)
		writeTextSummaryCompact(issues, w)
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

func writeTextIssuesCompact(issues []types.Issue, w io.Writer) {
	for _, issue := range issues {
		_, _ = fmt.Fprintf(w, "%s:%d:%d: %s (%s)\n",
			issue.File, issue.Range.Start.Line, issue.Range.Start.Column, issue.Message, issue.RuleID)
	}
}

func writeTextSummaryCompact(issues []types.Issue, w io.Writer) {
	var suffix string
	if len(issues) == 1 {
		suffix = ""
	} else {
		suffix = "s"
	}
	_, _ = fmt.Fprintf(w, "Summary: %d issue%s\n", len(issues), suffix)
}

func writeJson(issues []types.Issue, w io.Writer) error {
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
		rule, err := core.FindById(issue.RuleID)
		var severity, docsUrl string
		if err != nil {
			severity = "UNKNOWN"
			docsUrl = "about:blank"
		} else {
			rulesMeta := rule.META()
			severity = rulesMeta.Severity
			docsUrl = fmt.Sprintf(ruleDocsFormat, rulesMeta.DocsURL)
		}

		result = append(result, issueOutput{
			File:     issue.File,
			Line:     issue.Range.Start.Line,
			Column:   issue.Range.Start.Column,
			Message:  issue.Message,
			RuleId:   issue.RuleID,
			Severity: severity,
			//Category: "?",  // TODO later: implement rule category
			DocsUrl: docsUrl,
		})
	}

	return result
}
