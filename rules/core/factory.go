package core

import (
	"github.com/Marcel2603/tfcoach/cmd/config"
	"github.com/Marcel2603/tfcoach/internal/types"
)

const (
	rulePrefix = "core"
)

var (
	rules = []types.Rule{NamingConventionRule(), FileNamingRule(), RequiredProviderMustBeDeclaredRule()}
)

func All() []types.Rule {
	var enabledRules []types.Rule
	for _, rule := range rules {
		if config.GetConfigByRuleID(rule.ID()).Enabled {
			enabledRules = append(enabledRules, rule)
		}
	}
	return enabledRules
}
