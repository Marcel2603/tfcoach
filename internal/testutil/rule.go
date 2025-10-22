//go:build test

package testutil

import (
	"strings"

	"github.com/Marcel2603/tfcoach/internal/constants"
	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type AlwaysFlag struct {
	RuleID  string
	Message string
	Match   string
}

func (r *AlwaysFlag) ID() string { return r.RuleID }

func (r *AlwaysFlag) META() types.RuleMeta {
	return types.RuleMeta{
		Title:       "AlwaysFlag",
		Description: r.Message,
		Severity:    constants.SeverityHigh,
		DocsURI:     "tbd",
	}
}

func (r *AlwaysFlag) Apply(file string, f *hcl.File) []types.Issue {
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

func (r *AlwaysFlag) Finish() []types.Issue {
	return []types.Issue{}
}

type NeverFlag struct {
	RuleID  string
	Message string
}

func (r *NeverFlag) ID() string { return r.RuleID }

func (r *NeverFlag) META() types.RuleMeta {
	return types.RuleMeta{
		Title:       "NeverFlag",
		Description: r.Message,
		Severity:    constants.SeverityHigh,
		DocsURI:     "tbd",
	}
}

func (r *NeverFlag) Apply(_ string, f *hcl.File) []types.Issue {
	body, _ := f.Body.(*hclsyntax.Body)
	if body == nil {
		return nil
	}
	return []types.Issue{}
}

func (r *NeverFlag) Finish() []types.Issue {
	return []types.Issue{}
}

type FlagOnFinish struct {
	RuleID  string
	Message string
}

func (r *FlagOnFinish) ID() string { return r.RuleID }

func (r *FlagOnFinish) META() types.RuleMeta {
	return types.RuleMeta{
		Title:       "FlagOnFinish",
		Description: r.Message,
		Severity:    constants.SeverityHigh,
		DocsURI:     "tbd",
	}
}

func (r *FlagOnFinish) Apply(_ string, f *hcl.File) []types.Issue {
	body, _ := f.Body.(*hclsyntax.Body)
	if body == nil {
		return nil
	}
	return []types.Issue{}
}

func (r *FlagOnFinish) Finish() []types.Issue {
	return []types.Issue{{
		File:    "somefile.tf",
		Range:   hcl.Range{Filename: "somefile.tf", Start: hcl.Pos{Line: 1, Column: 2, Byte: 3}, End: hcl.Pos{Line: 4, Column: 5, Byte: 6}},
		Message: r.Message,
		RuleID:  r.RuleID,
	}}
}
