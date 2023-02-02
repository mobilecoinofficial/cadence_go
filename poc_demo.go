package cadence_go

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bxcodec/faker/v3"
	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

type POCDemoWorkflowResult struct {
	Output string
}

func POCDemoActivity1(ctx context.Context, input string) (string, error) {
	//workflowID := activity.GetInfo(ctx).WorkflowExecution.ID

	sentence := faker.Sentence()
	hostname, _ := os.Hostname()

	msg := fmt.Sprintf("POCDemoActivity1 activity1 is running on %s, sentence: %s\n", hostname, sentence)
	activity.GetLogger(ctx).Info(
		msg,
		zap.String("input", input),
	)

	activitySummary := fmt.Sprintf("POCDemoActivity1, sentence: %s\n", sentence)

	return activitySummary, nil
}

func POCDemoActivity2(ctx context.Context, input string) (string, error) {
	//workflowID := activity.GetInfo(ctx).WorkflowExecution.ID

	sentence := faker.Sentence()
	hostname, _ := os.Hostname()

	msg := fmt.Sprintf("POCDemo POCDemoActivity2 is running on %s, sentence: %s\n", hostname, sentence)
	activity.GetLogger(ctx).Info(
		msg,
		zap.String("input", input),
	)

	activitySummary := fmt.Sprintf("POCDemoActivity2, sentence: %s\n", sentence)

	return activitySummary, nil
}

func POCChildWorkflow1(ctx workflow.Context, input string) (*POCDemoWorkflowResult, error) {
	log.Printf("POCChildWorkflow1 is triggered with input: %s\n", input)

	res := POCDemoWorkflowResult{}

	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute * 3,
	}

	nctx := workflow.WithActivityOptions(ctx, ao)

	username := faker.Name()

	err := workflow.ExecuteActivity(
		nctx,
		POCDemoActivity1,
		username,
	).Get(ctx, nil)
	if err != nil {
		workflow.GetLogger(ctx).Error("POCChildWorkflow1 failed.", zap.Error(err))
		return &res, err
	}

	log.Printf("POCChildWorkflow1 is done for input: %s, generated username: %s\n", input, username)
	return &POCDemoWorkflowResult{
		Output: fmt.Sprintf("%s/%s", input, username),
	}, nil

}

func POCChildWorkflow2(ctx workflow.Context, input string) (*POCDemoWorkflowResult, error) {
	log.Println("POCChildWorkflow2 is triggered with input: ", input)

	res := POCDemoWorkflowResult{}

	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute * 3,
	}

	nctx := workflow.WithActivityOptions(ctx, ao)

	err := workflow.ExecuteActivity(
		nctx,
		POCDemoActivity2,
		input,
	).Get(ctx, nil)
	if err != nil {
		workflow.GetLogger(ctx).Error("POCChildWorkflow2 failed.", zap.Error(err))
		return &res, err
	}

	final := fmt.Sprintf("FINAL-%s", input)
	log.Printf("POCChildWorkflow2 is done for input: %s, generated final: %s\n", input, final)
	return &POCDemoWorkflowResult{
		Output: final,
	}, nil

}

func POCWorkflow(ctx workflow.Context) (*POCDemoWorkflowResult, error) {
	log.Println("POCDemoWorkflow is triggered...")

	// ao := workflow.ActivityOptions{
	// 	ScheduleToStartTimeout: time.Minute,
	// 	StartToCloseTimeout:    time.Minute * 3,
	// }

	//nctx := workflow.WithActivityOptions(ctx, ao)

	log.Println("POCDemoWorkflow. Waiting for Completion...")
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

	log.Println("POCDemoWorkflow. Starting workflow...")
	//endTime := workflow.Now(ctx)

	workflowID := workflow.GetInfo(ctx).WorkflowExecution.ID
	log.Println("Received POCDemoWorkflow. WorkflowID: ", workflowID)

	cwo1 := workflow.ChildWorkflowOptions{
		//WorkflowID:                   workflowID,
		ExecutionStartToCloseTimeout: time.Minute * 3,
	}
	ctx1 := workflow.WithChildOptions(ctx, cwo1)

	log.Printf("Starting Child Workflow1\n")

	var result1 POCDemoWorkflowResult
	var err error
	future1 := workflow.ExecuteChildWorkflow(ctx1, POCChildWorkflow1, "input1")
	err = future1.Get(ctx, &result1)
	if err != nil {
		log.Printf("Error: %s", err)
		os.Exit(1)
	}
	log.Printf("Child Workflow1 ended with result: %s\n", result1.Output)

	cwo2 := workflow.ChildWorkflowOptions{
		//WorkflowID:                   workflowID,
		ExecutionStartToCloseTimeout: time.Minute * 3,
	}
	ctx2 := workflow.WithChildOptions(ctx, cwo2)

	var result2 POCDemoWorkflowResult
	future2 := workflow.ExecuteChildWorkflow(ctx2, POCChildWorkflow2, result1.Output)
	err = future2.Get(ctx, &result2)
	if err != nil {
		log.Printf("Error: %s", err)
		os.Exit(1)
	}
	log.Printf("Child Workflow2 ended with result: %s\n", result2.Output)

	// username := faker.Username()

	// var err error
	// err = workflow.ExecuteActivity(
	// 	nctx,
	// 	POCDemoActivity1,
	// 	username,
	// ).Get(ctx, nil)
	// if err != nil {
	// 	workflow.GetLogger(ctx).Error("POCDemoActivity1 failed.", zap.Error(err))
	// 	return nil, err
	// }

	// err = workflow.ExecuteActivity(
	// 	nctx,
	// 	POCDemoActivity2,
	// 	username,
	// ).Get(ctx, nil)
	// if err != nil {
	// 	workflow.GetLogger(ctx).Error("POCDemoActivity2 failed.", zap.Error(err))
	// 	return nil, err
	// }

	return &POCDemoWorkflowResult{
		Output: "POCDemoWorkflow is done for username: " + fmt.Sprintf("%s/%s", result1.Output, result2.Output),
	}, nil
}
