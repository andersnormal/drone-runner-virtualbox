package config

import (
	"os"
	"runtime"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

// Runner ...
type Runner struct {
	Procs int64
}

// Config contains a configuration for drone-runner-virtualbox
type Config struct {
	// LogLevel is the level with with to log for this config
	LogLevel log.Level

	// ReloadSignal
	ReloadSignal syscall.Signal

	// TermSignal
	TermSignal syscall.Signal

	// KillSignal
	KillSignal syscall.Signal

	// Timeout of the runtime
	Timeout time.Duration

	// DroneRPCAddress ...
	DroneRPCAddress string

	// DroneRPCSecret ...
	DroneRPCSecret string

	// DroneRPCCapacity ...
	DroneRPCCapacity int

	// Runner ...
	Runner Runner
}

const (
	// DefaultLogLevel is the default logging level.
	DefaultLogLevel = log.InfoLevel

	// DefaultTermSignal is the signal to term the agent.
	DefaultTermSignal = syscall.SIGTERM

	// DefaultReloadSignal is the default signal for reload.
	DefaultReloadSignal = syscall.SIGHUP

	// DefaultKillSignal is the default signal for termination.
	DefaultKillSignal = syscall.SIGINT

	// DefaultDroneRPCAddress is the default address of the Drone server
	DefaultDroneRPCAddress = "http://localhost"

	// DefaultDroneRPCSecret is the default secret of the Drone server
	DefaultDroneRPCSecret = "magic_secret"

	// DefaultDroneRPCCapacity ...
	DefaultDroneRPCCapacity = 1
)

// New returns a new Config
func New() *Config {
	return &Config{
		LogLevel:         DefaultLogLevel,
		ReloadSignal:     DefaultReloadSignal,
		TermSignal:       DefaultTermSignal,
		KillSignal:       DefaultKillSignal,
		DroneRPCAddress:  DefaultDroneRPCAddress,
		DroneRPCSecret:   DefaultDroneRPCSecret,
		DroneRPCCapacity: DefaultDroneRPCCapacity,
		Runner:           Runner{},
	}
}

// Name ...
func (cfg *Config) Name() string {
	name, _ := os.Hostname()

	return name
}

// OS ...
func (cfg *Config) OS() string {
	return runtime.GOOS
}

// Arch ...
func (cfg *Config) Arch() string {
	return runtime.GOARCH
}
