package test

import (
	"go.temporal.io/sdk/workflow"
)

func Workflow(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)

	logger.Error("test workflow fired")

	return nil
}
