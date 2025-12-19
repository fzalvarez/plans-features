package features

import (
	"context"
	"errors"
	"plans-features/internal/domain/projects"

	"github.com/google/uuid"
)

type FeatureService interface {
	ListFeatures(ctx context.Context, projectID uuid.UUID) ([]FeatureResponse, error)
	CreateFeature(ctx context.Context, projectID uuid.UUID, req CreateFeatureRequest) (*FeatureResponse, error)
	GetFeature(ctx context.Context, projectID uuid.UUID, featureID uuid.UUID) (*FeatureResponse, error)
	UpdateFeature(ctx context.Context, projectID uuid.UUID, featureID uuid.UUID, req UpdateFeatureRequest) (*FeatureResponse, error)
}

type featureService struct {
	repo        FeatureRepository
	projectRepo projects.ProjectRepository
}

func NewFeatureService(repo FeatureRepository, projectRepo projects.ProjectRepository) FeatureService {
	return &featureService{repo: repo, projectRepo: projectRepo}
}

func isValidType(t string) bool {
	switch t {
	case "flag", "numeric", "value":
		return true
	default:
		return false
	}
}

func (s *featureService) ListFeatures(ctx context.Context, projectID uuid.UUID) ([]FeatureResponse, error) {
	return s.repo.List(ctx, projectID)
}

func (s *featureService) CreateFeature(ctx context.Context, projectID uuid.UUID, req CreateFeatureRequest) (*FeatureResponse, error) {
	// validate project exists
	if _, err := s.projectRepo.GetByID(ctx, projectID); err != nil {
		return nil, errors.New("project not found")
	}
	// name required
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	// type valid
	if !isValidType(req.Type) {
		return nil, errors.New("invalid type")
	}
	// code unique within project
	existing, err := s.repo.List(ctx, projectID)
	if err != nil {
		return nil, err
	}
	for _, f := range existing {
		if f.Code == normalizeCode(req.Code) {
			return nil, errors.New("feature code already exists")
		}
	}
	// default isActive handled by repo
	return s.repo.Create(ctx, projectID, req)
}

func (s *featureService) GetFeature(ctx context.Context, projectID uuid.UUID, featureID uuid.UUID) (*FeatureResponse, error) {
	return s.repo.GetByID(ctx, projectID, featureID)
}

func (s *featureService) UpdateFeature(ctx context.Context, projectID uuid.UUID, featureID uuid.UUID, req UpdateFeatureRequest) (*FeatureResponse, error) {
	// validate project exists
	if _, err := s.projectRepo.GetByID(ctx, projectID); err != nil {
		return nil, errors.New("project not found")
	}
	// ignore code changes (UpdateFeatureRequest has no Code)
	return s.repo.Update(ctx, projectID, featureID, req)
}
