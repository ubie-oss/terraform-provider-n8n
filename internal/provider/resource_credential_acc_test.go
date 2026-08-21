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
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/controllers"
)

func TestAccN8nCredential_basic(t *testing.T) {
	if !isIntegrationTestMode() {
		t.Skip("Skipping acceptance test unless TF_ACC=1")
	}

	name := acctest.RandomWithPrefix("tf-n8n-cred")

	createConfig, err := ReadAccTestResource([]string{"resources", "n8n_credential", "010_create.tf"})
	if err != nil {
		t.Fatalf("read create fixture: %v", err)
	}
	updateConfig, err := ReadAccTestResource([]string{"resources", "n8n_credential", "020_update.tf"})
	if err != nil {
		t.Fatalf("read update fixture: %v", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckCredentials(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_11_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: getProviderConfig() + replaceAccName(createConfig, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_credential.test", "name", name),
					resource.TestCheckResourceAttr("n8n_credential.test", "type", "httpHeaderAuth"),
					resource.TestCheckResourceAttr("n8n_credential.test", "delete_protection", "false"),
					resource.TestCheckResourceAttr("n8n_credential.test", "data_version", "1"),
					resource.TestCheckResourceAttrSet("n8n_credential.test", "id"),
					resource.TestCheckNoResourceAttr("n8n_credential.test", "data"),
				),
			},
			{
				ResourceName:            "n8n_credential.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"data", "data_version", "is_partial_data", "delete_protection"},
			},
			{
				Config: getProviderConfig() + replaceAccName(updateConfig, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_credential.test", "name", name+"-renamed"),
					resource.TestCheckResourceAttr("n8n_credential.test", "data_version", "2"),
					resource.TestCheckNoResourceAttr("n8n_credential.test", "data"),
				),
			},
		},
	})
}

func TestAccN8nCredential_dataSources(t *testing.T) {
	if !isIntegrationTestMode() {
		t.Skip("Skipping acceptance test unless TF_ACC=1")
	}

	name := acctest.RandomWithPrefix("tf-n8n-cred-ds")
	fixture, err := ReadAccTestResource([]string{"data_sources", "n8n_credential", "010_data.tf"})
	if err != nil {
		t.Fatalf("read data fixture: %v", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckCredentials(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_11_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: getProviderConfig() + replaceAccName(fixture, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("n8n_credential.test", "id", "data.n8n_credential.by_id", "id"),
					resource.TestCheckResourceAttrPair("n8n_credential.test", "name", "data.n8n_credential.by_id", "name"),
					resource.TestCheckResourceAttrSet("data.n8n_credentials.all", "id"),
					resource.TestCheckResourceAttrSet("data.n8n_credentials.all", "credentials.#"),
					resource.TestCheckResourceAttrSet("data.n8n_credential_schema.header", "schema_json"),
				),
			},
		},
	})
}

func TestAccN8nCredential_deleteProtection(t *testing.T) {
	if !isIntegrationTestMode() {
		t.Skip("Skipping acceptance test unless TF_ACC=1")
	}

	name := acctest.RandomWithPrefix("tf-n8n-cred-dp")

	protectConfig, err := ReadAccTestResource([]string{"resources", "n8n_credential", "030_delete_protection.tf"})
	if err != nil {
		t.Fatalf("read protect fixture: %v", err)
	}
	allowDestroyConfig, err := ReadAccTestResource([]string{"resources", "n8n_credential", "040_allow_destroy.tf"})
	if err != nil {
		t.Fatalf("read allow-destroy fixture: %v", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckCredentials(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_11_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: getProviderConfig() + replaceAccName(protectConfig, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_credential.test", "delete_protection", "true"),
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
					resource.TestCheckResourceAttr("n8n_credential.test", "delete_protection", "false"),
				),
			},
		},
	})
}

func testAccPreCheckCredentials(t *testing.T) {
	t.Helper()
	testAccPreCheck(t)

	client, err := n8n.New(os.Getenv("N8N_ENDPOINT"), os.Getenv("N8N_API_KEY"), nil)
	if err != nil {
		t.Fatalf("configure n8n client: %v", err)
	}
	_, err = controllers.NewCredentialController(client).List(context.Background(), controllers.ListCredentialsOptions{})
	if err != nil {
		t.Fatalf("probe GET /credentials: %v", err)
	}
}

func testAccCheckCredentialDestroy(s *terraform.State) error {
	client, err := n8n.New(os.Getenv("N8N_ENDPOINT"), os.Getenv("N8N_API_KEY"), nil)
	if err != nil {
		return err
	}
	ctrl := controllers.NewCredentialController(client)
	for _, rs := range s.RootModule().Resources {
		if !isCredentialManagedResource(rs.Type) {
			continue
		}
		_, err := ctrl.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("credential %q still exists", rs.Primary.ID)
		}
		if !n8n.IsNotFound(err) {
			return fmt.Errorf("checking destroy of credential %q: %w", rs.Primary.ID, err)
		}
	}
	return nil
}
