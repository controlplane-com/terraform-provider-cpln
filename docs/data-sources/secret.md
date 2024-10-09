---
page_title: "cpln_secret Data Source - terraform-provider-cpln"
subcategory: "Secret"
description: |-
---

# cpln_secret (Data Source)

Use this data source to access information about a [Secret](https://docs.controlplane.com/reference/secret) within Control Plane.

## Required

- **name** (String) Name of the secret.

## Outputs

The following attributes are exported:

- **cpln_id** (String) The ID, in GUID format, of the secret.
- **name** (String) Name of the secret.
- **description** (String) Description of the secret.
- **tags** (Map of String) Key-value map of resource tags.
- **self_link** (String) Full link to this resource. Can be referenced by other resources.
- **secret_link** (String) Output used when linking a secret to an environment variable or volume, in the format: `cpln://secret/SECRET_NAME`.
- **dictionary_as_envs** (Map of String) If a dictionary secret is defined, this output will be a key-value map in the following format: `key = cpln://secret/SECRET_NAME.key`.
- **aws** (Block List, Max: 1) ([see below](#nestedblock--aws)) [Reference Page](https://docs.controlplane.com/reference/secret#aws).
- **azure_connector** (Block List, Max: 1) ([see below](#nestedblock--azure_connector)) [Reference Page](https://docs.controlplane.com/reference/secret#azure-connector).
- **azure_sdk** (String, Sensitive) JSON string containing the Docker secret. [Reference Page](https://docs.controlplane.com/reference/secret#azure).
- **dictionary** (Map of String) List of unique key-value pairs. [Reference Page](https://docs.controlplane.com/reference/secret#dictionary).
- **docker** (String, Sensitive) JSON string containing the Docker secret. [Reference Page](https://docs.controlplane.com/reference/secret#docker).
- **ecr** (Block List, Max: 1) ([see below](#nestedblock--ecr)).
- **gcp** (String, Sensitive) JSON string containing the GCP secret. [Reference Page](https://docs.controlplane.com/reference/secret#gcp)
- **keypair** (Block List, Max: 1) ([see below](#nestedblock--keypair)) [Reference Page](https://docs.controlplane.com/reference/secret#keypair).
- **nats_account** (Block List, Max: 1) ([see below](#nestedblock--nats-account)) [Reference Page](https://docs.controlplane.com/reference/secret#nats-account).
- **opaque** (Block List, Max: 1) ([see below](#nestedblock--opaque)) [Reference Page](https://docs.controlplane.com/reference/secret#opaque).
- **tls** (Block List, Max: 1) ([see below](#nestedblock--tls)) [Reference Page](https://docs.controlplane.com/reference/secret#tls).
- **userpass** (Block List, Max: 1) ([see below](#nestedblock--userpass)) [Reference Page](https://docs.controlplane.com/reference/secret#username).

<a id="nestedblock--aws"></a>

### `aws`

Optional:

- **access_key** (String, Sensitive) Access Key provided by AWS.
- **role_arn** (String) Role ARN provided by AWS.
- **secret_key** (String, Sensitive) Secret Key provided by AWS.
- **external_id** (String) AWS IAM Role External ID.

<a id="nestedblock--azure_connector"></a>

### `azure_connector`

Optional:

- **code** (String, Sensitive) Code/Key to authenticate to deployment URL.
- **url** (String, Sensitive) Deployment URL.

<a id="nestedblock--ecr"></a>

### `ecr`

[Reference Page](https://docs.controlplane.com/reference/secret#ecr)

Optional:

- **access_key** (String) Access Key provided by AWS.
- **repos** (Set of String) List of ECR repositories.
- **role_arn** (String) Role ARN provided by AWS.
- **secret_key** (String, Sensitive) Secret Key provided by AWS.
- **external_id** (String) AWS IAM Role External ID. Used when setting up cross-account access to your ECR repositories.

<a id="nestedblock--keypair"></a>

### `keypair`

Optional:

- **passphrase** (String, Sensitive) Passphrase for private key.
- **public_key** (String) Public Key.
- **secret_key** (String, Sensitive) Secret/Private Key.

<a id="nestedblock--nats-account"></a>

### `nats_account`

Required:

- **account_id** (String) Account ID.
- **private_key** (String) Private Key.

<a id="nestedblock--opaque"></a>

### `opaque`

Optional:

- **encoding** (String) Available encodings: `plain`, `base64`. Default: `plain`.
- **payload** (String, Sensitive) Plain text or base64 encoded string. Use `encoding` attribute to specify encoding.

<a id="nestedblock--tls"></a>

### `tls`

Optional:

- **cert** (String) Public Certificate.
- **chain** (String) Chain Certificate.
- **key** (String, Sensitive) Private Certificate.

<a id="nestedblock--userpass"></a>

### `userpass`

Optional:

- **encoding** (String) Available encodings: `plain`, `base64`. Default: `plain`.
- **password** (String, Sensitive) Password.
- **username** (String) Username.

## Example Usage

```terraform
data "cpln_secret" "example" {
  name = "example-secret-opaque"
}

output "example-secret-payload" {
  value     = data.cpln_secret.example.opaque.payload
  sensitive = true
}
```
