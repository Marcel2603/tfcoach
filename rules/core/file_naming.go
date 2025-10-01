package core

import (
	"fmt"
	"path"
	"strings"

	"github.com/Marcel2603/tfcoach/internal/engine"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type FileNaming struct {
	id string
}

var typeToFile = map[string]string{
	"output":    "outputs.tf",
	"variable":  "variables.tf",
	"locals":    "locals.tf",
	"provider":  "providers.tf",
	"terraform": "terraform.tf",
	//"backend": "backend.tf
	"data": "data.tf",
}

func FileNamingRule() FileNaming {
	return FileNaming{
		id: rulePrefix + ".file_naming",
	}
}

func (r FileNaming) ID() string {
	return r.id
}

func (r FileNaming) META() engine.RuleMeta {
	return engine.RuleMeta{
		Title:       "File Naming",
		Description: "File naming should follow a strict convention.",
		Severity:    "HIGH",
		DocsURL:     strings.ReplaceAll(r.id, ".", "/"),
	}
}

func (r FileNaming) Apply(file string, f *hcl.File) []engine.Issue {
	body, ok := f.Body.(*hclsyntax.Body)
	if !ok {
		return nil
	}
	var out []engine.Issue
	for _, blk := range body.Blocks {
		blkType := blk.Type
		fileName := path.Base(file)
		if compliantFile, ok := typeToFile[blkType]; ok {
			if fileName != compliantFile {
				out = append(out, r.createIssue(file, compliantFile, blkType, blk.Range()))
			}
		}
	}
	return out
}

func (r FileNaming) createIssue(file string, compliantFile string, hclType string, hclRange hcl.Range) engine.Issue {
	return engine.Issue{
		File:    file,
		Range:   hclRange,
		Message: fmt.Sprintf("All %s should be inside of %s.", hclType, compliantFile),
		RuleID:  r.id,
	}
}
