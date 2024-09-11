package main

import (
	"github.com/xdire/temporal-async/activities"
	"github.com/xdire/temporal-async/util"
	"github.com/xdire/temporal-async/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	log := util.NewZeroLogForName("worker", "x1", "info")
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create temporal client")
	}
	defer c.Close()

	w := worker.New(c, "saga-process", worker.Options{
		MaxConcurrentActivityExecutionSize: 10,
	})

	w.RegisterWorkflow(workflows.StartProcess)
	w.RegisterActivity(activities.CreateProcess)
	w.RegisterActivity(activities.UpdateProcessStage)
	w.RegisterActivity(activities.UpdateProcessCost)

	// start the worker
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatal().Err(err).Msg("unable to start Worker")
	}
}
