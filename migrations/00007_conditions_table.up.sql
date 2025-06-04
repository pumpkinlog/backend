CREATE TABLE conditions (
    id SERIAL PRIMARY KEY,
    region_id TEXT NOT NULL,
    prompt TEXT NOT NULL,
    type TEXT NOT NULL,
    FOREIGN KEY (region_id) REFERENCES regions(id) ON DELETE CASCADE
);