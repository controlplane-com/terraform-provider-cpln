---
page_title: "cpln_mk8s Resource - terraform-provider-cpln"
subcategory: "Mk8s"
description: |-
---

# cpln_mk8s (Resource)

Manages a Mk8s's [Mk8s](https://docs.controlplane.com/mk8s/overview).

## Declaration

### Required

- **name** (String) Name of the Mk8s.
- **version** (String) TODO: Add description

~> **Note** Only one of the providers listed below can be included in a resource.

- **generic_provider** (Block List, Max: 1) ([see below](#nestedblock--generic_provider))
- **hetzner_provider** (Block List, Max: 1) ([see below](#nestedblock--hetzner_provider))
- **aws_provider** (Block List, Max: 1) ([see below](#nestedblock--aws_provider))

### Optional

- **description** (String) Description of the Mk8s.
- **tags** (Map of String) Key-value map of resource tags.
- **firewall** (Block List, Max: 1) ([see below](#nestedblock--firewall))
- **add_ons** (Block List, Max: 1) ([see below](#nestedblock--add_ons))

<a id="nestedblock--generic_provider"></a>

### `generic_provider`

TODO: Add description

Required:

- **location** (String) TODO: Add description

Optional:

- **networking** (Block List, Max: 1) ([see below](#nestedblock--generic_provider--networking))
- **node_pool** (Block List) ([see below](#nestedblock--generic_provider--node_pool))

<a id="nestedblock--generic_provider--networking"></a>

### `generic_provider.networking`

TODO: Add description

Optional:

- **service_network** (String) TODO: Add description
- **pod_network** (String) TODO: Add description

<a id="nestedblock--generic_provider--node_pool"></a>

### `generic_provider.node_pool`

TODO: Add description

Required:

- **name** (String) TODO: Add description

Optional:

- **labels** (Map of String) TODO: Add description
- **taint** (Block List) ([see below](#nestedblock--generic_provider--node_pool--taint))

<a id="nestedblock--generic_provider--node_pool--taint"></a>

### `generic_provider.node_pool.taint`

TODO: Add description

Optional:

- **key** (String) TODO: Add description
- **value** (String) TODO: Add description
- **effect** (String) TODO: Add description

<a id="nestedblock--hetzner_provider"></a>

### `hetzner_provider`

TODO: Add description

Required:

- **region** (String) TODO: Add description
- **token_secret_link** (String) TODO: Add description
- **network_id** (String) TODO: Add description

Optional:

- **networking** (Block List, Max: 1) ([see below](#nestedblock--generic_provider--networking))
- **pre_install_script** (String) TODO: Add description
- **firewall_id** (String) TODO: Add description
- **node_pool** (Block List) ([see below](#nestedblock--hetzner_provider--node_pool))
- **dedicated_server_node_pool** (Block List) ([see below](#nestedblock--generic_provider--node_pool))
- **image** (String) TODO: Add description
- **ssh_key** (String) TODO: Add description
- **autoscaler** (Block List, Max: 1) ([see below](#nestedblock--autoscaler))

<a id="nestedblock--hetzner_provider--node_pool"></a>

### `hetzner_provider.node_pool`

TODO: Add description

Required:

- **name** (String) TODO: Add description
- **server_type** (String) TODO: Add description

Optional:

- **labels** (Map of String) TODO: Add description
- **taint** (Block List) ([see below](#nestedblock--generic_provider--node_pool--taint))
- **override_image** (String) TODO: Add description
- **min_size** (Number) TODO: Add description
- **max_size** (Number) TODO: Add description

<a id="nestedblock--aws_provider"></a>

### `aws_provider`

TODO: Add description

Required:

- **region** (String) TODO: Add description
- **skip_create_roles** (Boolean) TODO: Add description
- **token_secret_link** (String) TODO: Add description
- **network_id** (String) TODO: Add description
- **image** (Block List, Max: 1) ([see below](#nestedblock--aws_provider--ami))
- **deploy_role_arn** (String) Control Plane will set up the cluster by assuming this role
- **vpc_id** (String) TODO: Add description

Optional:

- **networking** (Block List, Max: 1) ([see below](#nestedblock--generic_provider--networking))
- **pre_install_script** (String) TODO: Add description
- **key_pair** (String) TODO: Add description
- **disk_encryption_key_arn** (String) TODO: Add description
- **security_group_ids** (List of String) TODO: Add description
- **node_pool** (Block List) ([see below](#nestedblock--aws_provider--node_pool))
- **autoscaler** (Block List, Max: 1) ([see below](#nestedblock--autoscaler))

<a id="nestedblock--aws_provider--node_pool"></a>

### `aws_provider.node_pool`

TODO: Add description

Required:

- **name** (String) TODO: Add description
- **instance_types** (List of String) TODO: Add description
- **override_image** (Block List, Max: 1) ([see below](#nestedblock--aws_provider--ami))
- **subnet_ids** (List of String) TODO: Add description

Optional:

- **labels** (Map of String) TODO: Add description
- **taint** (Block List) ([see below](#nestedblock--generic_provider--node_pool--taint))
- **boot_disk_size** (Number) TODO: Add description
- **min_size** (Number) TODO: Add description
- **max_size** (Number) TODO: Add description
- **on_demand_base_capacity** (Number) TODO: Add description
- **on_demand_percentage_above_base_capacity** (Number) TODO: Add description
- **spot_allocation_strategy** (String) TODO: Add description
- **extra_security_group_ids** (List of String) TODO: Add description

<a id="nestedblock--aws_provider--ami"></a>

### `ami`

TODO: Add description

Required:

~> **Note** Only one of the following listed below can be included.

- **recommended** (String) TODO: Add description
- **exact** (String) TODO: Add description

<a id="nestedblock--autoscaler"></a>

### `autoscaler`

TODO: Add description

Optional:

- **expander** (List of String) TODO: Add description
- **unneeded_time** (String) TODO: Add description
- **unready_time** (String) TODO: Add description
- **utilization_threshold** (Float64) TODO: Add description

<a id="nestedblock--firewall"></a>

### `firewall`

TODO: Add description

Required:

- **source_cidr** (String) TODO: Add description

Optional:

- **description** (String) TODO: Add description

<a id="nestedblock--add_ons"></a>

### `add_ons`

TODO: Add description

Optional:

- **dashboard** (Boolean) TODO: Add description
- **azure_workload_identity** (Block List, Max: 1) ([see below](#nestedblock--add_ons--azure_workload_identity))
- **aws_workload_identity** (Boolean) TODO: Add description
- **local_path_storage** (Boolean) TODO: Add description
- **metrics** (Block List, Max: 1) ([see below](#nestedblock--add_ons--metrics))
- **logs** (Block List, Max: 1) ([see below](#nestedblock--add_ons--logs))
- **nvidia** (Block List, Max: 1) ([see below](#nestedblock--add_ons--nvidia))
- **aws_efs** (Block List, Max: 1) ([see below](#nestedblock--add_ons--aws))
- **aws_ecr** (Block List, Max: 1) ([see below](#nestedblock--add_ons--aws))
- **aws_elb** (Block List, Max: 1) ([see below](#nestedblock--add_ons--aws))
- **azure_acr** (Block List, Max: 1) ([see below](#nestedblock--add_ons--azure_acr))

<a id="nestedblock--add_ons--azure_workload_identity"></a>

### `add_ons.azure_workload_identity`

TODO: Add description

Required:

- **tenant_id** (String) TODO: Add description

<a id="nestedblock--add_ons--metrics"></a>

### `add_ons.metrics`

TODO: Add description

Optional:

- **kube_state** (Boolean) TODO: Add description
- **core_dns** (Boolean) TODO: Add description
- **kubelet** (Boolean) TODO: Add description
- **api_server** (Boolean) TODO: Add description
- **node_exporter** (Boolean) TODO: Add description
- **cadvisor** (Boolean) TODO: Add description
- **scrape_annotated** (Block List, Max: 1) ([see below](#nestedblock--add_ons--metrics--scrape-annotated))

<a id="nestedblock--add_ons--metrics--scrape-annotated"></a>

### `add_ons.metrics.scrape-annotated`

TODO: Add description

Optional:

- **interval_seconds** (Number) TODO: Add description
- **include_namespaces** (String) TODO: Add description
- **exclude_namespaces** (String) TODO: Add description
- **retain_labels** (String) TODO: Add description

<a id="nestedblock--add_ons--logs"></a>

### `add_ons.logs`

TODO: Add description

Optional:

- **audit_enabled** (Boolean) TODO: Add description
- **include_namespaces** (String) TODO: Add description
- **exclude_namespaces** (String) TODO: Add description

<a id="nestedblock--add_ons--nvidia"></a>

### `add_ons.nvidia`

TODO: Add description

Required:

- **taint_gpu_nodes** (Boolean) TODO: Add description

<a id="nestedblock--add_ons--aws"></a>

### `add_ons.aws`

TODO: Add description

Required:

- **role_arn** (String) TODO: Add description

<a id="nestedblock--add_ons--azure_acr"></a>

### `add_ons.azure_acr`

TODO: Add description

Required:

- **client_id** (String) TODO: Add description

## Outputs

The following attributes are exported:

- **cpln_id** (String) The ID, in GUID format, of the Mk8s.
- **alias** (String) The alias name of the Mk8s.
- **self_link** (String) Full link to this resource. Can be referenced by other resources.
- **status** (Block List, Max: 1) ([see below](#nestedblock--status)).

<a id="nestedblock--status"></a>

### `status`

Status of the mk8s.

Read-Only:

- **oidc_provider_url** (String) TODO: Add description
- **server_url** (String) TODO: Add description
- **home_location** (String) TODO: Add description
- **add_ons** (Block List, Max: 1) ([see below](#nestedblock--status--add_ons))

<a id="nestedblock--status--add_ons"></a>

### `status.add_ons`

TODO: Add description

Read-Only:

- **dashboard** (Block List, Max: 1) ([see below](#nestedblock--status--add_ons--dashobard))
- **aws_workload_identity** (Block List, Max: 1) ([see below](#nestedblock--status--add_ons--aws_workload_identity))
- **metrics** (Block List, Max: 1) ([see below](#nestedblock--status--add_ons--metrics))
- **logs** (Block List, Max: 1) ([see below](#nestedblock--status--add_ons--logs))
- **aws_ecr** (Block List, Max: 1) ([see below](#nestedblock--status--add_ons--aws))
- **aws_efs** (Block List, Max: 1) ([see below](#nestedblock--status--add_ons--aws))
- **aws_elb** (Block List, Max: 1) ([see below](#nestedblock--status--add_ons--aws))

<a id="nestedblock--status--add_ons--dashobard"></a>

### `status.add_ons.dashboard`

TODO: Add description

Read-Only:

- **url** (String) TODO: Add description

<a id="nestedblock--status--add_ons--aws_workload_identity"></a>

### `status.add_ons.aws_workload_identity`

TODO: Add description

Read-Only:

- **oidc_provider_config** (Block List, Max: 1) ([see below](#nestedblock--status--add_ons--aws_workload_identity--oidc_provider_config))
- **trust_policy** (String) TODO: Add description

<a id="nestedblock--status--add_ons--aws_workload_identity--oidc_provider_config"></a>

### `status.add_ons.aws_workload_identity.oidc_provider_config`

TODO: Add description

Read-Only:

- **provider_url** (String) TODO: Add description
- **audience** (String) TODO: Add description

<a id="nestedblock--status--add_ons--metrics"></a>

### `status.add_ons.metrics`

TODO: Add description

Read-Only:

- **prometheus_endpoint** (String) TODO: Add description
- **remote_write_config** (String) TODO: Add description

<a id="nestedblock--status--add_ons--logs"></a>

### `status.add_ons.logs`

TODO: Add description

Read-Only:

- **loki_address** (String) TODO: Add description

<a id="nestedblock--status--add_ons--aws"></a>

### `status.add_ons.aws`

TODO: Add description

Read-Only:

- **trust_policy** (String) TODO: Add description

## Example Usage - Generic Provider

```terraform
resource "cpln_mk8s" "generic" {

    name        = "demo-mk8s-generic-provider"
    description = "demo-mk8s-generic-provider"

    tags = {
        terraform_generated = "true"
        acceptance_test     = "true"
    }

    version = "1.28.4"

    firewall {
        source_cidr = "192.168.1.255"
        description = "hello world"
    }
    
    generic_provider {
        location = "aws-eu-central-1"
        
        networking {
            service_network = "10.43.0.0/16"
            pod_network 	= "10.42.0.0/16"
        }
        
        node_pool {
            name = "my-node-pool"

            labels = {
                hello = "world"
            }

            taint {
                key    = "hello"
                value  = "world"
                effect = "NoSchedule"
            }
        }
    }

    add_ons {
        dashboard = true

        azure_workload_identity {
            tenant_id = "7f43458a-a34e-4bfa-9e56-e2289e49c4ec"
        }

        aws_workload_identity = true
        local_path_storage    = true

        metrics {
            kube_state    = true
            core_dns      = true
            kubelet       = true
            api_server    = true
            node_exporter = true
            cadvisor      = true

            scrape_annotated {
                interval_seconds   = 30
                include_namespaces = "^\\d+$"
                exclude_namespaces  = "^[a-z]$"
                retain_labels      = "^\\w+$"
            }
        }

        logs {
            audit_enabled      = true
            include_namespaces = "^\\d+$"
            exclude_namespaces  = "^[a-z]$"
        }

        nvidia {
            taint_gpu_nodes = true
        }

        azure_acr {
            client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
        }
    }
}
```

## Example Usage - Hetzner Provider

```terraform
resource "cpln_mk8s" "hetzner" {
    
    name        = "demo-mk8s-hetzner-provider"
    description = "demo-mk8s-hetzner-provider"

    tags = {
        terraform_generated = "true"
        acceptance_test     = "true"
    }

    version = "1.28.4"

    firewall {
        source_cidr = "192.168.1.255"
        description = "hello world"
    }
    
    hetzner_provider {
        
        region = "fsn1"

        networking {
            service_network = "10.43.0.0/16"
            pod_network 	= "10.42.0.0/16"
        }

        pre_install_script = "#! echo hello world"
        token_secret_link  = "/org/terraform-test-org/secret/hetzner"
        network_id 		   = "2808575"

        node_pool {
            name = "my-hetzner-node-pool"

            labels = {
                hello = "world"
            }

            taint {
                key    = "hello"
                value  = "world"
                effect = "NoSchedule"
            }

            server_type    = "cx11"
            override_image = "debian-11"
            min_size 	   = 0
            max_size 	   = 0
        }

        dedicated_server_node_pool {
            name = "my-node-pool"

            labels = {
                hello = "world"
            }

            taint {
                key    = "hello"
                value  = "world"
                effect = "NoSchedule"
            }
        }

        image 	= "centos-7"
        ssh_key = "10925607"

        autoscaler {
            expander 	  		  = ["most-pods"]
            unneeded_time         = "10m"
            unready_time  		  = "20m"
            utilization_threshold = 0.7
        }
    }

    add_ons {
        dashboard = true

        azure_workload_identity {
            tenant_id = "7f43458a-a34e-4bfa-9e56-e2289e49c4ec"
        }

        aws_workload_identity = true
        local_path_storage    = true

        metrics {
            kube_state    = true
            core_dns      = true
            kubelet       = true
            api_server    = true
            node_exporter = true
            cadvisor      = true

            scrape_annotated {
                interval_seconds   = 30
                include_namespaces = "^\\d+$"
                exclude_namespaces  = "^[a-z]$"
                retain_labels      = "^\\w+$"
            }
        }

        logs {
            audit_enabled      = true
            include_namespaces = "^\\d+$"
            exclude_namespaces  = "^[a-z]$"
        }

        nvidia {
            taint_gpu_nodes = true
        }

        azure_acr {
            client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
        }
    }
}
```

## Example Usage - AWS Provider

```terraform
resource "cpln_mk8s" "aws" {

    name        = "demo-mk8s-aws-provider"
    description = "demo-mk8s-aws-provider"

    tags = {
        terraform_generated = "true"
        acceptance_test     = "true"
    }

    version = "1.28.4"

    firewall {
        source_cidr = "192.168.1.255"
        description = "hello world"
    }

    aws_provider {

        region            = "eu-central-1"
        skip_create_roles = false

        networking {
            service_network = "10.43.0.0/16"
            pod_network 	= "10.42.0.0/16"
        }

        pre_install_script = "#! echo hello world"

        image {
            recommended = "amazon/al2023"
        }

        deploy_role_arn         = "arn:aws:iam::12345678901:role/cpln"
        vpc_id                  = "vpc-03105bd4dc058d3a8"
        key_pair                = "cem_uzak"
        disk_encryption_key_arn = "arn:aws:kms:eu-central-1:12345678901:key/0a1bcd23-4567-8901-e2fg-3h4i5jk678lm"

        security_group_ids = ["sg-031480aa7a1e6e38b"]

        node_pool {
            name = "my-aws-node-pool"

            labels = {
                hello = "world"
            }

            taint {
                key    = "hello"
                value  = "world"
                effect = "NoSchedule"
            }

            instance_types = ["t4g.nano"]

            override_image {
                exact = "ami-123"
            }

            boot_disk_size                           = 20
            min_size                                 = 0
            max_size                                 = 0
            on_demand_base_capacity                  = 0
            on_demand_percentage_above_base_capacity = 0
            spot_allocation_strategy                 = "lowest-price"

            subnet_ids               = ["subnet-0e564a042e2a45009"]
            extra_security_group_ids = ["sg-031480aa7a1e6e38b"]
        }

        autoscaler {
            expander 	  		  = ["most-pods"]
            unneeded_time         = "10m"
            unready_time  		  = "20m"
            utilization_threshold = 0.7
        }
    }

    add_ons {
        dashboard = true

        azure_workload_identity {
            tenant_id = "7f43458a-a34e-4bfa-9e56-e2289e49c4ec"
        }

        aws_workload_identity = true
        local_path_storage    = true

        metrics {
            kube_state    = true
            core_dns      = true
            kubelet       = true
            api_server    = true
            node_exporter = true
            cadvisor      = true

            scrape_annotated {
                interval_seconds   = 30
                include_namespaces = "^\\d+$"
                exclude_namespaces  = "^[a-z]$"
                retain_labels      = "^\\w+$"
            }
        }

        logs {
            audit_enabled      = true
            include_namespaces = "^\\d+$"
            exclude_namespaces  = "^[a-z]$"
        }

        nvidia {
            taint_gpu_nodes = true
        }

        aws_efs {
            role_arn = "arn:aws:iam::123456789012:role/aws-efs-role"
        }

        aws_ecr {
            role_arn = "arn:aws:iam::123456789012:role/aws-ecr-role"
        }

        aws_elb {
            role_arn = "arn:aws:iam::123456789012:role/aws-elb-role"
        }

        azure_acr {
            client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
        }
    }
}
```