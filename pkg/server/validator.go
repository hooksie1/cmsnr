package server

import (
	"context"
	"github.com/open-policy-agent/opa/rego"
	api "gitlab.com/hooksie1/cmsnr/api/v1alpha1"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type Validator struct {
	Client  client.Client
	decoder *admission.Decoder
}

func (v *Validator) Handle(ctx context.Context, req admission.Request) admission.Response {
	opa := &api.OpaPolicy{}

	if err := v.decoder.Decode(req, opa); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	rego := rego.New(
		rego.Query("data.test"),
		rego.Module("example.rego",
			opa.Spec.Policy,
		))

	_, err := rego.Compile(ctx)
	if err != nil {
		return admission.Denied(err.Error())
	}

	return admission.Allowed("policy is valid")
}

func (v *Validator) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}
