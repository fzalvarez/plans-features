-- 004_create_features.down.sql
BEGIN;

-- 1. Drop trigger
DROP TRIGGER IF EXISTS update_features_updated_at ON features;

-- 2. Drop índices
DROP INDEX IF EXISTS idx_features_project_active;
DROP INDEX IF EXISTS idx_features_project_id;
DROP INDEX IF EXISTS idx_features_project_code_unique;

-- 3. Drop tabla (CASCADE borra dependencias)
DROP TABLE IF EXISTS features;

COMMIT;
