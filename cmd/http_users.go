package cmd

import (
	usertransport "barbar/domain/users/transport"
	"github.com/spf13/cobra"
)

func ServeHTTPUser() *cobra.Command {
	return &cobra.Command{
		Use:   "http-user",
		Short: "use http-user",
		Long:  `Http User serve can lister API user service`,
		Run: func(cmd *cobra.Command, args []string) {
			usertransport.RunHttp()
		},
	}
}
