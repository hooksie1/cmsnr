package cmd

import (
	"os"

	"github.com/hooksie1/cmsnr/pkg/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// validatingCmd represents the validating command
var validatingCmd = &cobra.Command{
	Use:   "validating",
	Short: "Starts the cmsnr validating webhook",
	Run:   validateServer,
}

func init() {
	startCmd.AddCommand(validatingCmd)
}

func validateServer(cmd *cobra.Command, args []string) {
	port := viper.GetInt("port")
	log.Debugf("validating webhook port: %d", port)
	log.Info("setting up webhook server")
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	validator := server.Validator{
		Client: mgr.GetClient(),
	}

	mgrServer := mgr.GetWebhookServer()
	mgrServer.Port = port
	mgrServer.CertDir = "/var/lib/cmsnr"
	mgrServer.Register("/validate", &webhook.Admission{
		Handler: &validator,
	})

	log.Info("starting webhook server")
	if err := mgrServer.Start(signals.SetupSignalHandler()); err != nil {
		log.Errorf("unable to start webhook server: %s", err)
		os.Exit(1)
	}

}
