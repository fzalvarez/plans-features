package tenantplans

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type TenantPlanRepository interface {
	ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]TenantPlanResponse, error)
	Create(ctx context.Context, tenantID uuid.UUID, req CreateTenantPlanRequest) (*TenantPlanResponse, error)
	Update(ctx context.Context, tenantID uuid.UUID, assignmentID uuid.UUID, req UpdateTenantPlanRequest) (*TenantPlanResponse, error)
	GetByTenantAndProject(ctx context.Context, tenantID uuid.UUID, projectID uuid.UUID) (*TenantPlanResponse, error)
	UpsertByTenantAndProject(ctx context.Context, tenantID uuid.UUID, projectID uuid.UUID, planID uuid.UUID) (*TenantPlanResponse, error)
}

type tenantPlanRepository struct {
	db *sql.DB
}

func NewTenantPlanRepository(db *sql.DB) TenantPlanRepository {
	return &tenantPlanRepository{db: db}
}

func (r *tenantPlanRepository) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]TenantPlanResponse, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, tenant_id, project_id, plan_id, created_at, updated_at
         FROM tenant_plans 
         WHERE tenant_id = $1 
         ORDER BY created_at DESC`,
		tenantID)
	if err != nil {
		return nil, fmt.Errorf("list tenant plans: %w", err)
	}
	defer rows.Close()

	var results []TenantPlanResponse
	for rows.Next() {
		var tp TenantPlanResponse
		var createdAt, updatedAt time.Time
		err := rows.Scan(&tp.ID, &tp.TenantID, &tp.ProjectID, &tp.PlanID, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan tenant plan: %w", err)
		}
		results = append(results, tp)
	}
	return results, rows.Err()
}

func (r *tenantPlanRepository) Create(ctx context.Context, tenantID uuid.UUID, req CreateTenantPlanRequest) (*TenantPlanResponse, error) {
	id := uuid.New()

	tp := &TenantPlanResponse{}
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO tenant_plans (id, tenant_id, project_id, plan_id)
         VALUES ($1, $2, $3, $4)
         RETURNING id, tenant_id, project_id, plan_id`,
		id, tenantID, req.ProjectCode, req.PlanCode).
		Scan(&tp.ID, &tp.TenantID, &tp.ProjectID, &tp.PlanID)

	if err != nil {
		return nil, fmt.Errorf("create tenant plan: %w", err)
	}
	return tp, nil
}

func (r *tenantPlanRepository) Update(ctx context.Context, tenantID uuid.UUID, assignmentID uuid.UUID, req UpdateTenantPlanRequest) (*TenantPlanResponse, error) {
	if req.PlanCode == nil {
		// Nada que actualizar
		return r.GetByID(ctx, assignmentID)
	}

	tp := &TenantPlanResponse{}
	err := r.db.QueryRowContext(ctx,
		`UPDATE tenant_plans 
         SET plan_id = $1, updated_at = NOW()
         WHERE id = $2 AND tenant_id = $3
         RETURNING id, tenant_id, project_id, plan_id`,
		req.PlanCode, assignmentID, tenantID).
		Scan(&tp.ID, &tp.TenantID, &tp.ProjectID, &tp.PlanID)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("tenant plan not found")
	}
	if err != nil {
		return nil, fmt.Errorf("update tenant plan: %w", err)
	}

	return tp, nil
}

func (r *tenantPlanRepository) GetByTenantAndProject(ctx context.Context, tenantID uuid.UUID, projectID uuid.UUID) (*TenantPlanResponse, error) {
	tp := &TenantPlanResponse{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, tenant_id, project_id, plan_id
         FROM tenant_plans 
         WHERE tenant_id = $1 AND project_id = $2`,
		tenantID, projectID).
		Scan(&tp.ID, &tp.TenantID, &tp.ProjectID, &tp.PlanID)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("tenant plan not found")
	}
	if err != nil {
		return nil, fmt.Errorf("get tenant plan: %w", err)
	}

	return tp, nil
}

func (r *tenantPlanRepository) UpsertByTenantAndProject(ctx context.Context, tenantID uuid.UUID, projectID uuid.UUID, planID uuid.UUID) (*TenantPlanResponse, error) {
	tp := &TenantPlanResponse{}
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO tenant_plans (tenant_id, project_id, plan_id)
         VALUES ($1, $2, $3)
         ON CONFLICT (tenant_id, project_id)
         DO UPDATE SET 
             plan_id = EXCLUDED.plan_id,
             updated_at = NOW()
         RETURNING id, tenant_id, project_id, plan_id`,
		tenantID, projectID, planID).
		Scan(&tp.ID, &tp.TenantID, &tp.ProjectID, &tp.PlanID)

	if err != nil {
		return nil, fmt.Errorf("upsert tenant plan: %w", err)
	}

	return tp, nil
}

// Helper para GetByID (si lo necesitas en otros métodos)
func (r *tenantPlanRepository) GetByID(ctx context.Context, id uuid.UUID) (*TenantPlanResponse, error) {
	tp := &TenantPlanResponse{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, tenant_id, project_id, plan_id
         FROM tenant_plans WHERE id = $1`,
		id).
		Scan(&tp.ID, &tp.TenantID, &tp.ProjectID, &tp.PlanID)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("tenant plan not found")
	}
	if err != nil {
		return nil, fmt.Errorf("get tenant plan: %w", err)
	}

	return tp, nil
}
