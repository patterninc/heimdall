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
	apply(spec *v1beta2.SparkApplicationSpec) error
}

func buildArguments(extra []string, appName, queryURI, user, resultURI string, returnResult bool) []string {
	if !returnResult {
		resultURI = ``
	}
	return append([]string{appName, queryURI, user, resultURI}, extra...)
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

func (s jarEntrypointStrategy) apply(spec *v1beta2.SparkApplicationSpec) error {
	mainClass := strings.TrimSpace(s.mainClass)
	if mainClass == "" {
		return ErrMissingEntryPoint
	}
	spec.Type = s.appType
	spec.MainClass = &mainClass

	spec.Arguments = buildArguments(s.arguments, s.appName, s.queryURI, s.user, s.resultURI, s.returnResult)
	return nil
}

type sqlWrapperEntrypointStrategy struct {
	appName      string
	queryURI     string
	user         string
	resultURI    string
	returnResult bool
	arguments    []string
}

func (s sqlWrapperEntrypointStrategy) apply(spec *v1beta2.SparkApplicationSpec) error {
	spec.Arguments = buildArguments(s.arguments, s.appName, s.queryURI, s.user, s.resultURI, s.returnResult)
	return nil
}

type pysparkEntrypointStrategy struct {
	appName       string
	queryURI      string
	user          string
	resultURI     string
	returnResult  bool
	arguments     []string
	scriptURI     string // pre-uploaded .py; job passes parameters.script_uri
	bundleURI     string // command bundle_uri; job never sets this
	bundleVersion string // job passes parameters.bundle_version
	entryPoint    string // job passes parameters.entry_point (path inside the zip)
}

func (s pysparkEntrypointStrategy) apply(spec *v1beta2.SparkApplicationSpec) error {
	if s.scriptURI != "" {
		extra := append([]string{s.scriptURI, ""}, s.arguments...)
		spec.Arguments = buildArguments(extra, s.appName, s.queryURI, s.user, s.resultURI, s.returnResult)
		return nil
	}

	if strings.TrimSpace(s.bundleVersion) == "" {
		return ErrMissingBundleVersion
	}
	entryPoint := strings.TrimSpace(s.entryPoint)
	if entryPoint == "" {
		return ErrMissingBundleEntry
	}

	bundleZipURI := updateS3ToS3aURI(strings.TrimRight(s.bundleURI, "/") + "/" + strings.TrimSpace(s.bundleVersion) + ".zip")
	extra := append([]string{bundleZipURI, entryPoint}, s.arguments...)
	spec.Arguments = buildArguments(extra, s.appName, s.queryURI, s.user, s.resultURI, s.returnResult)
	return nil
}

func validateScriptURI(scriptURI, allowedPrefix string) error {
	scriptURI = strings.TrimSpace(scriptURI)
	if scriptURI == "" {
		return ErrInvalidScriptURI
	}
	if !strings.HasPrefix(scriptURI, s3Prefix) && !strings.HasPrefix(scriptURI, s3aPrefix) {
		return ErrInvalidScriptURI
	}
	if strings.Contains(scriptURI, "..") {
		return ErrInvalidScriptURI
	}
	if !strings.HasSuffix(strings.ToLower(scriptURI), ".py") {
		return ErrInvalidScriptURI
	}

	normalized := updateS3ToS3aURI(scriptURI)
	prefix := updateS3ToS3aURI(strings.TrimRight(strings.TrimSpace(allowedPrefix), "/")) + "/"
	if !strings.HasPrefix(normalized, prefix) {
		return ErrInvalidScriptURI
	}
	return nil
}

// entrypointFactory builds the entrypoint strategy for a job from its execution context.
type entrypointFactory func(execCtx *executionContext) entrypointStrategy

var entrypointStrategiesByExt = map[string]entrypointFactory{
	".jar": newJarEntrypointStrategy,
	".py":  newSQLWrapperEntrypointStrategy,
}

var defaultEntrypointFactory entrypointFactory = newSQLWrapperEntrypointStrategy

func newEntrypointStrategy(execCtx *executionContext) entrypointStrategy {
	if execCtx.commandContext.BundleURI != "" {
		return newPySparkEntrypointStrategy(execCtx)
	}

	ext := strings.ToLower(path.Ext(execCtx.commandContext.WrapperURI))
	factory, ok := entrypointStrategiesByExt[ext]
	if !ok {
		factory = defaultEntrypointFactory
	}
	return factory(execCtx)
}

func newPySparkEntrypointStrategy(execCtx *executionContext) entrypointStrategy {
	jobContext := execCtx.jobContext

	s := pysparkEntrypointStrategy{
		appName:      execCtx.appName,
		queryURI:     execCtx.s3aQueryURI,
		user:         execCtx.job.User,
		resultURI:    execCtx.s3aResultURI,
		returnResult: jobContext.ReturnResult,
		arguments:    jobContext.Arguments,
		bundleURI: execCtx.commandContext.BundleURI,
	}

	if jobContext.Parameters != nil && strings.TrimSpace(jobContext.Parameters.ScriptURI) != "" {
		s.scriptURI = updateS3ToS3aURI(strings.TrimSpace(jobContext.Parameters.ScriptURI))
		return s
	}

	if jobContext.Parameters != nil {
		s.entryPoint = jobContext.Parameters.EntryPoint
		s.bundleVersion = jobContext.Parameters.BundleVersion
	}
	return s
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
