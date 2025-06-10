CREATE TABLE conditions (
    id VARCHAR(128) PRIMARY KEY,
    region_id TEXT NOT NULL,
    prompt TEXT NOT NULL,
    type TEXT NOT NULL,
    FOREIGN KEY (region_id) REFERENCES regions(id) ON DELETE CASCADE
);