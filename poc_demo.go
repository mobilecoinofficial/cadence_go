package cadence_go

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bxcodec/faker/v3"
	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

type POCDemoWorkflowResult struct {
	//EndTime time.Time
	Output string
}

func POCDemoActivity1(ctx context.Context, input string) (string, error) {
	workflowID := activity.GetInfo(ctx).WorkflowExecution.ID

	sentence := faker.Sentence()

	msg := fmt.Sprintf("POCDemo activity1 is running for workflow: %s, sentence: %s\n", workflowID, sentence)
	activity.GetLogger(ctx).Info(
		msg,
		zap.String("input", input),
	)

	activitySummary := fmt.Sprintf("POCDemoActivity1, sentence: %s\n", sentence)

	return activitySummary, nil
}

func POCDemoWorkflow(ctx workflow.Context) (*POCDemoWorkflowResult, error) {
	log.Println("POCDemoWorkflow is triggered...")

	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute * 3,
	}

	nctx := workflow.WithActivityOptions(ctx, ao)

	log.Println("Waiting for Completion...")
	// if workflow.HasLastCompletionResult(ctx) {
	// 	log.Println("Found last completion result...")
	// 	var lastResult POCDemoWorkflowResult
	// 	err := workflow.GetLastCompletionResult(ctx, &lastResult)
	// 	if err != nil {
	// 		log.Println("Error getting last completion result...")
	// 		return nil, err
	// 	}
	// 	log.Println("Last completion result: ", lastResult)
	// 	//startTime = lastResult.EndTime
	// }

	log.Println("Starting workflow...")
	//endTime := workflow.Now(ctx)

	username := faker.Username()

	err := workflow.ExecuteActivity(
		nctx,
		POCDemoActivity1,
		username,
	).Get(ctx, nil)
	if err != nil {
		workflow.GetLogger(ctx).Error("POCDemoActivity1 failed.", zap.Error(err))
		return nil, err
	}

	return &POCDemoWorkflowResult{
		Output: "POCDemoWorkflow is done for username: " + username,
	}, nil
}
