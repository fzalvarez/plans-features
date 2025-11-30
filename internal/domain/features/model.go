// TODO: struct Feature y DTOs

package features

type CreateFeatureRequest struct {
	Code        string
	Name        string
	Description string
	IsActive    bool
}

type UpdateFeatureRequest struct {
	Name        *string
	Description *string
	IsActive    *bool
}

type FeatureResponse struct {
	ID          string
	ProjectID   string
	Code        string
	Name        string
	Description string
	IsActive    bool
}
