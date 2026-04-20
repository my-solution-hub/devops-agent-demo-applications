-- Smart Home Cat Demo — PostgreSQL initialization script
-- Creates all tables required by the application services.

-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Cat Profiles
CREATE TABLE IF NOT EXISTS cat_profiles (
    cat_id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id             VARCHAR(255) NOT NULL,
    name                 VARCHAR(255) NOT NULL,
    breed                VARCHAR(255),
    age_months           INTEGER,
    weight_kg            NUMERIC(5,2) NOT NULL,
    dietary_restrictions TEXT[],
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_cat_profiles_owner ON cat_profiles(owner_id);

-- Devices
CREATE TABLE IF NOT EXISTS devices (
    device_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_type VARCHAR(50)  NOT NULL,
    name        VARCHAR(255) NOT NULL,
    status      VARCHAR(20)  NOT NULL DEFAULT 'offline',
    config      JSONB        DEFAULT '{}',
    last_seen   TIMESTAMPTZ,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Device Assignments
CREATE TABLE IF NOT EXISTS device_assignments (
    assignment_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cat_id        UUID        NOT NULL REFERENCES cat_profiles(cat_id),
    device_id     UUID        NOT NULL REFERENCES devices(device_id),
    device_type   VARCHAR(50) NOT NULL,
    assigned_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(cat_id, device_type)
);

-- Feeding Schedules
CREATE TABLE IF NOT EXISTS feeding_schedules (
    schedule_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cat_id               UUID         NOT NULL REFERENCES cat_profiles(cat_id),
    device_id            UUID         NOT NULL REFERENCES devices(device_id),
    meal_times           TEXT[]       NOT NULL,
    portion_grams        NUMERIC(6,1) NOT NULL,
    max_daily_grams      NUMERIC(6,1),
    min_interval_minutes INTEGER,
    enabled              BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at           TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_feeding_schedules_cat ON feeding_schedules(cat_id);
