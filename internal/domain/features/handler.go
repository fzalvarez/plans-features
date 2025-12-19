// TODO: handlers HTTP para features

package features

import (
	"encoding/json"
	"net/http"

	"plans-features/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type FeatureHandler struct {
	service FeatureService
}

func NewFeatureHandler(service FeatureService) *FeatureHandler {
	return &FeatureHandler{service: service}
}

// ListFeatures godoc
// @Summary List features for a project
// @Description List features available for the project identified by the API key
// @Tags features
// @Produce json
// @Param X-API-Key header string true "API Key"
// @Success 200 {array} features.FeatureResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/features [get]
func (h *FeatureHandler) ListFeatures(w http.ResponseWriter, r *http.Request) {
	projectIDStr, ok := r.Context().Value("project_id").(string)
	if !ok || projectIDStr == "" {
		utils.Error(w, http.StatusUnauthorized, "missing project context")
		return
	}

	projectID, err := uuid.Parse(projectIDStr) // ← Parse middleware string
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid project ID")
		return
	}
	fs, err := h.service.ListFeatures(r.Context(), projectID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, fs)
}

// CreateFeature godoc
// @Summary Create a feature for a project
// @Description Create a new feature for the project identified by the API key
// @Tags features
// @Accept json
// @Produce json
// @Param X-API-Key header string true "API Key"
// @Param feature body features.CreateFeatureRequest true "Create feature"
// @Success 201 {object} features.FeatureResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/features [post]
func (h *FeatureHandler) CreateFeature(w http.ResponseWriter, r *http.Request) {
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
	var req CreateFeatureRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.Code == "" || req.Name == "" || req.Type == "" {
		utils.Error(w, http.StatusBadRequest, "code, type and name are required")
		return
	}
	f, err := h.service.CreateFeature(r.Context(), projectID, req)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusCreated, f)
}

// GetFeature godoc
// @Summary Get a feature by ID
// @Description Retrieve a feature for the project identified by the API key
// @Tags features
// @Produce json
// @Param X-API-Key header string true "API Key"
// @Param featureId path string true "Feature ID"
// @Success 200 {object} features.FeatureResponse
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/features/{featureId} [get]
func (h *FeatureHandler) GetFeature(w http.ResponseWriter, r *http.Request) {
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

	featureIDStr := chi.URLParam(r, "featureId")
	if !ok || featureIDStr == "" {
		utils.Error(w, http.StatusUnauthorized, "missing project context")
		return
	}
	featureID, err := uuid.Parse(featureIDStr)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid feature ID")
		return
	}

	f, err := h.service.GetFeature(r.Context(), projectID, featureID)
	if err != nil {
		if err.Error() == "not found" {
			utils.Error(w, http.StatusNotFound, "not found")
			return
		}
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, f)
}

// UpdateFeature godoc
// @Summary Update a feature
// @Description Update fields of a feature for the project identified by the API key
// @Tags features
// @Accept json
// @Produce json
// @Param X-API-Key header string true "API Key"
// @Param featureId path string true "Feature ID"
// @Param feature body features.UpdateFeatureRequest true "Update feature"
// @Success 200 {object} features.FeatureResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/features/{featureId} [put]
func (h *FeatureHandler) UpdateFeature(w http.ResponseWriter, r *http.Request) {
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

	featureIDStr := chi.URLParam(r, "featureId")
	if !ok || featureIDStr == "" {
		utils.Error(w, http.StatusUnauthorized, "missing project context")
		return
	}

	featureID, err := uuid.Parse(featureIDStr)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid feature ID")
		return
	}
	var req UpdateFeatureRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	f, err := h.service.UpdateFeature(r.Context(), projectID, featureID, req)
	if err != nil {
		if err.Error() == "not found" {
			utils.Error(w, http.StatusNotFound, "not found")
			return
		}
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, f)
}
