package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	temporalClient "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	logAdapter "logur.dev/adapter/logrus"
	"logur.dev/logur"
	"net/http"
	"temporal/cmd/config"
	"temporal/workflows/test"
	"time"
)

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	//ctx := context.Background()

	confFile := flag.String("config", "config.yml", "The YAML configuration file.")
	flag.Parse()

	conf := config.Config{}
	if err := conf.Load(*confFile); err != nil {
		logger.WithError(err).WithField("config", *confFile).Fatal("failed to load configuration")
	}

	if err := conf.Validate(); err != nil {
		logger.WithError(err).WithField("config", *confFile).Fatal("invalid configuration")
	}

	logrusLevel, err := logrus.ParseLevel(conf.Log.Level)
	if err != nil {
		logger.WithError(err).Error("invalid log level provided")
		logrusLevel = logrus.InfoLevel
	}
	logger.SetLevel(logrusLevel)

	// ---

	// health endpoint
	logger.Info("spawning health endpoint")
	go func() {
		if err := http.ListenAndServe("0.0.0.0:8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})); err != nil {
			logger.WithError(err).Fatal("could not run health endpoint")
		}
	}()

	// temporal client
	logger.Info("creating temporal client")
	tClient, err := temporalClient.Dial(temporalClient.Options{
		HostPort:  fmt.Sprintf("%s:%d", conf.Temporal.Host, conf.Temporal.Port),
		Namespace: conf.Temporal.Namespace,
		Logger:    logur.LoggerToKV(logAdapter.New(logger)),
		ConnectionOptions: temporalClient.ConnectionOptions{
			TLS:                  nil,
			EnableKeepAliveCheck: true,
			KeepAliveTime:        time.Minute,
			KeepAliveTimeout:     time.Second * 10,
		},
	})

	if err != nil {
		logger.WithError(err).Fatal("could not dial temporal client")
	}

	defer tClient.Close()

	logger.Info("creating worker instance")

	tWorker := worker.New(tClient, conf.Temporal.Queue, worker.Options{
		MaxConcurrentActivityExecutionSize: int(conf.Temporal.MaxConcurrent),
		WorkerStopTimeout:                  time.Second * 30,
		LocalActivityWorkerOnly:            false,
		MaxHeartbeatThrottleInterval:       time.Second * 10,
		OnFatalError: func(e error) {
			logger.WithError(err).Error("fatal worker error encountered")
		},
		MaxConcurrentEagerActivityExecutionSize: 0,
	})

	logger.Info("registering workflows")

	if err := RegisterWorkflows(logger, conf, tWorker); err != nil {
		logger.WithError(err).Fatal("could not register workflows")
	}

	logger.Info("started worker")

	go func() {
		time.Sleep(time.Second * 5)

		logger.Info("running workflow")

		run, err := tClient.ExecuteWorkflow(context.Background(), temporalClient.StartWorkflowOptions{
			ID:                                       "",
			TaskQueue:                                conf.Temporal.Queue,
			WorkflowExecutionTimeout:                 0,
			WorkflowRunTimeout:                       0,
			WorkflowTaskTimeout:                      0,
			WorkflowIDReusePolicy:                    0,
			WorkflowExecutionErrorWhenAlreadyStarted: false,
			RetryPolicy:                              nil,
			CronSchedule:                             "",
			Memo:                                     nil,
			SearchAttributes:                         nil,
		}, test.Workflow, "myArgument")
		if err != nil {
			logger.WithError(err).Fatal("could not run workflow")
		}

		wfID := run.GetID()
		logger.WithField("id", wfID).Info("workflow ran")
	}()

	if err := tWorker.Run(worker.InterruptCh()); err != nil {
		logger.WithError(err).Fatal("worker returned error")
	}

	logger.Info("stopped worker")
}
