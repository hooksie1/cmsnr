package server

//
//import (
//	"context"
//	api "gitlab.com/hooksie1/cmsnr/api/v1alpha1"
//	"net/http"
//	"sigs.k8s.io/controller-runtime/pkg/client"
//	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
//)
//
//type Validator struct {
//	Client  client.Client
//	decoder *admission.Decoder
//}
//
//func (v *Validator) Handle(ctx context.Context, req admission.Request) admission.Response {
//	opa := &api.OpaPolicy{}
//
//	if err := v.decoder.Decode(req, opa); err != nil {
//		return admission.Errored(http.StatusBadRequest, err)
//	}
//
//}
