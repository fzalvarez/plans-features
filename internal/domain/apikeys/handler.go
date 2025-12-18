package apikeys

import (
	"encoding/json"
	"net/http"

	"plans-features/internal/utils"

	"github.com/go-chi/chi/v5"
)

type APIKeyHandler struct {
	service APIKeyService
}

func NewAPIKeyHandler(s APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{service: s}
}

// CreateKey godoc
// @Summary Create API key for a project
// @Description Create a new API key for the specified project (admin)
// @Tags apikeys
// @Accept json
// @Produce json
// @Param X-API-Key header string true "Admin API Key"
// @Param projectId path string true "Project ID"
// @Success 201 {object} apikeys.CreateAPIKeyResult
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/projects/{projectId}/apikeys [post]
func (h *APIKeyHandler) CreateKey(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectId")
	res, err := h.service.CreateKey(r.Context(), projectID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	// raw key returned only once
	utils.JSON(w, http.StatusCreated, res)
}

// RotateKey godoc
// @Summary Rotate API key for a project
// @Description Rotate the API key for the specified project (admin)
// @Tags apikeys
// @Accept json
// @Produce json
// @Param X-API-Key header string true "Admin API Key"
// @Param projectId path string true "Project ID"
// @Success 200 {object} apikeys.CreateAPIKeyResult
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/projects/{projectId}/apikeys/rotate [post]
func (h *APIKeyHandler) RotateKey(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectId")
	res, err := h.service.RotateKey(r.Context(), projectID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, res)
}

// RevokeKey godoc
// @Summary Revoke API key(s) for a project
// @Description Revoke API key(s) for a project (admin). Optionally provide key_prefix in body.
// @Tags apikeys
// @Accept json
// @Produce json
// @Param X-API-Key header string true "Admin API Key"
// @Param projectId path string true "Project ID"
// @Param body body apikeys.RevokeAPIKeyRequest false "Revoke options"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/projects/{projectId}/apikeys/revoke [post]
func (h *APIKeyHandler) RevokeKey(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectId")
	var req RevokeAPIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// allow empty body
	}
	var prefix *string
	if req.KeyPrefix != nil {
		prefix = req.KeyPrefix
	}
	if err := h.service.RevokeKey(r.Context(), projectID, prefix); err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
