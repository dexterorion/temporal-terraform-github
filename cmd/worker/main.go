package main

import (
	"log"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/dexterorion/meetupgo/config/env"
	"github.com/dexterorion/meetupgo/worker/workflows"
)

func main() {
	cfg := env.MustGetConfig()

	serviceClient, err := client.NewLazyClient(client.Options{
		Namespace: cfg.TemporalNamespace,
		HostPort:  cfg.TemporalHostPort,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	temporalWorker := worker.New(serviceClient, cfg.TemporalTaskQueue, worker.Options{
		WorkerStopTimeout: 30 * time.Second,
	})

	log.Print("registering workflows")
	workflows.Register(temporalWorker)

	if err := temporalWorker.Run(worker.InterruptCh()); err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
