package deployment

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewDeployment(name, namespace string, port int) *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": name,
				},
			},
			Template: getTemplate(name, port),
		},
	}
}

func getTemplate(name string, port int) corev1.PodTemplateSpec {
	return corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"app": name,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Image:           "hooksie1/cmsnr",
					ImagePullPolicy: "Always",
					Name:            name,
					Args:            []string{"server", "start", "mutating"},
					Ports: []corev1.ContainerPort{
						{
							Name:          "https",
							ContainerPort: int32(port),
						},
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "webhook-certs",
							MountPath: "/var/lib/cmsnr",
						},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: "webhook-certs",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "cmsnr-secret",
						},
					},
				},
			},
		},
	}
}
