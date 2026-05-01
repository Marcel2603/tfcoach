package core

import (
	"cmp"
	"fmt"
	"log/slog"
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
	supportedBlocks = []string{
		"resource",
		"data",
		"module",
		"ephemeral",
		"output",
		"variable",
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

type EnforceParameterOrder struct {
	id string
}

func EnforceParameterOrderRule() *EnforceParameterOrder {
	return &EnforceParameterOrder{
		id: rulePrefix + ".enforce_parameter_order",
	}
}

func (e *EnforceParameterOrder) ID() string {
	return e.id
}

func (e *EnforceParameterOrder) META() types.RuleMeta {
	return types.RuleMeta{
		Title:       "Enforce Parameter Order",
		Description: "Enforce parameters should follow a consistent order",
		Severity:    constants.SeverityMedium,
		DocsURI:     strings.ReplaceAll(e.id, ".", "/"),
	}
}

func (e *EnforceParameterOrder) Apply(path string, f *hcl.File) []types.Issue {
	body, ok := f.Body.(*hclsyntax.Body)
	if !ok {
		return nil
	}
	var out []types.Issue
	for _, blk := range body.Blocks {
		if slices.Contains(supportedBlocks, blk.Type) {
			if !isParameterOrderCorrect(blk.Body) {
				out = append(out, types.Issue{
					File:    path,
					Range:   blk.Range(),
					Message: fmt.Sprintf("Parameter order in %s block \"%s\" is incorrect", blk.Type, nameOf(blk)),
					RuleID:  e.id,
				})
			}
		}
	}
	return out
}

func (*EnforceParameterOrder) Finish() []types.Issue {
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
		slog.Warn("category not found", "type", param.paramType)
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
