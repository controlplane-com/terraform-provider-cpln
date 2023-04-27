---
page_title: "cpln_memcache Resource - terraform-provider-cpln"
subcategory: "Memcache Cluster"
description: |-

---

# cpln_memcache (Resource)

//TODO: Add resource description + add documentation link.

## Declaration

### Required

- **name** (String) Name of the Memcache Cluster.
- **node_count** (Number) //TODO: Add description.
- **node_size** (Float) //TODO: Add description.

### Optional

- **description** (String) Description of the Memcache Cluster.
- **tags** (Map of String) Key-value map of resource tags.
- **version** (String) Either 1.6.17 or 1.5.22 //TODO: Check this description.
- **options** (Block List, Max: 1) Memcache Cluster Options ([see below](#nestedblock--options)).

<a id="nestedblock--options"></a>

### `options`

Optional:

- **eviction_disabled** (Boolean) //TODO: Add description.
- **idle_timeout_seconds** (Number) //TODO: Add description.
- **max_item_size** (Number) //TODO: Add description.
- **max_connections** (Number) //TODO: Add description.

## Outputs

The following attributes are exported:

- **self_link** (String) Full link to this resource. Can be referenced by other resources. 

## Example Usage

```terraform
resource "cpln_memcache" "example" {
    name 		= "memcache-example"
    description = "Memcache description for memcache-example" 
    
    tags = {
        terraform_generated = "true"
        acceptance_test 	= "true"
    }

    node_count = 1
    node_size  = 0.3
    version    = "1.5.22"

    options {
        eviction_disabled 	 = true
        idle_timeout_seconds = 600
        max_item_size 		 = 1024
        max_connections      = 1024
    }

    locations  = ["/org/{your-org}/location/aws-us-west-2"]
}
```