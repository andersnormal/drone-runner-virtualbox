package cmd

import (
	"fmt"
	"os"

	"github.com/andersnormal/drone-runner-virtualbox/config"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfg *config.Config
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "drone-runner-virtualbox",
	Short: "drone virtualbox runner",
	Long: `
    This starts a runner that controls VirtualBox on the machine.
    It may clone or create a new virtual machine for the pipeline.
  `,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: runE,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// init config
	cfg = config.New()

	// silence on the root cmd
	RootCmd.SilenceErrors = true
	RootCmd.SilenceUsage = true

	// initialize cobra
	cobra.OnInitialize(initConfig)

	// adding flags
	addFlags(RootCmd, cfg)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Only log the warning severity or above.
	log.SetLevel(cfg.LogLevel)

	// if we should output verbose
	// if cfg.Verbose || cfg.Debug {
	// 	log.SetLevel(log.InfoLevel)
	// }
}
