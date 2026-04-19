CREATE TABLE IF NOT EXISTS resumes (
    id              BIGSERIAL    PRIMARY KEY,
    slug            VARCHAR(255) UNIQUE,
    user_id         BIGINT,
    region_id       BIGINT,
    district_id     BIGINT,
    mahalla_id      BIGINT,
    adress          TEXT,
    name            VARCHAR(255),
    photo           TEXT,
    title           VARCHAR(500),
    text            TEXT,
    contact         VARCHAR(255),
    price           BIGINT,
    experience_year INT,
    skills          TEXT,
    views_count     BIGINT       DEFAULT 0,
    is_active       BOOLEAN      DEFAULT TRUE,
    created_at      TIMESTAMPTZ  DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_resumes_user_id     ON resumes(user_id);
CREATE INDEX IF NOT EXISTS idx_resumes_region_id   ON resumes(region_id);
CREATE INDEX IF NOT EXISTS idx_resumes_district_id ON resumes(district_id);
CREATE INDEX IF NOT EXISTS idx_resumes_mahalla_id  ON resumes(mahalla_id);
CREATE INDEX IF NOT EXISTS idx_resumes_is_active   ON resumes(is_active);
CREATE INDEX IF NOT EXISTS idx_resumes_deleted_at  ON resumes(deleted_at);
CREATE INDEX IF NOT EXISTS idx_resumes_slug        ON resumes(slug);
CREATE INDEX IF NOT EXISTS idx_resumes_price       ON resumes(price);
