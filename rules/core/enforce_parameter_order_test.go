package core_test

import (
	"strings"
	"testing"

	"github.com/Marcel2603/tfcoach/internal/constants"
	"github.com/Marcel2603/tfcoach/internal/testutil"
	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/Marcel2603/tfcoach/rules/core"
)

func TestEnforceParameterOrder_ExpectedMETA(t *testing.T) {
	rule := core.EnforceParameterOrderRule()

	expectedMETA := types.RuleMeta{
		Title:       "Enforce Parameter Order",
		Description: "Enforce parameters should follow a consistent order",
		Severity:    constants.SeverityMedium,
		DocsURI:     strings.ReplaceAll(rule.ID(), ".", "/"),
	}

	if rule.META() != expectedMETA {
		t.Fatalf("meta mismatch; got %s, wanted %s", rule.META(), expectedMETA)
	}
}

func TestEnforceParameterOrder_AllGood(t *testing.T) {
	cases := []struct {
		name        string
		fileContent string
	}{
		{"resource_with_no_parameters",
			`resource "aws_instance" "web" {}`,
		},
		{"resource_with_count",
			`resource "aws_instance" "web" {
  count = 1
  ami = data.aws_ami.web.id
}`,
		},
		{
			"resource_with_for_each",
			`resource "aws_instance" "web" {
  for_each = var.instances
  ami = each.key
  availability_zone = each.value
}`,
		},
		{
			"resource_with_lifecycle",
			`resource "aws_instance" "web" {
  ami = 1234
  lifecycle {
    ignore_changes = [tags]
  }
}`,
		},
		{
			"resource_with_depends_on",
			`resource "aws_instance" "web" {
  ami = 1234
  depends_on = [
    aws_iam_role_policy.test
  ]
}`,
		},
		{
			"resource_with_non_block_and_block_parameters",
			`resource "aws_instance" "web" {
  ami = 1234
  availability_zone = var.az
  instance_market_options {
    market_type = "spot"
    spot_options {
      max_price = 0.002
    }
  }
}`,
		},
		{
			"resource_with_multiple_in_correct_order_1",
			`resource "aws_instance" "web" {
  count = 1
  ami = data.aws_ami.web.id
  instance_market_options {
    spot_options {
      max_price = 0.002
    }
  }
  lifecycle {
    ignore_changes = [tags]
  }
  depends_on = [
    aws_iam_role_policy.test
  ]
}`,
		},
		{
			"resource_with_multiple_in_correct_order_2",
			`resource "aws_instance" "web" {
  for_each = var.instances
  ami = each.key
  availability_zone = each.value
  instance_market_options {
    market_type = "spot"
  }
  depends_on = [
    aws_iam_role_policy.test
  ]
}`,
		},
		{
			"module_with_multiple_in_correct_order",
			`module "ec2_instance" {
  count   = length(local.instance_names)
  source  = "terraform-aws-modules/ec2-instance/aws"
  version = "6.0.2"

  name           = local.instance_names[count.index]
  ami            = data.aws_ami.latest_amazon_linux.id
  instance_type  = "t2.micro"

  depends_on = [aws_s3_bucket.example]
}`,
		},
		{
			"data_with_multiple_in_correct_order",
			`data "aws_ami" "latest" {
  count = 3
  name = "test"
  lifecycle {
    ignore_changes = [tags]
  }
  depends_on = [aws_s3_bucket.example]
}`,
		},
		{
			"ephemeral_with_multiple_in_correct_order",
			`ephemeral "aws_instance" "web" {
  count = 1
  ami = data.aws_ami.web.id
  instance_market_options {
    spot_options {
      max_price = 0.002
    }
  }
  lifecycle {
    ignore_changes = [tags]
  }
  depends_on = [
    aws_iam_role_policy.test
  ]
}`,
		},
		{
			"output_with_depends_on",
			`output "my_output" {
  value = "test"
  sensitive = true
  depends_on = [
    aws_iam_role_policy.test
  ]
}`,
		},
	}

	rule := core.EnforceParameterOrderRule()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			issues := rule.Apply("a.tf", testutil.ParseToHcl(t, "a.tf", tt.fileContent))
			if len(issues) != 0 {
				t.Fatalf("Issues found; expected none; got %d: %#v", len(issues), issues)
			}
		})
	}
}

func TestEnforceParameterOrder_ShouldComplain(t *testing.T) {
	cases := []struct {
		name        string
		issueCount  int
		fileContent string
	}{
		{"resource_count_not_first",
			1,
			`resource "aws_instance" "web" {
  ami = data.aws_ami.web.id
  count = 1
}`,
		},
		{
			"resource_for_each_not_first",
			1,
			`resource "aws_instance" "web" {
  instance_type = "t4g.nano"
  for_each = var.instances
  ami = each.key
  availability_zone = each.value
}`,
		},
		{
			"resource_lifecycle_too_high",
			1,
			`resource "aws_instance" "web" {
  ami = 1234
  lifecycle {
    ignore_changes = [tags]
  }
  tags = {
    Name = "test"
  }
}`,
		},
		{
			"resource_depends_on_too_high",
			1,
			`resource "aws_instance" "web" {
  ami = 1234
  depends_on = [
    aws_iam_role_policy.test
  ]
  tags = {
    Name = "test"
  }
}`,
		},
		{
			"resource_depends_on_before_lifecycle",
			1,
			`resource "aws_instance" "web" {
  count = 1
  ami = data.aws_ami.web.id
  depends_on = [
    aws_iam_role_policy.test
  ]
  lifecycle {
    ignore_changes = [tags]
  }
}`,
		},
		{
			// TODO #20: do we want to apply the rule for nested blocks?
			"resource_non_block_after_block_parameters",
			1,
			`resource "aws_instance" "web" {
  ami = 1234
  instance_market_options {
    market_type = "spot"
    spot_options {
      max_price = 0.002
    }
  }
  availability_zone = var.az
}`,
		},
		{
			"resource_multiple_non_compliant_resources",
			3,
			`resource "aws_instance" "non_compliant_1" {
  ami = 1234
  instance_market_options {
    market_type = "spot"
  }
  availability_zone = var.az
}

resource "aws_instance" "compliant" {
  count = 1
  ami = 5678
  availability_zone = var.az
  instance_market_options {
    market_type = "spot"
  }
  depends_on = [aws_iam_role_policy.test]
}

resource "aws_instance" "non_compliant_2" {
  ami = 4321
  depends_on = [aws_iam_role_policy.test]
  availability_zone = var.az
}

resource "aws_instance" "non_compliant_3" {
  ami = 8765
  availability_zone = var.az
  count = 1
}`,
		},
		{
			"module_depends_on_too_high",
			1,
			`module "ec2_instance" {
  count   = length(local.instance_names)
  source  = "terraform-aws-modules/ec2-instance/aws"
  version = "6.0.2"
  depends_on = [aws_s3_bucket.example]

  name           = local.instance_names[count.index]
  instance_type  = "t2.micro"
}`,
		},
		{
			"module_count_not_first",
			1,
			`module "ec2_instance" {
  source  = "terraform-aws-modules/ec2-instance/aws"
  version = "6.0.2"
  count   = length(local.instance_names)

  name           = local.instance_names[count.index]
  instance_type  = "t2.micro"
  depends_on = [aws_s3_bucket.example]
}`,
		},
		{
			"data_lifecycle_too_high",
			1,
			`data "aws_ami" "latest" {
  count = 3
  lifecycle {
    ignore_changes = [tags]
  }
  name = "test"
}`,
		},
		{
			"ephemeral_completely_wrong_order",
			1,
			`ephemeral "aws_instance" "web" {
  lifecycle {
    ignore_changes = [tags]
  }
  count = 1
  instance_market_options {
    spot_options {
      max_price = 0.002
    }
  }
  depends_on = [
    aws_iam_role_policy.test
  ]
  ami = data.aws_ami.web.id
}`,
		},
		{
			"output_depends_on_too_high",
			1,
			`output "my_output" {
  depends_on = [
    aws_iam_role_policy.test
  ]
  value = "test"
  sensitive = true
}`,
		},
	}

	rule := core.EnforceParameterOrderRule()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			issues := rule.Apply("a.tf", testutil.ParseToHcl(t, "a.tf", tt.fileContent))
			if len(issues) != tt.issueCount {
				t.Fatalf("Mismatch in reported issue count; want %d; got %d: %#v", tt.issueCount, len(issues), issues)
			}
		})
	}
}

func TestEnforceParameterOrder_IssueMessage(t *testing.T) {
	fileContent := `data "aws_ami" "latest" {
  count = 3
  lifecycle {
    ignore_changes = [tags]
  }
  name = "test"
}

resource "aws_instance" "web" {
  ami = data.aws_ami.web.id
  count = 1
}

output "my_output" {
  depends_on = [
    aws_iam_role_policy.test
  ]
  value = "test"
  sensitive = true
}`
	expectedIssueCount := 3
	expectedIssueMessages := []string{
		"Parameter order in data block \"latest\" is incorrect",
		"Parameter order in resource block \"web\" is incorrect",
		"Parameter order in output block \"my_output\" is incorrect",
	}

	rule := core.EnforceParameterOrderRule()

	issues := rule.Apply("a.tf", testutil.ParseToHcl(t, "a.tf", fileContent))

	if len(issues) != expectedIssueCount {
		t.Fatalf("Mismatch in reported issue count; want %d; got %d: %#v", expectedIssueCount, len(issues), issues)
	}

	for idx, issue := range issues {
		if issue.Message != expectedIssueMessages[idx] {
			t.Fatalf("Mismatch in issue message; want '%s'; got '%s'", expectedIssueMessages[idx], issue.Message)
		}
	}
}

func TestEnforceParameterOrder_FinishShouldDoNothing(t *testing.T) {
	rule := core.EnforceParameterOrderRule()

	issues := rule.Finish()
	if len(issues) != 0 {
		t.Fatalf("Issues found; expected none; got %d: %#v", len(issues), issues)
	}
}
