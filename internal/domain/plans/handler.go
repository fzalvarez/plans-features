package plans

import (
	"encoding/json"
	"net/http"

	"plans-features/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PlanHandler struct {
	service PlanService
}

func NewPlanHandler(service PlanService) *PlanHandler {
	return &PlanHandler{service: service}
}

// ListPlans godoc
// @Summary List plans for a project
// @Description List available plans for the project identified by the API key
// @Tags plans
// @Produce json
// @Param X-API-Key header string true "API Key"
// @Success 200 {array} plans.PlanResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/plans [get]
func (h *PlanHandler) ListPlans(w http.ResponseWriter, r *http.Request) {
	projectIDStr, ok := r.Context().Value("project_id").(string)
	if !ok || projectIDStr == "" {
		utils.Error(w, http.StatusUnauthorized, "missing project context")
		return
	}
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, "invalid project ID")
		return
	}
	ps, err := h.service.ListPlans(r.Context(), projectID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, ps)
}

// CreatePlan godoc
// @Summary Create a plan for a project
// @Description Create a new plan for the project identified by the API key
// @Tags plans
// @Accept json
// @Produce json
// @Param X-API-Key header string true "API Key"
// @Param plan body plans.CreatePlanRequest true "Create plan"
// @Success 201 {object} plans.PlanResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/plans [post]
func (h *PlanHandler) CreatePlan(w http.ResponseWriter, r *http.Request) {
	projectIDStr, ok := r.Context().Value("project_id").(string)
	if !ok || projectIDStr == "" {
		utils.Error(w, http.StatusUnauthorized, "missing project context")
		return
	}

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, "invalid project ID")
		return
	}

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
// @Description Retrieve a plan for the project identified by the API key
// @Tags plans
// @Produce json
// @Param X-API-Key header string true "API Key"
// @Param planId path string true "Plan ID"
// @Success 200 {object} plans.PlanResponse
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/plans/{planId} [get]
func (h *PlanHandler) GetPlan(w http.ResponseWriter, r *http.Request) {
	projectIDStr, ok := r.Context().Value("project_id").(string)
	if !ok || projectIDStr == "" {
		utils.Error(w, http.StatusUnauthorized, "missing project context")
		return
	}
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, "invalid project ID")
		return
	}

	planIDStr := chi.URLParam(r, "planId")
	planID, err := uuid.Parse(planIDStr)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid plan ID")
		return
	}

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
// @Description Update fields of a plan for the project identified by the API key
// @Tags plans
// @Accept json
// @Produce json
// @Param X-API-Key header string true "API Key"
// @Param planId path string true "Plan ID"
// @Param plan body plans.UpdatePlanRequest true "Update plan"
// @Success 200 {object} plans.PlanResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/plans/{planId} [put]
func (h *PlanHandler) UpdatePlan(w http.ResponseWriter, r *http.Request) {
	projectIDStr, ok := r.Context().Value("project_id").(string)
	if !ok || projectIDStr == "" {
		utils.Error(w, http.StatusUnauthorized, "missing project context")
		return
	}
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, "invalid project ID")
		return
	}

	planIDStr := chi.URLParam(r, "planId")
	planID, err := uuid.Parse(planIDStr)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid plan ID")
		return
	}

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
