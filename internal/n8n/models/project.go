package models

const (
	// ProjectTypeTeam is the only project type this provider manages.
	ProjectTypeTeam = "team"
	// ProjectTypePersonal is the owner's built-in project; it cannot be created or deleted as a team project.
	ProjectTypePersonal = "personal"
)

// Project is an n8n Public API project document.
// Create 201 returns extra fields (icon, description, role, scopes) that OpenAPI omits;
// unknown JSON keys are ignored. Terraform v1 only uses ID, Name, and Type.
type Project struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Icon        *string `json:"icon,omitempty"`
	Description *string `json:"description,omitempty"`
	Role        string  `json:"role,omitempty"`
	CreatorID   string  `json:"creatorId,omitempty"`
	CreatedAt   string  `json:"createdAt,omitempty"`
	UpdatedAt   string  `json:"updatedAt,omitempty"`
}

// IsTeam reports whether the project is a team project.
func (p Project) IsTeam() bool {
	return p.Type == ProjectTypeTeam
}

// ProjectList is the GET /projects pagination envelope.
type ProjectList struct {
	Data       []Project `json:"data"`
	NextCursor *string   `json:"nextCursor"`
}

// ProjectWrite is the request body for POST /projects and PUT /projects/{id}.
type ProjectWrite struct {
	Name string `json:"name"`
}
