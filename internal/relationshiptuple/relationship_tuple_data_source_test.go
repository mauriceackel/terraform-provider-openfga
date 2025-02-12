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

func TestAccRelationshipTupleDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccRelationshipTupleDataSourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.openfga_relationship_tuple.test",
						tfjsonpath.New("user"),
						knownvalue.StringExact("user:user-1"),
					),
					statecheck.ExpectKnownValue(
						"data.openfga_relationship_tuple.test",
						tfjsonpath.New("relation"),
						knownvalue.StringExact("viewer"),
					),
					statecheck.ExpectKnownValue(
						"data.openfga_relationship_tuple.test",
						tfjsonpath.New("object"),
						knownvalue.StringExact("document:dummy"),
					),
				},
			},
		},
	})
}

func testAccRelationshipTupleDataSourceConfig() string {
	return fmt.Sprintf(`
%[1]s

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

resource "openfga_relationship_tuple" "test" {
  store_id = openfga_store.test.id

  user     = "user:user-1"
  relation = "viewer"
  object   = "document:dummy"

  depends_on = [openfga_authorization_model.test]
}

data "openfga_relationship_tuple" "test" {
  store_id = openfga_store.test.id

  user     = "user:user-1"
  relation = "viewer"
  object   = "document:dummy"

  depends_on = [openfga_relationship_tuple.test]
}
`, acceptance.ProviderConfig)
}
