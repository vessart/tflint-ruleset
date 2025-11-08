package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformBackendTypeRule checks whether ...
type TerraformBackendTypeRule struct {
	tflint.DefaultRule
}

// NewTerraformBackendTypeRule returns a new rule
func NewTerraformBackendTypeRule() *TerraformBackendTypeRule {
	return &TerraformBackendTypeRule{}
}

// Name returns the rule name
func (r *TerraformBackendTypeRule) Name() string {
	return "terraform_backend_type"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformBackendTypeRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *TerraformBackendTypeRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *TerraformBackendTypeRule) Link() string {
	return ""
}

// Check checks whether ...
func (r *TerraformBackendTypeRule) Check(runner tflint.Runner) error {
	content, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: "terraform",
				Body: &hclext.BodySchema{
					Blocks: []hclext.BlockSchema{
						{
							Type:       "backend",
							LabelNames: []string{"type"},
							Body: &hclext.BodySchema{
								Attributes: []hclext.AttributeSchema{
									{Name: "address"},
									{Name: "lock_method"},
									{Name: "unlock_method"},
								},
							},
						},
					},
				},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, terraform := range content.Blocks {
		for _, backend := range terraform.Body.Blocks {
			backendType := backend.Labels[0]
			if backendType != "http" {
				return runner.EmitIssue(
					r,
					fmt.Sprintf("backend type must be 'http', but found '%s'", backendType),
					backend.DefRange,
				)
			}

			var lockMethod, unlockMethod string

			// Check lock_method
			lockMethodAttr := backend.Body.Attributes["lock_method"]
			if lockMethodAttr == nil {
				return runner.EmitIssue(
					r,
					`"lock_method" attribute is required`,
					backend.DefRange,
				)
			}
			if diag := runner.EvaluateExpr(lockMethodAttr.Expr, &lockMethod, nil); diag != nil {
				return diag
			}
			if lockMethod != "POST" {
				return runner.EmitIssue(
					r,
					`"lock_method" must be "POST"`,
					lockMethodAttr.Range,
				)
			}

			// Check unlock_method
			unlockMethodAttr := backend.Body.Attributes["unlock_method"]
			if unlockMethodAttr == nil {
				return runner.EmitIssue(
					r,
					`"unlock_method" attribute is required`,
					backend.DefRange,
				)
			}
			if diag := runner.EvaluateExpr(unlockMethodAttr.Expr, &unlockMethod, nil); diag != nil {
				return diag
			}
			if unlockMethod != "DELETE" {
				return runner.EmitIssue(
					r,
					`"unlock_method" must be "DELETE"`,
					unlockMethodAttr.Range,
				)
			}
		}
	}

	return nil
}
