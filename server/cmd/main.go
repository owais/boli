package main

import (
	"github.com/spf13/cobra"

	"github.com/owais/boli/server/cmd/cli"
	"github.com/owais/boli/server/cmd/net"
)

func main() {
	var cmd = &cobra.Command{
		Use:   "boli",
		Short: "Boli card game",
		Long:  `Boli is a six-player card game.`,
	}

	cmd.AddCommand(newCliCmd())
	cmd.AddCommand(serverCmd())

	if err := cmd.Execute(); err != nil {
		panic(err)
	}

}

func newCliCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "cli",
		Short: "CLI client for the game",
		Long:  `CLI client for the game`,
		Run: func(cmd *cobra.Command, args []string) {
			cli.Run()
		},
	}

	return cmd
}

func serverCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "server",
		Short: "Server for the game",
		Long:  `Server for the game`,
		Run: func(cmd *cobra.Command, args []string) {
			net.Run()
		},
	}
	return cmd
}
