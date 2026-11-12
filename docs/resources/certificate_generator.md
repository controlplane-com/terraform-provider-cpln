---
page_title: "cpln_certificate_generator Resource - terraform-provider-cpln"
subcategory: "Certificate Generator"
description: |-
---

# cpln_certificate_generator (Resource)

Manages an org's [Certificate Generator](https://docs.controlplane.com/reference/certificategenerator).

A certificate generator issues [domain](https://docs.controlplane.com/reference/domain) certificates from an ACME certificate authority of your choice instead of the built-in Let's Encrypt issuer. A domain uses it by setting `spec.certificate_generator_link`. Many domains can share one generator.

## Declaration

### Required

- **name** (String) Name of the Certificate Generator.
- **acme** (Attributes) Issues domain certificates from an ACME certificate authority. ([see below](#nestedblock--acme))

### Optional

- **description** (String) Description of the Certificate Generator.
- **tags** (Map of String) Key-value map of resource tags.

<a id="nestedblock--acme"></a>

### `acme`

Required:

- **server** (String) ACME directory URL of the CA, e.g. `https://acme-v02.api.letsencrypt.org/directory`. Must use `https`. Changing it recreates the certificate generator.
- **challenge** (Attributes) How ACME challenges are solved. Exactly one of `http01` or `dns01` must be set. Changing the challenge type recreates the certificate generator. ([see below](#nestedblock--acme--challenge))

Optional:

- **email** (String) Contact email registered with the ACME account.
- **external_account_binding** (Attributes) Binds the ACME account to an existing account at the CA. Required by most commercial CAs. ([see below](#nestedblock--acme--external_account_binding))

<a id="nestedblock--acme--external_account_binding"></a>

### `acme.external_account_binding`

Required:

- **key_id** (String) Key identifier issued by the CA.
- **secret_link** (String) Link to a secret of type `opaque` whose payload is the base64url encoded HMAC key exactly as issued by the CA (encoding `plain`).

<a id="nestedblock--acme--challenge"></a>

### `acme.challenge`

Optional:

- **http01** (Attributes) Solve ACME challenges over HTTP on the domain itself. Only domain-name certificates can be issued; wildcards require `dns01`. Set it to `{}`.
- **dns01** (Attributes) Solve ACME challenges with TXT records, either delegated to Control Plane (`cpln`) or written into DNS zones the org controls. Required for wildcard certificates. ([see below](#nestedblock--acme--challenge--dns01))

<a id="nestedblock--acme--challenge--dns01"></a>

### `acme.challenge.dns01`

Optional:

- **provider** (String) DNS provider that hosts the zones of the domains using this generator. `cpln`: domains CNAME `_acme-challenge` to Control Plane, no credentials needed. Valid values: `cpln`, `route53`, `cloudDns`, `azureDns`, `cloudflare`. Default: `cpln`.
- **route53** (Attributes) Route53 settings. Required when `provider` is `route53`. ([see below](#nestedblock--acme--challenge--dns01--route53))
- **cloud_dns** (Attributes) Google Cloud DNS settings. Required when `provider` is `cloudDns`. ([see below](#nestedblock--acme--challenge--dns01--cloud_dns))
- **azure_dns** (Attributes) Azure DNS settings. Required when `provider` is `azureDns`. ([see below](#nestedblock--acme--challenge--dns01--azure_dns))
- **cloudflare** (Attributes) Cloudflare settings. Required when `provider` is `cloudflare`. ([see below](#nestedblock--acme--challenge--dns01--cloudflare))

<a id="nestedblock--acme--challenge--dns01--route53"></a>

### `acme.challenge.dns01.route53`

Required:

- **region** (String) AWS region of the Route53 API to use, e.g. `us-east-1`.
- **secret_link** (String) Link to a secret of type `aws`. The access key must be allowed to change TXT records in the hosted zone.

Optional:

- **hosted_zone_id** (String) Restrict record updates to this hosted zone; discovered from the domain when omitted.

<a id="nestedblock--acme--challenge--dns01--cloud_dns"></a>

### `acme.challenge.dns01.cloud_dns`

Required:

- **project** (String) GCP project that owns the Cloud DNS zone.
- **secret_link** (String) Link to a secret of type `gcp` holding a service account key with the DNS Administrator role.

<a id="nestedblock--acme--challenge--dns01--azure_dns"></a>

### `acme.challenge.dns01.azure_dns`

Required:

- **resource_group_name** (String) Resource group that contains the DNS zone.
- **secret_link** (String) Link to a secret of type `azure-sdk`. Its service principal must be a DNS Zone Contributor on the zone.

Optional:

- **hosted_zone_name** (String) Restrict record updates to this DNS zone; discovered from the domain when omitted.

<a id="nestedblock--acme--challenge--dns01--cloudflare"></a>

### `acme.challenge.dns01.cloudflare`

Required:

- **secret_link** (String) Link to a secret of type `opaque` whose payload is a Cloudflare API token with Zone:DNS:Edit permission.

## Outputs

The following attributes are exported:

- **cpln_id** (String) The ID, in GUID format, of the Certificate Generator.
- **self_link** (String) Full link to this resource. Can be referenced by other resources.
- **status** (Block List, Max: 1) ([see below](#nestedblock--status)).

<a id="nestedblock--status"></a>

### `status`

The status of the Certificate Generator.

Read-Only:

- **ready** (Boolean) True once the generator is registered with the CA and can issue certificates.
- **message** (String) Reason the generator is not ready.
- **last_updated** (String) Timestamp of the last status update.

## Example Usage - HTTP-01

```terraform
resource "cpln_certificate_generator" "example_http01" {
  name        = "example-http01"
  description = "Certificate generator that solves HTTP-01 challenges"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  acme = {
    server = "https://acme-v02.api.letsencrypt.org/directory"
    email  = "ops@example.com"

    challenge = {
      http01 = {}
    }
  }
}
```

## Example Usage - DNS-01 Delegated to Control Plane

```terraform
resource "cpln_secret" "acme_eab_hmac" {
  name = "acme-eab-hmac"

  opaque {
    payload  = "HMAC_KEY_FROM_CA"
    encoding = "plain"
  }
}

resource "cpln_certificate_generator" "example_dns01_cpln" {
  name        = "example-dns01-cpln"
  description = "Certificate generator for a commercial CA that delegates DNS-01 challenges to Control Plane"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  acme = {
    server = "https://acme.example-ca.com/directory"
    email  = "ops@example.com"

    external_account_binding = {
      key_id      = "KEY_ID_FROM_CA"
      secret_link = cpln_secret.acme_eab_hmac.self_link
    }

    challenge = {
      dns01 = {
        provider = "cpln"
      }
    }
  }
}
```

## Example Usage - DNS-01 in Route53

```terraform
resource "cpln_secret" "route53_dns" {
  name = "route53-dns"

  aws {
    access_key = "AKIAIOSFODNN7EXAMPLE"
    secret_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
  }
}

resource "cpln_certificate_generator" "example_dns01_route53" {
  name        = "example-dns01-route53"
  description = "Certificate generator that writes DNS-01 challenge records into Route53"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  acme = {
    server = "https://acme-v02.api.letsencrypt.org/directory"
    email  = "ops@example.com"

    challenge = {
      dns01 = {
        provider = "route53"

        route53 = {
          region         = "us-east-1"
          hosted_zone_id = "Z0123456789ABCDEFGHIJ"
          secret_link    = cpln_secret.route53_dns.self_link
        }
      }
    }
  }
}

resource "cpln_domain" "example" {
  name = "app.example.com"

  spec {
    dns_mode                   = "cname"
    gvc_link                   = "/org/myorg/gvc/mygvc"
    certificate_generator_link = cpln_certificate_generator.example_dns01_route53.self_link

    ports {
      tls {}
    }
  }
}
```

## Import Syntax

The `terraform import` command is used to bring existing infrastructure resources, created outside of Terraform, into the Terraform state file, enabling their management through Terraform going forward.

To update a statefile with an existing certificate generator resource, execute the following import command:

```terraform
terraform import cpln_certificate_generator.RESOURCE_NAME CERTIFICATE_GENERATOR_NAME
```

-> 1. Substitute RESOURCE_NAME with the same string that is defined in the HCL file.<br/>2. Substitute CERTIFICATE_GENERATOR_NAME with the corresponding certificate generator defined in the resource.
