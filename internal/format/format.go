package format

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"

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

func writePretty(issues []types.Issue, allowEmojis bool, w io.Writer) error {
	preparedIssues := toIssueOutputs(issues)
	issuesGroupedByFile := make(map[string][]issueOutput)
	longestFilePath := 10 // for padding
	for _, issue := range preparedIssues {
		issuesGroupedByFile[issue.File] = append(issuesGroupedByFile[issue.File], issue)
		longestFilePath = max(longestFilePath, len(issue.File))
	}

	_, err := fmt.Fprintf(
		w,
		"Summary: %s issue%s found in %s file%s\n\n",
		boldFont.Sprint(len(issues)),
		condPlural(len(issues)),
		boldFont.Sprint(len(issuesGroupedByFile)),
		condPlural(len(issuesGroupedByFile)),
	)
	if err != nil {
		return err
	}

	var ruleMessagePrefix, docsPrefix string
	if allowEmojis {
		ruleMessagePrefix = "💡  "
		docsPrefix = "📑  "
	} else {
		ruleMessagePrefix = ""
		docsPrefix = "Docs: "
	}
	for _, fileName := range slices.Sorted(maps.Keys(issuesGroupedByFile)) {
		issuesInFile := issuesGroupedByFile[fileName]
		slices.SortStableFunc(issuesInFile, func(a, b issueOutput) int {
			return a.Severity.Cmp(b.Severity)
		})

		padding := strings.Repeat("─", longestFilePath-len(fileName))
		_, err = fmt.Fprintf(w, "─── %s %s─────────\n\n", boldFont.Sprint(fileName), padding)
		if err != nil {
			return err
		}
		for _, issue := range issuesInFile {
			_, err = fmt.Fprintf(
				w,
				"  %d:%d\t%s\t%s\n\t%s%s\n\t%s%s\n\n",
				issue.Line,
				issue.Column,
				boldFont.Sprint("["+issue.RuleID+"]"),
				color.New(issue.Severity.Color(), color.Bold).Sprint(issue.Severity),
				ruleMessagePrefix,
				issue.Message,
				docsPrefix,
				issue.DocsURL,
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func writeEducational(issues []types.Issue, allowEmojis bool, w io.Writer) error {
	preparedIssues := toIssueOutputs(issues)
	issuesGroupedByRuleID := make(map[string][]issueOutput)
	longestRuleID := 10 // for padding
	for _, issue := range preparedIssues {
		issuesGroupedByRuleID[issue.RuleID] = append(issuesGroupedByRuleID[issue.RuleID], issue)
		longestRuleID = max(longestRuleID, len(issue.RuleID))
	}

	_, err := fmt.Fprintf(
		w,
		"Summary: %s rule%s broken (%s issue%s total)\n",
		boldFont.Sprint(len(issuesGroupedByRuleID)),
		condPlural(len(issuesGroupedByRuleID)),
		boldFont.Sprint(len(issues)),
		condPlural(len(issues)),
	)
	if err != nil {
		return err
	}

	var brokenRules []types.Rule
	for ruleID := range issuesGroupedByRuleID {
		var rule types.Rule
		rule, err = core.FindByID(ruleID)
		if err != nil {
			rule = core.UnknownRule{PseudoID: ruleID}
		}
		brokenRules = append(brokenRules, rule)
	}
	slices.SortStableFunc(brokenRules, func(a, b types.Rule) int {
		return a.META().Severity.Cmp(b.META().Severity)
	})

	for _, rule := range brokenRules {
		ruleID := rule.ID()
		ruleMeta := rule.META()
		issuesForRule := issuesGroupedByRuleID[ruleID]
		slices.SortStableFunc(issuesForRule, func(a, b issueOutput) int {
			return strings.Compare(a.File, b.File)
		})

		padding := strings.Repeat("─", longestRuleID-len(ruleID))
		var explanationPrefix, docsPrefix, idPrefix, brokenListPrefix, ruleMessagePrefix, ruleMessageInfix string
		if allowEmojis {
			explanationPrefix = "💡  "
			idPrefix = "🆔  "
			docsPrefix = "📑  "
			brokenListPrefix = "⚠️  "
			ruleMessagePrefix = "🔹 "
			ruleMessageInfix = " ➡️  "
		} else {
			explanationPrefix = "Explanation: "
			idPrefix = "ID: "
			docsPrefix = "Read more: "
			brokenListPrefix = ""
			ruleMessagePrefix = "- "
			ruleMessageInfix = " ─ "
		}

		var docsURL string
		if ruleMeta.DocsURI == "about:blank" {
			docsURL = "about:blank"
		} else {
			docsURL = fmt.Sprintf(ruleDocsFormat, ruleMeta.DocsURI)
		}

		_, err = fmt.Fprintf(
			w,
			"\n─── %s (Severity %s) %s─────────\n\n%s%s\n\n%s%s\n%s%s\n\n%sBroken at:\n",
			boldFont.Sprint(ruleMeta.Title),
			color.New(ruleMeta.Severity.Color(), color.Bold).Sprint(ruleMeta.Severity),
			padding,
			explanationPrefix,
			ruleMeta.Description,
			idPrefix,
			greyColor.Sprint("["+ruleID+"]"),
			docsPrefix,
			docsURL,
			brokenListPrefix,
		)
		if err != nil {
			return err
		}
		for _, issue := range issuesForRule {
			_, err = fmt.Fprintf(
				w,
				"%s%s:%d:%d%s%s\n",
				ruleMessagePrefix,
				issue.File,
				issue.Line,
				issue.Column,
				ruleMessageInfix,
				issue.Message,
			)
			if err != nil {
				return err
			}
		}
		_, err = fmt.Fprint(w, "\n")
		if err != nil {
			return err
		}
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
