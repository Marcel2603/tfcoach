package format

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/Marcel2603/tfcoach/rules/core"
	"github.com/hashicorp/hcl/v2"
)

const ruleDocsFormat = "https://github.com/Marcel2603/tfcoach/tree/main/docs/pages/rules/%s.md"

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

type placeholderRule struct{}

func (r *placeholderRule) ID() string {
	return ""
}

func (r *placeholderRule) META() types.RuleMeta {
	return types.RuleMeta{
		Title:       "",
		Description: "",
		Severity:    "UNKNOWN",
		DocsURL:     "about:blank",
	}
}

func (r *placeholderRule) Apply(_ string, _ *hcl.File) []types.Issue {
	return []types.Issue{}
}

func (r *placeholderRule) Finish() []types.Issue {
	return []types.Issue{}
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
		if err != nil {
			rule = &placeholderRule{}
		}
		ruleMeta := rule.META()

		result = append(result, issueOutput{
			File:     issue.File,
			Line:     issue.Range.Start.Line,
			Column:   issue.Range.Start.Column,
			Message:  issue.Message,
			RuleId:   issue.RuleID,
			Severity: ruleMeta.Severity,
			//Category: "?",  // TODO later: implement rule category
			DocsUrl: fmt.Sprintf(ruleDocsFormat, ruleMeta.DocsURL),
		})
	}
	return result
}
