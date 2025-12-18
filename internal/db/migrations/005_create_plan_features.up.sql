-- 005_create_plan_features.up.sql
BEGIN;

-- 1. Función trigger (simple, sin DO $$)
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 2. Tabla junction plan_features
CREATE TABLE plan_features (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    plan_id UUID NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
    feature_id UUID NOT NULL REFERENCES features(id) ON DELETE CASCADE,
    value_json JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- 3. UNIQUE: 1 feature por plan
CREATE UNIQUE INDEX idx_plan_features_plan_feature_unique 
    ON plan_features (plan_id, feature_id);

-- 4. Índices
CREATE INDEX idx_plan_features_plan_id ON plan_features (plan_id);
CREATE INDEX idx_plan_features_project_id ON plan_features (project_id);
CREATE INDEX idx_plan_features_feature_id ON plan_features (feature_id);

-- 5. Trigger
CREATE TRIGGER update_plan_features_updated_at
    BEFORE UPDATE ON plan_features
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMIT;
