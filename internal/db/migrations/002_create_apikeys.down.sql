-- 002_create_apikeys.down.sql
BEGIN;

DROP INDEX IF EXISTS idx_api_keys_revoked;
DROP INDEX IF EXISTS idx_api_keys_hash;
DROP INDEX IF EXISTS idx_api_keys_project_id;
DROP TABLE IF EXISTS api_keys;

COMMIT;
