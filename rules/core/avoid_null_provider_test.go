package core_test

import (
	"strings"
	"testing"

	"github.com/Marcel2603/tfcoach/internal/constants"
	"github.com/Marcel2603/tfcoach/internal/testutil"
	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/Marcel2603/tfcoach/rules/core"
)

func TestAvoidNullProvider_ExpectedMETA(t *testing.T) {
	rule := core.AvoidNullProviderRule()

	expectedMETA := types.RuleMeta{
		Title:       "Avoid using hashicorp/null provider",
		Description: "With newer Terraform version, use locals and terraform_data as native replacement for hashicorp/null",
		Severity:    constants.SeverityMedium,
		DocsURI:     strings.ReplaceAll(rule.ID(), ".", "/"),
	}

	if rule.META() != expectedMETA {
		t.Fatalf("meta mismatch; got %+s, wanted %+s", rule.META(), expectedMETA)
	}
}

func TestAvoidNullProvider_AllGood(t *testing.T) {
	t.Parallel()

	rule := core.AvoidNullProviderRule()

	cases := []struct {
		name     string
		filename string
		resource string
	}{
		{"terraform_data",
			"main.tf",
			`resource "terraform_data" "this" {}`,
		},
		{"locals",
			"locals.tf",
			`locals {}`,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			filename := tt.filename
			resource := tt.resource
			issues := rule.Apply(filename, testutil.ParseToHcl(t, filename, resource))
			if len(issues) != 0 {
				t.Fatalf("Issues found; expected none; got %d: %#v", len(issues), issues)
			}
		})
	}
}

func TestAvoidNullProvider_FailOnUsingHashicorpNull(t *testing.T) {
	t.Parallel()

	rule := core.AvoidNullProviderRule()

	cases := []struct {
		name            string
		filename        string
		resource        string
		expectedMessage string
	}{
		{"null_resource",
			"main.tf",
			`resource "null_resource" "this" {}`,
			"Use terraform_data instead of null_resource",
		},
		{"null_data_source",
			"data.tf",
			`data "null_data_source" "this" {}`,
			"Use locals instead of null_data_source",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			filename := tt.filename
			resource := tt.resource
			issues := rule.Apply(filename, testutil.ParseToHcl(t, filename, resource))
			if len(issues) != 1 {
				t.Fatalf("expected 1 issue; got %d: %#v", len(issues), issues)
			}
			issueMessage := issues[0].Message
			if tt.expectedMessage != issueMessage {
				t.Fatalf("expected to found an issue with message %s, got %s", tt.expectedMessage, issueMessage)
			}
		})
	}
}

func TestAvoidNullProvider_FinishShouldDoNothing(t *testing.T) {
	rule := core.AvoidNullProviderRule()

	issues := rule.Finish()
	if len(issues) != 0 {
		t.Fatalf("Issues found; expected none; got %d: %#v", len(issues), issues)
	}
}
