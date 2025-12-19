package plans

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type PlanRepository interface {
	List(ctx context.Context, projectID uuid.UUID) ([]PlanResponse, error)
	Create(ctx context.Context, projectID uuid.UUID, req CreatePlanRequest) (*PlanResponse, error)
	GetByID(ctx context.Context, projectID uuid.UUID, planID uuid.UUID) (*PlanResponse, error)
	Update(ctx context.Context, projectID uuid.UUID, planID uuid.UUID, req UpdatePlanRequest) (*PlanResponse, error)
}

type planRepository struct {
	db *sql.DB
}

func NewPlanRepository(db *sql.DB) PlanRepository {
	return &planRepository{db: db}
}

func normalizeCode(code string) string {
	return strings.ToLower(strings.TrimSpace(code))
}

func (r *planRepository) List(ctx context.Context, projectID uuid.UUID) ([]PlanResponse, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, project_id, code, name, description, is_active, is_default, limits_json, created_at, updated_at
         FROM plans 
         WHERE project_id = $1 AND is_active = true 
         ORDER BY is_default DESC, created_at DESC`,
		projectID)
	if err != nil {
		return nil, fmt.Errorf("list plans: %w", err)
	}
	defer rows.Close()

	var plans []PlanResponse
	for rows.Next() {
		plan := &Plan{}
		var desc sql.NullString
		var limitsJSON []byte
		if err := rows.Scan(&plan.ID, &plan.ProjectID, &plan.Code, &plan.Name,
			&desc, &plan.IsActive, &plan.IsDefault, &limitsJSON,
			&plan.CreatedAt, &plan.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan plan: %w", err)
		}

		plan.Description = nullStringToPtr(desc)
		if err := json.Unmarshal(limitsJSON, &plan.Limits); err != nil {
			return nil, fmt.Errorf("unmarshal limits: %w", err)
		}

		plans = append(plans, *ToResponse(plan))
	}
	return plans, rows.Err()
}

func (r *planRepository) Create(ctx context.Context, projectID uuid.UUID, req CreatePlanRequest) (*PlanResponse, error) {
	id := uuid.New()
	var description *string
	if req.Description != "" {
		description = &req.Description
	}

	limitsJSON, err := json.Marshal(req.Limits)
	if err != nil {
		return nil, fmt.Errorf("marshal limits: %w", err)
	}

	plan := &Plan{}
	err = r.db.QueryRowContext(ctx,
		`INSERT INTO plans (id, project_id, code, name, description, is_active, is_default, limits_json)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
         RETURNING id, project_id, code, name, description, is_active, is_default, limits_json, created_at, updated_at`,
		id, projectID, normalizeCode(req.Code), req.Name, description,
		req.IsActive, req.IsDefault, limitsJSON).
		Scan(&plan.ID, &plan.ProjectID, &plan.Code, &plan.Name, &plan.Description,
			&plan.IsActive, &plan.IsDefault, &limitsJSON, &plan.CreatedAt, &plan.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("create plan: %w", err)
	}

	if err := json.Unmarshal(limitsJSON, &plan.Limits); err != nil {
		return nil, fmt.Errorf("unmarshal limits: %w", err)
	}

	return ToResponse(plan), nil
}

func (r *planRepository) GetByID(ctx context.Context, projectID uuid.UUID, planID uuid.UUID) (*PlanResponse, error) {
	plan := &Plan{}
	var desc sql.NullString
	var limitsJSON []byte

	err := r.db.QueryRowContext(ctx,
		`SELECT id, project_id, code, name, description, is_active, is_default, limits_json, created_at, updated_at
         FROM plans 
         WHERE project_id = $1 AND id = $2`,
		projectID, planID).
		Scan(&plan.ID, &plan.ProjectID, &plan.Code, &plan.Name,
			&desc, &plan.IsActive, &plan.IsDefault, &limitsJSON,
			&plan.CreatedAt, &plan.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("plan not found")
	}
	if err != nil {
		return nil, fmt.Errorf("get plan: %w", err)
	}

	plan.Description = nullStringToPtr(desc)
	if err := json.Unmarshal(limitsJSON, &plan.Limits); err != nil {
		return nil, fmt.Errorf("unmarshal limits: %w", err)
	}

	return ToResponse(plan), nil
}

func (r *planRepository) Update(ctx context.Context, projectID uuid.UUID, planID uuid.UUID, req UpdatePlanRequest) (*PlanResponse, error) {
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
	if req.IsDefault != nil {
		updates = append(updates, fmt.Sprintf("is_default = $%d", argIdx))
		args = append(args, *req.IsDefault)
		argIdx++
	}
	if req.Limits != nil {
		limitsJSON, err := json.Marshal(req.Limits)
		if err != nil {
			return nil, fmt.Errorf("marshal limits: %w", err)
		}
		updates = append(updates, fmt.Sprintf("limits_json = $%d", argIdx))
		args = append(args, limitsJSON)
		argIdx++
	}

	if len(updates) == 0 {
		return r.GetByID(ctx, projectID, planID)
	}

	updates = append(updates, fmt.Sprintf("id = $%d", argIdx))
	args = append(args, planID)

	query := fmt.Sprintf(
		`UPDATE plans 
         SET %s, updated_at = NOW()
         WHERE project_id = $%d AND %s
         RETURNING id, project_id, code, name, description, is_active, is_default, limits_json, created_at, updated_at`,
		strings.Join(updates[:len(updates)-1], ", "),
		argIdx, updates[len(updates)-1])

	plan := &Plan{}
	var desc sql.NullString
	var limitsJSON []byte
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&plan.ID, &plan.ProjectID, &plan.Code, &plan.Name,
		&desc, &plan.IsActive, &plan.IsDefault, &limitsJSON,
		&plan.CreatedAt, &plan.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("plan not found")
	}
	if err != nil {
		return nil, fmt.Errorf("update plan: %w", err)
	}

	plan.Description = nullStringToPtr(desc)
	if err := json.Unmarshal(limitsJSON, &plan.Limits); err != nil {
		return nil, fmt.Errorf("unmarshal limits: %w", err)
	}

	return ToResponse(plan), nil
}

// Helpers
func nullStringToPtr(ns sql.NullString) *string {
	if ns.Valid && ns.String != "" {
		s := ns.String
		return &s
	}
	return nil
}
