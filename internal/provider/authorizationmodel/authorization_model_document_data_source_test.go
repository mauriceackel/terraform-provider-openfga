package authorizationmodel_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/openfga/terraform-provider-openfga/internal/provider/acceptance"
)

func TestAccAuthorizationModelDocumentDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test DSL
			{
				Config: testAccAuthorizationModelDocumentDataSourceConfigDsl(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.openfga_authorization_model_document.test",
						tfjsonpath.New("result"),
						knownvalue.StringExact(expectedAuthorizationModelDocumentDataSourceResult),
					),
				},
			},
			// Test mod file
			{
				Config: testAccAuthorizationModelDocumentDataSourceConfigModFile(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.openfga_authorization_model_document.test",
						tfjsonpath.New("result"),
						knownvalue.StringExact(expectedModularAuthorizationModelDocumentDataSourceResult),
					),
				},
			},
			// Test JSON
			{
				Config: testAccAuthorizationModelDocumentDataSourceConfigJson(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.openfga_authorization_model_document.test",
						tfjsonpath.New("result"),
						knownvalue.StringExact(expectedAuthorizationModelDocumentDataSourceResult),
					),
				},
			},
			// Test model
			{
				Config: testAccAuthorizationModelDocumentDataSourceConfigModel(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.openfga_authorization_model_document.test",
						tfjsonpath.New("result"),
						knownvalue.StringExact(expectedAuthorizationModelDocumentDataSourceResult),
					),
				},
			},
		},
	})
}

const expectedAuthorizationModelDocumentDataSourceResult = `{"conditions":{"larger_than":{"expression":"a \u003e b","name":"larger_than","parameters":{"a":{"generic_types":[],"type_name":"TYPE_NAME_INT"},"b":{"generic_types":[],"type_name":"TYPE_NAME_INT"}}}},"schema_version":"1.1","type_definitions":[{"relations":{},"type":"user"},{"metadata":{"module":"","relations":{"viewer":{"directly_related_user_types":[{"condition":"","type":"user"}],"module":""}}},"relations":{"viewer":{"this":{}}},"type":"document"}]}`
const expectedModularAuthorizationModelDocumentDataSourceResult = `{"conditions":{"larger_than":{"expression":"a \u003e b","metadata":{"module":"conditions","source_info":{"file":"conditions/larger_than.fga"}},"name":"larger_than","parameters":{"a":{"generic_types":[],"type_name":"TYPE_NAME_INT"},"b":{"generic_types":[],"type_name":"TYPE_NAME_INT"}}}},"schema_version":"1.2","type_definitions":[{"metadata":{"module":"user","relations":{},"source_info":{"file":"user.fga"}},"relations":{},"type":"user"},{"metadata":{"module":"document","relations":{"viewer":{"directly_related_user_types":[{"condition":"","type":"user"}],"module":""}},"source_info":{"file":"document.fga"}},"relations":{"viewer":{"this":{}}},"type":"document"}]}`

const authorizationModelResource = `
resource "openfga_store" "test" {
	name = "test"
}

resource "openfga_authorization_model" "test" {
  store_id = openfga_store.test.id

  model_json = data.openfga_authorization_model_document.test.result
}`

func testAccAuthorizationModelDocumentDataSourceConfigDsl() string {
	return fmt.Sprintf(`
%[1]s

data "openfga_authorization_model_document" "test" {
	dsl = <<EOT
model
	schema 1.1

type user

type document
	relations
		define viewer: [user]

condition larger_than(a: int, b: int) {
	a > b
}
	EOT
}

%[2]s
`, acceptance.ProviderConfig, authorizationModelResource)
}

func testAccAuthorizationModelDocumentDataSourceConfigModFile() string {
	return fmt.Sprintf(`
%[1]s

data "openfga_authorization_model_document" "test" {
	mod_file_path = "${path.root}/../acceptance/modularmodel/fga.mod"
}

%[2]s
`, acceptance.ProviderConfig, authorizationModelResource)
}

func testAccAuthorizationModelDocumentDataSourceConfigJson() string {
	return fmt.Sprintf(`
%[1]s

data "openfga_authorization_model_document" "test" {
	json = <<EOT
{
	"type_definitions": [
		{ "type": "user" },
		{
			"type": "document",
			"relations": {
				"viewer": {
					"this": {}
				}
			},
			"metadata": {
				"relations":{"viewer":{"directly_related_user_types":[{"type":"user"}]}}
			}
    	}
	],
	"schema_version": "1.1",
	"conditions":{"larger_than":{"expression":"a > b","name":"larger_than","parameters":{"a":{"type_name":"TYPE_NAME_INT"},"b":{"type_name":"TYPE_NAME_INT"}}}}
}
	EOT
}

%[2]s
`, acceptance.ProviderConfig, authorizationModelResource)
}

func testAccAuthorizationModelDocumentDataSourceConfigModel() string {
	return fmt.Sprintf(`
%[1]s

data "openfga_authorization_model_document" "test" {
	model = {
		schema_version   = "1.1"
		type_definitions = [
			{
				type = "user"
			},
			{
				type      = "document"
				relations = {
					viewer = {
						this = {}
					}
				}
				metadata  = {
					relations = {
						viewer = {
							directly_related_user_types = [
								{ type = "user" }
							]
						}
					}
				}
			},
		]
		conditions       = {
			larger_than = {
				name       = "larger_than"
				expression = "a > b"
				parameters = {
					a = {
						type_name = "TYPE_NAME_INT"
					}
					b = {
						type_name = "TYPE_NAME_INT"
					} 
				}
			}
		}
	}
}

%[2]s
`, acceptance.ProviderConfig, authorizationModelResource)
}
