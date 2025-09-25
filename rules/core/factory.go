package core

import (
	"github.com/Marcel2603/tfcoach/internal/engine"
)

const (
	rulePrefix = "core"
)

var (
	rules = []engine.Rule{NamingConventionRule()}
)

func All() []engine.Rule {
	return rules
}
