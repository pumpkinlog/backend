CREATE TABLE rules_conditions (
    rule_id UUID NOT NULL,
    condition_id UUID NOT NULL,
    PRIMARY KEY (rule_id, condition_id),
    FOREIGN KEY (rule_id) REFERENCES rules(id),
    FOREIGN KEY (condition_id) REFERENCES conditions(id)
); 