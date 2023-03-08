package test

import (
	"fmt"
	"go.temporal.io/sdk/workflow"
	"time"
)

func Workflow(ctx workflow.Context, arg string) error {
	logger := workflow.GetLogger(ctx)

	logger.Info("test workflow fired", arg)

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:              "",
		ScheduleToStartTimeout: 0,
		StartToCloseTimeout:    time.Second,
		HeartbeatTimeout:       0,
		WaitForCancellation:    false,
		ActivityID:             "",
		RetryPolicy:            nil,
		DisableEagerExecution:  false,
	})

	run := workflow.ExecuteActivity(ctx, Activity, "newARGUMENT")
	if err := run.Get(ctx, nil); err != nil {
		return fmt.Errorf("could not run activity: %v", err)
	}

	return nil
}
