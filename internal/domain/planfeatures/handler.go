package planfeatures

import (
	"encoding/json"
	"net/http"

	"plans-features/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PlanFeatureHandler struct {
	service PlanFeatureService
}

func NewPlanFeatureHandler(s PlanFeatureService) *PlanFeatureHandler {
	return &PlanFeatureHandler{service: s}
}

// List godoc
// @Summary List features assigned to a plan
// @Description List features and their values assigned to the given plan
// @Tags planfeatures
// @Produce json
// @Param X-API-Key header string true "API Key"
// @Param planId path string true "Plan ID"
// @Success 200 {array} planfeatures.PlanFeatureResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/plans/{planId}/features [get]
func (h *PlanFeatureHandler) List(w http.ResponseWriter, r *http.Request) {
	projectIDStr, ok := r.Context().Value("project_id").(string)
	if !ok || projectIDStr == "" {
		utils.Error(w, http.StatusUnauthorized, "missing project context")
		return
	}

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid project ID format")
		return
	}

	planIDStr := chi.URLParam(r, "planId")
	if !ok || planIDStr == "" {
		utils.Error(w, http.StatusUnauthorized, "missing project context")
		return
	}

	planID, err := uuid.Parse(planIDStr)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid plan ID format")
		return
	}

	res, err := h.service.ListByPlan(r.Context(), projectID, planID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, res)
}

// Assign godoc
// @Summary Assign a feature to a plan
// @Description Assign a feature to a plan with a value. Feature and plan must belong to the same project.
// @Tags planfeatures
// @Accept json
// @Produce json
// @Param X-API-Key header string true "API Key"
// @Param planId path string true "Plan ID"
// @Param body body planfeatures.AssignFeatureRequest true "Assign feature"
// @Success 201 {object} planfeatures.PlanFeatureResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/plans/{planId}/features [post]
func (h *PlanFeatureHandler) Assign(w http.ResponseWriter, r *http.Request) {
	projectIDStr, ok := r.Context().Value("project_id").(string)
	if !ok || projectIDStr == "" {
		utils.Error(w, http.StatusUnauthorized, "missing project context")
		return
	}
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid project ID format")
		return
	}

	planIDStr := chi.URLParam(r, "planId")
	planID, err := uuid.Parse(planIDStr)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid plan ID format")
		return
	}

	var req AssignFeatureRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.FeatureID == uuid.Nil {
		utils.Error(w, http.StatusBadRequest, "feature_id is required")
		return
	}
	res, err := h.service.AssignFeature(r.Context(), projectID, planID, req)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.JSON(w, http.StatusCreated, res)
}
