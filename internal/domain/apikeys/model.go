package apikeys

import "time"

// CreateAPIKeyResult contiene la clave raw (entregada solo una vez) y metadatos guardados
type CreateAPIKeyResult struct {
	RawKey string
	Key    APIKeyResponse
}

// CreateAPIKeyRequest no necesita campos por ahora
type CreateAPIKeyRequest struct{}

// RevokeAPIKeyRequest puede opcionalmente especificar un prefijo de key a revocar
type RevokeAPIKeyRequest struct {
	KeyPrefix *string `json:"key_prefix,omitempty"`
}

// APIKeyResponse representa la información almacenada de una API key (sin el raw)
type APIKeyResponse struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	KeyPrefix string    `json:"key_prefix"`
	CreatedAt time.Time `json:"created_at"`
	Revoked   bool      `json:"revoked"`
}

// entidad interna
type apiKeyEntity struct {
	ID        string
	ProjectID string
	KeyHash   string
	KeyPrefix string
	CreatedAt time.Time
	Revoked   bool
}
