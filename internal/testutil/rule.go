//go:build test

package testutil

import (
	"strings"

	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type AlwaysFlag struct {
	RuleID  string
	Message string
	Match   string
}

func (r AlwaysFlag) ID() string { return r.RuleID }

func (r AlwaysFlag) META() types.RuleMeta {
	return types.RuleMeta{
		Title:       "AlwaysFlag",
		Description: r.Message,
		Severity:    "HIGH",
		DocsURL:     "tbd",
	}
}

func (r AlwaysFlag) Apply(file string, f *hcl.File) []types.Issue {
	body, _ := f.Body.(*hclsyntax.Body)
	if body == nil {
		return nil
	}

	src := ""
	for _, b := range body.Blocks {
		src += b.DefRange().String()
	}
	if r.Match == "" || strings.Contains(src, r.Match) {
		return []types.Issue{{
			File:    file,
			Range:   body.Range(),
			Message: r.Message,
			RuleID:  r.RuleID,
		}}
	}
	return nil
}

func (r AlwaysFlag) Finish() []types.Issue {
	return make([]types.Issue, 0)
}

type NeverFlag struct {
	RuleID  string
	Message string
}

func (r NeverFlag) ID() string { return r.RuleID }

func (r NeverFlag) META() types.RuleMeta {
	return types.RuleMeta{
		Title:       "NeverFlag",
		Description: r.Message,
		Severity:    "HIGH",
		DocsURL:     "tbd",
	}
}

func (r NeverFlag) Apply(_ string, f *hcl.File) []types.Issue {
	body, _ := f.Body.(*hclsyntax.Body)
	if body == nil {
		return nil
	}
	return []types.Issue{}
}

func (r NeverFlag) Finish() []types.Issue {
	return make([]types.Issue, 0)
}
