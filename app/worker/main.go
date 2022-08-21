package main

import (
	"OWTAssignment/app/adapters/cadence"
	"OWTAssignment/app/config"
	"OWTAssignment/app/server"
	"OWTAssignment/app/worker/workflows"
	"context"
	"go.uber.org/cadence/worker"
	"go.uber.org/zap"
)

var (
	appConfig     = config.AppConfig
	cadenceClient = cadence.CadenceAdapter
	logger        *zap.Logger
)

func init() {
	appConfig.LoadConfiguration()
	cadenceClient.Setup(&appConfig.Cadence)
	logger = appConfig.Logger
	server.Logger = *server.LogConfigRef
}

func startWorkers(h *cadence.Adapter, taskList string) {
	// Configure worker options.
	workerOptions := worker.Options{
		MetricsScope: h.Scope,
		Logger:       h.Logger,
	}

	cadenceWorker := worker.New(h.ServiceClient, h.Config.Domain, taskList, workerOptions)
	err := cadenceWorker.Start()
	if err != nil {
		h.Logger.Error("Failed to start workers.", zap.Error(err))
		panic("Failed to start workers")
	}
}

func startServer(c context.Context, cancelCtx func()) {
	if err := server.Server.Run("0.0.0.0:8080"); err != nil {
		logger.Error("Could not start server")
		logger.Panic("Cancelling server thread")
		cancelCtx()

	}

}

func main() {
	serverCtx := context.Background()
	serverCtx, cancelServerCtx := context.WithCancel(serverCtx)

	logger.Info("Configurations successfully loaded")
	logger.Info("Starting worker...")
	startWorkers(cadenceClient, workflows.TaskListName)

	logger.Info("Starting API server...")

	server.LoadRoutesAndMiddlewares()

	go startServer(serverCtx, cancelServerCtx)

	select {
	case <-serverCtx.Done():
		return

	}

}
