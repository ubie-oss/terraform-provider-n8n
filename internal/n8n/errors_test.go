package n8n

import (
	"errors"
	"net/http"
	"testing"
)

func TestIsFeatureUnlicensed(t *testing.T) {
	tests := []struct {
		name    string
		check   func(error) bool
		feature string
		other   string
	}{
		{
			name:    "projectRoleAdmin",
			check:   IsProjectRoleAdminUnlicensed,
			feature: "feat:projectRole:admin",
			other:   "feat:folders",
		},
		{
			name:    "folders",
			check:   IsFoldersUnlicensed,
			feature: "feat:folders",
			other:   "feat:projectRole:admin",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &APIError{
				StatusCode: http.StatusForbidden,
				Body:       `{"message":"Your license does not allow for ` + tt.feature + `. To enable ` + tt.feature + `, please upgrade to a license that supports this feature."}`,
			}
			if !tt.check(err) {
				t.Fatal("expected licensed-feature 403")
			}
			if !tt.check(errors.Join(errors.New("wrap"), err)) {
				t.Fatal("expected wrap to match")
			}
			if tt.check(&APIError{StatusCode: http.StatusForbidden, Body: `{"message":"` + tt.other + `"}`}) {
				t.Fatal("other feature 403 must not match")
			}
			if tt.check(&APIError{StatusCode: http.StatusForbidden, Body: `{"message":"Forbidden"}`}) {
				t.Fatal("missing-scope 403 must not be treated as unlicensed")
			}
			if tt.check(&APIError{StatusCode: http.StatusUnauthorized, Body: tt.feature}) {
				t.Fatal("401 must not match")
			}
		})
	}
}

func TestIsNotFound(t *testing.T) {
	if !IsNotFound(&NotFoundError{Resource: "project", ID: "x"}) {
		t.Fatal("expected not found")
	}
	if IsNotFound(&APIError{StatusCode: 404}) {
		t.Fatal("APIError 404 is not NotFoundError")
	}
}
