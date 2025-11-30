// @title Plans & Features Service API
// @version 1.0
// @description API para administrar Projects, Plans, Features y Tenant Plans.
// @BasePath /admin

//go:generate swag init -g cmd/api/main.go -o internal/docs

package main

import (
	"fmt"
	"log"
	"net/http"

	"plans-features/internal/router"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "plans-features/internal/docs" // Swagger docs — IMPORTANTE
)

func main() {
	// your API router
	r := router.NewRouter()

	// root router
	root := chi.NewRouter()

	// swagger UI
	root.Get("/swagger/*", httpSwagger.WrapHandler)

	// mount your internal router
	root.Mount("/", r)

	fmt.Println("server running on :8080")
	if err := http.ListenAndServe(":8080", root); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
