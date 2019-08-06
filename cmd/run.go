package cmd

import (
	"context"

	"github.com/andersnormal/drone-runner-virtualbox/run"

	"github.com/andersnormal/pkg/server"
	"github.com/drone/runner-go/client"
	"github.com/drone/runner-go/logger"
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

	// create new runner
	r := run.New(cfg, c, root.logger)
	s.Listen(r, true)

	// listen for the server and wait for it to fail,
	// or for sys interrupts
	if err := s.Wait(); err != nil {
		root.logger.Error(err)
	}

	// noop
	return nil
}
