package projects

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type ProjectRepository interface {
	List(ctx context.Context) ([]ProjectResponse, error)
	Create(ctx context.Context, req CreateProjectRequest) (*ProjectResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*ProjectResponse, error)
	GetByCode(ctx context.Context, code string) (*ProjectResponse, error)
	Update(ctx context.Context, id uuid.UUID, req UpdateProjectRequest) (*ProjectResponse, error)
}

type projectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) ProjectRepository {
	return &projectRepository{
		db: db,
	}
}

func (r *projectRepository) List(ctx context.Context) ([]ProjectResponse, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id::text, code, name, description, is_active 
         FROM projects 
         WHERE is_active = true 
         ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	defer rows.Close()

	var projects []ProjectResponse
	for rows.Next() {
		proj := &Project{}
		var desc sql.NullString
		if err := rows.Scan(&proj.ID, &proj.Code, &proj.Name, &desc, &proj.IsActive); err != nil {
			return nil, fmt.Errorf("scan project: %w", err)
		}
		proj.Description = nullStringToPtr(desc)
		projects = append(projects, *ToResponse(proj))
	}
	return projects, rows.Err()
}

func (r *projectRepository) Create(ctx context.Context, req CreateProjectRequest) (*ProjectResponse, error) {
	id := uuid.New()

	// Normalizar code
	code := strings.ToLower(strings.TrimSpace(req.Code))

	// Default is_active
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	// Convertir description string → *string (nullable)
	var description *string
	if req.Description != "" {
		description = &req.Description
	}

	// INSERT con RETURNING (transacción atómica)
	proj := &Project{}
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO projects (id, code, name, description, is_active) 
         VALUES ($1, $2, $3, $4, $5) 
         RETURNING id, code, name, description, is_active, created_at, updated_at`,
		id, code, req.Name, description, isActive).
		Scan(&proj.ID, &proj.Code, &proj.Name, &proj.Description,
			&proj.IsActive, &proj.CreatedAt, &proj.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}

	return ToResponse(proj), nil
}

func (r *projectRepository) GetByID(ctx context.Context, id uuid.UUID) (*ProjectResponse, error) {
	proj := &Project{}
	var desc sql.NullString

	err := r.db.QueryRowContext(ctx,
		`SELECT id::text, code, name, description, is_active, created_at, updated_at 
         FROM projects WHERE id = $1`, id).
		Scan(&proj.ID, &proj.Code, &proj.Name, &desc, &proj.IsActive,
			&proj.CreatedAt, &proj.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("project not found")
	}
	if err != nil {
		return nil, fmt.Errorf("get project by ID: %w", err)
	}

	if desc.Valid {
		proj.Description = &desc.String
	}
	return ToResponse(proj), nil
}

func (r *projectRepository) GetByCode(ctx context.Context, code string) (*ProjectResponse, error) {
	code = strings.ToLower(strings.TrimSpace(code))

	proj := &Project{}
	var desc sql.NullString
	err := r.db.QueryRowContext(ctx,
		`SELECT id, code, name, description, is_active, created_at, updated_at 
         FROM projects WHERE LOWER(TRIM(code)) = $1`, code).
		Scan(&proj.ID, &proj.Code, &proj.Name, &desc, &proj.IsActive,
			&proj.CreatedAt, &proj.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("project not found")
	}
	if err != nil {
		return nil, fmt.Errorf("get project by code: %w", err)
	}

	if desc.Valid {
		proj.Description = &desc.String
	}
	return ToResponse(proj), nil
}

func (r *projectRepository) Update(ctx context.Context, id uuid.UUID, req UpdateProjectRequest) (*ProjectResponse, error) {
	updates := []string{}
	args := []interface{}{}
	argIdx := 1

	if req.Name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", argIdx))
		args = append(args, *req.Name)
		argIdx++
	}
	if req.Description != nil {
		updates = append(updates, fmt.Sprintf("description = $%d", argIdx))
		args = append(args, req.Description)
		argIdx++
	}
	if req.IsActive != nil {
		updates = append(updates, fmt.Sprintf("is_active = $%d", argIdx))
		args = append(args, *req.IsActive)
		argIdx++
	}

	if len(updates) == 0 {
		return r.GetByID(ctx, id)
	}

	// Agregar WHERE id
	updates = append(updates, fmt.Sprintf("id = $%d", argIdx))
	args = append(args, id)

	query := fmt.Sprintf(
		`UPDATE projects 
         SET %s, updated_at = NOW() 
         WHERE %s 
         RETURNING id::text, code, name, description, is_active, created_at, updated_at`,
		strings.Join(updates[:len(updates)-1], ", "),
		updates[len(updates)-1])

	proj := &Project{}
	var desc sql.NullString
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&proj.ID, &proj.Code, &proj.Name, &desc, &proj.IsActive,
		&proj.CreatedAt, &proj.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("project not found")
	}
	if err != nil {
		return nil, fmt.Errorf("update project: %w", err)
	}

	if desc.Valid {
		proj.Description = &desc.String
	}
	return ToResponse(proj), nil
}
