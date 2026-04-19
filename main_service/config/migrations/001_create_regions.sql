CREATE TABLE IF NOT EXISTS regions (
    id          BIGSERIAL    PRIMARY KEY,
    parent_id   BIGINT       NOT NULL DEFAULT 0,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    is_active   BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_regions_parent_id  ON regions(parent_id);
CREATE INDEX IF NOT EXISTS idx_regions_name       ON regions(name);
CREATE INDEX IF NOT EXISTS idx_regions_is_active  ON regions(is_active);
CREATE INDEX IF NOT EXISTS idx_regions_deleted_at ON regions(deleted_at);
