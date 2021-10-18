package main

import (
	"barbar/cmd"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use: "barbar",
	}

	// Add Command
	rootCmd.AddCommand(cmd.ServeHTTPUser(), cmd.ServeGRPCUser(), cmd.ServeHTTPAuth(), cmd.InitIndexMongo())

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
