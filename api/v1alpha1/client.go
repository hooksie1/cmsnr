package v1alpha1

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type OpaPolicyInterface interface {
	List(opts metav1.ListOptions) (*OpaPolicyList, error)
	Get(name string, options metav1.GetOptions) (*OpaPolicy, error)
	Create(*OpaPolicy) (*OpaPolicy, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

type opaPolicyClient struct {
	restClient rest.Interface
	ns         string
}

func (c *opaPolicyClient) List(opts metav1.ListOptions) (*OpaPolicyList, error) {
	result := OpaPolicyList{}
	ctx := context.Background()
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("opapolicies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *opaPolicyClient) Get(name string, opts metav1.GetOptions) (*OpaPolicy, error) {
	result := OpaPolicy{}
	ctx := context.Background()
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *opaPolicyClient) Create(opapolicy *OpaPolicy) (*OpaPolicy, error) {
	result := OpaPolicy{}
	ctx := context.Background()
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("opapolicies").
		Body(opapolicy).
		Do(ctx).
		Into(&result)

	return &result, err

}

func (c *opaPolicyClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	ctx := context.Background()
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("opapolicies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)

}
