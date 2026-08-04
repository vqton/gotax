-- Salary templates and template components
-- salary_components table already exists in 000019

CREATE TABLE IF NOT EXISTS salary_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    name VARCHAR(200) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(company_id, name)
);

CREATE INDEX IF NOT EXISTS idx_salary_templates_company ON salary_templates(company_id);

CREATE TABLE IF NOT EXISTS salary_template_components (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    salary_template_id UUID NOT NULL REFERENCES salary_templates(id) ON DELETE CASCADE,
    component_id UUID NOT NULL REFERENCES salary_components(id) ON DELETE CASCADE,
    default_value NUMERIC(15,2) DEFAULT 0,
    formula TEXT,
    sort_order INT DEFAULT 0,
    UNIQUE(salary_template_id, component_id)
);

CREATE INDEX IF NOT EXISTS idx_stc_template ON salary_template_components(salary_template_id);
CREATE INDEX IF NOT EXISTS idx_stc_component ON salary_template_components(component_id);
