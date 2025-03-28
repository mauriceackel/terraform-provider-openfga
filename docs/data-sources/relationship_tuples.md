---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "openfga_relationship_tuples Data Source - openfga"
subcategory: ""
description: |-
  Provides the ability to list and retrieve details of existing relationship tuples in a specific store.
---

# openfga_relationship_tuples (Data Source)

Provides the ability to list and retrieve details of existing relationship tuples in a specific store.

## Example Usage

```terraform
data "openfga_relationship_tuples" "all" {
  store_id = "example_store_id"
}

data "openfga_relationship_tuples" "query" {
  store_id = "example_store_id"

  query = {
    user     = "user:user-1"
    relation = "viewer"
    object   = "document:"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `store_id` (String) The unique ID of the store to list relationship tuples for.

### Optional

- `query` (Attributes) A query to filter the returned relationship tuples. Can be left blank to retrieve all relationship tuples. (see [below for nested schema](#nestedatt--query))

### Read-Only

- `relationship_tuples` (Attributes List) List of existing relationship tuples in the specific store, matching the query. (see [below for nested schema](#nestedatt--relationship_tuples))

<a id="nestedatt--query"></a>
### Nested Schema for `query`

Required:

- `object` (String) The object of the resulting relationship tuples.

Optional:

- `relation` (String) The relation of the resulting relationship tuples.
- `user` (String) The user of the resulting relationship tuples.


<a id="nestedatt--relationship_tuples"></a>
### Nested Schema for `relationship_tuples`

Read-Only:

- `condition` (Attributes) A condition of the relationship tuple. (see [below for nested schema](#nestedatt--relationship_tuples--condition))
- `object` (String) The object of the relationship tuple.
- `relation` (String) The relation of the relationship tuple.
- `user` (String) The user of the relationship tuple.

<a id="nestedatt--relationship_tuples--condition"></a>
### Nested Schema for `relationship_tuples.condition`

Read-Only:

- `context_json` (String) The (partial) context under which the condition is evaluated.
- `name` (String) The name of the condition.
