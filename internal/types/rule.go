package types

import "github.com/hashicorp/hcl/v2"

type RuleMeta struct {
	Title       string
	Description string
	Severity    Severity
	DocsURL     string
}

type Rule interface {
	ID() string
	META() RuleMeta
	Apply(file string, f *hcl.File) []Issue
	Finish() []Issue
}
