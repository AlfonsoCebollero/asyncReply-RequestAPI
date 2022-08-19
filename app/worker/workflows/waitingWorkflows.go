package workflows

import (
	"context"
	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
	"time"
)

/**
 * This is the hello world workflow sample.
 */

// ApplicationName is the task list for this sample
const TaskListName = "asyncAPI"

// This is registration process where you register all your workflows
// and activity function handlers.
func init() {
	workflow.Register(WaitingWorkflow)
	activity.RegisterWithOptions(basicActivity, activity.RegisterOptions{
		EnableAutoHeartbeat: true,
	})
}

var activityOptions = workflow.ActivityOptions{
	ScheduleToStartTimeout: time.Minute,
	StartToCloseTimeout:    time.Minute,
	HeartbeatTimeout:       time.Second * 20,
}

func basicActivity(ctx context.Context) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("basicActivity started")
	return "Sleep interval completed", nil
}

func WaitingWorkflow(ctx workflow.Context, sleepTime int) (string, error) {
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	logger := workflow.GetLogger(ctx)
	logger.Info("waiting workflow started")
	interval := time.Second * time.Duration(sleepTime)

	err := workflow.Sleep(ctx, interval)
	if err != nil {
		return "", err
	}

	var activityResult string
	err = workflow.ExecuteActivity(ctx, basicActivity).Get(ctx, &activityResult)
	if err != nil {
		logger.Error("Activity failed.", zap.Error(err))
		return "", err
	}
	logger.Info("Activity Completed", zap.String("Activity Result", activityResult))

	return "Completed!", nil
}
