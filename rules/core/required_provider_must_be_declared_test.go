package core_test

import (
	"testing"

	"github.com/Marcel2603/tfcoach/internal/constants"
	"github.com/Marcel2603/tfcoach/internal/testutil"
	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/Marcel2603/tfcoach/rules/core"
)

func TestRequiredProviderMustBeDeclared_ExpectedMeta(t *testing.T) {
	rule := core.RequiredProviderMustBeDeclaredRule()

	expectedMETA := types.RuleMeta{
		Title:       "Required Provider Must Be Declared",
		Description: "All providers used in resources or data sources are declared in the terraform.required_providers block.",
		Severity:    constants.SeverityMedium,
		DocsURI:     "core/required_provider_must_be_declared",
	}

	actualMeta := rule.META()
	if actualMeta != expectedMETA {
		t.Fatalf("meta mismatch; got %s, wanted %s", actualMeta, expectedMETA)
	}
}

func TestRequiredProviderMustBeDeclared_ShouldFindIssuesInOneFileAtFinish(t *testing.T) {
	f := testutil.ParseToHcl(t, "a.tf", `
		resource "aws_s3_bucket" "this" {}
`)
	rule := core.RequiredProviderMustBeDeclaredRule()
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

func TestRequiredProviderMustBeDeclared_ShouldNotComplainWhenRequiredProviderIsPresent(t *testing.T) {
	f := testutil.ParseToHcl(t, "a.tf", `
	  terraform {
  	    required_providers {
		  aws = {
		    source  = "hashicorp/aws"
		    version = "~> 5.0"
		  }
	    }
	  }
	
	  resource "aws_s3_bucket" "this" {}
`)
	rule := core.RequiredProviderMustBeDeclaredRule()
	rule.Apply("a.tf", f)
	issues := rule.Finish()

	if len(issues) != 0 {
		t.Fatalf("Issues found; expected none; got %d: %#v", len(issues), issues)
	}
}

func TestRequiredProviderMustBeDeclared_ShouldNotComplainWhenRequiredProviderIsDeclaredInDifferentFile(t *testing.T) {
	fileA := testutil.ParseToHcl(t, "fileA.tf", `
	  resource "aws_s3_bucket" "this" {}
`)
	fileB := testutil.ParseToHcl(t, "fileB.tf", `
	  terraform {
  	    required_providers {
		  aws = {
		    source  = "hashicorp/aws"
		    version = "~> 5.0"
		  }
	    }
	  }
`)

	rule := core.RequiredProviderMustBeDeclaredRule()
	rule.Apply("fileA.tf", fileA)
	rule.Apply("fileB.tf", fileB)
	issues := rule.Finish()

	if len(issues) != 0 {
		t.Fatalf("Issues found; expected none; got %d: %#v", len(issues), issues)
	}
}

func TestRequiredProviderMustBeDeclared_ShouldNotComplainForMultipleProviders(t *testing.T) {
	fileA := testutil.ParseToHcl(t, "fileA.tf", `
      terraform {
  	    required_providers {
		  test = {
		    source  = "test"
		    version = "~> 1.0"
		  }
	    }
	  }

	  resource "aws_s3_bucket" "this" {}
      resource "test_a" "a" {}
      resource "test_b" "b" {}
	  resource "azurerm_resource_group" "this" {}
`)
	fileB := testutil.ParseToHcl(t, "fileB.tf", `
	  terraform {
  	    required_providers {
		  aws = {
		    source  = "hashicorp/aws"
		    version = "~> 5.0"
		  }
		  azurerm = {
		    source  = "hashicorp/azurerm"
		    version = "~> 4.0"
		  }
	    }
	  }
`)

	rule := core.RequiredProviderMustBeDeclaredRule()
	rule.Apply("fileA.tf", fileA)
	rule.Apply("fileB.tf", fileB)
	issues := rule.Finish()

	if len(issues) != 0 {
		t.Fatalf("Issues found; expected none; got %d: %#v", len(issues), issues)
	}
}

func TestRequiredProviderMustBeDeclared_ShouldFindOneIssuePerBlockUsingUndeclaredProvider(t *testing.T) {
	f := testutil.ParseToHcl(t, "a.tf", `
		terraform {
          required_providers {
            test = {
			  source  = "test"
			  version = "~> 1.0"
	        }
          }
        }

		resource "aws_s3_bucket" "bucket1" {}
		data "aws_apigatewayv2_api" "apigw1" {}
		resource "aws_cloudfront_distribution" "cf1" {}
		resource "test_a" "a" {}
        resource "azurerm_resource_group" "this" {}
        resource "aws_s3_bucket" "bucket2" {}
`)
	rule := core.RequiredProviderMustBeDeclaredRule()
	rule.Apply("a.tf", f)
	issues := rule.Finish()

	if len(issues) != 5 {
		t.Fatalf("issues mismatch; got %d, wanted 1", len(issues))
	}
}
