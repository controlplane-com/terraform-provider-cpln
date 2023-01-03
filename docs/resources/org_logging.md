---
page_title: "cpln_org_logging Resource - terraform-provider-cpln"
subcategory: "Org Logging"
description: |-
  
---

# cpln_org_logging (Resource)

Manages an [Org's external logging configuration](https://docs.controlplane.com/reference/org#external-logging).


## Declaration

### Required 

At least one of the following logging blocks are required:
  
- **coralogix_logging** (Block List, Max: 1) ([see below](#nestedblock--coralogix_logging)).
- **datadog_logging** (Block List, Max: 1) ([see below](#nestedblock--datadog_logging)).
- **logzio_logging** (Block List, Max: 1) ([see below](#nestedblock--logzio_logging)).
- **s3_logging** (Block List, Max: 1) ([see below](#nestedblock--s3_logging)).



<a id="nestedblock--coralogix_logging"></a>
 ### `coralogix_logging`

[Documentation Reference](https://docs.controlplane.com/reference/org#coralogix)

Required:

- **cluster** (String) Coralogix cluster URI. 
- **credentials** (String) Full link to referenced Opaque Secret. 
- **app** (String) App name to be displayed in Coralogix dashboard. 
- **subsystem** (String) Subsystem name to be displayed in Coralogix dashboard. 

~> **Note** Valid clusters: `coralogix.com`, `coralogix.us`, `app.coralogix.in`, `app.eu2.coralogix.com`, `app.coralogixsg.com`.

~> **Note** Supported variables for App and Subsystem are: freeformed or `{org}`, `{gvc}`, `{workload}`, `{location}`.



<a id="nestedblock--datadog_logging"></a>
 ### `datadog_logging`

[Documentation Reference](https://docs.controlplane.com/reference/org#coralogix)

Required:

- **host** (String) Datadog host URI. 
- **credentials** (String) Full link to referenced Opaque Secret. 

~> **Note** Valid Hosts: `http-intake.logs.datadoghq.com`, `http-intake.logs.us3.datadoghq.com`, `http-intake.logs.us5.datadoghq.com`, `http-intake.logs.datadoghq.eu`.

<a id="nestedblock--logzio_logging"></a>
 ### `logzio_logging`

[Documentation Reference](https://docs.controlplane.com/reference/org#logzio)

 Required:

- **listener_host** (String) Logzio listener host URI. 
- **credentials** (String) Full link to referenced Opaque Secret. 

~> **Note** Valid listener hosts: `listener.logz.io`, `listener-nl.logz.io`


<a id="nestedblock--s3_logging"></a>
 ### `s3_logging`

[Documentation Reference](https://docs.controlplane.com/reference/org#s3)

Required:

- **bucket** (String) Name of S3 bucket. 
- **region** (String) AWS region where bucket is located. 
- **prefix** (String) Bucket path prefix. Default: "/".
- **credentials** (String) Full link to referenced AWS Secret. 

## Outputs

The following attributes are exported:

- **cpln_id** (String) ID, in GUID format, of the Identity.
- **name** (String) The name of Org.
- **description** (String) The description of Org.
- **tags** (Map of String) Key-value map of the Org's tags.

## Example Usage

### Coralogix

```terraform

   resource "cpln_secret" "opaque" {

	    name = "opaque-random-coralogix-tbd"
		description = "opaque description opaque-random-tbd" 
		
		tags = {
			terraform_generated = "true"
			acceptance_test = "true"
			secret_type = "opaque"
		}

		opaque {
			payload = "opaque_secret_payload"
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
			app = "{workload}"
			subsystem = "{org}"
		}
	}	  	

```

### Datadog

```terraform

    resource "cpln_secret" "opaque" {

		name = "opaque-random-datadog-tbd"
		description = "opaque description" 
		
		tags = {
			terraform_generated = "true"
			acceptance_test = "true"
			secret_type = "opaque"
		}

		opaque {
			payload = "opaque_secret_payload"
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

### Logzio

```terraform

    resource "cpln_secret" "opaque" {

		name = "opaque-random-datadog-tbd"
		description = "opaque description" 
		
		tags = {
			terraform_generated = "true"
			acceptance_test = "true"
			secret_type = "opaque"
		}

		opaque {
			payload = "opaque_secret_payload"
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

### S3

```terraform

    resource "cpln_secret" "aws" {

		name = "aws-random-tbd"
		description = "aws description aws-random-tbd" 
				
		tags = {
			terraform_generated = "true"
			acceptance_test = "true"
			secret_type = "aws"
		} 
		
		aws {
			secret_key = "AKIAIOSFODNN7EXAMPLE"
			access_key = "AKIAwJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
			role_arn = "arn:awskey" 
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