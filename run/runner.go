package run

import (
	"context"
	"fmt"
	"sync"

	"github.com/andersnormal/drone-runner-virtualbox/config"

	"github.com/andersnormal/pkg/server"
	"github.com/drone/runner-go/client"

	// s "github.com/drone/runner-go/server"
	log "github.com/sirupsen/logrus"
)

// Runner ...
type Runner interface {
	server.Listener
}

type runner struct {
	opts *Opts
	cfg  *config.Config

	client *client.HTTPClient
	logger *log.Entry

	exit    chan struct{}
	errOnce sync.Once
	err     error
	wg      sync.WaitGroup
}

// Opt ...
type Opt func(*Opts)

// Opts ...
type Opts struct {
	Cap int
}

// New ...
func New(cfg *config.Config, client *client.HTTPClient, logger *log.Entry, opts ...Opt) Runner {
	options := new(Opts)

	r := new(runner)
	r.cfg = cfg
	r.opts = options
	r.logger = logger
	r.client = client

	configure(r, opts...)

	return r
}

// Start ...
func (r *runner) Start(ctx context.Context, ready func()) func() error {
	return func() error {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		//   engine := engine.New()
		// remote := remote.New(cli)
		// tracer := history.New(remote)
		// hook := loghistory.New()
		// logrus.AddHook(hook)

		//   server := server.Server{
		//     Addr: "8080",
		//     Handler: router.New(tracer, hook, router.Config{
		//       Username: config.Dashboard.Username,
		//       Password: config.Dashboard.Password,
		//       Realm:    config.Dashboard.Realm,
		//     }),
		//   }

		for i := 0; i < r.opts.Cap; i++ {
			r.run(r.poll(ctx))
		}

		ready()

		if err := r.wait(); err != nil {
			return err
		}

		return nil
	}
}

// Stop is stopping the queue
func (r *runner) Stop() error {
	return nil
}

func (r *runner) poll(ctx context.Context) func() error {
	return func() error {
		stage, err := r.client.Request(ctx, nil)
		if err != nil {
			return err
		}

		fmt.Println(stage)

		for {
			select {
			case <-ctx.Done():
				return nil
			}
		}

		//   log := logger.FromContext(ctx).WithField("thread", thread)
		// log.WithField("thread", thread).Debug("request stage from remote server")

		// request a new build stage for execution from the central
		// build server.
		// stage, err := p.Client.Request(ctx, p.Filter)
		// if err != nil {
		// 	log.WithError(err).Error("cannot request stage")
		// 	return err
		// }

		// exit if a nil or empty stage is returned from the system
		// and allow the runner to retry.
		// if stage == nil || stage.ID == 0 {
		// 	return nil
		// }

		// return p.Runner.Run(
		// 	logger.WithContext(noContext, log), stage)

		return nil
	}
}

func (r *runner) wait() error {
	r.wg.Wait()

	return r.err
}

func (r *runner) run(f func() error) {
	r.wg.Add(1)

	go func() {
		defer r.wg.Done()

		if err := f(); err != nil {
			r.errOnce.Do(func() {
				r.err = err
			})
		}
	}()
}

func configure(r *runner, opts ...Opt) error {
	for _, o := range opts {
		o(r.opts)
	}

	if r.opts.Cap == 0 {
		r.opts.Cap = 1
	}

	return nil
}
