package engine

import (
	"sort"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

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

	var issues []Issue
	for _, file := range files {
		bytes, err := e.src.ReadFile(file)
		if err != nil {
			issues = append(issues, Issue{
				File:    file,
				Message: "read error: " + err.Error(),
				RuleID:  "io",
			})
			continue
		}

		hclFile, diagnostics := hclsyntax.ParseConfig(bytes, file, hcl.InitialPos)
		if diagnostics.HasErrors() {
			issues = append(issues, Issue{
				File:    file,
				Message: "parse error: " + diagnostics.Error(),
				RuleID:  "parser",
			})
			continue
		}

		for _, rule := range e.rules {
			issues = append(issues, rule.Apply(file, hclFile)...)
		}
	}

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
