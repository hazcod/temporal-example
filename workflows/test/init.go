package test

import (
	"go.temporal.io/api/enums/v1"
	temporalClient "go.temporal.io/sdk/client"
	"time"
)

const (
	workflowName = "test"
)

func GetWorkflowOptions(queueName string) temporalClient.StartWorkflowOptions {
	return temporalClient.StartWorkflowOptions{
		ID:                                       workflowName,
		TaskQueue:                                queueName,
		WorkflowExecutionTimeout:                 time.Second * 10,
		WorkflowRunTimeout:                       time.Second * 5,
		WorkflowTaskTimeout:                      time.Second * 1,
		WorkflowIDReusePolicy:                    enums.WORKFLOW_ID_REUSE_POLICY_TERMINATE_IF_RUNNING,
		WorkflowExecutionErrorWhenAlreadyStarted: false,
		RetryPolicy:                              nil,
		CronSchedule:                             "",
		Memo:                                     nil,
		SearchAttributes:                         nil,
	}
}
