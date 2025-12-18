package planfeatures

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type PlanFeatureRepository interface {
	ListByPlan(ctx context.Context, projectID string, planID string) ([]PlanFeatureResponse, error)
	Create(ctx context.Context, projectID string, planID string, req AssignFeatureRequest) (*PlanFeatureResponse, error)
	Exists(ctx context.Context, projectID string, planID string, featureID string) (bool, error)
}

type planFeatureRepository struct {
	data map[string]planFeatureEntity
}

func NewPlanFeatureRepository() PlanFeatureRepository {
	return &planFeatureRepository{data: make(map[string]planFeatureEntity)}
}

func (r *planFeatureRepository) ListByPlan(ctx context.Context, projectID string, planID string) ([]PlanFeatureResponse, error) {
	res := make([]PlanFeatureResponse, 0)
	for _, e := range r.data {
		if e.PlanID == planID && e.ProjectID == projectID {
			res = append(res, PlanFeatureResponse{
				ID:        e.ID,
				PlanID:    e.PlanID,
				ProjectID: e.ProjectID,
				FeatureID: e.FeatureID,
				Value:     e.Value,
			})
		}
	}
	return res, nil
}

func (r *planFeatureRepository) Exists(ctx context.Context, projectID string, planID string, featureID string) (bool, error) {
	for _, e := range r.data {
		if e.PlanID == planID && e.ProjectID == projectID && e.FeatureID == featureID {
			return true, nil
		}
	}
	return false, nil
}

func (r *planFeatureRepository) Create(ctx context.Context, projectID string, planID string, req AssignFeatureRequest) (*PlanFeatureResponse, error) {
	// ensure not exists
	exists, _ := r.Exists(ctx, projectID, planID, req.FeatureID)
	if exists {
		return nil, errors.New("feature already assigned to plan")
	}
	id := uuid.New().String()
	ent := planFeatureEntity{
		ID:        id,
		PlanID:    planID,
		ProjectID: projectID,
		FeatureID: req.FeatureID,
		Value:     req.Value,
	}
	r.data[id] = ent
	res := PlanFeatureResponse{
		ID:        ent.ID,
		PlanID:    ent.PlanID,
		ProjectID: ent.ProjectID,
		FeatureID: ent.FeatureID,
		Value:     ent.Value,
	}
	return &res, nil
}
