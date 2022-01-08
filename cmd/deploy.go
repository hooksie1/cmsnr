package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/hooksie1/cmsnr/pkg/deployment"
	"os"
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
	vService := "cmsnr-validating-webhook"
	//secret := viper.GetString("secret")
	port := viper.GetInt("port")
	namespace := viper.GetString("namespace")

	mCert, mKey, err := deployment.GenerateCertificate(mService, namespace)
	if err != nil {
		log.Error(err)
		os.Exit(2)
	}

	vCert, vKey, err := deployment.GenerateCertificate(vService, namespace)
	if err != nil {
		log.Error(err)
		os.Exit(2)
	}

	mw := deployment.NewMutatingWebhookServer().NamespacedName(mService, namespace).MutatingWebhook(port, mCert).Rules()
	vw := deployment.NewValidatingWebhookServer().NamespacedName(vService, namespace).ValidatingWebhook(port, vCert).Rules()

	printKind(deployment.NewSA(namespace))
	printKind(deployment.NewClusterRole())
	printKind(deployment.NewClusterRolebinding(namespace))
	fmt.Println(deployment.NewCRD())
	//printKind(deployment.NewMutatingDeployment(service, ))
	//printKind(deployment.NewMutatingService())
	printKind(deployment.CertAsSecret(mCert, mKey, "mutating-secret", namespace))
	printKind(deployment.CertAsSecret(vCert, vKey, "validating-secret", namespace))
	printKind(mw)
	printKind(vw)

}
