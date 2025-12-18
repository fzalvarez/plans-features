package projects

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
)

type ProjectRepository interface {
	List(ctx context.Context) ([]ProjectResponse, error)
	Create(ctx context.Context, req CreateProjectRequest) (*ProjectResponse, error)
	GetByID(ctx context.Context, id string) (*ProjectResponse, error)
	GetByCode(ctx context.Context, code string) (*ProjectResponse, error)
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

func normalizeCode(code string) string {
	return strings.ToLower(strings.TrimSpace(code))
}

func (r *projectRepository) Create(ctx context.Context, req CreateProjectRequest) (*ProjectResponse, error) {
	id := uuid.New().String()
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	p := ProjectResponse{
		ID:          id,
		Code:        normalizeCode(req.Code),
		Name:        req.Name,
		Description: req.Description,
		IsActive:    isActive,
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

func (r *projectRepository) GetByCode(ctx context.Context, code string) (*ProjectResponse, error) {
	n := normalizeCode(code)
	for _, p := range r.data {
		if p.Code == n {
			return &p, nil
		}
	}
	return nil, errors.New("not found")
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
