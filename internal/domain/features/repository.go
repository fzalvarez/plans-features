package features

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type FeatureRepository interface {
	List(ctx context.Context, projectID uuid.UUID) ([]FeatureResponse, error)
	Create(ctx context.Context, projectID uuid.UUID, req CreateFeatureRequest) (*FeatureResponse, error)
	GetByID(ctx context.Context, projectID uuid.UUID, featureID uuid.UUID) (*FeatureResponse, error)
	Update(ctx context.Context, projectID uuid.UUID, featureID uuid.UUID, req UpdateFeatureRequest) (*FeatureResponse, error)
}

type featureRepository struct {
	db *sql.DB
}

func NewFeatureRepository(db *sql.DB) FeatureRepository {
	return &featureRepository{db: db}
}

func normalizeCode(code string) string {
	return strings.ToLower(strings.TrimSpace(code))
}

func (r *featureRepository) List(ctx context.Context, projectID uuid.UUID) ([]FeatureResponse, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, project_id, code, type, name, description, is_active, created_at, updated_at
         FROM features 
         WHERE project_id = $1 AND is_active = true 
         ORDER BY created_at DESC`,
		projectID)
	if err != nil {
		return nil, fmt.Errorf("list features: %w", err)
	}
	defer rows.Close()

	var features []FeatureResponse
	for rows.Next() {
		feat := &Feature{}
		var desc sql.NullString
		if err := rows.Scan(&feat.ID, &feat.ProjectID, &feat.Code, &feat.Type,
			&feat.Name, &desc, &feat.IsActive, &feat.CreatedAt, &feat.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan feature: %w", err)
		}
		feat.Description = nullStringToPtr(desc)
		features = append(features, *ToResponse(feat))
	}
	return features, rows.Err()
}

func (r *featureRepository) Create(ctx context.Context, projectID uuid.UUID, req CreateFeatureRequest) (*FeatureResponse, error) {
	id := uuid.New()
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	var description *string
	if req.Description != "" {
		description = &req.Description
	}

	feat := &Feature{}
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO features (id, project_id, code, type, name, description, is_active)
         VALUES ($1, $2, $3, $4, $5, $6, $7)
         RETURNING id, project_id, code, type, name, description, is_active, created_at, updated_at`,
		id, projectID, normalizeCode(req.Code), req.Type, req.Name, description, isActive).
		Scan(&feat.ID, &feat.ProjectID, &feat.Code, &feat.Type, &feat.Name,
			&feat.Description, &feat.IsActive, &feat.CreatedAt, &feat.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("create feature: %w", err)
	}
	return ToResponse(feat), nil
}

func (r *featureRepository) GetByID(ctx context.Context, projectID uuid.UUID, featureID uuid.UUID) (*FeatureResponse, error) {
	feat := &Feature{}
	var desc sql.NullString
	err := r.db.QueryRowContext(ctx,
		`SELECT id, project_id, code, type, name, description, is_active, created_at, updated_at
         FROM features 
         WHERE project_id = $1 AND id = $2`,
		projectID, featureID).
		Scan(&feat.ID, &feat.ProjectID, &feat.Code, &feat.Type, &feat.Name,
			&desc, &feat.IsActive, &feat.CreatedAt, &feat.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("feature not found")
	}
	if err != nil {
		return nil, fmt.Errorf("get feature: %w", err)
	}

	if desc.Valid {
		feat.Description = &desc.String
	}
	return ToResponse(feat), nil
}

func (r *featureRepository) Update(ctx context.Context, projectID uuid.UUID, featureID uuid.UUID, req UpdateFeatureRequest) (*FeatureResponse, error) {
	// Verificar existencia primero
	if _, err := r.GetByID(ctx, projectID, featureID); err != nil {
		return nil, err
	}

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
		return r.GetByID(ctx, projectID, featureID)
	}

	updates = append(updates, fmt.Sprintf("id = $%d", argIdx))
	args = append(args, featureID)

	query := fmt.Sprintf(
		`UPDATE features 
         SET %s, updated_at = NOW()
         WHERE project_id = $%d AND %s
         RETURNING id, project_id, code, type, name, description, is_active, created_at, updated_at`,
		strings.Join(updates[:len(updates)-1], ", "),
		argIdx, updates[len(updates)-1])

	feat := &Feature{}
	var desc sql.NullString
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&feat.ID, &feat.ProjectID, &feat.Code, &feat.Type, &feat.Name,
		&desc, &feat.IsActive, &feat.CreatedAt, &feat.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("feature not found")
	}
	if err != nil {
		return nil, fmt.Errorf("update feature: %w", err)
	}

	if desc.Valid {
		feat.Description = &desc.String
	}
	return ToResponse(feat), nil
}
