package heimdall

import (
	"context"
	"crypto/rand"
	_ "embed"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/babourine/x/pkg/set"
	"github.com/gorilla/mux"
	"github.com/hladush/go-telemetry/pkg/telemetry"

	"github.com/patterninc/heimdall/internal/pkg/aws"
	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/command"
	"github.com/patterninc/heimdall/pkg/object/job"
	jobStatus "github.com/patterninc/heimdall/pkg/object/job/status"
	"github.com/patterninc/heimdall/pkg/object/status"
	"github.com/patterninc/heimdall/pkg/plugin"
)

const (
	defaultPairsLength = 3
	formatFileNotFound = `unknown file: %s`
	resultFile         = `result`
	resultFilename     = resultFile + `.json`
	formatUserAgent    = `heimdall/%s`
	separator          = `/`
	jobFileFormat      = `%s/%s/%s`
	s3Prefix           = `s3://`
)

var (
	ErrCommandClusterPairNotFound = fmt.Errorf(`command-cluster pair is not found`)
	ErrJobCancelFailed            = fmt.Errorf(`async job unrecognized or already in final state`)
	runJobMethod                  = telemetry.NewMethod("runJob", "heimdall")
)

//go:embed queries/job/status_cancel_update.sql
var queryJobCancelUpdate string

type commandOnCluster struct {
	command *command.Command
	cluster *cluster.Cluster
}

func (h *Heimdall) submitJob(ctx context.Context, j *job.Job) (any, error) {

	// set / add job properties
	if err := j.Init(); err != nil {
		return nil, err
	}

	// let's determine command that we'll be running and on what compute...
	command, cluster, err := h.resolveJob(j.CommandCriteria, j.ClusterCriteria)
	if err != nil {
		return j, err
	}

	// let's set the mode in which we'll execute our job
	// ...and if our job is an async job, we just log it and return
	if j.IsSync = command.IsSync; !j.IsSync {
		if _, err := h.insertJob(j, cluster.ID, command.ID); err != nil {
			return nil, err
		}
		return j, nil
	}

	// let's run the job
	err = h.runJob(ctx, j, command, cluster)

	// before we process the error, we'll make the best effort to record this job in the database
	go h.insertJob(j, cluster.ID, command.ID)

	return j, err

}

func (h *Heimdall) runJob(ctx context.Context, j *job.Job, command *command.Command, cluster *cluster.Cluster) error {

	defer runJobMethod.RecordLatency(time.Now(), command.Name, cluster.Name)
	runJobMethod.CountRequest(command.Name, cluster.Name)

	// let's set environment
	runtime := &plugin.Runtime{
		WorkingDirectory: h.JobsDirectory + separator + j.ID,
		ArchiveDirectory: h.ArchiveDirectory + separator + j.ID,
		ResultDirectory:  h.ResultDirectory + separator + j.ID,
		Version:          h.Version,
		UserAgent:        fmt.Sprintf(formatUserAgent, h.Version),
	}

	// we're done with funtime...
	defer runtime.Close()

	if err := runtime.Set(); err != nil {
		return err
	}

	// set keepalive logic that will maintain the last timestamp the runJob function was active
	// we set a channel that we close when runJob completes
	// closing the channel will notify keepalive function that updates the timestamp in the db to complete
	keepaliveActive := make(chan struct{})
	defer close(keepaliveActive)

	// ...and now we just start keepalive function for this job
	go h.jobKeepalive(keepaliveActive, j.SystemID, h.agentName)

	// Create channels for coordination between plugin execution and cancellation monitoring
	jobDone := make(chan error, 1)
	cancelMonitorDone := make(chan struct{})

	// Create cancellable context for the job
	pluginCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Start plugin execution in goroutine
	go func() {
		defer close(cancelMonitorDone) // signal monitoring to stop
		err := h.commandHandlers[command.ID](pluginCtx, runtime, j, cluster)
		jobDone <- err
	}()

	// Start cancellation monitoring for async jobs
	if !j.IsSync {
		go func() {
			ticker := time.NewTicker(10 * time.Second)
			defer ticker.Stop()

			for {
				select {
				// plugin finished, stop monitoring
				case <-cancelMonitorDone:
					return
				case <-ticker.C:
					// If job is in cancelling state, trigger context cancellation
					result, err := h.getJobStatus(ctx, &jobRequest{ID: j.ID})
					if err == nil {
						if job, ok := result.(*job.Job); ok && job.Status == jobStatus.Cancelling {
							cancel()
							return
						}
					}
				}
			}
		}()
	}

	// Wait for job execution to complete
	jobErr := <-jobDone

	// Check if context was cancelled and mark status appropriately
	if pluginCtx.Err() != nil {
		j.Status = jobStatus.Cancelling // janitor will update to cancelled when resources are cleaned up
		runJobMethod.LogAndCountError(pluginCtx.Err(), command.Name, cluster.Name)
		return nil
	}

	// Handle plugin execution result (only if not cancelled)
	if jobErr != nil {
		j.Status = jobStatus.Failed
		j.Error = jobErr.Error()
		runJobMethod.LogAndCountError(jobErr, command.Name, cluster.Name)
		return jobErr
	}

	if j.StoreResultSync || !j.IsSync {
		h.storeResults(runtime, j)
	} else {
		go h.storeResults(runtime, j)
	}

	j.Status = jobStatus.Succeeded

	runJobMethod.CountSuccess(command.Name, cluster.Name)
	return nil

}

func (h *Heimdall) storeResults(runtime *plugin.Runtime, j *job.Job) error {
	// do we have result to be written?
	if j.Result == nil {
		return nil
	}

	// prepare result
	data, err := json.Marshal(j.Result)
	if err != nil {

		return err
	}

	// write result
	writeFileFunc := os.WriteFile
	if strings.HasPrefix(runtime.ResultDirectory, s3Prefix) {
		writeFileFunc = func(name string, data []byte, perm os.FileMode) error {
			return aws.WriteToS3(context.Background(), name, data, perm)
		}
	}
	if err := writeFileFunc(runtime.ResultDirectory+separator+resultFilename, data, 0600); err != nil {

		return err
	}

	return nil
}

func (h *Heimdall) cancelJob(ctx context.Context, req *jobRequest) (any, error) {

	sess, err := h.Database.NewSession(false)
	if err != nil {
		return nil, err
	}
	defer sess.Close()

	// Attempt to cancel
	rowsAffected, err := sess.Exec(queryJobCancelUpdate, req.ID, req.User)
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, ErrJobCancelFailed
	}

	// return job status
	return &job.Job{
		Status: jobStatus.Cancelling,
	}, nil
}

func (h *Heimdall) getJobFile(w http.ResponseWriter, r *http.Request) {

	// get vars
	vars := mux.Vars(r)
	jobID := vars[`id`]
	filename := vars[`file`]

	// do we have requested file in the allow list?
	allowFiles := set.New([]string{`stdout`, `stderr`, resultFile})

	if !allowFiles.Has(filename) {
		writeAPIError(w, fmt.Errorf(formatFileNotFound, filename), nil)
		return
	}

	// let's validate jobID we got
	if _, err := h.getJobStatus(r.Context(), &jobRequest{ID: jobID}); err != nil {
		writeAPIError(w, err, nil)
		return
	}

	// set context of the requested file
	contentType := contentTypePlain
	sourceDirectory := h.ArchiveDirectory

	if filename == resultFile {
		filename = resultFilename
		contentType = contentTypeJSON
		sourceDirectory = h.ResultDirectory
	}

	// get the file content
	readFileFunc := os.ReadFile
	filenamePath := fmt.Sprintf(jobFileFormat, sourceDirectory, jobID, filename)
	if strings.HasPrefix(filenamePath, s3Prefix) {
		readFileFunc = func(path string) ([]byte, error) {
			return aws.ReadFromS3(r.Context(), path)
		}
	}

	// get file's content
	data, err := readFileFunc(filenamePath)
	if err != nil {
		writeAPIError(w, err, nil)
		return
	}

	// return to the user...
	w.Header().Add(contentTypeKey, contentType)
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

// resolve the job criteria into command-cluster pair
func (h *Heimdall) resolveJob(commandCriteria, clusterCriteria *set.Set[string]) (*command.Command, *cluster.Cluster, error) {

	// let's set the struct and a slice in which we'll collect all the matching pairs
	pairs := make([]*commandOnCluster, 0, defaultPairsLength)

	// to resolve the job request we need to find a command-cluster pair that satisfies the job criteria
	// we start with findding commands list (multiple commands may satisfy the criteria)
	for _, command := range h.Commands {
		// we only use active commands...
		if command.Status != status.Active {
			continue
		}
		// ...where tags match...
		if commandCriteria != nil && command.Tags.Contains(commandCriteria) {
			// ...we found command-candidate, let's find all corresponding clusters
			for _, cluster := range h.Clusters {
				// ...cluster must be active...
				if cluster.Status != status.Active {
					continue
				}
				// ...and matchcluster tags from the selected command...
				if cluster.Tags.Contains(command.ClusterTags) {
					// ...and also match cluster criteria for the job!
					if clusterCriteria != nil && cluster.Tags.Contains(clusterCriteria) {
						// we found matching pair! let's add it to the list of potential candidates
						pairs = append(pairs, &commandOnCluster{command, cluster})
					}
				}
			}
		}
	}

	// let's select command that we execute...
	// with te cluster we execute it on...
	pairIndex := 0
	if l := int64(len(pairs)); l == 0 {
		// we're here because we did not find a command-cluster pair to run our job
		return nil, nil, ErrCommandClusterPairNotFound
	} else if l > 1 {
		// TODO: we need to support a custom selector that user can supply to choose
		// compute when there are multiple matches
		// (for example case when we want ot send 1% of traffic to the next version of cluster)
		// random selection for now should provide close to even distribution
		n, err := rand.Int(rand.Reader, big.NewInt(l))
		if err != nil {
			return nil, nil, err
		}
		pairIndex = int(n.Int64())
	}

	// if there was only one pair found, pairIndex will stay zero...
	return pairs[pairIndex].command, pairs[pairIndex].cluster, nil

}
