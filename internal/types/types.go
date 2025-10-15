package types

import (
	"cmp"
	"encoding/json"

	"github.com/hashicorp/hcl/v2"
)

type Issue struct {
	File    string
	Range   hcl.Range
	Message string
	RuleID  string
}

// TODO #13: move severity stuff to separate file?
type Severity struct {
	str      string
	priority int
}

var (
	SeverityHigh    = Severity{"HIGH", 1}
	SeverityMedium  = Severity{"MEDIUM", 2}
	SeverityLow     = Severity{"LOW", 3}
	SeverityUnknown = Severity{"UNKNOWN", 99}
)

func (s Severity) Cmp(other Severity) int {
	return cmp.Compare(s.priority, other.priority)
}

func (s Severity) String() string {
	return s.str
}

func (s Severity) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.str)
}

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
