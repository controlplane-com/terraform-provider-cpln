# Terraform Plugin Framework: Code Snippets

A collection of copy-&-pasteable Go templates for common Terraform-Plugin-Framework patterns.

---

## 1. BuildList (Single-Block)

Use this when your schema block only ever produces a single item. Just change the `Name`, `Model`, `Operator`, `ClientType` and field mappings.

```go
// buildName constructs a ClientType from the given Terraform state.
func (opr *Operator) buildName(state types.List) *client.ClientType {
  // Convert Terraform list into model blocks using generic helper
  blocks, ok := BuildList[models.Model](opr.Ctx, opr.Diags, state)

	// Return nil if conversion failed or list was empty
  if !ok {
    return nil
  }

  // Take the first (and only) block
  block := blocks[0]

  // Construct and return the output
  return &client.ClientType{}
}
```


## 2. BuildList (Iterate-Over-Blocks)

When your list may contain multiple items and you need to transform each into a Go value:

```go
// buildName constructs a []client.ClientType from the given Terraform state.
func (opr *Operator) buildName(state types.List) *[]client.ClientType {
	// Convert Terraform list into model blocks using generic helper
  blocks, ok := BuildList[models.Model](opr.Ctx, opr.Diags, state)

  // Return nil if conversion failed or list was empty
  if !ok {
    return nil
  }

  // Prepare the output slice
  output := []client.ClientType{}

  // Iterate over each block and construct an output item
  for _, block := range blocks {
    // Construct the item
    item := client.ClientType{}

    // Add the item to the output slice
    output = append(output, item)
  }

	// Return a pointer to the output
  return &output
}
```

## 3. FlattenList (Single-Block)

Reverse of the single-block `BuildList`—wrap one struct in a `types.List`:

```go
// flatten{{Name}} transforms *client.{{ClientType}} into a Terraform types.List.
func (opr *{{Operator}}) flatten{{Name}}(input *client.{{ClientType}}) types.List {
	// Get attribute types
  elementType := models.{{Model}}Model{}.AttributeTypes()

	// Check if the input is nil
  if input == nil {
		// Return a null list
    return types.ListNull(elementType)
  }

  // Build a single block
  block := models.{{Model}}Model{}

	// Return the successfully created types.List
  return FlattenList(opr.Ctx, opr.Diags, []models.{{Model}}Model{block})
}
```

## 4. FlattenList (Iterate-Over-Input)

When you need to turn a `[]client.Foo` or `*[]client.Foo` into a Terraform `types.List` of many blocks:

```go
// flatten{{Name}} transforms *[]client.{{ClientType}} into a Terraform types.List.
func (opr *{{Operator}}) flatten{{Name}}(input *[]client.{{ClientType}}) types.List {
	// Get attribute types
  elementType := models.{{Model}}Model{}.AttributeTypes()

	// Check if the input is nil
  if input == nil {
		// Return a null list
    return types.ListNull(elementType)
  }

	// Define the blocks slice
  var blocks []models.{{Model}}Model

	// Iterate over the slice and construct the blocks
  for _, item := range *input {
		// Construct a block
    block := models.{{Model}}Model{}

		// Append the constructed block to the blocks slice
    blocks = append(blocks, block)
  }

	// Return the successfully created types.List
  return FlattenList(opr.Ctx, opr.Diags, blocks)
}
```
