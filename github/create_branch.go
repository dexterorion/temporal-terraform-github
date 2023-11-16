package github

import (
	"context"
	"fmt"
	"time"

	"github.com/dexterorion/meetupgo/config/env"
	"github.com/dexterorion/meetupgo/orchestrator/terraform"
	"github.com/dexterorion/meetupgo/orchestrator/tfactivity"
	"github.com/dexterorion/meetupgo/orchestrator/tfexec"
	"github.com/dexterorion/meetupgo/orchestrator/tfworkspace"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type (
	CreateBranchInput struct {
		Repository string
		Branch     string
	}
)

func CreateBranchWorkflow(ctx workflow.Context, input CreateBranchInput) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour,
		HeartbeatTimeout:    time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    5 * time.Second,
			BackoffCoefficient: 1.3,
			MaximumInterval:    30 * time.Second,
		},
	})

	// Create branch
	err := workflow.ExecuteActivity(ctx, CreateBranchActivity, input).Get(ctx, nil)
	if err != nil {

		return err
	}

	return nil
}

func CreateBranchActivity(ctx context.Context, input CreateBranchInput) error {
	cfg := env.MustGetConfig()

	// Terraform wrapper
	tfa := tfactivity.New(tfworkspace.Config{
		TerraformPath: "github/branch",
		TerraformFS:   terraform.GH,
		Backend: tfexec.BackendConfig{
			Bucket: cfg.TfState.Bucket,
			Key:    fmt.Sprintf("repository-branch-%s.tfstate", input.Branch),
		},
	})

	_, err := tfa.Apply(ctx, tfworkspace.ApplyInput{
		Vars: map[string]interface{}{
			"repository": input.Repository,
			"branch":     input.Branch,
		},
		Env: map[string]string{
			"GITHUB_TOKEN": cfg.GithubToken,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
