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

var categoryOrder = map[string]int{
	"count":      0,
	"for_each":   1,
	"non_block":  2,
	"block":      3,
	"lifecycle":  4,
	"depends_on": 5,
}

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
	var foundCategories []int

	// TODO #20: need to sort attributes for consistent ordering here?
	for _, attr := range *attributes {
		if slices.Contains(metaArguments, attr.Name) {
			prio, ok := categoryOrder[attr.Name]
			if !ok {
				// unexpected attribute! probably misconfiguration of this rule
				return false
			}
			foundCategories = append(foundCategories, prio)
		} else if isNonBlockExpression(attr.Expr) {
			foundCategories = append(foundCategories, categoryOrder["non_block"])
		} else {
			foundCategories = append(foundCategories, categoryOrder["block"])
		}
	}

	fmt.Println(foundCategories)
	// check if the list of categories in order of appearance is correctly sorted
	previous := 0
	for _, foundCategory := range foundCategories {
		if foundCategory < previous {
			return false
		}
		previous = foundCategory
	}
	return true
}

func isNonBlockExpression(_ hclsyntax.Expression) bool {
	// TODO #20: switch case or reflect or better?
	return false
}
