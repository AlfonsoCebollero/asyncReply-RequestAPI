package server

import (
	"OWTAssignment/app/adapters/cadence"
	"OWTAssignment/app/config"
	_ "OWTAssignment/app/docs"
	"OWTAssignment/app/server/entities"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/cadence/.gen/go/shared"
	"go.uber.org/cadence/client"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const (
	wfPrefix             = "OWTAssignment/app/worker/workflows.%s"
	workflowAPIBaseRoute = "/api/v1/workflow/%s"
)

var (
	cadenceAdapter = cadence.CadenceAdapter

	Server       = gin.Default()
	Logger       *zap.Logger
	LogConfigRef **zap.Logger
)

func init() {
	LogConfigRef = &config.AppConfig.Logger
}

// createWorkflow Endpoint in charge of the creation of new workflows
// @Summary Creates new workflow if params are valid
// @Schemes
// @Description Receives a path param name, a query duration and, optionally, a query param task list name. With these a new workflow of the type that matches the name is created with the requested duration
// @Tags        WFCreation
// @Produce     json
// @Success     220          {object} entities.APIResponse{response=entities.WFCreationSuccessfulResponse}
// @Failure     400          {object} entities.APIError{}
// @Failure     500          {object} entities.APIError{}
// @Param       workflowName path     string true  "Workflow Name"
// @Param       duration     query    int    true  "Duration"
// @Param       taskListName query    string false "Task list name"
// @Router      /workflow/{workflowName} [post]
func createWorkflow(c *gin.Context) {
	Logger.Info("Validating request...")
	wfs := config.AppConfig.Cadence.Workflows
	taskListName, name, duration, err := workflowCreationRequestValidation(c, wfs)
	if err != nil {
		switch err {
		case entities.ValidationError:
			message := "A duration >= 30 must be supplied"
			c.JSON(http.StatusBadRequest, getValidationErrorResponse(message))
			return
		case entities.BadWorkflow:
			message := fmt.Sprintf("%s is not a valid workflow name, available options: %v",
				name, wfs)
			c.JSON(http.StatusBadRequest, getValidationErrorResponse(message))
			return
		}

	}
	Logger.Info("Request validated!")

	wo := client.StartWorkflowOptions{
		TaskList:                     taskListName,
		ExecutionStartToCloseTimeout: time.Hour * 12,
	}

	execution, err := cadenceAdapter.CadenceClient.StartWorkflow(
		context.Background(),
		wo,
		fmt.Sprintf(wfPrefix, name),
		duration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, getWorkflowCreationError(err))
		return
	}
	Logger.Info("Started work flow!",
		zap.String("WorkflowId", execution.ID),
		zap.String("RunId", execution.RunID))

	c.JSON(http.StatusAccepted, entities.APIResponse{
		Timestamp: time.Now().String(),
		Response: entities.WFCreationSuccessfulResponse{
			WorkflowID:     execution.ID,
			WorkflowRunID:  execution.RunID,
			WorkflowStatus: "Created",
			Message:        "Workflow was created successfully",
			StatusCode:     "ACCEPTED",
			Href:           config.AppConfig.Cadence.ServerBaseUrl + fmt.Sprintf(workflowAPIBaseRoute, execution.ID),
		},
	})

}

// retrieveStatus Retrieves the status of the workflow which ID is passed as a path param.
// @Summary Retrieves status of the workflow associated to the given ID
// @Schemes
// @Description Responds with the status of the searched workflow
// @Tags        WFStatusRetrieving
// @Produce     json
// @Success     200          {object} entities.APIResponse{response=entities.WFRetrievingSuccessfulResponse}
// @Failure     400          {object} entities.APIError{}
// @Failure     500          {object} entities.APIError{}
// @Param       workflowID path     string true  "Workflow ID"
// @Router      /workflow/status/{workflowID} [get]
func retrieveStatus(c *gin.Context) {
	workflowID := c.Param("id")
	Logger.Info("Retrieving workflow closeStatus...", zap.String("workflow-id", workflowID))

	workflowExecution, err := cadenceAdapter.CadenceClient.DescribeWorkflowExecution(context.Background(), workflowID, "")
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	closeStatus := workflowExecution.WorkflowExecutionInfo.CloseStatus

	var status string

	if closeStatus == nil {
		if workflowExecution.GetPendingDecision() != nil {
			status = "PENDING"
		} else {
			status = "RUNNING"
		}
		c.JSON(http.StatusOK, getWorkflowStatusResponse(status))
		return
	}

	switch *workflowExecution.WorkflowExecutionInfo.CloseStatus {
	case shared.WorkflowExecutionCloseStatusCompleted:
		status = "COMPLETED"
	case shared.WorkflowExecutionCloseStatusCanceled:
		status = "CANCELED"
	case shared.WorkflowExecutionCloseStatusFailed:
		status = "FAILED"
	case shared.WorkflowExecutionCloseStatusTerminated:
		status = "TERMINATED"
	default:
		status = "COMPLETED"

	}

	c.JSON(http.StatusOK, getWorkflowStatusResponse(status))
	return

}

// LoadRoutesAndMiddlewares sets all routes and middlewares selected to run within the API
func LoadRoutesAndMiddlewares() {
	Server.POST(fmt.Sprintf(workflowAPIBaseRoute, ":name"), createWorkflow)
	Server.GET(fmt.Sprintf(workflowAPIBaseRoute, "status/:id"), retrieveStatus)

	Server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	Logger.Info("Routes and middlewares successfully loaded")
}

// getValidationErrorResponse conforms the response when a validation error occurs
func getValidationErrorResponse(message string) entities.APIError {
	return entities.APIError{
		Timestamp: time.Now().String(),
		Response: entities.ErrorResponse{
			StatusCode: "BAD REQUEST",
			Message:    message,
		},
	}
}

// getWorkflowCreationError conforms the response when a workflow creation error occurs
func getWorkflowCreationError(err error) entities.APIError {
	return entities.APIError{
		Timestamp: time.Now().String(),
		Response: entities.ErrorResponse{
			StatusCode: "BAD REQUEST",
			Message:    fmt.Sprintf("There was an error while creating the workflow: %v", err),
		},
	}
}

func getWorkflowStatusResponse(status string) entities.APIResponse {
	return entities.APIResponse{
		Timestamp: time.Now().String(),
		Response: entities.WFRetrievingSuccessfulResponse{
			WorkflowStatus: status,
			StatusCode:     "OK",
		},
	}
}
