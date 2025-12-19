package projects

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type ProjectService interface {
	ListProjects(ctx context.Context) ([]ProjectResponse, error)
	CreateProject(ctx context.Context, req CreateProjectRequest) (*ProjectResponse, error)
	GetProject(ctx context.Context, id uuid.UUID) (*ProjectResponse, error)
	UpdateProject(ctx context.Context, id uuid.UUID, req UpdateProjectRequest) (*ProjectResponse, error)
}

type projectService struct {
	repo ProjectRepository
}

func NewProjectService(repo ProjectRepository) ProjectService {
	return &projectService{repo: repo}
}

func (s *projectService) ListProjects(ctx context.Context) ([]ProjectResponse, error) {
	return s.repo.List(ctx)
}

func (s *projectService) CreateProject(ctx context.Context, req CreateProjectRequest) (*ProjectResponse, error) {
	// Validation: code required
	if req.Code == "" {
		return nil, errors.New("code is required")
	}
	// Validation: name required
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	// Validation: code unique (use repo GetByCode)
	if _, err := s.repo.GetByCode(ctx, req.Code); err == nil {
		return nil, errors.New("project code already exists")
	}
	// Default is_active handled by repository
	return s.repo.Create(ctx, req)
}

func (s *projectService) GetProject(ctx context.Context, id uuid.UUID) (*ProjectResponse, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *projectService) UpdateProject(ctx context.Context, id uuid.UUID, req UpdateProjectRequest) (*ProjectResponse, error) {
	// Do not allow changing Code: UpdateProjectRequest does not include Code, so ignore if provided in payload
	return s.repo.Update(ctx, id, req)
}

// TODO: lógica de negocio (validaciones, reglas)
