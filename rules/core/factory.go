package core

import (
	"github.com/Marcel2603/tfcoach/internal/types"
)

const (
	rulePrefix = "core"
)

var (
	rules = []types.Rule{NamingConventionRule(), FileNamingRule()}
)

func All() []types.Rule {
	return rules
}
