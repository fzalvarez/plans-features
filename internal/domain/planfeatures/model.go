package planfeatures

// Create assignment request: feature id and dynamic value
type AssignFeatureRequest struct {
	FeatureID string      `json:"feature_id"`
	Value     interface{} `json:"value"`
}

// PlanFeatureResponse returned to clients
type PlanFeatureResponse struct {
	ID        string      `json:"id"`
	PlanID    string      `json:"plan_id"`
	ProjectID string      `json:"project_id"`
	FeatureID string      `json:"feature_id"`
	Value     interface{} `json:"value"`
}

// internal entity
type planFeatureEntity struct {
	ID        string      `db:"id"`
	PlanID    string      `db:"plan_id"`
	ProjectID string      `db:"project_id"`
	FeatureID string      `db:"feature_id"`
	Value     interface{} `db:"value"`
}
