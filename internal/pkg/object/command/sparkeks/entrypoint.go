package sparkeks

import (
	"github.com/kubeflow/spark-operator/v2/api/v1beta2"
)

// jarApplicationType is the SparkApplication.Spec.Type used for JAR/--class entrypoint
// jobs (Parameters.EntryPoint set). Chipmunk's Writer is Scala; if a Java JAR entrypoint
// is ever needed this would need to become configurable.
const jarApplicationType = v1beta2.SparkApplicationTypeScala

// entrypointStrategy configures the entrypoint-specific fields of a SparkApplication spec
// (Type, MainClass, Arguments) — the part that differs between a JAR/--class job and the
// default SQL-wrapper job. MainApplicationFile is common to both and set by the caller.
type entrypointStrategy interface {
	apply(spec *v1beta2.SparkApplicationSpec)
}

// jarEntrypointStrategy runs a custom-class (Scala/Java) JAR. Its main() receives
// [appName, user, query, ...arguments, (resultURI)] — the contract chipmunk's
// com.pattern.chipmunk.Writer expects: args(0)=appName, args(1)=user, args(2)=<literal SQL
// text, NOT a URI>, args(3)=s3 target, optional args(4)=resultURI.
type jarEntrypointStrategy struct {
	mainClass    string
	appName      string
	user         string
	query        string
	arguments    []string
	resultURI    string
	returnResult bool
}

func (s jarEntrypointStrategy) apply(spec *v1beta2.SparkApplicationSpec) {
	mainClass := s.mainClass
	spec.Type = jarApplicationType
	spec.MainClass = &mainClass

	args := []string{s.appName, s.user, s.query}
	args = append(args, s.arguments...)
	if s.returnResult {
		args = append(args, s.resultURI)
	}
	spec.Arguments = args
}

// sqlWrapperEntrypointStrategy runs the default .py SQL wrapper, whose main() receives
// [appName, queryURI, user, (resultURI)] and reads the query from queryURI. This is the
// existing (pre-JAR-support) behavior, preserved byte-for-byte.
type sqlWrapperEntrypointStrategy struct {
	appName      string
	queryURI     string
	user         string
	resultURI    string
	returnResult bool
}

func (s sqlWrapperEntrypointStrategy) apply(spec *v1beta2.SparkApplicationSpec) {
	if s.returnResult {
		spec.Arguments = []string{s.appName, s.queryURI, s.user, s.resultURI}
	} else {
		spec.Arguments = []string{s.appName, s.queryURI, s.user}
	}
}

// newEntrypointStrategy selects the entrypoint strategy for a job: the JAR/--class strategy
// when a main class is configured (jobContext.Parameters.EntryPoint set), otherwise the
// default SQL-wrapper strategy.
func newEntrypointStrategy(execCtx *executionContext) entrypointStrategy {
	jobContext := execCtx.jobContext
	s3aResultURI := updateS3ToS3aURI(execCtx.resultURI)

	if jobContext.Parameters != nil && jobContext.Parameters.EntryPoint != "" {
		return jarEntrypointStrategy{
			mainClass:    jobContext.Parameters.EntryPoint,
			appName:      execCtx.appName,
			user:         execCtx.job.User,
			query:        jobContext.Query,
			arguments:    jobContext.Arguments,
			resultURI:    s3aResultURI,
			returnResult: jobContext.ReturnResult,
		}
	}

	return sqlWrapperEntrypointStrategy{
		appName:      execCtx.appName,
		queryURI:     updateS3ToS3aURI(execCtx.queryURI),
		user:         execCtx.job.User,
		resultURI:    s3aResultURI,
		returnResult: jobContext.ReturnResult,
	}
}
