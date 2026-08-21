package n8n

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// NotFoundError indicates the requested resource does not exist (HTTP 404).
type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s %q not found", e.Resource, e.ID)
}

// APIError is a non-success response from the n8n Public API.
type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	if e.Body == "" {
		return fmt.Sprintf("n8n API error: HTTP %d", e.StatusCode)
	}
	return fmt.Sprintf("n8n API error: HTTP %d: %s", e.StatusCode, e.Body)
}

// IsNotFound reports whether err is or wraps a NotFoundError.
func IsNotFound(err error) bool {
	var nf *NotFoundError
	return errors.As(err, &nf)
}

// IsProjectRoleAdminUnlicensed reports a 403 whose body names feat:projectRole:admin.
// Live Community n8n and n8n Public API tests both use FeatureNotLicensedError for
// every project HTTP operation when that license flag is off.
func IsProjectRoleAdminUnlicensed(err error) bool {
	return isFeatureUnlicensed(err, "feat:projectRole:admin")
}

// IsFoldersUnlicensed reports a 403 whose body names feat:folders.
func IsFoldersUnlicensed(err error) bool {
	return isFeatureUnlicensed(err, "feat:folders")
}

func isFeatureUnlicensed(err error, feature string) bool {
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusForbidden {
		return false
	}
	return strings.Contains(apiErr.Body, feature)
}
