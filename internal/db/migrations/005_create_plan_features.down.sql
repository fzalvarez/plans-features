-- 005_create_plan_features.down.sql
BEGIN;

-- 1. Drop trigger
DROP TRIGGER IF EXISTS update_plan_features_updated_at ON plan_features;

-- 2. Drop índices
DROP INDEX IF EXISTS idx_plan_features_feature_id;
DROP INDEX IF EXISTS idx_plan_features_project_id;
DROP INDEX IF EXISTS idx_plan_features_plan_id;
DROP INDEX IF EXISTS idx_plan_features_plan_feature_unique;

-- 3. Drop tabla
DROP TABLE IF EXISTS plan_features;

COMMIT;
