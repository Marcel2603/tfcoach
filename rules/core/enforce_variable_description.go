package core

import (
	"fmt"
	"strings"

	"github.com/Marcel2603/tfcoach/internal/constants"
	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type EnforceVariableDescription struct {
	id string
}

func EnforceVariableDescriptionRule() *EnforceVariableDescription {
	return &EnforceVariableDescription{
		id: rulePrefix + ".enforce_variable_description",
	}
}

func (n *EnforceVariableDescription) ID() string {
	return n.id
}

func (n *EnforceVariableDescription) META() types.RuleMeta {
	return types.RuleMeta{
		Title:       "Enforce Variable Description",
		Description: "To understand what that variable does (even if it seems trivial), always add a description",
		Severity:    constants.SeverityMedium,
		DocsURI:     strings.ReplaceAll(n.id, ".", "/"),
	}
}

func (n *EnforceVariableDescription) Apply(file string, f *hcl.File) []types.Issue {
	body, ok := f.Body.(*hclsyntax.Body)
	if !ok {
		return nil
	}
	var out []types.Issue
	for _, blk := range body.Blocks {
		if blk.Type == "variable" {
			if !isDescriptionPresent(&blk.Body.Attributes) {
				out = append(out, types.Issue{
					File:    file,
					Range:   blk.Range(),
					Message: fmt.Sprintf("Variable \"%s\" has no description", blk.Labels[0]),
					RuleID:  n.id,
				})
			}
		}
	}
	return out
}

func (*EnforceVariableDescription) Finish() []types.Issue {
	return []types.Issue{}
}

func isDescriptionPresent(attributes *hclsyntax.Attributes) bool {
	for _, attr := range *attributes {
		if attr.Name == "description" {
			value, err := attr.Expr.Value(&hcl.EvalContext{})
			if err != nil {
				fmt.Println("error while parsing block value, skipping: ", err)
				continue
			}
			if value.AsString() != "" {
				return true
			}
		}
	}

	return false
}
