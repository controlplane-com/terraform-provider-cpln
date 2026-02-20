---
page_title: "cpln_helm_template Data Source - terraform-provider-cpln"
subcategory: "Helm"
description: |-
---

# cpln_helm_template (Data Source)

Renders Helm chart templates using the `cpln helm template` command without installing. Useful for previewing rendered manifests or feeding them into other resources.

For more information about cpln helm, see the [Control Plane Helm Guide](https://docs.controlplane.com/guides/cpln-helm).

~> **Important** The `cpln` CLI and `helm` CLI must both be installed and available in the PATH for this data source to function.

## Declaration

### Required

- **name** (String) The release name to use for rendering the templates.
- **chart** (String) Path to the chart. This can be a local path to a chart directory or packaged chart, or a chart name when used with `repository`.

### Optional

- **gvc** (String) The GVC (Global Virtual Cloud) context for rendering the helm chart templates. Required only if the chart contains GVC-scoped resources and the GVC is not defined within the chart manifests.
- **repository** (String) Chart repository URL where to locate the requested chart. Can be a Helm repository URL or an OCI registry URL.
- **version** (String) Specify a version constraint for the chart version to use. This can be a specific tag (e.g., 1.1.1) or a valid range (e.g., ^2.0.0). If not specified, the latest version is used.
- **values** (List of String) List of values in raw YAML to pass to the helm chart. Each entry is equivalent to a separate `--values/-f` flag. Values are merged in order, with later entries taking precedence.
- **set** (Map of String) Set values on the command line. Map of key-value pairs. Equivalent to using `--set` flag.
- **set_string** (Map of String) Set STRING values on the command line. Map of key-value pairs. Equivalent to using `--set-string` flag.
- **set_file** (Map of String) Set values from files specified via the command line. Map of key to file path. Equivalent to using `--set-file` flag.
- **dependency_update** (Boolean) Update dependencies if they are missing before rendering the chart.
- **description** (String) Add a custom description.
- **verify** (Boolean) Verify the package before using it.
- **repository_username** (String) Chart repository username where to locate the requested chart.
- **repository_password** (String, Sensitive) Chart repository password where to locate the requested chart.
- **repository_ca_file** (String) Verify certificates of HTTPS-enabled servers using this CA bundle.
- **repository_cert_file** (String) Identify HTTPS client using this SSL certificate file.
- **repository_key_file** (String) Identify HTTPS client using this SSL key file.
- **insecure_skip_tls_verify** (Boolean) Skip TLS certificate checks for the chart download.
- **render_subchart_notes** (Boolean) If set, render subchart notes along with the parent.
- **postrender** (Block) Post-renderer configuration:
  - **binary_path** (String, Required) The path to an executable to be used for post rendering.
  - **args** (List of String, Optional) Arguments to the post-renderer.

## Outputs

The following attributes are exported:

- **id** (String) The unique identifier for this data source (same as name).
- **manifest** (String) The rendered manifest output from helm template.

## Example Usage

### Basic Template Rendering

```terraform
data "cpln_helm_template" "example" {
  name  = "my-release"
  gvc   = "my-gvc"
  chart = "./my-chart"

  values = [file("${path.module}/values.yaml")]
}

output "rendered_manifest" {
  value = data.cpln_helm_template.example.manifest
}
```

### Template from Repository

```terraform
data "cpln_helm_template" "nginx" {
  name       = "my-nginx"
  gvc        = "my-gvc"
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
data "cpln_helm_template" "app" {
  name  = "my-app"
  gvc   = "my-gvc"
  chart = "./my-chart"

  values = [
    file("${path.module}/values/base.yaml"),
    file("${path.module}/values/production.yaml"),
  ]
}
```
