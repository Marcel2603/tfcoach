package processor

import (
	"strings"

	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/Marcel2603/tfcoach/internal/utils"
	"github.com/codeglyph/go-dotignore/v2"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
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

type IgnoreIssuesProcessor interface {
	ProcessIssues(issues []types.Issue) []types.Issue
	ScanFile(bytes []byte, hclFile *hcl.File, path string)
}

type ignoreIssuesProcessorImpl struct {
	ignoredRulesAtBlockLevel *types.Set[ruleIgnore]
	ignoredRulesAtFileLevel  *types.Set[ruleIgnore]
	ignoredFiles             *types.Set[string]
	fileIgnorer              *dotignore.RepositoryMatcher
}

func NewIgnoreIssuesProcessor(rootPath string) (IgnoreIssuesProcessor, error) {
	ignorer, err := dotignore.NewRepositoryMatcherWithConfig(
		rootPath,
		&dotignore.RepositoryConfig{IgnoreFileName: ".tfcoachignore"},
	)
	if err != nil {
		return nil, err
	}

	return &ignoreIssuesProcessorImpl{
		ignoredRulesAtBlockLevel: &types.Set[ruleIgnore]{},
		ignoredRulesAtFileLevel:  &types.Set[ruleIgnore]{},
		ignoredFiles:             &types.Set[string]{},
		fileIgnorer:              ignorer,
	}, nil
}

func (p *ignoreIssuesProcessorImpl) ScanFile(bytes []byte, hclFile *hcl.File, path string) {
	shouldIgnore, _ := p.fileIgnorer.Matches(path)
	if shouldIgnore {
		p.ignoredFiles.Add(path)
		return
	}

	tokens, _ := hclsyntax.LexConfig(bytes, path, hcl.InitialPos)
	body, _ := hclFile.Body.(*hclsyntax.Body)
	for _, tok := range tokens {
		if tok.Type == hclsyntax.TokenComment {
			comment := string(tok.Bytes)
			comment = strings.Join(strings.Fields(comment), "")
			if strings.HasPrefix(comment, ignoreFileWord) {
				ignoredRulesForFile := computeIgnoredRulesForFile(comment, path)
				p.appendUniqueRuleIgnoresAtFileLevel(ignoredRulesForFile)
			} else {
				if strings.HasPrefix(comment, ignoreRuleWord) {
					ignoredRulesForBlock := computeIgnoredRulesForBlock(comment, path, tok.Range, body)
					p.appendUniqueRuleIgnoresAtBlockLevel(ignoredRulesForBlock)
				}
			}
		}
	}
}

func (p *ignoreIssuesProcessorImpl) ProcessIssues(issues []types.Issue) []types.Issue {
	processIssue := func(issue types.Issue) []types.Issue {
		if !p.shouldIgnore(issue) {
			return []types.Issue{issue}
		}
		return []types.Issue{}
	}

	return utils.FlatMap(issues, processIssue)
}

func (p *ignoreIssuesProcessorImpl) appendUniqueRuleIgnoresAtBlockLevel(additionalRuleIgnores *types.Set[ruleIgnore]) {
	for _, r := range additionalRuleIgnores.Values() {
		p.ignoredRulesAtBlockLevel.Add(r)
	}
}

func (p *ignoreIssuesProcessorImpl) appendUniqueRuleIgnoresAtFileLevel(additionalRuleIgnores *types.Set[ruleIgnore]) {
	for _, r := range additionalRuleIgnores.Values() {
		p.ignoredRulesAtFileLevel.Add(r)
	}
}

func computeIgnoredRulesForFile(comment string, path string) *types.Set[ruleIgnore] {
	ignoredRules := types.Set[ruleIgnore]{}

	commentSplit := strings.SplitN(comment, ":", 2)
	if len(commentSplit) != 2 {
		return &ignoredRules
	}

	for id := range strings.SplitSeq(commentSplit[1], ",") {
		ignoredRules.Add(ruleIgnore{ruleID: id, path: path})
	}
	return &ignoredRules
}

func computeIgnoredRulesForBlock(comment string, path string, hclRange hcl.Range, body *hclsyntax.Body) *types.Set[ruleIgnore] {
	ignoredRules := types.Set[ruleIgnore]{}

	commentSplit := strings.SplitN(comment, ":", 2)
	if len(commentSplit) != 2 {
		return &ignoredRules
	}

	nearestRange, found := findNearestBlock(body, hclRange.Start)
	if !found {
		return &ignoredRules
	}

	for id := range strings.SplitSeq(commentSplit[1], ",") {
		ignoredRules.Add(ruleIgnore{ruleID: id, path: path, hclRange: nearestRange})
	}
	return &ignoredRules
}

func (p *ignoreIssuesProcessorImpl) shouldIgnore(issue types.Issue) bool {
	if p.ignoredFiles.Has(issue.File) {
		return true
	}

	if p.ignoredRulesAtFileLevel.Has(ruleIgnore{ruleID: issue.RuleID, path: issue.File}) {
		return true
	}

	for _, ignoredRule := range p.ignoredRulesAtBlockLevel.Values() {
		if ignoredRule.path != issue.File {
			continue
		}

		if ignoredRule.ruleID != issue.RuleID {
			continue
		}

		if ignoredRule.hclRange.ContainsPos(issue.Range.Start) {
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
