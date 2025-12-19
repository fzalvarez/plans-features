package plans

import (
	"time"

	"github.com/google/uuid"
)

// Para DB (Scan)
type Plan struct {
	ID          uuid.UUID              `db:"id"`
	ProjectID   uuid.UUID              `db:"project_id"`
	Code        string                 `db:"code"`
	Name        string                 `db:"name"`
	Description *string                `db:"description"`
	IsActive    bool                   `db:"is_active"`
	IsDefault   bool                   `db:"is_default"`
	Limits      map[string]interface{} `db:"limits"`
	CreatedAt   time.Time              `db:"created_at"`
	UpdatedAt   time.Time              `db:"updated_at"`
}

type CreatePlanRequest struct {
	Code        string                 `json:"code"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	IsActive    bool                   `json:"is_active"`
	IsDefault   bool                   `json:"is_default"`
	Limits      map[string]interface{} `json:"limits,omitempty"`
}

type UpdatePlanRequest struct {
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	IsActive    *bool                  `json:"is_active,omitempty"`
	IsDefault   *bool                  `json:"is_default,omitempty"`
	Limits      map[string]interface{} `json:"limits,omitempty"`
}

type PlanResponse struct {
	ID          uuid.UUID              `json:"id"`
	ProjectID   uuid.UUID              `json:"project_id"`
	Code        string                 `json:"code"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	IsActive    bool                   `json:"is_active"`
	IsDefault   bool                   `json:"is_default"`
	Limits      map[string]interface{} `json:"limits"`
}

func ToResponse(plan *Plan) *PlanResponse {
	resp := &PlanResponse{
		ID:        plan.ID,
		ProjectID: plan.ProjectID,
		Code:      plan.Code,
		Name:      plan.Name,
		IsActive:  plan.IsActive,
		IsDefault: plan.IsDefault,
		Limits:    plan.Limits,
	}
	if plan.Description != nil {
		resp.Description = *plan.Description
	}
	return resp
}
