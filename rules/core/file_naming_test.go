package core_test

import (
	"fmt"
	"testing"

	"github.com/Marcel2603/tfcoach/internal/testutil"
	"github.com/Marcel2603/tfcoach/rules/core"
)

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
				t.Fatalf("Uncorrect number of issues; expected one; got %d: %#v", len(issues), issues)
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
	fmt.Println(resource)
	issues := rule.Apply(filename, testutil.ParseToHcl(t, filename, resource))
	if len(issues) != 6 {
		t.Fatalf("Uncorrect number of issues; expected one; got %d: %#v", len(issues), issues)
	}

}
