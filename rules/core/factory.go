package core

import (
	"github.com/Marcel2603/tfcoach/internal/engine"
)

var (
	rules = []engine.Rule{RequireThisRule()}
)

func All() []engine.Rule {
	return rules
}
