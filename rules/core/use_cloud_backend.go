package core

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Marcel2603/tfcoach/internal/constants"
	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type UseCloudBackend struct {
	id            string
	foundBackends detectedBackendBlocks
}

type detectedBackendBlocks struct {
	sync.RWMutex
	backendBlocks []types.DetectedBlock
	cloudBlocks   []types.DetectedBlock
}

func UseCloudBackendRule() *UseCloudBackend {
	return &UseCloudBackend{
		id: rulePrefix + ".use_cloud_backend",
		foundBackends: detectedBackendBlocks{
			backendBlocks: make([]types.DetectedBlock, 0),
			cloudBlocks:   make([]types.DetectedBlock, 0),
		},
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
				if child.Type == "backend" {
					fmt.Printf("Backend block found %s \n", child.Labels[0])
					u.addBackendBlock(child.Labels[0], file, child.Range())
				}
				if child.Type == "cloud" {
					fmt.Printf("Cloud block found \n")
					u.addCloudBlock("cloud", file, child.Range())
				}
			}
		}
	}
	// Issues will be emitted after all files are parsed
	return []types.Issue{}
}

func (*UseCloudBackend) Finish() []types.Issue {

	return []types.Issue{}
}

func (u *UseCloudBackend) addBackendBlock(name string, file string, blockRange hcl.Range) {
	u.foundBackends.Lock()
	u.foundBackends.backendBlocks = append(
		u.foundBackends.backendBlocks,
		types.DetectedBlock{Name: name, File: file, Range: blockRange},
	)
	u.foundBackends.Unlock()
}

func (u *UseCloudBackend) addCloudBlock(name string, file string, blockRange hcl.Range) {
	u.foundBackends.Lock()
	u.foundBackends.cloudBlocks = append(
		u.foundBackends.cloudBlocks,
		types.DetectedBlock{Name: name, File: file, Range: blockRange},
	)
	u.foundBackends.Unlock()
}
