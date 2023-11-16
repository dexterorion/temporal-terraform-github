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
	CreateRepositoryInput struct {
		Name        string
		Description string
	}

	CreateRepositoryOutput struct {
		URL string
	}
)

func CreateRepositoryWorkflow(ctx workflow.Context, input CreateRepositoryInput) (CreateRepositoryOutput, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour,
		HeartbeatTimeout:    time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    5 * time.Second,
			BackoffCoefficient: 1.3,
			MaximumInterval:    30 * time.Second,
		},
	})

	// Create repo
	var repo CreateRepositoryOutput
	err := workflow.ExecuteActivity(ctx, CreateRepositoryActivity, input).Get(ctx, &repo)
	if err != nil {

		return CreateRepositoryOutput{}, err
	}

	return repo, nil
}

func CreateRepositoryActivity(ctx context.Context, input CreateRepositoryInput) (CreateRepositoryOutput, error) {
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

	applyOutput, err := tfa.Apply(ctx, tfworkspace.ApplyInput{
		Vars: map[string]interface{}{
			"name":        input.Name,
			"description": input.Description,
		},
		Env: map[string]string{
			"GITHUB_TOKEN": cfg.GithubToken,
		},
	})
	if err != nil {
		return CreateRepositoryOutput{}, err
	}

	url, err := applyOutput.String("html_url")
	if err != nil {
		return CreateRepositoryOutput{}, err
	}

	return CreateRepositoryOutput{
		URL: url,
	}, nil
}
