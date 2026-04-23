-- pg_trgm: ILIKE '%...%' so'rovlari uchun GIN trigram indekslari
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Resumes: qidiruv maydonlari
CREATE INDEX IF NOT EXISTS idx_resumes_name_trgm    ON resumes USING GIN (name    gin_trgm_ops) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_resumes_title_trgm   ON resumes USING GIN (title   gin_trgm_ops) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_resumes_skills_trgm  ON resumes USING GIN (skills  gin_trgm_ops) WHERE deleted_at IS NULL;

-- Vacancies: qidiruv maydonlari
CREATE INDEX IF NOT EXISTS idx_vacancies_name_trgm  ON vacancies USING GIN (name  gin_trgm_ops) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_vacancies_title_trgm ON vacancies USING GIN (title gin_trgm_ops) WHERE deleted_at IS NULL;

-- Resumes: filter maydonlari
CREATE INDEX IF NOT EXISTS idx_resumes_active_region   ON resumes (region_id,   is_active) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_resumes_active_district ON resumes (district_id, is_active) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_resumes_user_id         ON resumes (user_id)                WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_resumes_price           ON resumes (price)                  WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_resumes_experience      ON resumes (experience_year)        WHERE deleted_at IS NULL;

-- Vacancies: filter maydonlari
CREATE INDEX IF NOT EXISTS idx_vacancies_active_region   ON vacancies (region_id,   is_active) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_vacancies_active_district ON vacancies (district_id, is_active) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_vacancies_user_id         ON vacancies (user_id)                WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_vacancies_price           ON vacancies (price)                  WHERE deleted_at IS NULL;

-- Category join tables
CREATE INDEX IF NOT EXISTS idx_cat_resume_resume_id   ON category_resume  (resume_id);
CREATE INDEX IF NOT EXISTS idx_cat_resume_cat_id      ON category_resume  (categorya_id);
CREATE INDEX IF NOT EXISTS idx_cat_vacancy_vacancy_id ON category_vacancy (vacancy_id);
CREATE INDEX IF NOT EXISTS idx_cat_vacancy_cat_id     ON category_vacancy (categorya_id);
