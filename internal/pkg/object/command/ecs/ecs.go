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
	CPU                    int                       `yaml:"cpu,omitempty" json:"cpu,omitempty"`
	Memory                 int                       `yaml:"memory,omitempty" json:"memory,omitempty"`
	ContainerOverrides     []types.ContainerOverride `yaml:"container_overrides,omitempty" json:"container_overrides,omitempty"`
	PollingInterval        duration.Duration         `yaml:"polling_interval,omitempty" json:"polling_interval,omitempty"`
	Timeout                duration.Duration         `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	MaxFailCount           int                       `yaml:"max_fail_count,omitempty" json:"max_fail_count,omitempty"` // max failures before giving up
}

// ECS cluster context structure
type ecsClusterContext struct {
	MaxCPU           int       `yaml:"max_cpu,omitempty" json:"max_cpu,omitempty"`
	MaxMemory        int       `yaml:"max_memory,omitempty" json:"max_memory,omitempty"`
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
	TaskNum       int // Original task number (0, 1, 2, etc.)
	Retries       int
	ExecutionTime float64
	FailedARNs    []string // History of ARNs for this position
	Completed     bool
}

// executionContext holds the final resolved configuration for job execution.
type executionContext struct {
	TaskCount             int                       `json:"task_count"`
	CPU                   int                       `json:"cpu"`
	Memory                int                       `json:"memory"`
	TaskDefinitionWrapper *taskDefinitionWrapper    `json:"task_definition_wrapper"`
	ContainerOverrides    []types.ContainerOverride `json:"container_overrides"`
	ClusterConfig         *ecsClusterContext        `json:"cluster_config"`

	PollingInterval duration.Duration `json:"polling_interval"`
	Timeout         duration.Duration `json:"timeout"`
	MaxFailCount    int               `json:"max_fail_count"`

	ecsClient  *ecs.Client
	taskDefARN *string
	tasks      map[string]*taskTracker
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
func (e *ecsCommandContext) handler(r *plugin.Runtime, job *job.Job, cluster *cluster.Cluster) error {

	// Build execution context with resolved configuration and loaded template
	execCtx, err := buildExecutionContext(e, job, cluster)
	if err != nil {
		return err
	}

	// register task definition
	if err := execCtx.registerTaskDefinition(); err != nil {
		return err
	}

	// Start tasks
	if err := execCtx.startTasks(job.ID); err != nil {
		return err
	}

	// Poll for completion
	if err := execCtx.pollForCompletion(); err != nil {
		return err
	}

	// Store results
	if err := storeResults(execCtx, job); err != nil {
		return err
	}

	return nil

}

// prepare and register task definition with ECS
func (execCtx *executionContext) registerTaskDefinition() error {
	registerInput := &ecs.RegisterTaskDefinitionInput{
		Family:                  aws.String(aws.ToString(execCtx.TaskDefinitionWrapper.TaskDefinition.Family)),
		RequiresCompatibilities: []types.Compatibility{types.CompatibilityFargate},
		NetworkMode:             types.NetworkModeAwsvpc,
		Cpu:                     aws.String(fmt.Sprintf("%d", execCtx.CPU)),
		Memory:                  aws.String(fmt.Sprintf("%d", execCtx.Memory)),
		ExecutionRoleArn:        aws.String(execCtx.ClusterConfig.ExecutionRoleARN),
		TaskRoleArn:             aws.String(execCtx.ClusterConfig.TaskRoleARN),
		ContainerDefinitions:    execCtx.TaskDefinitionWrapper.TaskDefinition.ContainerDefinitions,
	}

	registerOutput, err := execCtx.ecsClient.RegisterTaskDefinition(ctx, registerInput)
	if err != nil {
		return err
	}

	execCtx.taskDefARN = registerOutput.TaskDefinition.TaskDefinitionArn

	return nil

}

// startTasks launches all tasks and returns a map of task trackers
func (execCtx *executionContext) startTasks(jobID string) error {

	for i := 0; i < execCtx.TaskCount; i++ {
		taskARN, err := runTask(execCtx, fmt.Sprintf("%s%s-%d", startedByPrefix, jobID, i), i)
		if err != nil {
			return err
		}
		taskName := fmt.Sprintf("%s%s-%d", startedByPrefix, jobID, i)
		execCtx.tasks[taskName] = &taskTracker{
			Name:      taskName,
			ActiveARN: taskARN,
			TaskNum:   i,
		}
	}

	return nil
}

// monitor tasks until completion, faliure, or timeout
func (execCtx *executionContext) pollForCompletion() error {

	startTime := time.Now()
	stopTime := startTime.Add(time.Duration(execCtx.Timeout))

	// Poll until all tasks are complete or timeout
	for {
		// Describe the uncompleted tasks we're tracking
		var activeARNs []string
		for _, tracker := range execCtx.tasks {
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

		describeOutput, err := execCtx.ecsClient.DescribeTasks(ctx, describeInput)
		if err != nil {
			return err
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
			tracker, exists := execCtx.tasks[aws.ToString(task.StartedBy)]
			if !exists {
				return fmt.Errorf("could not find tracker for StartedBy tag %s", aws.ToString(task.StartedBy))
			}

			// Check for task failures based on exit code
			if isTaskSuccessful(task, execCtx) {
				// Update the tracker directly
				tracker.ExecutionTime = time.Since(startTime).Seconds() // Total time from start
				tracker.Completed = true
				continue
			}

			tracker.Retries++
			tracker.FailedARNs = append(tracker.FailedARNs, aws.ToString(task.TaskArn))

			// Exit if we've failed too many times
			if tracker.Retries >= execCtx.MaxFailCount {

				// Stop all other running tasks
				reason := fmt.Sprintf(errMaxFailCount, tracker.ActiveARN, tracker.Retries, execCtx.MaxFailCount)
				if err := stopAllTasks(execCtx, reason); err != nil {
					return err
				}

				return fmt.Errorf("%s", reason)
			}

			newTaskARN, err := runTask(execCtx, tracker.Name, tracker.TaskNum)
			if err != nil {
				return err
			}

			// Assign the new task ARN to the tracker
			tracker.ActiveARN = newTaskARN

			// Task failed but will be restarted, so mark as not complete
			allComplete = false
			continue
		}

		// If all tasks are complete, break out of the loop
		if allComplete {
			break
		}

		// Check if we've timed out
		if time.Now().After(stopTime) {
			// Collect ARNs of tasks that did not complete
			var incompleteARNs []string
			for _, tracker := range execCtx.tasks {
				if !tracker.Completed {
					incompleteARNs = append(incompleteARNs, tracker.ActiveARN)
				}
			}

			// Stop all remaining tasks
			reason := fmt.Sprintf(errPollingTimeout, incompleteARNs, execCtx.Timeout)
			if err := stopAllTasks(execCtx, reason); err != nil {
				return err
			}

			// Return error with information about incomplete tasks
			return fmt.Errorf("%s", reason)
		}

		// Sleep until next poll time
		time.Sleep(time.Duration(execCtx.PollingInterval))
	}

	// If you're here, all tasks are complete
	return nil

}

func buildExecutionContext(commandCtx *ecsCommandContext, j *job.Job, c *cluster.Cluster) (*executionContext, error) {

	execCtx := &executionContext{
		tasks: make(map[string]*taskTracker),
	}

	// Create a context from commandCtx and unmarshal onto execCtx (defaults)
	commandContext := context.New(commandCtx)
	if err := commandContext.Unmarshal(execCtx); err != nil {
		return nil, err
	}

	// Overlay job context (overrides command values)
	if j.Context != nil {
		if err := j.Context.Unmarshal(execCtx); err != nil {
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
	execCtx.ClusterConfig = clusterContext

	// Load task definition template
	taskDefWrapper, err := loadTaskDefinitionTemplate(commandCtx.TaskDefinitionTemplate)
	if err != nil {
		return nil, err
	}
	execCtx.TaskDefinitionWrapper = taskDefWrapper // Store the wrapper for polling

	// Build container overrides for all containers
	if err := buildContainerOverrides(execCtx); err != nil {
		return nil, err
	}

	// Validate the resolved configuration
	if err := validateExecutionContext(execCtx); err != nil {
		return nil, err
	}

	// initialize AWS session
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	execCtx.ecsClient = ecs.NewFromConfig(cfg)

	return execCtx, nil

}

// validateExecutionContext validates the final resolved configuration
func validateExecutionContext(ctx *executionContext) error {

	if ctx.TaskCount <= 0 || ctx.TaskCount > ctx.ClusterConfig.MaxTaskCount {
		return fmt.Errorf("task count (%d) needs to be greater than 0 and less than cluster max task count (%d)", ctx.TaskCount, ctx.ClusterConfig.MaxTaskCount)
	}

	if ctx.CPU <= 0 || ctx.CPU > ctx.ClusterConfig.MaxCPU {
		return fmt.Errorf("cpu (%d) needs to be greater than 0 and less than or equal to cluster max cpu (%d)", ctx.CPU, ctx.ClusterConfig.MaxCPU)
	}

	if ctx.Memory <= 0 || ctx.Memory > ctx.ClusterConfig.MaxMemory {
		return fmt.Errorf("memory (%d) needs to be greater than 0 and less than or equal to cluster max memory (%d)", ctx.Memory, ctx.ClusterConfig.MaxMemory)
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
func stopAllTasks(execCtx *executionContext, reason string) error {

	for _, t := range execCtx.tasks {
		if t.Completed {
			continue
		}
		stopInput := &ecs.StopTaskInput{
			Cluster: aws.String(execCtx.ClusterConfig.ClusterName),
			Task:    aws.String(t.ActiveARN),
			Reason:  aws.String(reason),
		}

		_, err := execCtx.ecsClient.StopTask(ctx, stopInput)
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
func runTask(execCtx *executionContext, startedBy string, taskNum int) (string, error) {

	// Create a copy of the overrides and add TASK_NAME and TASK_NUM env variables
	finalOverrides := append([]types.ContainerOverride{}, execCtx.ContainerOverrides...)

	for i := range finalOverrides {
		finalOverrides[i].Environment = append(finalOverrides[i].Environment,
			types.KeyValuePair{
				Name:  aws.String("TASK_NAME"),
				Value: aws.String(startedBy),
			},
			types.KeyValuePair{
				Name:  aws.String("TASK_NUM"),
				Value: aws.String(fmt.Sprintf("%d", taskNum)),
			},
		)
	}

	// build run task input
	runTaskInput := &ecs.RunTaskInput{
		Cluster:        aws.String(execCtx.ClusterConfig.ClusterName),
		TaskDefinition: execCtx.taskDefARN,
		LaunchType:     types.LaunchType(execCtx.ClusterConfig.LaunchType),
		Count:          aws.Int32(1),
		StartedBy:      aws.String(startedBy),
		Overrides: &types.TaskOverride{
			ContainerOverrides: finalOverrides,
		},
		NetworkConfiguration: &types.NetworkConfiguration{
			AwsvpcConfiguration: &types.AwsVpcConfiguration{
				Subnets:        execCtx.ClusterConfig.VPCConfig.Subnets,
				SecurityGroups: execCtx.ClusterConfig.VPCConfig.SecurityGroups,
				AssignPublicIp: types.AssignPublicIpDisabled,
			},
		},
	}

	runTaskOutput, err := execCtx.ecsClient.RunTask(ctx, runTaskInput)
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

// storeResults builds and stores the final result for the job.
func storeResults(execCtx *executionContext, j *job.Job) error {

	// Build result
	j.Result = &result.Result{}
	j.Result.Columns = []*column.Column{
		{Name: "task_arn", Type: "string"},
		{Name: "duration", Type: "float"},
		{Name: "retries", Type: "int"},
		{Name: "failed_arns", Type: "string"},
	}

	// Create result data from task results
	j.Result.Data = make([][]interface{}, 0, len(execCtx.tasks))
	for _, tracker := range execCtx.tasks {
		j.Result.Data = append(j.Result.Data, []interface{}{
			tracker.ActiveARN,
			tracker.ExecutionTime,
			tracker.Retries,
			strings.Join(tracker.FailedARNs, ","),
		})
	}

	return nil

}
