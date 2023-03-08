package test

import (
	"context"
	"go.temporal.io/sdk/activity"
)

func Activity(ctx context.Context, arg string) error {
	logger := activity.GetLogger(ctx)

	logger.Info("test activity fired", arg)

	return nil
}
