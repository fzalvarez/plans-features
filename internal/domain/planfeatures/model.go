package planfeatures

import "github.com/google/uuid"

// Create assignment request: feature id and dynamic value
type AssignFeatureRequest struct {
	FeatureID uuid.UUID   `json:"feature_id"`
	Value     interface{} `json:"value"`
}

// PlanFeatureResponse returned to clients
type PlanFeatureResponse struct {
	ID        uuid.UUID   `json:"id"`
	PlanID    uuid.UUID   `json:"plan_id"`
	ProjectID uuid.UUID   `json:"project_id"`
	FeatureID uuid.UUID   `json:"feature_id"`
	Value     interface{} `json:"value"`
}

// internal entity
type planFeatureEntity struct {
	ID        uuid.UUID   `db:"id"`
	PlanID    uuid.UUID   `db:"plan_id"`
	ProjectID uuid.UUID   `db:"project_id"`
	FeatureID uuid.UUID   `db:"feature_id"`
	Value     interface{} `db:"value"`
}
