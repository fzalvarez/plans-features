package projects

import (
	"encoding/json"
	"net/http"

	"plans-features/internal/utils"

	"github.com/go-chi/chi/v5"
)

type ProjectHandler struct {
	service ProjectService
}

func NewProjectHandler(service ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

// ListProjects godoc
// @Summary List projects
// @Description Retrieve the list of projects
// @Tags projects
// @Produce json
// @Param X-API-Key header string true "API Key"
// @Success 200 {array} projects.ProjectResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/projects [get]
func (h *ProjectHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	ps, err := h.service.ListProjects(r.Context())
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "internal error")
		return
	}
	utils.JSON(w, http.StatusOK, ps)
}

// CreateProject godoc
// @Summary Create a project
// @Description Create a new project with unique code
// @Tags projects
// @Accept json
// @Produce json
// @Param X-API-Key header string true "API Key"
// @Param project body projects.CreateProjectRequest true "Create project"
// @Success 201 {object} projects.ProjectResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/projects [post]
func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.Code == "" || req.Name == "" {
		utils.Error(w, http.StatusBadRequest, "code and name are required")
		return
	}

	p, err := h.service.CreateProject(r.Context(), req)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusCreated, p)
}

// GetProject godoc
// @Summary Get a project by ID
// @Description Retrieve a project by its ID
// @Tags projects
// @Produce json
// @Param X-API-Key header string true "API Key"
// @Param projectId path string true "Project ID"
// @Success 200 {object} projects.ProjectResponse
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/projects/{projectId} [get]
func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "projectId")
	p, err := h.service.GetProject(r.Context(), id)
	if err != nil {
		if err.Error() == "not found" {
			utils.Error(w, http.StatusNotFound, "not found")
			return
		}
		utils.Error(w, http.StatusInternalServerError, "internal error")
		return
	}
	utils.JSON(w, http.StatusOK, p)
}

// UpdateProject godoc
// @Summary Update a project
// @Description Update name/description/is_active of a project (code cannot be changed)
// @Tags projects
// @Accept json
// @Produce json
// @Param X-API-Key header string true "API Key"
// @Param projectId path string true "Project ID"
// @Param project body projects.UpdateProjectRequest true "Update project"
// @Success 200 {object} projects.ProjectResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/projects/{projectId} [put]
func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "projectId")
	var req UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	p, err := h.service.UpdateProject(r.Context(), id, req)
	if err != nil {
		if err.Error() == "not found" {
			utils.Error(w, http.StatusNotFound, "not found")
			return
		}
		utils.Error(w, http.StatusInternalServerError, "internal error")
		return
	}
	utils.JSON(w, http.StatusOK, p)
}
