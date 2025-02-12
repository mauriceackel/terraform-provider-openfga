package relationshiptuple_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/mauriceackel/terraform-provider-openfga/internal/acceptance"
)

func TestAccRelationshipTuplesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test Empty
			{
				Config: testAccRelationshipTuplesDataSourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.openfga_relationship_tuples.all",
						tfjsonpath.New("relationship_tuples"),
						knownvalue.ListSizeExact(0),
					),
				},
			},
			// Setup relationship tuples
			{
				Config: testAccRelationshipTuplesDataSourceConfig(
					"document-1",
					"document-2",
				),
			},
			// Read testing
			{
				Config: testAccRelationshipTuplesDataSourceConfig(
					"document-1",
					"document-2",
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.openfga_relationship_tuples.all",
						tfjsonpath.New("relationship_tuples"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"store_id": knownvalue.NotNull(),
								"user":     knownvalue.StringExact("user:user-1"),
								"relation": knownvalue.StringExact("viewer"),
								"object":   knownvalue.StringExact("document:document-1"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"store_id": knownvalue.NotNull(),
								"user":     knownvalue.StringExact("user:user-1"),
								"relation": knownvalue.StringExact("viewer"),
								"object":   knownvalue.StringExact("document:document-2"),
							}),
						}),
					),
				},
			},
			// Query testing
			{
				Config: testAccRelationshipTuplesDataSourceConfig(
					"document-1",
					"document-2",
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.openfga_relationship_tuples.query",
						tfjsonpath.New("relationship_tuples"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"store_id": knownvalue.NotNull(),
								"user":     knownvalue.StringExact("user:user-1"),
								"relation": knownvalue.StringExact("viewer"),
								"object":   knownvalue.StringExact("document:document-1"),
							}),
						}),
					),
				},
			},
		},
	})
}

func testAccRelationshipTuplesDataSourceConfig(objectNames ...string) string {
	var resources string = `
resource "openfga_store" "test" {
  name = "test"
}

data "openfga_authorization_model_document" "test" {
  dsl = <<EOT
model
  schema 1.1

type user

type document
  relations
    define viewer: [user]
  EOT
}

resource "openfga_authorization_model" "test" {
  store_id = openfga_store.test.id

  model_json = data.openfga_authorization_model_document.test.result
}
	`

	for idx, objectName := range objectNames {
		var dependsOn string
		if idx > 0 {
			dependsOn = fmt.Sprintf(`depends_on = [openfga_relationship_tuple.tuple_%[1]d]`, idx-1)
		}

		resources += fmt.Sprintf(`
resource "openfga_relationship_tuple" "tuple_%[1]d" {
  store_id = openfga_store.test.id

  user     = "user:user-1"
  relation = "viewer"
  object   = "document:%[2]s"
  %[3]s
}
`, idx, objectName, dependsOn)
	}

	return fmt.Sprintf(`
%[1]s

%[2]s

data "openfga_relationship_tuples" "all" {
  store_id = openfga_store.test.id
}

data "openfga_relationship_tuples" "query" {
  store_id = openfga_store.test.id

  query = {
    object = "document:document-1" 
  }
}
`, acceptance.ProviderConfig, resources)
}
