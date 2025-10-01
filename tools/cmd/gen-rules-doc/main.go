package main

import (
	"fmt"

	"github.com/Marcel2603/tfcoach/rules/core"
)

func main() {
	fmt.Println("# Rules")
	fmt.Println("## Core")
	fmt.Println("| ID | Rule | Summary |")
	fmt.Println("|------|--------|---------|")
	for _, r := range core.All() {
		meta := r.META()
		fmt.Printf("| [%s](%s) | %s | %s |\n", r.ID(), meta.DocsURL, meta.Title, meta.Description)
	}
}
