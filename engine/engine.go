package engine

import (
	"context"
	"fmt"
	"io"

	e "github.com/drone-runners/drone-runner-exec/engine"
)

type engine struct {
}

func New() e.Engine {
	ee := new(engine)

	return ee
}

// Setup the pipeline environment.
func (ee *engine) Setup(ctx context.Context, spec *e.Spec) error {
	fmt.Println("setup")

	return nil
}

// Run runs the pipeine step.
func (ee *engine) Run(ctx context.Context, spec *e.Spec, step *e.Step, w io.Writer) (*e.State, error) {
	state := &e.State{
		ExitCode:  0,
		Exited:    true,
		OOMKilled: false,
	}

	return state, nil
}

// Create creates the pipeline state.
func (ee *engine) Create(ctx context.Context, spec *e.Spec, step *e.Step) error {
	return nil
}

// Start the pipeline step.
func (ee *engine) Start(ctx context.Context, spec *e.Spec, step *e.Step) error {
	return nil
}

// Wait for the pipeline step to complete and returns the completion results.
func (ee *engine) Wait(ctx context.Context, spec *e.Spec, step *e.Step) (*e.State, error) {
	return nil, nil
}

// Tail the pipeline step logs.
func (ee *engine) Tail(ctx context.Context, spec *e.Spec, step *e.Step) (io.ReadCloser, error) {
	return nil, nil
}

// Destroy the pipeline environment.
func (ee *engine) Destroy(ctx context.Context, spec *e.Spec) error {
	return nil
}
