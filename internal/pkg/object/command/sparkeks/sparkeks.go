package sparkeks

import (
	ct "context"
	"fmt"
	"os"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	"github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/job"
	"github.com/patterninc/heimdall/pkg/plugin"
	"github.com/patterninc/heimdall/pkg/result"
	"github.com/patterninc/heimdall/pkg/result/column"
)

// sparkEksCommandContext represents the Spark EKS command context
type sparkEksCommandContext struct {
	QueriesURI    string            `yaml:"queries_uri,omitempty" json:"queries_uri,omitempty"`
	ResultsURI    string            `yaml:"results_uri,omitempty" json:"results_uri,omitempty"`
	LogsURI       *string           `yaml:"logs_uri,omitempty" json:"logs_uri,omitempty"`
	Properties    map[string]string `yaml:"properties,omitempty" json:"properties,omitempty"`
	KubeNamespace string            `yaml:"kube_namespace,omitempty" json:"kube_namespace,omitempty"`
}

// sparkEksJobContext represents the context for a spark job
type sparkEksJobContext struct {
	Query        string            `yaml:"query,omitempty" json:"query,omitempty"`
	Properties   map[string]string `yaml:"properties,omitempty" json:"properties,omitempty"`
	ReturnResult bool              `yaml:"return_result,omitempty" json:"return_result,omitempty"`
}

// sparkEksClusterContext represents the context for a spark cluster
type sparkEksClusterContext struct {
	ExecutionRoleArn *string           `yaml:"execution_role_arn,omitempty" json:"execution_role_arn,omitempty"`
	RoleARN          *string           `yaml:"role_arn,omitempty" json:"role_arn,omitempty"`
	Properties       map[string]string `yaml:"properties,omitempty" json:"properties,omitempty"`
	Image            *string           `yaml:"image,omitempty" json:"image,omitempty"`
	Region           *string           `yaml:"region,omitempty" json:"region,omitempty"`
	TemplateFile     string            `yaml:"template_file,omitempty" json:"template_file,omitempty"`
}

// SparkApplication represents the structure of the Spark Kubernetes application
type SparkApplication struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
	Spec struct {
		Type                string `yaml:"type"`
		Mode                string `yaml:"mode"`
		Image               string `yaml:"image"`
		ImagePullPolicy     string `yaml:"imagePullPolicy"`
		MainApplicationFile string `yaml:"mainApplicationFile"`
		SparkVersion        string `yaml:"sparkVersion"`
		RestartPolicy       struct {
			Type string `yaml:"type"`
		} `yaml:"restartPolicy"`
		TimeToLiveSeconds int `yaml:"timeToLiveSeconds"`
		Driver            struct {
			Cores          int               `yaml:"cores"`
			Memory         string            `yaml:"memory"`
			ServiceAccount string            `yaml:"serviceAccount"`
			EnvVars        map[string]string `yaml:"envVars"`
		} `yaml:"driver"`
		Executor struct {
			Instances int    `yaml:"instances"`
			Cores     int    `yaml:"cores"`
			Memory    string `yaml:"memory"`
		} `yaml:"executor"`
		SparkConf map[string]string `yaml:"sparkConf"`
		Deps      struct {
			Packages []string `yaml:"packages"`
		} `yaml:"deps"`
		Arguments []string `yaml:"arguments,omitempty"`
	} `yaml:"spec"`
}

const (
	jobCheckInterval    = 5                                // seconds
	jobTimeout          = (5 * 60 * 60) / jobCheckInterval // 5 hours
	defaultNamespace    = "default"
	applicationPrefix   = "spark-sql-job"
	defaultTemplateFile = "configs/spark-application.yaml"
	wrapperScript       = "s3://pattern-dl-ops/contrib/spark/spark-sql-s3-wrapper.py"
)

var (
	ctx               = ct.Background()
	assumeRoleSession = aws.String("AssumeRoleSession")
	rxS3              = regexp.MustCompile(`^s3://([^/]+)/(.*)$`)
)

var (
	ErrUnknownCluster  = fmt.Errorf(`unknown cluster`)
	ErrJobCanceled     = fmt.Errorf(`job canceled`)
	ErrJobSubmission   = fmt.Errorf(`job submission failed`)
	ErrKubeConfig      = fmt.Errorf(`failed to update kubeconfig`)
	ErrApplicationSpec = fmt.Errorf(`failed to load application template`)
	ErrTemplateFile    = fmt.Errorf(`failed to read template file`)
)

// New creates a new Spark EKS plugin handler.
func New(commandContext *context.Context) (plugin.Handler, error) {
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

// Handler for the Spark EKS connectivity and cluster listing.
func (s *sparkEksCommandContext) handler(r *plugin.Runtime, j *job.Job, c *cluster.Cluster) (err error) {
	// Unmarshal job context
	jobContext := &sparkEksJobContext{}
	if j.Context != nil {
		if err := j.Context.Unmarshal(jobContext); err != nil {
			return err
		}
	}

	// Unmarshal cluster context
	clusterContext := &sparkEksClusterContext{}
	if c.Context != nil {
		if err := c.Context.Unmarshal(clusterContext); err != nil {
			return err
		}
	}

	// Initialize clusterContext.Properties if nil
	if clusterContext.Properties == nil {
		clusterContext.Properties = make(map[string]string)
	}

	// Load AWS configuration with region from cluster context if specified
	var awsConfig aws.Config
	if clusterContext.Region != nil {
		awsConfig, err = config.LoadDefaultConfig(ctx, config.WithRegion(*clusterContext.Region))
	} else {
		// Use AWS_REGION environment variable if set, otherwise use default
		if region := os.Getenv("AWS_REGION"); region != "" {
			awsConfig, err = config.LoadDefaultConfig(ctx, config.WithRegion(region))
		} else {
			awsConfig, err = config.LoadDefaultConfig(ctx)
		}
	}
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Handle role assumption if specified
	if clusterContext.RoleARN != nil {
		stsSvc := sts.NewFromConfig(awsConfig)
		assumeRoleOutput, err := stsSvc.AssumeRole(ctx, &sts.AssumeRoleInput{
			RoleArn:         clusterContext.RoleARN,
			RoleSessionName: assumeRoleSession,
		})
		if err != nil {
			return fmt.Errorf("failed to assume role %s: %w", *clusterContext.RoleARN, err)
		}

		// Create new config with assumed role credentials
		awsConfig.Credentials = credentials.NewStaticCredentialsProvider(
			*assumeRoleOutput.Credentials.AccessKeyId,
			*assumeRoleOutput.Credentials.SecretAccessKey,
			*assumeRoleOutput.Credentials.SessionToken,
		)
	}

	// Create EKS client
	eksClient := eks.NewFromConfig(awsConfig)

	// Create result data with cluster information
	clusterData := make([]map[string]any, 0)

	// Describe the specific "spark-eks" cluster
	clusterName := "spark-eks"
	describeOutput, err := eksClient.DescribeCluster(ctx, &eks.DescribeClusterInput{
		Name: &clusterName,
	})
	if err != nil {
		return fmt.Errorf("failed to describe EKS cluster %s: %w", clusterName, err)
	}

	cluster := describeOutput.Cluster
	clusterInfo := map[string]any{
		"name":      clusterName,
		"status":    string(cluster.Status),
		"version":   aws.ToString(cluster.Version),
		"endpoint":  aws.ToString(cluster.Endpoint),
		"createdAt": cluster.CreatedAt,
	}

	if cluster.RoleArn != nil {
		clusterInfo["roleArn"] = *cluster.RoleArn
	}

	clusterData = append(clusterData, clusterInfo)

	// If the job should return results, create a result
	if jobContext.ReturnResult {
		// Create columns for the result
		columns := []*column.Column{
			{Name: "cluster_name", Type: column.Type("string")},
			{Name: "status", Type: column.Type("string")},
			{Name: "version", Type: column.Type("string")},
			{Name: "endpoint", Type: column.Type("string")},
			{Name: "role_arn", Type: column.Type("string")},
			{Name: "created_at", Type: column.Type("string")},
		}

		// Create rows for the result
		data := make([][]any, 0)
		for _, cluster := range clusterData {
			roleArn := ""
			if cluster["roleArn"] != nil {
				roleArn = fmt.Sprintf("%v", cluster["roleArn"])
			}

			row := []any{
				cluster["name"],
				cluster["status"],
				cluster["version"],
				cluster["endpoint"],
				roleArn,
				fmt.Sprintf("%v", cluster["createdAt"]),
			}
			data = append(data, row)
		}

		// Create the result
		j.Result = &result.Result{
			Columns: columns,
			Data:    data,
		}
	}

	return nil
}

// Comment out unused functions for now
/*
// updateKubeConfig updates the kubectl configuration to access the EKS cluster
func updateKubeConfig(r *plugin.Runtime, clusterName, region string) error {
	cmd := exec.Command("aws", "eks", "update-kubeconfig", "--region", region, "--name", clusterName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		r.Stderr.WriteString(fmt.Sprintf("Failed to update kubeconfig: %v\n%s", err, output))
		return ErrKubeConfig
	}
	r.Stdout.WriteString(fmt.Sprintf("Updated kubeconfig: %s\n", output))
	return nil
}

// generateSparkAppSpec generates the Spark application spec from the template
func generateSparkAppSpec(r *plugin.Runtime, templateFile string, s *sparkEksCommandContext, jobContext *sparkEksJobContext, clusterContext *sparkEksClusterContext, appName, queryURI, resultURI string) (string, error) {
	// Read the template file
	r.Stdout.WriteString(fmt.Sprintf("Loading template file: %s\n", templateFile))
	templateData, err := ioutil.ReadFile(templateFile)
	if err != nil {
		r.Stderr.WriteString(fmt.Sprintf("Failed to read template file: %v\n", err))
		return "", ErrTemplateFile
	}

	// Parse the YAML template
	var sparkApp SparkApplication
	if err := yaml.Unmarshal(templateData, &sparkApp); err != nil {
		r.Stderr.WriteString(fmt.Sprintf("Failed to parse template YAML: %v\n", err))
		return "", ErrApplicationSpec
	}

	// Update the application with our values
	sparkApp.Metadata.Name = appName
	sparkApp.Metadata.Namespace = s.KubeNamespace

	// Set main application file
	mainAppFile := queryURI
	if jobContext.ReturnResult {
		mainAppFile = wrapperScript
		// Add arguments for the wrapper script
		sparkApp.Spec.Arguments = []string{queryURI, resultURI}
	}
	sparkApp.Spec.MainApplicationFile = mainAppFile

	// Set image if provided in cluster context
	if clusterContext.Image != nil {
		sparkApp.Spec.Image = *clusterContext.Image
	}

	// Set region in environment variables
	if clusterContext.Region != nil {
		// Initialize EnvVars map if it's nil
		if sparkApp.Spec.Driver.EnvVars == nil {
			sparkApp.Spec.Driver.EnvVars = make(map[string]string)
		}
		sparkApp.Spec.Driver.EnvVars["AWS_REGION"] = *clusterContext.Region
	}

	// Update driver configuration from command properties
	if driverCores, ok := jobContext.Properties["spark.driver.cores"]; ok {
		// Parse string to int, handle error if needed
		var cores int
		if _, err := fmt.Sscanf(driverCores, "%d", &cores); err == nil {
			sparkApp.Spec.Driver.Cores = cores
		}
		delete(jobContext.Properties, "spark.driver.cores")
	}

	if driverMemory, ok := jobContext.Properties["spark.driver.memory"]; ok {
		sparkApp.Spec.Driver.Memory = driverMemory
		delete(jobContext.Properties, "spark.driver.memory")
	}

	// Update executor configuration from command properties
	if executorInstances, ok := jobContext.Properties["spark.executor.instances"]; ok {
		// Parse string to int, handle error if needed
		var instances int
		if _, err := fmt.Sscanf(executorInstances, "%d", &instances); err == nil {
			sparkApp.Spec.Executor.Instances = instances
		}
		delete(jobContext.Properties, "spark.executor.instances")
	}

	if executorCores, ok := jobContext.Properties["spark.executor.cores"]; ok {
		// Parse string to int, handle error if needed
		var cores int
		if _, err := fmt.Sscanf(executorCores, "%d", &cores); err == nil {
			sparkApp.Spec.Executor.Cores = cores
		}
		delete(jobContext.Properties, "spark.executor.cores")
	}

	if executorMemory, ok := jobContext.Properties["spark.executor.memory"]; ok {
		sparkApp.Spec.Executor.Memory = executorMemory
		delete(jobContext.Properties, "spark.executor.memory")
	}

	// Add cluster properties to sparkConf if any are defined
	if clusterContext.Properties != nil {
		for k, v := range clusterContext.Properties {
			sparkApp.Spec.SparkConf[k] = v
		}
	}

	// Initialize SparkConf map if it's nil
	if sparkApp.Spec.SparkConf == nil {
		sparkApp.Spec.SparkConf = make(map[string]string)
	}

	// Add any remaining job properties to sparkConf
	for k, v := range jobContext.Properties {
		sparkApp.Spec.SparkConf[k] = v
	}

	// Set logs URI if provided
	if s.LogsURI != nil {
		sparkApp.Spec.SparkConf["spark.eventLog.dir"] = *s.LogsURI
	}

	// Convert the updated spec back to YAML
	updatedYAML, err := yaml.Marshal(sparkApp)
	if err != nil {
		return "", fmt.Errorf("failed to convert application spec to YAML: %w", err)
	}

	return string(updatedYAML), nil
}

// uploadFileToS3 uploads content to S3
func uploadFileToS3(fileURI, content string, awsConfig aws.Config) error {
	// Get bucket name and prefix
	s3Parts := rxS3.FindAllStringSubmatch(fileURI, -1)
	if len(s3Parts) == 0 || len(s3Parts[0]) < 3 {
		return fmt.Errorf("unexpected queries key: %v", s3Parts)
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
*/
