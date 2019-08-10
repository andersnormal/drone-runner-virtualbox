// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package engine

import (
	"context"
	"io"
	"os"
)

// New returns a new engine.
func New() Engine {
	return new(engine)
}

type engine struct{}

// Setup the pipeline environment.
func (e *engine) Setup(ctx context.Context, spec *Spec) error {
	return nil
}

// Destroy the pipeline environment.
func (e *engine) Destroy(ctx context.Context, spec *Spec) error {
	return os.RemoveAll(spec.Root)
}

// Run runs the pipeline step.
func (e *engine) Run(ctx context.Context, spec *Spec, step *Step, output io.Writer) (*State, error) {
  state := &State{
		ExitCode:  0,
		Exited:    true,
		OOMKilled: false,
	}

	return state, nil
}

//
// Not Implemented
//

// Create creates the pipeline step.
func (e *engine) Create(ctx context.Context, spec *Spec, step *Step) error {
	return nil // no-op for bash implementation
}

// Start the pipeline step.
func (e *engine) Start(context.Context, *Spec, *Step) error {
	return nil // no-op for bash implementation
}

// Wait for the pipeline step to complete and returns the completion results.
func (e *engine) Wait(context.Context, *Spec, *Step) (*State, error) {
	return nil, nil // no-op for bash implementation
}

// Tail the pipeline step logs.
func (e *engine) Tail(context.Context, *Spec, *Step) (io.ReadCloser, error) {
	return nil, nil // no-op for bash implementation
}
