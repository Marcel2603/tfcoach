package core

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Marcel2603/tfcoach/internal/constants"
	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type ResourceParameterOrder struct {
	id string
}

func ResourceParameterOrderRule() *ResourceParameterOrder {
	return &ResourceParameterOrder{
		id: rulePrefix + ".resource_parameter_order",
	}
}

func (r *ResourceParameterOrder) ID() string {
	return r.id
}

func (r *ResourceParameterOrder) META() types.RuleMeta {
	return types.RuleMeta{
		Title:       "Resource Parameter Order",
		Description: "Resource parameters should follow a consistent order",
		Severity:    constants.SeverityMedium,
		DocsURI:     strings.ReplaceAll(r.id, ".", "/"),
	}
}

func (r *ResourceParameterOrder) Apply(path string, f *hcl.File) []types.Issue {
	body, ok := f.Body.(*hclsyntax.Body)
	if !ok {
		return nil
	}
	var out []types.Issue
	for _, blk := range body.Blocks {
		if blk.Type == "resource" {
			// TODO #20: iterating on attributes is not enough! ignores lifecycle and block parameters
			if !isParameterOrderCorrect(&blk.Body.Attributes) {
				out = append(out, types.Issue{
					File:    path,
					Range:   blk.Range(),
					Message: fmt.Sprintf("Parameter order in resource \"%s\" is incorrect", blk.Labels[1]),
					RuleID:  r.id,
				})
			}
		}
	}
	return out
}

func (*ResourceParameterOrder) Finish() []types.Issue {
	return []types.Issue{}
}

func isParameterOrderCorrect(attributes *hclsyntax.Attributes) bool {
	metaArguments := []string{"count", "for_each", "lifecycle", "depends_on"}
	var categorizedAttributes []string

	// TODO #20: need to sort attributes for consistent ordering here?
	for _, attr := range *attributes {
		if slices.Contains(metaArguments, attr.Name) {
			categorizedAttributes = append(categorizedAttributes, attr.Name)
		} else if isNonBlockExpression(attr.Expr) {
			categorizedAttributes = append(categorizedAttributes, "non_block")
		} else {
			categorizedAttributes = append(categorizedAttributes, "block")
		}
	}

	fmt.Println(categorizedAttributes)
	correctCategoryOrder := []string{"count", "for_each", "non_block", "block", "lifecycle", "depends_on"}
	refCatIdx := 0
	for _, currentCat := range categorizedAttributes {
		orderIdx := slices.IndexFunc(correctCategoryOrder, func(c string) bool { return c == currentCat })
		if orderIdx == -1 {
			// TODO #20: not found, should we test this?
			return false
		}
		if orderIdx < refCatIdx {
			// repetition of "earlier" category: wrong order
			return false
		}
		// we expect only this category or "later" from now on
		refCatIdx = orderIdx
	}
	return true
}

func isNonBlockExpression(_ hclsyntax.Expression) bool {
	// TODO #20: switch case or reflect or better?
	return false
}
