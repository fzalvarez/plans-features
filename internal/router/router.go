package router

import (
	"net/http"

	"plans-features/internal/domain/features"
	"plans-features/internal/domain/plans"
	"plans-features/internal/domain/projects"
	"plans-features/internal/domain/tenantplans"

	"github.com/go-chi/chi/v5"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	// -------------------------
	// SINGLETON repositories
	// -------------------------
	projectRepo := projects.NewProjectRepository()
	planRepo := plans.NewPlanRepository()
	featureRepo := features.NewFeatureRepository()
	tenantPlanRepo := tenantplans.NewTenantPlanRepository()

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

	// -------------------------
	// Handlers
	// -------------------------

	projectHandler := projects.NewProjectHandler(projectService)
	planHandler := plans.NewPlanHandler(planService)
	featureHandler := features.NewFeatureHandler(featureService)
	tenantPlanHandler := tenantplans.NewTenantPlanHandler(tenantPlanService)

	// -------------------------
	// Routes
	// -------------------------

	r.Route("/admin", func(r chi.Router) {

		// Projects
		r.Route("/projects", func(r chi.Router) {
			r.Get("/", projectHandler.ListProjects)
			r.Post("/", projectHandler.CreateProject)
			r.Get("/{projectId}", projectHandler.GetProject)
			r.Patch("/{projectId}", projectHandler.UpdateProject)

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

	return r
}
