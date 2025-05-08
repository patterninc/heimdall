package janitor

import (
	"fmt"
	"time"

	"github.com/patterninc/heimdall/internal/pkg/database"
)

type Janitor struct {
	Keepalive int `yaml:"keepalive,omitempty" json:"keepalive,omitempty"`
	StaleJob  int `yaml:"stale_job,omitempty" json:"stale_job,omitempty"`
	db        *database.Database
}

func (j *Janitor) Start(d *database.Database) error {

	// record database context
	j.db = d

	// let's run jobs cleanup once before we start it as a go routine
	if err := j.cleanupStaleJobs(); err != nil {
		return err
	}

	// start cleanup loop
	go func() {

		for {

			if err := j.cleanupStaleJobs(); err != nil {
				fmt.Println(`Janitor error:`, err)
			}

			time.Sleep(60 * time.Second)

		}

	}()

	return nil

}
