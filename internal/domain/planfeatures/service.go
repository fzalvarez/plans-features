package planfeatures

import (
	"context"
	"errors"

	"plans-features/internal/domain/features"
	"plans-features/internal/domain/plans"
	"plans-features/internal/domain/projects"
)

type PlanFeatureService interface {
	ListByPlan(ctx context.Context, projectID string, planID string) ([]PlanFeatureResponse, error)
	AssignFeature(ctx context.Context, projectID string, planID string, req AssignFeatureRequest) (*PlanFeatureResponse, error)
}

type planFeatureService struct {
	repo        PlanFeatureRepository
	planRepo    plans.PlanRepository
	featureRepo features.FeatureRepository
	projectRepo projects.ProjectRepository
}

func NewPlanFeatureService(repo PlanFeatureRepository, planRepo plans.PlanRepository, featureRepo features.FeatureRepository, projectRepo projects.ProjectRepository) PlanFeatureService {
	return &planFeatureService{repo: repo, planRepo: planRepo, featureRepo: featureRepo, projectRepo: projectRepo}
}

func (s *planFeatureService) ListByPlan(ctx context.Context, projectID string, planID string) ([]PlanFeatureResponse, error) {
	return s.repo.ListByPlan(ctx, projectID, planID)
}

func (s *planFeatureService) AssignFeature(ctx context.Context, projectID string, planID string, req AssignFeatureRequest) (*PlanFeatureResponse, error) {
	// validate project exists
	if _, err := s.projectRepo.GetByID(ctx, projectID); err != nil {
		return nil, errors.New("project not found")
	}
	// validate plan exists and belongs to project
	p, err := s.planRepo.GetByID(ctx, projectID, planID)
	if err != nil {
		return nil, errors.New("plan not found")
	}
	if p.ProjectID != projectID {
		return nil, errors.New("plan does not belong to project")
	}
	// validate feature exists and belongs to project
	f, err := s.featureRepo.GetByID(ctx, projectID, req.FeatureID)
	if err != nil {
		return nil, errors.New("feature not found")
	}
	if f.ProjectID != projectID {
		return nil, errors.New("feature does not belong to project")
	}
	// validate value type based on feature.Type
	switch f.Type {
	case "flag":
		_, ok := req.Value.(bool)
		if !ok {
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
		_, ok := req.Value.(string)
		if !ok {
			return nil, errors.New("value must be string for value feature")
		}
	default:
		return nil, errors.New("invalid feature type")
	}
	// ensure not duplicate in plan
	exists, _ := s.repo.Exists(ctx, projectID, planID, req.FeatureID)
	if exists {
		return nil, errors.New("feature already assigned to plan")
	}
	// create
	return s.repo.Create(ctx, projectID, planID, req)
}
