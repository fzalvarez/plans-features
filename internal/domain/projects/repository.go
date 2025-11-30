package projects

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type ProjectRepository interface {
	List(ctx context.Context) ([]ProjectResponse, error)
	Create(ctx context.Context, req CreateProjectRequest) (*ProjectResponse, error)
	GetByID(ctx context.Context, id string) (*ProjectResponse, error)
	Update(ctx context.Context, id string, req UpdateProjectRequest) (*ProjectResponse, error)
}

type projectRepository struct {
	data map[string]ProjectResponse
}

func NewProjectRepository() ProjectRepository {
	return &projectRepository{
		data: make(map[string]ProjectResponse),
	}
}

func (r *projectRepository) List(ctx context.Context) ([]ProjectResponse, error) {
	res := make([]ProjectResponse, 0, len(r.data))
	for _, v := range r.data {
		res = append(res, v)
	}
	return res, nil
}

func (r *projectRepository) Create(ctx context.Context, req CreateProjectRequest) (*ProjectResponse, error) {
	id := uuid.New().String()
	p := ProjectResponse{
		ID:          id,
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
	}
	r.data[id] = p
	return &p, nil
}

func (r *projectRepository) GetByID(ctx context.Context, id string) (*ProjectResponse, error) {
	p, ok := r.data[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return &p, nil
}

func (r *projectRepository) Update(ctx context.Context, id string, req UpdateProjectRequest) (*ProjectResponse, error) {
	p, ok := r.data[id]
	if !ok {
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
	r.data[id] = p
	return &p, nil
}
