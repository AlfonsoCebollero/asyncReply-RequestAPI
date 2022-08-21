package useCases

import (
	"github.com/AlfonsoCebollero/asyncReply-RequestAPI/app/config"
	"github.com/AlfonsoCebollero/asyncReply-RequestAPI/app/useCases/interfaces"
	"go.uber.org/zap"
)

func CreateWorkflow(manager interfaces.IWorkflowManager, taskListName, duration, nameParam string) (string, string, error) {
	logger := config.AppConfig.Logger
	wfs := config.AppConfig.Cadence.Workflows
	logger.Info("Validating request...")
	taskListName, name, interval, err := WorkflowCreationRequestValidation(taskListName, duration, nameParam, wfs)
	if err != nil {
		return "", "", err
	}
	logger.Info("Request validated!")

	return manager.StartWorkflow(name, taskListName, interval)
}

func RetrieveWorkflowStatus(manager interfaces.IWorkflowManager, workflowID string) (string, error) {
	logger := config.AppConfig.Logger
	logger.Info("Retrieving workflow closeStatus...", zap.String("workflow-id", workflowID))

	status, pendingDecision, err := manager.RetrieveWFStatus(workflowID)

	if err != nil {
		return "", err
	}

	switch status {
	case -1:
		if pendingDecision {
			return "PENDING", nil
		} else {
			return "RUNNING", nil
		}
	case 0: // WorkflowExecutionCloseStatusCompleted:
		return "COMPLETED", nil
	case 1: // WorkflowExecutionCloseStatusCanceled:
		return "FAILED", nil
	case 2: // WorkflowExecutionCloseStatusFailed:
		return "CANCELED", nil
	case 3: // WorkflowExecutionCloseStatusTerminated:
		return "TERMINATED", nil
	case 4: // WorkflowExecutionCloseStatusContinuedAsNew
		return "CONTINUED", nil
	default: // WorkflowExecutionCloseStatusTimedOut
		return "TIMED OUT", nil

	}
}
