package spark

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/emrcontainers"
	"github.com/aws/aws-sdk-go-v2/service/emrcontainers/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/babourine/x/pkg/set"

	heimdallContext "github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/job"
	"github.com/patterninc/heimdall/pkg/plugin"
	"github.com/patterninc/heimdall/pkg/result"
)

type sparkSubmitParameters struct {
	Properties map[string]string `yaml:"properties,omitempty" json:"properties,omitempty"`
	EntryPoint string            `yaml:"entry_point,omitempty" json:"entry_point,omitempty"`
}

// spark represents the Spark command context
type sparkCommandContext struct {
	QueriesURI string            `yaml:"queries_uri,omitempty" json:"queries_uri,omitempty"`
	ResultsURI string            `yaml:"results_uri,omitempty" json:"results_uri,omitempty"`
	LogsURI    *string           `yaml:"logs_uri,omitempty" json:"logs_uri,omitempty"`
	WrapperURI *string           `yaml:"wrapper_uri,omitempty" json:"wrapper_uri,omitempty"`
	Properties map[string]string `yaml:"properties,omitempty" json:"properties,omitempty"`
}

// sparkJobContext represents the context for a spark job
type sparkJobContext struct {
	Query        string                 `yaml:"query,omitempty" json:"query,omitempty"`
	Arguments    []string               `yaml:"arguments,omitempty" json:"arguments,omitempty"`
	Parameters   *sparkSubmitParameters `yaml:"parameters,omitempty" json:"parameters,omitempty"`
	ReturnResult bool                   `yaml:"return_result,omitempty" json:"return_result,omitempty"`
}

// sparkClusterContext represents the context for a spark cluster
type sparkClusterContext struct {
	ExecutionRoleArn *string           `yaml:"execution_role_arn,omitempty" json:"execution_role_arn,omitempty"`
	EMRReleaseLabel  *string           `yaml:"emr_release_label,omitempty" json:"emr_release_label,omitempty"`
	RoleARN          *string           `yaml:"role_arn,omitempty" json:"role_arn,omitempty"`
	Properties       map[string]string `yaml:"properties,omitempty" json:"properties,omitempty"`
}

const (
	driverMemoryProperty = `spark.driver.memory`
	jobCheckInterval     = 5                                // seconds
	jobTimeout           = (5 * 60 * 60) / jobCheckInterval // 5 hours
	noStateDetails       = `no state details provided`
)

var (
	sparkDefaults     = aws.String(`spark-defaults`)
	assumeRoleSession = aws.String("AssumeRoleSession")
	runtimeStates     = set.New([]types.JobRunState{types.JobRunStateCompleted, types.JobRunStateFailed, types.JobRunStateCancelled})
	rxS3              = regexp.MustCompile(`^s3://([^/]+)/(.*)$`)
)

var (
	ErrUnknownCluster = fmt.Errorf(`unknown cluster`)
	ErrJobCanceled    = fmt.Errorf(`job canceled`)
)

// New creates a new Spark plugin handler.
func New(commandContext *heimdallContext.Context) (plugin.Handler, error) {

	s := &sparkCommandContext{}

	if commandContext != nil {
		if err := commandContext.Unmarshal(s); err != nil {
			return nil, err
		}
	}

	return s.handler, nil

}

// Handler for the Spark job submission.
func (s *sparkCommandContext) handler(ctx context.Context, r *plugin.Runtime, j *job.Job, c *cluster.Cluster) (err error) {

	// let's unmarshal job context
	jobContext := &sparkJobContext{}
	if j.Context != nil {
		if err := j.Context.Unmarshal(jobContext); err != nil {
			return err
		}
	}

	// let's unmarshal cluster context
	clusterContext := &sparkClusterContext{}
	if c.Context != nil {
		if err := c.Context.Unmarshal(clusterContext); err != nil {
			return err
		}
	}

	if jobContext.Parameters == nil {
		jobContext.Parameters = &sparkSubmitParameters{}
	}

	// let's prepare job properties
	if jobContext.Parameters.Properties == nil {
		jobContext.Parameters.Properties = make(map[string]string)
	}
	for k, v := range s.Properties {
		if _, found := jobContext.Parameters.Properties[k]; !found {
			jobContext.Parameters.Properties[k] = v
		}
	}

	// do we have driver memory setting in the job properties?
	if value, found := jobContext.Parameters.Properties[driverMemoryProperty]; found {
		clusterContext.Properties[driverMemoryProperty] = value
		delete(jobContext.Parameters.Properties, driverMemoryProperty)
	}

	// setting AWS client
	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	// let's set empty options function,...
	assumeRoleOptions := func(_ *emrcontainers.Options) {}

	// ...and, if we have assume role ARN set, let's establish creds...
	if clusterContext.RoleARN != nil {

		stsSvc := sts.NewFromConfig(awsConfig)

		assumeRoleOutput, err := stsSvc.AssumeRole(ctx, &sts.AssumeRoleInput{
			RoleArn:         clusterContext.RoleARN,
			RoleSessionName: assumeRoleSession,
		})
		if err != nil {
			return err
		}

		assumeRoleOptions = func(o *emrcontainers.Options) {
			o.Credentials = credentials.NewStaticCredentialsProvider(
				*assumeRoleOutput.Credentials.AccessKeyId,
				*assumeRoleOutput.Credentials.SecretAccessKey,
				*assumeRoleOutput.Credentials.SessionToken,
			)
		}

	}

	svc := emrcontainers.NewFromConfig(awsConfig, assumeRoleOptions)

	// let's get the cluster ID
	clusterID, err := getClusterID(ctx, svc, c.Name)
	if err != nil {
		return err
	}

	// let's set the result uri
	resultURI := fmt.Sprintf("%s/%s", s.ResultsURI, j.ID)

	// upload query to s3 here...
	queryURI := fmt.Sprintf("%s/%s/query.sql", s.QueriesURI, j.ID)
	if err := uploadFileToS3(ctx, queryURI, jobContext.Query); err != nil {
		return err
	}

	// let's set job driver
	jobDriver := &types.JobDriver{}
	s.setJobDriver(jobContext, jobDriver, queryURI, resultURI)

	// let's prepare job payload
	jobPayload := &emrcontainers.StartJobRunInput{
		Name:             aws.String(j.ID),
		VirtualClusterId: clusterID,
		ExecutionRoleArn: clusterContext.ExecutionRoleArn,
		ReleaseLabel:     clusterContext.EMRReleaseLabel,
		JobDriver:        jobDriver,
		ConfigurationOverrides: &types.ConfigurationOverrides{
			ApplicationConfiguration: []types.Configuration{{
				Classification: sparkDefaults,
				Properties:     clusterContext.Properties,
			}},
			MonitoringConfiguration: &types.MonitoringConfiguration{
				PersistentAppUI: types.PersistentAppUIEnabled,
				S3MonitoringConfiguration: &types.S3MonitoringConfiguration{
					LogUri: s.LogsURI,
				},
			},
		},
	}

	// record the payload so we could easier understand what was submitted
	jobPayloadJSON, err := json.MarshalIndent(jobPayload, ``, `  `)
	if err != nil {
		return err
	}
	r.Stdout.WriteString(string(jobPayloadJSON) + "\n\n")

	// start the job
	outputStartJobRun, err := svc.StartJobRun(ctx, jobPayload)
	if err != nil {
		return err
	}

	// TODO: cleanup at some point, once the command is stable
	r.Stdout.WriteString(fmt.Sprintf("Cluster Job ID: %v\n", *outputStartJobRun.Id))
	// spew.Fdump(r.Stdout, s, clusterContext, jobContext)

	// keep checking until job succeeded or failed...
timeoutLoop:
	for i := 0; i < jobTimeout; i++ {
		time.Sleep(jobCheckInterval * time.Second)
		describeJobOutput, err := svc.DescribeJobRun(ctx, &emrcontainers.DescribeJobRunInput{
			Id:               outputStartJobRun.Id,
			VirtualClusterId: clusterID,
		})
		if err != nil {
			// TODO: log error if it's persistent
			r.Stderr.WriteString(fmt.Sprintf("job error: %v", err))
		}
		if describeJobOutput != nil {
			// print state every ~30 seconds...
			if state := describeJobOutput.JobRun.State; i%6 == 0 || runtimeStates.Has(state) {
				printState(r.Stdout, state)
				switch state {
				case types.JobRunStateCompleted:
					break timeoutLoop
				case types.JobRunStateFailed:
					stateDetails := noStateDetails
					if sd := describeJobOutput.JobRun.StateDetails; sd != nil {
						stateDetails = *sd
					}
					return fmt.Errorf("job failed [%v]: %v", describeJobOutput.JobRun.FailureReason, stateDetails)
				case types.JobRunStateCancelled:
					return ErrJobCanceled
				}
			}
		}
	}

	// TODO: set dummy result for now, return actual result...
	if j.Result, err = result.FromAvro(resultURI); err != nil {
		return err
	}

	return nil

}

func (s *sparkCommandContext) setJobDriver(jobContext *sparkJobContext, jobDriver *types.JobDriver, queryURI string, resultURI string) {
	jobParameters := getSparkSubmitParameters(jobContext)
	if jobContext.Arguments != nil {
		jobDriver.SparkSubmitJobDriver = &types.SparkSubmitJobDriver{
			EntryPoint:            s.WrapperURI,
			EntryPointArguments:   jobContext.Arguments,
			SparkSubmitParameters: jobParameters,
		}
		return
	}
	if jobContext.ReturnResult {
		jobDriver.SparkSubmitJobDriver = &types.SparkSubmitJobDriver{
			EntryPoint:            s.WrapperURI,
			EntryPointArguments:   []string{queryURI, resultURI},
			SparkSubmitParameters: jobParameters,
		}
		return
	}

	jobDriver.SparkSqlJobDriver = &types.SparkSqlJobDriver{
		EntryPoint:         &queryURI,
		SparkSqlParameters: jobParameters,
	}

}

func getClusterID(ctx context.Context, svc *emrcontainers.Client, clusterName string) (*string, error) {

	// let's get the cluster ID
	outputListClusters, err := svc.ListVirtualClusters(ctx, &emrcontainers.ListVirtualClustersInput{
		States:                []types.VirtualClusterState{types.VirtualClusterStateRunning},
		ContainerProviderType: types.ContainerProviderTypeEks,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list virtual clusters: %w", err)
	}

	for _, vc := range outputListClusters.VirtualClusters {
		if *vc.Name == clusterName {
			return vc.Id, nil
		}
	}

	return nil, fmt.Errorf("cluster %s: %w", clusterName, ErrUnknownCluster)

}

func getSparkSubmitParameters(context *sparkJobContext) *string {
	properties := context.Parameters.Properties
	conf := make([]string, 0, len(properties))

	for k, v := range properties {
		conf = append(conf, fmt.Sprintf("--conf %s=%s", k, v))
	}
	if context.Parameters.EntryPoint != "" {
		conf = append(conf, fmt.Sprintf("--class %s", context.Parameters.EntryPoint))
	}
	return aws.String(strings.Join(conf, ` `))

}

func printState(stdout *os.File, state types.JobRunState) {
	stdout.WriteString(fmt.Sprintf("%v - job is still running. latest status: %v\n", time.Now(), state))
}

func uploadFileToS3(ctx context.Context, fileURI, content string) error {

	// get bucket name and prefix
	s3Parts := rxS3.FindAllStringSubmatch(fileURI, -1)
	if len(s3Parts) == 0 || len(s3Parts[0]) < 3 {
		return fmt.Errorf("unexpected queries key: %v", s3Parts)
	}

	// upload file
	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	// Create an S3 client
	svc := s3.NewFromConfig(awsConfig)
	uploader := manager.NewUploader(svc)

	// Upload the string content to S3
	if _, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: &s3Parts[0][1],
		Key:    &s3Parts[0][2],
		Body:   strings.NewReader(content),
	}); err != nil {
		return err
	}

	return nil

}
