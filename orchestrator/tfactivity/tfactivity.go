package tfactivity

import (
	"context"
	"time"

	"go.temporal.io/sdk/activity"

	"github.com/dexterorion/meetupgo/orchestrator/heartbeat"
	"github.com/dexterorion/meetupgo/orchestrator/tfworkspace"
)

type Activity struct {
	config tfworkspace.Config
}

func New(wsConfig tfworkspace.Config) *Activity {
	return &Activity{config: wsConfig}
}

func (a *Activity) Apply(ctx context.Context, input tfworkspace.ApplyInput) (tfworkspace.ApplyOutput, error) {
	logger := activity.GetLogger(ctx)
	ctx, cancel := heartbeat.Begin(ctx, 30*time.Second)
	defer cancel()

	logger.Info("terraform activity apply", "TerraformPath", a.config.TerraformPath,
		"StateBucket", a.config.Backend.Bucket, "StateKey", a.config.Backend.Key)

	// Blocking call that returns when terraform exits
	return tfworkspace.New(a.config).Apply(ctx, input)
}

func (a *Activity) Destroy(ctx context.Context, input tfworkspace.DestroyInput) error {
	logger := activity.GetLogger(ctx)
	ctx, cancel := heartbeat.Begin(ctx, 30*time.Second)
	defer cancel()

	logger.Info("terraform activity destroy", "TerraformPath", a.config.TerraformPath,
		"StateBucket", a.config.Backend.Bucket, "StateKey", a.config.Backend.Key)

	// Blocking call that returns when terraform exits
	return tfworkspace.New(a.config).Destroy(ctx, input)
}
