package client

import (
	"time"

	"github.com/hooksie1/cmsnr/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

// Watch resources watches for new OPA resources. If they are found, they are passed
// to the event handler functions that send them as a message to the client.
func (c *Client) WatchResources() {
	_, opaController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return c.ClientSet.OpaPolicies(c.Namespace).List(lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return c.ClientSet.OpaPolicies(c.Namespace).Watch(lo)
			},
		},
		&v1alpha1.OpaPolicy{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{
			DeleteFunc: func(obj interface{}) {
				if r, ok := obj.(*v1alpha1.OpaPolicy); ok {
					m := v1alpha1.OpaMessage{
						Method:    "delete",
						OpaPolicy: *r,
					}

					c.Queue <- m
				}
			},
			UpdateFunc: func(oldObj interface{}, newObj interface{}) {
				if r, ok := newObj.(*v1alpha1.OpaPolicy); ok {
					m := v1alpha1.OpaMessage{
						Method:    "update",
						OpaPolicy: *r,
					}

					c.Queue <- m
				}
			},
			AddFunc: func(obj interface{}) {
				if r, ok := obj.(*v1alpha1.OpaPolicy); ok {
					m := v1alpha1.OpaMessage{
						Method:    "add",
						OpaPolicy: *r,
					}

					c.Queue <- m
				}
			},
		},
	)

	go opaController.Run(wait.NeverStop)
}
