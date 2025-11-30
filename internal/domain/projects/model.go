// TODO: struct Project y DTOs de request/response

package projects

type CreateProjectRequest struct {
	Code        string
	Name        string
	Description string
	IsActive    bool
}

type UpdateProjectRequest struct {
	Name        *string
	Description *string
	IsActive    *bool
}

type ProjectResponse struct {
	ID          string
	Code        string
	Name        string
	Description string
	IsActive    bool
}
