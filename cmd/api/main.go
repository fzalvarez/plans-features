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

	"plans-features/internal/config"
	"plans-features/internal/db"
	"plans-features/internal/router"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "plans-features/internal/docs"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// 1. Cargar configuración
	cfg := config.Load()

	// 2. Inicializar PostgreSQL
	dbConn, err := db.New(cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Printf("DB close error: %v", err)
		}
	}()

	if err := runMigrations(cfg.Database.URL); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}

	// 3. Crear router con DB inyectada
	r := router.NewRouter(dbConn) // ← PASAR DB AQUÍ

	// 4. Root router (sin cambios)
	root := chi.NewRouter()
	root.Get("/swagger/*", httpSwagger.WrapHandler)
	root.Mount("/", r)

	fmt.Println("Server running on :8080 with PostgreSQL")
	if err := http.ListenAndServe(":8080", root); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func runMigrations(dbURL string) error {
	m, err := migrate.New(
		"file://internal/db/migrations",
		dbURL,
	)
	if err != nil {
		return fmt.Errorf("new migrate: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("up migrations: %w", err)
	}

	fmt.Println("Migrations applied successfully")
	return nil
}
