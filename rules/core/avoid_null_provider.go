package core

import (
	"fmt"
	"strings"

	"github.com/Marcel2603/tfcoach/internal/constants"
	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type AvoidNullProvider struct {
	id string
}

func AvoidNullProviderRule() *AvoidNullProvider {
	return &AvoidNullProvider{
		id: rulePrefix + ".avoid_null_provider",
	}
}

func (r *AvoidNullProvider) ID() string {
	return r.id
}

func (r *AvoidNullProvider) META() types.RuleMeta {
	return types.RuleMeta{
		Title:       "Avoid using hashicorp/null provider",
		Description: "With newer Terraform version, use locals and terraform_data as native replacement for hashicorp/null",
		Severity:    constants.SeverityMedium,
		DocsURI:     strings.ReplaceAll(r.id, ".", "/"),
	}
}

func (r *AvoidNullProvider) Apply(file string, f *hcl.File) []types.Issue {
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

		configurationType := blk.Labels[0]

		if issue := r.checkConfigurationType(configurationType, file, blk); issue != nil {
			out = append(out, *issue)
		}
	}
	return out
}

func (*AvoidNullProvider) Finish() []types.Issue {
	return []types.Issue{}
}

func (r *AvoidNullProvider) checkConfigurationType(configurationType string, file string, blk *hclsyntax.Block) *types.Issue {
	if configurationType == "null_data_source" {
		return &types.Issue{
			File:    file,
			Range:   blk.Range(),
			Message: fmt.Sprintf("Use locals instead of %s", configurationType),
			RuleID:  r.id,
		}
	}

	if configurationType == "null_resource" {
		return &types.Issue{
			File:    file,
			Range:   blk.Range(),
			Message: fmt.Sprintf("Use terraform_data instead of %s", configurationType),
			RuleID:  r.id,
		}
	}
	return nil
}
