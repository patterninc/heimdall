package trino

import "github.com/patterninc/heimdall/pkg/object/job"

type jobContext struct {
	Query string `yaml:"query"`
}

func newJobContext(j *job.Job) (*jobContext, error) {
	if j.Context == nil {
		return nil, ErrEmptyJobContext
	}

	ctx := &jobContext{}
	if err := j.Context.Unmarshal(ctx); err != nil {
		return nil, err
	}
	return ctx, nil
}
