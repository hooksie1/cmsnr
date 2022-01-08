package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "cmsnrctl server commands",
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.PersistentFlags().String("service", "cmsnr-webhook", "Service name")
	viper.BindPFlag("service", serverCmd.PersistentFlags().Lookup("service"))
	serverCmd.PersistentFlags().String("secret", "cmsnr-secret", "Secret name")
	viper.BindPFlag("secret", serverCmd.PersistentFlags().Lookup("secret"))
	serverCmd.PersistentFlags().String("namespace", "default", "Namespace")
	viper.BindPFlag("namespace", serverCmd.PersistentFlags().Lookup("namespace"))

}

type webhookServer struct {
	serverType string
	service    string
	namespace  string
	name       string
	port       int
	cert       []byte
	key        []byte
	certPath   string
	print      func(interface{}, func(i interface{}))
}
