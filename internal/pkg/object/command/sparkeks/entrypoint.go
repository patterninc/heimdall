package sparkeks

import (
	"path"
	"strings"

	"github.com/kubeflow/spark-operator/v2/api/v1beta2"
)

// --- JAR application type (JVM language) resolution ---
// Add an entry to support another JVM application type.
//e.g. Combination: applicationType (Java) -> EntrypointStrategy (JAR) + (Class name)

var jarApplicationTypes = map[string]v1beta2.SparkApplicationType{
	"scala": v1beta2.SparkApplicationTypeScala,
	"java":  v1beta2.SparkApplicationTypeJava,
}

const defaultJarApplicationType = v1beta2.SparkApplicationTypeScala

// resolveJarApplicationType maps an optional configured application type (e.g. "Java", "Scala",
// case-insensitive) to its SparkApplication type, falling back to the default (Scala) when
// unset or unrecognized.
func resolveJarApplicationType(applicationType string) v1beta2.SparkApplicationType {
	if t, ok := jarApplicationTypes[strings.ToLower(strings.TrimSpace(applicationType))]; ok {
		return t
	}
	return defaultJarApplicationType
}

type entrypointStrategy interface {
	apply(spec *v1beta2.SparkApplicationSpec)
}

func buildArguments(override []string, appName, queryURI, user, resultURI string, returnResult bool) []string {
	if len(override) > 0 {
		return override
	}
	args := []string{appName, queryURI, user}
	if returnResult {
		args = append(args, resultURI)
	}
	return args
}

type jarEntrypointStrategy struct {
	appType      v1beta2.SparkApplicationType
	mainClass    string
	appName      string
	queryURI     string
	user         string
	resultURI    string
	returnResult bool
	arguments    []string
}

func (s jarEntrypointStrategy) apply(spec *v1beta2.SparkApplicationSpec) {
	mainClass := s.mainClass
	spec.Type = s.appType
	spec.MainClass = &mainClass

	spec.Arguments = buildArguments(s.arguments, s.appName, s.queryURI, s.user, s.resultURI, s.returnResult)
}

type sqlWrapperEntrypointStrategy struct {
	appName      string
	queryURI     string
	user         string
	resultURI    string
	returnResult bool
	arguments    []string
}

func (s sqlWrapperEntrypointStrategy) apply(spec *v1beta2.SparkApplicationSpec) {
	spec.Arguments = buildArguments(s.arguments, s.appName, s.queryURI, s.user, s.resultURI, s.returnResult)
}

// entrypointFactory builds the entrypoint strategy for a job from its execution context.
type entrypointFactory func(execCtx *executionContext) entrypointStrategy

var entrypointStrategiesByExt = map[string]entrypointFactory{
	".jar": newJarEntrypointStrategy,
	".py":  newSQLWrapperEntrypointStrategy,
}

var defaultEntrypointFactory entrypointFactory = newSQLWrapperEntrypointStrategy
 
func newEntrypointStrategy(execCtx *executionContext) entrypointStrategy {
	ext := strings.ToLower(path.Ext(execCtx.commandContext.WrapperURI))
	factory, ok := entrypointStrategiesByExt[ext]
	if !ok {
		factory = defaultEntrypointFactory
	}
	return factory(execCtx)
}

func newJarEntrypointStrategy(execCtx *executionContext) entrypointStrategy {
	jobContext := execCtx.jobContext

	var mainClass, applicationType string
	if jobContext.Parameters != nil {
		mainClass = jobContext.Parameters.EntryPoint
		applicationType = jobContext.Parameters.ApplicationType
	}

	return jarEntrypointStrategy{
		appType:      resolveJarApplicationType(applicationType),
		mainClass:    mainClass,
		appName:      execCtx.appName,
		queryURI:     execCtx.s3aQueryURI,
		user:         execCtx.job.User,
		resultURI:    execCtx.s3aResultURI,
		returnResult: jobContext.ReturnResult,
		arguments:    jobContext.Arguments,
	}
}

func newSQLWrapperEntrypointStrategy(execCtx *executionContext) entrypointStrategy {
	jobContext := execCtx.jobContext

	return sqlWrapperEntrypointStrategy{
		appName:      execCtx.appName,
		queryURI:     execCtx.s3aQueryURI,
		user:         execCtx.job.User,
		resultURI:    execCtx.s3aResultURI,
		returnResult: jobContext.ReturnResult,
		arguments:    jobContext.Arguments,
	}
}
