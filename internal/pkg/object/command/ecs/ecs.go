package ecs

import (
	ct "context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	hc "github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/job"
	"github.com/patterninc/heimdall/pkg/plugin"
	"github.com/patterninc/heimdall/pkg/result"
	"github.com/patterninc/heimdall/pkg/result/column"
)

// ECS command context structure
type ecsCommandContext struct {
	TaskDefinitionTemplate string              `yaml:"task_definition_template,omitempty" json:"task_definition_template,omitempty"`
	TaskCount              int                 `yaml:"task_count,omitempty" json:"task_count,omitempty"`
	ContainerOverrides     []containerOverride `yaml:"container_overrides,omitempty" json:"container_overrides,omitempty"`
	PollingInterval        int                 `yaml:"polling_interval,omitempty" json:"polling_interval,omitempty"` // in seconds
	Timeout                int                 `yaml:"timeout,omitempty" json:"timeout,omitempty"`                   // in seconds
	MaxFailCount           int                 `yaml:"max_fail_count,omitempty" json:"max_fail_count,omitempty"`     // max failures before giving up
}

// ECS job context structure
type ecsJobContext struct {
	TaskCount          int                 `yaml:"task_count,omitempty" json:"task_count,omitempty"`
	ContainerOverrides []containerOverride `yaml:"container_overrides,omitempty" json:"container_overrides,omitempty"`
}

// ECS cluster context structure
type ecsClusterContext struct {
	CPU              int       `yaml:"cpu,omitempty" json:"cpu,omitempty"`
	Memory           int       `yaml:"memory,omitempty" json:"memory,omitempty"`
	MaxTaskCount     int       `yaml:"max_task_count,omitempty" json:"max_task_count,omitempty"`
	ExecutionRoleARN string    `yaml:"execution_role_arn,omitempty" json:"execution_role_arn,omitempty"`
	TaskRoleARN      string    `yaml:"task_role_arn,omitempty" json:"task_role_arn,omitempty"`
	ClusterName      string    `yaml:"cluster_name,omitempty" json:"cluster_name,omitempty"`
	LaunchType       string    `yaml:"launch_type,omitempty" json:"launch_type,omitempty"`
	VPCConfig        vpcConfig `yaml:"vpc_config,omitempty" json:"vpc_config,omitempty"`
}

// VPC configuration structure
type vpcConfig struct {
	Subnets        []string `yaml:"subnets,omitempty" json:"subnets,omitempty"`
	SecurityGroups []string `yaml:"security_groups,omitempty" json:"security_groups,omitempty"`
}

// Container override structure
type containerOverride struct {
	Name    string   `yaml:"name,omitempty" json:"name,omitempty"`
	Command []string `yaml:"command,omitempty" json:"command,omitempty"`
}

// Task result structure
type taskResult struct {
	TaskARN       string  `json:"task_arn"`
	ExitCode      int     `json:"exit_code"`
	ExecutionTime float64 `json:"execution_time"` // in seconds
	Status        string  `json:"status"`
	Retries       int     `json:"retries"` // Added for retries
}

const (
	defaultPollingInterval = 30   // seconds
	defaultTaskTimeout     = 3600 // seconds (1 hour)
	defaultMaxFailCount    = 1
)

var (
	ctx = ct.Background()
)

// New creates a new ECS plugin handler
func New(commandContext *hc.Context) (plugin.Handler, error) {

	e := &ecsCommandContext{}

	if commandContext != nil {
		if err := commandContext.Unmarshal(e); err != nil {
			return nil, err
		}
	}

	return e.handler, nil

}

// handler implements the main ECS plugin logic
func (e *ecsCommandContext) handler(r *plugin.Runtime, j *job.Job, c *cluster.Cluster) error {

	// Prepare job context with command defaults
	jobContext := &ecsJobContext{}
	if j.Context != nil {
		if err := j.Context.Unmarshal(jobContext); err != nil {
			return fmt.Errorf("failed to unmarshal job context: %w", err)
		}
	}

	// Fill in command defaults where job context values are not specified
	if jobContext.TaskCount <= 0 {
		jobContext.TaskCount = e.TaskCount
	}
	if len(jobContext.ContainerOverrides) == 0 {
		jobContext.ContainerOverrides = e.ContainerOverrides
	}

	// Set default task count if not specified
	if jobContext.TaskCount <= 0 {
		jobContext.TaskCount = 1
	}

	// Set default polling interval if not specified
	if e.PollingInterval == 0 {
		e.PollingInterval = defaultPollingInterval
	}

	// Set default timeout if not specified
	if e.Timeout == 0 {
		e.Timeout = defaultTaskTimeout
	}

	// unmarshal cluster context
	clusterContext := &ecsClusterContext{}
	if c.Context != nil {
		if err := c.Context.Unmarshal(clusterContext); err != nil {
			return err
		}
	}

	// Validate task count
	if jobContext.TaskCount > clusterContext.MaxTaskCount {
		return fmt.Errorf("task count (%d) exceeds cluster max task count (%d)",
			jobContext.TaskCount, clusterContext.MaxTaskCount)
	}

	// Initialize AWS session
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	ecsClient := ecs.NewFromConfig(cfg)

	// Load task definition template
	taskDef, err := loadTaskDefinitionTemplate(e.TaskDefinitionTemplate)
	if err != nil {
		return fmt.Errorf("failed to load task definition template: %w", err)
	}

	// Apply overrides from command and cluster context
	if err := applyOverrides(taskDef, jobContext); err != nil {
		return err
	}

	// Register the task definition
	registerInput := &ecs.RegisterTaskDefinitionInput{
		Family:                  aws.String(fmt.Sprintf("heimdall-task-%s", j.ID)),
		RequiresCompatibilities: []types.Compatibility{types.CompatibilityFargate},
		NetworkMode:             types.NetworkModeAwsvpc,
		Cpu:                     aws.String(fmt.Sprintf("%d", clusterContext.CPU)),
		Memory:                  aws.String(fmt.Sprintf("%d", clusterContext.Memory)),
		ExecutionRoleArn:        aws.String(clusterContext.ExecutionRoleARN),
		TaskRoleArn:             aws.String(clusterContext.TaskRoleARN),
		ContainerDefinitions:    taskDef.ContainerDefinitions,
	}

	registerOutput, err := ecsClient.RegisterTaskDefinition(ctx, registerInput)
	if err != nil {
		return fmt.Errorf("failed to register task definition: %w", err)
	}

	// Run tasks
	_, err = runTasks(ecsClient, registerOutput.TaskDefinition, clusterContext, jobContext, j.ID)
	if err != nil {
		return fmt.Errorf("failed to run tasks: %w", err)
	}

	// Poll for task completion using startedBy tag
	taskResults, err := pollTaskCompletion(ecsClient, clusterContext.ClusterName, j.ID, jobContext, e, taskDef, clusterContext, r)
	if err != nil {
		return fmt.Errorf("failed to poll task completion: %w", err)
	}

	j.Result = &result.Result{}

	// Set columns
	j.Result.Columns = []*column.Column{
		{Name: "task_arn", Type: "string"},
		{Name: "duration", Type: "float"},
		{Name: "retries", Type: "int"},
		{Name: "exit_code", Type: "int"},
		{Name: "status", Type: "string"},
	}

	// Create result data from task results
	j.Result.Data = make([][]interface{}, 0, len(taskResults))
	for _, taskResult := range taskResults {
		j.Result.Data = append(j.Result.Data, []interface{}{
			taskResult.TaskARN,
			taskResult.ExecutionTime,
			taskResult.Retries,
			taskResult.ExitCode,
			taskResult.Status,
		})
	}

	return nil

}

func loadTaskDefinitionTemplate(templatePath string) (*types.TaskDefinition, error) {
	if templatePath == "" {
		// Return a basic template if none provided
		return &types.TaskDefinition{
			ContainerDefinitions: []types.ContainerDefinition{
				{
					Name:  aws.String("main"),
					Image: aws.String("alpine:latest"),
				},
			},
		}, nil
	}

	data, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	var taskDef types.TaskDefinition
	if err := json.Unmarshal(data, &taskDef); err != nil {
		return nil, fmt.Errorf("failed to unmarshal task definition: %w", err)
	}

	return &taskDef, nil
}

func applyOverrides(taskDef *types.TaskDefinition, jobContext *ecsJobContext) error {
	// Apply container overrides
	if len(jobContext.ContainerOverrides) > 0 {
		for _, override := range jobContext.ContainerOverrides {
			for i, container := range taskDef.ContainerDefinitions {
				if aws.ToString(container.Name) == override.Name {
					if len(override.Command) > 0 {
						taskDef.ContainerDefinitions[i].Command = override.Command
					}
					break
				}
			}
		}
	}

	return nil
}

// buildRunTaskInput creates a RunTaskInput with the given parameters
func buildRunTaskInput(clusterName string, taskDef *types.TaskDefinition, clusterContext *ecsClusterContext, startedBy string) *ecs.RunTaskInput {
	return &ecs.RunTaskInput{
		Cluster:        aws.String(clusterName),
		TaskDefinition: taskDef.TaskDefinitionArn,
		LaunchType:     types.LaunchType(clusterContext.LaunchType),
		StartedBy:      aws.String(startedBy),
		NetworkConfiguration: &types.NetworkConfiguration{
			AwsvpcConfiguration: &types.AwsVpcConfiguration{
				Subnets:        clusterContext.VPCConfig.Subnets,
				SecurityGroups: clusterContext.VPCConfig.SecurityGroups,
				AssignPublicIp: types.AssignPublicIpDisabled,
			},
		},
	}
}

// runTasks runs the specified number of tasks
func runTasks(ecsClient *ecs.Client, taskDef *types.TaskDefinition, clusterContext *ecsClusterContext, jobContext *ecsJobContext, jobID string) ([]string, error) {

	var taskARNs []string

	for i := 0; i < jobContext.TaskCount; i++ {
		startedBy := fmt.Sprintf("heimdall-job-%s-%d", jobID, i+1)

		runTaskInput := buildRunTaskInput(clusterContext.ClusterName, taskDef, clusterContext, startedBy)

		runTaskOutput, err := ecsClient.RunTask(ctx, runTaskInput)
		if err != nil {
			return taskARNs, fmt.Errorf("failed to run task %d: %w", i+1, err)
		}

		if len(runTaskOutput.Tasks) > 0 {
			taskARNs = append(taskARNs, aws.ToString(runTaskOutput.Tasks[0].TaskArn))
		}
	}

	return taskARNs, nil

}

// pollTaskCompletion polls for task completion
func pollTaskCompletion(ecsClient *ecs.Client, clusterName string, jobID string, jobContext *ecsJobContext, commandContext *ecsCommandContext, taskDef *types.TaskDefinition, clusterContext *ecsClusterContext, r *plugin.Runtime) ([]taskResult, error) {

	fmt.Println("Polling for task completion")

	var taskResults []taskResult

	if commandContext.PollingInterval == 0 {
		commandContext.PollingInterval = defaultPollingInterval
	}

	if commandContext.Timeout == 0 {
		commandContext.Timeout = defaultTaskTimeout
	}

	if commandContext.MaxFailCount == 0 {
		commandContext.MaxFailCount = defaultMaxFailCount
	}

	// Start polling time
	startTime := time.Now()
	startedByPrefix := fmt.Sprintf("heimdall-job-%s", jobID)

	// Track failures per task
	taskFailures := make(map[string]int)

	// Poll until all tasks are complete or timeout
	for time.Since(startTime) < time.Duration(commandContext.Timeout)*time.Second {
		// List tasks with our startedBy prefix (includes original and restarted tasks)
		listInput := &ecs.ListTasksInput{
			Cluster:   aws.String(clusterName),
			StartedBy: aws.String(startedByPrefix),
		}

		listOutput, err := ecsClient.ListTasks(ctx, listInput)
		if err != nil {
			return nil, err
		}

		// Tasks may not be running yet, so we need to poll again
		if len(listOutput.TaskArns) == 0 {
			r.Stdout.WriteString(fmt.Sprintf("No tasks found with startedBy prefix: %s\n", startedByPrefix))
			time.Sleep(time.Duration(commandContext.PollingInterval) * time.Second)
			continue
		}

		// Describe all tasks
		describeInput := &ecs.DescribeTasksInput{
			Cluster: aws.String(clusterName),
			Tasks:   listOutput.TaskArns,
		}

		describeOutput, err := ecsClient.DescribeTasks(ctx, describeInput)
		if err != nil {
			return nil, err
		}

		// Check if all tasks are complete
		allComplete := true
		taskResults = []taskResult{} // Reset results

		for _, task := range describeOutput.Tasks {
			taskARN := aws.ToString(task.TaskArn)
			lastStatus := aws.ToString(task.LastStatus)
			r.Stdout.WriteString(fmt.Sprintf("Task %s status: %s\n", taskARN, lastStatus))

			if lastStatus == "RUNNING" || lastStatus == "PENDING" {
				allComplete = false
				continue
			}

			// Task is complete (STOPPED)
			var exitCode int
			var executionTime float64

			if len(task.Containers) > 0 {
				container := task.Containers[0]
				if container.ExitCode != nil {
					exitCode = int(*container.ExitCode)
				}

				// Calculate execution time if we have start/stop times
				if task.StartedAt != nil && task.StoppedAt != nil {
					executionTime = task.StoppedAt.Sub(*task.StartedAt).Seconds()
				}
			}

			// Check for task failures
			if exitCode != 0 {
				taskFailures[taskARN]++
				failCount := taskFailures[taskARN]
				r.Stdout.WriteString(fmt.Sprintf("Task %s failed with exit code %d (failure #%d)\n", taskARN, exitCode, failCount))

				// Exit if we've failed too many times
				if failCount >= commandContext.MaxFailCount {
					return taskResults, fmt.Errorf("task %s failed %d times (max: %d), giving up", taskARN, failCount, commandContext.MaxFailCount)
				}

				// Restart the failed task
				r.Stdout.WriteString(fmt.Sprintf("Restarting task %s (attempt %d)\n", taskARN, failCount+1))

				// Create new startedBy tag for the restart
				restartStartedBy := fmt.Sprintf("heimdall-job-%s-restart-%d", jobID, failCount+1)

				runTaskInput := buildRunTaskInput(clusterName, taskDef, clusterContext, restartStartedBy)

				restartOutput, err := ecsClient.RunTask(ctx, runTaskInput)
				if err != nil {
					return nil, err
				}

				if len(restartOutput.Tasks) > 0 {
					newTaskARN := aws.ToString(restartOutput.Tasks[0].TaskArn)
					r.Stdout.WriteString(fmt.Sprintf("Task restarted: %s -> %s\n", taskARN, newTaskARN))
					// Update the failure count for the new task ARN
					taskFailures[newTaskARN] = failCount
				}

				// Task failed but will be restarted, so mark as not complete
				allComplete = false
				continue

			}

			// Task completed successfully
			taskResults = append(taskResults, taskResult{
				TaskARN:       taskARN,
				ExitCode:      exitCode,
				ExecutionTime: executionTime,
				Status:        lastStatus,
				Retries:       taskFailures[taskARN], // Add retries to the result
			})
		}

		if allComplete {
			r.Stdout.WriteString("All tasks completed successfully\n")
			break
		}

		// Wait before next poll
		time.Sleep(time.Duration(commandContext.PollingInterval) * time.Second)
	}

	// Check if we timed out
	if time.Since(startTime) >= time.Duration(commandContext.Timeout)*time.Second {
		return taskResults, fmt.Errorf("polling timed out after %v", time.Duration(commandContext.Timeout)*time.Second)
	}

	// Return results without erroring on non-zero exit codes
	// The task restart logic should be handled elsewhere
	return taskResults, nil

}
