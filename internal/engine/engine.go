package engine

import (
	"cmp"
	"slices"
	"strings"
	"sync"

	"github.com/Marcel2603/tfcoach/internal/engine/processor"
	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/Marcel2603/tfcoach/internal/utils"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type Engine struct {
	src   Source
	rules []types.Rule
}

func New(src Source) *Engine {
	return &Engine{src: src, rules: []types.Rule{}}
}

func (e *Engine) Register(r types.Rule) {
	e.rules = append(e.rules, r)
}

func (e *Engine) RegisterMany(r []types.Rule) {
	for _, rule := range r {
		e.Register(rule)
	}
}

func (e *Engine) Run(root string) ([]types.Issue, error) {
	files, err := e.src.List(root)
	if err != nil {
		return nil, err
	}

	// TODO #42: pass .tfcoachignore infos to processor
	ignoreIssuesProcessor := processor.NewIgnoreIssuesProcessor()

	issuesAfterApply := utils.ProcessInParallelChan(files, func(path string, issuesChan chan<- types.Issue) {
		e.processFile(path, issuesChan, ignoreIssuesProcessor)
	})

	issuesAfterFinish := utils.ProcessInParallel(e.rules, types.Rule.Finish)

	issues := ignoreIssuesProcessor.ProcessIssues(slices.Concat(issuesAfterApply, issuesAfterFinish))

	// sort for deterministic output
	slices.SortStableFunc(issues, func(a, b types.Issue) int {
		if a.File != b.File {
			strings.Compare(a.File, b.File)
		}
		if a.Range.Start.Line != b.Range.Start.Line {
			return cmp.Compare(a.Range.Start.Line, b.Range.Start.Line)
		}
		if a.Range.Start.Column != b.Range.Start.Column {
			return cmp.Compare(a.Range.Start.Column, b.Range.Start.Column)
		}
		if a.RuleID != b.RuleID {
			return strings.Compare(a.RuleID, b.RuleID)
		}
		return strings.Compare(a.Message, b.Message)
	})

	return issues, nil
}

func (e *Engine) processFile(path string, issuesChan chan<- types.Issue, postProcessor processor.IgnoreIssuesProcessor) {
	bytes, err := e.src.ReadFile(path)
	if err != nil {
		issuesChan <- types.Issue{
			File:    path,
			Message: "read error: " + err.Error(),
			RuleID:  "io",
		}
		return
	}

	hclFile, diagnostics := hclsyntax.ParseConfig(bytes, path, hcl.InitialPos)
	if diagnostics.HasErrors() {
		issuesChan <- types.Issue{
			File:    path,
			Message: "parse error: " + diagnostics.Error(),
			RuleID:  "parser",
		}
		return
	}
	var fileProcessingGroup sync.WaitGroup
	fileProcessingGroup.Go(func() {
		postProcessor.ScanFile(bytes, hclFile, path)
	})

	applyOnFile := func(r types.Rule) []types.Issue {
		return r.Apply(path, hclFile)
	}
	for _, issue := range utils.ProcessInParallel(e.rules, applyOnFile) {
		issuesChan <- issue
	}

	fileProcessingGroup.Wait()
}
