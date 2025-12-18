-- 003_create_plans.down.sql
BEGIN;

DROP TRIGGER IF EXISTS update_plans_updated_at ON plans;
DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE;
DROP INDEX IF EXISTS idx_plans_project_default;
DROP INDEX IF EXISTS idx_plans_project_id;
DROP INDEX IF EXISTS idx_plans_project_code_unique;
DROP TABLE IF EXISTS plans;

COMMIT;
