---
page_title: "cpln_org_logging Resource - terraform-provider-cpln"
subcategory: "Org Logging"
description: |-
---

# cpln_org_logging (Resource)

Manages an [Org's external logging configuration](https://docs.controlplane.com/external-logging/overview).

## Declaration

### Required

You can define up to **four** logging blocks:

- **s3_logging** (Block List, Max: 1) ([see below](#nestedblock--s3_logging)).
- **coralogix_logging** (Block List, Max: 1) ([see below](#nestedblock--coralogix_logging)).
- **datadog_logging** (Block List, Max: 1) ([see below](#nestedblock--datadog_logging)).
- **elastic_logging** (Block List, Max: 1) ([see below](#nestedblock--elastic_logging)).
- **logzio_logging** (Block List, Max: 1) ([see below](#nestedblock--logzio_logging)).

<a id="nestedblock--s3_logging"></a>

### `s3_logging`

[Documentation Reference](https://docs.controlplane.com/external-logging/s3)

Required:

- **bucket** (String) Name of S3 bucket.
- **region** (String) AWS region where bucket is located.
- **prefix** (String) Bucket path prefix. Default: "/".
- **credentials** (String) Full link to referenced AWS Secret.

<a id="nestedblock--coralogix_logging"></a>

### `coralogix_logging`

[Documentation Reference](https://docs.controlplane.com/external-logging/coralogix)

Required:

- **cluster** (String) Coralogix cluster URI.
- **credentials** (String) Full link to referenced Opaque Secret.
- **app** (String) App name to be displayed in Coralogix dashboard.
- **subsystem** (String) Subsystem name to be displayed in Coralogix dashboard.

~> **Note** Valid clusters: `coralogix.com`, `coralogix.us`, `app.coralogix.in`, `app.eu2.coralogix.com`, `app.coralogixsg.com`.

~> **Note** Supported variables for App and Subsystem are: freeformed or `{org}`, `{gvc}`, `{workload}`, `{location}`.

<a id="nestedblock--datadog_logging"></a>

### `datadog_logging`

[Documentation Reference](https://docs.controlplane.com/external-logging/datadog)

Required:

- **host** (String) Datadog host URI.
- **credentials** (String) Full link to referenced Opaque Secret.

~> **Note** Valid Hosts: `http-intake.logs.datadoghq.com`, `http-intake.logs.us3.datadoghq.com`, `http-intake.logs.us5.datadoghq.com`, `http-intake.logs.datadoghq.eu`.

<a id="nestedblock--elastic_logging"></a>

### `elastic_logging`

<!-- [Documentation Reference](https://docs.controlplane.com/external-logging/elastic-aws)
[Documentation Reference](https://docs.controlplane.com/external-logging/elastic-co) -->

Required:
At least one of the following logging blocks are required:

- **aws** (Block List, Max: 1) ([see below](#nestedblock--elastic_logging--aws)).
- **elastic_cloud** (Block List, Max: 1) ([see below](#nestedblock--elastic_logging--elastic_cloud)).
- **generic** (Block List, Max: 1) ([see below](#nestedblock--elastic_logging--generic)).

<a id="nestedblock--elastic_logging--aws"></a>

Required:

- **host** (String) A valid AWS ElasticSearch hostname (must end with es.amazonaws.com).
- **port** (Number) Port. Default: 443
- **index** (String) Logging Index.
- **type** (String) Logging Type.
- **credentials** (String) Full Link to a secret of type `aws`.
- **region** (String) Valid AWS region.

<a id="nestedblock--elastic_logging--elastic_cloud"></a>

Required:

- **index** (String) Logging Index.
- **type** (String) Logging Type.
- **credentials** (String) Full Link to a secret of type `userpass`.
- **cloud_id** (String) [Cloud ID](https://www.elastic.co/guide/en/cloud/current/ec-cloud-id.html)

<a id="nestedblock--elastic_logging--generic"></a>

Required:

- **host** (String) A valid Elastic Search provider hostname.
- **port** (Number) Port. Default: 443
- **path** (String) Logging path.
- **index** (String) Logging Index.
- **type** (String) Logging Type.
- **credentials** (String) Full Link to a secret of type `userpass`.

<a id="nestedblock--logzio_logging"></a>

### `logzio_logging`

[Documentation Reference](https://docs.controlplane.com/external-logging/logz-io)

Required:

- **listener_host** (String) Logzio listener host URI.
- **credentials** (String) Full link to referenced Opaque Secret.

~> **Note** Valid listener hosts: `listener.logz.io`, `listener-nl.logz.io`

## Outputs

The following attributes are exported:

- **cpln_id** (String) ID, in GUID format, of the Org.
- **name** (String) The name of Org.
- **description** (String) The description of Org.
- **tags** (Map of String) Key-value map of the Org's tags.

## Example Usage

### S3

```terraform
resource "cpln_secret" "aws" {

  name        = "aws-random-tbd"
  description = "aws description aws-random-tbd"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    secret_type         = "aws"
  }

  aws {
    secret_key = "AKIAIOSFODNN7EXAMPLE"
    access_key = "AKIAwJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
    role_arn   = "arn:awskey"
  }
}

resource "cpln_org_logging" "new" {

  s3_logging {

    bucket = "test-bucket"
    region = "us-east1"
    prefix = "/"

    // AWS Secret Only
    credentials = cpln_secret.aws.self_link
  }
}
```

### Coralogix

```terraform
resource "cpln_secret" "opaque" {

  name        = "opaque-random-coralogix-tbd"
  description = "opaque description opaque-random-tbd"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    secret_type         = "opaque"
  }

  opaque {
    payload  = "opaque_secret_payload"
    encoding = "plain"
  }
}

resource "cpln_org_logging" "new" {

  coralogix_logging {

    // Valid clusters
    // coralogix.com, coralogix.us, app.coralogix.in, app.eu2.coralogix.com, app.coralogixsg.com
    cluster = "coralogix.com"

    // Opaque Secret Only
    credentials = cpln_secret.opaque.self_link

    // Supported variables for App and Subsystem are:
    // {org}, {gvc}, {workload}, {location}
    app       = "{workload}"
    subsystem = "{org}"
  }
}
```

### Datadog

```terraform
resource "cpln_secret" "opaque" {

  name        = "opaque-random-datadog-tbd"
  description = "opaque description"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    secret_type         = "opaque"
  }

  opaque {
    payload  = "opaque_secret_payload"
    encoding = "plain"
  }
}

resource "cpln_org_logging" "new" {

  datadog_logging {

    // Valid Hosts
    // http-intake.logs.datadoghq.com, http-intake.logs.us3.datadoghq.com,
    // http-intake.logs.us5.datadoghq.com, http-intake.logs.datadoghq.eu
    host = "http-intake.logs.datadoghq.com"

    // Opaque Secret Only
    credentials = cpln_secret.opaque.self_link
  }
}
```

### Logz.io

```terraform
resource "cpln_secret" "opaque" {

  name        = "opaque-random-datadog-tbd"
  description = "opaque description"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    secret_type         = "opaque"
  }

  opaque {
    payload  = "opaque_secret_payload"
    encoding = "plain"
  }
}

resource "cpln_org_logging" "new" {

  logzio_logging {

    // Valid Listener Hosts
    // listener.logz.io, listener-nl.logz.io
    listener_host = "listener.logz.io"

    // Opaque Secret Only
    credentials = cpln_secret.opaque.self_link
  }
}
```

### Elastic - AWS

```terraform
resource "cpln_secret" "aws" {

  name        = "aws-random-elastic-logging-aws"
  description = "aws description"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    secret_type         = "aws"
  }

  aws {
    secret_key = "AKIAIOSFODNN7EXAMPLE"
    access_key = "AKIAwJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
    role_arn   = "arn:awskey"
  }
}

resource "cpln_org_logging" "new" {

  elastic_logging {
    aws {
      host        = "es.amazonaws.com"
      port        = 8080
      index       = "my-index"
      type        = "my-type"

      // AWS Secret Only
      credentials = cpln_secret.aws.self_link
      region      = "us-east-1"
    }
  }
}
```

### Elastic - Elastic Cloud

```terraform
resource "cpln_secret" "userpass" {
  name = "example"

  userpass {
    username = "cpln_username"
    password = "cpln_password"
    encoding = "plain"
  }
}

resource "cpln_org_logging" "new" {

  elastic_logging {
    elastic_cloud {
      index       = "my-index"
      type        = "my-type"

      // UserPass Secret Only
      credentials = cpln_secret.userpass.self_link
      cloud_id    = "my-cloud-id"
    }
  }
}
```

### Elastic - Generic

```terraform
resource "cpln_secret" "userpass" {
  name = "example"

  userpass {
    username = "cpln_username"
    password = "cpln_password"
    encoding = "plain"
  }
}

resource "cpln_org_logging" "new" {

  elastic_logging {
    generic {
      host  = "example.com"
      port  = 9200
      path  = "/var/log/elasticsearch/"
      index = "my-index"
      type  = "my-type"

      // UserPass Secret Only
      credentials = cpln_secret.userpass.self_link
    }
  }
}
```

### Use Of Three Unique Loggings


```terraform
resource "cpln_secret" "aws" {

  name        = "aws-random-tbd"
  description = "aws description aws-random-tbd"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    secret_type         = "aws"
  }

  aws {
    secret_key = "AKIAIOSFODNN7EXAMPLE"
    access_key = "AKIAwJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
    role_arn   = "arn:awskey"
  }
}

resource "cpln_secret" "opaque-coralogix" {

  name        = "opaque-random-coralogix-tbd"
  description = "opaque description opaque-random-tbd"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    secret_type         = "opaque"
  }

  opaque {
    payload  = "opaque_secret_payload"
    encoding = "plain"
  }
}

resource "cpln_secret" "opaque-datadog" {

  name        = "opaque-random-datadog-tbd"
  description = "opaque description"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    secret_type         = "opaque"
  }

  opaque {
    payload  = "opaque_secret_payload"
    encoding = "plain"
  }
}

resource "cpln_org_logging" "tf-logging" {
  s3_logging {

    bucket = "test-bucket"
    region = "us-east1"
    prefix = "/"

    // AWS Secret Only
    credentials = cpln_secret.aws.self_link
  }

  coralogix_logging {

    // Valid clusters
    // coralogix.com, coralogix.us, app.coralogix.in, app.eu2.coralogix.com, app.coralogixsg.com
    cluster = "coralogix.com"

    // Opaque Secret Only
    credentials = cpln_secret.opaque.self_link

    // Supported variables for App and Subsystem are:
    // {org}, {gvc}, {workload}, {location}
    app       = "{workload}"
    subsystem = "{org}"
  }

  datadog_logging {

    // Valid Hosts
    // http-intake.logs.datadoghq.com, http-intake.logs.us3.datadoghq.com,
    // http-intake.logs.us5.datadoghq.com, http-intake.logs.datadoghq.eu
    host = "http-intake.logs.datadoghq.com"

    // Opaque Secret Only
    credentials = cpln_secret.opaque.self_link
  }
}
```

## Import Syntax

The `terraform import` command is used to bring existing infrastructure resources, created outside of Terraform, into the Terraform state file, enabling their management through Terraform going forward.

To update a statefile with an existing org logging resource, execute the following import command:

```terraform
terraform import cpln_org_logging.RESOURCE_NAME ORG_NAME
```

-> 1. Substitute RESOURCE_NAME with the same string that is defined in the HCL file.<br/>2. Substitute ORG_NAME with the target org.
