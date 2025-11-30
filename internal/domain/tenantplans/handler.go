// TODO: handlers HTTP para tenantplans

package tenantplans

import (
	"encoding/json"
	"net/http"

	"plans-features/internal/utils"

	"github.com/go-chi/chi/v5"
)

type TenantPlanHandler struct {
	service TenantPlanService
}

func NewTenantPlanHandler(service TenantPlanService) *TenantPlanHandler {
	return &TenantPlanHandler{service: service}
}

// ListAssignments godoc
// @Summary List tenant assignments
// @Tags TenantPlans
// @Produce json
// @Param tenantId path string true "Tenant ID"
// @Success 200 {array} tenantplans.TenantPlanResponse
// @Router /admin/tenants/{tenantId}/assignments [get]
func (h *TenantPlanHandler) ListAssignments(w http.ResponseWriter, r *http.Request) {
	tenantID := chi.URLParam(r, "tenantId")
	as, err := h.service.ListAssignments(r.Context(), tenantID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, as)
}

// CreateAssignment godoc
// @Summary Create a tenant assignment
// @Tags TenantPlans
// @Accept json
// @Produce json
// @Param tenantId path string true "Tenant ID"
// @Param assignment body tenantplans.CreateTenantPlanRequest true "Create assignment"
// @Success 201 {object} tenantplans.TenantPlanResponse
// @Router /admin/tenants/{tenantId}/assignments [post]
func (h *TenantPlanHandler) CreateAssignment(w http.ResponseWriter, r *http.Request) {
	tenantID := chi.URLParam(r, "tenantId")
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
// @Tags TenantPlans
// @Accept json
// @Produce json
// @Param tenantId path string true "Tenant ID"
// @Param assignmentId path string true "Assignment ID"
// @Param assignment body tenantplans.UpdateTenantPlanRequest true "Update assignment"
// @Success 200 {object} tenantplans.TenantPlanResponse
// @Router /admin/tenants/{tenantId}/assignments/{assignmentId} [patch]
func (h *TenantPlanHandler) UpdateAssignment(w http.ResponseWriter, r *http.Request) {
	tenantID := chi.URLParam(r, "tenantId")
	assignmentID := chi.URLParam(r, "assignmentId")
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
