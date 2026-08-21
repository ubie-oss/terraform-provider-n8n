package provider

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/controllers"
)

func TestAccN8nProject_basic(t *testing.T) {
	if !isIntegrationTestMode() {
		t.Skip("Skipping acceptance test unless TF_ACC=1")
	}

	createName := acctest.RandomWithPrefix("tf-n8n-proj")
	updateName := createName + "-upd"

	createConfig, err := ReadAccTestResource([]string{"resources", "n8n_project", "010_create.tf"})
	if err != nil {
		t.Fatalf("read create fixture: %v", err)
	}
	updateConfig, err := ReadAccTestResource([]string{"resources", "n8n_project", "020_update.tf"})
	if err != nil {
		t.Fatalf("read update fixture: %v", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckProjects(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: getProviderConfig() + replaceAccName(createConfig, createName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_project.test", "name", createName),
					resource.TestCheckResourceAttr("n8n_project.test", "type", "team"),
					resource.TestCheckResourceAttr("n8n_project.test", "delete_protection", "false"),
					resource.TestCheckResourceAttrSet("n8n_project.test", "id"),
				),
			},
			{
				ResourceName:            "n8n_project.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_protection"},
			},
			{
				Config: getProviderConfig() + replaceAccName(updateConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_project.test", "name", updateName),
					resource.TestCheckResourceAttr("n8n_project.test", "type", "team"),
					resource.TestCheckResourceAttr("n8n_project.test", "delete_protection", "false"),
				),
			},
		},
	})
}

func TestAccN8nProject_dataSources(t *testing.T) {
	if !isIntegrationTestMode() {
		t.Skip("Skipping acceptance test unless TF_ACC=1")
	}

	name := acctest.RandomWithPrefix("tf-n8n-proj-ds")
	fixture, err := ReadAccTestResource([]string{"data_sources", "n8n_project", "010_data.tf"})
	if err != nil {
		t.Fatalf("read data fixture: %v", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckProjects(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: getProviderConfig() + replaceAccName(fixture, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("n8n_project.test", "id", "data.n8n_project.by_id", "id"),
					resource.TestCheckResourceAttrPair("n8n_project.test", "name", "data.n8n_project.by_id", "name"),
					resource.TestCheckResourceAttrPair("n8n_project.test", "type", "data.n8n_project.by_id", "type"),
					resource.TestCheckResourceAttr("data.n8n_projects.all", "id", "projects"),
					resource.TestCheckResourceAttrSet("data.n8n_projects.all", "projects.#"),
				),
			},
		},
	})
}

func TestAccN8nProject_deleteProtection(t *testing.T) {
	if !isIntegrationTestMode() {
		t.Skip("Skipping acceptance test unless TF_ACC=1")
	}

	name := acctest.RandomWithPrefix("tf-n8n-proj-dp")

	protectConfig, err := ReadAccTestResource([]string{"resources", "n8n_project", "030_delete_protection.tf"})
	if err != nil {
		t.Fatalf("read protect fixture: %v", err)
	}
	allowDestroyConfig, err := ReadAccTestResource([]string{"resources", "n8n_project", "040_allow_destroy.tf"})
	if err != nil {
		t.Fatalf("read allow-destroy fixture: %v", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckProjects(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: getProviderConfig() + replaceAccName(protectConfig, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_project.test", "delete_protection", "true"),
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
					resource.TestCheckResourceAttr("n8n_project.test", "delete_protection", "false"),
				),
			},
		},
	})
}

func testAccPreCheckProjects(t *testing.T) {
	t.Helper()
	testAccPreCheck(t)

	client, err := n8n.New(os.Getenv("N8N_ENDPOINT"), os.Getenv("N8N_API_KEY"), nil)
	if err != nil {
		t.Fatalf("configure n8n client: %v", err)
	}
	_, err = controllers.NewProjectController(client).List(context.Background())
	if n8n.IsProjectRoleAdminUnlicensed(err) {
		t.Skip("n8n instance is not licensed for feat:projectRole:admin; skip project acceptance tests")
	}
	if err != nil {
		t.Fatalf("probe GET /projects: %v", err)
	}
}

func testAccCheckProjectDestroy(s *terraform.State) error {
	client, err := n8n.New(os.Getenv("N8N_ENDPOINT"), os.Getenv("N8N_API_KEY"), nil)
	if err != nil {
		return err
	}
	ctrl := controllers.NewProjectController(client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "n8n_project" {
			continue
		}
		_, err := ctrl.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("project %q still exists", rs.Primary.ID)
		}
		if !errors.Is(err, controllers.ErrProjectNotFound) {
			return fmt.Errorf("checking destroy of project %q: %w", rs.Primary.ID, err)
		}
	}
	return nil
}

func replaceAccName(config, name string) string {
	return strings.ReplaceAll(config, "{{NAME}}", name)
}
