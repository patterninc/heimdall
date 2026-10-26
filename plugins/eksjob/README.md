# EKS Job Plugin

The `eksjob` plugin runs a Heimdall job as a **Kubernetes custom resource** on an AWS EKS cluster. It does not know about any particular CRD: the command context names the resource and how to fill it in, and any operator whose resource exposes the [status contract](#status-contract) below can be driven by it.

For each job the plugin:

1. Assumes the cluster's IAM role, calls `eks:DescribeCluster` for the endpoint and CA, and authenticates to the Kubernetes API with a presigned STS token (the token `aws eks get-token` produces). No kubeconfig or `aws` CLI is needed.
2. Builds the resource: `spec_defaults`, then the job's `spec` (top-level keys replace), then the job `size`'s resources at `spec.<size_field>`.
3. Creates it as `heimdall-<job id>`. If a resource with that name already exists (a retry of the same Heimdall job), the plugin attaches to it instead of creating a second run.
4. Polls `status.phase` every `poll_interval` until it is `Succeeded` or `Failed`, while streaming the logs of the run's pods to the job's stdout.
5. Succeeds on `Succeeded`. On `Failed` it prints the run's recent Warning events to stderr and fails the job with `status.failureReason` in the error.

Killing the Heimdall job cancels the run by merge-patching `spec.<cancel_field>: <cancel_value>`. The plugin never deletes the resource; the operator owns its lifecycle (for example deleting finished runs after a TTL).

## Configuration

### Command context

| Field | Required | Default | Description |
|---|---|---|---|
| `crd_group` | yes | | API group of the custom resource, e.g. `example.com` |
| `crd_version` | yes | | API version, e.g. `v1alpha1` |
| `crd_kind` | yes | | Kind, e.g. `ExampleRun` |
| `crd_plural` | yes | | Resource (plural) name, e.g. `exampleruns` |
| `poll_interval` | no | `15s` | How often `status.phase` is read (Go duration, minimum `1s`) |
| `spec_defaults` | no | | Object merged first into the resource `spec` |
| `required_spec_fields` | no | | Top-level spec keys that must be set after the merge; the job fails before anything is created otherwise |
| `size_field` | no | `resources` | Spec key that receives the size profile |
| `size_profiles` | no | | Map of size name to a `corev1.ResourceRequirements` (`requests`, `limits`). When set, every job must name one of them |
| `cancel_field` | no | | Spec key set to cancel a run. When empty, killing the job does not cancel the run |
| `cancel_value` | no | `cancel` | Value written to `cancel_field` |
| `pod_selector_label` | no | | Pods labelled `<pod_selector_label>=<resource name>` are the run's pods; their logs are streamed. When empty, no logs are streamed |

`New` only checks the shape of this context (types, the four `crd_*` fields, every size profile is an object) and makes no network calls, because a `New` error stops Heimdall from starting. Values such as quantities in a size profile or the image are checked when a job uses them, and a bad value fails that job with a clear error.

```yaml
commands:
  - name: example-run-0.1.0
    status: active
    plugin: eksjob
    tags: [type:example-run]
    cluster_tags: [type:example-eks]
    context:
      crd_group: example.com
      crd_version: v1alpha1
      crd_kind: ExampleRun
      crd_plural: exampleruns
      poll_interval: 15s
      spec_defaults:
        image: registry.example.com/runner:1.0.0
        serviceAccountName: runner
        ttl: 1h
        nodeSelector:
          workload_type: runs
        tolerations:
          - key: dedicated
            operator: Equal
            value: runs
            effect: NoSchedule
      size_field: resources
      size_profiles:
        s: {requests: {cpu: "1", memory: 2Gi}, limits: {memory: 2Gi}}
        m: {requests: {cpu: "2", memory: 4Gi}, limits: {memory: 4Gi}}
      cancel_field: setState
      cancel_value: cancel
      pod_selector_label: example.com/run
```

### Cluster context

| Field | Required | Default | Description |
|---|---|---|---|
| `cluster_name` | yes | | EKS cluster name |
| `region` | no | `AWS_REGION` / default AWS config | AWS region of the cluster |
| `role_arn` | no | | IAM role assumed for `eks:DescribeCluster` and the Kubernetes token. Without it the Heimdall task's own credentials are used |
| `namespace` | no | `default` | Namespace the resources are created in |

```yaml
clusters:
  - name: example-eks
    status: active
    health_check: true
    tags: [type:example-eks]
    context:
      role_arn: arn:aws:iam::123456789012:role/example
      region: us-west-2
      cluster_name: example
      namespace: runs
```

The role must be mapped to a Kubernetes identity (EKS access entry or `aws-auth`) that can, in the namespace:

| Resource | Verbs | Used for |
|---|---|---|
| `<crd_plural>` | create, get, list, patch | Creating the run, polling its status, the health check, cancelling |
| `pods` | list | Finding the run's pods |
| `pods/log` | get | Streaming their logs |
| `events` | list | Reporting Warning events when a run fails |

plus IAM `eks:DescribeCluster` on the cluster. The plugin never deletes, updates or watches anything, and reads status through `get` on the resource itself, so it needs no `<crd_plural>/status` access.

### Job context

| Field | Required | Description |
|---|---|---|
| `size` | when the command has `size_profiles` | One of the command's size profiles |
| `labels` | no | Copied to the resource's `metadata.labels` after sanitising to label syntax (invalid characters become `_`, values are cut to 63 characters, keys that cannot be made valid are dropped with a warning) |
| `spec` | no | Merged over `spec_defaults`; top-level keys replace |

Unknown fields in any context are ignored.

```json
{
  "name": "example-run",
  "version": "0.0.1",
  "command_criteria": ["type:example-run"],
  "cluster_criteria": ["type:example-eks"],
  "context": {
    "size": "s",
    "labels": {"owner": "jane.doe@example.com", "dag_id": "ingest_orders", "task_id": "load"},
    "spec": {
      "pipeline": "tasks:\n  - name: echo\n",
      "taskCount": 1,
      "timeout": "7200s",
      "retries": 0
    }
  }
}
```

The size profile is authoritative: `spec.<size_field>` from the job is replaced by the profile. Besides the job labels, the resource gets `app.kubernetes.io/managed-by: heimdall` and `heimdall/job-id: <job id>`, which job labels cannot override.

The resource stays well under the Kubernetes 1 MiB object limit only if inline spec values (such as a pipeline document) do; large payloads belong in object storage, referenced from the spec.

## Status contract

Any custom resource run by `eksjob` must expose:

| Field | Values | Used for |
|---|---|---|
| `status.phase` | `Pending`, `Running`, `Succeeded`, `Failed` | `Succeeded` and `Failed` are terminal; anything else (including empty) keeps polling |
| `status.failureReason` | string | Included in the Heimdall job error when the phase is `Failed` |
| `status.startTime` | RFC 3339 timestamp | Printed with the run duration when the run finishes |
| `status.completionTime` | RFC 3339 timestamp | Printed with the run duration when the run finishes |

and, to support cancellation, react to `spec.<cancel_field>: <cancel_value>` by stopping the run and moving it to a terminal phase.

## Logs

The plugin lists the run's pods (`<pod_selector_label>=<resource name>`) every few seconds and follows the log of every started container (`follow`, with timestamps). Pods that appear late are picked up by a later listing, and a restarted container is streamed again as a new instance. If a stream drops while its container is still running it is resumed from the last timestamp written, without repeating lines. Once more than one pod has been seen, lines are prefixed with `[<pod>]` (`[<pod>/<container>]` for multi-container pods). Containers stuck in a waiting state such as `ImagePullBackOff` are reported once on stderr.

Logs are streamed to stdout; the plugin's own progress lines go to stdout and warnings to stderr. When the run finishes the plugin waits up to 30 seconds for open streams to reach the end of the log.

Logs of a container instance that started and restarted between two listings are not backfilled. Only `spec.containers` are streamed: init containers (including native sidecars) are neither streamed nor reported when waiting.

## Outputs

The plugin publishes `run_name` and `run_namespace` as job outputs, so cluster `attributes` can link to the run, e.g. `{{ .Outputs.run_name }}`.

## Health check

With `health_check: true` on the cluster, the probe assumes the role, checks the EKS cluster is `ACTIVE` and lists `<crd_plural>` in the namespace.
