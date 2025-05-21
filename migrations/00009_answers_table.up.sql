CREATE TABLE answers (
    user_id UUID NOT NULL,
    condition_id UUID NOT NULL,
    value JSONB NOT NULL,
    PRIMARY KEY (user_id, condition_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (condition_id) REFERENCES conditions(id)
); 