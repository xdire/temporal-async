package activities

import (
	"context"
	"fmt"
	"github.com/lib/pq"
	"github.com/xdire/temporal-async/messaging"
	"github.com/xdire/temporal-async/util"
	"go.temporal.io/sdk/temporal"
	"time"
)

func CreateProcess(ctx context.Context, process *messaging.CreateProcess) (*messaging.Process, error) {
	if process.WorkflowID == "" {
		return nil, temporal.
			NewNonRetryableApplicationError(
				"service error",
				util.ProcessingErrorType.Error(),
				fmt.Errorf("no sync-id parameter presented"))
	}

	dbc := &util.DB{}
	err := dbc.Connect()
	if err != nil {
		return nil, err
	}

	proc := &messaging.Process{
		UUID:    process.WorkflowID,
		Name:    process.Name,
		Desc:    process.Desc,
		Stage:   "init",
		Cost:    0,
		Parts:   pq.StringArray{},
		Created: time.Now().UTC(),
	}

	tx := dbc.Conn().Create(proc)
	if tx.Error != nil {
		return nil, temporal.NewNonRetryableApplicationError("service error", util.ProcessingErrorType.Error(), err)
	}

	return proc, nil
}

func UpdateProcessStage(ctx context.Context, process *messaging.Process, stage string) (*messaging.Process, error) {
	if process.UUID == "" {
		return nil, temporal.
			NewNonRetryableApplicationError(
				"service error",
				util.ProcessingErrorType.Error(),
				fmt.Errorf("no sync-id parameter presented"))
	}

	dbc := &util.DB{}
	err := dbc.Connect()
	if err != nil {
		return nil, err
	}

	proc := &messaging.Process{}
	records := dbc.Conn().Where("uuid = ?", process.UUID).First(proc)
	if records.Error != nil {
		return nil, temporal.NewNonRetryableApplicationError(fmt.Errorf("service error, %w", err).Error(), util.ProcessingErrorType.Error(), nil)
	}

	proc.Stage = stage
	dbc.Conn().Model(proc).Where("uuid = ?", process.UUID).Update("stage", stage)
	if records.Error != nil {
		return nil, records.Error
	}

	return proc, nil
}

func UpdateProcessCost(ctx context.Context, process *messaging.Process, cost float32) (*messaging.Process, error) {
	if process.UUID == "" {
		return nil, temporal.
			NewNonRetryableApplicationError(
				"service error",
				util.ProcessingErrorType.Error(),
				fmt.Errorf("no sync-id parameter presented"))
	}

	dbc := &util.DB{}
	err := dbc.Connect()
	if err != nil {
		return nil, err
	}

	proc := &messaging.Process{}
	records := dbc.Conn().Where("uuid = ?", process.UUID).First(proc)
	if records.Error != nil {
		return nil, temporal.NewNonRetryableApplicationError(fmt.Errorf("service error, %w", err).Error(), util.ProcessingErrorType.Error(), nil)
	}

	proc.Cost = cost
	dbc.Conn().Model(proc).Where("uuid = ?", process.UUID).Update("cost", cost)
	if records.Error != nil {
		return nil, records.Error
	}

	return proc, nil
}
