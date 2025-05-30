CREATE TABLE region_conditions (
    region_id TEXT NOT NULL REFERENCES regions(id),
    condition_id UUID NOT NULL REFERENCES conditions(id),
    PRIMARY KEY (region_id, condition_id)
);