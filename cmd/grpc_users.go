package cmd

import (
	usertransport "barbar/domain/users/transport"
	"github.com/spf13/cobra"
)

func ServeGRPCUser() *cobra.Command {
	return &cobra.Command{
		Use:   "grpc-user",
		Short: "use grpc-user",
		Long:  `Grpc User serve can lister API user service`,
		Run: func(cmd *cobra.Command, args []string) {
			usertransport.ServeGRPC()
		},
	}
}
