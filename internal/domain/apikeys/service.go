package apikeys

import (
	"context"
	"errors"
	"fmt"
	"time"

	"plans-features/internal/domain/projects"

	"github.com/google/uuid"
)

type APIKeyService interface {
	CreateKey(ctx context.Context, projectID string) (*CreateAPIKeyResult, error)
	RotateKey(ctx context.Context, projectID string) (*CreateAPIKeyResult, error)
	RevokeKey(ctx context.Context, projectID string, keyPrefix *string) error
	ValidateKey(ctx context.Context, rawKey string) (string, error) // returns projectID
}

type apiKeyService struct {
	repo        APIKeyRepository
	projectRepo projects.ProjectRepository
}

// NewAPIKeyService ahora requiere projectRepo para validar existencia de proyectos
func NewAPIKeyService(repo APIKeyRepository, projectRepo projects.ProjectRepository) APIKeyService {
	return &apiKeyService{repo: repo, projectRepo: projectRepo}
}

func genRawKey() string {
	// simple generator: uuid + timestamp
	return fmt.Sprintf("%s.%d", uuid.New().String(), time.Now().Unix())
}

func (s *apiKeyService) CreateKey(ctx context.Context, projectID string) (*CreateAPIKeyResult, error) {
	// validate project exists
	if _, err := s.projectRepo.GetByID(ctx, projectID); err != nil {
		return nil, errors.New("project not found")
	}
	// revoke existing active keys for project (only 1 active allowed)
	_ = s.repo.Revoke(ctx, projectID, nil)

	raw := genRawKey()
	res, err := s.repo.Create(ctx, projectID, raw)
	if err != nil {
		return nil, err
	}
	return &CreateAPIKeyResult{RawKey: raw, Key: *res}, nil
}

func (s *apiKeyService) RotateKey(ctx context.Context, projectID string) (*CreateAPIKeyResult, error) {
	// validate project exists
	if _, err := s.projectRepo.GetByID(ctx, projectID); err != nil {
		return nil, errors.New("project not found")
	}
	// revoke previous
	_ = s.repo.Revoke(ctx, projectID, nil)

	raw := genRawKey()
	res, err := s.repo.Rotate(ctx, projectID, raw)
	if err != nil {
		return nil, err
	}
	return &CreateAPIKeyResult{RawKey: raw, Key: *res}, nil
}

func (s *apiKeyService) RevokeKey(ctx context.Context, projectID string, keyPrefix *string) error {
	// validate project exists
	if _, err := s.projectRepo.GetByID(ctx, projectID); err != nil {
		return errors.New("project not found")
	}
	return s.repo.Revoke(ctx, projectID, keyPrefix)
}

func (s *apiKeyService) ValidateKey(ctx context.Context, rawKey string) (string, error) {
	res, err := s.repo.Validate(ctx, rawKey)
	if err != nil {
		return "", errors.New("invalid api key")
	}
	return res.ProjectID, nil
}
