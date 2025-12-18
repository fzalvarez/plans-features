package plans

import (
	"context"
	"errors"

	"plans-features/internal/domain/projects"
)

// PlanService defines business operations for plans
type PlanService interface {
	ListPlans(ctx context.Context, projectID string) ([]PlanResponse, error)
	CreatePlan(ctx context.Context, projectID string, req CreatePlanRequest) (*PlanResponse, error)
	GetPlan(ctx context.Context, projectID string, planID string) (*PlanResponse, error)
	UpdatePlan(ctx context.Context, projectID string, planID string, req UpdatePlanRequest) (*PlanResponse, error)
}

type planService struct {
	repo        PlanRepository
	projectRepo projects.ProjectRepository
}

func NewPlanService(repo PlanRepository, projectRepo projects.ProjectRepository) PlanService {
	return &planService{repo: repo, projectRepo: projectRepo}
}

func (s *planService) ListPlans(ctx context.Context, projectID string) ([]PlanResponse, error) {
	return s.repo.List(ctx, projectID)
}

func (s *planService) CreatePlan(ctx context.Context, projectID string, req CreatePlanRequest) (*PlanResponse, error) {
	// validate project exists
	if _, err := s.projectRepo.GetByID(ctx, projectID); err != nil {
		return nil, errors.New("project not found")
	}
	// name required
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	// code unique within project
	plans, err := s.repo.List(ctx, projectID)
	if err != nil {
		return nil, err
	}
	for _, p := range plans {
		if p.Code == normalizeCode(req.Code) {
			return nil, errors.New("plan code already exists")
		}
	}
	// if first plan for project and IsDefault == false, force it to true
	if len(plans) == 0 && !req.IsDefault {
		req.IsDefault = true
	}
	// if IsDefault true, unset others
	if req.IsDefault {
		for _, p := range plans {
			if p.IsDefault {
				falseVal := false
				_, _ = s.repo.Update(ctx, projectID, p.ID, UpdatePlanRequest{IsDefault: &falseVal})
			}
		}
	}
	return s.repo.Create(ctx, projectID, req)
}

func (s *planService) GetPlan(ctx context.Context, projectID string, planID string) (*PlanResponse, error) {
	return s.repo.GetByID(ctx, projectID, planID)
}

func (s *planService) UpdatePlan(ctx context.Context, projectID string, planID string, req UpdatePlanRequest) (*PlanResponse, error) {
	// validate project exists
	if _, err := s.projectRepo.GetByID(ctx, projectID); err != nil {
		return nil, errors.New("project not found")
	}
	// If IsDefault true, unset others
	if req.IsDefault != nil && *req.IsDefault {
		plans, err := s.repo.List(ctx, projectID)
		if err != nil {
			return nil, err
		}
		for _, p := range plans {
			if p.IsDefault && p.ID != planID {
				falseVal := false
				_, _ = s.repo.Update(ctx, projectID, p.ID, UpdatePlanRequest{IsDefault: &falseVal})
			}
		}
	}
	// ignore any Code changes (UpdatePlanRequest does not have Code)
	return s.repo.Update(ctx, projectID, planID, req)
}
