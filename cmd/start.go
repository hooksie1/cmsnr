package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start is the entrypoint for the admission webhooks",
}

func init() {
	serverCmd.AddCommand(startCmd)

	startCmd.PersistentFlags().IntP("port", "p", 8443, "Port the server should use")
	viper.BindPFlag("port", startCmd.PersistentFlags().Lookup("port"))
}
