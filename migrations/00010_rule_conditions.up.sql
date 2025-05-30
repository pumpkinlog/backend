CREATE TABLE rule_conditions (
    rule_id UUID NOT NULL REFERENCES rules(id),
    condition_id UUID NOT NULL REFERENCES conditions(id),
    weight INTEGER NOT NULL DEFAULT 100,
    expected JSONB NOT NULL,
    comparator TEXT NOT NULL,
    PRIMARY KEY (rule_id, condition_id)
);