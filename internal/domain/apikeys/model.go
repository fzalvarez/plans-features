package apikeys

import (
	"time"

	"github.com/google/uuid"
)

// Para DB (Scan)
type APIKey struct {
	ID        uuid.UUID `db:"id"`
	ProjectID uuid.UUID `db:"project_id"`
	KeyHash   string    `db:"key_hash"`
	KeyPrefix string    `db:"key_prefix"`
	Revoked   bool      `db:"revoked"`
	CreatedAt time.Time `db:"created_at"`
}

type CreateAPIKeyResult struct {
	RawKey string         `json:"raw_key"`
	Key    APIKeyResponse `json:"key"`
}

type CreateAPIKeyRequest struct{} // Vacío por ahora

type RevokeAPIKeyRequest struct {
	KeyPrefix *string `json:"key_prefix,omitempty"`
}

type APIKeyResponse struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"project_id"`
	KeyPrefix string    `json:"key_prefix"`
	CreatedAt time.Time `json:"created_at"`
	Revoked   bool      `json:"revoked"`
}

// Interno para DB (sin uuid.UUID para Scan simple)
type apiKeyEntity struct {
	ID        uuid.UUID
	ProjectID uuid.UUID
	KeyHash   string
	KeyPrefix string
	CreatedAt time.Time
	Revoked   bool
}

func ToResponse(apiKey *APIKey) *APIKeyResponse {
	return &APIKeyResponse{
		ID:        apiKey.ID,
		ProjectID: apiKey.ProjectID,
		KeyPrefix: apiKey.KeyPrefix,
		CreatedAt: apiKey.CreatedAt,
		Revoked:   apiKey.Revoked,
	}
}
