package sparkeks

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/kubeflow/spark-operator/v2/api/v1beta2"
	sparkClientSet "github.com/kubeflow/spark-operator/v2/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"

	heimdallContext "github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/job"
	"github.com/patterninc/heimdall/pkg/plugin"
	"github.com/patterninc/heimdall/pkg/result"
)

const (
	jobCheckInterval          = 5 * time.Second
	jobTimeout                = 5 * time.Hour
	statusReportInterval      = 30 * time.Second
	defaultNamespace          = "default"
	applicationPrefix         = "spark-sql-job"
	awsRegionEnvVar           = "AWS_REGION"
	s3Prefix                  = "s3://"
	s3aPrefix                 = "s3a://"
	stdoutLogSuffix           = "stdout.log"
	stderrLogSuffix           = "stderr.log"
	defaultSparkAppAPIVersion = "sparkoperator.k8s.io/v1beta2"
	defaultSparkAppKind       = "SparkApplication"

	queriesPath = "queries"
	resultsPath = "results"
	logsPath    = "logs"

	queryFileName     = "query.sql"
	avroFileExtension = "avro"

	// Spark configuration related constants
	sparkEventLogDirProperty              = "spark.eventLog.dir"
	sparkAppNameProperty                  = "spark.app.name"
	sparkKubernetesSubmitInDriverProperty = "spark.kubernetes.submitInDriver"
	sparkDriverCoresKey                   = "spark.driver.cores"
	sparkDriverMemoryKey                  = "spark.driver.memory"
	sparkExecutorInstancesKey             = "spark.executor.instances"
	sparkExecutorCoresKey                 = "spark.executor.cores"
	sparkExecutorMemoryKey                = "spark.executor.memory"
	sparkAppLabelSelectorFormat           = "sparkoperator.k8s.io/app-name"

	unknownErrorMsg             = "Unknown error"
	sparkJobSubmissionFailedMsg = "spark job submission failed"
	sparkAppUnknownStateMsg     = "spark app entered unknown state"
)

var (
	ctx           = context.Background()
	rxS3          = regexp.MustCompile(`^s3://([^/]+)/(.*)$`)
	runtimeStates = []v1beta2.ApplicationStateType{
		v1beta2.ApplicationStateCompleted,
		v1beta2.ApplicationStateFailed,
		v1beta2.ApplicationStateFailedSubmission,
		v1beta2.ApplicationStateUnknown,
	}
)

var (
	ErrJobCanceled          = fmt.Errorf("job was canceled before completion")
	ErrJobSubmission        = fmt.Errorf("failed to submit Spark application to Kubernetes cluster")
	ErrKubeConfig           = fmt.Errorf("failed to configure Kubernetes client: ensure EKS cluster access is properly configured")
	ErrApplicationSpec      = fmt.Errorf("failed to load or parse SparkApplication template")
	ErrSparkApplicationFile = fmt.Errorf("failed to read SparkApplication application template file: check file path and permissions")
)

type sparkEksCommandContext struct {
	JobsURI       string            `yaml:"jobs_uri,omitempty" json:"jobs_uri,omitempty"`
	WrapperURI    string            `yaml:"wrapper_uri,omitempty" json:"wrapper_uri,omitempty"`
	Properties    map[string]string `yaml:"properties,omitempty" json:"properties,omitempty"`
	KubeNamespace string            `yaml:"kube_namespace,omitempty" json:"kube_namespace,omitempty"`
}

type sparkEksJobParameters struct {
	Properties map[string]string `yaml:"properties,omitempty" json:"properties,omitempty"`
}

type sparkEksJobContext struct {
	Query        string                 `yaml:"query,omitempty" json:"query,omitempty"`
	Parameters   *sparkEksJobParameters `yaml:"parameters,omitempty" json:"parameters,omitempty"`
	ReturnResult bool                   `yaml:"return_result,omitempty" json:"return_result,omitempty"`
}

type sparkEksClusterContext struct {
	RoleARN              *string           `yaml:"role_arn,omitempty" json:"role_arn,omitempty"`
	Properties           map[string]string `yaml:"properties,omitempty" json:"properties,omitempty"`
	Image                *string           `yaml:"image,omitempty" json:"image,omitempty"`
	Region               *string           `yaml:"region,omitempty" json:"region,omitempty"`
	SparkApplicationFile string            `yaml:"spark_application_file,omitempty" json:"spark_application_file,omitempty"`
}

// executionContext holds the final resolved configuration and clients for a job execution.
type executionContext struct {
	runtime        *plugin.Runtime
	job            *job.Job
	cluster        *cluster.Cluster
	commandContext *sparkEksCommandContext
	jobContext     *sparkEksJobContext
	clusterContext *sparkEksClusterContext

	sparkClient *sparkClientSet.Clientset
	kubeClient  *kubernetes.Clientset

	awsConfig aws.Config

	appName      string
	queryURI     string
	resultURI    string
	logURI       string
	sparkApp     *v1beta2.SparkApplication
	submittedApp *v1beta2.SparkApplication
}

// New creates a new Spark EKS plugin handler.
func New(commandContext *heimdallContext.Context) (plugin.Handler, error) {
	s := &sparkEksCommandContext{
		KubeNamespace: defaultNamespace,
	}

	if commandContext != nil {
		if err := commandContext.Unmarshal(s); err != nil {
			return nil, err
		}
	}

	return s.handler, nil
}

// handler executes the Spark EKS job submission and execution.
func (s *sparkEksCommandContext) handler(r *plugin.Runtime, j *job.Job, c *cluster.Cluster) error {
	// 1. Build execution context, create URIs, and upload query
	execCtx, err := buildExecutionContextAndURI(r, j, c, s)
	if err != nil {
		return fmt.Errorf("failed to build execution context: %w", err)
	}

	// 2. Submit the Spark Application to the cluster
	if err := execCtx.submitSparkApp(); err != nil {
		return err
	}

	// 3. Monitor the job until completion and collect logs
	monitorErr := execCtx.monitorJobAndCollectLogs()

	// 4. Cleanup any resources that are still pending
	if err := execCtx.cleanupSparkApp(); err != nil {
		// Log cleanup error but don't override the main monitoring error
		execCtx.runtime.Stderr.WriteString(fmt.Sprintf("Warning: failed to cleanup application %s: %v\n", execCtx.submittedApp.Name, err))
	}

	// If monitoring failed, return that error after attempting cleanup
	if monitorErr != nil {
		return fmt.Errorf("failed to monitor Spark job: %w", monitorErr)
	}

	// 5. Get and store results if required
	if err := execCtx.getAndStoreResults(); err != nil {
		return err
	}

	return nil
}

// buildExecutionContextAndURI prepares the context, merges configurations, and uploads the query.
func buildExecutionContextAndURI(r *plugin.Runtime, j *job.Job, c *cluster.Cluster, s *sparkEksCommandContext) (*executionContext, error) {
	execCtx := &executionContext{
		runtime:        r,
		job:            j,
		cluster:        c,
		commandContext: s,
	}

	// Parse job context
	jobContext := &sparkEksJobContext{}
	if j.Context != nil {
		if err := j.Context.Unmarshal(jobContext); err != nil {
			return nil, fmt.Errorf("failed to unmarshal job context: %w", err)
		}
	}
	execCtx.jobContext = jobContext

	// Parse cluster context
	clusterContext := &sparkEksClusterContext{}
	if c.Context != nil {
		if err := c.Context.Unmarshal(clusterContext); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cluster context: %w", err)
		}
	}
	execCtx.clusterContext = clusterContext

	// Initialize and merge properties from command -> job
	if execCtx.jobContext.Parameters == nil {
		execCtx.jobContext.Parameters = &sparkEksJobParameters{
			Properties: make(map[string]string),
		}
	}
	if execCtx.jobContext.Parameters.Properties == nil {
		execCtx.jobContext.Parameters.Properties = make(map[string]string)
	}
	for k, v := range s.Properties {
		if _, found := execCtx.jobContext.Parameters.Properties[k]; !found {
			execCtx.jobContext.Parameters.Properties[k] = v
		}
	}
	if execCtx.clusterContext.Properties == nil {
		execCtx.clusterContext.Properties = make(map[string]string)
	}

	// Load AWS config once and store in execCtx
	awsConfig, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}
	execCtx.awsConfig = awsConfig

	// Set URIs and App Name
	execCtx.appName = fmt.Sprintf("%s-%s", applicationPrefix, j.ID)
	execCtx.queryURI = fmt.Sprintf("%s/%s/%s/%s", s.JobsURI, j.ID, queriesPath, queryFileName)
	execCtx.resultURI = fmt.Sprintf("%s/%s/%s", s.JobsURI, j.ID, resultsPath)
	execCtx.logURI = fmt.Sprintf("%s/%s/%s", s.JobsURI, j.ID, logsPath)

	// Upload query to S3
	if err := uploadFileToS3(execCtx.awsConfig, execCtx.queryURI, execCtx.jobContext.Query); err != nil {
		return nil, fmt.Errorf("failed to upload query to S3: %w", err)
	}

	return execCtx, nil
}

// submitSparkApp creates clients, generates the spec, and submits it to Kubernetes.
func (e *executionContext) submitSparkApp() error {
	// Create Kubernetes and Spark Operator clients
	if err := createSparkClients(e); err != nil {
		return fmt.Errorf("failed to create Spark Operator client: %w", err)
	}

	// Generate Spark application spec from template and context
	if err := generateSparkApp(e); err != nil {
		return fmt.Errorf("failed to generate Spark application: %w", err)
	}

	// Record the payload for debugging
	sparkAppJSON, err := json.MarshalIndent(e.sparkApp, ``, `  `)
	if err != nil {
		return err
	}
	e.runtime.Stdout.WriteString(string(sparkAppJSON) + "\n\n")

	// Log Spark submit parameters if available
	if sparkSubmitParams := getSparkSubmitParameters(e.jobContext); sparkSubmitParams != nil {
		e.runtime.Stdout.WriteString(fmt.Sprintf("Spark Submit Parameters: %s\n", *sparkSubmitParams))
	}

	// Submit the application
	submittedApp, err := e.sparkClient.SparkoperatorV1beta2().SparkApplications(e.sparkApp.Namespace).Create(ctx, e.sparkApp, metav1.CreateOptions{})
	if err != nil {
		e.runtime.Stderr.WriteString(fmt.Sprintf("Failed to submit Spark application: %v\n", err))
		return ErrJobSubmission
	}
	e.submittedApp = submittedApp

	e.runtime.Stdout.WriteString(fmt.Sprintf("Successfully submitted Spark application: %s/%s\n", submittedApp.Namespace, submittedApp.Name))
	return nil
}

// cleanupSparkApp removes the SparkApplication from the cluster if it still exists.
func (e *executionContext) cleanupSparkApp() error {
	if e.submittedApp == nil {
		return nil
	}

	name, namespace := e.submittedApp.Name, e.submittedApp.Namespace
	_, getErr := e.sparkClient.SparkoperatorV1beta2().SparkApplications(namespace).Get(ctx, name, metav1.GetOptions{})
	if getErr == nil {
		// Application exists, proceed with deletion
		cleanupErr := e.sparkClient.SparkoperatorV1beta2().SparkApplications(namespace).Delete(ctx, name, metav1.DeleteOptions{})
		if cleanupErr != nil && !strings.Contains(cleanupErr.Error(), "not found") {
			return cleanupErr // Return error if deletion fails for reasons other than "not found"
		}
		e.runtime.Stdout.WriteString(fmt.Sprintf("Cleaned up Spark application: %s/%s\n", namespace, name))
	}
	return nil // "not found" is a successful cleanup state
}

// getAndStoreResults fetches the job output from S3 and stores it.
func (e *executionContext) getAndStoreResults() error {
	if !e.jobContext.ReturnResult {
		return nil
	}

	returnResultFileURI, err := getS3FileURI(e.awsConfig, e.resultURI, avroFileExtension)
	if err != nil {
		e.runtime.Stdout.WriteString(fmt.Sprintf("failed to find .avro file in results directory %s: %s", e.resultURI, err))
		return fmt.Errorf("failed to find .avro file in results directory %s: %w", e.resultURI, err)
	}

	if e.job.Result, err = result.FromAvro(returnResultFileURI); err != nil {
		return fmt.Errorf("failed to fetch results from S3: %w", err)
	}

	return nil
}

// uploadFileToS3 uploads content to S3.
func uploadFileToS3(awsConfig aws.Config, fileURI, content string) error {
	s3Parts := rxS3.FindAllStringSubmatch(fileURI, -1)
	if len(s3Parts) == 0 || len(s3Parts[0]) < 3 {
		return fmt.Errorf("unexpected S3 URI format: %s", fileURI)
	}
	bucket, key := s3Parts[0][1], s3Parts[0][2]

	uploader := manager.NewUploader(s3.NewFromConfig(awsConfig))

	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   strings.NewReader(content),
	})
	return err
}

// updateS3ToS3aURI replaces s3:// with s3a:// for Hadoop FS compliance.
func updateS3ToS3aURI(uri string) string {
	return strings.ReplaceAll(uri, s3Prefix, s3aPrefix)
}

// getS3FileURI finds a file in an S3 directory that matches the given extension.
func getS3FileURI(awsConfig aws.Config, directoryURI, matchingExtension string) (string, error) {
	s3Parts := rxS3.FindAllStringSubmatch(directoryURI, -1)
	if len(s3Parts) == 0 || len(s3Parts[0]) < 3 {
		return "", fmt.Errorf("invalid S3 URI format: %s", directoryURI)
	}
	bucket, prefix := s3Parts[0][1], s3Parts[0][2]

	svc := s3.NewFromConfig(awsConfig)

	result, err := svc.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return "", fmt.Errorf("failed to list S3 objects: %w", err)
	}

	for _, obj := range result.Contents {
		if strings.Contains(*obj.Key, matchingExtension) {
			return fmt.Sprintf("s3://%s/%s", bucket, *obj.Key), nil
		}
	}

	return "", fmt.Errorf("no file found with extension '%s' in directory: %s", matchingExtension, directoryURI)
}

// printState logs the current job status.
func printState(writer io.Writer, state v1beta2.ApplicationStateType) {
	fmt.Fprintf(writer, "%v - job is still running. latest status: %v\n", time.Now(), state)
}

// getSparkSubmitParameters returns Spark submit parameters as a string.
func getSparkSubmitParameters(context *sparkEksJobContext) *string {
	if context.Parameters == nil || len(context.Parameters.Properties) == 0 {
		return nil
	}
	conf := make([]string, 0, len(context.Parameters.Properties))
	for k, v := range context.Parameters.Properties {
		conf = append(conf, fmt.Sprintf("--conf %s=%s", k, v))
	}
	result := strings.Join(conf, ` `)
	return &result
}

// getSparkApplicationPods returns the list of pods associated with a Spark application.
func getSparkApplicationPods(kubeClient *kubernetes.Clientset, appName, namespace string) ([]corev1.Pod, error) {
	labelSelector := fmt.Sprintf("%s=%s", sparkAppLabelSelectorFormat, appName)
	podList, err := kubeClient.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods for Spark application %s: %w", appName, err)
	}
	return podList.Items, nil
}

// isPodInValidPhase returns true if the pod is in a running, succeeded, or failed phase.
func isPodInValidPhase(pod corev1.Pod) bool {
	return pod.Status.Phase == corev1.PodRunning ||
		pod.Status.Phase == corev1.PodSucceeded ||
		pod.Status.Phase == corev1.PodFailed
}

// writeLastLinesToStderr writes the last 100 lines of log content to stderr if it's from a driver pod
func writeLastLinesToStderr(execCtx *executionContext, pod corev1.Pod, logContent string) {
	// Check if this is a driver pod by looking for "driver" in the pod name
	if !strings.Contains(pod.Name, "driver") {
		return
	}

	lines := strings.Split(strings.TrimSpace(logContent), "\n")
	if len(lines) == 0 {
		return
	}

	// Get last 100 lines or all lines if fewer than 100
	startIndex := 0
	if len(lines) > 100 {
		startIndex = len(lines) - 100
	}

	lastLines := lines[startIndex:]
	if len(lastLines) > 0 {
		execCtx.runtime.Stderr.WriteString(fmt.Sprintf("\n=== Last %d lines from driver pod %s ===\n", len(lastLines), pod.Name))
		execCtx.runtime.Stderr.WriteString(strings.Join(lastLines, "\n"))
		execCtx.runtime.Stderr.WriteString("\n=== End of driver logs ===\n\n")
	}
}

// getAndUploadPodContainerLogs fetches logs from a specific container in a pod and uploads them to S3.
func getAndUploadPodContainerLogs(execCtx *executionContext, pod corev1.Pod, container corev1.Container, previous bool, logType string, writeToStderr bool) {
	logOptions := &corev1.PodLogOptions{Container: container.Name, Previous: previous}
	req := execCtx.kubeClient.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, logOptions)
	logs, err := req.Stream(ctx)
	if err != nil {
		return
	}
	defer logs.Close()

	logContent, err := io.ReadAll(logs)
	if err != nil {
		return
	}
	if string(logContent) != "" {
		// Write to stderr if requested and this is stdout logs (not previous/stderr logs)
		if writeToStderr && !previous && logType == stdoutLogSuffix {
			writeLastLinesToStderr(execCtx, pod, string(logContent))
		}

		logURI := fmt.Sprintf("%s/%s-%s", execCtx.logURI, pod.Name, logType)
		if err := uploadFileToS3(execCtx.awsConfig, logURI, string(logContent)); err != nil {
			execCtx.runtime.Stderr.WriteString(fmt.Sprintf("Pod %s, container %s: %s upload error: %v\n", pod.Name, container.Name, logType, err))
		}
	}
}

// getSparkApplicationPodLogs fetches logs from pods and uploads them to S3.
func getSparkApplicationPodLogs(execCtx *executionContext, pods []corev1.Pod, writeToStderr bool) error {
	for _, pod := range pods {
		if !isPodInValidPhase(pod) {
			continue
		}
		for _, container := range pod.Spec.Containers {
			// Get current logs and upload
			getAndUploadPodContainerLogs(execCtx, pod, container, false, stdoutLogSuffix, writeToStderr)
			// Get logs from previous (failed) runs and upload
			getAndUploadPodContainerLogs(execCtx, pod, container, true, stderrLogSuffix, false)
		}
	}
	return nil
}

// createSparkClients creates Kubernetes and Spark clients for the EKS cluster.
func createSparkClients(execCtx *executionContext) error {
	region := os.Getenv(awsRegionEnvVar)
	if execCtx.clusterContext.Region != nil {
		region = *execCtx.clusterContext.Region
	}

	// Update kubeconfig using AWS CLI
	cmd := exec.Command("aws", "eks", "update-kubeconfig", "--region", region, "--name", execCtx.cluster.Name)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to update kubeconfig: ensure AWS CLI is installed and configured: %w", err)
	}

	// Load from updated kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		return fmt.Errorf("failed to load kubeconfig after update: %w", err)
	}

	// Create Spark Operator client
	sparkClient, err := sparkClientSet.NewForConfig(config)
	if err != nil {
		execCtx.runtime.Stderr.WriteString(fmt.Sprintf("Failed to create Spark Operator client: %v\n", err))
		return ErrKubeConfig
	}
	execCtx.sparkClient = sparkClient

	// Create Kubernetes client
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		execCtx.runtime.Stderr.WriteString(fmt.Sprintf("Failed to create Kubernetes client: %v\n", err))
		return ErrKubeConfig
	}
	execCtx.kubeClient = kubeClient

	execCtx.runtime.Stdout.WriteString(fmt.Sprintf("Successfully created Spark Operator and Kubernetes clients for cluster: %s\n", execCtx.cluster.Name))
	return nil
}

// applySparkOperatorConfig consolidates all Spark Operator configuration updates and overrides.
func applySparkOperatorConfig(execCtx *executionContext) {
	sparkApp := execCtx.sparkApp
	jobContext := execCtx.jobContext
	clusterContext := execCtx.clusterContext

	// Set application name and namespace
	sparkApp.ObjectMeta.Name = execCtx.appName
	sparkApp.ObjectMeta.Namespace = execCtx.commandContext.KubeNamespace

	// Set main application file and arguments
	if execCtx.commandContext.WrapperURI != "" {
		s3aWrapperURI := updateS3ToS3aURI(execCtx.commandContext.WrapperURI)
		s3aQueryURI := updateS3ToS3aURI(execCtx.queryURI)
		s3aResultURI := updateS3ToS3aURI(execCtx.resultURI)
		mainAppFile := s3aWrapperURI
		if jobContext.ReturnResult {
			sparkApp.Spec.Arguments = []string{execCtx.appName, s3aQueryURI, execCtx.job.User, s3aResultURI}
		} else {
			sparkApp.Spec.Arguments = []string{execCtx.appName, s3aQueryURI, execCtx.job.User}
		}
		sparkApp.Spec.MainApplicationFile = &mainAppFile
	}

	if sparkApp.Spec.SparkConf == nil {
		sparkApp.Spec.SparkConf = make(map[string]string)
	}

	sparkApp.Spec.SparkConf[sparkAppNameProperty] = execCtx.appName

	if execCtx.logURI != "" {
		logURI := updateS3ToS3aURI(execCtx.logURI)
		sparkApp.Spec.SparkConf[sparkEventLogDirProperty] = logURI
	}

	if sparkSubmitParams := getSparkSubmitParameters(jobContext); sparkSubmitParams != nil {
		if _, exists := sparkApp.Spec.SparkConf[sparkKubernetesSubmitInDriverProperty]; !exists {
			sparkApp.Spec.SparkConf[sparkKubernetesSubmitInDriverProperty] = "false"
		}
	}

	if clusterContext.Image != nil {
		sparkApp.Spec.Image = clusterContext.Image
	}

	if clusterContext.Region != nil {
		if sparkApp.Spec.Driver.EnvVars == nil {
			sparkApp.Spec.Driver.EnvVars = make(map[string]string)
		}
		sparkApp.Spec.Driver.EnvVars[awsRegionEnvVar] = *clusterContext.Region
	}

	// Driver and Executor resources are handled by deleting from job properties after use
	// to avoid them being added to sparkConf directly.
	if driverCores := jobContext.Parameters.Properties[sparkDriverCoresKey]; driverCores != "" {
		if cores, err := strconv.ParseInt(driverCores, 10, 32); err == nil {
			cores32 := int32(cores)
			sparkApp.Spec.Driver.Cores = &cores32
		}
		delete(jobContext.Parameters.Properties, sparkDriverCoresKey)
	}
	if driverMemory := jobContext.Parameters.Properties[sparkDriverMemoryKey]; driverMemory != "" {
		sparkApp.Spec.Driver.Memory = &driverMemory
		delete(jobContext.Parameters.Properties, sparkDriverMemoryKey)
	}
	if executorInstances := jobContext.Parameters.Properties[sparkExecutorInstancesKey]; executorInstances != "" {
		if instances, err := strconv.ParseInt(executorInstances, 10, 32); err == nil {
			instances32 := int32(instances)
			sparkApp.Spec.Executor.Instances = &instances32
		}
		delete(jobContext.Parameters.Properties, sparkExecutorInstancesKey)
	}
	if executorCores := jobContext.Parameters.Properties[sparkExecutorCoresKey]; executorCores != "" {
		if cores, err := strconv.ParseInt(executorCores, 10, 32); err == nil {
			cores32 := int32(cores)
			sparkApp.Spec.Executor.Cores = &cores32
		}
		delete(jobContext.Parameters.Properties, sparkExecutorCoresKey)
	}
	if executorMemory := jobContext.Parameters.Properties[sparkExecutorMemoryKey]; executorMemory != "" {
		sparkApp.Spec.Executor.Memory = &executorMemory
		delete(jobContext.Parameters.Properties, sparkExecutorMemoryKey)
	}

	for k, v := range clusterContext.Properties {
		sparkApp.Spec.SparkConf[k] = v
	}
	for k, v := range jobContext.Parameters.Properties {
		sparkApp.Spec.SparkConf[k] = v
	}
}

// generateSparkApp generates the Spark application from the template.
func generateSparkApp(execCtx *executionContext) error {
	// Load and parse template
	sparkApp, err := loadTemplate(execCtx)
	if err != nil {
		return err
	}
	execCtx.sparkApp = sparkApp

	// Apply all configurations and overrides from contexts
	applySparkOperatorConfig(execCtx)

	return nil
}

// loadTemplate loads and parses the YAML template file.
func loadTemplate(execCtx *executionContext) (*v1beta2.SparkApplication, error) {
	sparkApplicationFile := execCtx.clusterContext.SparkApplicationFile
	if sparkApplicationFile == "" {
		sparkApp := &v1beta2.SparkApplication{
			Spec: v1beta2.SparkApplicationSpec{
				Driver:   v1beta2.DriverSpec{},
				Executor: v1beta2.ExecutorSpec{},
			},
		}
		sparkApp.APIVersion = defaultSparkAppAPIVersion
		sparkApp.Kind = defaultSparkAppKind
		return sparkApp, nil
	}

	execCtx.runtime.Stdout.WriteString(fmt.Sprintf("Loading template file: %s\n", sparkApplicationFile))
	yamlContent, err := os.ReadFile(sparkApplicationFile)
	if err != nil {
		execCtx.runtime.Stderr.WriteString(fmt.Sprintf("Failed to read template file '%s': %v\n", sparkApplicationFile, err))
		return nil, ErrSparkApplicationFile
	}

	var sparkApp v1beta2.SparkApplication
	if err := yaml.Unmarshal(yamlContent, &sparkApp); err != nil {
		execCtx.runtime.Stderr.WriteString(fmt.Sprintf("Failed to parse YAML template file '%s': %v\n", sparkApplicationFile, err))
		return nil, ErrApplicationSpec
	}

	if sparkApp.APIVersion == "" {
		sparkApp.APIVersion = defaultSparkAppAPIVersion
	}
	if sparkApp.Kind == "" {
		sparkApp.Kind = defaultSparkAppKind
	}

	return &sparkApp, nil
}

// monitorJobAndCollectLogs monitors the Spark job until completion and collects logs.
func (e *executionContext) monitorJobAndCollectLogs() error {
	appName, namespace := e.submittedApp.Name, e.submittedApp.Namespace
	e.runtime.Stdout.WriteString(fmt.Sprintf("Monitoring Spark application: %s\n", appName))

	monitorCtx, cancel := context.WithTimeout(ctx, jobTimeout)
	defer cancel()
	lastReport := time.Now()
	var finalSparkApp *v1beta2.SparkApplication

	for {
		if monitorCtx.Err() != nil {
			if finalSparkApp != nil {
				collectSparkApplicationLogs(e, finalSparkApp, true)
			}
			if monitorCtx.Err() == context.DeadlineExceeded {
				return fmt.Errorf("spark job timed out after %v", jobTimeout)
			}
			return monitorCtx.Err()
		}

		time.Sleep(jobCheckInterval)

		sparkApp, err := e.sparkClient.SparkoperatorV1beta2().SparkApplications(namespace).Get(monitorCtx, appName, metav1.GetOptions{})
		if err != nil {
			e.runtime.Stderr.WriteString(fmt.Sprintf("Spark application %s/%s not found or deleted externally: %v\n", namespace, appName, err))
			if finalSparkApp != nil {
				collectSparkApplicationLogs(e, finalSparkApp, true)
			}
			return fmt.Errorf("spark application %s/%s not found: %w", namespace, appName, err)
		}

		finalSparkApp = sparkApp
		state := sparkApp.Status.AppState.State

		if time.Since(lastReport) >= statusReportInterval || isTerminalState(state) {
			printState(e.runtime.Stdout, state)
			lastReport = time.Now()
		}

		if state == v1beta2.ApplicationStateRunning {
			collectSparkApplicationLogs(e, sparkApp, false)
			continue
		}

		switch state {
		case v1beta2.ApplicationStateCompleted:
			collectSparkApplicationLogs(e, sparkApp, false)
			e.runtime.Stdout.WriteString("Spark job completed successfully\n")
			return nil
		case v1beta2.ApplicationStateFailed:
			collectSparkApplicationLogs(e, sparkApp, true)
			errorMessage := sparkApp.Status.AppState.ErrorMessage
			if errorMessage == "" {
				errorMessage = unknownErrorMsg
			}
			e.runtime.Stderr.WriteString(fmt.Sprintf("Spark job failed: %s\n", errorMessage))
			return fmt.Errorf("spark job failed: %s", errorMessage)
		case v1beta2.ApplicationStateFailedSubmission, v1beta2.ApplicationStateUnknown:
			collectSparkApplicationLogs(e, finalSparkApp, true)
			msg := sparkJobSubmissionFailedMsg
			if state == v1beta2.ApplicationStateUnknown {
				msg = sparkAppUnknownStateMsg
			}
			return fmt.Errorf("%s", msg)
		}
	}
}

// collectSparkApplicationLogs collects logs from Spark application pods.
func collectSparkApplicationLogs(execCtx *executionContext, sparkApp *v1beta2.SparkApplication, writeToStderr bool) {
	if sparkApp == nil {
		return
	}
	pods, err := getSparkApplicationPods(execCtx.kubeClient, sparkApp.Name, sparkApp.Namespace)
	if err != nil {
		execCtx.runtime.Stderr.WriteString(fmt.Sprintf("Warning: failed to get Spark application pods: %v\n", err))
		return
	}

	if err := getSparkApplicationPodLogs(execCtx, pods, writeToStderr); err != nil {
		execCtx.runtime.Stderr.WriteString(fmt.Sprintf("Warning: failed to collect pod logs: %v\n", err))
	}
}

// isTerminalState checks if the given state is a terminal state.
func isTerminalState(state v1beta2.ApplicationStateType) bool {
	for _, terminalState := range runtimeStates {
		if state == terminalState {
			return true
		}
	}
	return false
}
