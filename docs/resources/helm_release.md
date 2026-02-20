---
page_title: "cpln_helm_release Resource - terraform-provider-cpln"
subcategory: "Helm"
description: |-
---

# cpln_helm_release (Resource)

Manages Helm chart deployments on Control Plane using the `cpln helm` command. This resource allows you to install, upgrade, and uninstall Helm charts that deploy Control Plane resources.

For more information about cpln helm, see the [Control Plane Helm Guide](https://docs.controlplane.com/guides/cpln-helm).

~> **Important** The `cpln` CLI and `helm` CLI must both be installed and available in the PATH for this resource to function.

~> **Important** The token or service account used must have permissions to create the resources defined in the helm chart, as well as `reveal` permission for secrets (to manage release state).

## Declaration

### Required

- **name** (String) The release name for this helm deployment.
- **chart** (String) Path to the chart. This can be a local path to a chart directory or packaged chart, or a chart name when used with `repository`.

### Optional

- **gvc** (String) The GVC (Global Virtual Cloud) to use for the helm deployment. Required only if the chart deploys GVC-scoped resources and the GVC is not defined within the chart manifests.
- **repository** (String) Chart repository URL where to locate the requested chart. Can be a Helm repository URL or an OCI registry URL.
- **version** (String) Specify a version constraint for the chart version to use. This can be a specific tag (e.g., 1.1.1) or a valid range (e.g., ^2.0.0). If not specified, the latest version is used.
- **values** (List of String) List of values in raw YAML to pass to the helm chart. Each entry is equivalent to a separate `--values/-f` flag. Values are merged in order, with later entries taking precedence.
- **set** (Map of String) Set values on the command line. Map of key-value pairs. Equivalent to using `--set` flag.
- **set_string** (Map of String) Set STRING values on the command line. Map of key-value pairs. Equivalent to using `--set-string` flag.
- **set_file** (Map of String) Set values from files specified via the command line. Map of key to file path. Equivalent to using `--set-file` flag.
- **wait** (Boolean) If set to true, will wait until all Workloads are in a ready state before marking the release as successful. Default: `false`.
- **timeout** (Number) The amount of seconds to wait for workloads to be ready before timing out. Only used when wait is true. Default: `300`.
- **dependency_update** (Boolean) Update dependencies if they are missing before installing the chart. Default: `false`.
- **description** (String) Add a custom description for the release.
- **verify** (Boolean) Verify the package before using it. Default: `false`.
- **max_history** (Number) Maximum number of revisions saved per release. Use 0 for no limit. Default: `10`. Only used on upgrade.
- **repository_username** (String) Chart repository username where to locate the requested chart.
- **repository_password** (String, Sensitive) Chart repository password where to locate the requested chart.
- **repository_ca_file** (String) Verify certificates of HTTPS-enabled servers using this CA bundle.
- **repository_cert_file** (String) Identify HTTPS client using this SSL certificate file.
- **repository_key_file** (String) Identify HTTPS client using this SSL key file.
- **insecure_skip_tls_verify** (Boolean) Skip TLS certificate checks for the chart download. Default: `false`.
- **render_subchart_notes** (Boolean) If set, render subchart notes along with the parent on install/upgrade. Default: `false`.
- **postrender** (Block) Post-renderer configuration:
  - **binary_path** (String, Required) The path to an executable to be used for post rendering.
  - **args** (List of String, Optional) Arguments to the post-renderer.

~> **Note** The `name` field requires resource replacement if changed.

## Outputs

The following attributes are exported:

- **id** (String) The unique identifier for this helm release (same as name).
- **status** (String) The current status of the helm release (e.g., "deployed", "failed").
- **revision** (Number) The current revision number of the helm release.
- **manifest** (String) The rendered manifest of the helm release.
- **resources** (Map of String) Rendered manifests keyed by resource identity (e.g., `workload/my-gvc/my-workload` for GVC-scoped resources, or `secret/my-secret` for org-scoped resources). Each value contains the rendered YAML manifest for that resource.

## Example Usage

### Basic Local Chart

```terraform
resource "cpln_gvc" "example" {
  name = "my-gvc"
}

resource "cpln_helm_release" "example" {
  name  = "my-release"
  gvc   = cpln_gvc.example.name
  chart = "./my-chart"

  values = [file("${path.module}/values.yaml")]
}
```

### Chart from Repository

```terraform
resource "cpln_gvc" "example" {
  name = "my-gvc"
}

resource "cpln_helm_release" "nginx" {
  name       = "my-nginx"
  gvc        = cpln_gvc.example.name
  chart      = "nginx"
  repository = "https://charts.bitnami.com/bitnami"
  version    = "15.0.0"

  values = [file("${path.module}/values.yaml")]

  set = {
    "image.tag" = "1.25.0"
  }
}
```

### Multiple Values Files

```terraform
resource "cpln_helm_release" "app" {
  name  = "my-app"
  gvc   = cpln_gvc.example.name
  chart = "./my-chart"

  values = [
    file("${path.module}/values/base.yaml"),
    file("${path.module}/values/production.yaml"),
  ]
}
```

### Chart from OCI Registry

```terraform
resource "cpln_gvc" "example" {
  name = "my-gvc"
}

resource "cpln_helm_release" "app" {
  name       = "my-app"
  gvc        = cpln_gvc.example.name
  chart      = "my-chart"
  repository = "oci://registry.example.com/charts"
  version    = "1.0.0"

  repository_username = "registry-user"
  repository_password = var.registry_password
}
```

### With Wait and Timeout

```terraform
resource "cpln_gvc" "example" {
  name = "my-gvc"
}

resource "cpln_helm_release" "critical_app" {
  name  = "critical-app"
  gvc   = cpln_gvc.example.name
  chart = "./critical-chart"

  wait    = true
  timeout = 600

  values = [file("${path.module}/values.yaml")]
}
```

### Using Set Values

```terraform
resource "cpln_gvc" "example" {
  name = "my-gvc"
}

resource "cpln_helm_release" "configurable" {
  name  = "configurable-app"
  gvc   = cpln_gvc.example.name
  chart = "./my-chart"

  values = [file("${path.module}/values.yaml")]

  set = {
    "app.environment"       = "production"
    "app.logLevel"          = "info"
    "resources.limits.cpu"  = "1000m"
  }

  set_string = {
    "app.version" = "1.2.3"
  }
}
```

### With Post-Renderer

```terraform
resource "cpln_helm_release" "app" {
  name  = "my-app"
  gvc   = cpln_gvc.example.name
  chart = "./my-chart"

  values = [file("${path.module}/values.yaml")]

  postrender = {
    binary_path = "/usr/local/bin/kustomize"
    args        = ["build", "./overlays/production"]
  }
}
```

### Accessing Created Resources

```terraform
resource "cpln_gvc" "example" {
  name = "my-gvc"
}

resource "cpln_helm_release" "app" {
  name  = "my-app"
  gvc   = cpln_gvc.example.name
  chart = "./my-chart"

  values = [file("${path.module}/values.yaml")]
}

output "helm_resources" {
  value = cpln_helm_release.app.resources
}

output "helm_manifest" {
  value = cpln_helm_release.app.manifest
}
```

## Import Syntax

To import an existing helm release into Terraform state:

```terraform
terraform import cpln_helm_release.RESOURCE_NAME RELEASE_NAME
```

-> 1. Substitute RESOURCE_NAME with the same string that is defined in the HCL file.<br/>2. Substitute RELEASE_NAME with the name of the installed helm release.

~> **Note** Import only recovers the computed attributes (`id`, `status`, `revision`, `manifest`, `resources`). Configuration attributes such as `chart`, `gvc`, `values`, `set`, `repository`, and all other options cannot be recovered from the release state. You must define the correct configuration in your Terraform file before or after importing. On the next `terraform apply`, Terraform will detect the difference and trigger a helm upgrade to reconcile the state.

## Lifecycle Considerations

- **Install**: On resource creation, `cpln helm install` is executed.
- **Upgrade**: On resource update, `cpln helm upgrade` is executed.
- **Uninstall**: On resource deletion, `cpln helm uninstall` is executed.
- **Rollback**: To rollback to a previous revision, you can use `cpln helm rollback` outside of Terraform, or manage the values/chart version in your Terraform configuration.

## Automatically Injected Values

The `cpln helm` command automatically injects the following values into all charts:

- `cpln.org`: Current organization name
- `cpln.gvc`: Current GVC name

Avoid defining top-level `cpln` keys in your values files, as they will be overwritten.
