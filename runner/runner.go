package runner

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/drone-runners/drone-runner-exec/engine"
	"github.com/drone-runners/drone-runner-exec/engine/compiler"
	"github.com/drone-runners/drone-runner-exec/engine/resource"
	"github.com/drone-runners/drone-runner-exec/runtime"

	"github.com/drone/drone-go/drone"
	"github.com/drone/envsubst"
	"github.com/drone/runner-go/client"
	"github.com/drone/runner-go/environ"
	"github.com/drone/runner-go/manifest"
	"github.com/drone/runner-go/pipeline"
	"github.com/drone/runner-go/secret"
	log "github.com/sirupsen/logrus"
)

// Runner ...
type Runner interface {
	// Run is running a build stage on the runner
	Run(ctx context.Context, stage *drone.Stage, details *client.Context) (*drone.Stage, error)
}

type runner struct {
	logger *log.Entry

	env map[string]string

	client   *client.HTTPClient
	exec     runtime.Execer
	reporter pipeline.Reporter
	secret   secret.Provider
	match    func(*drone.Repo, *drone.Build) bool
}

// New ...
func New(client *client.HTTPClient, reporter pipeline.Reporter, exec runtime.Execer, env map[string]string, secret secret.Provider, match func(*drone.Repo, *drone.Build) bool, ll *log.Entry) Runner {
	r := new(runner)

	r.logger = ll
	r.reporter = reporter
	r.match = match
	r.secret = secret
	r.env = env
	r.exec = exec
	r.client = client

	return r
}

// Run ...
func (r *runner) Run(ctx context.Context, stage *drone.Stage, details *client.Context) (*drone.Stage, error) {
	// construct some logging information
	ll := r.log().WithFields(log.Fields{
		"repo.id":        details.Repo.ID,
		"repo.namespace": details.Repo.Namespace,
		"repo.name":      details.Repo.Name,
		"build.id":       details.Build.ID,
		"build.number":   details.Build.Number,
	})

	// log ...
	ll.Infof("starting runner")

	// creating timeout according to the repo
	timeout := time.Duration(details.Repo.Timeout) * time.Minute
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// create the env for the run
	envs := environ.Combine(
		// s.Environ,
		environ.System(details.System),
		environ.Repo(details.Repo),
		environ.Build(details.Build),
		environ.Stage(stage),
		environ.Link(details.Repo, details.Build, details.System),
		details.Build.Params,
	)

	// this is to start the executer
	// string substitution function ensures that string
	// replacement variables are escaped and quoted if they
	// contain a newline character.
	subf := func(k string) string {
		v := envs[k]
		if strings.Contains(v, "\n") {
			v = fmt.Sprintf("%q", v)
		}
		return v
	}

	// setting up the state for the run
	state := &pipeline.State{
		Build:  details.Build,
		Stage:  stage,
		Repo:   details.Repo,
		System: details.System,
	}

	// evaluates whether or not the agent can process the
	// pipeline. An agent may choose to reject a repository
	// or build for security reasons.
	if r.match != nil && !r.match(details.Repo, details.Build) {
		state.FailAll(errors.New("insufficient permission to run the pipeline"))

		return nil, r.reporter.ReportStage(ctx, state)
	}

	// evaluates string replacement expressions and returns an
	// update configuration file string.
	config, err := envsubst.Eval(string(details.Config.Data), subf)
	if err != nil {
		state.FailAll(err)

		return nil, r.reporter.ReportStage(ctx, state)
	}

	// parse the yaml configuration file.
	manifest, err := manifest.ParseString(config)
	if err != nil {
		state.FailAll(err)

		return nil, r.reporter.ReportStage(ctx, state)
	}

	fmt.Println(stage.Name, manifest, config)

	// find the named stage in the yaml configuration file.
	resource, err := resource.Lookup(stage.Name, manifest)
	if err != nil {
		state.FailAll(err)

		fmt.Println(err)

		return nil, r.reporter.ReportStage(ctx, state)
	}

	secrets := secret.Combine(
		secret.Static(details.Secrets),
		secret.Encrypted(),
		r.secret,
	)

	// compile the yaml configuration file to an intermediate
	// representation, and then
	comp := &compiler.Compiler{
		Pipeline: resource,
		Manifest: manifest,
		Environ:  r.env,
		Build:    details.Build,
		Stage:    stage,
		Repo:     details.Repo,
		System:   details.System,
		Netrc:    details.Netrc,
		Secret:   secrets,
	}

	spec := comp.Compile(ctx)
	for _, src := range spec.Steps {
		// steps that are skipped are ignored and are not stored
		// in the drone database, nor displayed in the UI.
		if src.RunPolicy == engine.RunNever {
			continue
		}
		stage.Steps = append(stage.Steps, &drone.Step{
			Name:      src.Name,
			Number:    len(stage.Steps) + 1,
			StageID:   stage.ID,
			Status:    drone.StatusPending,
			ErrIgnore: src.IgnoreErr,
		})
	}

	stage.Started = time.Now().Unix()
	stage.Status = drone.StatusRunning

	if err := r.client.Update(ctx, stage); err != nil {
		return nil, err
	}

	// go to executer
	ll.Info("should execute here")

	if err = r.exec.Exec(ctx, spec, state); err != nil {
		return nil, err
	}

	return stage, nil
}

func (r *runner) log() *log.Entry {
	return r.logger
}
