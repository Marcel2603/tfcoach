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

var (
	categoryOrder = map[string]int{
		"count":      0,
		"for_each":   0,
		"non_block":  10,
		"block":      20,
		"lifecycle":  30,
		"depends_on": 40,
	}
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
	// detect parameters in all attributes and blocks
	var detectedParams []detectedParam
	for _, attr := range body.Attributes {
		detectedParams = append(detectedParams, detectFromAttribute(attr))
	}
	for _, blk := range body.Blocks {
		detectedParams = append(detectedParams, detectFromBlock(blk))
	}

	// order detected parameters by their position in the file
	slices.SortStableFunc(detectedParams, func(a, b detectedParam) int { return a.compare(b) })

	// assign expected order to each parameter
	var foundCategories []int
	for _, param := range detectedParams {
		order, ok := categoryOrder[param.paramType]
		if ok {
			foundCategories = append(foundCategories, order)
		}
		// TODO later: log warnings if not found
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

func detectFromAttribute(attr *hclsyntax.Attribute) detectedParam {
	var paramType string
	switch attr.Name {
	case "count", "for_each", "depends_on":
		paramType = attr.Name
	default:
		paramType = "non_block"
	}

	return detectedParam{
		paramType: paramType,
		startPos:  attr.Range().Start,
	}
}

func detectFromBlock(blk *hclsyntax.Block) detectedParam {
	var paramType string
	switch blk.Type {
	case "lifecycle":
		paramType = blk.Type
	default:
		paramType = "block"
	}

	return detectedParam{
		paramType: paramType,
		startPos:  blk.Range().Start,
	}
}
