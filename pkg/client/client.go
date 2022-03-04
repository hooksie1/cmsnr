package client

import (
	"github.com/hooksie1/cmsnr/api/v1alpha1"
	"k8s.io/client-go/rest"
)

// Client holds the information for the cmsnr client.
type Client struct {
	// Namespace is the namespace where the client watches for policies.
	Namespace string

	// ClientSet is the k8s clientset
	ClientSet v1alpha1.OpaV1Alpha1Interface

	// Queue is the channel where messages are passed back to the client
	Queue chan v1alpha1.OpaMessage
}

// NewClient returns a new client
func NewClient(config *rest.Config, namespace string) (*Client, error) {

	clientSet, err := v1alpha1.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		Namespace: namespace,
		ClientSet: clientSet,
		Queue:     make(chan v1alpha1.OpaMessage),
	}, nil

}
