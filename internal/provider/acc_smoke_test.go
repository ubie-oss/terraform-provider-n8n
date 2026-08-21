package provider

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
)

// workflowListProbe is the CE-available GET /workflows envelope (data + optional cursor).
type workflowListProbe struct {
	Data []any `json:"data"`
}

func TestAccN8n_publicAPIKeyWorks(t *testing.T) {
	if !isIntegrationTestMode() {
		t.Skip("Skipping acceptance test unless TF_ACC=1")
	}
	testAccPreCheck(t)

	client, err := n8n.New(os.Getenv("N8N_ENDPOINT"), os.Getenv("N8N_API_KEY"), nil)
	if err != nil {
		t.Fatalf("configure n8n client: %v", err)
	}

	var out workflowListProbe
	err = client.DoJSON(context.Background(), http.MethodGet, client.URL("workflows")+"?limit=1", nil, &out, "workflow", "")
	if err != nil {
		t.Fatalf("GET /api/v1/workflows with bootstrapped API key: %v", err)
	}
}
