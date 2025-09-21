-- dev/setup.sql
\set ON_ERROR_STOP on

-- 1) Create the database if it doesn't exist (psql-specific \gexec trick)
SELECT 'CREATE DATABASE service_catalog WITH TEMPLATE=template0 ENCODING ''UTF8'''
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'service_catalog')\gexec

-- 2) Connect to the new (or existing) database
\connect service_catalog

-- 3) Enable UUID generation
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- 4) Table + indexes
CREATE TABLE IF NOT EXISTS services (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name        TEXT        NOT NULL,
  description TEXT,
  version     TEXT        NOT NULL,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now()
);

-- Make seeds idempotent (avoid dupes across repeated runs)
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_constraint
    WHERE conname = 'uq_services_name'
  ) THEN
    ALTER TABLE services ADD CONSTRAINT uq_services_name UNIQUE (name);
  END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_services_name_lower ON services (LOWER(name));

-- 5) Dummy data (idempotent upsert on UNIQUE(name))
INSERT INTO services (name, description, version)
VALUES
  ('Auth',     'Authentication service',       'v1.3.2'),
  ('Payments', 'Payment orchestration',        'v2.1.0'),
  ('Catalog',  'Service catalog metadata',     'v0.9.5')
ON CONFLICT (name) DO UPDATE
SET description = EXCLUDED.description,
    version     = EXCLUDED.version,
    updated_at  = now();


-- one-time setup
ALTER TABLE services
  ADD COLUMN search_doc tsvector GENERATED ALWAYS AS (
    setweight(to_tsvector('english', coalesce(name,'')), 'A') ||
    setweight(to_tsvector('english', coalesce(description,'')), 'B')
  ) STORED;
CREATE INDEX IF NOT EXISTS idx_services_fts ON services USING GIN (search_doc);

