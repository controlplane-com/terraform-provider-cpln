# CLAUDE.md — terraform-provider-cpln

Project-specific guidance for Claude. **Match the existing code's style, structure, comments, and positioning exactly** — do not invent new conventions. Before writing any new file or function, find the closest existing analog and mirror it.

---

## 1. What this project is

Terraform provider for Control Plane, built on the **Terraform Plugin Framework** (NOT SDKv2). It exposes Control Plane resources (`cpln_gvc`, `cpln_workload`, `cpln_domain`, `cpln_mk8s`, `cpln_secret`, …) and their data sources. The provider proxies a typed Go client (`internal/provider/client/`) to the Control Plane REST API.

Module: `github.com/controlplane-com/terraform-provider-cpln`. Vendored deps live in `vendor/`.

## 2. Directory layout (authoritative)

```
internal/provider/
  client/                    # API client structs (one .go per kind: domain.go, mk8s.go, …)
  models/<kind>/<kind>.go    # Terraform plan/state models keyed by tfsdk tags
  modifiers/                 # Plan modifiers (description_modifier.go, tag_modifier.go, …)
  validators/                # Schema validators (name_validator.go, link_validator.go, …)
  common.go                  # EntityBaseModel, EntityBase, EntityOperator, generic CRUD
  helper.go                  # BuildX/FlattenX helpers, pointer helpers, misc utilities
  resource_<kind>.go         # Resource implementation
  resource_<kind>_test.go    # Acceptance tests
  data_source_<kind>.go      # Data source implementation
  data_source_<kind>_test.go # Data source acceptance tests
  provider.go                # Resource/data-source registry
docs/resources/<kind>.md     # User-facing resource docs
docs/data-sources/<kind>.md  # User-facing data source docs
templates/                   # Code skeletons — start here for new resources
  resource_skeleton.txt
  resource_test_skeleton.txt
  list/build_single.txt … set/flatten_multi.txt   # Builder/Flattener templates
  SNIPPETS.md
```

When creating anything new, first read the matching template under `templates/`, then read the closest sibling resource for non-skeleton patterns.

## 3. Architectural primitives (do not reinvent)

| Primitive | Where | Purpose |
|---|---|---|
| `EntityBaseModel` | `common.go` | Shared tfsdk fields (id, cpln_id, name, description, tags, self_link). Embed in every `<Kind>ResourceModel`. |
| `EntityBase` | `common.go` | Holds `*client.Client`. Embed in every `<Kind>Resource` / `<Kind>DataSource`. |
| `EntityOperator[Plan]` | `common.go` | Holds `Plan`, `Ctx`, `Diags`, `Client`. Embed in every `<Kind>ResourceOperator`. |
| `EntityOperations[Plan, API]` | `common.go` | Wires operator to generic CRUD. |
| `CreateGeneric` / `ReadGeneric` / `UpdateGeneric` / `DeleteGeneric` | `common.go` | The Create/Read/Update/Delete methods are one-line delegations. |
| `client.Base` | `client/base.go` | Embedded by every API kind; holds Name, Description, Tags, etc. |

The operator pattern: every resource implements an `EntityOperatorInterface` via four methods on `*<Kind>ResourceOperator`:

```
NewAPIRequest(isUpdate bool) client.<Kind>
MapResponseToState(apiResp *client.<Kind>, isCreate bool) <Kind>ResourceModel
InvokeCreate / InvokeRead / InvokeUpdate / InvokeDelete
```

## 4. Resource file structure (resource_<kind>.go)

Order is fixed — match it exactly. Section markers are mandatory.

```go
package cpln

import (
    // stdlib first
    "context"

    // module-internal next
    client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
    models "github.com/controlplane-com/terraform-provider-cpln/internal/provider/models/<kind>"

    // hashicorp last (validators, framework, types)
    ...
)

// Ensure resource implements required interfaces.
var (
    _ resource.Resource                = &<Kind>Resource{}
    _ resource.ResourceWithImportState = &<Kind>Resource{}
    // Add ResourceWithModifyPlan / ResourceWithValidateConfig only if implemented
)

/*** Resource Model ***/

// <Kind>ResourceModel holds the Terraform state for the resource.
type <Kind>ResourceModel struct {
    EntityBaseModel
    ... // top-level attributes here
}

/*** Resource Configuration ***/

// <Kind>Resource is the resource implementation.
type <Kind>Resource struct {
    EntityBase
    Operations EntityOperations[<Kind>ResourceModel, client.<Kind>]
}

// New<Kind>Resource returns a new instance of the resource implementation.
func New<Kind>Resource() resource.Resource { ... }

// Configure configures the resource before use.
func (xr *<Kind>Resource) Configure(...)

// ImportState sets up the import operation to map the imported ID to the "id" attribute in the state.
func (xr *<Kind>Resource) ImportState(...)

// Metadata provides the resource type name.
func (xr *<Kind>Resource) Metadata(...) { resp.TypeName = "cpln_<kind>" }

// Schema defines the schema for the resource.
func (xr *<Kind>Resource) Schema(...) { ... }

// Create creates the resource.
func (xr *<Kind>Resource) Create(...) { CreateGeneric(ctx, req, resp, xr.Operations) }
// Read fetches the current state of the resource.
func (xr *<Kind>Resource) Read(...)   { ReadGeneric(...) }
// Update modifies the resource.
func (xr *<Kind>Resource) Update(...) { UpdateGeneric(...) }
// Delete removes the resource.
func (xr *<Kind>Resource) Delete(...) { DeleteGeneric(...) }

/*** Plan Modifiers ***/        // OPTIONAL — only if the resource has plan-modifier helpers (RequiresReplace, etc.)

/*** Schemas ***/                // OPTIONAL — for shared sub-schemas referenced from Schema()

/*** Schemas Defaults ***/       // OPTIONAL — for default<X>Value() helpers used with objectdefault.StaticValue

/*** Resource Operator ***/

// <Kind>ResourceOperator is the operator for managing the state.
type <Kind>ResourceOperator struct {
    EntityOperator[<Kind>ResourceModel]
}

// NewAPIRequest creates a request payload from a state model.
func (xro *<Kind>ResourceOperator) NewAPIRequest(isUpdate bool) client.<Kind> { ... }

// MapResponseToState constructs the Terraform state model from the API response payload.
func (xro *<Kind>ResourceOperator) MapResponseToState(apiResp *client.<Kind>, isCreate bool) <Kind>ResourceModel { ... }

// InvokeCreate / InvokeRead / InvokeUpdate / InvokeDelete    // one-liners delegating to xro.Client.<Kind>(...)

// Builders //

// build<Field> constructs a <client type> from the given Terraform state.
func (xro *<Kind>ResourceOperator) build<Field>(state types.Object) *client.<X> { ... }

// Flatteners //

// flatten<Field> transforms <client type> into a Terraform <types.X>.
func (xro *<Kind>ResourceOperator) flatten<Field>(input *client.<X>) types.Object { ... }

// Helpers //                    // OPTIONAL — resource-local helpers
```

**Major sections** use `/*** Title ***/`, **subsections** (Builders, Flatteners, Helpers, Test Cases, Configs) use `// Title //`. Both have a blank line before and after.

**Operator method receiver name** is the resource initials + `o`: `MK8S` → `mro`, `Domain` → `dro`, `<Kind>ResourceOperator` → `xro` (per the templates). Resource receiver omits the `o`: `mr`, `dr`.

## 5. Comment style (rigid)

Every function gets a one-line `//` comment immediately above. The comment is descriptive — never references PRs, tasks, or call sites.

```go
// build<Field> constructs a <ClientType> from the given Terraform state.
func (...) build<Field>(...) ...
```

Every logical step inside a function gets a short narrating comment. Not narrating thoughts — narrating what the next line does. Mirror the cadence in `resource_mk8s.go`:

```go
// Convert Terraform list into model blocks using generic helper
blocks, ok := BuildList[models.NetworkingModel](mro.Ctx, mro.Diags, state)

// Return nil if conversion failed or list was empty
if !ok {
    return nil
}

// Take the first (and only) block
block := blocks[0]

// Construct and return the output
return &client.Mk8sNetworkingConfig{ ... }
```

Comments live on their own line, never trailing. They are short (one line). The entire codebase uses this density — do not strip comments to "be concise"; do not bunch multiple steps under one comment.

## 6. Schema shape preference — Object over List-as-Block

**For new code, prefer single-nested objects (`schema.SingleNestedAttribute` + `types.Object`) over list-with-Max-1 blocks (`schema.ListNestedBlock` + `types.List`).** The `ListNestedBlock` style is older and reads as a block in HCL (`block { ... }`); the object style reads as an attribute (`block = { ... }`) — cleaner, less ambiguous, and matches the API shape better when the underlying field is a single object.

| Use this... | When the field is... | Terraform model type |
|---|---|---|
| `schema.SingleNestedAttribute` | a single nested object | `types.Object` |
| `schema.ListNestedAttribute` (or `ListNestedBlock` only when matching legacy) | an ordered collection / list of objects | `types.List` |
| `schema.SetNestedAttribute` | an unordered, duplicate-free collection | `types.Set` |
| `schema.MapAttribute` | a `map[string]string` (or other primitive map) | `types.Map` |
| primitives (`StringAttribute`, `Int32Attribute`, etc.) | scalars | `types.String`, `types.Int32`, etc. |

Use `schema.ListNestedBlock` only when (a) the existing resource you are editing already uses it everywhere and you must match for consistency, or (b) Terraform-side ergonomics genuinely require `block { ... }` syntax. Otherwise, prefer the object form.

When you reach for an object, mirror the BYOK config in `resource_mk8s.go` (`schema.SingleNestedAttribute` with `Optional: true`, `Default: objectdefault.StaticValue(default<X>Value())`) and the matching `*Model` with `AttributeTypes()` returning a `types.ObjectType`.

## 6.1 Builders & Flatteners — copy from templates

For every nested API field, write one builder and one flattener. Use the matching skeleton in `templates/`:

| Terraform shape | Build template | Flatten template |
|---|---|---|
| `types.Object` (single nested object — preferred) | use `BuildObject[T]` directly (see `helper.go`) | use `FlattenObject[T]` directly |
| `types.List` (single block, legacy) | `templates/list/build_single.txt` | `templates/list/flatten_single.txt` |
| `types.List` (multiple) | `templates/list/build_multi.txt` | `templates/list/flatten_multi.txt` |
| `types.Set` (single) | `templates/set/build_single.txt` | `templates/set/flatten_single.txt` |
| `types.Set` (multiple) | `templates/set/build_multi.txt` | `templates/set/flatten_multi.txt` |
| `[]models.X` already in hand | `templates/known_list/build_single.txt` | `templates/known_list/flatten_single.txt` |

Naming: `buildFooBarBaz` / `flattenFooBarBaz` — camelCase, mirroring the path through the model. For nested blocks, reuse the parent name as a prefix: `buildSpec`, `buildSpecPorts`, `buildSpecPortsCors`, `buildSpecPortsCorsAllowOrigins`.

Builder signature returns a **pointer** to the client type or to a slice. Returning `nil` on empty/failed conversion is the convention. For object builders the canonical body is:

```go
func (xro *<Kind>ResourceOperator) build<Field>(state types.Object) *client.<X> {
    block, ok := BuildObject[models.<X>Model](xro.Ctx, xro.Diags, state)
    if !ok || block == nil {
        return nil
    }
    return &client.<X>{
        Foo: BuildString(block.Foo),
        Bar: BuildBool(block.Bar),
    }
}
```

And the matching flattener:

```go
func (xro *<Kind>ResourceOperator) flatten<Field>(input *client.<X>) types.Object {
    elementType := models.<X>Model{}.AttributeTypes().(types.ObjectType)
    if input == nil {
        return types.ObjectNull(elementType.AttrTypes)
    }
    block := models.<X>Model{
        Foo: types.StringPointerValue(input.Foo),
        Bar: types.BoolPointerValue(input.Bar),
    }
    return FlattenObject(xro.Ctx, xro.Diags, &block)
}
```

(Reference implementation: `buildAddOnByokJuicefs` / `flattenAddOnByokJuicefs` in `resource_mk8s.go`.)

## 7. Models (`internal/provider/models/<kind>/<kind>.go`)

One file per kind, package = kind. Layout:

```go
package <kind>

import (
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

/*** Main Models ***/

// <Top> //

type <Top>Model struct {
    Field types.String `tfsdk:"field"`
    Nested types.List  `tfsdk:"nested"`
}

func (m <Top>Model) AttributeTypes() attr.Type {
    return types.ObjectType{
        AttrTypes: map[string]attr.Type{
            "field":  types.StringType,
            "nested": types.ListType{ElemType: NestedModel{}.AttributeTypes()},
        },
    }
}

// <Top> -> Nested //

type NestedModel struct { ... }
func (m NestedModel) AttributeTypes() attr.Type { ... }
```

Rules:
- **Every model has an `AttributeTypes() attr.Type` method.** No exceptions — `BuildList`/`BuildSet`/`BuildObject` and their flatteners rely on it.
- **Section divider for each model**: `// Path -> To -> Field //` describing the schema path.
- `tfsdk` tag uses **snake_case**; struct field uses **PascalCase**.
- Pick the model field type from the schema attribute kind (see Section 6):
  - `SingleNestedAttribute` → `types.Object` (preferred for new single nested objects)
  - `ListNestedAttribute` / `ListNestedBlock` → `types.List`
  - `SetNestedAttribute` → `types.Set`
  - `MapAttribute` → `types.Map`
  - `StringAttribute` / `Int32Attribute` / `Float32Attribute` / `Float64Attribute` / `BoolAttribute` → matching `types.<Scalar>`
- Joi schema mapping (per the API spec it mirrors): `.number()` → `types.Float64`, `.number().integer()` → `types.Int32`, `.string()` → `types.String`, `.boolean()` → `types.Bool`.

## 8. Client structs (`internal/provider/client/<kind>.go`)

```go
package cpln

type <Kind> struct {
    Base
    Spec        *<Kind>Spec `json:"spec,omitempty"`
    SpecReplace *<Kind>Spec `json:"$replace/spec,omitempty"` // Only when the API uses $replace
    Status      *<Kind>Status `json:"status,omitempty"`
}

type <Kind>Spec struct {
    Field    *string                 `json:"field,omitempty"`     // Enum: "a", "b"
    Nested   *Nested<Kind>Type       `json:"nested,omitempty"`
    List     *[]Nested<Kind>Type     `json:"list,omitempty"`
}
```

Rules:
- **Every field is a pointer** with `json:"<name>,omitempty"` — no exceptions, even bools and ints. This lets `omitempty` distinguish "absent" from "zero".
- Slice fields are `*[]T`, never `[]T`.
- For string enums, append a trailing comment listing valid values: `// Enum: "http", "http2", "tcp"`.
- Section dividers: `/*** Spec ***/`, `/*** Status ***/`, `/*** Spec Related ***/`.
- For PUT/replace semantics that the API supports, include both `Spec *T \`json:"spec,..."\`` and `SpecReplace *T \`json:"$replace/spec,..."\`` and have `NewAPIRequest` set whichever applies based on `isUpdate`.

## 9. Schema authoring rules

- Use `MergeAttributes(xr.EntityBaseAttributes("<Kind>"), map[string]schema.Attribute{ ... })` so the base id/cpln_id/name/description/tags/self_link attributes come for free.
- Every attribute has a `Description` string ending with a period. For enums, document valid values inline: `Description: "... Valid values: \`a\`, \`b\`. Default: \`a\`."`.
- Pull validators from `internal/provider/validators` (`NameValidator`, `LinkValidator`, `DescriptionValidator`, `TagValidator`) and from `terraform-plugin-framework-validators` (`stringvalidator.LengthAtLeast`, `int32validator.AtLeast`, etc.).
- Pull plan modifiers from `internal/provider/modifiers` and from `*planmodifier`.
- For computed-only fields keep them under `status` as a `ListNestedAttribute` that mirrors the API status block (see `resource_domain.go`, `resource_mk8s.go`).
- Use `objectdefault.StaticValue(default<X>Value())` for nested object defaults — write the helper under `/*** Schemas Defaults ***/`.

## 10. Test files (resource_<kind>_test.go)

Layout (canonical, see `resource_mk8s_test.go`):

```go
package cpln

import (
    "errors"; "fmt"; "testing"
    "github.com/hashicorp/terraform-plugin-log/tflog"
    "github.com/hashicorp/terraform-plugin-testing/terraform"
    "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
    "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

/*** Acceptance Test ***/

// TestAccControlPlane<Kind>_basic performs an acceptance test for the resource.
func TestAccControlPlane<Kind>_basic(t *testing.T) {
    resourceTest := New<Kind>ResourceTest()
    resource.Test(t, resource.TestCase{
        PreCheck:                 func() { testAccPreCheck(t, "<KIND-UPPER>") },
        ProtoV6ProviderFactories: GetProviderServer(),
        CheckDestroy:             resourceTest.CheckDestroy,
        Steps:                    resourceTest.Steps,
    })
}

/*** Resource Test ***/

type <Kind>ResourceTest struct {
    Steps      []resource.TestStep
    RandomName string
}

func New<Kind>ResourceTest() <Kind>ResourceTest { ... }

// CheckDestroy verifies that all resources have been destroyed.
func (xrt *<Kind>ResourceTest) CheckDestroy(s *terraform.State) error { ... }

// Test Scenarios //
// One func per scenario: NewMk8sGenericProviderScenario, NewMk8sAwsProviderScenario, …

// Test Cases //
// Build<Scenario>TestStep, Build<Scenario>Update1TestStep, …

// Configs //
// <Scenario>RequiredOnlyHcl, <Scenario>Update1Hcl, …
```

Test rules:
- Acceptance tests use `resource.TestCase` with `Steps`. State checks use `resource.ComposeAggregateTestCheckFunc`.
- Use the project helpers `c.GetDefaultChecks(...)`, `c.TestCheckResourceAttr(...)`, `c.TestCheckNestedBlocks(...)`, `c.TestCheckMapObjectAttr(...)`, `c.TestCheckObjectAttr(...)` from `provider_test.go`. Do NOT call `resource.TestCheckResourceAttr` directly when a method exists on the test case.
- HCL configs use `fmt.Sprintf` with `c.ResourceName`, `c.Name`, `c.DescriptionUpdate` placeholders. Match the indentation and positioning of existing scenarios — including the multiple-update step ordering (`initialStep, importStep, update1, update2, update3, update2, update1, initialStep`).
- When adding a new optional API field, update **every** HCL config that exercises the parent block AND the corresponding `Expected*` checker so the assertion is real. Do not leave `// TODO: Add <field> test here` placeholders behind — fill them in.
- Use `ExpectNonEmptyPlan: true` only when a known status-attribute drift is in play.

## 11. Reuse policy — `helper.go` and `common.go` are first stops

**Before writing any utility, build, or flatten function, look in `helper.go` and `common.go`.** These two files exist precisely to keep resources DRY. Reaching for `os`, `strconv`, `strings`, or hand-rolling a converter is almost always wrong — the project already has a helper.

### `internal/provider/helper.go` — generic, type-level utilities

Use this file when the work is **not specific to any one Control Plane resource** — it converts Terraform types ↔ Go primitives, formats links, or does generic string/JSON/HCL work.

**Build (Plan → API):**
`BuildString`, `BuildInt`, `BuildFloat32`, `BuildFloat64`, `BuildBool`, `BuildTags`, `BuildMapString`, `BuildSetString`, `BuildListString`, `BuildSetInt`, `BuildList[T]`, `BuildSet[T]`, `BuildObject[T]`.

**Flatten (API → State):**
`FlattenInt`, `FlattenFloat64`, `FlattenSelfLink`, `FlattenTags`, `FlattenMapString`, `FlattenSetString`, `FlattenListString`, `FlattenSetInt`, `FlattenList[T]`, `FlattenSet[T]`, `FlattenObject[T]`.

For built-in pointer values use the framework directly: `types.StringPointerValue`, `types.BoolPointerValue`, `types.Int32PointerValue`. For going the other way without a wrapper, the `Build*` helpers above are the project convention.

**Misc helpers in `helper.go`:** `MergeAttributes`, `GetNameFromSelfLink`, `GetSelfLink`, `GetSelfLinkWithGvc`, `BoolPointer`, `StringPointer`, `IntPointer`, `Float64Pointer`, `Float32Pointer`, `IsGvcScopedResource`, `GetInterface`, `ParseValueAndUnit`, `ToStringSlice`, `PreserveJSONFormatting`, `StringSliceToString`, `IntSliceToString`, `IntSliceToStringSlice`, `MapToHCL`.

**Add to `helper.go`** when:
- The function is type-level / not tied to any specific resource shape (e.g. a new `Build<X>` or `Flatten<X>` for a Terraform/Go primitive pair, a new pointer helper, a generic string formatter).
- The helper would be useful from any resource or data source.

### `internal/provider/common.go` — shared cross-resource builders/flatteners and base types

Use this file when the same builder/flattener would be **used by more than one resource or data source**. Existing examples already shared:

- **Base entity infra:** `EntityBaseModel` (Fill/From/GetID), `EntityBase` (`EntityBaseConfigure`, `EntityBaseAttributes`, `NameSchema`, `DescriptionSchema`), `EntityOperator[Plan]` (`Init`, `BuildSetString`, `BuildSetInt`, `BuildListString`, `BuildMapString`, `BuildQuery`, `BuildTracing`, `BuildLoadBalancerIpSet`, `FormatIpSetPath`, `FlattenQuery`, `FlattenTracing`, `FlattenLoadBalancerIpSet`), `EntityOperations` + `NewEntityOperations`.
- **Generic CRUD:** `CreateGeneric`, `ReadGeneric`, `UpdateGeneric`, `DeleteGeneric`.
- **Cross-resource shared shapes:** `BuildQuery`, `BuildTracing` (lightstep/otel/cpln), `BuildRouteHeaders`, `BuildRouteHeadersRequest`, `BuildRouteMirror`, plus their `Flatten*` counterparts. These are used by GVCs, workloads, domains, identities, etc., so live in `common.go`.
- **Shared schema fragments:** `QuerySchema`, `LightstepTracingSchema`, `OtelTracingSchema`, `ControlPlaneTracingSchema`, `CustomTagsTracingSchema`, `GetPortValidators`.

**Add to `common.go`** when:
- A builder, flattener, or schema fragment is (or would be) used by **two or more** resources/data sources.
- A field shape (like `query`, `tracing`, route mirrors, route headers) appears across multiple kinds in the API.
- You're tempted to copy/paste a builder between two resource files — extract it to `common.go` instead.

**Promote, don't duplicate.** If you find yourself adding the same builder to a second resource, move the existing one to `common.go` (or `helper.go` if it's type-level), then have both call sites use it. The cost of a small refactor now beats two-call-sites drift later.

### `internal/provider/models/common/` — shared models

When a model struct is referenced from `common.go`'s shared builders/flatteners (e.g. `LightstepTracingModel`, `OtelTracingModel`, `ControlPlaneTracingModel`, `QueryModel`), it lives under `internal/provider/models/common/`. Resource-local models stay under `internal/provider/models/<kind>/`.

### Writing-order checklist (before adding new code)

1. **Search `helper.go`** — is there already a `Build<X>` / `Flatten<X>` / pointer / formatter that does this?
2. **Search `common.go`** — is there already a shared builder/flattener/schema for this domain shape?
3. **Search `models/common/`** — is there a shared model already?
4. **Search the closest sibling resource** — is there a buildable analog you can mirror or extract?
5. Only then write something new. If it could ever be reused, put it in `common.go` (cross-resource) or `helper.go` (type-level) from day one.

## 12. Validators / modifiers — extend, don't fork

Existing validators (`internal/provider/validators/*.go`): `NameValidator`, `LinkValidator`, `DescriptionValidator`, `TagValidator`, `DisallowListValidator`, `DisallowPrefixValidator`, `PrefixStringValidator`. New validators follow the same shape: a struct with `Description`, `MarkdownDescription`, `Validate<Type>` methods.

Existing modifiers (`internal/provider/modifiers/*.go`): description, dictionary-as-envs, suppress-diff-on-equal-json, tag. Same skeleton when adding new ones.

## 13. Documentation (`docs/resources/<kind>.md`)

Every user-facing change to a resource requires a doc update. Template:

```markdown
---
page_title: "cpln_<kind> Resource - terraform-provider-cpln"
subcategory: "<Kind Title Case>"
description: |-
---

# cpln_<kind> (Resource)

<one-paragraph description, link to docs.controlplane.com>

## Declaration

### Required
- **name** (String) Name of the <Kind>.

### Optional
- **description** (String) ...
- **tags** (Map of String) Key-value map of resource tags.
- **<block>** (Block List, Max: 1) ([see below](#nestedblock--<block>))

## Outputs

The following attributes are exported:
- **cpln_id** (String) The ID, in GUID format, of the <Kind>.
- **self_link** (String) Full link to this resource. Can be referenced by other resources.

## Example Usage

```terraform
resource "cpln_<kind>" "example" {
  name = "example"
  ...
}
```

## Import Syntax

To update a statefile with an existing <kind> resource, execute the following import command:

```terraform
terraform import cpln_<kind>.RESOURCE_NAME <KIND>_NAME
```
```

Nested blocks get an `<a id="nestedblock--path--to--block"></a>` anchor and a `### \`path.to.block\`` header followed by `Required:` / `Optional:` lists, matching `docs/resources/mk8s.md` for the deepest schema example.

When you add a new optional field, also add it to **every** `## Example Usage` block that already configures the parent — pick a sensible default value (often the schema default).

## 14. Field / resource addition checklist

When you add a new field to an existing resource — or a whole new resource or data source — touch every layer in this order. Skipping any layer is the most common cause of half-finished work in this repo.

### Adding an optional field to an existing resource

1. **Client struct** (`internal/provider/client/<kind>.go`) — add the new field as `*T` with `json:"name,omitempty"`. For string enums append a trailing `// Enum: "a", "b"` comment (§8).
2. **Model** (`internal/provider/models/<kind>/<kind>.go`) — add the matching `types.X` field with `tfsdk:"snake_name"` and update the parent's `AttributeTypes()` map. (§7)
3. **Resource schema** (`internal/provider/resource_<kind>.go`) — add the schema attribute (prefer `SingleNestedAttribute` for new nested objects; see §6). If the parent block has a `default<X>Value()` helper, update that too.
4. **Builder & flattener** (`resource_<kind>.go` `// Builders //` and `// Flatteners //` sections) — read the value from the model into the client struct on the way in, and back to the model on the way out. Reuse helpers from `helper.go` / `common.go` first (§11). **If a builder/flattener already lives in `common.go` (e.g. `BuildRouteMirror`), edit it there — don't fork.**
5. **Resource tests** (`internal/provider/resource_<kind>_test.go`) — exercise the new field in **at least one** existing HCL config and update the matching `Expected*` checker / `TestCheckNestedBlocks` / `TestCheckResourceAttr` assertion. For Optional fields, leave at least one peer config without the field set so both the "set" and "unset" paths get coverage. No `// TODO` placeholders. (§10)
6. **Data source mirror** — if a `data_source_<kind>.go` exists, surface the field there too (with matching `data_source_<kind>_test.go` and `docs/data-sources/<kind>.md` updates).
7. **Docs** (`docs/resources/<kind>.md`, plus `docs/data-sources/<kind>.md` if applicable) — add the field to the Required/Optional list under the right nested block, and to **every** `## Example Usage` block that already configures the parent (pick a sensible default value). (§13)
8. **Verify** before each commit: `go build ./... && go vet ./... && go test -race ./internal/provider/ -count=1 -timeout=120s`.

### Adding a new resource or data source

In addition to running the field checklist above for **every** field, also:

- **Provider registry** (`internal/provider/provider.go`) — register the new `New<Kind>Resource` / `New<Kind>DataSource` in the `Resources()` / `DataSources()` slice.
- **New files** (resources): `internal/provider/client/<kind>.go`, `internal/provider/models/<kind>/<kind>.go`, `internal/provider/resource_<kind>.go`, `internal/provider/resource_<kind>_test.go`, `docs/resources/<kind>.md`.
- **New files** (data sources): `internal/provider/data_source_<kind>.go`, `internal/provider/data_source_<kind>_test.go`, `docs/data-sources/<kind>.md`. The client struct under `internal/provider/client/<kind>.go` is shared with the resource (don't duplicate it).
- **Templates**: start from `templates/resource_skeleton.txt` and `templates/resource_test_skeleton.txt`. Mirror the closest sibling resource for everything not in the skeleton.

### One-pass mental check before pushing

- [ ] **Client** — field added with correct `json` tag
- [ ] **Model** — field added; `AttributeTypes()` updated
- [ ] **Resource schema** — attribute exposed (with proper validators/defaults)
- [ ] **Builder** reads it; **flattener** writes it (in the right file — `resource_<kind>.go` or `common.go`)
- [ ] **Resource tests** — at least one HCL config sets it; the matching checker asserts it
- [ ] **Data source** — mirrored if applicable
- [ ] **Docs** — listed under the correct block; example usage shows it
- [ ] `go build`, `go vet`, and unit tests are green

## 15. Build & test commands

```bash
go build ./...                                              # full build
go vet ./...                                                # vet
go test -race ./internal/provider/ -count=1 -timeout=120s   # unit tests
make test                                                   # short test target
make testacc                                                # acceptance tests (requires CPLN_* envs, runs against real API)
make install                                                # build + install plugin to ~/.terraform.d
```

The CI workflow pre-installs the Terraform CLI before running acceptance tests.

## 16. Workflow expectations

- **Reuse first.** Before writing any utility, builder, flattener, validator, or schema fragment, check `helper.go`, `common.go`, `models/common/`, `validators/`, and `modifiers/`. If something analogous exists, use it. If you'd be copy/pasting between resources, extract to `common.go`. (Section 11.)
- **Mirror, don't invent.** Before writing a resource, read at least two existing resources of similar shape (e.g., for a CRUD-only resource read `resource_audit_context.go`; for a complex multi-provider resource read `resource_mk8s.go`).
- **Prefer object over list-as-block** for new single-nested attributes. `schema.SingleNestedAttribute` + `types.Object` is the modern shape; reach for `ListNestedBlock` only when matching legacy code in the same resource. (Section 6.)
- **Templates are starting points, not stop-points.** `templates/resource_skeleton.txt` and `templates/resource_test_skeleton.txt` give the bones; add Builders/Flatteners/Schemas sections as needed by reading neighboring files.
- **Comment density is non-negotiable.** Every function carries a one-line description, every step inside it carries a comment. This is the project's norm — do not "trim".
- **Plural surgery, plural fix.** If you add an attribute to one place, find every analog (HCL configs, test checkers, defaults helper, doc examples, etc.) and update them all. Half-finished additions are the most common style break in this repo (see prior `// TODO: Add byok test here` placeholders). Use the §14 checklist as the canonical list of files to touch.
- **Joi → Go primitive types** are tracked via the upstream API schema. When in doubt, search the corresponding `controlplane/nodelibs/schema/src/*.ts` file for the Joi definition and translate per Section 7.
- **No emojis in code, comments, commits, or docs unless the user asks.**
- **Never delete or rewrite working code as a shortcut.** If a hook or test fails, fix the root cause.

## 17. Commit / branch conventions (set globally, repeated here)

- Branch: `majid/<3-4 kebab-case words>`.
- Commit message: one capitalized present-tense line, no body, no type prefix.
  - `Add juicefs to mk8s BYOK addon configuration`
  - `Make cpln_domain TLS optional`

## 18. When unsure

1. Search `helper.go` for an existing type-level utility (Build/Flatten/pointer/formatter).
2. Search `common.go` for an existing shared builder/flattener/schema/CRUD primitive.
3. Search `internal/provider/models/common/` for an existing shared model.
4. Open `templates/SNIPPETS.md` and the matching `templates/<list|set|known_list>/*.txt`.
5. Read the closest sibling resource in `internal/provider/resource_<kind>.go`.
6. Read its model file `internal/provider/models/<kind>/<kind>.go`.
7. Read its test file `internal/provider/resource_<kind>_test.go`.
8. Mirror.

If a pattern is genuinely missing for what you need, ask — don't invent.
