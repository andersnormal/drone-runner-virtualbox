package cmd

import (
	c "github.com/andersnormal/drone-runner-virtualbox/config"

	"github.com/spf13/cobra"
)

func addFlags(cmd *cobra.Command, cfg *c.Config) {
	// Drone RPC address ...
	cmd.Flags().StringVar(&cfg.DroneRPCAddress, "rpc-host", c.DefaultDroneRPCAddress, "drone rpc host")

	// Drone RPC address ....
	cmd.Flags().StringVar(&cfg.DroneRPCSecret, "rpc-secret", c.DefaultDroneRPCAddress, "drone rpc secret")

	// Drone RPC capacity ....
	cmd.Flags().IntVar(&cfg.DroneRPCCapacity, "rpc-capacity", c.DefaultDroneRPCCapacity, "drone rpc capcity")
}
