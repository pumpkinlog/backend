CREATE TABLE rule_conditions (
    rule_id INTEGER,
    condition_id INTEGER,
    region_id TEXT NOT NULL,
    comparator TEXT NOT NULL,
    expected JSONB NOT NULL,
    PRIMARY KEY (rule_id, condition_id),
    FOREIGN KEY (rule_id) REFERENCES rules(id) ON DELETE CASCADE,
    FOREIGN KEY (condition_id) REFERENCES conditions(id) ON DELETE CASCADE,
    FOREIGN KEY (region_id) REFERENCES regions(id) ON DELETE CASCADE
);