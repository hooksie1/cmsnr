package deployment

func NewCRD() string {
	return `
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: opapolicies.cmsnr.com
spec:
  group: "cmsnr.com"
  names:
    kind: "OpaPolicy"
    listKind: "OpaPolicyList"
    singular: "opapolicy"
    plural: "opapolicies"
  scope: "Namespaced"
  versions: 
    - name: v1alpha1
      schema:
        openAPIV3Schema:
          description: Schema for the OpaPolicy
          properties:
            apiVersion:
              description: 'defines the versioned schema'
              type: string
            kind:
              description: 'Kind is a value representing the object'
              type: string
            metadata:
              type: object
            spec:
              description: OpaPolicySpec defines the desired state of the OpaPolicy
              properties:
                deploymentName:
                  description: The name of the matching deployment
                  type: string
                policyName:
                  description: The name of the policy as it should be in OPA
                  type: string
                policy:
                  description: The rego policy contents
                  type: string
              type: object
          type: object
      served: true
      storage: true`
}
