//go:build test

package testutil

import "github.com/Marcel2603/tfcoach/internal/engine"

func NewEngineWith(src engine.Source, rules []engine.Rule) *engine.Engine {
	e := engine.New(src)
	e.RegisterMany(rules)
	return e
}
