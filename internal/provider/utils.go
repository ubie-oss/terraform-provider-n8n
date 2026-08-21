// Copyright 2026 Ubie, inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"context"
	"embed"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

//go:embed docs/**/*.md
var embeddedDocs embed.FS

const (
	integrationTestModeEnvVar = "TF_ACC"
)

func isIntegrationTestMode() bool {
	return os.Getenv(integrationTestModeEnvVar) == "1"
}

func getPathToAccTests() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get current file path")
	}
	accTestsPath := path.Join(path.Dir(filename), "acc_tests")
	if _, err := os.Stat(accTestsPath); os.IsNotExist(err) {
		return "", fmt.Errorf("acc_tests directory does not exist at %s", accTestsPath)
	}
	return accTestsPath, nil
}

func getPathToAccTestResource(elements []string) (string, error) {
	pathToAccTests, err := getPathToAccTests()
	if err != nil {
		return "", err
	}
	allElements := append([]string{pathToAccTests}, elements...)
	accTestResourcePath := path.Join(allElements...)
	cleanedAccTestsPath := path.Clean(pathToAccTests)
	cleanedResourcePath := path.Clean(accTestResourcePath)
	if !strings.HasPrefix(cleanedResourcePath, cleanedAccTestsPath) {
		return "", fmt.Errorf("attempted to access file outside acc_tests directory: %s", accTestResourcePath)
	}
	if _, err := os.Stat(accTestResourcePath); os.IsNotExist(err) {
		return "", fmt.Errorf("acc_tests resource does not exist at %s", accTestResourcePath)
	}
	return accTestResourcePath, nil
}

// ReadAccTestResource reads a .tf fixture from internal/provider/acc_tests.
func ReadAccTestResource(elements []string) (string, error) {
	p, err := getPathToAccTestResource(elements)
	if err != nil {
		return "", err
	}
	resource, err := os.ReadFile(filepath.Clean(p))
	if err != nil {
		return "", err
	}
	return string(resource), nil
}

func getProviderConfig() string {
	return `
provider "n8n" {
  # endpoint and api_key from N8N_ENDPOINT / N8N_API_KEY
}
`
}

const (
	folderResourceIDPrefix = "projects/"
	folderResourceIDMid    = "/folders/"
)

func formatFolderResourceID(projectID, folderID string) string {
	return folderResourceIDPrefix + projectID + folderResourceIDMid + folderID
}

func parseFolderResourceID(id string) (projectID, folderID string, err error) {
	id = strings.TrimSpace(id)
	rest, ok := strings.CutPrefix(id, folderResourceIDPrefix)
	if !ok {
		return "", "", fmt.Errorf("folder import id %q must be projects/{project_id}/folders/{folder_id}", id)
	}
	projectID, folderID, ok = strings.Cut(rest, folderResourceIDMid)
	if !ok || projectID == "" || folderID == "" || strings.Contains(folderID, "/") {
		return "", "", fmt.Errorf("folder import id %q must be projects/{project_id}/folders/{folder_id}", id)
	}
	return projectID, folderID, nil
}

func optionalStringPointer(v types.String) *string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	s := strings.TrimSpace(v.ValueString())
	if s == "" {
		return nil
	}
	return &s
}

func optionalStringValue(v *string) types.String {
	if v == nil || strings.TrimSpace(*v) == "" {
		return types.StringNull()
	}
	return types.StringValue(*v)
}

func optionalBoolPointer(v types.Bool) *bool {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	b := v.ValueBool()
	return &b
}

// optionalStringsEqual reports whether two optional Terraform strings match after
// trim, treating null/unknown/empty as the same absent value.
func optionalStringsEqual(a, b types.String) bool {
	ap, bp := optionalStringPointer(a), optionalStringPointer(b)
	if ap == nil || bp == nil {
		return ap == bp
	}
	return *ap == *bp
}

// readMarkdownDescription reads the content of a markdown file from the embedded filesystem.
// The filename parameter should be in the format "internal/provider/docs/..." or "docs/..."
func readMarkdownDescription(ctx context.Context, filename string) (string, error) {
	// Extract the path relative to the 'internal/provider/' prefix
	relativePathInProvider := strings.TrimPrefix(filename, "internal/provider/")

	// If the path doesn't start with "docs/", the TrimPrefix didn't do anything,
	// meaning the input didn't have the "internal/provider/" prefix.
	// In that case, use the filename as-is (for backward compatibility).
	if !strings.HasPrefix(relativePathInProvider, "docs/") && strings.HasPrefix(filename, "docs/") {
		relativePathInProvider = filename
	}

	tflog.Debug(ctx, fmt.Sprintf("Attempting to read embedded markdown file: %s (original filename: %s)", relativePathInProvider, filename))

	// Read from the embedded filesystem
	content, err := embeddedDocs.ReadFile(relativePathInProvider)
	if err != nil {
		// Log the error with the path that failed
		tflog.Error(ctx, fmt.Sprintf("Error reading embedded markdown file %s: %s", relativePathInProvider, err.Error()))
		return "", fmt.Errorf("failed to read markdown file %s (tried embedded path: %s): %w", filename, relativePathInProvider, err)
	}

	return string(content), nil
}
