CREATE TABLE IF NOT EXISTS vacancies (
    id          BIGSERIAL    PRIMARY KEY,
    slug        VARCHAR(255) UNIQUE,
    user_id     BIGINT,
    region_id   BIGINT,
    district_id BIGINT,
    mahalla_id  BIGINT,
    adress      TEXT,
    name        VARCHAR(255),
    title       VARCHAR(500),
    text        TEXT,
    contact     VARCHAR(255),
    price       BIGINT,
    views_count BIGINT       DEFAULT 0,
    is_active   BOOLEAN      DEFAULT TRUE,
    created_at  TIMESTAMPTZ  DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_vacancies_user_id     ON vacancies(user_id);
CREATE INDEX IF NOT EXISTS idx_vacancies_region_id   ON vacancies(region_id);
CREATE INDEX IF NOT EXISTS idx_vacancies_district_id ON vacancies(district_id);
CREATE INDEX IF NOT EXISTS idx_vacancies_mahalla_id  ON vacancies(mahalla_id);
CREATE INDEX IF NOT EXISTS idx_vacancies_is_active   ON vacancies(is_active);
CREATE INDEX IF NOT EXISTS idx_vacancies_deleted_at  ON vacancies(deleted_at);
CREATE INDEX IF NOT EXISTS idx_vacancies_slug        ON vacancies(slug);
CREATE INDEX IF NOT EXISTS idx_vacancies_price       ON vacancies(price);
