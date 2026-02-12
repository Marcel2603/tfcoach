package formatter

import (
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"

	"github.com/fatih/color"
)

func writePretty(issues []issueOutput, allowEmojis bool, w io.Writer) error {
	issuesGroupedByFile := make(map[string][]issueOutput)
	longestFilePath := 10 // for padding
	for _, issue := range issues {
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

	symbols := getPrettyFormatSymbols(allowEmojis)

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
				symbols.ruleMessagePrefix,
				issue.Message,
				symbols.docsPrefix,
				issue.DocsURL,
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type prettyFormatSymbols struct {
	ruleMessagePrefix string
	docsPrefix        string
}

func getPrettyFormatSymbols(allowEmojis bool) prettyFormatSymbols {
	if allowEmojis {
		return prettyFormatSymbols{
			ruleMessagePrefix: "💡  ",
			docsPrefix:        "📑  ",
		}
	}
	return prettyFormatSymbols{
		ruleMessagePrefix: "",
		docsPrefix:        "Docs: ",
	}
}
