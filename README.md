# plans-features

Servicio backend escrito en Go para gestionar planes, características, proyectos y la asignación de planes por tenant. Proporciona la estructura de API, acceso a base de datos PostgreSQL y capas de dominio (handlers, services, repositories) para las entidades principales.

## Qué hace el proyecto

- Expone una API (entrada en `cmd/api/main.go`) para CRUD y operaciones relacionadas con:
  - `plans` (planes)
  - `features` (características)
  - `projects` (proyectos)
  - `tenantplans` (asignación de planes a tenants)
- Conexión a PostgreSQL con migraciones en `internal/db/migrations`.
- Estructura de carpetas por dominio en `internal/domain` con `handler.go`, `service.go`, `repository.go` y `model.go`.
- Utilidades comunes en `internal/utils` y logging en `pkg/logger`.

## Qué ya está implementado

- Estructura de proyecto y módulos principales para cada dominio (`plans`, `features`, `projects`, `tenantplans`).
- Handlers HTTP, servicios y repositorios esbozados para las entidades principales.
- Conexión a PostgreSQL (`internal/db/postgres.go`) y archivo `sqlc.yaml` para generación de consultas si se utiliza `sqlc`.
- Sistema básico de configuración en `internal/config/config.go`.
- Módulo de logging en `pkg/logger` y utilidades para respuestas y validaciones en `internal/utils`.
- Router de la API en `internal/router/router.go`.

## Qué falta / próximos pasos

- Implementar la lógica de negocio completa en los servicios y repositorios (si hay métodos aún por desarrollar).
- Pruebas unitarias e integración para handlers, servicios y repositorios.
- Autenticación y autorización (JWT, OAuth u otro método) si la API debe ser protegida.
- Documentación de la API (OpenAPI/Swagger) y ejemplos de uso.
- Integración continua (CI) y despliegue (CD) automatizados.
- Revisión y aplicación de migraciones pendientes en `internal/db/migrations`.
- Validaciones más completas y manejo de errores robusto según necesidades de negocio.

## Cómo ejecutar (rápido)

1. Configurar las variables de entorno o archivo de configuración según `internal/config/config.go`.
2. Preparar una base de datos PostgreSQL y ejecutar las migraciones en `internal/db/migrations`.
3. Construir y ejecutar la API:

```powershell
go build -o bin/api ./cmd/api
./bin/api
```

O ejecutar directamente:

```powershell
go run ./cmd/api
```

Ajustar los comandos según el entorno y la forma de despliegue deseada.

---

Si quieres, puedo añadir más detalles (endpoints disponibles, ejemplo de .env, pasos para ejecutar migraciones, o generar Swagger) o actualizar el README con información específica que prefieras.