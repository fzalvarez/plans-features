package planfeatures

import (
	"context"
	"errors"
	"fmt"

	"plans-features/internal/domain/features"
	"plans-features/internal/domain/plans"
	"plans-features/internal/domain/projects"

	"github.com/google/uuid"
)

type PlanFeatureService interface {
	ListByPlan(ctx context.Context, projectID uuid.UUID, planID uuid.UUID) ([]PlanFeatureResponse, error)
	AssignFeature(ctx context.Context, projectID uuid.UUID, planID uuid.UUID, req AssignFeatureRequest) (*PlanFeatureResponse, error)
}

type planFeatureService struct {
	repo        PlanFeatureRepository
	planRepo    plans.PlanRepository
	featureRepo features.FeatureRepository
	projectRepo projects.ProjectRepository
}

func NewPlanFeatureService(repo PlanFeatureRepository, planRepo plans.PlanRepository, featureRepo features.FeatureRepository, projectRepo projects.ProjectRepository) PlanFeatureService {
	return &planFeatureService{
		repo:        repo,
		planRepo:    planRepo,
		featureRepo: featureRepo,
		projectRepo: projectRepo,
	}
}

func (s *planFeatureService) ListByPlan(ctx context.Context, projectID uuid.UUID, planID uuid.UUID) ([]PlanFeatureResponse, error) {
	return s.repo.ListByPlan(ctx, projectID, planID)
}

func (s *planFeatureService) AssignFeature(ctx context.Context, projectID uuid.UUID, planID uuid.UUID, req AssignFeatureRequest) (*PlanFeatureResponse, error) {
	// 1. Validate project exists
	if _, err := s.projectRepo.GetByID(ctx, projectID); err != nil {
		return nil, errors.New("project not found")
	}

	// 2. Validate plan exists and belongs to project
	plan, err := s.planRepo.GetByID(ctx, projectID, planID) // ← UUIDs
	if err != nil {
		return nil, errors.New("plan not found")
	}
	if plan.ProjectID != projectID {
		return nil, errors.New("plan does not belong to project")
	}

	// 3. Validate feature exists and belongs to project
	feature, err := s.featureRepo.GetByID(ctx, projectID, req.FeatureID) // ← UUIDs
	if err != nil {
		return nil, errors.New("feature not found")
	}
	if feature.ProjectID != projectID {
		return nil, errors.New("feature does not belong to project")
	}

	// 4. Validate value type based on feature.Type
	switch feature.Type {
	case "flag":
		if _, ok := req.Value.(bool); !ok {
			return nil, errors.New("value must be boolean for flag feature")
		}
	case "numeric":
		switch req.Value.(type) {
		case float64, float32, int, int64, int32:
			// ok
		default:
			return nil, errors.New("value must be numeric for numeric feature")
		}
	case "value":
		if _, ok := req.Value.(string); !ok {
			return nil, errors.New("value must be string for value feature")
		}
	default:
		return nil, errors.New("invalid feature type")
	}

	// 5. Ensure not duplicate in plan (repo ya valida UNIQUE constraint)
	exists, err := s.repo.Exists(ctx, projectID, planID, req.FeatureID)
	if err != nil {
		return nil, fmt.Errorf("check exists: %w", err)
	}
	if exists {
		return nil, errors.New("feature already assigned to plan")
	}

	// 6. Create assignment
	return s.repo.Create(ctx, projectID, planID, req)
}
