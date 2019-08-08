package executer

import (
	"context"

	"github.com/drone-runners/drone-runner-exec/engine"
	"github.com/drone/runner-go/pipeline"
)

// Executer ...
type Executer interface {
	Exec(context.Context, *engine.Spec, *pipeline.State) error
}

type executer struct {
}

// New ...
func New() Executer {
	e := new(executer)

	return e
}

// Exec ...
func (e *executer) Exec(ctx context.Context, spec *engine.Spec, state *pipeline.State) error {
	return nil
}
