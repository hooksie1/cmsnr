package deployment

import (
	"context"
	"fmt"
	admissionv1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func NewMutatingWebhookConfig(name, namespace string, port int, cert []byte) *admissionv1.MutatingWebhookConfiguration {
	none := admissionv1.SideEffectClassNone
	var path string
	path = "/mutate"
	var p int32
	p = int32(port)
	seconds := int32(5)

	webhook := admissionv1.MutatingWebhook{
		Name: fmt.Sprintf("%s.%s.svc.cluster.local", name, namespace),
		AdmissionReviewVersions: []string{
			"v1",
			"v1alpha1",
		},
		SideEffects:    &none,
		TimeoutSeconds: &seconds,
		ClientConfig: admissionv1.WebhookClientConfig{
			Service: &admissionv1.ServiceReference{
				Name:      fmt.Sprintf("%s", name),
				Namespace: namespace,
				Path:      &path,
				Port:      &p,
			},
			CABundle: cert,
		},
		Rules: getRules(),
		ObjectSelector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"cmsnr.com/inject": "enabled",
			},
		},
	}

	return &admissionv1.MutatingWebhookConfiguration{
		TypeMeta: metav1.TypeMeta{
			Kind:       "MutatingWebhookConfiguration",
			APIVersion: "admissionregistration.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Webhooks: []admissionv1.MutatingWebhook{
			webhook,
		},
	}

}

func getOperations() []admissionv1.OperationType {
	return []admissionv1.OperationType{
		"CREATE",
		"UPDATE",
	}
}

func getRules() []admissionv1.RuleWithOperations {
	return []admissionv1.RuleWithOperations{
		{
			Operations: getOperations(),
			Rule: admissionv1.Rule{
				APIGroups:   []string{""},
				APIVersions: []string{"v1"},
				Resources:   []string{"pods"},
			},
		},
	}
}

func CreateWebhook(c *admissionv1.MutatingWebhookConfiguration) error {
	ctx := context.Background()
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	_, err = clientSet.AdmissionregistrationV1().MutatingWebhookConfigurations().Create(ctx, c, metav1.CreateOptions{})

	return err

}
