package server

import (
	"context"
	"fmt"
	"net/http"

	api "github.com/hooksie1/cmsnr/api/v1alpha1"
	"github.com/open-policy-agent/opa/ast"
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

	_, err := ast.ParseModule("example.rego", opa.Spec.Policy)
	if err != nil {
		e := fmt.Sprintf("eval error: %v", err.Error())
		return admission.Denied(e)

	}

	return admission.Allowed("policy is valid")
}

func (v *Validator) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}
