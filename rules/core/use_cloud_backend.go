package core

import (
	"slices"
	"strings"

	"github.com/Marcel2603/tfcoach/internal/constants"
	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type UseCloudBackend struct {
	id            string
	foundBackends *types.Set[types.DetectedBlock]
}

func UseCloudBackendRule() *UseCloudBackend {
	return &UseCloudBackend{
		id:            rulePrefix + ".use_cloud_backend",
		foundBackends: &types.Set[types.DetectedBlock]{},
	}
}

func (u *UseCloudBackend) ID() string {
	return u.id
}

func (u *UseCloudBackend) META() types.RuleMeta {
	return types.RuleMeta{
		Title:       "Use a cloud backend to store the state",
		Description: "To store the Terraform state securely, define a cloud backend",
		Severity:    constants.SeverityHigh,
		DocsURI:     strings.ReplaceAll(u.id, ".", "/"),
	}
}

func (u *UseCloudBackend) Apply(file string, f *hcl.File) []types.Issue {
	body, ok := f.Body.(*hclsyntax.Body)
	if !ok {
		return nil
	}
	for _, blk := range body.Blocks {
		if blk.Type == "terraform" {
			for _, child := range blk.Body.Blocks {
				if child.Type == constants.DetectedBlockTypeBackend.Value && len(child.Labels) > 0 {
					u.addBlock(constants.DetectedBlockTypeBackend, child.Labels[0], file, child.Range())
				}
				if child.Type == constants.DetectedBlockTypeCloud.Value {
					u.addBlock(constants.DetectedBlockTypeCloud, "cloud", file, child.Range())
				}
			}
		}
	}
	// Issues will be emitted after all files are parsed
	return []types.Issue{}
}

func (u *UseCloudBackend) Finish() []types.Issue {
	blocks := u.foundBackends
	if blocks.Len() == 0 {
		return []types.Issue{
			{
				RuleID:  u.id,
				Message: "No backend configured. State will not be stored remotely",
				Range:   hcl.Range{},
			},
		}
	}
	blockValues := blocks.Values()
	localBackendIndex := slices.IndexFunc(blockValues, func(b types.DetectedBlock) bool { return b.Name == "local" })
	if localBackendIndex != -1 {
		return []types.Issue{
			{
				RuleID:  u.id,
				File:    blockValues[localBackendIndex].File,
				Message: "Local backend configured. State will be stored locally",
				Range:   blockValues[localBackendIndex].Range,
			},
		}
	}
	return []types.Issue{}
}

func (u *UseCloudBackend) addBlock(detectedBlockType types.DetectedBlockType, name string, file string, blockRange hcl.Range) {
	u.foundBackends.Add(types.DetectedBlock{Name: name, File: file, Range: blockRange, Type: detectedBlockType})
}
