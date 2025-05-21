CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    favorite_regions TEXT[],
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);