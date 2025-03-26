---
page_title: "cpln_mk8s_kubeconfig Resource - terraform-provider-cpln"
subcategory: "Mk8s"
description: |-
---

# cpln_mk8s_kubeconfig (Resource)

Obtain the Cluster’s Kubeconfig.

## Declaration

### Required

- **name** (String) Name of the Mk8s.

~> **Note** Only one of the below can be included in the resource.

- **profile** (String) The name of the cpln profile used to generate the kubeconfig file for authenticating with your Kubernetes cluster.
- **service_account** (String) The name of an existing service account for which a key will be generated, enabling kubeconfig-based authentication with your Kubernetes cluster.

## Outputs

The following attributes are exported:

- **kubeconfig** (String) The Kubeconfig in YAML format.

## Example Usage - Profile

```terraform
resource "cpln_mk8s_kubeconfig" "new" {
  name    = "generic-cluster"
  profile = "default"
}

output "generic-cluster-kubeconfig" {
  value = cpln_mk8s_kubeconfig.new.kubeconfig
  sensitive = true
}
```

## Example Usage - Service Account

```terraform
resource "cpln_mk8s_kubeconfig" "new" {
  name            = "generic-cluster"
  service_account = "devops-sa"
}

output "generic-cluster-kubeconfig" {
  value = cpln_mk8s_kubeconfig.new.kubeconfig
  sensitive = true
}
```
