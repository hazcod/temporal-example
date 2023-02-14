package test

import (
	"context"
	"go.temporal.io/sdk/activity"
)

func Activity(ctx context.Context) error {
	logger := activity.GetLogger(ctx)

	logger.Debug("test activity fired")

	return nil
}
