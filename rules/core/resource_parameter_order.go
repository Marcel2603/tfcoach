package core

import (
	"strings"

	"github.com/Marcel2603/tfcoach/internal/constants"
	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/hashicorp/hcl/v2"
)

type ResourceParameterOrder struct {
	id string
}

func ResourceParameterOrderRule() *ResourceParameterOrder {
	return &ResourceParameterOrder{
		id: rulePrefix + ".resource_parameter_order",
	}
}

func (r *ResourceParameterOrder) ID() string {
	return r.id
}

func (r *ResourceParameterOrder) META() types.RuleMeta {
	return types.RuleMeta{
		Title:       "Resource Parameter Order",
		Description: "Resource parameters should follow a consistent order",
		Severity:    constants.SeverityMedium,
		DocsURI:     strings.ReplaceAll(r.id, ".", "/"),
	}
}

func (r *ResourceParameterOrder) Apply(_ string, _ *hcl.File) []types.Issue {
	return []types.Issue{}
}

func (*ResourceParameterOrder) Finish() []types.Issue {
	return []types.Issue{}
}
