//go:build exclude

package main

import (
	"fmt"

	"github.com/Marcel2603/tfcoach/rules/core"
)

func main() {
	fmt.Println("# Rules")
	fmt.Println("## Core")
	fmt.Println("| Rule ID | Summary |")
	fmt.Println("|--------|---------|")
	for _, r := range core.All() {
		fmt.Printf("| %s |  |\\n", r.ID())
	}
}
