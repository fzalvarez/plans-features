package planfeatures

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type PlanFeatureRepository interface {
	ListByPlan(ctx context.Context, projectID uuid.UUID, planID uuid.UUID) ([]PlanFeatureResponse, error)
	Create(ctx context.Context, projectID uuid.UUID, planID uuid.UUID, req AssignFeatureRequest) (*PlanFeatureResponse, error)
	Exists(ctx context.Context, projectID uuid.UUID, planID uuid.UUID, featureID uuid.UUID) (bool, error)
}

type planFeatureRepository struct {
	db *sql.DB
}

func NewPlanFeatureRepository(db *sql.DB) PlanFeatureRepository {
	return &planFeatureRepository{db: db}
}

func (r *planFeatureRepository) ListByPlan(ctx context.Context, projectID uuid.UUID, planID uuid.UUID) ([]PlanFeatureResponse, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, project_id, plan_id, feature_id, value_json 
         FROM plan_features 
         WHERE project_id = $1 AND plan_id = $2 
         ORDER BY id`,
		projectID, planID)
	if err != nil {
		return nil, fmt.Errorf("list plan features: %w", err)
	}
	defer rows.Close()

	var results []PlanFeatureResponse
	for rows.Next() {
		var pf PlanFeatureResponse
		var valueJSON []byte
		err := rows.Scan(&pf.ID, &pf.ProjectID, &pf.PlanID, &pf.FeatureID, &valueJSON)
		if err != nil {
			return nil, fmt.Errorf("scan plan feature: %w", err)
		}

		// JSONB → interface{}
		if err := json.Unmarshal(valueJSON, &pf.Value); err != nil {
			return nil, fmt.Errorf("unmarshal value_json: %w", err)
		}

		results = append(results, pf)
	}
	return results, rows.Err()
}

func (r *planFeatureRepository) Exists(ctx context.Context, projectID uuid.UUID, planID uuid.UUID, featureID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS (
            SELECT 1 FROM plan_features 
            WHERE project_id = $1 AND plan_id = $2 AND feature_id = $3
        )`, projectID, planID, featureID).
		Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("check exists: %w", err)
	}
	return exists, nil
}

func (r *planFeatureRepository) Create(ctx context.Context, projectID uuid.UUID, planID uuid.UUID, req AssignFeatureRequest) (*PlanFeatureResponse, error) {
	// Verificar unique constraint primero
	exists, err := r.Exists(ctx, projectID, planID, req.FeatureID)
	if err != nil {
		return nil, fmt.Errorf("check exists: %w", err)
	}
	if exists {
		return nil, errors.New("feature already assigned to this plan")
	}

	id := uuid.New()

	// interface{} → JSONB
	valueJSON, err := json.Marshal(req.Value)
	if err != nil {
		return nil, fmt.Errorf("marshal value: %w", err)
	}

	pf := &PlanFeatureResponse{}
	err = r.db.QueryRowContext(ctx,
		`INSERT INTO plan_features (id, project_id, plan_id, feature_id, value_json)
         VALUES ($1, $2, $3, $4, $5)
         RETURNING id, project_id, plan_id, feature_id, value_json`,
		id, projectID, planID, req.FeatureID, valueJSON).
		Scan(&pf.ID, &pf.ProjectID, &pf.PlanID, &pf.FeatureID, &valueJSON)

	if err != nil {
		return nil, fmt.Errorf("create plan feature: %w", err)
	}

	// Re-unmarshal para Value interface{}
	if err := json.Unmarshal(valueJSON, &pf.Value); err != nil {
		return nil, fmt.Errorf("unmarshal response value: %w", err)
	}

	return pf, nil
}
