// Copyright 2026 Ubie, inc.
// Boilerplate: copy patterns into internal/provider when adding acceptance tests.
package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccExampleBasic = `
provider "n8n" {
  # endpoint and api_key from N8N_ENDPOINT / N8N_API_KEY
}

resource "n8n_project" "test" {
  name = "acceptance-example"
}
`

func TestAccExample_basic(t *testing.T) {
	if !isIntegrationTestMode() {
		t.Skip("Skipping acceptance test unless TF_ACC=1")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccExampleBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_project.test", "name", "acceptance-example"),
				),
			},
			{
				ResourceName:      "n8n_project.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
