package core_test

import (
	"fmt"
	"testing"

	"github.com/Marcel2603/tfcoach/internal/constants"
	"github.com/Marcel2603/tfcoach/internal/testutil"
	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/Marcel2603/tfcoach/rules/core"
)

func TestUseCloudBackend_META(t *testing.T) {
	rule := core.UseCloudBackendRule()

	expectedMETA := types.RuleMeta{
		Title:       "Use a cloud backend to store the state",
		Description: "To store the Terraform state securely, define a cloud backend",
		Severity:    constants.SeverityHigh,
		DocsURI:     "core/use_cloud_backend",
	}
	actualMETA := rule.META()

	if actualMETA != expectedMETA {
		t.Errorf("meta mismatch; got %+v, want %+v", actualMETA, expectedMETA)
	}
}

func TestUseCloudBackend_ShouldFineOneIssueWithoutBackendBlock(t *testing.T) {
	f := testutil.ParseToHcl(t, "a.tf", `resource "terraform_data" "this" {}`)
	rule := core.UseCloudBackendRule()
	rule.Apply("a.tf", f)
	issues := rule.Finish()

	if len(issues) != 1 {
		t.Fatalf("issues mismatch; got %d, wanted 1", len(issues))
	}

	actualRuleID := issues[0].RuleID
	if actualRuleID != rule.ID() {
		t.Fatalf("rule id mismatch; got %s, want %s", actualRuleID, rule.ID())
	}
}

func TestUseCloudBackend_ShouldFineOneIssueForLocalBackend(t *testing.T) {
	f := testutil.ParseToHcl(t, "a.tf", `
	terraform {
		backend "local" {
			path = "relative/path/to/terraform.tfstate"
		}
	}
	`)
	rule := core.UseCloudBackendRule()
	rule.Apply("a.tf", f)
	issues := rule.Finish()

	if len(issues) != 1 {
		t.Fatalf("issues mismatch; got %d, wanted 1", len(issues))
	}

	actualRuleID := issues[0].RuleID
	if actualRuleID != rule.ID() {
		t.Fatalf("rule id mismatch; got %s, want %s", actualRuleID, rule.ID())
	}
}

func TestUseCloudBackend_ShouldComplainOnMultipleFiles(t *testing.T) {
	f := testutil.ParseToHcl(t, "a.tf", `
	resource "aws_s3_bucket" "this" {}
	`)
	backend := testutil.ParseToHcl(t, "backend.tf", `
	terraform {
		backend "s3" {
		}
	}
`)
	rule := core.UseCloudBackendRule()
	rule.Apply("a.tf", f)
	rule.Apply("backend.tf", backend)
	issues := rule.Finish()

	if len(issues) != 0 {
		t.Fatalf("issues mismatch; got %d, wanted 1", len(issues))
	}
}

func TestUseCloudBackend_ShouldComplainOnHashicorpCloud(t *testing.T) {
	f := testutil.ParseToHcl(t, "a.tf", `
	terraform {
		cloud {
		}
	}
	`)
	rule := core.UseCloudBackendRule()
	rule.Apply("a.tf", f)
	issues := rule.Finish()

	if len(issues) != 0 {
		t.Fatalf("issues mismatch; got %d, wanted 1", len(issues))
	}
}

func TestUseCloudBackend_ShouldComplainOnCloudBackends(t *testing.T) {
	tests := []struct {
		backendType string
	}{
		{backendType: "remote"},
		{backendType: "azurerm"},
		{backendType: "consul"},
		{backendType: "cos"},
		{backendType: "gcs"},
		{backendType: "http"},
		{backendType: "kubernetes"},
		{backendType: "oci"},
		{backendType: "oss"},
		{backendType: "pg"},
		{backendType: "s3"},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s backend", tt.backendType), func(t *testing.T) {
			f := testutil.ParseToHcl(t, "a.tf", fmt.Sprintf(`
			terraform {
				backend "%s" {
				}
			}`, tt.backendType))
			rule := core.UseCloudBackendRule()
			rule.Apply("a.tf", f)
			issues := rule.Finish()

			if len(issues) != 0 {
				t.Fatalf("issues mismatch; got %d, wanted 1", len(issues))
			}
		})
	}

}
