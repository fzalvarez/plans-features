package tenantplans

// CreateTenantPlanRequest representa la petición para asignar un plan a un tenant
type CreateTenantPlanRequest struct {
	ProjectCode string
	PlanCode    string
}

// UpdateTenantPlanRequest representa la petición para actualizar la asignación
type UpdateTenantPlanRequest struct {
	PlanCode *string
}

// TenantPlanResponse representa la respuesta de una asignación
type TenantPlanResponse struct {
	ID        string
	TenantID  string
	ProjectID string
	PlanID    string
}
