CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL,
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    date TIMESTAMP NOT NULL,
    FOREIGN KEY (device_id) REFERENCES devices(id)
);
