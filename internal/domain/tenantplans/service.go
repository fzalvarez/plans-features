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

	// API methods
	GetTenantPlan(ctx context.Context, tenantID string, projectID string) (*TenantPlanResponse, error)
	AssignTenantPlan(ctx context.Context, tenantID string, projectID string, planID string) (*TenantPlanResponse, error)
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

// GetTenantPlan devuelve la asignación efectiva del tenant para el proyecto o el plan por defecto
func (s *tenantPlanService) GetTenantPlan(ctx context.Context, tenantID string, projectID string) (*TenantPlanResponse, error) {
	// try explicit assignment
	if tp, err := s.repo.GetByTenantAndProject(ctx, tenantID, projectID); err == nil {
		return tp, nil
	}
	// else fallback to default plan for project
	plans, err := s.planRepo.List(ctx, projectID)
	if err != nil {
		return nil, err
	}
	for _, p := range plans {
		if p.IsDefault {
			// return a synthetic assignment (not persisted)
			res := TenantPlanResponse{
				ID:        "",
				TenantID:  tenantID,
				ProjectID: projectID,
				PlanID:    p.ID,
			}
			return &res, nil
		}
	}
	return nil, errors.New("no plan available")
}

// AssignTenantPlan asigna o actualiza la asignación del tenant para el proyecto
func (s *tenantPlanService) AssignTenantPlan(ctx context.Context, tenantID string, projectID string, planID string) (*TenantPlanResponse, error) {
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
	return s.repo.UpsertByTenantAndProject(ctx, tenantID, projectID, planID)
}
