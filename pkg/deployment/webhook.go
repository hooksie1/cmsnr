package deployment

import (
	"fmt"
	admissionv1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type admissionServer struct {
	Name      string
	Namespace string
	Port      int
	Cert      []byte
}
type MutatingServer struct {
	admissionServer
	Config *admissionv1.MutatingWebhookConfiguration
}

type ValidatingServer struct {
	admissionServer
	Config *admissionv1.ValidatingWebhookConfiguration
}

func NewMutatingWebhookServer() *MutatingServer {
	return &MutatingServer{
		Config: &admissionv1.MutatingWebhookConfiguration{
			TypeMeta: metav1.TypeMeta{
				Kind:       "MutatingWebhookConfiguration",
				APIVersion: "admissionregistration.k8s.io/v1",
			},
		},
	}
}

func (m *MutatingServer) NamespacedName(name, namespace string) *MutatingServer {
	m.Name = name
	m.Namespace = namespace

	m.Config.ObjectMeta = metav1.ObjectMeta{
		Name:      m.Name,
		Namespace: m.Namespace,
	}
	return m
}

func (m *MutatingServer) MutatingWebhook(port int, cert []byte) *MutatingServer {
	m.Port = port
	m.Cert = cert

	m.Config.Webhooks = []admissionv1.MutatingWebhook{
		m.newMutatingWebhook(),
	}

	return m
}

func (m *MutatingServer) newMutatingWebhook() admissionv1.MutatingWebhook {
	none := admissionv1.SideEffectClassNone
	var path string
	path = "/mutate"
	var p int32
	p = int32(m.Port)
	seconds := int32(5)

	return admissionv1.MutatingWebhook{
		Name: fmt.Sprintf("%s.%s.svc.cluster.local", m.Name, m.Namespace),
		AdmissionReviewVersions: []string{
			"v1",
			"v1alpha1",
		},
		SideEffects:    &none,
		TimeoutSeconds: &seconds,
		ClientConfig: admissionv1.WebhookClientConfig{
			Service: &admissionv1.ServiceReference{
				Name:      fmt.Sprintf("%s", m.Name),
				Namespace: m.Namespace,
				Path:      &path,
				Port:      &p,
			},
			CABundle: m.Cert,
		},
		ObjectSelector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"cmsnr.com/inject": "enabled",
			},
		},
	}
}

func (m *MutatingServer) Rules() *MutatingServer {
	m.Config.Webhooks[0].Rules = []admissionv1.RuleWithOperations{
		{
			Operations: []admissionv1.OperationType{
				"CREATE",
				"UPDATE",
			},
			Rule: admissionv1.Rule{
				APIGroups:   []string{""},
				APIVersions: []string{"v1"},
				Resources:   []string{"pods"},
			},
		},
	}

	return m
}

func NewValidatingWebhookServer() *ValidatingServer {
	return &ValidatingServer{
		Config: &admissionv1.ValidatingWebhookConfiguration{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ValidatingWebhookConfiguration",
				APIVersion: "admissionregistration.k8s.io/v1",
			},
		},
	}
}

func (v *ValidatingServer) NamespacedName(name, namespace string) *ValidatingServer {
	v.Name = name
	v.Namespace = namespace

	v.Config.ObjectMeta = metav1.ObjectMeta{
		Name:      v.Name,
		Namespace: v.Namespace,
	}
	return v
}

func (v *ValidatingServer) ValidatingWebhook(port int, cert []byte) *ValidatingServer {
	v.Port = port
	v.Cert = cert

	v.Config.Webhooks = []admissionv1.ValidatingWebhook{
		v.newValidatingWebhook(),
	}

	return v
}

func (v *ValidatingServer) newValidatingWebhook() admissionv1.ValidatingWebhook {
	none := admissionv1.SideEffectClassNone
	var path string
	path = "/validate"
	var p int32
	p = int32(v.Port)
	seconds := int32(5)

	return admissionv1.ValidatingWebhook{
		Name: fmt.Sprintf("%s.%s.svc.cluster.local", v.Name, v.Namespace),
		AdmissionReviewVersions: []string{
			"v1",
			"v1alpha1",
		},
		SideEffects:    &none,
		TimeoutSeconds: &seconds,
		ClientConfig: admissionv1.WebhookClientConfig{
			Service: &admissionv1.ServiceReference{
				Name:      fmt.Sprintf("%s", v.Name),
				Namespace: v.Namespace,
				Path:      &path,
				Port:      &p,
			},
			CABundle: v.Cert,
		},
	}
}

func (v *ValidatingServer) Rules() *ValidatingServer {
	v.Config.Webhooks[0].Rules = []admissionv1.RuleWithOperations{
		{
			Operations: []admissionv1.OperationType{
				"CREATE",
				"UPDATE",
			},
			Rule: admissionv1.Rule{
				APIGroups:   []string{"cmsnr.com"},
				APIVersions: []string{"v1alpha1"},
				Resources:   []string{"opapolicies"},
			},
		},
	}

	return v
}
