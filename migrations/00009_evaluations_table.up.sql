CREATE TABLE evaluations (
    user_id INTEGER NOT NULL,
    region_id TEXT NOT NULL,
    passed boolean NOT NULL,
    details JSONB NOT NULL,
    point_in_time TIMESTAMPTZ NOT NULL,
    evaluated_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (user_id, region_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (region_id) REFERENCES regions(id) ON DELETE CASCADE
); 