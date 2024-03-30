---
page_title: "cpln_image Data Source - terraform-provider-cpln"
subcategory: "Image"
description: |-
---

# cpln_image (Data Source)

Use this data source to access information about an [Image](https://docs.controlplane.com/reference/image) within Control Plane.

## Required

- **name** (String) Name of the image.

## Outputs

The following attributes are exported:

- **cpln_id** (String) The ID, in GUID format, of the Image.
- **name** (String) Name of the Image.
- **tags** (Map of String) Key-value map of resource tags.
- **self_link** (String) Full link to this resource. Can be referenced by other resources.
- **tag** (String) Tag of the image.
- **repository** (String) Respository name of the image.
- **digest** (String) // TODO: Add description
- **manifest** (Block List, Max: 1) ([see below](#nestedblock--manifest))

<a id="nestedblock--manifest"></a>

### `manifest`

// TODO: Add description

- **config** (Block List, Max: 1) ([see below](#nestedblock--config)).
- **layers** (Block List) ([see below](#nestedblock--config)).
- **media_type** (String) // TODO: Add description.
- **schema_version** (Number) // TODO: Add description.

<a id="nestedblock--config"></a>

### `config`

// TODO: Add description

- **size** (Number) // TODO: Add description.
- **digest** (String) // TODO: Add description.
- **media_type** (String) // TODO: Add description.

## Example Usage

```terraform
data "cpln_image" "image" {
    name = "IMAGE_NAME:TAG"
}

output "image" {
  value = data.cpln_image.image
}

output "image_tag" {
  value = data.cpln_image.image.tag
}
```
