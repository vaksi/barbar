package cmd

import (
	authtransport "barbar/domain/auth/transport"
	"github.com/spf13/cobra"
)

func ServeHTTPAuth() *cobra.Command {
	return &cobra.Command{
		Use:   "http-auth",
		Short: "use http-auth",
		Long:  `Http Auth serve can lister API auth service`,
		Run: func(cmd *cobra.Command, args []string) {
			authtransport.RunHttp()
		},
	}
}