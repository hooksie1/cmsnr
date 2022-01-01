package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/hooksie1/cmsnr/api/v1alpha1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch will watch the cluster for changes in OPA CRDs.",
	Run:   watchPolicies,
}

var versions map[string]string

func init() {
	versions = make(map[string]string)

	opaCmd.AddCommand(watchCmd)

	watchCmd.Flags().StringP("deployment", "d", "", "The deployment name of the OPA policy.")
	viper.BindPFlag("deployment", watchCmd.Flags().Lookup("deployment"))
	watchCmd.MarkFlagRequired("deployment")
}

func watchPolicies(cmd *cobra.Command, args []string) {
	var config *rest.Config
	var err error

	deploymentName := viper.GetString("deployment")

	config, err = rest.InClusterConfig()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	v1alpha1.BuildScheme(scheme.Scheme)

	clientSet, err := v1alpha1.NewForConfig(config)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	store := WatchResources(clientSet, namespace)

	log.Infof("watching for policy changes with deployment name: %s", deploymentName)
	for {
		var name string
		var version string
		var depName string
		opaFromStore := store.List()
		for _, v := range opaFromStore {
			name = v.(*v1alpha1.OpaPolicy).Name
			depName = v.(*v1alpha1.OpaPolicy).Spec.DeploymentName
			version = v.(*v1alpha1.OpaPolicy).ResourceVersion

			if depName != deploymentName {
				continue
			}

			if depName == "" {
				versions[name] = version
			}

			if versions[name] == version {
				continue
			}

			log.Infof("Found policy change for policy %s, updating policy", name)
			versions[name] = version
			if err := updatePolicy(v.(*v1alpha1.OpaPolicy)); err != nil {
				log.Error(err)
				os.Exit(1)
			}

		}
		time.Sleep(2 * time.Second)
	}
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
