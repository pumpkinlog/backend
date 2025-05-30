CREATE TABLE evaluations (
    user_id UUID NOT NULL,
    region_id TEXT NOT NULL,
    passed boolean NOT NULL,
    region JSONB NOT NULL,
    evaluations JSONB NOT NULL,
    evaluated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    PRIMARY KEY (user_id, region_id)
); 