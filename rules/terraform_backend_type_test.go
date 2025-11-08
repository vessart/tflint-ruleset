package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformBackendTypeRule_Check(t *testing.T) {
	tests := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "valid backend with correct methods",
			Content: `
terraform {
	backend "http" {
		address        = "https://example.com/state"
		lock_method    = "POST"
		unlock_method  = "DELETE"
	}
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "invalid backend type",
			Content: `
terraform {
	backend "s3" {
		bucket = "example-bucket"
		lock_method = "POST"
		unlock_method = "DELETE"
	}
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformBackendTypeRule(),
					Message: "backend type must be 'http', but found 's3'",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 15},
					},
				},
			},
		},
		{
			Name: "invalid lock_method",
			Content: `
terraform {
	backend "http" {
		address        = "https://example.com/state"
		lock_method    = "GET"
		unlock_method  = "DELETE"
	}
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformBackendTypeRule(),
					Message: `"lock_method" must be "POST"`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 5},
						End:      hcl.Pos{Line: 5, Column: 29},
					},
				},
			},
		},
		{
			Name: "invalid unlock_method",
			Content: `
terraform {
	backend "http" {
		address        = "https://example.com/state"
		lock_method    = "POST"
		unlock_method  = "GET"
	}
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformBackendTypeRule(),
					Message: `"unlock_method" must be "DELETE"`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 6, Column: 5},
						End:      hcl.Pos{Line: 6, Column: 29},
					},
				},
			},
		},
		{
			Name: "missing lock_method",
			Content: `
terraform {
	backend "http" {
		address        = "https://example.com/state"
		unlock_method  = "DELETE"
	}
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformBackendTypeRule(),
					Message: `"lock_method" attribute is required`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 6, Column: 4},
					},
				},
			},
		},
		{
			Name: "missing unlock_method",
			Content: `
terraform {
	backend "http" {
		address        = "https://example.com/state"
		lock_method    = "POST"
	}
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformBackendTypeRule(),
					Message: `"unlock_method" attribute is required`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 6, Column: 4},
					},
				},
			},
		},
	}

	rule := NewTerraformBackendTypeRule()

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"resource.tf": test.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, test.Expected, runner.Issues)
		})
	}
}
