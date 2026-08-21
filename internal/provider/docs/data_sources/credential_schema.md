Reads the n8n credential data schema for a type name.

The body is a JSON Schema-like object from `GET /credentials/schema/{type}`. Unknown types produce a not-found error. Schema shapes include n8n-specific property types and are instance-defined (core and community nodes).
