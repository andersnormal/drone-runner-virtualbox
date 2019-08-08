package poller

import (
	"context"
	"sync"
	"time"

	"github.com/andersnormal/drone-runner-virtualbox/config"
	"github.com/andersnormal/drone-runner-virtualbox/runner"

	"github.com/andersnormal/pkg/server"
	"github.com/drone/drone-go/drone"
	"github.com/drone/runner-go/client"

	// s "github.com/drone/runner-go/server"
	log "github.com/sirupsen/logrus"
)

// Poller ...
type Poller interface {
	server.Listener
}

type poller struct {
	opts *Opts
	cfg  *config.Config

	client *client.HTTPClient
	logger *log.Entry

	runner runner.Runner

	stage   chan *drone.Stage
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
func New(cfg *config.Config, runner runner.Runner, client *client.HTTPClient, logger *log.Entry, opts ...Opt) Poller {
	options := new(Opts)

	p := new(poller)
	p.cfg = cfg
	p.opts = options
	p.logger = logger
	p.client = client
	p.runner = runner

	// configure channel
	p.stage = make(chan *drone.Stage)

	configure(p, opts...)

	return p
}

// Start ...
func (p *poller) Start(ctx context.Context, ready func()) func() error {
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

		// start

		for i := 0; i < p.opts.Cap; i++ {
			p.run(p.watch(ctx))
		}

		// call for being ready, and start the next service
		time.Sleep(1 * time.Second)
		ready()

		// wait for the group
		if err := p.wait(); err != nil {
			return err
		}

		return nil
	}
}

// Stop is stopping the queue
func (p *poller) Stop() error {
	return nil
}

func (p *poller) poll(ctx context.Context) func() error {
	return func() error {
		for {
			// this is to direct the request to the runner
			stage, err := p.client.Request(ctx, &client.Filter{
				Kind: "pipeline",
				Type: "virtualbox",
			})
			if err != nil {
				return err
			}

			p.stage <- stage
		}
	}
}

func (p *poller) watch(ctx context.Context) func() error {
	return func() error {
		// start to poll
		p.run(p.poll(ctx))

		// looking for channel
		for {
			select {
			case stage, ok := <-p.stage:
				if !ok {
					return nil
				}

				if stage == nil || stage.ID == 0 {
					continue
				}

				// this is sync, running in its own loop
				p.staging(ctx, stage)

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
	}
}

func (p *poller) staging(ctx context.Context, stage *drone.Stage) error {
	stage.Machine = "test"

	err := p.client.Accept(ctx, stage)
	if err != nil {
		log.WithError(err).Error("cannot accept stage")
		return err
	}

	// get data
	data, err := p.client.Detail(ctx, stage)
	if err != nil {
		log.WithError(err).Error("cannot get stage details")

		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	p.run(func() error {
		done, err := p.client.Watch(ctx, data.Build.ID)
		if err != nil {
			return err
		}

		if done {
			cancel()
		}

		return nil
	})

	// create new run
	s, err := p.runner.Run(ctx, stage, data)
	if err != nil {
		return err
	}

	// update the the stage on the server
	if err := p.client.Update(ctx, s); err != nil {
		return err
	}

	return nil
}

func (p *poller) wait() error {
	p.wg.Wait()

	return p.err
}

func (p *poller) run(f func() error) {
	p.wg.Add(1)

	go func() {
		defer p.wg.Done()

		if err := f(); err != nil {
			p.errOnce.Do(func() {
				p.err = err
			})
		}
	}()
}

func configure(p *poller, opts ...Opt) error {
	for _, o := range opts {
		o(p.opts)
	}

	if p.opts.Cap == 0 {
		p.opts.Cap = 1
	}

	return nil
}
