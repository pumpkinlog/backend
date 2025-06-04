CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    favorite_regions TEXT[] NOT NULL DEFAULT '{}',
    want_residency TEXT[] NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);