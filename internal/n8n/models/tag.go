package models

// Tag is an n8n Public API global tag document (GET/POST/PUT/DELETE /tags).
type Tag struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
}

// TagList is the GET /tags cursor envelope.
type TagList struct {
	Data       []Tag   `json:"data"`
	NextCursor *string `json:"nextCursor"`
}

// TagWrite is the request body for POST /tags and PUT /tags/{id}.
type TagWrite struct {
	Name string `json:"name"`
}
