CREATE TABLE presences (
    user_id UUID NOT NULL,
    region_id TEXT NOT NULL,
    date DATE NOT NULL,
    device_id UUID,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    PRIMARY KEY (user_id, region_id, date),
    FOREIGN KEY (region_id) REFERENCES regions(id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (device_id) REFERENCES devices(id)
);
