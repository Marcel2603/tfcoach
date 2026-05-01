package processor

import (
	"fmt"
	"path/filepath"
	"slices"
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
	fileMatchers             map[string]*dotignore.PatternMatcher // dir -> matcher
}

func NewIgnoreIssuesProcessor(ignoreFiles []string) (IgnoreIssuesProcessor, error) {
	matchers := make(map[string]*dotignore.PatternMatcher, len(ignoreFiles))
	for _, f := range ignoreFiles {
		abs, err := filepath.Abs(f)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve ignore file path: %w", err)
		}
		m, err := dotignore.NewPatternMatcherFromFile(abs)
		if err != nil {
			return nil, fmt.Errorf("failed to parse ignore file %s: %w", abs, err)
		}
		matchers[filepath.Dir(abs)] = m
	}

	return &ignoreIssuesProcessorImpl{
		ignoredRulesAtBlockLevel: &types.Set[ruleIgnore]{},
		ignoredRulesAtFileLevel:  &types.Set[ruleIgnore]{},
		ignoredFiles:             &types.Set[string]{},
		fileMatchers:             matchers,
	}, nil
}

func (p *ignoreIssuesProcessorImpl) matchesIgnoreFile(path string) bool {
	if len(p.fileMatchers) == 0 {
		return false
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	var dirs []string
	for dir := filepath.Dir(absPath); filepath.Dir(dir) != dir; dir = filepath.Dir(dir) {
		dirs = append(dirs, dir)
		if _, ok := p.fileMatchers[dir]; ok {
			break
		}
	}
	slices.Reverse(dirs)

	matched := false
	for _, dir := range dirs {
		matcher, ok := p.fileMatchers[dir]
		if !ok {
			continue
		}
		rel, err := filepath.Rel(dir, absPath)
		if err != nil {
			continue
		}
		if isMatch, _, err := matcher.MatchesWithTracking(rel); err == nil {
			matched = isMatch
		}
	}
	return matched
}

func (p *ignoreIssuesProcessorImpl) ScanFile(bytes []byte, hclFile *hcl.File, path string) {
	if p.matchesIgnoreFile(path) {
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
			} else if strings.HasPrefix(comment, ignoreRuleWord) {
				ignoredRulesForBlock := computeIgnoredRulesForBlock(comment, path, tok.Range, body)
				p.appendUniqueRuleIgnoresAtBlockLevel(ignoredRulesForBlock)
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
