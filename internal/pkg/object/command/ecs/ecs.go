package ecs

import (
	ct "context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/duration"
	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/job"
	"github.com/patterninc/heimdall/pkg/plugin"
	"github.com/patterninc/heimdall/pkg/result"
	"github.com/patterninc/heimdall/pkg/result/column"
)

// ECS command context structure
type ecsCommandContext struct {
	TaskDefinitionTemplate string                    `yaml:"task_definition_template,omitempty" json:"task_definition_template,omitempty"`
	TaskCount              int                       `yaml:"task_count,omitempty" json:"task_count,omitempty"`
	ContainerOverrides     []types.ContainerOverride `yaml:"container_overrides,omitempty" json:"container_overrides,omitempty"`
	PollingInterval        duration.Duration         `yaml:"polling_interval,omitempty" json:"polling_interval,omitempty"`
	Timeout                duration.Duration         `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	MaxFailCount           int                       `yaml:"max_fail_count,omitempty" json:"max_fail_count,omitempty"` // max failures before giving up
}

// ECS job context structure
type ecsJobContext struct {
	TaskCount          int                       `yaml:"task_count,omitempty" json:"task_count,omitempty"`
	ContainerOverrides []types.ContainerOverride `yaml:"container_overrides,omitempty" json:"container_overrides,omitempty"`
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

// Task definition wrapper with pre-computed essential containers map
type taskDefinitionWrapper struct {
	TaskDefinition      *types.TaskDefinition
	EssentialContainers map[string]bool
}

// Task position tracker structure
type taskTracker struct {
	Name          string
	ActiveARN     string
	Retries       int
	ExecutionTime float64
	FailedARNs    []string // History of ARNs for this position
	Completed     bool
}

// executionContext holds the final resolved configuration for job execution.
type executionContext struct {
	TaskCount             int                       `json:"task_count"`
	ContainerOverrides    []types.ContainerOverride `json:"container_overrides"`
	PollingInterval       duration.Duration         `json:"polling_interval"`
	Timeout               duration.Duration         `json:"timeout"`
	MaxFailCount          int                       `json:"max_fail_count"`
	TaskDefinitionWrapper *taskDefinitionWrapper    `json:"task_definition_wrapper"`
	ClusterConfig         *ecsClusterContext        `json:"cluster_config"`
}

const (
	defaultPollingInterval = duration.Duration(30 * time.Second)
	defaultTaskTimeout     = duration.Duration(1 * time.Hour)
	defaultMaxFailCount    = 1
	defaultTaskCount       = 1
	startedByPrefix        = "heimdall-job-"
	errMaxFailCount        = "task %s failed %d times (max: %d), giving up"
	errPollingTimeout      = "polling timed out for arns %v after %v"
)

var (
	ctx                = ct.Background()
	errMissingTemplate = fmt.Errorf("task definition template is required")
)

func New(commandContext *context.Context) (plugin.Handler, error) {

	e := &ecsCommandContext{
		PollingInterval: defaultPollingInterval,
		Timeout:         defaultTaskTimeout,
		MaxFailCount:    defaultMaxFailCount,
		TaskCount:       defaultTaskCount,
	}

	if commandContext != nil {
		if err := commandContext.Unmarshal(e); err != nil {
			return nil, err
		}
	}

	return e.handler, nil

}

// handler implements the main ECS plugin logic
func (e *ecsCommandContext) handler(r *plugin.Runtime, j *job.Job, c *cluster.Cluster) error {

	// Build execution context with resolved configuration and loaded template
	execCtx, err := buildExecutionContext(e, j, c)
	if err != nil {
		return err
	}

	// initialize AWS session
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	ecsClient := ecs.NewFromConfig(cfg)

	// register task definition
	taskDefARN, err := registerTaskDefinition(ecsClient, execCtx)
	if err != nil {
		return err
	}

	// Start tasks
	tasks, err := startTasks(ecsClient, taskDefARN, execCtx, j.ID)
	if err != nil {
		return err
	}

	// Poll for completion
	tasks, err = pollForCompletion(ecsClient, taskDefARN, execCtx, tasks)
	if err != nil {
		return err
	}

	// Build result
	j.Result = &result.Result{}
	j.Result.Columns = []*column.Column{
		{Name: "task_arn", Type: "string"},
		{Name: "duration", Type: "float"},
		{Name: "retries", Type: "int"},
		{Name: "failed_arns", Type: "string"},
	}

	// Create result data from task results
	j.Result.Data = make([][]interface{}, 0, len(tasks))
	for _, tracker := range tasks {
		j.Result.Data = append(j.Result.Data, []interface{}{
			tracker.ActiveARN,
			tracker.ExecutionTime,
			tracker.Retries,
			strings.Join(tracker.FailedARNs, ","),
		})
	}

	return nil

}

// prepare and register task definition with ECS
func registerTaskDefinition(ecsClient *ecs.Client, execCtx *executionContext) (*string, error) {
	registerInput := &ecs.RegisterTaskDefinitionInput{
		Family:                  aws.String(aws.ToString(execCtx.TaskDefinitionWrapper.TaskDefinition.Family)),
		RequiresCompatibilities: []types.Compatibility{types.CompatibilityFargate},
		NetworkMode:             types.NetworkModeAwsvpc,
		Cpu:                     aws.String(fmt.Sprintf("%d", execCtx.ClusterConfig.CPU)),
		Memory:                  aws.String(fmt.Sprintf("%d", execCtx.ClusterConfig.Memory)),
		ExecutionRoleArn:        aws.String(execCtx.ClusterConfig.ExecutionRoleARN),
		TaskRoleArn:             aws.String(execCtx.ClusterConfig.TaskRoleARN),
		ContainerDefinitions:    execCtx.TaskDefinitionWrapper.TaskDefinition.ContainerDefinitions,
	}

	registerOutput, err := ecsClient.RegisterTaskDefinition(ctx, registerInput)
	if err != nil {
		return nil, err
	}

	return registerOutput.TaskDefinition.TaskDefinitionArn, nil

}

// startTasks launches all tasks and returns a map of task trackers
func startTasks(ecsClient *ecs.Client, taskDefARN *string, execCtx *executionContext, jobID string) (map[string]*taskTracker, error) {

	tasks := make(map[string]*taskTracker)

	for i := 0; i < execCtx.TaskCount; i++ {
		taskARN, err := runTask(ecsClient, taskDefARN, execCtx.ClusterConfig, execCtx.ContainerOverrides, fmt.Sprintf("%s%s-%d", startedByPrefix, jobID, i))
		if err != nil {
			return nil, err
		}
		taskName := fmt.Sprintf("%s%s-%d", startedByPrefix, jobID, i)
		tasks[taskName] = &taskTracker{
			Name:      taskName,
			ActiveARN: taskARN,
		}
	}

	return tasks, nil
}

// monitor tasks until completion, faliure, or timeout
func pollForCompletion(ecsClient *ecs.Client, taskDefARN *string, execCtx *executionContext, tasks map[string]*taskTracker) (map[string]*taskTracker, error) {

	startTime := time.Now()
	stopTime := startTime.Add(time.Duration(execCtx.Timeout))

	// Poll until all tasks are complete or timeout
	for {
		// Describe the uncompleted tasks we're tracking
		var activeARNs []string
		for _, tracker := range tasks {
			if !tracker.Completed {
				activeARNs = append(activeARNs, tracker.ActiveARN)
			}
		}

		// If no active tasks, we're done
		if len(activeARNs) == 0 {
			break
		}

		describeInput := &ecs.DescribeTasksInput{
			Cluster: aws.String(execCtx.ClusterConfig.ClusterName),
			Tasks:   activeARNs,
		}

		describeOutput, err := ecsClient.DescribeTasks(ctx, describeInput)
		if err != nil {
			return nil, err
		}

		// Check if all tasks are complete
		allComplete := true

		for _, task := range describeOutput.Tasks {

			// If the task is not stopped, it's not complete
			if aws.ToString(task.LastStatus) != "STOPPED" {
				allComplete = false
				continue
			}

			// If task has stopped, grab its tracker to start updating
			tracker, exists := tasks[aws.ToString(task.StartedBy)]
			if !exists {
				return nil, fmt.Errorf("could not find tracker for StartedBy tag %s", aws.ToString(task.StartedBy))
			}

			// Check for task failures based on exit code
			if !isTaskSuccessful(task, execCtx) {

				tracker.Retries++
				tracker.FailedARNs = append(tracker.FailedARNs, aws.ToString(task.TaskArn))

				// Exit if we've failed too many times
				if tracker.Retries >= execCtx.MaxFailCount {

					// Stop all other running tasks
					reason := fmt.Sprintf(errMaxFailCount, tracker.ActiveARN, tracker.Retries, execCtx.MaxFailCount)
					if err := stopAllTasks(ecsClient, execCtx.ClusterConfig.ClusterName, tasks, reason); err != nil {
						return nil, err
					}

					return nil, fmt.Errorf("%s", reason)
				}

				newTaskARN, err := runTask(ecsClient, taskDefARN, execCtx.ClusterConfig, execCtx.ContainerOverrides, tracker.Name)
				if err != nil {
					return nil, err
				}

				// Assign the new task ARN to the tracker
				tracker.ActiveARN = newTaskARN

				// Task failed but will be restarted, so mark as not complete
				allComplete = false
				continue
			}

			// Update the tracker directly
			tracker.ExecutionTime = time.Since(startTime).Seconds() // Total time from start
			tracker.Completed = true
		}

		// If all tasks are complete, break out of the loop
		if allComplete {
			break
		}

		// Check if we've timed out
		if time.Now().After(stopTime) {
			// Collect ARNs of tasks that did not complete
			var incompleteARNs []string
			for _, tracker := range tasks {
				if !tracker.Completed {
					incompleteARNs = append(incompleteARNs, tracker.ActiveARN)
				}
			}

			// Stop all remaining tasks
			reason := fmt.Sprintf(errPollingTimeout, incompleteARNs, execCtx.Timeout)
			if err := stopAllTasks(ecsClient, execCtx.ClusterConfig.ClusterName, tasks, reason); err != nil {
				return nil, err
			}

			// Return error with information about incomplete tasks
			return nil, fmt.Errorf("%s", reason)
		}

		// Sleep until next poll time
		time.Sleep(time.Duration(execCtx.PollingInterval))
	}

	// If you're here, all tasks are complete
	return tasks, nil

}

func buildExecutionContext(commandCtx *ecsCommandContext, j *job.Job, c *cluster.Cluster) (*executionContext, error) {

	ctx := &executionContext{}

	// Create a context from commandCtx and unmarshal onto execCtx (defaults)
	commandContext := context.New(commandCtx)
	if err := commandContext.Unmarshal(ctx); err != nil {
		return nil, err
	}

	// Overlay job context (overrides command values)
	if j.Context != nil {
		if err := j.Context.Unmarshal(ctx); err != nil {
			return nil, err
		}
	}

	// Add cluster config (no overlapping values)
	clusterContext := &ecsClusterContext{}
	if c.Context != nil {
		if err := c.Context.Unmarshal(clusterContext); err != nil {
			return nil, err
		}
	}
	ctx.ClusterConfig = clusterContext

	// Load task definition template
	taskDefWrapper, err := loadTaskDefinitionTemplate(commandCtx.TaskDefinitionTemplate)
	if err != nil {
		return nil, err
	}
	ctx.TaskDefinitionWrapper = taskDefWrapper // Store the wrapper for polling

	// Build container overrides for all containers
	if err := buildContainerOverrides(ctx); err != nil {
		return nil, err
	}

	// Validate the resolved configuration
	if err := validateExecutionContext(ctx); err != nil {
		return nil, err
	}

	return ctx, nil

}

// validateExecutionContext validates the final resolved configuration
func validateExecutionContext(ctx *executionContext) error {

	if ctx.TaskCount <= 0 || ctx.TaskCount > ctx.ClusterConfig.MaxTaskCount {
		return fmt.Errorf("task count (%d) needs to be greater than 0 and less than cluster max task count (%d)", ctx.TaskCount, ctx.ClusterConfig.MaxTaskCount)
	}

	return nil

}

// buildContainerOverrides processes container overrides and builds the final overrides for all containers
func buildContainerOverrides(execCtx *executionContext) error {

	// Create a map of container names
	existingContainers := make(map[string]bool)
	for _, container := range execCtx.TaskDefinitionWrapper.TaskDefinition.ContainerDefinitions {
		existingContainers[aws.ToString(container.Name)] = true
	}

	// Create a map of execution context overrides
	containerOverridesMap := make(map[string]types.ContainerOverride)
	for _, containerOverride := range execCtx.ContainerOverrides {
		// Validate that the container name exists
		if !existingContainers[aws.ToString(containerOverride.Name)] {
			return fmt.Errorf("container override '%s' not found in task definition template", aws.ToString(containerOverride.Name))
		}
		containerOverridesMap[aws.ToString(containerOverride.Name)] = containerOverride
	}

	// Build container overrides for all containers in the task definition
	var containerOverrides []types.ContainerOverride
	for _, container := range execCtx.TaskDefinitionWrapper.TaskDefinition.ContainerDefinitions {
		containerName := aws.ToString(container.Name)

		// Use existing override if it exists, otherwise create a blank one
		if override, exists := containerOverridesMap[containerName]; exists {
			containerOverrides = append(containerOverrides, override)
		} else {
			containerOverrides = append(containerOverrides, types.ContainerOverride{
				Name: aws.String(containerName),
			})
		}
	}

	execCtx.ContainerOverrides = containerOverrides

	return nil

}

// stopAllTasks stops all non-completed tasks with the given reason
func stopAllTasks(ecsClient *ecs.Client, clusterName string, tasks map[string]*taskTracker, reason string) error {

	for _, t := range tasks {
		if t.Completed {
			continue
		}
		stopInput := &ecs.StopTaskInput{
			Cluster: aws.String(clusterName),
			Task:    aws.String(t.ActiveARN),
			Reason:  aws.String(reason),
		}

		_, err := ecsClient.StopTask(ctx, stopInput)
		if err != nil {
			return err
		}
	}

	return nil

}

func loadTaskDefinitionTemplate(templatePath string) (*taskDefinitionWrapper, error) {

	if templatePath == `` {
		return nil, errMissingTemplate
	}

	data, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, err
	}

	var taskDef types.TaskDefinition
	if err := json.Unmarshal(data, &taskDef); err != nil {
		return nil, err
	}

	// Pre-compute essential containers map
	essentialContainers := make(map[string]bool)
	for _, containerDef := range taskDef.ContainerDefinitions {
		if containerDef.Essential != nil && *containerDef.Essential {
			essentialContainers[aws.ToString(containerDef.Name)] = true
		}
	}

	return &taskDefinitionWrapper{
		TaskDefinition:      &taskDef,
		EssentialContainers: essentialContainers,
	}, nil

}

// runTask runs a single task and returns the task ARN
func runTask(ecsClient *ecs.Client, taskDefARN *string, clusterContext *ecsClusterContext, containerOverrides []types.ContainerOverride, startedBy string) (string, error) {

	// Create a copy of the overrides and add TASK_NAME env variable
	finalOverrides := append([]types.ContainerOverride{}, containerOverrides...)
	for i := range finalOverrides {
		finalOverrides[i].Environment = append(finalOverrides[i].Environment, types.KeyValuePair{
			Name:  aws.String("TASK_NAME"),
			Value: aws.String(startedBy),
		})
	}

	// build run task input
	runTaskInput := &ecs.RunTaskInput{
		Cluster:        aws.String(clusterContext.ClusterName),
		TaskDefinition: taskDefARN,
		LaunchType:     types.LaunchType(clusterContext.LaunchType),
		Count:          aws.Int32(1),
		StartedBy:      aws.String(startedBy),
		Overrides: &types.TaskOverride{
			ContainerOverrides: finalOverrides,
		},
		NetworkConfiguration: &types.NetworkConfiguration{
			AwsvpcConfiguration: &types.AwsVpcConfiguration{
				Subnets:        clusterContext.VPCConfig.Subnets,
				SecurityGroups: clusterContext.VPCConfig.SecurityGroups,
				AssignPublicIp: types.AssignPublicIpDisabled,
			},
		},
	}

	runTaskOutput, err := ecsClient.RunTask(ctx, runTaskInput)
	if err != nil {
		return ``, err
	}

	taskARN := aws.ToString(runTaskOutput.Tasks[0].TaskArn)

	return taskARN, nil

}

// tasks are successful if all essential containers exit with a zero exit code
func isTaskSuccessful(task types.Task, execCtx *executionContext) bool {
	// Check all containers in the running task
	for _, container := range task.Containers {
		containerName := aws.ToString(container.Name)

		if execCtx.TaskDefinitionWrapper.EssentialContainers[containerName] {
			if container.ExitCode != nil && *container.ExitCode != 0 {
				return false
			}
		}
	}

	return true

}
