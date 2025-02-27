---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "openfga_authorization_models Data Source - openfga"
subcategory: ""
description: |-
  Provides the ability to list and retrieve details of existing authorization models in a specific store.
---

# openfga_authorization_models (Data Source)

Provides the ability to list and retrieve details of existing authorization models in a specific store.

## Example Usage

```terraform
data "openfga_authorization_models" "example" {
  store_id = "example_store_id"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `store_id` (String) The unique ID of the store to list authorization models for.

### Read-Only

- `authorization_models` (Attributes List) List of existing authorization models in the specific store. (see [below for nested schema](#nestedatt--authorization_models))

<a id="nestedatt--authorization_models"></a>
### Nested Schema for `authorization_models`

Read-Only:

- `id` (String) The unique ID of the authorization model.
- `model_json` (String) The authorization model definition in JSON format.
