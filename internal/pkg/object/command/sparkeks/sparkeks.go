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
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"

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
	jobCheckInterval                      = 5 * time.Second
	jobTimeout                            = 5 * time.Hour
	statusReportInterval                  = 30 * time.Second
	defaultNamespace                      = "default"
	applicationPrefix                     = "spark-sql-job"
	awsRegionEnvVar                       = "AWS_REGION"
	sparkEventLogDirProperty              = "spark.eventLog.dir"
	sparkAppNameProperty                  = "spark.app.name"
	sparkKubernetesSubmitInDriverProperty = "spark.kubernetes.submitInDriver"
	s3Prefix                              = "s3://"
	s3aPrefix                             = "s3a://"
	stdoutLogSuffix                       = "stdout.log"
	stderrLogSuffix                       = "stderr.log"
	defaultSparkAppAPIVersion             = "sparkoperator.k8s.io/v1beta2"
	defaultSparkAppKind                   = "SparkApplication"
)

var (
	ctx               = context.Background()
	assumeRoleSession = aws.String("AssumeRoleSession")
	rxS3              = regexp.MustCompile(`^s3://([^/]+)/(.*)$`)
	runtimeStates     = []v1beta2.ApplicationStateType{
		v1beta2.ApplicationStateCompleted,
		v1beta2.ApplicationStateFailed,
		v1beta2.ApplicationStateFailedSubmission,
		v1beta2.ApplicationStateInvalidating,
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
	QueriesURI    string            `yaml:"queries_uri,omitempty" json:"queries_uri,omitempty"`
	ResultsURI    string            `yaml:"results_uri,omitempty" json:"results_uri,omitempty"`
	LogsURI       string            `yaml:"logs_uri,omitempty" json:"logs_uri,omitempty"`
	WrapperURI    string            `yaml:"wrapper_uri,omitempty" json:"wrapper_uri,omitempty"`
	Properties    map[string]string `yaml:"properties,omitempty" json:"properties,omitempty"`
	KubeNamespace string            `yaml:"kube_namespace,omitempty" json:"kube_namespace,omitempty"`
}

type sparkEksJobContext struct {
	Query        string            `yaml:"query,omitempty" json:"query,omitempty"`
	Properties   map[string]string `yaml:"properties,omitempty" json:"properties,omitempty"`
	ReturnResult bool              `yaml:"return_result,omitempty" json:"return_result,omitempty"`
}

type sparkEksClusterContext struct {
	RoleARN              *string           `yaml:"role_arn,omitempty" json:"role_arn,omitempty"`
	Properties           map[string]string `yaml:"properties,omitempty" json:"properties,omitempty"`
	Image                *string           `yaml:"image,omitempty" json:"image,omitempty"`
	Region               *string           `yaml:"region,omitempty" json:"region,omitempty"`
	SparkApplicationFile string            `yaml:"spark_application_file,omitempty" json:"spark_application_file,omitempty"`
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

// uploadFileToS3 uploads content to S3
func uploadFileToS3(fileURI, content string) error {
	// get bucket name and prefix
	s3Parts := rxS3.FindAllStringSubmatch(fileURI, -1)
	if len(s3Parts) == 0 || len(s3Parts[0]) < 3 {
		return fmt.Errorf("unexpected queries key: %v", s3Parts)
	}

	// upload file
	awsConfig, err := awsconfig.LoadDefaultConfig(ctx)
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

// updateS3ToS3aURI replaces s3:// with s3a:// for Hadoop FS compliance
func updateS3ToS3aURI(uri string) string {
	return strings.ReplaceAll(uri, s3Prefix, s3aPrefix)
}

// getS3FileURI finds a file in S3 directory that matches the given extension
func getS3FileURI(directoryURI, matchingExtension string) (string, error) {
	// Parse S3 URI to get bucket and prefix
	s3Parts := rxS3.FindAllStringSubmatch(directoryURI, -1)
	if len(s3Parts) == 0 || len(s3Parts[0]) < 3 {
		return "", fmt.Errorf("invalid S3 URI format: %s", directoryURI)
	}

	bucket := s3Parts[0][1]
	prefix := s3Parts[0][2]

	// Load AWS config
	awsConfig, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client
	svc := s3.NewFromConfig(awsConfig)

	// List objects in the directory
	listInput := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}

	result, err := svc.ListObjectsV2(ctx, listInput)
	if err != nil {
		return "", fmt.Errorf("failed to list S3 objects: %w", err)
	}

	// Find the first file that contains the matching extension
	for _, obj := range result.Contents {
		if strings.Contains(*obj.Key, matchingExtension) {
			return fmt.Sprintf("s3://%s/%s", bucket, *obj.Key), nil
		}
	}

	return "", fmt.Errorf("no file found with extension '%s' in directory: %s", matchingExtension, directoryURI)
}

// printState logs the current job status
func printState(stdout *os.File, state v1beta2.ApplicationStateType) {
	stdout.WriteString(fmt.Sprintf("%v - job is still running. latest status: %v\n", time.Now(), state))
}

// getSparkSubmitParameters returns Spark submit parameters as a string
func getSparkSubmitParameters(context *sparkEksJobContext) *string {
	if len(context.Properties) == 0 {
		return nil
	}
	properties := context.Properties
	conf := make([]string, 0, len(properties))
	for k, v := range properties {
		conf = append(conf, fmt.Sprintf("--conf %s=%s", k, v))
	}
	result := strings.Join(conf, ` `)
	return &result
}

// getSparkApplicationPods returns the list of pods associated with a Spark application
func getSparkApplicationPods(kubeClient *kubernetes.Clientset, appName, namespace string) ([]corev1.Pod, error) {
	labelSelector := fmt.Sprintf("sparkoperator.k8s.io/app-name=%s", appName)

	listOptions := metav1.ListOptions{
		LabelSelector: labelSelector,
	}

	podList, err := kubeClient.CoreV1().Pods(namespace).List(ctx, listOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to list pods for Spark application %s: %w", appName, err)
	}

	return podList.Items, nil
}

// getSparkApplicationPodLogs fetches logs from Spark application pods and uploads them to S3
func getSparkApplicationPodLogs(kubeClient *kubernetes.Clientset, pods []corev1.Pod, logURI string, r *plugin.Runtime) error {
	for _, pod := range pods {
		if pod.Status.Phase != corev1.PodRunning && pod.Status.Phase != corev1.PodSucceeded && pod.Status.Phase != corev1.PodFailed {
			continue
		}
		for _, container := range pod.Spec.Containers {
			stdoutLogs, err := getPodContainerLogs(kubeClient, pod.Name, pod.Namespace, container.Name, false)
			if err != nil {
				r.Stderr.WriteString(fmt.Sprintf("Pod %s, container %s: stdout log error: %v\n", pod.Name, container.Name, err))
				continue
			}
			if stdoutLogs != "" {
				stdoutURI := fmt.Sprintf("%s/%s-%s", logURI, pod.Name, stdoutLogSuffix)
				if err := uploadFileToS3(stdoutURI, stdoutLogs); err != nil {
					r.Stderr.WriteString(fmt.Sprintf("Pod %s, container %s: stdout upload error: %v\n", pod.Name, container.Name, err))
				}
			}
			stderrLogs, err := getPodContainerLogs(kubeClient, pod.Name, pod.Namespace, container.Name, true)
			if err != nil {
				r.Stderr.WriteString(fmt.Sprintf("Pod %s, container %s: stderr log error: %v\n", pod.Name, container.Name, err))
				continue
			}
			if stderrLogs != "" {
				stderrURI := fmt.Sprintf("%s/%s-%s", logURI, pod.Name, stderrLogSuffix)
				if err := uploadFileToS3(stderrURI, stderrLogs); err != nil {
					r.Stderr.WriteString(fmt.Sprintf("Pod %s, container %s: stderr upload error: %v\n", pod.Name, container.Name, err))
				}
			}
		}
	}
	return nil
}

// getPodContainerLogs is a helper function to fetch logs from a specific container in a pod
func getPodContainerLogs(kubeClient *kubernetes.Clientset, podName, namespace, containerName string, previous bool) (string, error) {
	logOptions := &corev1.PodLogOptions{
		Container: containerName,
		Previous:  previous,
	}

	req := kubeClient.CoreV1().Pods(namespace).GetLogs(podName, logOptions)
	logs, err := req.Stream(ctx)
	if err != nil {
		return "", err
	}
	defer logs.Close()

	var logContent strings.Builder
	buffer := make([]byte, 4096)

	for {
		n, err := logs.Read(buffer)
		if n > 0 {
			logContent.Write(buffer[:n])
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
	}

	return logContent.String(), nil
}

// handler executes the Spark EKS job submission and execution.
func (s *sparkEksCommandContext) handler(r *plugin.Runtime, j *job.Job, c *cluster.Cluster) error {
	// Parse contexts
	jobContext := &sparkEksJobContext{}
	if j.Context != nil {
		if err := j.Context.Unmarshal(jobContext); err != nil {
			return fmt.Errorf("failed to unmarshal job context: %w", err)
		}
	}

	clusterContext := &sparkEksClusterContext{}
	if c.Context != nil {
		if err := c.Context.Unmarshal(clusterContext); err != nil {
			return fmt.Errorf("failed to unmarshal cluster context: %w", err)
		}
	}

	// Initialize and merge properties
	if jobContext.Properties == nil {
		jobContext.Properties = make(map[string]string)
	}
	for k, v := range s.Properties {
		if _, found := jobContext.Properties[k]; !found {
			jobContext.Properties[k] = v
		}
	}
	if clusterContext.Properties == nil {
		clusterContext.Properties = make(map[string]string)
	}

	// Setup AWS configuration
	region := os.Getenv(awsRegionEnvVar)
	if clusterContext.Region != nil {
		region = *clusterContext.Region
	}

	var awsConfig aws.Config
	var err error
	if region != "" {
		awsConfig, err = awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(region))
	} else {
		awsConfig, err = awsconfig.LoadDefaultConfig(ctx)
	}
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Handle role assumption if specified
	if clusterContext.RoleARN != nil {
		stsSvc := sts.NewFromConfig(awsConfig)
		assumeRoleOutput, err := stsSvc.AssumeRole(ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(*clusterContext.RoleARN),
			RoleSessionName: assumeRoleSession,
		})
		if err != nil {
			return fmt.Errorf("failed to assume role %s: %w", *clusterContext.RoleARN, err)
		}

		awsConfig.Credentials = credentials.NewStaticCredentialsProvider(
			*assumeRoleOutput.Credentials.AccessKeyId,
			*assumeRoleOutput.Credentials.SecretAccessKey,
			*assumeRoleOutput.Credentials.SessionToken,
		)
	}

	// Create spark operator clients
	sparkClient, kubeClient, err := s.createSparkClient(r, c.Name, region, clusterContext.Region)
	if err != nil {
		return fmt.Errorf("failed to create Spark Operator client: %w", err)
	}

	// Upload query to S3
	queryURI := fmt.Sprintf("%s/%s/query.sql", s.QueriesURI, j.ID)
	if err := uploadFileToS3(queryURI, jobContext.Query); err != nil {
		return fmt.Errorf("failed to upload query to S3: %w", err)
	}

	// Create log URI directory if LogsURI is specified
	logURI := fmt.Sprintf("%s/%s/.keepdir", s.LogsURI, j.ID)
	// Create a zero-byte file to establish the directory structure in S3
	if err := uploadFileToS3(logURI, ""); err != nil {
		return fmt.Errorf("failed to create log directory in S3: %w", err)
	}

	logURI = strings.TrimSuffix(logURI, "/.keepdir")

	// Generate and submit Spark application
	resultURI := fmt.Sprintf("%s/%s", s.ResultsURI, j.ID)
	appName := fmt.Sprintf("%s-%s", applicationPrefix, j.ID)

	sparkApp, err := s.generateSparkApp(r, clusterContext, jobContext, appName, queryURI, resultURI, logURI, s.WrapperURI)
	if err != nil {
		return fmt.Errorf("failed to generate Spark application: %w", err)
	}

	// record the payload so we could easier understand what was submitted
	sparkAppJSON, err := json.MarshalIndent(sparkApp, ``, `  `)
	if err != nil {
		return err
	}
	r.Stdout.WriteString(string(sparkAppJSON) + "\n\n")

	// Log Spark submit parameters if available
	if sparkSubmitParams := getSparkSubmitParameters(jobContext); sparkSubmitParams != nil {
		r.Stdout.WriteString(fmt.Sprintf("Spark Submit Parameters: %s\n", *sparkSubmitParams))
	}

	submittedApp, err := sparkClient.SparkoperatorV1beta2().SparkApplications(sparkApp.Namespace).Create(ctx, sparkApp, metav1.CreateOptions{})
	if err != nil {
		r.Stderr.WriteString(fmt.Sprintf("Failed to submit Spark application: %v\n", err))
		return ErrJobSubmission
	}

	r.Stdout.WriteString(fmt.Sprintf("Successfully submitted Spark application: %s/%s\n", submittedApp.Namespace, submittedApp.Name))

	// Monitor the job until completion and collect logs before cleanup
	monitorErr := s.monitorJobAndCollectLogs(r, sparkClient, kubeClient, submittedApp.Name, submittedApp.Namespace, logURI)

	// Always attempt cleanup, even if monitorJobAndCollectLogs returns error
	cleanupErr := sparkClient.SparkoperatorV1beta2().SparkApplications(submittedApp.Namespace).Delete(ctx, submittedApp.Name, metav1.DeleteOptions{})
	if cleanupErr != nil {
		r.Stderr.WriteString(fmt.Sprintf("Warning: failed to cleanup application %s: %v\n", submittedApp.Name, cleanupErr))
	} else {
		r.Stdout.WriteString(fmt.Sprintf("Cleaned up Spark application: %s/%s\n", submittedApp.Namespace, submittedApp.Name))
	}

	if monitorErr != nil {
		return fmt.Errorf("failed to monitor Spark job: %w", monitorErr)
	}

	// If job should return results, fetch from S3
	if jobContext.ReturnResult {
		// Get the actual .avro file URI from the results directory
		returnResultFileURI, err := getS3FileURI(resultURI, "avro")
		if err != nil {
			r.Stdout.WriteString(fmt.Sprintf("failed to find .avro file in results directory %s: %s", resultURI, err))
			return fmt.Errorf("failed to find .avro file in results directory %s: %w", resultURI, err)
		}

		if j.Result, err = result.FromAvro(returnResultFileURI); err != nil {
			return fmt.Errorf("failed to fetch results from S3: %w", err)
		}
	}

	return nil
}

// createSparkClient creates Kubernetes and Spark clients for the EKS cluster
func (s *sparkEksCommandContext) createSparkClient(r *plugin.Runtime, clusterName, defaultRegion string, clusterRegion *string) (*sparkClientSet.Clientset, *kubernetes.Clientset, error) {
	region := defaultRegion
	if clusterRegion != nil {
		region = *clusterRegion
	}

	// Update kubeconfig using AWS CLI
	cmd := exec.Command("aws", "eks", "update-kubeconfig", "--region", region, "--name", clusterName)
	if err := cmd.Run(); err != nil {
		return nil, nil, fmt.Errorf("failed to update kubeconfig: please ensure AWS CLI is installed and configured, and you have access to the EKS cluster: %w", err)
	}

	// Load from updated kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load kubeconfig after update: %w", err)
	}

	// Create Spark Operator client
	sparkClient, err := sparkClientSet.NewForConfig(config)
	if err != nil {
		r.Stderr.WriteString(fmt.Sprintf("Failed to create Spark Operator client: %v\n", err))
		return nil, nil, ErrKubeConfig
	}

	// Create Kubernetes client
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		r.Stderr.WriteString(fmt.Sprintf("Failed to create Kubernetes client: %v\n", err))
		return nil, nil, ErrKubeConfig
	}

	r.Stdout.WriteString(fmt.Sprintf("Successfully created Spark Operator and Kubernetes clients for cluster: %s\n", clusterName))
	return sparkClient, kubeClient, nil
}

// applySparkOperatorConfig consolidates all Spark Operator configuration updates and overrides
func applySparkOperatorConfig(sparkApp *v1beta2.SparkApplication, r *plugin.Runtime, appName, queryURI, resultURI, logURI, wrapperURI string, jobContext *sparkEksJobContext, clusterContext *sparkEksClusterContext, kubeNamespace string) {
	// Set application name and namespace
	sparkApp.ObjectMeta.Name = appName
	sparkApp.ObjectMeta.Namespace = kubeNamespace

	// Set main application file and arguments
	if wrapperURI != "" {
		s3aWrapperURI := updateS3ToS3aURI(wrapperURI)
		s3aQueryURI := updateS3ToS3aURI(queryURI)
		s3aResultURI := updateS3ToS3aURI(resultURI)
		mainAppFile := s3aWrapperURI
		if jobContext.ReturnResult {
			sparkApp.Spec.Arguments = []string{appName, s3aQueryURI, s3aResultURI}
		} else {
			sparkApp.Spec.Arguments = []string{appName, s3aQueryURI}
		}
		sparkApp.Spec.MainApplicationFile = &mainAppFile
	}

	// Ensure SparkConf is initialized
	if sparkApp.Spec.SparkConf == nil {
		sparkApp.Spec.SparkConf = make(map[string]string)
	}

	// Set consistent app name for both driver and executor pods
	sparkApp.Spec.SparkConf[sparkAppNameProperty] = appName

	// Set log URI
	if logURI != "" {
		logURI = updateS3ToS3aURI(logURI)
		sparkApp.Spec.SparkConf[sparkEventLogDirProperty] = logURI
	}

	// Add Spark submit parameters if available
	if sparkSubmitParams := getSparkSubmitParameters(jobContext); sparkSubmitParams != nil {
		if sparkApp.Spec.SparkConf == nil {
			sparkApp.Spec.SparkConf = make(map[string]string)
		}
		if _, exists := sparkApp.Spec.SparkConf[sparkKubernetesSubmitInDriverProperty]; !exists {
			sparkApp.Spec.SparkConf[sparkKubernetesSubmitInDriverProperty] = "false"
		}
		r.Stdout.WriteString(fmt.Sprintf("Applied Spark submit parameters: %s\n", *sparkSubmitParams))
	}

	// Set container image
	if clusterContext.Image != nil {
		sparkApp.Spec.Image = clusterContext.Image
	}

	// Set environment variables
	if clusterContext.Region != nil {
		if sparkApp.Spec.Driver.EnvVars == nil {
			sparkApp.Spec.Driver.EnvVars = make(map[string]string)
		}
		sparkApp.Spec.Driver.EnvVars[awsRegionEnvVar] = *clusterContext.Region
	}

	// Configure Spark resources and properties (driver/executor)
	// Driver resources
	if driverCores := jobContext.Properties["spark.driver.cores"]; driverCores != "" {
		if cores, err := strconv.ParseInt(driverCores, 10, 32); err == nil {
			cores32 := int32(cores)
			sparkApp.Spec.Driver.Cores = &cores32
		}
		delete(jobContext.Properties, "spark.driver.cores")
	}
	if driverMemory := jobContext.Properties["spark.driver.memory"]; driverMemory != "" {
		sparkApp.Spec.Driver.Memory = &driverMemory
		delete(jobContext.Properties, "spark.driver.memory")
	}
	// Executor resources
	if executorInstances := jobContext.Properties["spark.executor.instances"]; executorInstances != "" {
		if instances, err := strconv.ParseInt(executorInstances, 10, 32); err == nil {
			instances32 := int32(instances)
			sparkApp.Spec.Executor.Instances = &instances32
		}
		delete(jobContext.Properties, "spark.executor.instances")
	}
	if executorCores := jobContext.Properties["spark.executor.cores"]; executorCores != "" {
		if cores, err := strconv.ParseInt(executorCores, 10, 32); err == nil {
			cores32 := int32(cores)
			sparkApp.Spec.Executor.Cores = &cores32
		}
		delete(jobContext.Properties, "spark.executor.cores")
	}
	if executorMemory := jobContext.Properties["spark.executor.memory"]; executorMemory != "" {
		sparkApp.Spec.Executor.Memory = &executorMemory
		delete(jobContext.Properties, "spark.executor.memory")
	}

	// Add cluster properties to sparkConf
	for k, v := range clusterContext.Properties {
		sparkApp.Spec.SparkConf[k] = v
	}
	// Add remaining job properties to sparkConf
	for k, v := range jobContext.Properties {
		sparkApp.Spec.SparkConf[k] = v
	}
}

// generateSparkApp generates the Spark application from the template
func (s *sparkEksCommandContext) generateSparkApp(r *plugin.Runtime, clusterContext *sparkEksClusterContext, jobContext *sparkEksJobContext, appName, queryURI, resultURI, logURI, wrapperURI string) (*v1beta2.SparkApplication, error) {
	// Load and parse template
	sparkApp, err := s.loadTemplate(r, clusterContext.SparkApplicationFile)
	if err != nil {
		return nil, err
	}

	// Consolidated configuration (now also sets namespace)
	applySparkOperatorConfig(sparkApp, r, appName, queryURI, resultURI, logURI, wrapperURI, jobContext, clusterContext, s.KubeNamespace)
	return sparkApp, nil
}

// loadTemplate loads and parses the YAML template file
func (s *sparkEksCommandContext) loadTemplate(r *plugin.Runtime, sparkApplicationFile string) (*v1beta2.SparkApplication, error) {
	if sparkApplicationFile == "" {
		sparkApp := &v1beta2.SparkApplication{
			Spec: v1beta2.SparkApplicationSpec{
				Driver:   v1beta2.DriverSpec{},
				Executor: v1beta2.ExecutorSpec{},
			},
		}
		// Set defaults if not present
		sparkApp.APIVersion = defaultSparkAppAPIVersion
		sparkApp.Kind = defaultSparkAppKind
		return sparkApp, nil
	}

	r.Stdout.WriteString(fmt.Sprintf("Loading template file: %s\n", sparkApplicationFile))

	yamlContent, err := os.ReadFile(sparkApplicationFile)
	if err != nil {
		r.Stderr.WriteString(fmt.Sprintf("Failed to read template file '%s': %v\n", sparkApplicationFile, err))
		return nil, ErrSparkApplicationFile
	}

	var sparkApp v1beta2.SparkApplication
	if err := yaml.Unmarshal(yamlContent, &sparkApp); err != nil {
		r.Stderr.WriteString(fmt.Sprintf("Failed to parse YAML template file '%s': %v\n", sparkApplicationFile, err))
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

// monitorJobAndCollectLogs monitors the Spark job until completion and collects logs before cleanup
func (s *sparkEksCommandContext) monitorJobAndCollectLogs(r *plugin.Runtime, sparkClient *sparkClientSet.Clientset, kubeClient *kubernetes.Clientset, appName, namespace, logURI string) error {
	r.Stdout.WriteString(fmt.Sprintf("Monitoring Spark application: %s\n", appName))

	timeout := time.Now().Add(jobTimeout)
	lastReport := time.Now()
	var finalSparkApp *v1beta2.SparkApplication

	for time.Now().Before(timeout) {
		time.Sleep(jobCheckInterval)

		sparkApp, err := sparkClient.SparkoperatorV1beta2().SparkApplications(namespace).Get(ctx, appName, metav1.GetOptions{})
		if err != nil {
			// If the SparkApplication is deleted manually, attempt cleanup and return error
			r.Stderr.WriteString(fmt.Sprintf("Failed to get application status: %v\n", err))
			_ = sparkClient.SparkoperatorV1beta2().SparkApplications(namespace).Delete(ctx, appName, metav1.DeleteOptions{})
			return fmt.Errorf("spark application not found or deleted: %w", err)
		}

		finalSparkApp = sparkApp
		state := sparkApp.Status.AppState.State

		if time.Since(lastReport) >= statusReportInterval || s.isTerminalState(string(state)) {
			printState(r.Stdout, state)
			lastReport = time.Now()
		}

		if state == v1beta2.ApplicationStateRunning {
			s.collectSparkApplicationLogs(r, kubeClient, sparkApp, logURI)
		}

		switch state {
		case v1beta2.ApplicationStateCompleted, v1beta2.ApplicationStateFailed:
			s.collectSparkApplicationLogs(r, kubeClient, sparkApp, logURI)
			// Attempt cleanup before returning
			_ = sparkClient.SparkoperatorV1beta2().SparkApplications(namespace).Delete(ctx, appName, metav1.DeleteOptions{})
			if state == v1beta2.ApplicationStateCompleted {
				r.Stdout.WriteString("Spark job completed successfully\n")
				return nil
			} else {
				errorMessage := sparkApp.Status.AppState.ErrorMessage
				if errorMessage == "" {
					errorMessage = "Unknown error"
				}
				r.Stderr.WriteString(fmt.Sprintf("Spark job failed: %s\n", errorMessage))
				return fmt.Errorf("spark job failed: %s", errorMessage)
			}
		case v1beta2.ApplicationStateFailedSubmission:
			if finalSparkApp != nil {
				s.collectSparkApplicationLogs(r, kubeClient, finalSparkApp, logURI)
			}
			_ = sparkClient.SparkoperatorV1beta2().SparkApplications(namespace).Delete(ctx, appName, metav1.DeleteOptions{})
			return fmt.Errorf("spark job submission failed")
		case v1beta2.ApplicationStateInvalidating:
			if finalSparkApp != nil {
				s.collectSparkApplicationLogs(r, kubeClient, finalSparkApp, logURI)
			}
			_ = sparkClient.SparkoperatorV1beta2().SparkApplications(namespace).Delete(ctx, appName, metav1.DeleteOptions{})
			return ErrJobCanceled
		}
	}

	// Timeout case - try to collect logs and cleanup
	if finalSparkApp != nil {
		s.collectSparkApplicationLogs(r, kubeClient, finalSparkApp, logURI)
	}
	_ = sparkClient.SparkoperatorV1beta2().SparkApplications(namespace).Delete(ctx, appName, metav1.DeleteOptions{})
	return fmt.Errorf("spark job timed out after %v", jobTimeout)
}

// collectSparkApplicationLogs collects logs from Spark application pods
func (s *sparkEksCommandContext) collectSparkApplicationLogs(r *plugin.Runtime, kubeClient *kubernetes.Clientset, sparkApp *v1beta2.SparkApplication, logURI string) {
	pods, err := getSparkApplicationPods(kubeClient, sparkApp.Name, sparkApp.Namespace)
	if err != nil {
		r.Stderr.WriteString(fmt.Sprintf("Warning: failed to get Spark application pods: %v\n", err))
		return
	}

	if err := getSparkApplicationPodLogs(kubeClient, pods, logURI, r); err != nil {
		r.Stderr.WriteString(fmt.Sprintf("Warning: failed to collect pod logs: %v\n", err))
	}
}

// isTerminalState checks if the given state is a terminal state
func (s *sparkEksCommandContext) isTerminalState(state string) bool {
	appState := v1beta2.ApplicationStateType(state)
	for _, terminalState := range runtimeStates {
		if appState == terminalState {
			return true
		}
	}
	return false
}
