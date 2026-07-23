package sparkeks

import (
	"path"
	"strings"

	"github.com/kubeflow/spark-operator/v2/api/v1beta2"
)

// --- JAR application type (JVM language) resolution ---

// jarApplicationTypes maps a configured application-type name (lower-cased) to its
// SparkApplication type. A JAR entrypoint is a JVM job, so only Java and Scala apply
// (Python/R are not JAR-based). Add an entry to support another JVM application type.
var jarApplicationTypes = map[string]v1beta2.SparkApplicationType{
	"scala": v1beta2.SparkApplicationTypeScala,
	"java":  v1beta2.SparkApplicationTypeJava,
}

// defaultJarApplicationType is used when a JAR job does not set Parameters.ApplicationType, or
// sets one that isn't recognized. Chipmunk's Writer is Scala.
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

// --- Entrypoint strategies ---

// entrypointStrategy configures the entrypoint-specific fields of a SparkApplication spec
// (Type, MainClass, Arguments) — the part that differs between a JAR/--class job and the
// default SQL-wrapper job. MainApplicationFile is common to both and set by the caller.
type entrypointStrategy interface {
	apply(spec *v1beta2.SparkApplicationSpec)
}

// jarEntrypointStrategy runs a custom-class JVM (Java/Scala) JAR. Its main() receives
// [appName, user, query, ...arguments, (resultURI)] — the contract chipmunk's
// com.pattern.chipmunk.Writer expects: args(0)=appName, args(1)=user, args(2)=<literal SQL
// text, NOT a URI>, args(3)=s3 target, optional args(4)=resultURI. appType is the resolved
// SparkApplication type (Scala/Java/...).
type jarEntrypointStrategy struct {
	appType      v1beta2.SparkApplicationType
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
	spec.Type = s.appType
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

// --- Extension-based strategy selection ---

// entrypointFactory builds the entrypoint strategy for a job from its execution context.
type entrypointFactory func(execCtx *executionContext) entrypointStrategy

// entrypointStrategiesByExt maps a wrapper-file extension (from commandContext.WrapperURI) to
// the entrypoint strategy that runs it. Add an entry to support a new wrapper type.
var entrypointStrategiesByExt = map[string]entrypointFactory{
	".jar": newJarEntrypointStrategy,
	".py":  newSQLWrapperEntrypointStrategy,
}

// defaultEntrypointFactory runs when the wrapper extension is unrecognized: the existing
// SQL-wrapper behavior (a safe default — it never sets a JAR MainClass/Type).
var defaultEntrypointFactory entrypointFactory = newSQLWrapperEntrypointStrategy

// newEntrypointStrategy selects the entrypoint strategy from the wrapper file's extension
// (e.g. ".jar" -> JAR/--class, ".py" -> SQL wrapper), falling back to the SQL wrapper for any
// unrecognized extension.
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
		user:         execCtx.job.User,
		query:        jobContext.Query,
		arguments:    jobContext.Arguments,
		resultURI:    updateS3ToS3aURI(execCtx.resultURI),
		returnResult: jobContext.ReturnResult,
	}
}

func newSQLWrapperEntrypointStrategy(execCtx *executionContext) entrypointStrategy {
	jobContext := execCtx.jobContext

	return sqlWrapperEntrypointStrategy{
		appName:      execCtx.appName,
		queryURI:     updateS3ToS3aURI(execCtx.queryURI),
		user:         execCtx.job.User,
		resultURI:    updateS3ToS3aURI(execCtx.resultURI),
		returnResult: jobContext.ReturnResult,
	}
}
