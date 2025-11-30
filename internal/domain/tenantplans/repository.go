// TODO: interface + implementación con sqlc

package tenantplans

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type TenantPlanRepository interface {
	ListByTenant(ctx context.Context, tenantID string) ([]TenantPlanResponse, error)
	Create(ctx context.Context, tenantID string, req CreateTenantPlanRequest) (*TenantPlanResponse, error)
	Update(ctx context.Context, tenantID string, assignmentID string, req UpdateTenantPlanRequest) (*TenantPlanResponse, error)
}

type tenantPlanRepository struct {
	data map[string]TenantPlanResponse
}

func NewTenantPlanRepository() TenantPlanRepository {
	return &tenantPlanRepository{data: make(map[string]TenantPlanResponse)}
}

func (r *tenantPlanRepository) ListByTenant(ctx context.Context, tenantID string) ([]TenantPlanResponse, error) {
	res := make([]TenantPlanResponse, 0)
	for _, v := range r.data {
		if v.TenantID == tenantID {
			res = append(res, v)
		}
	}
	return res, nil
}

func (r *tenantPlanRepository) Create(ctx context.Context, tenantID string, req CreateTenantPlanRequest) (*TenantPlanResponse, error) {
	id := uuid.New().String()
	tp := TenantPlanResponse{
		ID:        id,
		TenantID:  tenantID,
		ProjectID: req.ProjectCode,
		PlanID:    req.PlanCode,
	}
	r.data[id] = tp
	return &tp, nil
}

func (r *tenantPlanRepository) Update(ctx context.Context, tenantID string, assignmentID string, req UpdateTenantPlanRequest) (*TenantPlanResponse, error) {
	tp, ok := r.data[assignmentID]
	if !ok || tp.TenantID != tenantID {
		return nil, errors.New("not found")
	}
	if req.PlanCode != nil {
		tp.PlanID = *req.PlanCode
	}
	r.data[assignmentID] = tp
	return &tp, nil
}
