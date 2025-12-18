package plans

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
)

type PlanRepository interface {
	List(ctx context.Context, projectID string) ([]PlanResponse, error)
	Create(ctx context.Context, projectID string, req CreatePlanRequest) (*PlanResponse, error)
	GetByID(ctx context.Context, projectID string, planID string) (*PlanResponse, error)
	Update(ctx context.Context, projectID string, planID string, req UpdatePlanRequest) (*PlanResponse, error)
}

type planRepository struct {
	data map[string]PlanResponse
}

func NewPlanRepository() PlanRepository {
	return &planRepository{
		data: make(map[string]PlanResponse),
	}
}

func (r *planRepository) List(ctx context.Context, projectID string) ([]PlanResponse, error) {
	res := make([]PlanResponse, 0)
	for _, v := range r.data {
		if v.ProjectID == projectID {
			res = append(res, v)
		}
	}
	return res, nil
}

func normalizeCode(code string) string {
	return strings.ToLower(strings.TrimSpace(code))
}

func (r *planRepository) Create(ctx context.Context, projectID string, req CreatePlanRequest) (*PlanResponse, error) {
	id := uuid.New().String()
	p := PlanResponse{
		ID:          id,
		ProjectID:   projectID,
		Code:        normalizeCode(req.Code),
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
		IsDefault:   req.IsDefault,
		Limits:      req.Limits,
	}
	r.data[id] = p
	return &p, nil
}

func (r *planRepository) GetByID(ctx context.Context, projectID string, planID string) (*PlanResponse, error) {
	p, ok := r.data[planID]
	if !ok || p.ProjectID != projectID {
		return nil, errors.New("not found")
	}
	return &p, nil
}

func (r *planRepository) Update(ctx context.Context, projectID string, planID string, req UpdatePlanRequest) (*PlanResponse, error) {
	p, ok := r.data[planID]
	if !ok || p.ProjectID != projectID {
		return nil, errors.New("not found")
	}
	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.Description != nil {
		p.Description = *req.Description
	}
	if req.IsActive != nil {
		p.IsActive = *req.IsActive
	}
	if req.IsDefault != nil {
		p.IsDefault = *req.IsDefault
	}
	if req.Limits != nil {
		p.Limits = req.Limits
	}
	r.data[planID] = p
	return &p, nil
}
