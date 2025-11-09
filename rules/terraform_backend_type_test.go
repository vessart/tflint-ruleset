package rules

import (
	"testing"

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
