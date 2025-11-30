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

	repo := projects.NewProjectRepository()
	service := projects.NewProjectService(repo)
	ph := projects.NewProjectHandler(service)

	planRepo := plans.NewPlanRepository()
	planService := plans.NewPlanService(planRepo)
	planHandler := plans.NewPlanHandler(planService)

	featureRepo := features.NewFeatureRepository()
	featureService := features.NewFeatureService(featureRepo)
	featureHandler := features.NewFeatureHandler(featureService)

	tenantRepo := tenantplans.NewTenantPlanRepository()
	tenantService := tenantplans.NewTenantPlanService(tenantRepo)
	tenantHandler := tenantplans.NewTenantPlanHandler(tenantService)

	r.Route("/admin", func(r chi.Router) {
		r.Route("/projects", func(r chi.Router) {
			r.Get("/", ph.ListProjects)
			r.Post("/", ph.CreateProject)
			r.Get("/{projectId}", ph.GetProject)
			r.Patch("/{projectId}", ph.UpdateProject)

			r.Route("/{projectId}/plans", func(r chi.Router) {
				r.Get("/", planHandler.ListPlans)
				r.Post("/", planHandler.CreatePlan)
				r.Get("/{planId}", planHandler.GetPlan)
				r.Patch("/{planId}", planHandler.UpdatePlan)
			})

			r.Route("/{projectId}/features", func(r chi.Router) {
				r.Get("/", featureHandler.ListFeatures)
				r.Post("/", featureHandler.CreateFeature)
				r.Get("/{featureId}", featureHandler.GetFeature)
				r.Patch("/{featureId}", featureHandler.UpdateFeature)
			})
		})

		r.Route("/tenants", func(r chi.Router) {
			r.Route("/{tenantId}/assignments", func(r chi.Router) {
				r.Get("/", tenantHandler.ListAssignments)
				r.Post("/", tenantHandler.CreateAssignment)
				r.Patch("/{assignmentId}", tenantHandler.UpdateAssignment)
			})
		})
	})

	return r
}
