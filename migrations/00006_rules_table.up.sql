CREATE TABLE rules (
    id UUID PRIMARY KEY,
    region_id TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    rule_type TEXT NOT NULL,
    period_type TEXT NOT NULL,
    threshold SMALLINT NOT NULL,
    year_start_month SMALLINT NOT NULL,
    year_start_day SMALLINT NOT NULL,
    offset_years SMALLINT NOT NULL,
    years SMALLINT NOT NULL,
    rolling_days SMALLINT NOT NULL,
    rolling_months SMALLINT NOT NULL,
    rolling_years SMALLINT NOT NULL,
    FOREIGN KEY (region_id) REFERENCES regions(id)
);
