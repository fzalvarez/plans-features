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
// @Tags Projects
// @Produce json
// @Success 200 {array} projects.ProjectResponse
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
// @Tags Projects
// @Accept json
// @Produce json
// @Param project body projects.CreateProjectRequest true "Create project"
// @Success 201 {object} projects.ProjectResponse
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
		utils.Error(w, http.StatusInternalServerError, "internal error")
		return
	}
	utils.JSON(w, http.StatusCreated, p)
}

// GetProject godoc
// @Summary Get a project by ID
// @Tags Projects
// @Produce json
// @Param projectId path string true "Project ID"
// @Success 200 {object} projects.ProjectResponse
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
// @Tags Projects
// @Accept json
// @Produce json
// @Param projectId path string true "Project ID"
// @Param project body projects.UpdateProjectRequest true "Update project"
// @Success 200 {object} projects.ProjectResponse
// @Router /admin/projects/{projectId} [patch]
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
