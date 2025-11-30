package features

import (
	"context"
	"errors"

	"plans-features/internal/domain/projects"
)

// TODO: lógica de negocio para features

type FeatureService interface {
	ListFeatures(ctx context.Context, projectID string) ([]FeatureResponse, error)
	CreateFeature(ctx context.Context, projectID string, req CreateFeatureRequest) (*FeatureResponse, error)
	GetFeature(ctx context.Context, projectID string, featureID string) (*FeatureResponse, error)
	UpdateFeature(ctx context.Context, projectID string, featureID string, req UpdateFeatureRequest) (*FeatureResponse, error)
}

type featureService struct {
	repo        FeatureRepository
	projectRepo projects.ProjectRepository
}

func NewFeatureService(repo FeatureRepository, projectRepo projects.ProjectRepository) FeatureService {
	return &featureService{repo: repo, projectRepo: projectRepo}
}

func (s *featureService) ListFeatures(ctx context.Context, projectID string) ([]FeatureResponse, error) {
	return s.repo.List(ctx, projectID)
}

func (s *featureService) CreateFeature(ctx context.Context, projectID string, req CreateFeatureRequest) (*FeatureResponse, error) {
	// validate project exists
	if _, err := s.projectRepo.GetByID(ctx, projectID); err != nil {
		return nil, errors.New("project not found")
	}
	// name required
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	// code unique within project
	existing, err := s.repo.List(ctx, projectID)
	if err != nil {
		return nil, err
	}
	for _, f := range existing {
		if f.Code == req.Code {
			return nil, errors.New("feature code already exists")
		}
	}
	return s.repo.Create(ctx, projectID, req)
}

func (s *featureService) GetFeature(ctx context.Context, projectID string, featureID string) (*FeatureResponse, error) {
	return s.repo.GetByID(ctx, projectID, featureID)
}

func (s *featureService) UpdateFeature(ctx context.Context, projectID string, featureID string, req UpdateFeatureRequest) (*FeatureResponse, error) {
	// validate project exists
	if _, err := s.projectRepo.GetByID(ctx, projectID); err != nil {
		return nil, errors.New("project not found")
	}
	// ignore code changes (UpdateFeatureRequest has no Code)
	return s.repo.Update(ctx, projectID, featureID, req)
}
