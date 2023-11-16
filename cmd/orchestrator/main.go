package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/dexterorion/meetupgo/config/env"
	"github.com/dexterorion/meetupgo/worker/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
)

var (
	workflow  string
	taskqueue string
	input     string
)

func init() {
	flag.StringVar(&workflow, "workflow", "", "Workflow name")
	flag.StringVar(&taskqueue, "taskqueue", "", "Task queue name")
	flag.StringVar(&input, "input", "", "Input JSON file")

	flag.Parse()
}

func main() {
	cfg := env.MustGetConfig()

	serviceClient, err := client.NewLazyClient(client.Options{
		Namespace: cfg.TemporalNamespace,
		HostPort:  cfg.TemporalHostPort,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx := context.Background()

	in, err := workflows.MapInputToWorkflow(workflow, input)
	if err != nil {
		log.Fatal(err.Error(), "wfname", workflow)
	}

	run, err := serviceClient.ExecuteWorkflow(ctx, client.StartWorkflowOptions{
		TaskQueue:          taskqueue,
		WorkflowRunTimeout: time.Second * 300,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Second * 100,
			MaximumAttempts:    3,
		},
	}, workflow, in)

	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("Started workflow WFID: %s, RID: %s", run.GetID(), run.GetRunID())
}
