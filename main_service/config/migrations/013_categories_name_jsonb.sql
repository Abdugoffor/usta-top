ALTER TABLE categories
  ALTER COLUMN name TYPE JSONB
  USING CASE
    WHEN name IS NOT NULL AND name != '' THEN jsonb_build_object('default', name)
    ELSE '{}'::jsonb
  END;
