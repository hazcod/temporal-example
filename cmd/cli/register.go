package main

import (
	"github.com/sirupsen/logrus"
	"go.temporal.io/sdk/worker"
	"temporal/cmd/config"
	"temporal/workflows/test"
)

func RegisterWorkflows(_ *logrus.Logger, _ config.Config, worker worker.Worker) error {
	// test
	worker.RegisterActivity(test.Activity)
	worker.RegisterWorkflow(test.Workflow)

	return nil
}
