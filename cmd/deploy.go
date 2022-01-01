package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/hooksie1/cmsnr/pkg/deployment"
	"sigs.k8s.io/yaml"
)

// deployCmd represents the cert command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "generates certificates for webhook server",
	Run:   generate,
}

func init() {
	serverCmd.AddCommand(deployCmd)
}

func printKind(i interface{}) {
	o, err := yaml.Marshal(i)
	if err != nil {
		log.Error("error printing object: %s", err)
		os.Exit(2)
	}

	fmt.Printf("---\n%s\n", o)
}

func generate(cmd *cobra.Command, args []string) {
	mService := "cmsnr-mutating-webhook"
	name := viper.GetString("secret")
	port := viper.GetInt("port")
	namespace := viper.GetString("namespace")

	mCert, mKey, err := deployment.GenerateCertificate(mService, namespace)
	if err != nil {
		log.Error(err)
		os.Exit(2)
	}

	mw := webhookServer{
		service:   mService,
		namespace: namespace,
		name:      name,
		port:      port,
		cert:      mCert,
		key:       mKey,
	}

	mw.printServiceAccount()

	mw.printClusterRole()

	mw.printClusterRoleBinding()

	mw.printCRD()

	mw.printMutatingDeployment()

	mw.printMutatingService()

	mw.printMutatingSecret()

	mw.printMutatingWebhook()

}

func (w *webhookServer) printServiceAccount() {
	printKind(deployment.NewSA(w.namespace))
}

func (w *webhookServer) printClusterRole() {
	printKind(deployment.NewClusterRole())
}

func (w *webhookServer) printClusterRoleBinding() {
	printKind(deployment.NewClusterRolebinding(w.namespace))
}

func (w *webhookServer) printCRD() {
	fmt.Println(deployment.NewCRD())
}

func (w *webhookServer) printMutatingDeployment() {
	printKind(deployment.NewDeployment(w.service, w.namespace, w.port))
}

func (w *webhookServer) printMutatingService() {
	printKind(deployment.NewService(w.service, w.namespace, w.port))
}

func (w *webhookServer) printMutatingSecret() {
	printKind(deployment.CertAsSecret(w.cert, w.key, w.name, w.namespace))
}

func (w *webhookServer) printMutatingWebhook() {
	printKind(deployment.NewMutatingWebhookConfig(w.service, w.namespace, w.port, w.cert))
}
