//go:build tfcoach_tools

package main

import "fmt"

func main() {
	fmt.Println("Generate Usage-page")
	GenerateUsage("docs/pages/getting-started/usage.md")
	fmt.Println("Usage-page generated")
	fmt.Println("Generate Rules Overview")
	GenerateRulesOverview("docs/pages/rules/index.md")
	fmt.Println("Rules Overview generated")
}
