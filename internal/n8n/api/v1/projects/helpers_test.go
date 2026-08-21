package projects

import (
	"net/http"
	"testing"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/api/apitest"
)

func testClient(t *testing.T, handler http.HandlerFunc) *n8n.Client {
	return apitest.TestClient(t, handler)
}

func writeFixture(t *testing.T, w http.ResponseWriter, status int, name string) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(apitest.Testdata(t, name)); err != nil {
		t.Errorf("write fixture: %v", err)
	}
}

func decodeBody(t *testing.T, r *http.Request, out any) {
	apitest.DecodeBody(t, r, out)
}

func assertAPIKey(t *testing.T, r *http.Request) {
	apitest.AssertAPIKey(t, r)
}
