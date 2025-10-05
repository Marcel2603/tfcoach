package engine

import (
	"sort"
	"sync"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

const issuesChanBufSize = 3 // TODO later: choose appropriate buffer size (balance performance vs resource usage)

type Engine struct {
	src   Source
	rules []Rule
}

func New(src Source) *Engine {
	return &Engine{src: src, rules: []Rule{}}
}

func (e *Engine) Register(r Rule) {
	e.rules = append(e.rules, r)
}

func (e *Engine) RegisterMany(r []Rule) {
	for _, rule := range r {
		e.Register(rule)
	}
}

func (e *Engine) Run(root string) ([]Issue, error) {
	files, err := e.src.List(root)
	if err != nil {
		return nil, err
	}

	issuesChan := make(chan Issue, issuesChanBufSize)
	fileDoneChan := make(chan struct{})
	var wg sync.WaitGroup
	for _, path := range files {
		wg.Go(func() {
			e.processFile(path, issuesChan)
			fileDoneChan <- struct{}{}
		})
	}
	wg.Go(func() {
		closeAfterSignalCount(len(files), fileDoneChan)
		close(issuesChan)
	})

	var issues []Issue
	for issue := range issuesChan {
		issues = append(issues, issue)
	}
	wg.Wait()

	// sort for deterministic output
	sort.SliceStable(issues, func(i, j int) bool {
		a, b := issues[i], issues[j]
		if a.File != b.File {
			return a.File < b.File
		}
		if a.Range.Start.Line != b.Range.Start.Line {
			return a.Range.Start.Line < b.Range.Start.Line
		}
		if a.Range.Start.Column != b.Range.Start.Column {
			return a.Range.Start.Column < b.Range.Start.Column
		}
		if a.RuleID != b.RuleID {
			return a.RuleID < b.RuleID
		}
		return a.Message < b.Message
	})

	return issues, nil
}

func (e *Engine) processFile(path string, issuesChan chan<- Issue) {
	bytes, err := e.src.ReadFile(path)
	if err != nil {
		issuesChan <- Issue{
			File:    path,
			Message: "read error: " + err.Error(),
			RuleID:  "io",
		}
		return
	}

	hclFile, diagnostics := hclsyntax.ParseConfig(bytes, path, hcl.InitialPos)
	if diagnostics.HasErrors() {
		issuesChan <- Issue{
			File:    path,
			Message: "parse error: " + diagnostics.Error(),
			RuleID:  "parser",
		}
		return
	}

	var fileWg sync.WaitGroup
	ruleDoneChan := make(chan struct{})
	for _, rule := range e.rules {
		fileWg.Go(func() {
			for _, issue := range rule.Apply(path, hclFile) {
				issuesChan <- issue
			}
			ruleDoneChan <- struct{}{}
		})
	}

	fileWg.Go(func() {
		closeAfterSignalCount(len(e.rules), ruleDoneChan)
	})

	fileWg.Wait()
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
