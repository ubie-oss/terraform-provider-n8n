package provider

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/controllers"
)

func TestAccN8nTag_basic(t *testing.T) {
	if !isIntegrationTestMode() {
		t.Skip("Skipping acceptance test unless TF_ACC=1")
	}

	// n8n tag names are limited to 24 characters (longer names return a misleading 409).
	createName := fmt.Sprintf("t-%s", acctest.RandString(8))
	updateName := createName + "-u"

	createConfig, err := ReadAccTestResource([]string{"resources", "n8n_tag", "010_create.tf"})
	if err != nil {
		t.Fatalf("read create fixture: %v", err)
	}
	updateConfig, err := ReadAccTestResource([]string{"resources", "n8n_tag", "020_update.tf"})
	if err != nil {
		t.Fatalf("read update fixture: %v", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckTags(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: getProviderConfig() + replaceAccName(createConfig, createName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_tag.test", "name", createName),
					resource.TestCheckResourceAttr("n8n_tag.test", "delete_protection", "false"),
					resource.TestCheckResourceAttrSet("n8n_tag.test", "id"),
					resource.TestCheckResourceAttrSet("n8n_tag.test", "created_at"),
					resource.TestCheckResourceAttrSet("n8n_tag.test", "updated_at"),
				),
			},
			{
				ResourceName:            "n8n_tag.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_protection"},
			},
			{
				Config: getProviderConfig() + replaceAccName(updateConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_tag.test", "name", updateName),
					resource.TestCheckResourceAttr("n8n_tag.test", "delete_protection", "false"),
				),
			},
		},
	})
}

func TestAccN8nTag_dataSources(t *testing.T) {
	if !isIntegrationTestMode() {
		t.Skip("Skipping acceptance test unless TF_ACC=1")
	}

	name := fmt.Sprintf("t-%s", acctest.RandString(8))
	fixture, err := ReadAccTestResource([]string{"data_sources", "n8n_tag", "010_data.tf"})
	if err != nil {
		t.Fatalf("read data fixture: %v", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckTags(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: getProviderConfig() + replaceAccName(fixture, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("n8n_tag.test", "id", "data.n8n_tag.by_id", "id"),
					resource.TestCheckResourceAttrPair("n8n_tag.test", "name", "data.n8n_tag.by_id", "name"),
					resource.TestCheckResourceAttr("data.n8n_tags.all", "id", "tags"),
					resource.TestCheckResourceAttrSet("data.n8n_tags.all", "tags.#"),
				),
			},
		},
	})
}

func TestAccN8nTag_deleteProtection(t *testing.T) {
	if !isIntegrationTestMode() {
		t.Skip("Skipping acceptance test unless TF_ACC=1")
	}

	name := fmt.Sprintf("t-%s", acctest.RandString(8))

	protectConfig, err := ReadAccTestResource([]string{"resources", "n8n_tag", "030_delete_protection.tf"})
	if err != nil {
		t.Fatalf("read protect fixture: %v", err)
	}
	allowDestroyConfig, err := ReadAccTestResource([]string{"resources", "n8n_tag", "040_allow_destroy.tf"})
	if err != nil {
		t.Fatalf("read allow-destroy fixture: %v", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckTags(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: getProviderConfig() + replaceAccName(protectConfig, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_tag.test", "delete_protection", "true"),
				),
			},
			{
				Config:      getProviderConfig() + replaceAccName(protectConfig, name),
				Destroy:     true,
				ExpectError: regexp.MustCompile("delete protection is enabled"),
			},
			{
				Config: getProviderConfig() + replaceAccName(allowDestroyConfig, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_tag.test", "delete_protection", "false"),
				),
			},
		},
	})
}

func testAccPreCheckTags(t *testing.T) {
	t.Helper()
	testAccPreCheck(t)

	client, err := n8n.New(os.Getenv("N8N_ENDPOINT"), os.Getenv("N8N_API_KEY"), nil)
	if err != nil {
		t.Fatalf("configure n8n client: %v", err)
	}
	_, err = controllers.NewTagController(client).List(context.Background())
	if err != nil {
		t.Fatalf("probe GET /tags: %v", err)
	}
}

func testAccCheckTagDestroy(s *terraform.State) error {
	client, err := n8n.New(os.Getenv("N8N_ENDPOINT"), os.Getenv("N8N_API_KEY"), nil)
	if err != nil {
		return err
	}
	ctrl := controllers.NewTagController(client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "n8n_tag" {
			continue
		}
		_, err := ctrl.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("tag %q still exists", rs.Primary.ID)
		}
		if !n8n.IsNotFound(err) {
			return fmt.Errorf("checking destroy of tag %q: %w", rs.Primary.ID, err)
		}
	}
	return nil
}
