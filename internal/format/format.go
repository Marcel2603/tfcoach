//revive:disable:var-naming For now it's okay to have a generic name
package format

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/Marcel2603/tfcoach/internal/constants"
	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/Marcel2603/tfcoach/rules/core"
	"github.com/fatih/color"
)

const ruleDocsFormat = "https://marcel2603.github.io/tfcoach/rules/%s"

var (
	boldFont  = color.New(color.Bold)
	greyColor = color.RGB(90, 90, 90)
)

type issueOutput struct {
	File     string         `json:"file"`
	Line     int            `json:"line"`
	Column   int            `json:"column"`
	Message  string         `json:"message"`
	RuleID   string         `json:"rule_id"`
	Severity types.Severity `json:"severity"`
	Category string         `json:"category"`
	DocsURL  string         `json:"docs_url"`
}

type jsonOutput struct {
	IssueCount int           `json:"issue_count"`
	Issues     []issueOutput `json:"issues"`
}

func WriteResults(issues []types.Issue, w io.Writer, outputFormat string, allowEmojis bool) error {
	preparedIssues := toIssueOutputs(issues)
	return writeResultsFromPrepared(preparedIssues, w, outputFormat, allowEmojis)
}

func ReformatResults(srcReport []byte, w io.Writer, outputFormat string, allowEmojis bool) error {
	var report jsonOutput
	err := json.Unmarshal(srcReport, &report)
	if err != nil {
		return err
	}
	return writeResultsFromPrepared(report.Issues, w, outputFormat, allowEmojis)
}

func writeResultsFromPrepared(issues []issueOutput, w io.Writer, outputFormat string, allowEmojis bool) error {
	switch outputFormat {
	case "compact":
		writeTextIssuesCompact(issues, w)
		writeTextSummaryCompact(issues, w)
	case "json":
		err := writeJSON(issues, w)
		if err != nil {
			return err
		}
	case "pretty":
		err := writePretty(issues, allowEmojis, w)
		if err != nil {
			return err
		}
	case "educational":
		err := writeEducational(issues, allowEmojis, w)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown output format: %s", outputFormat)
	}
	return nil
}

func toIssueOutputs(issues []types.Issue) []issueOutput {
	var result []issueOutput

	for _, issue := range issues {
		rule, err := core.FindByID(issue.RuleID)
		var severity types.Severity
		var docsURL string
		if err != nil {
			severity = constants.SeverityUnknown
			docsURL = "about:blank"
		} else {
			rulesMeta := rule.META()
			severity = rulesMeta.Severity
			docsURL = fmt.Sprintf(ruleDocsFormat, rulesMeta.DocsURI)
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

func condPlural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
