package cmd

import (
	"time"

	"gitlab.com/hooksie1/cmsnr/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

func WatchResources(clientSet v1alpha1.OpaV1Alpha1Interface, namespace string, inform chan<- v1alpha1.OpaMessage) {
	_, opaController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return clientSet.OpaPolicies(namespace).List(lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return clientSet.OpaPolicies(namespace).Watch(lo)
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

					inform <- m
				}
			},
			UpdateFunc: func(oldObj interface{}, newObj interface{}) {
				if r, ok := newObj.(*v1alpha1.OpaPolicy); ok {
					m := v1alpha1.OpaMessage{
						Method:    "update",
						OpaPolicy: *r,
					}

					inform <- m
				}
			},
			AddFunc: func(obj interface{}) {
				if r, ok := obj.(*v1alpha1.OpaPolicy); ok {
					m := v1alpha1.OpaMessage{
						Method:    "add",
						OpaPolicy: *r,
					}

					inform <- m
				}
			},
		},
	)

	go opaController.Run(wait.NeverStop)
}
