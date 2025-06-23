package dynamo

import (
	ct "context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	"github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/job"
	"github.com/patterninc/heimdall/pkg/plugin"
	"github.com/patterninc/heimdall/pkg/result"
)

// dynamoJobContext represents the context for a dynamo job
type dynamoJobContext struct {
	Query string `yaml:"query,omitempty" json:"query,omitempty"`
	Limit int    `yaml:"limit,omitempty" json:"limit,omitempty"`
}

// dynamoClusterContext represents the context for a dynamo endpoint
type dynamoClusterContext struct {
	RoleARN *string `yaml:"role_arn,omitempty" json:"role_arn,omitempty"`
}

// dynamoCommandContext represents the dynamo command context
type dynamoCommandContext struct{}

var (
	ctx               = ct.Background()
	assumeRoleSession = aws.String("AssumeRoleSession")
)

// New creates a new dynamo plugin handler.
func New(_ *context.Context) (plugin.Handler, error) {

	s := &dynamoCommandContext{}
	return s.handler, nil

}

// Handler for the Spark job submission.
func (d *dynamoCommandContext) handler(r *plugin.Runtime, j *job.Job, c *cluster.Cluster) (err error) {

	// let's unmarshal job context
	jobContext := &dynamoJobContext{}
	if j.Context != nil {
		if err := j.Context.Unmarshal(jobContext); err != nil {
			return err
		}
	}

	// let's unmarshal cluster context
	clusterContext := &dynamoClusterContext{}
	if c.Context != nil {
		if err := c.Context.Unmarshal(clusterContext); err != nil {
			return err
		}
	}

	// setting AWS client
	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	// let's set empty options function,...
	assumeRoleOptions := func(_ *dynamodb.Options) {}

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

		assumeRoleOptions = func(o *dynamodb.Options) {
			o.Credentials = credentials.NewStaticCredentialsProvider(
				*assumeRoleOutput.Credentials.AccessKeyId,
				*assumeRoleOutput.Credentials.SecretAccessKey,
				*assumeRoleOutput.Credentials.SessionToken,
			)
		}

	}

	svc := dynamodb.NewFromConfig(awsConfig, assumeRoleOptions)

	executeStatementInput := &dynamodb.ExecuteStatementInput{
		Statement: aws.String(jobContext.Query),
	}

	// do we have limit in job context?
	if jobContext.Limit > 0 {
		executeStatementInput.Limit = aws.Int32(int32(jobContext.Limit))
	}

	// paginate
	for {

		executeStatementOutput, err := svc.ExecuteStatement(ctx, executeStatementInput)
		if err != nil {
			return err
		}

		// before we check for the result, let's assume we do not have the next page...
		executeStatementInput.NextToken = nil

		if executeStatementOutput != nil {
			page, err := result.FromDynamo(executeStatementOutput.Items)
			if err != nil {
				return err
			}
			executeStatementInput.NextToken = executeStatementOutput.NextToken
			if page != nil {
				if j.Result == nil {
					j.Result = page
				} else {
					j.Result.Data = append(j.Result.Data, page.Data...)
				}
			}
			if jobContext.Limit > 0 && len(j.Result.Data) >= jobContext.Limit {
				break
			}
		}

		if executeStatementInput.NextToken == nil {
			break
		}

	}

	return nil

}
