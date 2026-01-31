package processor

import (
	"strings"

	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"golang.org/x/sync/syncmap"
)

const (
	ignoreFileWord = "#tfcoach-ignore-file"
	ignoreRuleWord = "#tfcoach-ignore"
)

type ruleIgnore struct {
	ruleID   string
	path     string
	hclRange hcl.Range
}

type ruleIgnoreSet struct {
	m syncmap.Map
}

func (r *ruleIgnoreSet) add(i ruleIgnore) {
	r.m.Store(i, struct{}{})
}

func (r *ruleIgnoreSet) values() []ruleIgnore {
	var result []ruleIgnore
	r.m.Range(func(k, _ interface{}) bool {
		rule, ok := k.(ruleIgnore)
		if !ok {
			return false
		}
		result = append(result, rule)
		return true
	})
	return result
}

type IgnoreIssuesProcessor interface {
	ProcessIssues(issues []types.Issue) []types.Issue
	ScanFile(bytes []byte, hclFile *hcl.File, path string)
}

type ignoreIssuesProcessorImpl struct {
	ignoreRules *ruleIgnoreSet
}

func NewIgnoreIssuesProcessor() IgnoreIssuesProcessor {
	return &ignoreIssuesProcessorImpl{
		ignoreRules: &ruleIgnoreSet{},
	}
}

func (p *ignoreIssuesProcessorImpl) ScanFile(bytes []byte, hclFile *hcl.File, path string) {
	tokens, _ := hclsyntax.LexConfig(bytes, path, hcl.InitialPos)
	body, _ := hclFile.Body.(*hclsyntax.Body)
	for _, tok := range tokens {
		if tok.Type == hclsyntax.TokenComment {
			comment := string(tok.Bytes)
			comment = strings.Join(strings.Fields(comment), "")
			if strings.HasPrefix(comment, ignoreFileWord) {
				ignoredFileRules := p.processIgnoreFile(comment, path)
				p.appendUniqueRuleIgnores(p.ignoreRules, ignoredFileRules)
			} else {
				if strings.HasPrefix(comment, ignoreRuleWord) {
					ignoredRules := p.processIgnoreRule(comment, path, tok.Range, body)
					p.appendUniqueRuleIgnores(p.ignoreRules, ignoredRules)
				}
			}
		}
	}
}

func (p *ignoreIssuesProcessorImpl) ProcessIssues(issues []types.Issue) []types.Issue {
	filteredIssues := issues[:0]
	for _, issue := range issues {
		if p.shouldIgnore(issue) {
			continue
		}
		filteredIssues = append(filteredIssues, issue)
	}
	return filteredIssues
}

func (*ignoreIssuesProcessorImpl) appendUniqueRuleIgnores(current *ruleIgnoreSet, additionalRuleIgnores *ruleIgnoreSet) {
	for _, r := range additionalRuleIgnores.values() {
		current.add(r)
	}
}

func (p *ignoreIssuesProcessorImpl) processIgnoreFile(comment string, path string) *ruleIgnoreSet {
	ignoredRules := ruleIgnoreSet{}

	commentSplit := strings.SplitN(comment, ":", 2)
	if len(commentSplit) != 2 {
		return &ignoredRules
	}

	for id := range strings.SplitSeq(commentSplit[1], ",") {
		ignoredRules.add(ruleIgnore{ruleID: id, path: path})
	}
	return &ignoredRules
}

func (p *ignoreIssuesProcessorImpl) processIgnoreRule(comment string, path string, hclRange hcl.Range, body *hclsyntax.Body) *ruleIgnoreSet {
	ignoredRules := ruleIgnoreSet{}

	commentSplit := strings.SplitN(comment, ":", 2)
	if len(commentSplit) != 2 {
		return &ignoredRules
	}

	nearestRange, found := findNearestBlock(body, hclRange.Start)
	if !found {
		return &ignoredRules
	}

	for id := range strings.SplitSeq(commentSplit[1], ",") {
		ignoredRules.add(ruleIgnore{ruleID: id, path: path, hclRange: nearestRange})
	}
	return &ignoredRules
}

func (p *ignoreIssuesProcessorImpl) shouldIgnore(issue types.Issue) bool {
	for _, ignoredRule := range p.ignoreRules.values() {
		if ignoredRule.path != issue.File {
			continue
		}

		if ignoredRule.ruleID != issue.RuleID {
			continue
		}

		if ignoredRule.hclRange == (hcl.Range{}) {
			// rule ignored for whole file
			return true
		}

		if ignoredRule.hclRange.ContainsPos(issue.Range.Start) {
			// rule ignored for block containing the issue
			return true
		}
	}
	return false
}

func findNearestBlock(body *hclsyntax.Body, pos hcl.Pos) (hcl.Range, bool) {
	var nearestRange hcl.Range

	for _, block := range body.Blocks {
		start := block.Range().Start

		if pos.Line > start.Line {
			if block.Range().End.Line > pos.Line {
				return nearestRange, false
			}
			continue
		}

		return block.Range(), true
	}
	return nearestRange, false
}
