package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	CRDGroup   string = "cmsnr.com"
	CRDVersion string = "v1alpha1"
)

var SchemeGroupVersion = schema.GroupVersion{Group: CRDGroup, Version: CRDVersion}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(BuildScheme)
	AddToScheme   = SchemeBuilder.AddToScheme
)

func BuildScheme(scheme *runtime.Scheme) error {
	s := schema.GroupVersion{
		Group:   CRDGroup,
		Version: CRDVersion,
	}
	scheme.AddKnownTypes(s,
		&OpaPolicy{},
		&OpaPolicyList{},
	)

	metav1.AddToGroupVersion(scheme, s)
	return nil
}
