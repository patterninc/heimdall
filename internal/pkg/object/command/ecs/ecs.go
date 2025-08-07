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
	hc "github.com/patterninc/heimdall/pkg/context"
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
	PollingInterval        int                       `yaml:"polling_interval,omitempty" json:"polling_interval,omitempty"` // in seconds
	Timeout                int                       `yaml:"timeout,omitempty" json:"timeout,omitempty"`                   // in seconds
	MaxFailCount           int                       `yaml:"max_fail_count,omitempty" json:"max_fail_count,omitempty"`     // max failures before giving up
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

const (
	defaultPollingInterval = 30   // seconds
	defaultTaskTimeout     = 3600 // seconds (1 hour)
	defaultMaxFailCount    = 1
	defaultTaskCount       = 1
	startedByPrefix        = "heimdall-job-"
	errMaxFailCount        = "task %s failed %d times (max: %d), giving up"
	errPollingTimeout      = "polling timed out after after %v"
)

var (
	ctx                = ct.Background()
	errMissingTemplate = fmt.Errorf("task definition template is required")
)

func New(commandContext *hc.Context) (plugin.Handler, error) {

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

	// Prepare job context with command defaults
	jobContext := &ecsJobContext{}
	if j.Context != nil {
		if err := j.Context.Unmarshal(jobContext); err != nil {
			return err
		}
	}

	// fill in command defaults where job context values are not specified
	if len(jobContext.ContainerOverrides) == 0 {
		jobContext.ContainerOverrides = e.ContainerOverrides
	}

	// unmarshal cluster context
	clusterContext := &ecsClusterContext{}
	if c.Context != nil {
		if err := c.Context.Unmarshal(clusterContext); err != nil {
			return err
		}
	}

	// validate task count
	if jobContext.TaskCount > clusterContext.MaxTaskCount || jobContext.TaskCount <= 0 {
		return fmt.Errorf("task count (%d) needs to be above 0 and less than or equal to the cluster max task count (%d)",
			jobContext.TaskCount, clusterContext.MaxTaskCount)
	}

	// initialize AWS session
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	ecsClient := ecs.NewFromConfig(cfg)

	// load task definition template
	taskDefWrapper, err := loadTaskDefinitionTemplate(e.TaskDefinitionTemplate)
	if err != nil {
		return err
	}

	// Build container overrides for all containers
	if err := buildContainerOverrides(taskDefWrapper.TaskDefinition, jobContext); err != nil {
		return err
	}

	// register the task definition
	registerInput := &ecs.RegisterTaskDefinitionInput{
		Family:                  aws.String(aws.ToString(taskDefWrapper.TaskDefinition.Family)),
		RequiresCompatibilities: []types.Compatibility{types.CompatibilityFargate},
		NetworkMode:             types.NetworkModeAwsvpc,
		Cpu:                     aws.String(fmt.Sprintf("%d", clusterContext.CPU)),
		Memory:                  aws.String(fmt.Sprintf("%d", clusterContext.Memory)),
		ExecutionRoleArn:        aws.String(clusterContext.ExecutionRoleARN),
		TaskRoleArn:             aws.String(clusterContext.TaskRoleARN),
		ContainerDefinitions:    taskDefWrapper.TaskDefinition.ContainerDefinitions,
	}

	registerOutput, err := ecsClient.RegisterTaskDefinition(ctx, registerInput)
	if err != nil {
		return err
	}

	// Start tracking tasks
	tasks := make(map[string]*taskTracker)
	startTime := time.Now()

	for i := 0; i < jobContext.TaskCount; i++ {
		taskARN, err := runTask(ecsClient, registerOutput.TaskDefinition, clusterContext, jobContext.ContainerOverrides, fmt.Sprintf("%s%s-%d", startedByPrefix, j.ID, i))
		if err != nil {
			return err
		}
		taskName := fmt.Sprintf("%s%s-%d", startedByPrefix, j.ID, i)
		tasks[taskName] = &taskTracker{
			Name:      taskName,
			ActiveARN: taskARN,
		}
	}

	// Poll until all tasks are complete or timeout
	for time.Since(startTime) < time.Duration(e.Timeout)*time.Second {

		// Describe the uncompleted tasks we're tracking
		var activeARNs []string
		for _, tracker := range tasks {
			if !tracker.Completed {
				activeARNs = append(activeARNs, tracker.ActiveARN)
			}
		}

		describeInput := &ecs.DescribeTasksInput{
			Cluster: aws.String(clusterContext.ClusterName),
			Tasks:   activeARNs,
		}

		describeOutput, err := ecsClient.DescribeTasks(ctx, describeInput)
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
			tracker, exists := tasks[aws.ToString(task.StartedBy)]
			if !exists {
				return fmt.Errorf("could not find tracker for StartedBy tag %s", aws.ToString(task.StartedBy))
			}

			// Check for task failures based on exit code
			if !isTaskSuccessful(task, taskDefWrapper) {

				tracker.Retries++
				tracker.FailedARNs = append(tracker.FailedARNs, aws.ToString(task.TaskArn))

				// Exit if we've failed too many times
				if tracker.Retries >= e.MaxFailCount {
					// Stop all other running tasks
					if err := stopAllTasks(ecsClient, clusterContext.ClusterName, tasks); err != nil {
						return err
					}

					return fmt.Errorf(errMaxFailCount, tracker.Name, tracker.Retries+1, e.MaxFailCount)
				}

				newTaskARN, err := runTask(ecsClient, registerOutput.TaskDefinition, clusterContext, jobContext.ContainerOverrides, tracker.Name)
				if err != nil {
					return err
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

		// Wait before next poll
		time.Sleep(time.Duration(e.PollingInterval) * time.Second)

	}

	// If you're here you've either timed out or all tasks are complete
	// Create result data from task results

	j.Result = &result.Result{}
	j.Result.Columns = []*column.Column{
		{Name: "task_arn", Type: "string"},
		{Name: "duration", Type: "float"},
		{Name: "retries", Type: "int"},
		{Name: "failed_arns", Type: "string"},
	}

	j.Result.Data = make([][]interface{}, 0, len(tasks))
	for _, tracker := range tasks {
		j.Result.Data = append(j.Result.Data, []interface{}{
			tracker.ActiveARN,
			tracker.ExecutionTime,
			tracker.Retries,
			strings.Join(tracker.FailedARNs, ","),
		})
	}

	// Return error if we timed out
	if time.Since(startTime) >= time.Duration(e.Timeout)*time.Second {
		if err := stopAllTasks(ecsClient, clusterContext.ClusterName, tasks); err != nil {
			return err
		}

		return fmt.Errorf(errPollingTimeout, time.Duration(e.Timeout)*time.Second)
	}

	return nil

}

func buildContainerOverrides(taskDef *types.TaskDefinition, jobContext *ecsJobContext) error {

	// Create a map of container names
	existingContainers := make(map[string]bool)
	for _, container := range taskDef.ContainerDefinitions {
		existingContainers[aws.ToString(container.Name)] = true
	}

	// Create a map of jobContext overrides
	containerOverridesMap := make(map[string]types.ContainerOverride)
	for _, containerOverride := range jobContext.ContainerOverrides {
		// Validate that the container name exists in the task definition template
		if !existingContainers[aws.ToString(containerOverride.Name)] {
			return fmt.Errorf("container override '%s' not found in task definition template", aws.ToString(containerOverride.Name))
		}
		containerOverridesMap[aws.ToString(containerOverride.Name)] = containerOverride
	}

	// Build container overrides for all containers in the task definition
	var containerOverrides []types.ContainerOverride
	for _, container := range taskDef.ContainerDefinitions {
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

	jobContext.ContainerOverrides = containerOverrides

	return nil

}

// stopAllTasks stops all non-completed tasks
func stopAllTasks(ecsClient *ecs.Client, clusterName string, tasks map[string]*taskTracker) error {

	for _, t := range tasks {
		if !t.Completed {
			stopInput := &ecs.StopTaskInput{
				Cluster: aws.String(clusterName),
				Task:    aws.String(t.ActiveARN),
			}

			_, err := ecsClient.StopTask(ctx, stopInput)
			if err != nil {
				return err
			}
		}
	}

	return nil

}

func loadTaskDefinitionTemplate(templatePath string) (*taskDefinitionWrapper, error) {

	if templatePath == "" {
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
func runTask(ecsClient *ecs.Client, taskDef *types.TaskDefinition, clusterContext *ecsClusterContext, containerOverrides []types.ContainerOverride, startedBy string) (string, error) {

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
		TaskDefinition: taskDef.TaskDefinitionArn,
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

	if len(runTaskOutput.Tasks) == 0 {
		return ``, fmt.Errorf("no tasks were created")
	}

	taskARN := aws.ToString(runTaskOutput.Tasks[0].TaskArn)

	return taskARN, nil

}

// tasks are successful if all essential containers exit with a zero exit code
func isTaskSuccessful(task types.Task, taskDefWrapper *taskDefinitionWrapper) bool {
	// Check all containers in the running task
	for _, container := range task.Containers {
		containerName := aws.ToString(container.Name)

		if essential, exists := taskDefWrapper.EssentialContainers[containerName]; exists && essential {
			if container.ExitCode != nil && *container.ExitCode != 0 {
				return false
			}
		}
	}

	return true

}
