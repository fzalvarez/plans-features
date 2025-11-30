package tenantplans

import (
	"context"
	"errors"

	"plans-features/internal/domain/plans"
	"plans-features/internal/domain/projects"
)

type TenantPlanService interface {
	ListAssignments(ctx context.Context, tenantID string) ([]TenantPlanResponse, error)
	CreateAssignment(ctx context.Context, tenantID string, req CreateTenantPlanRequest) (*TenantPlanResponse, error)
	UpdateAssignment(ctx context.Context, tenantID string, assignmentID string, req UpdateTenantPlanRequest) (*TenantPlanResponse, error)
}

type tenantPlanService struct {
	repo        TenantPlanRepository
	projectRepo projects.ProjectRepository
	planRepo    plans.PlanRepository
}

func NewTenantPlanService(
	repo TenantPlanRepository,
	projectRepo projects.ProjectRepository,
	planRepo plans.PlanRepository,
) TenantPlanService {
	return &tenantPlanService{
		repo:        repo,
		projectRepo: projectRepo,
		planRepo:    planRepo,
	}
}

func (s *tenantPlanService) ListAssignments(ctx context.Context, tenantID string) ([]TenantPlanResponse, error) {
	return s.repo.ListByTenant(ctx, tenantID)
}

func (s *tenantPlanService) CreateAssignment(ctx context.Context, tenantID string, req CreateTenantPlanRequest) (*TenantPlanResponse, error) {

	// project code required
	if req.ProjectCode == "" {
		return nil, errors.New("project code is required")
	}

	// plan code required
	if req.PlanCode == "" {
		return nil, errors.New("plan code is required")
	}

	// (Validaciones futuras con BD: validar projectCode y planCode)
	// Por ahora solo validaciones internas:

	// tenant can have only one plan per project
	assignments, err := s.repo.ListByTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	for _, a := range assignments {
		if a.ProjectID == req.ProjectCode {
			return nil, errors.New("project already assigned")
		}
	}

	return s.repo.Create(ctx, tenantID, req)
}

func (s *tenantPlanService) UpdateAssignment(ctx context.Context, tenantID string, assignmentID string, req UpdateTenantPlanRequest) (*TenantPlanResponse, error) {
	// ignore project changes (not allowed yet)
	return s.repo.Update(ctx, tenantID, assignmentID, req)
}
