package engine

import (
	"cmp"
	"slices"
	"strings"
	"sync"

	"github.com/Marcel2603/tfcoach/internal/engine/processor"
	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

const issuesChanBufSize = 5 // TODO later: choose appropriate buffer size (balance performance vs resource usage)

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

	issuesChan := make(chan types.Issue, issuesChanBufSize)
	fileDoneChan := make(chan struct{})
	ruleFinishDoneChan := make(chan struct{})
	ignoreIssuesProcessor := processor.NewIgnoreIssuesProcessor()
	var wg sync.WaitGroup

	for _, path := range files {
		wg.Go(func() {
			e.processFile(path, issuesChan, ignoreIssuesProcessor)
			fileDoneChan <- struct{}{}
		})
	}

	wg.Go(func() {
		// wait for all files to have been processed before triggering rule finish
		closeAfterSignalCount(len(files), fileDoneChan)
		for _, rule := range e.rules {
			wg.Go(func() {
				for _, issue := range rule.Finish() {
					issuesChan <- issue
				}
				ruleFinishDoneChan <- struct{}{}
			})
		}
	})

	wg.Go(func() {
		closeAfterSignalCount(len(e.rules), ruleFinishDoneChan)
		close(issuesChan)
	})

	issues := collectAllFromChannel(issuesChan)
	wg.Wait()

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

	issues = ignoreIssuesProcessor.ProcessIssues(issues)

	return issues, nil
}

func (e *Engine) processFile(path string, issuesChan chan<- types.Issue, postProcessor *processor.IgnoreIssuesProcessor) {
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
	ruleApplyDoneChan := make(chan struct{})
	for _, rule := range e.rules {
		fileProcessingGroup.Go(func() {
			for _, issue := range rule.Apply(path, hclFile) {
				issuesChan <- issue
			}
			ruleApplyDoneChan <- struct{}{}
		})
	}

	fileProcessingGroup.Go(func() {
		closeAfterSignalCount(len(e.rules), ruleApplyDoneChan)
	})

	fileProcessingGroup.Wait()
}

func closeAfterSignalCount(target int, signalChannel chan struct{}) {
	defer close(signalChannel)

	if target == 0 {
		return
	}

	signalCount := 0
	for {
		select {
		case <-signalChannel:
			signalCount++
			if signalCount >= target {
				return
			}
		}
	}
}

func collectAllFromChannel(issuesChan <-chan types.Issue) []types.Issue {
	var issues []types.Issue
	for issue := range issuesChan {
		issues = append(issues, issue)
	}
	return issues
}
