package features

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type FeatureRepository interface {
	List(ctx context.Context, projectID string) ([]FeatureResponse, error)
	Create(ctx context.Context, projectID string, req CreateFeatureRequest) (*FeatureResponse, error)
	GetByID(ctx context.Context, projectID string, featureID string) (*FeatureResponse, error)
	Update(ctx context.Context, projectID string, featureID string, req UpdateFeatureRequest) (*FeatureResponse, error)
}

type featureRepository struct {
	data map[string]FeatureResponse
}

func NewFeatureRepository() FeatureRepository {
	return &featureRepository{
		data: make(map[string]FeatureResponse),
	}
}

func (r *featureRepository) List(ctx context.Context, projectID string) ([]FeatureResponse, error) {
	res := make([]FeatureResponse, 0)
	for _, v := range r.data {
		if v.ProjectID == projectID {
			res = append(res, v)
		}
	}
	return res, nil
}

func (r *featureRepository) Create(ctx context.Context, projectID string, req CreateFeatureRequest) (*FeatureResponse, error) {
	id := uuid.New().String()
	f := FeatureResponse{
		ID:          id,
		ProjectID:   projectID,
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
	}
	r.data[id] = f
	return &f, nil
}

func (r *featureRepository) GetByID(ctx context.Context, projectID string, featureID string) (*FeatureResponse, error) {
	f, ok := r.data[featureID]
	if !ok || f.ProjectID != projectID {
		return nil, errors.New("not found")
	}
	return &f, nil
}

func (r *featureRepository) Update(ctx context.Context, projectID string, featureID string, req UpdateFeatureRequest) (*FeatureResponse, error) {
	f, ok := r.data[featureID]
	if !ok || f.ProjectID != projectID {
		return nil, errors.New("not found")
	}
	if req.Name != nil {
		f.Name = *req.Name
	}
	if req.Description != nil {
		f.Description = *req.Description
	}
	if req.IsActive != nil {
		f.IsActive = *req.IsActive
	}
	r.data[featureID] = f
	return &f, nil
}
