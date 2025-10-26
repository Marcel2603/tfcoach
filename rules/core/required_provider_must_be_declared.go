package core

import (
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/Marcel2603/tfcoach/internal/constants"
	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type detectedBlock struct {
	file         string
	resourceName string
	blockRange   hcl.Range
}

type requiredProviders struct {
	sync.RWMutex
	m map[string][]detectedBlock
}

type RequiredProviderMustBeDeclared struct {
	id                string
	requiredProviders requiredProviders
	foundProviders    []string
}

func RequiredProviderMustBeDeclaredRule() *RequiredProviderMustBeDeclared {
	return &RequiredProviderMustBeDeclared{
		id:                rulePrefix + ".required_provider_must_be_declared",
		requiredProviders: requiredProviders{m: make(map[string][]detectedBlock)},
		foundProviders:    []string{"terraform"}, // built-in provider "terraform" always counts as present
	}
}

func (r *RequiredProviderMustBeDeclared) ID() string {
	return r.id
}

func (r *RequiredProviderMustBeDeclared) META() types.RuleMeta {
	return types.RuleMeta{
		Title:       "Required Provider Must Be Declared",
		Description: "All providers used in resources or data sources are declared in the terraform.required_providers block.",
		Severity:    constants.SeverityMedium, // TODO #13: revert to HIGH
		DocsURI:     strings.ReplaceAll(r.id, ".", "/"),
	}
}

func (r *RequiredProviderMustBeDeclared) Apply(file string, f *hcl.File) []types.Issue {
	body, ok := f.Body.(*hclsyntax.Body)
	if !ok {
		return nil
	}
	for _, blk := range body.Blocks {
		switch blk.Type {
		case "resource", "data":
			name := blk.Labels[0]
			provider := strings.Split(name, "_")
			r.addBlockToRequiredProvider(provider[0], file, name, blk.Range())
		case "terraform":
			for _, child := range blk.Body.Blocks {
				if child.Type != "required_providers" {
					continue
				}
				r.addFoundProviders(child.Body)
			}
		}
	}

	// only report issues after all files have been checked
	return []types.Issue{}
}

func (r *RequiredProviderMustBeDeclared) Finish() []types.Issue {
	var issues []types.Issue
	for requiredProvider, detectedBlocks := range r.requiredProviders.m {
		if slices.Contains(r.foundProviders, requiredProvider) {
			continue
		}
		for _, block := range detectedBlocks {
			issues = append(issues, types.Issue{
				File:    block.file,
				Range:   block.blockRange,
				Message: fmt.Sprintf("Block \"%s\" requires provider \"%s\" which is not declared.", block.resourceName, requiredProvider),
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
		detectedBlock{file, resourceName, blockRange},
	)
	r.requiredProviders.Unlock()
}

func (r *RequiredProviderMustBeDeclared) addFoundProviders(requiredProvidersBody *hclsyntax.Body) {
	for _, provider := range requiredProvidersBody.Attributes {
		r.foundProviders = append(r.foundProviders, provider.Name)
	}
}
