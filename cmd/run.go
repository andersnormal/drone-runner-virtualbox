package cmd

import (
	"context"

	"github.com/andersnormal/drone-runner-virtualbox/executer"
	"github.com/andersnormal/drone-runner-virtualbox/poller"
  "github.com/andersnormal/drone-runner-virtualbox/runner"
  "github.com/andersnormal/drone-runner-virtualbox/match"

	"github.com/andersnormal/pkg/server"
	"github.com/drone/runner-go/client"
	"github.com/drone/runner-go/logger"
	"github.com/drone/runner-go/pipeline/history"
	"github.com/drone/runner-go/pipeline/remote"
	"github.com/drone/runner-go/secret"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type root struct {
	logger *log.Entry
}

func runE(cmd *cobra.Command, args []string) error {
	// this is the main loop
	// create a new root
	root := new(root)

	// init logger
	root.logger = log.WithFields(log.Fields{})

	// create root context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create server
	s := server.NewServer(ctx)

	// create new client
	c := client.New(
		cfg.DroneRPCAddress,
		cfg.DroneRPCSecret,
		true,
	)

	// set logger for the client
	c.Logger = logger.Logrus(
		root.logger,
	)

	// create match function
	m := match.Func(
		[]string{}, // todo, replace with config
		[]string{},
		false,
	)

	// create secret provider
	sec := secret.External(
		"",
		"",
		true,
	)

	// create run env...
	env := map[string]string{}

	// create executer ...
	exec := executer.New()

	// create a new runner
	// engine := engine.New()
	remote := remote.New(c)
	reporter := history.New(remote)
	r := runner.New(reporter, exec, env, sec, m, root.logger)

	// create new poller
	p := poller.New(cfg, r, c, root.logger)
	s.Listen(p, true)

	// listen for the server and wait for it to fail,
	// or for sys interrupts
	if err := s.Wait(); err != nil {
		root.logger.Error(err)
	}

	// noop
	return nil
}
