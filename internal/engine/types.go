package engine

import "github.com/hashicorp/hcl/v2"

type Issue struct {
	File     string
	Range    hcl.Range
	Message  string
	RuleID   string
	Severity string
}

type Rule interface {
	ID() string
	Apply(file string, f *hcl.File) []Issue
}
