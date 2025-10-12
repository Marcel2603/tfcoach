package core_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/Marcel2603/tfcoach/internal/testutil"
	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/Marcel2603/tfcoach/rules/core"
)

func TestFileNaming_ExpectedMeta(t *testing.T) {
	rule := core.FileNamingRule()

	expectedMETA := types.RuleMeta{
		Title:       "File Naming",
		Description: "File naming should follow a strict convention.",
		Severity:    "HIGH",
		DocsURL:     strings.ReplaceAll(rule.ID(), ".", "/"),
	}

	if rule.META() != expectedMETA {
		t.Fatalf("meta mismatch; got %s, wanted %s", rule.META(), expectedMETA)
	}
}

func TestFileNaming_ShouldIssueSameID(t *testing.T) {
	f := testutil.ParseToHcl(t, "bad.tf", `
		data "test" "test" {}
	`)

	rule := core.FileNamingRule()
	issues := rule.Apply("bad.tf", f)

	if issues[0].RuleID != rule.ID() {
		t.Fatalf("rule id mismatch; got %s, want %s", issues[0].RuleID, rule.ID())
	}
}

func TestFileNaming_MustComplain(t *testing.T) {
	rule := core.FileNamingRule()

	cases := []struct {
		name     string
		filename string
		resource string
	}{
		{"variable",
			"variables.tf",
			`variable "test" {}`,
		},
		{"output",
			"outputs.tf",
			`output "test" {}`,
		},
		{"locals",
			"locals.tf",
			`locals {}`,
		},
		{"data",
			"data.tf",
			`data "test" {}`,
		},
		{"provider",
			"providers.tf",
			`provider "test" {}`,
		},
		{"terraform",
			"terraform.tf",
			`terraform {}`,
		},
		{"resource",
			"anything.tf",
			`resource "anything" "test" {}`,
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

func TestFileNaming_ShouldFailOnWrongFiles(t *testing.T) {
	rule := core.FileNamingRule()

	cases := []struct {
		name     string
		filename string
		resource string
	}{
		{"variable",
			"kdfsfs.tf",
			`variable "test" {}`,
		},
		{"output",
			"sdfsdfc.tf",
			`output "test" {}`,
		},
		{"locals",
			"alocals.tf",
			`locals {}`,
		},
		{"data",
			"main.tf",
			`data "test" {}`,
		},
		{"provider",
			"terraform.tf",
			`provider "test" {}`,
		},
		{"terraform",
			"data.tf",
			`terraform {}`,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			filename := tt.filename
			resource := tt.resource
			issues := rule.Apply(filename, testutil.ParseToHcl(t, filename, resource))
			if len(issues) != 1 {
				t.Fatalf("Incorrect number of issues; expected one; got %d: %#v", len(issues), issues)
			}
		})
	}
}

func TestFileNaming_ShouldFailOnFileWithMultipleIssues(t *testing.T) {
	rule := core.FileNamingRule()

	filename := "main.tf"
	resource := `
	data "archive_file" "zip" {}

	locals {}
	
	resource "null_resource" "test" {}
		
	output "test" {
	  value = "test"
	}
	
	provider "aws" {}
	
	terraform {}
	
	variable "test" {}
	`
	issues := rule.Apply(filename, testutil.ParseToHcl(t, filename, resource))
	if len(issues) != 6 {
		t.Fatalf("Incorrect number of issues; expected one; got %d: %#v", len(issues), issues)
	}
}

func TestFileNaming_FinishShouldDoNothing(t *testing.T) {
	rule := core.FileNamingRule()

	issues := rule.Finish()
	if len(issues) != 0 {
		t.Fatalf("Issues found; expected none; got %d: %#v", len(issues), issues)
	}
}

func TestFileNaming_ShouldFailOnMisconfiguredTerraformBlocks(t *testing.T) {
	rule := core.FileNamingRule()

	terraformWithBackend := `
	terraform {
		backend "s3" {}
		cloud {
		}
	}`

	terraformInNonComplaintFile := `
		terraform {
			required_version = "0.12.0"
		}`
	terraformProviderVersion := `
		terraform {
			required_providers {
				pro {
					  version = "<version-constraint>"
					  source  = "<provider-address>"
					}
			}
		}`

	cases := []struct {
		name     string
		filename string
		resource string
		issues   []string
	}{
		{"terraformWithBackendInTerraform.tf",
			"terraform.tf",
			terraformWithBackend,
			[]string{`Block "backend" should be inside of backend.tf.`,
				`Block "cloud" should be inside of backend.tf.`},
		},
		{
			"terraformBlockInNonComplaintFile",
			"main.tf",
			terraformInNonComplaintFile,
			[]string{`Block "terraform" should be inside of [backend.tf terraform.tf].`,
				`Attribute "required_version" should be inside of terraform.tf.`},
		},
		{
			"terraformProviderVersionInTerraform.tf",
			"providers.tf",
			terraformProviderVersion,
			[]string{`Block "terraform" should be inside of [backend.tf terraform.tf].`,
				`Block "required_providers" should be inside of terraform.tf.`},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			filename := tt.filename
			resource := tt.resource
			issues := rule.Apply(filename, testutil.ParseToHcl(t, filename, resource))
			if len(issues) != len(tt.issues) {
				t.Fatalf("Issues found; expected %d; got %d: %#v", len(tt.issues), len(issues), issues)
			}
			for _, issue := range issues {
				if !slices.Contains(tt.issues, issue.Message) {
					t.Fatalf("Found unexpected Issue %s", issue.Message)
				}
			}
		})
	}

}
