package wfManager

import (
	"OWTAssignment/app/config"
	"context"
	"fmt"
	"go.uber.org/cadence/client"
	"go.uber.org/zap"
	"time"
)

type CadenceWFManager struct {
	CadenceClient *client.Client
	WFPrefix      string
}

func (c *CadenceWFManager) StartWorkflow(name string, taskList string, duration int) (string, string, error) {
	logger := config.AppConfig.Logger

	wo := client.StartWorkflowOptions{
		TaskList:                     taskList,
		ExecutionStartToCloseTimeout: time.Hour * 12,
	}
	cClient := *c.CadenceClient

	execution, err := cClient.StartWorkflow(
		context.Background(),
		wo,
		fmt.Sprintf(c.WFPrefix, name),
		duration)
	if err != nil {
		return "", "", err
	}

	logger.Info("Started work flow!",
		zap.String("WorkflowId", execution.ID),
		zap.String("RunId", execution.RunID))

	return execution.ID, execution.RunID, nil
}

func (c *CadenceWFManager) RetrieveWFStatus(workflowID string) (int, bool, error) {
	cClient := *c.CadenceClient
	workflowExecution, err := cClient.DescribeWorkflowExecution(context.Background(), workflowID, "")
	if err != nil {
		return 0, false, err
	}

	closeStatus := workflowExecution.WorkflowExecutionInfo.CloseStatus
	pendingDecision := workflowExecution.GetPendingDecision()

	if closeStatus == nil {
		if pendingDecision != nil {
			return -1, true, nil
		}
		return -1, false, nil
	}

	return int(*workflowExecution.WorkflowExecutionInfo.CloseStatus), false, nil

}
