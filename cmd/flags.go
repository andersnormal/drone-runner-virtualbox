package cmd

import (
	c "github.com/andersnormal/drone-runner-virtualbox/config"

	"github.com/spf13/cobra"
)

func addFlags(cmd *cobra.Command, cfg *c.Config) {
	// enable verbose output
	cmd.Flags().BoolVar(&cfg.Verbose, "verbose", c.DefaultVerbose, "enable verbose output")

	// enable debug options
	cmd.Flags().BoolVar(&cfg.Debug, "debug", c.DefaultDebug, "enable debug")

	// addr to connect to
	cmd.Flags().StringVar(&cfg.Addr, "addr", c.DefaultAddr, "address of the Xcode server")
}
