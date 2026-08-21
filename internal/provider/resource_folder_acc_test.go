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

func TestAccN8nFolder_basic(t *testing.T) {
	if !isIntegrationTestMode() {
		t.Skip("Skipping acceptance test unless TF_ACC=1")
	}

	name := acctest.RandomWithPrefix("tf-n8n-fld")

	createConfig, err := ReadAccTestResource([]string{"resources", "n8n_folder", "010_create.tf"})
	if err != nil {
		t.Fatalf("read create fixture: %v", err)
	}
	updateConfig, err := ReadAccTestResource([]string{"resources", "n8n_folder", "020_update.tf"})
	if err != nil {
		t.Fatalf("read update fixture: %v", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckFolders(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckFolderDestroy,
		Steps: []resource.TestStep{
			{
				Config: getProviderConfig() + replaceAccName(createConfig, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_folder.test", "name", name+"-root"),
					resource.TestCheckResourceAttr("n8n_folder.test", "delete_protection", "false"),
					resource.TestCheckResourceAttrSet("n8n_folder.test", "id"),
					resource.TestCheckResourceAttrSet("n8n_folder.test", "folder_id"),
					resource.TestCheckResourceAttrPair("n8n_folder.child", "parent_folder_id", "n8n_folder.test", "folder_id"),
					resource.TestCheckResourceAttr("n8n_folder.child", "name", name+"-child"),
				),
			},
			{
				ResourceName:            "n8n_folder.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_protection", "transfer_to_folder_id"},
			},
			{
				Config: getProviderConfig() + replaceAccName(updateConfig, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_folder.test", "name", name+"-renamed"),
					resource.TestCheckResourceAttrPair("n8n_folder.child", "parent_folder_id", "n8n_folder.test", "folder_id"),
				),
			},
		},
	})
}

func TestAccN8nFolder_dataSources(t *testing.T) {
	if !isIntegrationTestMode() {
		t.Skip("Skipping acceptance test unless TF_ACC=1")
	}

	name := acctest.RandomWithPrefix("tf-n8n-fld-ds")
	fixture, err := ReadAccTestResource([]string{"data_sources", "n8n_folder", "010_data.tf"})
	if err != nil {
		t.Fatalf("read data fixture: %v", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckFolders(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckFolderDestroy,
		Steps: []resource.TestStep{
			{
				Config: getProviderConfig() + replaceAccName(fixture, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("n8n_folder.test", "folder_id", "data.n8n_folder.by_id", "folder_id"),
					resource.TestCheckResourceAttrPair("n8n_folder.test", "name", "data.n8n_folder.by_id", "name"),
					resource.TestCheckResourceAttrSet("data.n8n_folders.all", "id"),
					resource.TestCheckResourceAttrSet("data.n8n_folders.all", "folders.#"),
				),
			},
		},
	})
}

func TestAccN8nFolder_deleteProtection(t *testing.T) {
	if !isIntegrationTestMode() {
		t.Skip("Skipping acceptance test unless TF_ACC=1")
	}

	name := acctest.RandomWithPrefix("tf-n8n-fld-dp")

	protectConfig, err := ReadAccTestResource([]string{"resources", "n8n_folder", "030_delete_protection.tf"})
	if err != nil {
		t.Fatalf("read protect fixture: %v", err)
	}
	allowDestroyConfig, err := ReadAccTestResource([]string{"resources", "n8n_folder", "040_allow_destroy.tf"})
	if err != nil {
		t.Fatalf("read allow-destroy fixture: %v", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckFolders(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckFolderDestroy,
		Steps: []resource.TestStep{
			{
				Config: getProviderConfig() + replaceAccName(protectConfig, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_folder.test", "delete_protection", "true"),
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
					resource.TestCheckResourceAttr("n8n_folder.test", "delete_protection", "false"),
				),
			},
		},
	})
}

func testAccPreCheckFolders(t *testing.T) {
	t.Helper()
	testAccPreCheck(t)

	client, err := n8n.New(os.Getenv("N8N_ENDPOINT"), os.Getenv("N8N_API_KEY"), nil)
	if err != nil {
		t.Fatalf("configure n8n client: %v", err)
	}
	_, err = controllers.NewFolderController(client).List(context.Background(), controllers.ListFoldersOptions{ProjectID: "personal"})
	if n8n.IsFoldersUnlicensed(err) {
		t.Skip("n8n instance is not licensed for feat:folders; skip folder acceptance tests")
	}
	if err != nil {
		t.Fatalf("probe GET /projects/personal/folders: %v", err)
	}

	_, err = controllers.NewProjectController(client).List(context.Background())
	if n8n.IsProjectRoleAdminUnlicensed(err) {
		t.Skip("n8n instance is not licensed for feat:projectRole:admin; skip folder acceptance tests")
	}
	if err != nil {
		t.Fatalf("probe GET /projects: %v", err)
	}
}

func testAccCheckFolderDestroy(s *terraform.State) error {
	client, err := n8n.New(os.Getenv("N8N_ENDPOINT"), os.Getenv("N8N_API_KEY"), nil)
	if err != nil {
		return err
	}
	folderCtrl := controllers.NewFolderController(client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "n8n_folder" {
			continue
		}
		_, err := folderCtrl.Get(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["folder_id"])
		if err == nil {
			return fmt.Errorf("folder %q still exists", rs.Primary.ID)
		}
		if !n8n.IsNotFound(err) {
			return fmt.Errorf("checking destroy of folder %q: %w", rs.Primary.ID, err)
		}
	}
	return testAccCheckProjectDestroy(s)
}
