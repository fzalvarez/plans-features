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
	ID        uuid.UUID `json:"id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	ProjectID uuid.UUID `json:"project_id"`
	PlanID    uuid.UUID `json:"plan_id"`
}

// API request para asignar plan usando project_id desde context
type PlanAssignRequest struct {
	PlanID uuid.UUID `json:"plan_id"`
}

// Interno para DB (Scan)
type TenantPlan struct {
	ID        uuid.UUID `db:"id"`
	TenantID  uuid.UUID `db:"tenant_id"`
	ProjectID uuid.UUID `db:"project_id"`
	PlanID    uuid.UUID `db:"plan_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
