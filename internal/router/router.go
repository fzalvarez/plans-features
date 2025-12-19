package router

import (
	"net/http"

	"plans-features/internal/db"
	"plans-features/internal/domain/apikeys"
	"plans-features/internal/domain/features"
	"plans-features/internal/domain/planfeatures"
	"plans-features/internal/domain/plans"
	"plans-features/internal/domain/projects"
	"plans-features/internal/domain/tenantplans"

	"github.com/go-chi/chi/v5"
)

// key for context value
type ctxKey string

const CtxProjectID ctxKey = "project_id"

func NewRouter(db *db.DB) http.Handler {
	r := chi.NewRouter()

	// -------------------------
	// SINGLETON repositories
	// -------------------------
	projectRepo := projects.NewProjectRepository(db.SQLDB())
	planRepo := plans.NewPlanRepository(db.SQLDB())
	featureRepo := features.NewFeatureRepository(db.SQLDB())
	tenantPlanRepo := tenantplans.NewTenantPlanRepository(db.SQLDB())
	apiKeyRepo := apikeys.NewAPIKeyRepository(db.SQLDB())
	planFeatureRepo := planfeatures.NewPlanFeatureRepository(db.SQLDB())

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

	planFeatureService := planfeatures.NewPlanFeatureService(planFeatureRepo, planRepo, featureRepo, projectRepo)

	// -------------------------
	// Handlers
	// -------------------------

	projectHandler := projects.NewProjectHandler(projectService)
	planHandler := plans.NewPlanHandler(planService)
	featureHandler := features.NewFeatureHandler(featureService)
	tenantPlanHandler := tenantplans.NewTenantPlanHandler(tenantPlanService)
	apiKeyHandler := apikeys.NewAPIKeyHandler(apiKeyService)
	planFeatureHandler := planfeatures.NewPlanFeatureHandler(planFeatureService)

	// -------------------------
	// Middleware: ApiKeyAuth
	// -------------------------
	/*apiKeyAuth := func(next http.Handler) http.Handler {
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
	r.Use(apiKeyAuth) */

	// -------------------------
	// Routes
	// -------------------------

	// -------------------------
	// ADMIN routes (management)
	// @Summary Admin endpoints
	// @Description Administrative endpoints to manage Projects, Plans, Features, Tenant assignments and API keys
	// @Tags projects, plans, features, apikeys, tenantplans
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
	// -------------------------
	// @Summary Public API endpoints (scoped by API key)
	// @Description API endpoints accessible with X-API-Key header. These endpoints operate within the project context derived from the API key.
	// @Tags plans, features, planfeatures, tenantplans
	// @Param X-API-Key header string true "API Key"
	// -------------------------
	r.Route("/api", func(r chi.Router) {
		r.Get("/plans", planHandler.ListPlans)
		r.Post("/plans", planHandler.CreatePlan)
		r.Get("/plans/{planId}", planHandler.GetPlan)
		r.Put("/plans/{planId}", planHandler.UpdatePlan)

		// Features API scoped by API key
		r.Get("/features", featureHandler.ListFeatures)
		r.Post("/features", featureHandler.CreateFeature)
		r.Get("/features/{featureId}", featureHandler.GetFeature)
		r.Put("/features/{featureId}", featureHandler.UpdateFeature)

		// PlanFeatures routes
		r.Route("/plans/{planId}/features", func(r chi.Router) {
			r.Get("/", planFeatureHandler.List)
			r.Post("/", planFeatureHandler.Assign)
		})

		// TenantPlans API: get effective plan and assign plan (scoped by API key)
		r.Get("/tenants/{tenantId}/plan", tenantPlanHandler.GetTenantPlan)
		r.Post("/tenants/{tenantId}/plan", tenantPlanHandler.AssignTenantPlan)
	})

	// -------------------------
	// Existing admin routes follow
	// -------------------------

	return r
}
