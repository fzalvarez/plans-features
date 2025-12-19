package features

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Para DB (Scan interno)
type Feature struct {
	ID          uuid.UUID `db:"id"`
	ProjectID   uuid.UUID `db:"project_id"`
	Code        string    `db:"code"`
	Type        string    `db:"type"`
	Name        string    `db:"name"`
	Description *string   `db:"description"`
	IsActive    bool      `db:"is_active"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type CreateFeatureRequest struct {
	Code        string `json:"code"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

type UpdateFeatureRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

type FeatureResponse struct {
	ID          uuid.UUID `json:"id"`
	ProjectID   uuid.UUID `json:"project_id"`
	Code        string    `json:"code"`
	Type        string    `json:"type"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
}

func ToResponse(feat *Feature) *FeatureResponse {
	resp := &FeatureResponse{
		ID:        feat.ID,
		ProjectID: feat.ProjectID,
		Code:      feat.Code,
		Type:      feat.Type,
		Name:      feat.Name,
		IsActive:  feat.IsActive,
	}
	if feat.Description != nil {
		resp.Description = *feat.Description
	}
	return resp
}

func nullStringToPtr(ns sql.NullString) *string {
	if ns.Valid && ns.String != "" {
		s := ns.String
		return &s
	}
	return nil
}
