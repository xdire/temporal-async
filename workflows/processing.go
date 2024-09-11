package workflows

import (
	"github.com/xdire/temporal-async/activities"
	"github.com/xdire/temporal-async/messaging"
	"github.com/xdire/temporal-async/util"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"time"
)

func StartProcess(ctx workflow.Context, params *messaging.CreateProcess) error {
	logger := workflow.GetLogger(ctx)
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 20,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Minute,
			BackoffCoefficient: 1.5,
			MaximumAttempts:    5,
			MaximumInterval:    time.Hour,
			NonRetryableErrorTypes: []string{
				util.ProcessingErrorType.Error(),
			},
		},
	}

	// Get ID of the workflow for bind process to the Workflow signaling
	params.WorkflowID = workflow.GetInfo(ctx).WorkflowExecution.ID

	// start workflow with activity options
	ctx = workflow.WithActivityOptions(ctx, options)

	proc := &messaging.Process{}
	err := workflow.ExecuteActivity(ctx, activities.CreateProcess, params).Get(ctx, proc)
	if err != nil {
		logger.Error("failed to create process", "error", err)
		return err
	}

	var stageUpdateDone bool
	var processStageUpdateName string
	var confirmProcessCompletionName string

	stageUpdateSelector := workflow.NewSelector(ctx)
	confirmUpdateSelector := workflow.NewSelector(ctx)

	stageUpdateSignals := workflow.GetSignalChannel(ctx, "processStageUpdates")
	processUpdateSignals := workflow.GetSignalChannel(ctx, "processUpdates")

	stageUpdateSelector.AddReceive(stageUpdateSignals, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, &processStageUpdateName)
		// Just log notify
		if processStageUpdateName == "finalize" {
			logger.Info("stage update selector completed")
			stageUpdateDone = true
		}
	})

	confirmUpdateSelector.AddReceive(processUpdateSignals, func(c workflow.ReceiveChannel, more bool) {
		confirmDone := c.Receive(ctx, &confirmProcessCompletionName)
		// Just log notify
		if confirmDone {
			logger.Info("confirm update selector completed")
		}
	})

	// Receive events until signal completion is received
	for !stageUpdateDone {
		stageUpdateSelector.Select(ctx)
		errStage := workflow.ExecuteActivity(ctx, activities.UpdateProcessStage, proc, processStageUpdateName).Get(ctx, proc)
		if errStage != nil {
			logger.Error("unable update stage in process", "Error", errStage.Error())
		}
		logger.Info("stage updated", "stage", processStageUpdateName)
	}

	err = workflow.ExecuteActivity(ctx, activities.UpdateProcessCost, proc, 99.9).Get(ctx, proc)
	if err != nil {
		return err
	}

	confirmUpdateSelector.Select(ctx)

	err = workflow.ExecuteActivity(ctx, activities.UpdateProcessStage, proc, "confirmed").Get(ctx, proc)
	if err != nil {
		logger.Error("unable update stage in process", "Error", err.Error())
		return err
	}

	logger.Info("Process created", "ID", proc.ID, "UUID", proc.UUID)
	return nil
}
