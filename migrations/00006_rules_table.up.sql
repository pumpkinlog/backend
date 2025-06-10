CREATE TABLE rules (
    id VARCHAR(128) PRIMARY KEY,
    region_id TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    node JSONB NOT NULL,
    FOREIGN KEY (region_id) REFERENCES regions(id) ON DELETE CASCADE
);
