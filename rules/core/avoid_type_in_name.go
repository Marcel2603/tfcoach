package core

import (
	"fmt"
	"strings"

	"github.com/Marcel2603/tfcoach/internal/constants"
	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type AvoidTypeInName struct {
	id string
}

func AvoidTypeInNameRule() *AvoidTypeInName {
	return &AvoidTypeInName{
		id: rulePrefix + ".avoid_type_in_name",
	}
}

func (r *AvoidTypeInName) ID() string {
	return r.id
}

func (r *AvoidTypeInName) META() types.RuleMeta {
	return types.RuleMeta{
		Title:       "Avoid Type in Name",
		Description: "Names shouldn't repeat their type.",
		Severity:    constants.SeverityHigh,
		DocsURI:     strings.ReplaceAll(r.id, ".", "/"),
	}
}

func (r *AvoidTypeInName) Apply(file string, f *hcl.File) []types.Issue {
	body, ok := f.Body.(*hclsyntax.Body)
	if !ok {
		return nil
	}
	var out []types.Issue
	for _, blk := range body.Blocks {
		blkType := blk.Type
		if blkType != "resource" && blkType != "data" {
			continue
		}
		resourceTypes := blk.Labels[0]
		resourceName := blk.Labels[1]
		for resourceType := range strings.SplitSeq(resourceTypes, "_") {
			if strings.Contains(resourceName, resourceType) {
				out = append(out, types.Issue{
					RuleID:  r.ID(),
					File:    file,
					Range:   blk.Range(),
					Message: fmt.Sprintf("Block \"%s\" violates naming convention, it should not repeat the type \"%s\"", resourceName, resourceType),
				})
			}
		}
	}
	return out
}

func (*AvoidTypeInName) Finish() []types.Issue {
	return []types.Issue{}
}
