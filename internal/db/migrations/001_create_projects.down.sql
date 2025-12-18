-- 001_create_projects.down.sql

BEGIN;

DROP TRIGGER IF EXISTS update_projects_updated_at ON projects;
DROP FUNCTION IF EXISTS update_projects_updated_at_column();
DROP INDEX IF EXISTS idx_projects_active;
DROP INDEX IF EXISTS idx_projects_code;
DROP TABLE IF EXISTS projects;

COMMIT;
