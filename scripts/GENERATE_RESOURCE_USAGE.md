# Resource Generator Script

This script helps you scaffold a new Terraform provider resource and its acceptance tests from predefined templates.

## Prerequisites

- Python 3.6+

## Keys and Their Meanings

1. `RESOURCE_TYPE_NAME`
 - **Description**: The lower-case, snake-style name of the resource (e.g., `identity`).  
 - **Usage**: Used in filenames (`resource_{RESOURCE_TYPE_NAME}.go`), Terraform type name (`cpln_{RESOURCE_TYPE_NAME}`), and test resource block.

2. `RESOURCE_NAME`
 - **Description**: The PascalCase name of the resource (e.g., `Identity`).  
 - **Usage**: Used for Go type names, struct names, and constructor names.

3. `API_OBJECT_NAME`
 - **Description**: The Go type of the API object in the client (e.g., `Identity`).  
 - **Usage**: Used in `client.{API_OBJECT_NAME}` references.

4. `RESOURCE_INSTANCE_NAME`
 - **Description**: The variable name for the resource receiver (e.g., `ir`).  
 - **Usage**: Used in Go method receivers (e.g., `func (ir *IdentityResource) ...`).

5. `RESOURCE_STRING_NAME`
 - **Description**: The string key passed to `ResourceBaseAttributes` (e.g., `"identity"`).  
 - **Usage**: Used in schema definition: `ResourceBaseAttributes("{RESOURCE_STRING_NAME}")`.

6. `RESOURCE_VAR_NAME`
 - **Description**: The lower-case prefix for test variables (e.g., `identity`).  
 - **Usage**: Used in test struct names and function calls (`identityTest.name`).

7. `RESOURCE_CAPITAL_NAME`
 - **Description**: The name used in descriptive text within tests (e.g., `IDENTITY`).  
 - **Usage**: Used in acceptance test names and descriptions.


## Generating Files

Run the script from the repository root:

```bash
scripts/generate_resource.py \
RESOURCE_TYPE_NAME=identity \
RESOURCE_NAME=Identity \
API_OBJECT_NAME=Identity \
RESOURCE_INSTANCE_NAME=r \
RESOURCE_STRING_NAME="identity" \
RESOURCE_VAR_NAME=identity \
RESOURCE_CAPITAL_NAME=Identity
```

If you omit any key, the script will prompt you interactively:

```bash
Missing parameters: API_OBJECT_NAME, RESOURCE_VAR_NAME
(1/7) Enter value for API_OBJECT_NAME: Identity
(2/7) Enter value for RESOURCE_VAR_NAME: identity
...
```

After successful execution, you will see:

```bash
Generated: /path/to/project/internal/provider/resource_identity.go
Generated: /path/to/project/internal/provider/resource_identity_test.go
```