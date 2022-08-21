package server

import (
	"OWTAssignment/app/adapters/cadence"
	"OWTAssignment/app/config"
	_ "OWTAssignment/app/docs"
	"OWTAssignment/app/server/entities"
	"OWTAssignment/app/useCases"
	"OWTAssignment/app/useCases/wfManager"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	cadenceWFManager = wfManager.CadenceWFManager{
		CadenceClient: &cadenceAdapter.CadenceClient,
		WFPrefix:      wfPrefix,
	}
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
	queryParams := c.Request.URL.Query()
	duration := queryParams.Get("duration")
	taskListName := queryParams.Get("taskList")

	executionID, executionRunID, err := useCases.CreateWorkflow(
		&cadenceWFManager,
		taskListName,
		duration,
		c.Param("name"),
	)

	switch err {
	case nil:
		c.JSON(http.StatusAccepted, entities.APIResponse{
			Timestamp: time.Now().String(),
			Response: entities.WFCreationSuccessfulResponse{
				WorkflowID:     executionID,
				WorkflowRunID:  executionRunID,
				WorkflowStatus: "Created",
				Message:        "Workflow was created successfully",
				StatusCode:     "ACCEPTED",
				Href: config.AppConfig.Cadence.ServerBaseUrl +
					fmt.Sprintf(workflowAPIBaseRoute, "status/"+executionID),
			},
		})
		return
	case entities.ValidationError:
		message := "A duration >= 30 must be supplied"
		c.JSON(http.StatusBadRequest, getValidationErrorResponse(message))
		return
	case entities.BadWorkflow:
		message := fmt.Sprintf("%s is not a valid workflow name, available options: %v",
			c.Param("name"), config.AppConfig.Cadence.Workflows)
		c.JSON(http.StatusInternalServerError, getValidationErrorResponse(message))
		return
	default:
		c.JSON(http.StatusInternalServerError, getWorkflowCreationError(err))
		return

	}

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
	status, err := useCases.RetrieveWorkflowStatus(&cadenceWFManager, workflowID)

	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
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

// getWorkflowStatusResponse conforms the response when a workflow retrieving operation completes successfully
func getWorkflowStatusResponse(status string) entities.APIResponse {
	return entities.APIResponse{
		Timestamp: time.Now().String(),
		Response: entities.WFRetrievingSuccessfulResponse{
			WorkflowStatus: status,
			StatusCode:     "OK",
		},
	}
}
