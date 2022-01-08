package server

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestCheckInject(t *testing.T) {
	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"cmsnr.com/inject": "enabled",
			},
		},
	}

	ok := checkInject(&pod)
	if !ok {
		t.Error("error validating injection")
	}

}
