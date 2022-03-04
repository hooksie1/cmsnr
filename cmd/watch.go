package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hooksie1/cmsnr/api/v1alpha1"
	"github.com/hooksie1/cmsnr/pkg/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch will watch the cluster for changes in OPA CRDs.",
	Run:   watchPolicies,
}

func init() {

	opaCmd.AddCommand(watchCmd)

	watchCmd.Flags().StringP("deployment", "d", "", "The deployment name of the OPA policy.")
	viper.BindPFlag("deployment", watchCmd.Flags().Lookup("deployment"))
	watchCmd.MarkFlagRequired("deployment")
}

func watchPolicies(cmd *cobra.Command, args []string) {
	deploymentName := viper.GetString("deployment")

	v1alpha1.BuildScheme(scheme.Scheme)

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	client, err := client.NewClient(config, namespace)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	client.WatchResources()

	log.Infof("watching for policy changes with deployment name: %s", deploymentName)
	for v := range client.Queue {
		if v.OpaPolicy.Spec.DeploymentName != deploymentName {
			continue
		}

		if v.OpaPolicy.Spec.DeploymentName == "" {
			continue
		}

		switch v.Method {
		case "add":
			log.Infof("found new opa policy %s for deployment %s", v.OpaPolicy.Spec.PolicyName, v.OpaPolicy.Spec.DeploymentName)
			if err := updatePolicy(&v.OpaPolicy); err != nil {
				log.Error(err)
				os.Exit(1)
			}
		case "update":
			log.Infof("found update for opa policy %s for deployment %s", v.OpaPolicy.Spec.PolicyName, v.OpaPolicy.Spec.DeploymentName)
			if err := updatePolicy(&v.OpaPolicy); err != nil {
				log.Error(err)
				os.Exit(1)
			}
		case "delete":
			log.Infof("deleting opa policy %s for deployment %s", v.OpaPolicy.Spec.PolicyName, v.OpaPolicy.Spec.DeploymentName)
			if err := deletePolicy(&v.OpaPolicy); err != nil {
				log.Error(err)
				os.Exit(1)
			}
		}

		time.Sleep(2 * time.Second)
	}

	close(client.Queue)

}

func updatePolicy(p *v1alpha1.OpaPolicy) error {
	policy := p.Spec.Policy
	url := fmt.Sprintf("http://localhost:8181/v1/policies/%s", p.Spec.PolicyName)
	client := http.Client{}

	req, err := http.NewRequest("PUT", url, strings.NewReader(policy))
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("error posting policy")
	}

	return nil
}

func deletePolicy(p *v1alpha1.OpaPolicy) error {
	url := fmt.Sprintf("http://localhost:8181/v1/policies/%s", p.Spec.PolicyName)
	client := http.Client{}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("error deleting policy")
	}

	return nil
}
