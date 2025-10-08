package core

import (
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type DetectedBlock struct {
	file         string
	resourceName string
	blockRange   hcl.Range
}

type RequiredProviders struct {
	sync.RWMutex
	m map[string][]DetectedBlock
}

type RequiredProviderMustBeDeclared struct {
	id                string
	requiredProviders RequiredProviders
	foundProviders    []string
}

func RequiredProviderMustBeDeclaredRule() *RequiredProviderMustBeDeclared {
	return &RequiredProviderMustBeDeclared{
		id:                rulePrefix + ".required_provider_must_be_declared",
		requiredProviders: RequiredProviders{m: make(map[string][]DetectedBlock)},
		foundProviders:    make([]string, 0),
	}
}

func (r *RequiredProviderMustBeDeclared) ID() string {
	return r.id
}

func (r *RequiredProviderMustBeDeclared) META() types.RuleMeta {
	return types.RuleMeta{
		Title:       "Required Provider Must Be Declared",
		Description: "All providers used in resources or data sources are declared in the terraform.required_providers block.",
		Severity:    "HIGH",
		DocsURL:     strings.ReplaceAll(r.id, ".", "/"),
	}
}

func (r *RequiredProviderMustBeDeclared) Apply(file string, f *hcl.File) []types.Issue {
	body, ok := f.Body.(*hclsyntax.Body)
	if !ok {
		return nil
	}
	for _, blk := range body.Blocks {
		blkType := blk.Type
		if blkType == "resource" || blkType == "data" {
			name := blk.Labels[0]
			provider := strings.Split(name, "_")
			r.addBlockToRequiredProvider(provider[0], file, name, blk.Range())
		} else if blkType == "terraform" {
			for _, child := range blk.Body.Blocks {
				if child.Type != "required_providers" {
					continue
				}
				r.addFoundProviders(child.Body)
			}
		}
	}

	// only report issues after all files have been checked
	return make([]types.Issue, 0)
}

func (r *RequiredProviderMustBeDeclared) Finish() []types.Issue {
	var issues []types.Issue
	for requiredProvider, detectedBlocks := range r.requiredProviders.m {
		if slices.Contains(r.foundProviders, requiredProvider) {
			continue
		}
		for _, detectedBlock := range detectedBlocks {
			issues = append(issues, types.Issue{
				File:    detectedBlock.file,
				Range:   detectedBlock.blockRange,
				Message: fmt.Sprintf("Block %s requires provider %s which is not declared.", detectedBlock.resourceName, requiredProvider),
				RuleID:  r.id,
			})
		}
	}
	return issues
}

func (r *RequiredProviderMustBeDeclared) addBlockToRequiredProvider(provider string, file string, resourceName string, blockRange hcl.Range) {
	r.requiredProviders.Lock()
	r.requiredProviders.m[provider] = append(
		r.requiredProviders.m[provider],
		DetectedBlock{file, resourceName, blockRange},
	)
	r.requiredProviders.Unlock()
}

func (r *RequiredProviderMustBeDeclared) addFoundProviders(requiredProvidersBody *hclsyntax.Body) {
	for _, provider := range requiredProvidersBody.Attributes {
		r.foundProviders = append(r.foundProviders, provider.Name)
	}
}
