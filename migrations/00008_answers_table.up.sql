CREATE TABLE answers (
    user_id INTEGER NOT NULL,
    condition_id VARCHAR(128) NOT NULL,
    region_id TEXT NOT NULL,
    value JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (user_id, condition_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (condition_id) REFERENCES conditions(id),
    FOREIGN KEY (region_id) REFERENCES regions(id) ON DELETE CASCADE
); 