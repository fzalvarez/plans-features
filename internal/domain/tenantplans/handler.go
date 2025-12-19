// TODO: handlers HTTP para tenantplans

package tenantplans

import (
	"encoding/json"
	"net/http"

	"plans-features/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type TenantPlanHandler struct {
	service TenantPlanService
}

func NewTenantPlanHandler(service TenantPlanService) *TenantPlanHandler {
	return &TenantPlanHandler{service: service}
}

// ListAssignments godoc
// @Summary List tenant assignments
// @Description Admin: list all plan assignments for a tenant
// @Tags tenantplans
// @Produce json
// @Param tenantId path string true "Tenant ID"
// @Success 200 {array} tenantplans.TenantPlanResponse
// @Failure 500 {object} map[string]string
// @Router /admin/tenants/{tenantId}/assignments [get]
func (h *TenantPlanHandler) ListAssignments(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := chi.URLParam(r, "tenantId")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid tenant ID")
		return
	}

	as, err := h.service.ListAssignments(r.Context(), tenantID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, as)
}

// CreateAssignment godoc
// @Summary Create a tenant assignment
// @Description Admin: create a tenant assignment for a given project and plan
// @Tags tenantplans
// @Accept json
// @Produce json
// @Param tenantId path string true "Tenant ID"
// @Param assignment body tenantplans.CreateTenantPlanRequest true "Create assignment"
// @Success 201 {object} tenantplans.TenantPlanResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/tenants/{tenantId}/assignments [post]
func (h *TenantPlanHandler) CreateAssignment(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := chi.URLParam(r, "tenantId")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid tenant ID")
		return
	}

	var req CreateTenantPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.ProjectCode == "" || req.PlanCode == "" {
		utils.Error(w, http.StatusBadRequest, "project_code and plan_code are required")
		return
	}
	p, err := h.service.CreateAssignment(r.Context(), tenantID, req)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusCreated, p)
}

// UpdateAssignment godoc
// @Summary Update a tenant assignment
// @Description Admin: update the plan for a tenant assignment
// @Tags tenantplans
// @Accept json
// @Produce json
// @Param tenantId path string true "Tenant ID"
// @Param assignmentId path string true "Assignment ID"
// @Param assignment body tenantplans.UpdateTenantPlanRequest true "Update assignment"
// @Success 200 {object} tenantplans.TenantPlanResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/tenants/{tenantId}/assignments/{assignmentId} [patch]
func (h *TenantPlanHandler) UpdateAssignment(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := chi.URLParam(r, "tenantId")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid tenant ID")
		return
	}

	assignmentIDStr := chi.URLParam(r, "assignmentId")
	assignmentID, err := uuid.Parse(assignmentIDStr)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid assignment ID")
		return
	}
	var req UpdateTenantPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	p, err := h.service.UpdateAssignment(r.Context(), tenantID, assignmentID, req)
	if err != nil {
		if err.Error() == "not found" {
			utils.Error(w, http.StatusNotFound, "not found")
			return
		}
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, p)
}

// API: GetTenantPlan godoc
// @Summary Get effective tenant plan for a project
// @Description Returns the assigned plan for tenant+project, falling back to project's default plan
// @Tags tenantplans
// @Produce json
// @Param X-API-Key header string true "API Key"
// @Param tenantId path string true "Tenant ID"
// @Success 200 {object} tenantplans.TenantPlanResponse
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tenants/{tenantId}/plan [get]
func (h *TenantPlanHandler) GetTenantPlan(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := chi.URLParam(r, "tenantId")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid tenant ID")
		return
	}

	projectIDStr, ok := r.Context().Value("project_id").(string)
	if !ok || projectIDStr == "" {
		utils.Error(w, http.StatusUnauthorized, "missing project context")
		return
	}
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid project ID")
		return
	}
	p, err := h.service.GetTenantPlan(r.Context(), tenantID, projectID)
	if err != nil {
		utils.Error(w, http.StatusNotFound, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, p)
}

// API: AssignTenantPlan godoc
// @Summary Assign a plan to tenant for the project from context
// @Description Assigns or updates a tenant's plan for the project identified by the API key
// @Tags tenantplans
// @Accept json
// @Produce json
// @Param X-API-Key header string true "API Key"
// @Param tenantId path string true "Tenant ID"
// @Param body body tenantplans.PlanAssignRequest true "Assign plan"
// @Success 200 {object} tenantplans.TenantPlanResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tenants/{tenantId}/plan [post]
func (h *TenantPlanHandler) AssignTenantPlan(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := chi.URLParam(r, "tenantId")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid tenant ID")
		return
	}
	projectIDStr, ok := r.Context().Value("project_id").(string)
	if !ok || projectIDStr == "" {
		utils.Error(w, http.StatusUnauthorized, "missing project context")
		return
	}
	projectID, err := uuid.Parse(projectIDStr)
	var req PlanAssignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.PlanID == uuid.Nil {
		utils.Error(w, http.StatusBadRequest, "plan_id is required")
		return
	}
	p, err := h.service.AssignTenantPlan(r.Context(), tenantID, projectID, req.PlanID)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, p)
}
