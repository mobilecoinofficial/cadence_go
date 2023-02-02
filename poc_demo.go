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
	hostname, _ := os.Hostname()
	username := faker.Username()

	msg := fmt.Sprintf("POCDemoActivity1 activity1 is running on %s\n", hostname)
	activity.GetLogger(ctx).Info(
		msg,
		zap.String("input", input),
	)

	return username, nil
}

func POCDemoActivity2(ctx context.Context, input string) (string, error) {
	hostname, _ := os.Hostname()

	msg := fmt.Sprintf("POCDemo POCDemoActivity2 is running on %s\n", hostname)
	activity.GetLogger(ctx).Info(
		msg,
		zap.String("input", input),
	)

	return fmt.Sprintf("activity2[%s]", input), nil
}

func POCDemoActivity3(ctx context.Context, input string) (string, error) {
	hostname, _ := os.Hostname()

	msg := fmt.Sprintf("POCDemo POCDemoActivity3 is running on %s\n", hostname)
	activity.GetLogger(ctx).Info(
		msg,
		zap.String("input", input),
	)

	return fmt.Sprintf("activity3{%s}", input), nil
}

func POCChildWorkflow1(ctx workflow.Context, input string) (*POCDemoWorkflowResult, error) {
	log.Printf("POCChildWorkflow1 is triggered with input: %s\n", input)

	res := POCDemoWorkflowResult{}

	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute * 3,
	}

	nctx := workflow.WithActivityOptions(ctx, ao)

	err := workflow.ExecuteActivity(
		nctx,
		POCDemoActivity1,
		input,
	).Get(ctx, nil)
	if err != nil {
		workflow.GetLogger(ctx).Error("POCChildWorkflow1 failed.", zap.Error(err))
		return &res, err
	}

	log.Printf("POCChildWorkflow1 is done for input: %s, result: %s\n", input, res.Output)
	return &POCDemoWorkflowResult{
		Output: fmt.Sprintf("%s/%s", input, res.Output),
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

	log.Printf("POCChildWorkflow2 is done for input: %s, generated final: %s\n", input, res.Output)
	return &POCDemoWorkflowResult{
		Output: res.Output,
	}, nil

}

func POCChildWorkflow3(ctx workflow.Context, input string) (*POCDemoWorkflowResult, error) {
	log.Println("POCChildWorkflow3 is triggered with input: ", input)

	res := POCDemoWorkflowResult{}

	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute * 3,
	}

	nctx := workflow.WithActivityOptions(ctx, ao)

	err := workflow.ExecuteActivity(
		nctx,
		POCDemoActivity3,
		input,
	).Get(ctx, nil)
	if err != nil {
		workflow.GetLogger(ctx).Error("POCChildWorkflow3 failed.", zap.Error(err))
		return &res, err
	}

	// rand.Intn(max - min) + min

	log.Printf("POCChildWorkflow3 is done for input: %s, result: %s\n", input, res.Output)
	return &POCDemoWorkflowResult{
		Output: res.Output,
	}, nil

}

func POCWorkflow(ctx workflow.Context) (*POCDemoWorkflowResult, error) {
	log.Println("POCDemoWorkflow is triggered...")

	workflowID := workflow.GetInfo(ctx).WorkflowExecution.ID
	log.Println("Received POCDemoWorkflow. WorkflowID: ", workflowID)

	cwo1 := workflow.ChildWorkflowOptions{
		//WorkflowID:                   workflowID,
		ExecutionStartToCloseTimeout: time.Minute * 3,
	}
	ctx1 := workflow.WithChildOptions(ctx, cwo1)

	log.Printf("Starting Child Workflow: POCChildWorkflow1\n")

	var result1 POCDemoWorkflowResult
	var err error
	future1 := workflow.ExecuteChildWorkflow(ctx1, POCChildWorkflow1, "input1")
	err = future1.Get(ctx, &result1)
	if err != nil {
		log.Printf("Error: %s", err)
		os.Exit(1)
	}
	log.Printf("POCChildWorkflow1 ended with result: %s\n", result1.Output)

	cwo2 := workflow.ChildWorkflowOptions{
		//WorkflowID:                   workflowID,
		ExecutionStartToCloseTimeout: time.Minute * 3,
	}
	ctx2 := workflow.WithChildOptions(ctx, cwo2)

	log.Printf("Starting Child Workflow: POCChildWorkflow2\n")

	var result2 POCDemoWorkflowResult
	future2 := workflow.ExecuteChildWorkflow(ctx2, POCChildWorkflow2, result1.Output)
	err = future2.Get(ctx, &result2)
	if err != nil {
		log.Printf("Error: %s", err)
		os.Exit(1)
	}
	log.Printf("POCChildWorkflow2 ended with result: %s\n", result2.Output)

	cwo3 := workflow.ChildWorkflowOptions{
		//WorkflowID:                   workflowID,
		ExecutionStartToCloseTimeout: time.Minute * 3,
	}
	ctx3 := workflow.WithChildOptions(ctx, cwo3)

	log.Printf("Starting Child Workflow: POCChildWorkflow3\n")

	var result3 POCDemoWorkflowResult
	future3 := workflow.ExecuteChildWorkflow(ctx3, POCChildWorkflow3, result2.Output)
	err = future3.Get(ctx, &result3)
	if err != nil {
		log.Printf("Error: %s", err)
		os.Exit(1)
	}
	log.Printf("POCChildWorkflow3 ended with result: %s\n", result3.Output)

	return &POCDemoWorkflowResult{
		Output: "POCDemoWorkflow is done for username: " + fmt.Sprintf("%s/%s/%s", result1.Output, result2.Output, result3.Output),
	}, nil
}
