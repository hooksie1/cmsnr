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

// mutatingCmd represents the mutating command
var mutatingCmd = &cobra.Command{
	Use:   "mutating",
	Short: "Mutating starts the cmsnr mutating webhook",
	Run:   mutateServer,
}

func init() {
	startCmd.AddCommand(mutatingCmd)
}

func mutateServer(cmd *cobra.Command, args []string) {
	port := viper.GetInt("port")
	registry := viper.GetString("registry")
	log.Debugf("mutating webhook port: %d", port)
	log.Info("setting up webhook server")
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	injector := server.SidecarInjector{
		Client:    mgr.GetClient(),
		Namespace: namespace,
		Registry:  registry,
	}
	log.Info("setting up server")
	mgrServer := mgr.GetWebhookServer()
	mgrServer.Port = port
	mgrServer.CertDir = "/var/lib/cmsnr"
	mgrServer.Register("/mutate", &webhook.Admission{
		Handler: &injector,
	})

	log.Info("starting webhook server")
	if err := mgrServer.Start(signals.SetupSignalHandler()); err != nil {
		log.Errorf("unable to start webhook server: %s", err)
		os.Exit(1)
	}
}
