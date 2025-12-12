package core

import (
	"fmt"

	"github.com/Marcel2603/tfcoach/cmd/config"
	"github.com/Marcel2603/tfcoach/internal/constants"
	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/hashicorp/hcl/v2"
)

const (
	rulePrefix = "core"
)

var (
	rules = []types.Rule{
		NamingConventionRule(),
		FileNamingRule(),
		RequiredProviderMustBeDeclaredRule(),
		EnforceVariableDescriptionRule(),
		AvoidNullProviderRule(),
		AvoidTypeInNameRule(),
	}
	ruleMap = mapRules(rules)
)

func All() []types.Rule {
	return rules
}

func EnabledRules() []types.Rule {
	var enabledRules []types.Rule
	for _, rule := range rules {
		if config.GetConfigByRuleID(rule.ID()).Enabled {
			enabledRules = append(enabledRules, rule)
		}
	}
	return enabledRules
}

func FindByID(id string) (types.Rule, error) {
	rule, ok := ruleMap[id]
	if !ok {
		return nil, fmt.Errorf("no rule found for ID %s", id)
	}
	return rule, nil
}

func mapRules(rulesList []types.Rule) map[string]types.Rule {
	result := make(map[string]types.Rule)
	for _, rule := range rulesList {
		result[rule.ID()] = rule
	}
	return result
}

type UnknownRule struct {
	PseudoID string
}

func (r *UnknownRule) ID() string {
	return r.PseudoID
}

func (*UnknownRule) META() types.RuleMeta {
	return types.RuleMeta{
		Title:       "Unknown",
		Description: "Unknown rule",
		Severity:    constants.SeverityUnknown,
		DocsURI:     "about:blank",
	}
}

func (*UnknownRule) Apply(_ string, _ *hcl.File) []types.Issue {
	return []types.Issue{}
}

func (*UnknownRule) Finish() []types.Issue {
	return []types.Issue{}
}
