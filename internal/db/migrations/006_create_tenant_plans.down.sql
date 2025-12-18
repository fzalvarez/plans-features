-- 006_create_tenant_plans.down.sql
BEGIN;

-- 1. Drop trigger
DROP TRIGGER IF EXISTS update_tenant_plans_updated_at ON tenant_plans;

-- 2. Drop índices
DROP INDEX IF EXISTS idx_tenant_plans_project_id;
DROP INDEX IF EXISTS idx_tenant_plans_tenant_project;
DROP INDEX IF EXISTS idx_tenant_plans_tenant_project_unique;

-- 3. Drop tabla
DROP TABLE IF EXISTS tenant_plans;

COMMIT;
