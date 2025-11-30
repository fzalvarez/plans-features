package plans

import (
	"encoding/json"
	"net/http"

	"plans-features/internal/utils"

	"github.com/go-chi/chi/v5"
)

type PlanHandler struct {
	service PlanService
}

func NewPlanHandler(service PlanService) *PlanHandler {
	return &PlanHandler{service: service}
}

// ListPlans godoc
// @Summary List plans for a project
// @Tags Plans
// @Produce json
// @Param projectId path string true "Project ID"
// @Success 200 {array} plans.PlanResponse
// @Router /admin/projects/{projectId}/plans [get]
func (h *PlanHandler) ListPlans(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectId")
	ps, err := h.service.ListPlans(r.Context(), projectID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, ps)
}

// CreatePlan godoc
// @Summary Create a plan for a project
// @Tags Plans
// @Accept json
// @Produce json
// @Param projectId path string true "Project ID"
// @Param plan body plans.CreatePlanRequest true "Create plan"
// @Success 201 {object} plans.PlanResponse
// @Router /admin/projects/{projectId}/plans [post]
func (h *PlanHandler) CreatePlan(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectId")
	var req CreatePlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.Code == "" || req.Name == "" {
		utils.Error(w, http.StatusBadRequest, "code and name are required")
		return
	}
	p, err := h.service.CreatePlan(r.Context(), projectID, req)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusCreated, p)
}

// GetPlan godoc
// @Summary Get a plan by ID
// @Tags Plans
// @Produce json
// @Param projectId path string true "Project ID"
// @Param planId path string true "Plan ID"
// @Success 200 {object} plans.PlanResponse
// @Router /admin/projects/{projectId}/plans/{planId} [get]
func (h *PlanHandler) GetPlan(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectId")
	planID := chi.URLParam(r, "planId")
	p, err := h.service.GetPlan(r.Context(), projectID, planID)
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

// UpdatePlan godoc
// @Summary Update a plan
// @Tags Plans
// @Accept json
// @Produce json
// @Param projectId path string true "Project ID"
// @Param planId path string true "Plan ID"
// @Param plan body plans.UpdatePlanRequest true "Update plan"
// @Success 200 {object} plans.PlanResponse
// @Router /admin/projects/{projectId}/plans/{planId} [patch]
func (h *PlanHandler) UpdatePlan(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectId")
	planID := chi.URLParam(r, "planId")
	var req UpdatePlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	p, err := h.service.UpdatePlan(r.Context(), projectID, planID, req)
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
