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
	DestroyRepositoryInput struct {
		WorkflowID  string
		Name        string
		Description string
	}

	DestroyRepositoryOutput struct {
		URL string
	}
)

func DestroyRepositoryWorkflow(ctx workflow.Context, input DestroyRepositoryInput) (DestroyRepositoryOutput, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour,
		HeartbeatTimeout:    time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    5 * time.Second,
			BackoffCoefficient: 1.3,
			MaximumInterval:    30 * time.Second,
		},
	})

	// Destroy repo
	var repo DestroyRepositoryOutput
	err := workflow.ExecuteActivity(ctx, DestroyRepositoryActivity, input).Get(ctx, &repo)
	if err != nil {
		return DestroyRepositoryOutput{}, err
	}

	return repo, nil
}

func DestroyRepositoryActivity(ctx context.Context, input DestroyRepositoryInput) (DestroyRepositoryOutput, error) {
	cfg := env.MustGetConfig()

	// Terraform wrapper
	tfa := tfactivity.New(tfworkspace.Config{
		TerraformPath: "github/repository",
		TerraformFS:   terraform.GH,
		Backend: tfexec.BackendConfig{
			Bucket: cfg.TfState.Bucket,
			Key:    fmt.Sprintf("repository-%s.tfstate", input.Name),
		},
	})

	err := tfa.Destroy(ctx, tfworkspace.DestroyInput{
		Vars: map[string]interface{}{
			"name":        input.Name,
			"description": input.Description,
		},
		Env: map[string]string{
			"GITHUB_TOKEN": cfg.GithubToken,
		},
	})

	return DestroyRepositoryOutput{}, err
}
