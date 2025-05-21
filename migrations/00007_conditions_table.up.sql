CREATE TABLE conditions (
    id UUID PRIMARY KEY,
    rule_id UUID NOT NULL,
    prompt TEXT NOT NULL,
    type TEXT NOT NULL,
    comparator TEXT NOT NULL,
    expected JSONB NOT NULL,
    FOREIGN KEY (rule_id) REFERENCES rules(id)
);