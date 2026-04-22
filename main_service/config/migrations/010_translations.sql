CREATE TABLE IF NOT EXISTS translations (
    id          BIGSERIAL PRIMARY KEY,
    slug        VARCHAR(150) UNIQUE NOT NULL,
    name        JSONB NOT NULL,
    is_active   BOOLEAN DEFAULT TRUE,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW()
);