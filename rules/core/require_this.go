package core

import (
	"github.com/Marcel2603/tfcoach/internal/engine"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type RequireThis struct {
}

func RequireThisRule() RequireThis {
	return RequireThis{}
}

const (
	requireThisId        = "core.test_rule"
	requireThisIdMessage = `resource name must be "this"`
)

func (RequireThis) ID() string { return requireThisId }

func (RequireThis) Apply(file string, f *hcl.File) []engine.Issue {
	body, ok := f.Body.(*hclsyntax.Body)
	if !ok {
		return nil
	}
	var out []engine.Issue
	for _, blk := range body.Blocks {
		if blk.Type != "resource" {
			continue
		}
		// resource "<type>" "<name>"
		if len(blk.Labels) >= 2 && blk.Labels[1] != "this" {
			out = append(out, engine.Issue{
				File:    file,
				Range:   blk.DefRange(),
				Message: requireThisIdMessage,
				RuleID:  requireThisId,
			})
		}
	}
	return out
}
