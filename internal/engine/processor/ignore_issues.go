package processor

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

const (
	ignoreFileWord = "#tfcoach-ignore-file"
	ignoreRuleWord = "#tfcoach-ignore"
)

type ruleIgnore struct {
	ruleID   string
	hclRange hcl.Range
	path     string
}

type IgnoreIssuesProcessor struct {
	ignoreFiles []ruleIgnore
	ignoreRules []ruleIgnore
}

func NewIgnoreIssuesProcessor() *IgnoreIssuesProcessor {
	return &IgnoreIssuesProcessor{ignoreFiles: []ruleIgnore{}, ignoreRules: []ruleIgnore{}}
}

func (p *IgnoreIssuesProcessor) ScanFile(bytes []byte, hclFile *hcl.File, path string) {
	tokens, _ := hclsyntax.LexConfig(bytes, path, hcl.InitialPos)
	body, _ := hclFile.Body.(*hclsyntax.Body)
	for _, tok := range tokens {
		if tok.Type == hclsyntax.TokenComment {
			comment := string(tok.Bytes)
			comment = strings.Join(strings.Fields(comment), "")
			if strings.HasPrefix(comment, ignoreFileWord) {
				ignoredFileRules := p.processIgnoreFile(comment, path)
				if len(ignoredFileRules) > 0 {
					p.ignoreFiles = p.appendUniqueRuleIgnores(p.ignoreFiles, ignoredFileRules)
				}
			} else {
				if strings.HasPrefix(comment, ignoreRuleWord) {
					ignoredRules := p.processIgnoreRule(comment, path, tok.Range, body)
					if len(ignoredRules) > 0 {
						p.ignoreRules = p.appendUniqueRuleIgnores(p.ignoreRules, ignoredRules)
					}
				}
			}
		}
	}
}

func (*IgnoreIssuesProcessor) appendUniqueRuleIgnores(ignoredRules []ruleIgnore, newRules []ruleIgnore) []ruleIgnore {
	for _, newRule := range newRules {
		if slices.Contains(ignoredRules, newRule) {
			continue
		}
		ignoredRules = append(ignoredRules, newRule)
	}
	return ignoredRules
}

func (p *IgnoreIssuesProcessor) ProcessIssues(issues []types.Issue) []types.Issue {
	filteredIssues := issues[:0]
	for _, issue := range issues {
		if p.containsFileIgnoreRule(issue) {
			fmt.Printf("Ignored Issue %s on File %s \n", issue.RuleID, issue.File)
			continue
		}
		if p.containsRuleIgnoreComment(issue) {
			fmt.Printf("Ignored Issue %s on Range %d \n", issue.RuleID, issue.Range.Start.Line)
			continue
		}
		filteredIssues = append(filteredIssues, issue)
	}
	return filteredIssues
}

func (p *IgnoreIssuesProcessor) processIgnoreFile(comment string, path string) []ruleIgnore {
	commentSplit := strings.SplitN(comment, ":", 2)
	if len(commentSplit) != 2 {
		return []ruleIgnore{}
	}
	var ignoredRules []ruleIgnore
	for id := range strings.SplitSeq(commentSplit[1], ",") {
		ignoredRules = append(ignoredRules, ruleIgnore{ruleID: id, path: path})
	}
	return ignoredRules
}

func (p *IgnoreIssuesProcessor) findNearestBlock(body *hclsyntax.Body, pos hcl.Pos) hcl.Range {
	var nearestRange hcl.Range

	for _, block := range body.Blocks {
		start := block.Range().Start

		if pos.Line > start.Line {
			if block.Range().End.Line > pos.Line {
				return nearestRange
			}
			continue
		}

		return block.Range()
	}
	return nearestRange
}

func (p *IgnoreIssuesProcessor) processIgnoreRule(comment string, path string, hclRange hcl.Range, body *hclsyntax.Body) []ruleIgnore {
	commentSplit := strings.SplitN(comment, ":", 2)
	if len(commentSplit) != 2 {
		return []ruleIgnore{}
	}
	nearestRange := p.findNearestBlock(body, hclRange.Start)
	if nearestRange == (hcl.Range{}) {
		return []ruleIgnore{}
	}
	var ignoredRules []ruleIgnore
	for id := range strings.SplitSeq(commentSplit[1], ",") {
		ignoredRules = append(ignoredRules, ruleIgnore{ruleID: id, path: path, hclRange: nearestRange})
	}
	return ignoredRules
}

func (p *IgnoreIssuesProcessor) containsFileIgnoreRule(issue types.Issue) bool {
	for _, ignoredRule := range p.ignoreFiles {
		if ignoredRule.path == issue.File && ignoredRule.ruleID == issue.RuleID {
			return true
		}
	}
	return false
}

func (p *IgnoreIssuesProcessor) containsRuleIgnoreComment(issue types.Issue) bool {
	for _, ignoredRule := range p.ignoreRules {
		if ignoredRule.path != issue.File {
			continue
		}

		if ignoredRule.hclRange.ContainsPos(issue.Range.Start) && ignoredRule.ruleID == issue.RuleID {
			return true
		}
	}
	return false
}
