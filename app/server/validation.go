package server

import (
	"OWTAssignment/app/server/entities"
	"OWTAssignment/app/worker/workflows"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

// workflowCreationRequestValidation validates if the request params are correct,
// it also checks if a task list has been selected. Returns task list name, task duration and error.
func workflowCreationRequestValidation(c *gin.Context, wfs map[string]string) (string, string, int, error) {
	queryParams := c.Request.URL.Query()
	Logger.Debug("Checking if a taskList has been given")
	queryParamList, exists := queryParams["taskList"]
	var taskListName string
	if exists {
		taskListName = queryParamList[0]
	} else {
		taskListName = workflows.TaskListName
	}
	Logger.Info("Task list set.", zap.String("TaskListName", taskListName))

	Logger.Debug("Checking duration param.")
	duration, exists := queryParams["duration"]
	if !exists {
		Logger.Error("Duration has not been passed intoi the request")
		return "", "", 0, entities.ValidationError
	}
	toIntDuration, err := strconv.Atoi(duration[0])
	if err != nil {
		Logger.Error("Duration must be an integer", zap.String("Duration", duration[0]))
		return "", "", 0, entities.ValidationError
	}

	if toIntDuration < 30 {
		Logger.Error("Duration must be higher than or equal to 30", zap.Int("Duration", toIntDuration))
		return "", "", 0, entities.ValidationError
	}

	Logger.Debug("Checking name param.")

	name := c.Param("name")
	wfNameValue, exists := wfs[strings.ToLower(name)]
	if !exists {
		Logger.Error(
			"This workflow is not available or does not exist",
			zap.Any("Available Options", wfs))
		return "", name, 0, entities.BadWorkflow
	}
	return taskListName, wfNameValue, toIntDuration, nil
}
