package run

import (
	"context"
	"time"

	"github.com/drone/drone-go/drone"
	"github.com/drone/runner-go/client"
	"github.com/drone/runner-go/environ"
	log "github.com/sirupsen/logrus"
)

// Stage ...
type Stage struct {
	stage   *drone.Stage
	details *client.Context
	logger  *log.Entry
}

// NewStage ...
func NewStage(stage *drone.Stage, details *client.Context, ll *log.Entry) *Stage {
	s := new(Stage)
	s.details = details
	s.logger = ll
	s.stage = stage

	return s
}

// Run ...
func (s *Stage) Run(ctx context.Context) error {
	// construct some logging information
	ll := s.log().WithFields(log.Fields{
		"repo.id":        s.details.Repo.ID,
		"repo.namespace": s.details.Repo.Namespace,
		"repo.name":      s.details.Repo.Name,
		"build.id":       s.details.Build.ID,
		"build.number":   s.details.Build.Number,
	})

	// creating timeout according to the repo
	timeout := time.Duration(s.details.Repo.Timeout) * time.Minute
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	envs := environ.Combine(
		// s.Environ,
		environ.System(s.details.System),
		environ.Repo(s.details.Repo),
		environ.Build(s.details.Build),
		environ.Stage(s.stage),
		environ.Link(s.details.Repo, s.details.Build, s.details.System),
		s.details.Build.Params,
	)

	return nil
}

func (s *Stage) log() *log.Entry {
	return s.logger
}
