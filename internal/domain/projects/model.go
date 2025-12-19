package projects

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID `db:"id"`
	Code        string    `db:"code"`
	Name        string    `db:"name"`
	Description *string   `db:"description"`
	IsActive    bool      `db:"is_active"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type CreateProjectRequest struct {
	Code        string `json:"code" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

type UpdateProjectRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

type ProjectResponse struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
}

func ToResponse(proj *Project) *ProjectResponse {
	resp := &ProjectResponse{
		ID:       proj.ID,
		Code:     proj.Code,
		Name:     proj.Name,
		IsActive: proj.IsActive,
	}
	if proj.Description != nil {
		resp.Description = *proj.Description
	}
	return resp
}

func ToCreateRequest(req *CreateProjectRequest) *Project {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	var desc *string
	if req.Description != "" {
		desc = &req.Description
	}
	return &Project{
		Code:        req.Code,
		Name:        req.Name,
		Description: desc,
		IsActive:    isActive,
	}
}

func nullStringToPtr(ns sql.NullString) *string {
	if ns.Valid && ns.String != "" {
		s := ns.String
		return &s
	}
	return nil
}
