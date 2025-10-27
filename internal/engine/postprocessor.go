package engine

import (
	"fmt"
	"strings"

	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

const (
	ignoreFileWord = "tfcoach-ignore-file"
	ignoreRuleWord = "tfcoach-ignore"
)

type ruleIgnore struct {
	ruleIDs  []string
	hclRange hcl.Range
	path     string
}

type Postprocessor struct {
	ignoreFiles []ruleIgnore
	ignoreRules []ruleIgnore
}

func NewPostProcessor() *Postprocessor {
	return &Postprocessor{ignoreFiles: []ruleIgnore{}, ignoreRules: []ruleIgnore{}}
}

func (p *Postprocessor) ScanFile(bytes []byte, path string) {
	toks, _ := hclsyntax.LexConfig(bytes, path, hcl.InitialPos)

	for _, tok := range toks {
		if tok.Type == hclsyntax.TokenComment {
			comment := string(tok.Bytes)
			comment = strings.ReplaceAll(comment, " ", "")
			if strings.Contains(comment, ignoreFileWord) {
				ignoredFile, ok := p.processIgnoreFile(comment, path)
				if ok {
					p.ignoreFiles = append(p.ignoreFiles, ignoredFile)
				}
			}
			if strings.Contains(comment, ignoreRuleWord) {
				ignoredRule, ok := p.processIgnoreRule(comment, path, tok.Range)
				if ok {
					p.ignoreRules = append(p.ignoreFiles, ignoredRule)
				}
			}
		}
	}
}

func (p *Postprocessor) processIgnoreFile(comment string, path string) (ruleIgnore, bool) {
	commentSplit := strings.SplitN(comment, ":", 2)
	if len(commentSplit) != 2 {
		return ruleIgnore{}, false
	}
	return ruleIgnore{ruleIDs: strings.Split(commentSplit[1], ","), path: path}, true
}

func (p *Postprocessor) processIgnoreRule(comment string, path string, hclRange hcl.Range) (ruleIgnore, bool) {
	commentSplit := strings.SplitN(comment, ":", 2)
	if len(commentSplit) != 2 {
		return ruleIgnore{}, false
	}
	return ruleIgnore{ruleIDs: strings.Split(commentSplit[1], ","), path: path, hclRange: hclRange}, true
}

func (p *Postprocessor) ProcessIssues(issues *[]types.Issue) {
	fmt.Printf("IgnoreFiles %+v \n", p.ignoreFiles)
	fmt.Printf("IgnoreRules %+v \n", p.ignoreRules)
}
