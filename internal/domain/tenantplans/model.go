package tenantplans

import (
	"time"

	"github.com/google/uuid"
)

// CreateTenantPlanRequest representa la petición para asignar un plan a un tenant (admin)
type CreateTenantPlanRequest struct {
	ProjectCode string `json:"project_code"`
	PlanCode    string `json:"plan_code"`
}

// UpdateTenantPlanRequest representa la petición para actualizar la asignación (admin)
type UpdateTenantPlanRequest struct {
	PlanCode *string `json:"plan_code,omitempty"`
}

// TenantPlanResponse representa la respuesta de una asignación
type TenantPlanResponse struct {
	ID        string `json:"id"`
	TenantID  string `json:"tenant_id"`
	ProjectID string `json:"project_id"`
	PlanID    string `json:"plan_id"`
}

// API request para asignar plan usando project_id desde context
type PlanAssignRequest struct {
	PlanID string `json:"plan_id"`
}

// Interno para DB (Scan)
type TenantPlan struct {
	ID        uuid.UUID `db:"id"`
	TenantID  string    `db:"tenant_id"`
	ProjectID uuid.UUID `db:"project_id"`
	PlanID    uuid.UUID `db:"plan_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
