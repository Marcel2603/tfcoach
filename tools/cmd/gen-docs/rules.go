//go:build tfcoach_tools

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/Marcel2603/tfcoach/rules/core"
)

func GenerateRulesOverview(filename string) {
	var buf bytes.Buffer
	buf.WriteString("# Rules \n")
	buf.WriteString("## Core \n")
	buf.WriteString("| Rule | Summary | \n")
	buf.WriteString("|--------|---------| \n")
	for _, r := range core.All() {
		meta := r.META()
		buf.WriteString(fmt.Sprintf("| [%s](%s.md) | %s |\n", meta.Title, meta.DocsURL, meta.Description))
	}
	if err := os.WriteFile(filename, buf.Bytes(), 0644); err != nil {
		log.Fatalf("failed to write rules overview: %v", err)
	}
}
