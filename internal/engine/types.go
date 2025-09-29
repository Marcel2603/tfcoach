package engine

import "github.com/hashicorp/hcl/v2"

type Issue struct {
	File    string
	Range   hcl.Range
	Message string
	RuleID  string
}

type RuleMeta struct {
	Title       string
	Description string
	Severity    string
	DocsURL     string
}

type Rule interface {
	ID() string
	META() RuleMeta
	Apply(file string, f *hcl.File) []Issue
}
