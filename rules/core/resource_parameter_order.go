package core

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/Marcel2603/tfcoach/internal/constants"
	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type detectedParam struct {
	paramType string
	startPos  hcl.Pos
}

func (d detectedParam) compare(other detectedParam) int {
	if d.startPos.Line != other.startPos.Line {
		return cmp.Compare(d.startPos.Line, other.startPos.Line)
	}
	return cmp.Compare(d.startPos.Column, other.startPos.Column)
}

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
			if !isParameterOrderCorrect(blk.Body) {
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

func isParameterOrderCorrect(body *hclsyntax.Body) bool {
	specialKeywords := []string{"count", "for_each", "lifecycle", "depends_on"}

	var detectedParams []detectedParam
	for _, attr := range body.Attributes {
		var paramType string
		if slices.Contains(specialKeywords, attr.Name) {
			paramType = attr.Name
		} else {
			paramType = "non_block"
		}

		detectedParams = append(detectedParams, detectedParam{
			paramType: paramType,
			startPos:  attr.Range().Start,
		})
	}
	for _, blk := range body.Blocks {
		var paramType string
		if slices.Contains(specialKeywords, blk.Type) {
			paramType = blk.Type
		} else {
			paramType = "block"
		}

		detectedParams = append(detectedParams, detectedParam{
			paramType: paramType,
			startPos:  blk.Range().Start,
		})
	}

	slices.SortStableFunc(detectedParams, func(a, b detectedParam) int { return a.compare(b) })

	var foundCategories []int
	for _, param := range detectedParams {
		foundCategories = append(foundCategories, categoryOrder[param.paramType])
	}

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
