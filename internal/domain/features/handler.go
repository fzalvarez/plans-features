// TODO: handlers HTTP para features

package features

import (
	"encoding/json"
	"net/http"

	"plans-features/internal/utils"

	"github.com/go-chi/chi/v5"
)

type FeatureHandler struct {
	service FeatureService
}

func NewFeatureHandler(service FeatureService) *FeatureHandler {
	return &FeatureHandler{service: service}
}

// ListFeatures godoc
// @Summary List features for a project
// @Tags Features
// @Produce json
// @Param projectId path string true "Project ID"
// @Success 200 {array} features.FeatureResponse
// @Router /admin/projects/{projectId}/features [get]
func (h *FeatureHandler) ListFeatures(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectId")
	fs, err := h.service.ListFeatures(r.Context(), projectID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, fs)
}

// CreateFeature godoc
// @Summary Create a feature for a project
// @Tags Features
// @Accept json
// @Produce json
// @Param projectId path string true "Project ID"
// @Param feature body features.CreateFeatureRequest true "Create feature"
// @Success 201 {object} features.FeatureResponse
// @Router /admin/projects/{projectId}/features [post]
func (h *FeatureHandler) CreateFeature(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectId")
	var req CreateFeatureRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.Code == "" || req.Name == "" {
		utils.Error(w, http.StatusBadRequest, "code and name are required")
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
// @Tags Features
// @Produce json
// @Param projectId path string true "Project ID"
// @Param featureId path string true "Feature ID"
// @Success 200 {object} features.FeatureResponse
// @Router /admin/projects/{projectId}/features/{featureId} [get]
func (h *FeatureHandler) GetFeature(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectId")
	featureID := chi.URLParam(r, "featureId")
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
// @Tags Features
// @Accept json
// @Produce json
// @Param projectId path string true "Project ID"
// @Param featureId path string true "Feature ID"
// @Param feature body features.UpdateFeatureRequest true "Update feature"
// @Success 200 {object} features.FeatureResponse
// @Router /admin/projects/{projectId}/features/{featureId} [patch]
func (h *FeatureHandler) UpdateFeature(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectId")
	featureID := chi.URLParam(r, "featureId")
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
