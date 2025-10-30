package format

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/Marcel2603/tfcoach/rules/core"
	"github.com/fatih/color"
)

func writeEducational(issues []types.Issue, allowEmojis bool, w io.Writer) error {
	issuesGroupedByRuleID := groupByRuleID(issues)

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

	brokenRules := extractRulesSortedBySeverity(issuesGroupedByRuleID, err)

	longestRuleTitle := 10 // min for padding
	for _, rule := range brokenRules {
		ruleMeta := rule.META()
		longestRuleTitle = max(longestRuleTitle, len(ruleMeta.Title)+len(ruleMeta.Severity.String()))
	}
	symbols := getEducationalFormatSymbols(allowEmojis)

	for _, rule := range brokenRules {
		ruleID := rule.ID()
		ruleMeta := rule.META()
		issuesForRule := issuesGroupedByRuleID[ruleID]
		slices.SortStableFunc(issuesForRule, func(a, b types.Issue) int {
			return strings.Compare(a.File, b.File)
		})

		padding := strings.Repeat("─", longestRuleTitle-len(ruleMeta.Title)-len(ruleMeta.Severity.String()))

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
			symbols.explanationPrefix,
			ruleMeta.Description,
			symbols.idPrefix,
			greyColor.Sprint("["+ruleID+"]"),
			symbols.docsPrefix,
			docsURL,
			symbols.brokenListPrefix,
		)
		if err != nil {
			return err
		}
		for _, issue := range issuesForRule {
			_, err = fmt.Fprintf(
				w,
				"%s%s:%d:%d%s%s\n",
				symbols.ruleMessagePrefix,
				issue.File,
				issue.Range.Start.Line,
				issue.Range.Start.Column,
				symbols.ruleMessageInfix,
				issue.Message,
			)
			if err != nil {
				return err
			}
		}
		_, err = fmt.Fprintln(w)
		if err != nil {
			return err
		}
	}
	return nil
}

func extractRulesSortedBySeverity(issuesGroupedByRuleID map[string][]types.Issue, err error) []types.Rule {
	var rules []types.Rule
	for ruleID := range issuesGroupedByRuleID {
		var rule types.Rule
		rule, err = core.FindByID(ruleID)
		if err != nil {
			rule = &core.UnknownRule{PseudoID: ruleID}
		}
		rules = append(rules, rule)
	}

	slices.SortStableFunc(rules, func(a, b types.Rule) int {
		aMeta := a.META()
		bMeta := b.META()
		if aMeta.Severity != bMeta.Severity {
			return aMeta.Severity.Cmp(bMeta.Severity)
		}
		return strings.Compare(aMeta.Title, bMeta.Title)
	})
	return rules
}

func groupByRuleID(issues []types.Issue) map[string][]types.Issue {
	issuesGroupedByRuleID := make(map[string][]types.Issue)
	for _, issue := range issues {
		issuesGroupedByRuleID[issue.RuleID] = append(issuesGroupedByRuleID[issue.RuleID], issue)
	}
	return issuesGroupedByRuleID
}

type educationalFormatSymbols struct {
	explanationPrefix string
	idPrefix          string
	docsPrefix        string
	brokenListPrefix  string
	ruleMessagePrefix string
	ruleMessageInfix  string
}

func getEducationalFormatSymbols(allowEmojis bool) educationalFormatSymbols {
	if allowEmojis {
		return educationalFormatSymbols{
			explanationPrefix: "💡  ",
			idPrefix:          "🆔  ",
			docsPrefix:        "📑  ",
			brokenListPrefix:  "⚠️  ",
			ruleMessagePrefix: "🔹 ",
			ruleMessageInfix:  " ➡️  ",
		}
	}
	return educationalFormatSymbols{
		explanationPrefix: "Explanation: ",
		idPrefix:          "ID: ",
		docsPrefix:        "Read more: ",
		brokenListPrefix:  "",
		ruleMessagePrefix: "- ",
		ruleMessageInfix:  " ─ ",
	}
}
