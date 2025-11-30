// TODO: struct Plan y DTOs

package plans

type CreatePlanRequest struct {
	Code        string
	Name        string
	Description string
	IsActive    bool
	IsDefault   bool
	Limits      map[string]interface{}
}

type UpdatePlanRequest struct {
	Name        *string
	Description *string
	IsActive    *bool
	IsDefault   *bool
	Limits      map[string]interface{}
}

type PlanResponse struct {
	ID          string
	ProjectID   string
	Code        string
	Name        string
	Description string
	IsActive    bool
	IsDefault   bool
	Limits      map[string]interface{}
}
