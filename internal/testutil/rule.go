//go:build test

package testutil

import (
	"strings"

	"github.com/Marcel2603/tfcoach/internal/engine"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type AlwaysFlag struct {
	RuleID  string
	Message string
	Match   string
}

func (r AlwaysFlag) ID() string { return r.RuleID }

func (r AlwaysFlag) META() engine.RuleMeta {
	return engine.RuleMeta{
		Title:       "AlwaysFlag",
		Description: r.Message,
		Severity:    "HIGH",
		DocsURL:     "tbd",
	}
}

func (r AlwaysFlag) Apply(file string, f *hcl.File) []engine.Issue {
	body, _ := f.Body.(*hclsyntax.Body)
	if body == nil {
		return nil
	}

	src := ""
	for _, b := range body.Blocks {
		src += b.DefRange().String()
	}
	if r.Match == "" || strings.Contains(src, r.Match) {
		return []engine.Issue{{
			File:    file,
			Range:   body.Range(),
			Message: r.Message,
			RuleID:  r.RuleID,
		}}
	}
	return nil
}
