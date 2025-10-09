package core

import (
	"fmt"
	"maps"
	"path"
	"slices"
	"strings"

	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/Marcel2603/tfcoach/internal/utils"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type FileNaming struct {
	id string
}

var (
	generalTypeToFile = map[string]string{
		"output":   "outputs.tf",
		"variable": "variables.tf",
		"locals":   "locals.tf",
		"provider": "providers.tf",
		"data":     "data.tf",
	}
	terraformBlkTypeToFile = map[string]string{
		"backend":            "backend.tf",
		"cloud":              "backend.tf",
		"required_providers": "terraform.tf",
		"provider_meta":      "terraform.tf",
	}
	defaultTerraformFilename = "terraform.tf"
)

func FileNamingRule() *FileNaming {
	return &FileNaming{
		id: rulePrefix + ".file_naming",
	}
}

func (r *FileNaming) ID() string {
	return r.id
}

func (r *FileNaming) META() types.RuleMeta {
	return types.RuleMeta{
		Title:       "File Naming",
		Description: "File naming should follow a strict convention.",
		Severity:    "HIGH",
		DocsURL:     strings.ReplaceAll(r.id, ".", "/"),
	}
}

func (r *FileNaming) Apply(file string, f *hcl.File) []types.Issue {
	body, ok := f.Body.(*hclsyntax.Body)
	if !ok {
		return nil
	}
	var out []types.Issue
	for _, blk := range body.Blocks {
		blkType := blk.Type
		fileName := path.Base(file)

		if blkType == "terraform" {
			issues := r.analyzeTerraformType(file, fileName, blk)
			if len(issues) > 0 {
				out = append(out, issues...)
			}
			continue
		}
		if compliantFile, ok := generalTypeToFile[blkType]; ok {
			if fileName != compliantFile {
				out = append(out, r.createIssue(file, compliantFile, blkType, "Block", blk.Range()))
			}
		}
	}
	return out
}

func (r *FileNaming) Finish() []types.Issue {
	return make([]types.Issue, 0)
}

func (r *FileNaming) createIssue(file string, compliantFile string, hclType string, hclDataType string, hclRange hcl.Range) types.Issue {
	return types.Issue{
		File:    file,
		Range:   hclRange,
		Message: fmt.Sprintf(`%s "%s" should be inside of %s.`, hclDataType, hclType, compliantFile),
		RuleID:  r.id,
	}
}

func (r *FileNaming) analyzeTerraformType(file string, fileName string, terraformBlk *hclsyntax.Block) []types.Issue {
	var issues []types.Issue
	issues = append(issues, r.analyzeAllowedFilenamesForTerraformBlock(file, fileName, terraformBlk)...)
	for _, blk := range terraformBlk.Body.Blocks {
		typeFile, ok := terraformBlkTypeToFile[blk.Type]
		compliantFilename := defaultTerraformFilename
		if ok {
			compliantFilename = typeFile
		}
		if fileName != compliantFilename {
			issues = append(issues, r.createIssue(file, compliantFilename, blk.Type, "Block", blk.Range()))
		}
	}

	for _, attr := range terraformBlk.Body.Attributes {
		if fileName != defaultTerraformFilename {
			issues = append(issues, r.createIssue(file, defaultTerraformFilename, attr.Name, "Attribute", attr.Range()))
		}
	}

	return issues
}

func (r *FileNaming) analyzeAllowedFilenamesForTerraformBlock(file string, fileName string, terraformBlk *hclsyntax.Block) []types.Issue {
	files := slices.Collect(maps.Values(terraformBlkTypeToFile))
	files = utils.SortAndDeduplicate(append(files, defaultTerraformFilename))
	var issues []types.Issue
	if !slices.Contains(files, fileName) {
		issues = append(issues, r.createIssue(file, fmt.Sprintf("%+v", files), terraformBlk.Type, "Block", terraformBlk.Range()))
	}
	return issues
}
