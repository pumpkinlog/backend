CREATE TABLE regions (
    id TEXT PRIMARY KEY,
    parent_region_id TEXT,
    region_type TEXT NOT NULL,
    name TEXT NOT NULL,
    continent TEXT NOT NULL,
    year_start_month SMALLINT NOT NULL DEFAULT 1,
    year_start_day SMALLINT NOT NULL DEFAULT 1,
    lat_lng DOUBLE PRECISION[] NOT NULL,
    sources JSONB NOT NULL,
    FOREIGN KEY (parent_region_id) REFERENCES regions(id) ON DELETE CASCADE
);