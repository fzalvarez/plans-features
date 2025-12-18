-- 004_create_features.up.sql
BEGIN;

-- 1. Crear función trigger (solo si no existe) - SIN DO $$
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 2. Tabla features
CREATE TABLE features (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    code TEXT NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('flag', 'numeric', 'value')),
    name TEXT NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- 3. UNIQUE index
CREATE UNIQUE INDEX idx_features_project_code_unique ON features (project_id, code) 
WHERE is_active = true;

-- 4. Índices
CREATE INDEX idx_features_project_id ON features (project_id);
CREATE INDEX idx_features_project_active ON features (project_id, is_active);

-- 5. Trigger
CREATE TRIGGER update_features_updated_at
    BEFORE UPDATE ON features
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMIT;
