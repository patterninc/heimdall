package eksjob

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
	"unicode"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation"
)

const (
	runNamePrefix  = `heimdall-`
	maxLabelLength = validation.LabelValueMaxLength

	labelManagedBy      = `app.kubernetes.io/managed-by`
	labelManagedByValue = `heimdall`
	labelJobID          = `heimdall/job-id`
)

var (
	rxRunNameInvalid     = regexp.MustCompile(`[^a-z0-9-]+`)
	rxLabelInvalid       = regexp.MustCompile(`[^A-Za-z0-9._-]`)
	rxLabelPrefixInvalid = regexp.MustCompile(`[^a-z0-9.-]`)
)

// runName derives the resource name from the Heimdall job ID, so a retry of the same
// Heimdall job finds the existing run instead of starting a second one.
func runName(jobID string) (string, error) {

	id := strings.Trim(rxRunNameInvalid.ReplaceAllString(strings.ToLower(jobID), `-`), `-`)
	if id == `` {
		return ``, fmt.Errorf("cannot derive a resource name from job ID %q", jobID)
	}

	name := runNamePrefix + id
	if len(name) > validation.DNS1123LabelMaxLength {
		name = strings.TrimRight(name[:validation.DNS1123LabelMaxLength], `-`)
	}

	if errs := validation.IsDNS1123Label(name); len(errs) > 0 {
		return ``, fmt.Errorf("invalid resource name %q derived from job ID %q: %s", name, jobID, strings.Join(errs, `; `))
	}

	return name, nil

}

// buildRun renders the custom resource for a job: spec_defaults, then the job spec (top-level
// keys replace), then the size profile's resources at spec.<size_field>.
func (c *commandContext) buildRun(name, namespace, jobID string, jc *jobContext, warnings io.Writer) (*unstructured.Unstructured, error) {

	spec := map[string]any{}
	if c.SpecDefaults != nil {
		spec = runtime.DeepCopyJSON(c.SpecDefaults)
	}
	for k, v := range jc.Spec {
		spec[k] = runtime.DeepCopyJSONValue(v)
	}

	resources, err := c.sizeResources(jc.Size)
	if err != nil {
		return nil, err
	}
	if resources != nil {
		spec[c.SizeField] = resources
	}

	var missing []string
	for _, field := range c.RequiredSpecFields {
		if v, found := spec[field]; !found || v == nil || v == `` {
			missing = append(missing, field)
		}
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("job context spec is missing required field(s): %s", strings.Join(missing, `, `))
	}

	labels, dropped := sanitizeLabels(jc.Labels)
	for _, key := range dropped {
		fmt.Fprintf(warnings, "eksjob: dropping label %q: cannot be made a valid Kubernetes label\n", key)
	}
	labels[labelManagedBy] = labelManagedByValue
	if v := sanitizeLabelValue(jobID); v != `` {
		labels[labelJobID] = v
	}

	obj := &unstructured.Unstructured{Object: map[string]any{`spec`: spec}}
	obj.SetAPIVersion(c.CRDGroup + `/` + c.CRDVersion)
	obj.SetKind(c.CRDKind)
	obj.SetName(name)
	obj.SetNamespace(namespace)
	obj.SetLabels(labels)

	return obj, nil

}

// sizeResources returns the size profile as an unstructured corev1.ResourceRequirements, or
// nil when the command defines no size profiles and the job names no size.
func (c *commandContext) sizeResources(size string) (map[string]any, error) {

	if size == `` {
		if len(c.SizeProfiles) == 0 {
			return nil, nil
		}
		return nil, fmt.Errorf("job context size is required; valid sizes: %s", strings.Join(c.sizeNames(), `, `))
	}

	raw, found := c.SizeProfiles[size]
	if !found {
		return nil, fmt.Errorf("unknown job size %q; valid sizes: %s", size, strings.Join(c.sizeNames(), `, `))
	}

	var rr corev1.ResourceRequirements
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&rr); err != nil {
		return nil, fmt.Errorf("invalid size profile %q in command context: %w", size, err)
	}

	resources, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&rr)
	if err != nil {
		return nil, fmt.Errorf("invalid size profile %q in command context: %w", size, err)
	}

	return resources, nil

}

func (c *commandContext) sizeNames() []string {
	names := make([]string, 0, len(c.SizeProfiles))
	for name := range c.SizeProfiles {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// sanitizeLabels turns arbitrary job labels into valid Kubernetes labels. Keys that cannot be
// made valid are returned in dropped.
func sanitizeLabels(in map[string]any) (labels map[string]string, dropped []string) {

	labels = make(map[string]string, len(in))
	for k, v := range in {
		key := sanitizeLabelKey(k)
		if key == `` {
			dropped = append(dropped, k)
			continue
		}
		value := ``
		if v != nil {
			value = sanitizeLabelValue(fmt.Sprint(v))
		}
		labels[key] = value
	}
	sort.Strings(dropped)

	return labels, dropped

}

// sanitizeLabelKey returns a valid label key ([prefix/]name), or "" when k cannot be made valid.
func sanitizeLabelKey(k string) string {

	prefix, name := ``, k
	if i := strings.LastIndex(k, `/`); i >= 0 {
		prefix, name = k[:i], k[i+1:]
	}

	name = sanitizeLabelValue(name)
	if name == `` {
		return ``
	}

	key := name
	if prefix != `` {
		prefix = strings.Trim(rxLabelPrefixInvalid.ReplaceAllString(strings.ToLower(prefix), `-`), `-.`)
		if len(validation.IsDNS1123Subdomain(prefix)) > 0 {
			return ``
		}
		key = prefix + `/` + name
	}

	if len(validation.IsQualifiedName(key)) > 0 {
		return ``
	}

	return key

}

// sanitizeLabelValue replaces characters not allowed in label values with "_", truncates to 63
// characters and trims non-alphanumeric characters from both ends.
func sanitizeLabelValue(v string) string {
	v = rxLabelInvalid.ReplaceAllString(v, `_`)
	if len(v) > maxLabelLength {
		v = v[:maxLabelLength]
	}
	return strings.TrimFunc(v, func(r rune) bool {
		return !(r < unicode.MaxASCII && (unicode.IsLetter(r) || unicode.IsDigit(r)))
	})
}
