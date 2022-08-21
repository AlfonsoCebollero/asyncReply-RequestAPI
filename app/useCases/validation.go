package useCases

import (
	"github.com/AlfonsoCebollero/asyncReply-RequestAPI/app/config"
	"github.com/AlfonsoCebollero/asyncReply-RequestAPI/app/server/entities"
	"github.com/AlfonsoCebollero/asyncReply-RequestAPI/app/worker/workflows"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

// WorkflowCreationRequestValidation validates if the request params are correct,
// it also checks if a task list has been selected. Returns task list name, task duration and error.
func WorkflowCreationRequestValidation(taskListName, duration, name string, wfs map[string]string) (string, string, int, error) {
	logger := config.AppConfig.Logger
	logger.Debug("Checking if a taskList has been given")
	if taskListName == "" {
		taskListName = workflows.TaskListName
	}

	logger.Info("Task list set.", zap.String("TaskListName", taskListName))

	logger.Debug("Checking duration param.")
	if duration == "" {
		logger.Error("Duration has not been passed into the request")
		return "", "", 0, entities.ValidationError
	}
	toIntDuration, err := strconv.Atoi(duration)
	if err != nil {
		logger.Error("Duration must be an integer", zap.String("Duration", duration))
		return "", "", 0, entities.ValidationError
	}

	if toIntDuration < 30 {
		logger.Error("Duration must be higher than or equal to 30", zap.Int("Duration", toIntDuration))
		return "", "", 0, entities.ValidationError
	}

	logger.Debug("Checking name param.")

	wfNameValue, exists := wfs[strings.ToLower(name)]
	if !exists {
		logger.Error(
			"This workflow is not available or does not exist",
			zap.Any("Available Options", wfs))
		return "", name, 0, entities.BadWorkflow
	}
	return taskListName, wfNameValue, toIntDuration, nil
}
