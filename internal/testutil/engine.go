//go:build test

package testutil

import (
	"github.com/Marcel2603/tfcoach/internal/engine"
	"github.com/Marcel2603/tfcoach/internal/types"
)

func NewEngineWith(src engine.Source, rules []types.Rule) *engine.Engine {
	e := engine.New(src)
	e.RegisterMany(rules)
	return e
}
