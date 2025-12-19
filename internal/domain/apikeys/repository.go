package apikeys

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type APIKeyRepository interface {
	Create(ctx context.Context, projectID uuid.UUID, rawKey string) (*APIKeyResponse, error)
	Rotate(ctx context.Context, projectID uuid.UUID, rawKey string) (*APIKeyResponse, error)
	Revoke(ctx context.Context, projectID uuid.UUID, keyPrefix *string) error
	Validate(ctx context.Context, rawKey string) (*APIKeyResponse, error)
}

type apiKeyRepository struct {
	db *sql.DB
}

func NewAPIKeyRepository(db *sql.DB) APIKeyRepository {
	return &apiKeyRepository{db: db}
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

func (r *apiKeyRepository) Create(ctx context.Context, projectID uuid.UUID, rawKey string) (*APIKeyResponse, error) {
	id := uuid.New()
	keyHash := hashKey(rawKey)
	keyPrefix := prefixOf(rawKey)

	apiKey := &APIKey{}
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO api_keys (id, project_id, key_hash, key_prefix, revoked) 
         VALUES ($1, $2, $3, $4, $5) 
         RETURNING id, project_id, key_hash, key_prefix, revoked, created_at`,
		id, projectID, keyHash, keyPrefix, false).
		Scan(&apiKey.ID, &apiKey.ProjectID, &apiKey.KeyHash, &apiKey.KeyPrefix,
			&apiKey.Revoked, &apiKey.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("create api key: %w", err)
	}

	return ToResponse(apiKey), nil
}

func (r *apiKeyRepository) Rotate(ctx context.Context, projectID uuid.UUID, rawKey string) (*APIKeyResponse, error) {
	return r.Create(ctx, projectID, rawKey)
}

func (r *apiKeyRepository) Revoke(ctx context.Context, projectID uuid.UUID, keyPrefix *string) error {
	if keyPrefix == nil {
		_, err := r.db.ExecContext(ctx,
			`UPDATE api_keys SET revoked = true WHERE project_id = $1`, projectID)
		return err
	}

	result, err := r.db.ExecContext(ctx,
		`UPDATE api_keys SET revoked = true WHERE project_id = $1 AND key_prefix = $2`,
		projectID, *keyPrefix)
	if err != nil {
		return fmt.Errorf("revoke api key: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("no api keys found with that prefix")
	}
	return nil
}

func (r apiKeyRepository) Validate(ctx context.Context, rawKey string) (*APIKeyResponse, error) {
	keyHash := hashKey(rawKey)

	apiKey := &APIKey{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, project_id, key_hash, key_prefix, revoked, created_at 
         FROM api_keys 
         WHERE key_hash = $1 AND revoked = false 
         LIMIT 1`,
		keyHash).
		Scan(&apiKey.ID, &apiKey.ProjectID, &apiKey.KeyHash,
			&apiKey.KeyPrefix, &apiKey.Revoked, &apiKey.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("invalid api key")
	}
	if err != nil {
		return nil, fmt.Errorf("validate api key: %w", err)
	}

	return ToResponse(apiKey), nil
}
