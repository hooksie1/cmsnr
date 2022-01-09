package server

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type Config struct {
	Containers []corev1.Container `yaml:"containers"`
}

type SidecarInjector struct {
	Client  client.Client
	decoder *admission.Decoder
}

func getContainers(depName string) []corev1.Container {
	return []corev1.Container{
		{
			Name:            "opa",
			Image:           "openpolicyagent/opa:latest-static",
			ImagePullPolicy: corev1.PullPolicy("IfNotPresent"),
			Args:            []string{"run", "--server"},
		},
		{
			Name:            "cmsnr-client",
			Image:           "hooksie1/cmsnr:latest",
			ImagePullPolicy: corev1.PullPolicy("IfNotPresent"),
			Args:            []string{"opa", "watch", fmt.Sprintf("-d=%s", depName)},
		},
	}
}

func checkInject(pod *corev1.Pod) bool {
	if pod.Labels["cmsnr.com/inject"] == "enabled" {
		return true
	}

	return false
}

func (s *SidecarInjector) Handle(ctx context.Context, r admission.Request) admission.Response {
	pod := &corev1.Pod{}

	err := s.decoder.Decode(r, pod)
	if err != nil {
		log.Errorf("error decoding: %s", err)
		return admission.Errored(http.StatusBadRequest, err)
	}

	if pod.Annotations == nil {
		pod.Annotations = map[string]string{}
	}

	if checkInject(pod) {
		log.Infof("Injecting sidecar for %s", pod.Name)
		pod.Spec.Containers = append(pod.Spec.Containers, getContainers(pod.Annotations["cmsnr.com/deploymentName"])...)
		if pod.Spec.ServiceAccountName == "default" {
			log.Info("no service account defined, adding cmsnr account")
			pod.Spec.ServiceAccountName = "cmsnr"
		}
		pod.Annotations["cmsnr.com/injected"] = "true"
	}

	marshaled, err := json.Marshal(pod)
	if err != nil {
		log.Errorf("error marshaling pod: %s", err)
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(r.Object.Raw, marshaled)
}

func (s *SidecarInjector) InjectDecoder(d *admission.Decoder) error {
	s.decoder = d
	return nil
}
