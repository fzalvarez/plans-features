package apikeys

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
)

type APIKeyRepository interface {
	Create(ctx context.Context, projectID string, rawKey string) (*APIKeyResponse, error)
	Rotate(ctx context.Context, projectID string, rawKey string) (*APIKeyResponse, error)
	Revoke(ctx context.Context, projectID string, keyPrefix *string) error
	Validate(ctx context.Context, rawKey string) (*APIKeyResponse, error)
}

type apiKeyRepository struct {
	data map[string]apiKeyEntity // id -> entity
}

func NewAPIKeyRepository() APIKeyRepository {
	return &apiKeyRepository{data: make(map[string]apiKeyEntity)}
}

func hashKey(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}

func prefixOf(raw string) string {
	if len(raw) <= 8 {
		return raw
	}
	return raw[:8]
}

func (r *apiKeyRepository) Create(ctx context.Context, projectID string, rawKey string) (*APIKeyResponse, error) {
	id := uuid.New().String()
	ent := apiKeyEntity{
		ID:        id,
		ProjectID: projectID,
		KeyHash:   hashKey(rawKey),
		KeyPrefix: prefixOf(rawKey),
		CreatedAt: time.Now().UTC(),
		Revoked:   false,
	}
	r.data[id] = ent
	res := APIKeyResponse{
		ID:        ent.ID,
		ProjectID: ent.ProjectID,
		KeyPrefix: ent.KeyPrefix,
		CreatedAt: ent.CreatedAt,
		Revoked:   ent.Revoked,
	}
	return &res, nil
}

// Rotate: create a new key and leave old keys as-is (could revoke old if needed)
func (r *apiKeyRepository) Rotate(ctx context.Context, projectID string, rawKey string) (*APIKeyResponse, error) {
	return r.Create(ctx, projectID, rawKey)
}

func (r *apiKeyRepository) Revoke(ctx context.Context, projectID string, keyPrefix *string) error {
	for id, e := range r.data {
		if e.ProjectID != projectID {
			continue
		}
		if keyPrefix != nil {
			if e.KeyPrefix != *keyPrefix {
				continue
			}
		}
		e.Revoked = true
		r.data[id] = e
	}
	return nil
}

func (r *apiKeyRepository) Validate(ctx context.Context, rawKey string) (*APIKeyResponse, error) {
	h := hashKey(rawKey)
	for _, e := range r.data {
		if e.KeyHash == h && !e.Revoked {
			res := APIKeyResponse{
				ID:        e.ID,
				ProjectID: e.ProjectID,
				KeyPrefix: e.KeyPrefix,
				CreatedAt: e.CreatedAt,
				Revoked:   e.Revoked,
			}
			return &res, nil
		}
	}
	return nil, errors.New("invalid api key")
}
