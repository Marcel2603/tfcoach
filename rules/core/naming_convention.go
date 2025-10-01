package core

import (
	"regexp"
	"strings"

	"github.com/Marcel2603/tfcoach/internal/engine"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

var nameFormatRegex = regexp.MustCompile(`^[a-z0-9_]+$`)

type NamingConvention struct {
	id      string
	message string
}

func NamingConventionRule() NamingConvention {
	return NamingConvention{
		id:      rulePrefix + ".naming_convention",
		message: "terraform names should only contain lowercase alphanumeric characters and underscores",
	}
}

func (n NamingConvention) ID() string {
	return n.id
}

func (n NamingConvention) META() engine.RuleMeta {
	return engine.RuleMeta{
		Title:       "NamingConvention",
		Description: n.message,
		Severity:    "HIGH",
		DocsURL:     strings.ReplaceAll(n.id, ".", "/"),
	}
}

func (n NamingConvention) Apply(file string, f *hcl.File) []engine.Issue {
	body, ok := f.Body.(*hclsyntax.Body)
	if !ok {
		return nil
	}
	var out []engine.Issue
	for _, blk := range body.Blocks {
		name := nameOf(blk)
		if name != "" && !nameFormatRegex.MatchString(name) {
			out = append(out, engine.Issue{
				File:    file,
				Range:   blk.Range(),
				Message: n.message,
				RuleID:  n.id,
			})
		}
	}
	return out
}

func nameOf(block *hclsyntax.Block) string {
	// <block_type> "<label1>" "<label2>"
	if len(block.Labels) == 0 {
		return ""
	}
	if block.Type == "resource" || block.Type == "data" {
		return block.Labels[1]
	}
	return block.Labels[0]
}
