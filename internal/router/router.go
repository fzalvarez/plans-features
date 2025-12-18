package router

import (
	"context"
	"net/http"

	"plans-features/internal/domain/apikeys"
	"plans-features/internal/domain/features"
	"plans-features/internal/domain/plans"
	"plans-features/internal/domain/projects"
	"plans-features/internal/domain/tenantplans"

	"plans-features/internal/utils"

	"github.com/go-chi/chi/v5"
)

// key for context value
type ctxKey string

const CtxProjectID ctxKey = "project_id"

func NewRouter() http.Handler {
	r := chi.NewRouter()

	// -------------------------
	// SINGLETON repositories
	// -------------------------
	projectRepo := projects.NewProjectRepository()
	planRepo := plans.NewPlanRepository()
	featureRepo := features.NewFeatureRepository()
	tenantPlanRepo := tenantplans.NewTenantPlanRepository()
	apiKeyRepo := apikeys.NewAPIKeyRepository()

	// -------------------------
	// Services with dependencies
	// -------------------------

	projectService := projects.NewProjectService(projectRepo)

	planService := plans.NewPlanService(planRepo, projectRepo)

	featureService := features.NewFeatureService(featureRepo, projectRepo)

	tenantPlanService := tenantplans.NewTenantPlanService(
		tenantPlanRepo,
		projectRepo,
		planRepo,
	)

	apiKeyService := apikeys.NewAPIKeyService(apiKeyRepo, projectRepo)

	// -------------------------
	// Handlers
	// -------------------------

	projectHandler := projects.NewProjectHandler(projectService)
	planHandler := plans.NewPlanHandler(planService)
	featureHandler := features.NewFeatureHandler(featureService)
	tenantPlanHandler := tenantplans.NewTenantPlanHandler(tenantPlanService)
	apiKeyHandler := apikeys.NewAPIKeyHandler(apiKeyService)

	// -------------------------
	// Middleware: ApiKeyAuth
	// -------------------------
	apiKeyAuth := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read X-API-Key or Authorization Bearer
			key := r.Header.Get("X-API-Key")
			if key == "" {
				auth := r.Header.Get("Authorization")
				if len(auth) > 7 && auth[:7] == "Bearer " {
					key = auth[7:]
				}
			}
			if key == "" {
				utils.Error(w, http.StatusUnauthorized, "missing api key")
				return
			}
			// validate
			projectID, err := apiKeyService.ValidateKey(r.Context(), key)
			if err != nil {
				utils.Error(w, http.StatusUnauthorized, "invalid api key")
				return
			}
			// set project id in context
			ctx := context.WithValue(r.Context(), CtxProjectID, projectID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	// register middleware
	r.Use(apiKeyAuth)

	// -------------------------
	// Routes
	// -------------------------

	r.Route("/admin", func(r chi.Router) {

		// Projects
		r.Route("/projects", func(r chi.Router) {
			r.Get("/", projectHandler.ListProjects)
			r.Post("/", projectHandler.CreateProject)
			r.Get("/{projectId}", projectHandler.GetProject)
			r.Put("/{projectId}", projectHandler.UpdateProject)

			// API keys for project
			r.Post("/{projectId}/apikeys", apiKeyHandler.CreateKey)
			r.Post("/{projectId}/apikeys/rotate", apiKeyHandler.RotateKey)
			r.Post("/{projectId}/apikeys/revoke", apiKeyHandler.RevokeKey)

			// Plans per project
			r.Route("/{projectId}/plans", func(r chi.Router) {
				r.Get("/", planHandler.ListPlans)
				r.Post("/", planHandler.CreatePlan)
				r.Get("/{planId}", planHandler.GetPlan)
				r.Patch("/{planId}", planHandler.UpdatePlan)
			})

			// Features per project
			r.Route("/{projectId}/features", func(r chi.Router) {
				r.Get("/", featureHandler.ListFeatures)
				r.Post("/", featureHandler.CreateFeature)
				r.Get("/{featureId}", featureHandler.GetFeature)
				r.Patch("/{featureId}", featureHandler.UpdateFeature)
			})
		})

		// Tenant plan assignments
		r.Route("/tenants/{tenantId}/assignments", func(r chi.Router) {
			r.Get("/", tenantPlanHandler.ListAssignments)
			r.Post("/", tenantPlanHandler.CreateAssignment)
			r.Patch("/{assignmentId}", tenantPlanHandler.UpdateAssignment)
		})
	})

	// API routes (require ApiKeyAuth middleware to set project_id in context)
	r.Route("/api", func(r chi.Router) {
		r.Get("/plans", planHandler.ListPlans)
		r.Post("/plans", planHandler.CreatePlan)
		r.Get("/plans/{planId}", planHandler.GetPlan)
		r.Put("/plans/{planId}", planHandler.UpdatePlan)
	})

	// -------------------------
	// Existing admin routes follow
	// -------------------------

	return r
}
