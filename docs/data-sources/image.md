---
page_title: "cpln_image Data Source - terraform-provider-cpln"
subcategory: "Image"
description: |-
---

# cpln_image (Data Source)

Use this data source to access information about an [Image](https://docs.controlplane.com/reference/image) within Control Plane.

## Required

- **name** (String) Name of the image. If the tag of the image is not specified, the latest image will be fetched for this data source.

## Outputs

The following attributes are exported:

- **cpln_id** (String) The ID, in GUID format, of the Image.
- **name** (String) Name of the Image.
- **tags** (Map of String) Key-value map of resource tags.
- **self_link** (String) Full link to this resource. Can be referenced by other resources.
- **tag** (String) Tag of the image.
- **repository** (String) Respository name of the image.
- **digest** (String) A unique SHA256 hash used to identify a specific image version within the image registry.
- **manifest** (Block List, Max: 1) ([see below](#nestedblock--manifest))

<a id="nestedblock--manifest"></a>

### `manifest`

 The manifest provides configuration and layers information about the image. It plays a crucial role in the Docker image distribution system, enabling image creation, verification, and replication in a consistent and secure manner.

- **config** (Block List, Max: 1) ([see below](#nestedblock--config--layers)).
- **layers** (Block List) ([see below](#nestedblock--config--layers)).
- **media_type** (String) Specifies the type of the content represented in the manifest, allowing Docker clients and registries to understand how to handle the document correctly.
- **schema_version** (Number) The version of the Docker Image Manifest format.

<a id="nestedblock--config--layers"></a>

### `config` and `layers`

The config is a JSON blob that contains the image configuration data which includes environment variables, default command to run, and other settings necessary to run the container based on this image.

Layers lists the digests of the image's layers. These layers are filesystem changes or additions made in each step of the Docker image's creation process. The layers are stored separately and pulled as needed, which allows for efficient storage and transfer of images. Each layer is represented by a SHA256 digest, ensuring the integrity and authenticity of the image.

- **size** (Number) The size of the image or layer in bytes. This helps in estimating the space required and the download time.
- **digest** (String) A unique SHA256 hash used to identify a specific image version within the image registry.
- **media_type** (String) Specifies the type of the content represented in the manifest, allowing Docker clients and registries to understand how to handle the document correctly.

## Example Usage

```terraform
# Get latest image
data "cpln_image" "image-name-only" {
    name = "IMAGE_NAME"
}

# Get Specific image
data "cpln_image" "image-name-with-tag" {
    name = "IMAGE_NAME:TAG"
}

output "latest_image" {
  value = data.cpln_image.image-name-only
}

output "specific_image" {
  value = data.cpln_image.image-name-with-tag
}
```
