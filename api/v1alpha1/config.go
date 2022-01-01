package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type OpaV1Alpha1Interface interface {
	OpaPolicies(namespace string) OpaPolicyInterface
}

type OpaV1Alpha1Client struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*OpaV1Alpha1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &OpaV1Alpha1Client{restClient: client}, nil
}

func (c *OpaV1Alpha1Client) OpaPolicies(namespace string) OpaPolicyInterface {
	return &opaPolicyClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}

func setConfigDefaults(config *rest.Config) error {
	gv := &schema.GroupVersion{Group: CRDGroup, Version: CRDVersion}
	config.ContentConfig.GroupVersion = gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}
