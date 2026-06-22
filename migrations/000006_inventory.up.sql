-- Material catalog
CREATE TABLE materials (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(200) NOT NULL,
    unit        VARCHAR(50)  NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_materials_name ON materials(name);

-- Head-office and project inventory (one row per material per location)
CREATE TABLE inventory_items (
    id            UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    material_id   UUID          NOT NULL REFERENCES materials(id) ON DELETE CASCADE,
    location_type VARCHAR(20)   NOT NULL CHECK (location_type IN ('head_office','project')),
    project_id    UUID          REFERENCES projects(id) ON DELETE CASCADE,
    quantity      NUMERIC(12,3) NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    unit_cost     NUMERIC(12,2),
    created_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);
-- NULL project_id for head_office, partial unique indexes handle the constraint
CREATE UNIQUE INDEX idx_inventory_ho      ON inventory_items(material_id)             WHERE location_type = 'head_office';
CREATE UNIQUE INDEX idx_inventory_project ON inventory_items(material_id, project_id) WHERE location_type = 'project';

-- Requisition status enum
CREATE TYPE requisition_status AS ENUM ('pending','approved','rejected','fulfilled');

-- Requisitions (project requests materials from head office)
CREATE TABLE requisitions (
    id           UUID               PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id   UUID               NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    requested_by UUID               NOT NULL REFERENCES users(id),
    status       requisition_status NOT NULL DEFAULT 'pending',
    notes        TEXT,
    created_at   TIMESTAMPTZ        NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ        NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_requisitions_project ON requisitions(project_id);
CREATE INDEX idx_requisitions_status  ON requisitions(status);

CREATE TABLE requisition_items (
    id                 UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    requisition_id     UUID          NOT NULL REFERENCES requisitions(id) ON DELETE CASCADE,
    material_id        UUID          NOT NULL REFERENCES materials(id),
    quantity_requested NUMERIC(12,3) NOT NULL CHECK (quantity_requested > 0),
    quantity_approved  NUMERIC(12,3)           CHECK (quantity_approved >= 0),
    created_at         TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    UNIQUE(requisition_id, material_id)
);

-- Full audit trail of every inventory movement
CREATE TABLE inventory_transactions (
    id               UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    material_id      UUID          NOT NULL REFERENCES materials(id),
    location_type    VARCHAR(20)   NOT NULL,
    project_id       UUID          REFERENCES projects(id),
    quantity_change  NUMERIC(12,3) NOT NULL,
    transaction_type VARCHAR(30)   NOT NULL,
    reference_id     UUID,
    reference_type   VARCHAR(50),
    notes            TEXT,
    created_by       UUID          REFERENCES users(id),
    created_at       TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_inv_tx_material ON inventory_transactions(material_id);
CREATE INDEX idx_inv_tx_project  ON inventory_transactions(project_id)   WHERE project_id IS NOT NULL;
CREATE INDEX idx_inv_tx_created  ON inventory_transactions(created_at DESC);

-- Materials consumed in daily logs (auto-deducts from project inventory)
CREATE TABLE log_materials (
    id            UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    log_id        UUID          NOT NULL REFERENCES daily_logs(id) ON DELETE CASCADE,
    material_id   UUID          NOT NULL REFERENCES materials(id),
    quantity_used NUMERIC(12,3) NOT NULL CHECK (quantity_used > 0),
    created_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    UNIQUE(log_id, material_id)
);
CREATE INDEX idx_log_materials_log ON log_materials(log_id);
