-- 006_create_tenant_plans.up.sql
BEGIN;

-- 1. Función trigger SIMPLE (SIN DO $$)
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 2. Tabla tenant_plans
CREATE TABLE tenant_plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id TEXT NOT NULL,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    plan_id UUID NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- 3. UNIQUE constraint
CREATE UNIQUE INDEX idx_tenant_plans_tenant_project_unique 
    ON tenant_plans (tenant_id, project_id);

-- 4. Índices
CREATE INDEX idx_tenant_plans_tenant_project ON tenant_plans (tenant_id, project_id);
CREATE INDEX idx_tenant_plans_project_id ON tenant_plans (project_id);

-- 5. Trigger
CREATE TRIGGER update_tenant_plans_updated_at
    BEFORE UPDATE ON tenant_plans
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMIT;
