//go:build tfcoach_tools

package main

import (
	"bytes"
	"cmp"
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/Marcel2603/tfcoach/rules/core"
)

func GenerateRulesOverview(filename string) {
	var buf bytes.Buffer
	buf.WriteString("# Rules\n")
	buf.WriteString("## Core\n")
	buf.WriteString("| Rule | Summary |\n")
	buf.WriteString("|--------|---------|\n")
	rules := core.All()

	slices.SortStableFunc(rules, func(a, b types.Rule) int {
		return cmp.Compare(a.META().Title, b.META().Title)
	})

	for _, r := range core.All() {
		meta := r.META()
		buf.WriteString(fmt.Sprintf("| [%s](%s.md) | %s |\n", meta.Title, meta.DocsURI, meta.Description))
	}
	if err := os.WriteFile(filename, buf.Bytes(), 0644); err != nil {
		log.Fatalf("failed to write rules overview: %v", err)
	}
}
