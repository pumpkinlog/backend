CREATE TABLE presences (
    user_id INTEGER,
    region_id TEXT ,
    date DATE,
    device_id INTEGER,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (user_id, region_id, date),
    FOREIGN KEY (region_id) REFERENCES regions(id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (device_id) REFERENCES devices(id)
);
